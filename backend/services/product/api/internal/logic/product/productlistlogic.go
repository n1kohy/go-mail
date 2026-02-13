package product

import (
	"context"

	"go-mail/common/xerr"
	"go-mail/services/product/api/internal/svc"
	"go-mail/services/product/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductListLogic {
	return &ProductListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ProductList 商品列表（支持分类筛选+分页）
func (l *ProductListLogic) ProductList(req *types.ProductListReq) (resp *types.ProductListResp, err error) {
	var total int64

	if req.CategoryId > 0 {
		// 按分类查询
		products, err := l.svcCtx.ProductModel.FindListByCategoryId(l.ctx, req.CategoryId, req.Page, req.Size)
		if err != nil {
			logx.Errorf("查询商品列表失败: %v", err)
			return nil, xerr.NewCodeError(xerr.ServerError)
		}
		total, _ = l.svcCtx.ProductModel.CountByCategoryId(l.ctx, req.CategoryId)

		var list []types.ProductListItem
		for _, p := range products {
			// 查第一个 SKU 的价格作为展示价
			skus, _ := l.svcCtx.ProductSkuModel.FindByProductId(l.ctx, p.Id)
			var price float64
			if len(skus) > 0 {
				price = skus[0].Price
			}
			list = append(list, types.ProductListItem{
				Id:        p.Id,
				Name:      p.Name,
				MainImage: p.MainImage,
				Price:     price,
			})
		}
		return &types.ProductListResp{List: list, Total: total}, nil
	}

	// 查全部
	products, err := l.svcCtx.ProductModel.FindAll(l.ctx, req.Page, req.Size)
	if err != nil {
		logx.Errorf("查询商品列表失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}
	total, _ = l.svcCtx.ProductModel.CountAll(l.ctx)

	var list []types.ProductListItem
	for _, p := range products {
		skus, _ := l.svcCtx.ProductSkuModel.FindByProductId(l.ctx, p.Id)
		var price float64
		if len(skus) > 0 {
			price = skus[0].Price
		}
		list = append(list, types.ProductListItem{
			Id:        p.Id,
			Name:      p.Name,
			MainImage: p.MainImage,
			Price:     price,
		})
	}
	return &types.ProductListResp{List: list, Total: total}, nil
}
