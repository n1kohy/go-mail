package logic

import (
	"context"

	"go-mail/services/order/rpc/internal/svc"
	"go-mail/services/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOrder 查询订单（供 Payment/Logistics 调用）
func (l *GetOrderLogic) GetOrder(in *order.GetOrderReq) (*order.GetOrderResp, error) {
	o, err := l.svcCtx.OrderMasterModel.FindOne(l.ctx, in.OrderId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "订单不存在")
	}

	return &order.GetOrderResp{
		OrderId:   o.Id,
		UserId:    o.UserId,
		PayAmount: o.PayAmount,
		Status:    int32(o.Status),
	}, nil
}
