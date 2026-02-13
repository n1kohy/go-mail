package logistics

import (
	"context"

	"go-mail/common/xerr"
	"go-mail/services/logistics/api/internal/svc"
	"go-mail/services/logistics/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogisticsDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogisticsDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogisticsDetailLogic {
	return &LogisticsDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// LogisticsDetail 物流详情
func (l *LogisticsDetailLogic) LogisticsDetail(orderId string) (resp *types.LogisticsDetailResp, err error) {
	shipping, err := l.svcCtx.ShippingOrderModel.FindByOrderId(l.ctx, orderId)
	if err != nil {
		return nil, xerr.NewCodeErrorMsg(xerr.NotFound, "物流信息不存在")
	}

	return &types.LogisticsDetailResp{
		OrderId:    shipping.OrderId,
		TrackingNo: shipping.TrackingNo,
		Carrier:    shipping.Carrier,
		Status:     int(shipping.Status),
		CreateTime: shipping.CreateTime.Format("2006-01-02 15:04:05"),
	}, nil
}
