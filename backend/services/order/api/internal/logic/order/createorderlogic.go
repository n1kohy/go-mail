package order

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-mail/common/xerr"
	"go-mail/services/cart/rpc/cartrpc"
	"go-mail/services/inventory/rpc/inventoryrpc"
	"go-mail/services/order/api/internal/svc"
	"go-mail/services/order/api/internal/types"
	"go-mail/services/order/model"
	"go-mail/services/promotion/rpc/promotionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateOrder 创建订单 — 完整跨服务 RPC 调用链
func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (resp *types.CreateOrderResp, err error) {
	userId, err := l.getUserId()
	if err != nil {
		return nil, err
	}

	// ========== 1. Cart RPC → 获取购物车选中商品 ==========
	cartResp, err := l.svcCtx.CartRpc.GetCartItems(l.ctx, &cartrpc.GetCartItemsReq{
		UserId: userId,
	})
	if err != nil {
		logx.Errorf("获取购物车失败: userId=%d, err=%v", userId, err)
		return nil, xerr.NewCodeErrorMsg(xerr.ServerError, "获取购物车失败")
	}
	if len(cartResp.Items) == 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.BadRequest, "购物车为空")
	}

	// ========== 2. 计算总金额（使用购物车商品价格） ==========
	var totalAmount float64
	var skuIds []int64
	for _, item := range cartResp.Items {
		totalAmount += item.Price * float64(item.Quantity)
		skuIds = append(skuIds, item.SkuId)
	}

	// ========== 3. Promotion RPC → 计算优惠 ==========
	var discountAmount float64
	payAmount := totalAmount
	if req.CouponId > 0 {
		promoResp, err := l.svcCtx.PromotionRpc.CalculateDiscount(l.ctx, &promotionrpc.CalculateDiscountReq{
			CouponId:    req.CouponId,
			TotalAmount: totalAmount,
		})
		if err != nil {
			logx.Errorf("计算优惠失败: couponId=%d, err=%v", req.CouponId, err)
		} else if promoResp.Valid {
			discountAmount = promoResp.DiscountAmount
			payAmount = promoResp.PayAmount
		}
	}

	// ========== 4. Inventory RPC → 扣减库存 ==========
	for _, item := range cartResp.Items {
		_, err := l.svcCtx.InventoryRpc.DeductStock(l.ctx, &inventoryrpc.DeductStockReq{
			SkuId:    item.SkuId,
			Quantity: int32(item.Quantity),
		})
		if err != nil {
			logx.Errorf("扣减库存失败: skuId=%d, quantity=%d, err=%v", item.SkuId, item.Quantity, err)
			return nil, xerr.NewCodeErrorMsg(xerr.BadRequest, fmt.Sprintf("商品 %s 库存不足", item.ProductName))
		}
	}

	// ========== 5. 写入 order_master + order_item ==========
	orderId := fmt.Sprintf("SN%s%04d", time.Now().Format("20060102150405"), userId%10000)
	expireTime := time.Now().Add(30 * time.Minute)

	_, err = l.svcCtx.OrderMasterModel.Insert(l.ctx, &model.OrderMaster{
		Id:              orderId,
		UserId:          userId,
		TotalAmount:     totalAmount,
		DiscountAmount:  discountAmount,
		FreightAmount:   0,
		PayAmount:       payAmount,
		CouponId:        req.CouponId,
		Status:          0, // 待支付
		AddressSnapshot: fmt.Sprintf(`{"address_id":%d}`, req.AddressId),
		ExpireTime:      expireTime,
	})
	if err != nil {
		logx.Errorf("创建订单失败: userId=%d, err=%v", userId, err)
		// 回滚库存
		l.rollbackStock(cartResp.Items)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	// 创建订单明细
	for _, item := range cartResp.Items {
		_, err = l.svcCtx.OrderItemModel.Insert(l.ctx, &model.OrderItem{
			OrderId:     orderId,
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			SkuId:       item.SkuId,
			SkuSpecs:    item.SkuSpecs,
			Price:       item.Price,
			Quantity:    int(item.Quantity),
		})
		if err != nil {
			logx.Errorf("创建订单明细失败: orderId=%s, skuId=%d, err=%v", orderId, item.SkuId, err)
		}
	}

	// ========== 6. 核销优惠券 ==========
	if req.CouponId > 0 && discountAmount > 0 {
		_, err = l.svcCtx.PromotionRpc.UseCoupon(l.ctx, &promotionrpc.UseCouponReq{
			UserId:      userId,
			CouponId:    req.CouponId,
			UsedOrderId: orderId,
		})
		if err != nil {
			logx.Errorf("核销优惠券失败: couponId=%d, err=%v", req.CouponId, err)
		}
	}

	// ========== 7. Cart RPC → 清理购物车 ==========
	_, err = l.svcCtx.CartRpc.ClearCartItems(l.ctx, &cartrpc.ClearCartItemsReq{
		UserId: userId,
		SkuIds: skuIds,
	})
	if err != nil {
		logx.Errorf("清理购物车失败: userId=%d, err=%v", userId, err)
	}

	logx.Infof("订单创建成功: orderId=%s, userId=%d, payAmount=%.2f, items=%d", orderId, userId, payAmount, len(cartResp.Items))

	return &types.CreateOrderResp{
		OrderId:    orderId,
		PayAmount:  payAmount,
		ExpireTime: expireTime.Format("2006-01-02T15:04:05"),
	}, nil
}

// rollbackStock 库存回滚（创建订单失败时）
func (l *CreateOrderLogic) rollbackStock(items []*cartrpc.CartItemInfo) {
	for _, item := range items {
		_, err := l.svcCtx.InventoryRpc.RollbackStock(l.ctx, &inventoryrpc.RollbackStockReq{
			SkuId:    item.SkuId,
			Quantity: int32(item.Quantity),
		})
		if err != nil {
			logx.Errorf("库存回滚失败: skuId=%d, quantity=%d, err=%v", item.SkuId, item.Quantity, err)
		}
	}
}

func (l *CreateOrderLogic) getUserId() (int64, error) {
	uid := l.ctx.Value("userId")
	if uid == nil {
		return 0, xerr.NewCodeError(xerr.Unauthorized)
	}
	userId, err := uid.(json.Number).Int64()
	if err != nil {
		return 0, xerr.NewCodeError(xerr.Unauthorized)
	}
	return userId, nil
}
