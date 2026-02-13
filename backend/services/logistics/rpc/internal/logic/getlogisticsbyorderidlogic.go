package logic

import (
	"context"

	"go-mail/services/logistics/rpc/internal/svc"
	"go-mail/services/logistics/rpc/logistics"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetLogisticsByOrderIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLogisticsByOrderIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogisticsByOrderIdLogic {
	return &GetLogisticsByOrderIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetLogisticsByOrderId 按订单号查询物流
func (l *GetLogisticsByOrderIdLogic) GetLogisticsByOrderId(in *logistics.GetLogisticsByOrderIdReq) (*logistics.GetLogisticsByOrderIdResp, error) {
	shipping, err := l.svcCtx.ShippingOrderModel.FindByOrderId(l.ctx, in.OrderId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "物流信息不存在")
	}

	return &logistics.GetLogisticsByOrderIdResp{
		Id:         shipping.Id,
		OrderId:    shipping.OrderId,
		TrackingNo: shipping.TrackingNo,
		Carrier:    shipping.Carrier,
		Status:     int32(shipping.Status),
	}, nil
}
