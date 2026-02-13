package order

import (
	"context"
	"encoding/json"

	"go-mail/common/xerr"
	"go-mail/services/inventory/rpc/inventoryrpc"
	"go-mail/services/order/api/internal/svc"
	"go-mail/services/order/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CancelOrder 取消订单（含库存回滚）
func (l *CancelOrderLogic) CancelOrder(req *types.CancelOrderReq) (resp *types.CancelOrderResp, err error) {
	userId, err := l.getUserId()
	if err != nil {
		return nil, err
	}

	order, err := l.svcCtx.OrderMasterModel.FindOne(l.ctx, req.OrderId)
	if err != nil {
		return nil, xerr.NewCodeErrorMsg(xerr.NotFound, "订单不存在")
	}

	// 校验订单归属
	if order.UserId != userId {
		return nil, xerr.NewCodeError(xerr.Unauthorized)
	}

	// 仅待支付状态可取消
	if order.Status != 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.BadRequest, "当前状态无法取消")
	}

	// 更新状态为已取消 (4)
	err = l.svcCtx.OrderMasterModel.UpdateStatus(l.ctx, req.OrderId, 4)
	if err != nil {
		logx.Errorf("取消订单失败: orderId=%s, err=%v", req.OrderId, err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	// Inventory RPC → 回滚库存
	items, err := l.svcCtx.OrderItemModel.FindByOrderId(l.ctx, req.OrderId)
	if err == nil {
		for _, item := range items {
			_, err := l.svcCtx.InventoryRpc.RollbackStock(l.ctx, &inventoryrpc.RollbackStockReq{
				SkuId:    item.SkuId,
				Quantity: int32(item.Quantity),
			})
			if err != nil {
				logx.Errorf("库存回滚失败: orderId=%s, skuId=%d, err=%v", req.OrderId, item.SkuId, err)
			}
		}
	}

	logx.Infof("订单已取消: orderId=%s, userId=%d, reason=%s", req.OrderId, userId, req.Reason)
	return &types.CancelOrderResp{Msg: "取消成功"}, nil
}

func (l *CancelOrderLogic) getUserId() (int64, error) {
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
