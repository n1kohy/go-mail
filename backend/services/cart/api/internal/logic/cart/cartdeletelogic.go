package cart

import (
	"context"
	"encoding/json"

	"go-mail/common/xerr"
	"go-mail/services/cart/api/internal/svc"
	"go-mail/services/cart/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartDeleteLogic {
	return &CartDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CartDelete 删除购物车商品
func (l *CartDeleteLogic) CartDelete(req *types.CartDeleteReq) (resp *types.CartDeleteResp, err error) {
	userId, err := l.getUserId()
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.CartItemModel.DeleteByUserIdAndSkuIds(l.ctx, userId, req.SkuIds)
	if err != nil {
		logx.Errorf("删除购物车失败: userId=%d, skuIds=%v, err=%v", userId, req.SkuIds, err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	return &types.CartDeleteResp{Msg: "删除成功"}, nil
}

func (l *CartDeleteLogic) getUserId() (int64, error) {
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
