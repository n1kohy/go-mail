package search

import (
	"context"

	"go-mail/services/search/api/internal/svc"
	"go-mail/services/search/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchHotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchHotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchHotLogic {
	return &SearchHotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SearchHot 热搜词（初版：静态数据，后续接入 Redis 热搜排行）
func (l *SearchHotLogic) SearchHot() (resp *types.SearchHotResp, err error) {
	// 初版使用静态热搜词，后续从 Redis ZREVRANGE 获取
	return &types.SearchHotResp{
		Keywords: []string{
			"iPhone",
			"连衣裙",
			"Switch",
			"笔记本电脑",
			"蓝牙耳机",
			"运动鞋",
			"面膜",
			"零食大礼包",
		},
	}, nil
}
