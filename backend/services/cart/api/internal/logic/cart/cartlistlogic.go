package cart

import (
	"context"
	"encoding/json"

	"go-mail/common/xerr"
	"go-mail/services/cart/api/internal/svc"
	"go-mail/services/cart/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartListLogic {
	return &CartListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CartList 购物车列表
func (l *CartListLogic) CartList() (resp *types.CartListResp, err error) {
	userId, err := l.getUserId()
	if err != nil {
		return nil, err
	}

	cartItems, err := l.svcCtx.CartItemModel.FindByUserId(l.ctx, userId)
	if err != nil {
		logx.Errorf("查询购物车失败: userId=%d, err=%v", userId, err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	var items []types.CartItem
	var totalAmount float64
	for _, item := range cartItems {
		items = append(items, types.CartItem{
			Id:          item.Id,
			SkuId:       item.SkuId,
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			Specs:       item.SkuSpecs,
			Price:       item.Price,
			Quantity:    item.Quantity,
			Selected:    item.Selected,
		})
		if item.Selected {
			totalAmount += item.Price * float64(item.Quantity)
		}
	}

	if items == nil {
		items = []types.CartItem{}
	}

	return &types.CartListResp{Items: items, TotalAmount: totalAmount}, nil
}

func (l *CartListLogic) getUserId() (int64, error) {
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
