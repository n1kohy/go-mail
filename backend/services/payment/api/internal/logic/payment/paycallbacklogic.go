package payment

import (
	"context"

	"go-mail/common/xerr"
	"go-mail/services/order/rpc/orderrpc"
	"go-mail/services/payment/api/internal/svc"
	"go-mail/services/payment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayCallbackLogic {
	return &PayCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PayCallback 支付回调（第三方调用）
func (l *PayCallbackLogic) PayCallback(req *types.PayCallbackReq) (resp *types.PayCallbackResp, err error) {
	// 1. 验签（初版跳过，后续对接真实支付 SDK 验签）

	// 2. 更新支付流水状态
	err = l.svcCtx.PaymentFlowModel.UpdateStatus(l.ctx, req.OrderId, req.Status, req.TradeNo)
	if err != nil {
		logx.Errorf("更新支付状态失败: orderId=%s, err=%v", req.OrderId, err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	// 3. 调用 Order RPC → 更新订单状态为已支付(1)
	if req.Status == 1 {
		_, err = l.svcCtx.OrderRpc.UpdateOrderStatus(l.ctx, &orderrpc.UpdateOrderStatusReq{
			OrderId: req.OrderId,
			Status:  1, // 已支付
		})
		if err != nil {
			logx.Errorf("更新订单状态为已支付失败: orderId=%s, err=%v", req.OrderId, err)
		}
	}

	logx.Infof("支付回调成功: orderId=%s, tradeNo=%s, status=%d", req.OrderId, req.TradeNo, req.Status)
	return &types.PayCallbackResp{Msg: "success"}, nil
}
