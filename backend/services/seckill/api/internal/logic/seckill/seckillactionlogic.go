package seckill

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-mail/common/xerr"
	"go-mail/services/seckill/api/internal/svc"
	"go-mail/services/seckill/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SeckillActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSeckillActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillActionLogic {
	return &SeckillActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SeckillAction 执行秒杀
// 初版：同步处理（后续接入 Redis 预减库存 + MQ 异步下单）
func (l *SeckillActionLogic) SeckillAction(req *types.SeckillActionReq) (resp *types.SeckillActionResp, err error) {
	userId, err := l.getUserId()
	if err != nil {
		return nil, err
	}

	// 初版：模拟秒杀逻辑
	// 后续流程：
	// 1. Redis DECR 预减库存
	// 2. 发送 MQ 消息异步处理
	// 3. 创建订单

	logx.Infof("秒杀请求: userId=%d, skuId=%d, time=%s", userId, req.SkuId, time.Now().Format("15:04:05.000"))

	// 模拟排队中响应
	return &types.SeckillActionResp{
		Code: 2002,
		Msg:  fmt.Sprintf("排队中，skuId=%d", req.SkuId),
	}, nil
}

func (l *SeckillActionLogic) getUserId() (int64, error) {
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
