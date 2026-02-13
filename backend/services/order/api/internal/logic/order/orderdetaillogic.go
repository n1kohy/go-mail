package order

import (
	"context"
	"encoding/json"

	"go-mail/common/xerr"
	"go-mail/services/order/api/internal/svc"
	"go-mail/services/order/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailLogic {
	return &OrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// OrderDetail 订单详情
func (l *OrderDetailLogic) OrderDetail(orderId string) (resp *types.OrderDetailResp, err error) {
	userId, err := l.getUserId()
	if err != nil {
		return nil, err
	}

	order, err := l.svcCtx.OrderMasterModel.FindOne(l.ctx, orderId)
	if err != nil {
		return nil, xerr.NewCodeErrorMsg(xerr.NotFound, "订单不存在")
	}

	// 校验订单归属
	if order.UserId != userId {
		return nil, xerr.NewCodeError(xerr.Unauthorized)
	}

	// 查询订单明细
	items, err := l.svcCtx.OrderItemModel.FindByOrderId(l.ctx, orderId)
	if err != nil {
		logx.Errorf("查询订单明细失败: orderId=%s, err=%v", orderId, err)
	}

	var orderItems []types.OrderItemInfo
	for _, item := range items {
		orderItems = append(orderItems, types.OrderItemInfo{
			SkuId:       item.SkuId,
			ProductName: item.ProductName,
			Specs:       item.SkuSpecs,
			Price:       item.Price,
			Quantity:    item.Quantity,
		})
	}

	if orderItems == nil {
		orderItems = []types.OrderItemInfo{}
	}

	return &types.OrderDetailResp{
		OrderId:        order.Id,
		Status:         int(order.Status),
		TotalAmount:    order.TotalAmount,
		DiscountAmount: order.DiscountAmount,
		FreightAmount:  order.FreightAmount,
		PayAmount:      order.PayAmount,
		Items:          orderItems,
		CreateTime:     order.CreateTime.Format("2006-01-02 15:04:05"),
	}, nil
}

func (l *OrderDetailLogic) getUserId() (int64, error) {
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
