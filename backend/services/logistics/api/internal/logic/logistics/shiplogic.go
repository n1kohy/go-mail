package logistics

import (
	"context"
	"encoding/json"

	"go-mail/common/xerr"
	"go-mail/services/logistics/api/internal/svc"
	"go-mail/services/logistics/api/internal/types"
	"go-mail/services/logistics/model"
	"go-mail/services/order/rpc/orderrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShipLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShipLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShipLogic {
	return &ShipLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Ship 发货
func (l *ShipLogic) Ship(req *types.ShipReq) (resp *types.ShipResp, err error) {
	_, err = l.getUserId()
	if err != nil {
		return nil, err
	}

	// 1. 调用 Order RPC → 验证订单状态（应为已支付=1）
	orderResp, err := l.svcCtx.OrderRpc.GetOrder(l.ctx, &orderrpc.GetOrderReq{
		OrderId: req.OrderId,
	})
	if err != nil {
		return nil, xerr.NewCodeErrorMsg(xerr.NotFound, "订单不存在")
	}
	if orderResp.Status != 1 {
		return nil, xerr.NewCodeErrorMsg(xerr.BadRequest, "订单状态无法发货")
	}

	// 2. 创建物流记录
	_, err = l.svcCtx.ShippingOrderModel.Insert(l.ctx, &model.ShippingOrder{
		OrderId:    req.OrderId,
		TrackingNo: req.TrackingNo,
		Carrier:    req.Carrier,
		Status:     1, // 已发货
	})
	if err != nil {
		logx.Errorf("创建物流记录失败: orderId=%s, err=%v", req.OrderId, err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	// 3. 调用 Order RPC → 更新订单状态为已发货(2)
	_, err = l.svcCtx.OrderRpc.UpdateOrderStatus(l.ctx, &orderrpc.UpdateOrderStatusReq{
		OrderId: req.OrderId,
		Status:  2,
	})
	if err != nil {
		logx.Errorf("更新订单状态为已发货失败: orderId=%s, err=%v", req.OrderId, err)
	}

	logx.Infof("发货成功: orderId=%s, carrier=%s, trackingNo=%s", req.OrderId, req.Carrier, req.TrackingNo)
	return &types.ShipResp{Msg: "发货成功"}, nil
}

func (l *ShipLogic) getUserId() (int64, error) {
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
