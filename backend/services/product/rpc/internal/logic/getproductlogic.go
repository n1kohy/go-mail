package logic

import (
	"context"

	"go-mail/services/product/model"
	"go-mail/services/product/rpc/internal/svc"
	"go-mail/services/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetProduct 根据 ID 查询商品
func (l *GetProductLogic) GetProduct(in *product.GetProductReq) (*product.GetProductResp, error) {
	p, err := l.svcCtx.ProductModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(codes.NotFound, "商品不存在")
		}
		logx.Errorf("查询商品失败: id=%d, err=%v", in.Id, err)
		return nil, status.Error(codes.Internal, "查询商品失败")
	}

	return &product.GetProductResp{
		Id:         p.Id,
		Name:       p.Name,
		MainImage:  p.MainImage,
		CategoryId: p.CategoryId,
		Status:     int32(p.Status),
	}, nil
}
