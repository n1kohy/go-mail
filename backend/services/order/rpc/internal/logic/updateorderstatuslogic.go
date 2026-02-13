package logic

import (
	"context"

	"go-mail/services/order/rpc/internal/svc"
	"go-mail/services/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateOrderStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderStatusLogic {
	return &UpdateOrderStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateOrderStatus 更新订单状态（供 Payment 回调使用）
func (l *UpdateOrderStatusLogic) UpdateOrderStatus(in *order.UpdateOrderStatusReq) (*order.UpdateOrderStatusResp, error) {
	err := l.svcCtx.OrderMasterModel.UpdateStatus(l.ctx, in.OrderId, int(in.Status))
	if err != nil {
		logx.Errorf("更新订单状态失败: orderId=%s, status=%d, err=%v", in.OrderId, in.Status, err)
		return nil, status.Error(codes.Internal, "更新订单状态失败")
	}

	logx.Infof("订单状态更新: orderId=%s → status=%d", in.OrderId, in.Status)
	return &order.UpdateOrderStatusResp{Success: true}, nil
}
