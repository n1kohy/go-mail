package logic

import (
	"context"

	"go-mail/services/cart/rpc/cart"
	"go-mail/services/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetCartItemsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCartItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCartItemsLogic {
	return &GetCartItemsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetCartItems 获取用户选中的购物车商品（供 Order 服务调用）
func (l *GetCartItemsLogic) GetCartItems(in *cart.GetCartItemsReq) (*cart.GetCartItemsResp, error) {
	items, err := l.svcCtx.CartItemModel.FindSelectedByUserId(l.ctx, in.UserId)
	if err != nil {
		logx.Errorf("查询购物车失败: userId=%d, err=%v", in.UserId, err)
		return nil, status.Error(codes.Internal, "查询购物车失败")
	}

	var list []*cart.CartItemInfo
	for _, item := range items {
		list = append(list, &cart.CartItemInfo{
			SkuId:       item.SkuId,
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			SkuSpecs:    item.SkuSpecs,
			Price:       item.Price,
			Quantity:    int32(item.Quantity),
		})
	}

	return &cart.GetCartItemsResp{Items: list}, nil
}
