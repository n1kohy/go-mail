package product

import (
	"context"

	"go-mail/common/xerr"
	"go-mail/services/product/api/internal/svc"
	"go-mail/services/product/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductDetailLogic {
	return &ProductDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ProductDetail 商品详情（SPU + SKU 联查）
func (l *ProductDetailLogic) ProductDetail(id int64) (resp *types.ProductDetailResp, err error) {
	// 查询商品
	product, err := l.svcCtx.ProductModel.FindOne(l.ctx, id)
	if err != nil {
		logx.Errorf("查询商品失败: id=%d, err=%v", id, err)
		return nil, xerr.NewCodeErrorMsg(xerr.NotFound, "商品不存在")
	}

	// 查询分类名称
	categoryName := ""
	category, err := l.svcCtx.CategoryModel.FindOne(l.ctx, product.CategoryId)
	if err == nil {
		categoryName = category.Name
	}

	// 查询 SKU 列表
	skus, err := l.svcCtx.ProductSkuModel.FindByProductId(l.ctx, id)
	if err != nil {
		logx.Errorf("查询 SKU 失败: productId=%d, err=%v", id, err)
	}

	var skuItems []types.SkuItem
	for _, sku := range skus {
		skuItems = append(skuItems, types.SkuItem{
			SkuId: sku.Id,
			Specs: sku.Specs,
			Price: sku.Price,
			Image: sku.Image,
		})
	}

	return &types.ProductDetailResp{
		Id:         product.Id,
		Name:       product.Name,
		SubTitle:   product.SubTitle,
		MainImage:  product.MainImage,
		DetailHtml: product.DetailHtml,
		Category:   categoryName,
		Status:     int(product.Status),
		Skus:       skuItems,
	}, nil
}
