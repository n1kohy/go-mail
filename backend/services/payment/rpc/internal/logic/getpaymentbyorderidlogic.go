package logic

import (
	"context"

	"go-mail/services/payment/rpc/internal/svc"
	"go-mail/services/payment/rpc/payment"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetPaymentByOrderIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaymentByOrderIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentByOrderIdLogic {
	return &GetPaymentByOrderIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetPaymentByOrderId 按订单号查询支付记录
func (l *GetPaymentByOrderIdLogic) GetPaymentByOrderId(in *payment.GetPaymentByOrderIdReq) (*payment.GetPaymentByOrderIdResp, error) {
	flow, err := l.svcCtx.PaymentFlowModel.FindByOrderId(l.ctx, in.OrderId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "支付记录不存在")
	}

	return &payment.GetPaymentByOrderIdResp{
		Id:      flow.Id,
		OrderId: flow.OrderId,
		TradeNo: flow.TradeNo,
		Amount:  flow.Amount,
		Channel: int32(flow.Channel),
		Status:  int32(flow.Status),
	}, nil
}
