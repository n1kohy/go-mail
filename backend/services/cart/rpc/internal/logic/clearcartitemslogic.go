package logic

import (
	"context"

	"go-mail/services/cart/rpc/cart"
	"go-mail/services/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClearCartItemsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearCartItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearCartItemsLogic {
	return &ClearCartItemsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ClearCartItems 清空购物车中指定商品（下单成功后调用）
func (l *ClearCartItemsLogic) ClearCartItems(in *cart.ClearCartItemsReq) (*cart.ClearCartItemsResp, error) {
	err := l.svcCtx.CartItemModel.DeleteByUserIdAndSkuIds(l.ctx, in.UserId, in.SkuIds)
	if err != nil {
		logx.Errorf("清空购物车失败: userId=%d, err=%v", in.UserId, err)
		return nil, status.Error(codes.Internal, "清空购物车失败")
	}

	logx.Infof("购物车清理成功: userId=%d, skuIds=%v", in.UserId, in.SkuIds)
	return &cart.ClearCartItemsResp{Success: true}, nil
}
