package payment

import (
	"context"
	"encoding/json"
	"fmt"

	"go-mail/common/xerr"
	"go-mail/services/payment/api/internal/svc"
	"go-mail/services/payment/api/internal/types"
	"go-mail/services/payment/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayLogic {
	return &PayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Pay 发起支付
func (l *PayLogic) Pay(req *types.PayReq) (resp *types.PayResp, err error) {
	_, err = l.getUserId()
	if err != nil {
		return nil, err
	}

	// 创建支付流水
	_, err = l.svcCtx.PaymentFlowModel.Insert(l.ctx, &model.PaymentFlow{
		OrderId: req.OrderId,
		Amount:  0, // 后续从 Order RPC 获取真实金额
		Channel: int64(req.Channel),
		Status:  0,
	})
	if err != nil {
		logx.Errorf("创建支付流水失败: orderId=%s, err=%v", req.OrderId, err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	// 初版：模拟支付链接（后续对接支付宝/微信 SDK）
	payUrl := fmt.Sprintf("https://pay.example.com/checkout?order_id=%s&channel=%d", req.OrderId, req.Channel)

	logx.Infof("支付发起: orderId=%s, channel=%d", req.OrderId, req.Channel)
	return &types.PayResp{PayUrl: payUrl}, nil
}

func (l *PayLogic) getUserId() (int64, error) {
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
