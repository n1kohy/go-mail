package payment

import (
	"context"
	"encoding/json"

	"go-mail/common/xerr"
	"go-mail/services/payment/api/internal/svc"
	"go-mail/services/payment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundLogic {
	return &RefundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Refund 退款
func (l *RefundLogic) Refund(req *types.RefundReq) (resp *types.RefundResp, err error) {
	_, err = l.getUserId()
	if err != nil {
		return nil, err
	}

	// 检查支付记录
	flow, err := l.svcCtx.PaymentFlowModel.FindByOrderId(l.ctx, req.OrderId)
	if err != nil {
		return nil, xerr.NewCodeErrorMsg(xerr.NotFound, "支付记录不存在")
	}

	// 仅已支付可退款
	if flow.Status != 1 {
		return nil, xerr.NewCodeErrorMsg(xerr.BadRequest, "当前状态无法退款")
	}

	// 退款金额校验
	if req.Amount > flow.Amount {
		return nil, xerr.NewCodeErrorMsg(xerr.BadRequest, "退款金额超过支付金额")
	}

	// 更新退款
	err = l.svcCtx.PaymentFlowModel.UpdateRefund(l.ctx, req.OrderId, req.Amount)
	if err != nil {
		logx.Errorf("退款失败: orderId=%s, err=%v", req.OrderId, err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	logx.Infof("退款成功: orderId=%s, amount=%.2f, reason=%s", req.OrderId, req.Amount, req.Reason)
	return &types.RefundResp{Msg: "退款申请成功"}, nil
}

func (l *RefundLogic) getUserId() (int64, error) {
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
