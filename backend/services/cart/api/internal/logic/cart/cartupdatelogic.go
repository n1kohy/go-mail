package cart

import (
	"context"
	"encoding/json"

	"go-mail/common/xerr"
	"go-mail/services/cart/api/internal/svc"
	"go-mail/services/cart/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartUpdateLogic {
	return &CartUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CartUpdate 修改购物车数量
func (l *CartUpdateLogic) CartUpdate(req *types.CartUpdateReq) (resp *types.CartUpdateResp, err error) {
	userId, err := l.getUserId()
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.CartItemModel.UpdateQuantity(l.ctx, userId, req.SkuId, req.Quantity)
	if err != nil {
		logx.Errorf("更新购物车数量失败: userId=%d, skuId=%d, err=%v", userId, req.SkuId, err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	return &types.CartUpdateResp{Msg: "更新成功"}, nil
}

func (l *CartUpdateLogic) getUserId() (int64, error) {
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
