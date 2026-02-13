package cart

import (
	"context"
	"encoding/json"

	"go-mail/common/xerr"
	"go-mail/services/cart/api/internal/svc"
	"go-mail/services/cart/api/internal/types"
	"go-mail/services/cart/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartAddLogic {
	return &CartAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CartAdd 加入购物车（Upsert：存在则累加数量）
func (l *CartAddLogic) CartAdd(req *types.CartAddReq) (resp *types.CartAddResp, err error) {
	// 从 JWT 中提取 userId
	userId, err := l.getUserId()
	if err != nil {
		return nil, err
	}

	// Upsert: 存在相同 user_id + sku_id 则累加数量
	err = l.svcCtx.CartItemModel.Upsert(l.ctx, &model.CartItem{
		UserId:      userId,
		ProductId:   0, // 初版暂不校验商品，后续接入 Product RPC
		SkuId:       req.SkuId,
		ProductName: "", // 后续从 Product RPC 获取
		SkuSpecs:    "",
		Price:       0, // 后续从 Product RPC 获取最新价格
		Quantity:    req.Quantity,
		Selected:    true,
	})
	if err != nil {
		logx.Errorf("加入购物车失败: userId=%d, skuId=%d, err=%v", userId, req.SkuId, err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	return &types.CartAddResp{Msg: "已加入购物车"}, nil
}

func (l *CartAddLogic) getUserId() (int64, error) {
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
