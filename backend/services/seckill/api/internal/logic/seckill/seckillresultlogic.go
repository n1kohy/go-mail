package seckill

import (
	"context"
	"encoding/json"

	"go-mail/common/xerr"
	"go-mail/services/seckill/api/internal/svc"
	"go-mail/services/seckill/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SeckillResultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSeckillResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillResultLogic {
	return &SeckillResultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SeckillResult 查询秒杀结果
func (l *SeckillResultLogic) SeckillResult(req *types.SeckillResultReq) (resp *types.SeckillResultResp, err error) {
	userId, err := l.getUserId()
	if err != nil {
		return nil, err
	}

	// 初版：模拟查询结果
	// 后续流程：查 Redis 秒杀结果 key = seckill:result:{userId}:{skuId}
	// status: 0=排队中, 1=成功, 2=失败

	logx.Infof("查询秒杀结果: userId=%d, skuId=%d", userId, req.SkuId)

	return &types.SeckillResultResp{
		Status:  0,  // 排队中
		OrderId: "", // 空表示未出结果
	}, nil
}

func (l *SeckillResultLogic) getUserId() (int64, error) {
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
