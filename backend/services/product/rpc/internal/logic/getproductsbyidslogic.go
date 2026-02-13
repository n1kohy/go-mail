package logic

import (
	"context"

	"go-mail/services/product/rpc/internal/svc"
	"go-mail/services/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetProductsByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductsByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductsByIdsLogic {
	return &GetProductsByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetProductsByIds 批量查询商品
func (l *GetProductsByIdsLogic) GetProductsByIds(in *product.GetProductsByIdsReq) (*product.GetProductsByIdsResp, error) {
	if len(in.Ids) == 0 {
		return &product.GetProductsByIdsResp{}, nil
	}

	products, err := l.svcCtx.ProductModel.FindByIds(l.ctx, in.Ids)
	if err != nil {
		logx.Errorf("批量查询商品失败: err=%v", err)
		return nil, status.Error(codes.Internal, "批量查询商品失败")
	}

	var list []*product.GetProductResp
	for _, p := range products {
		list = append(list, &product.GetProductResp{
			Id:         p.Id,
			Name:       p.Name,
			MainImage:  p.MainImage,
			CategoryId: p.CategoryId,
			Status:     int32(p.Status),
		})
	}

	return &product.GetProductsByIdsResp{Products: list}, nil
}
