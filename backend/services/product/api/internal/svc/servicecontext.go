package svc

import (
	"go-mail/services/product/api/internal/config"
	"go-mail/services/product/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	ProductModel    model.ProductModel
	ProductSkuModel model.ProductSkuModel
	CategoryModel   model.CategoryModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewSqlConn("postgres", c.DataSource)
	return &ServiceContext{
		Config:          c,
		ProductModel:    model.NewProductModel(conn, c.Cache),
		ProductSkuModel: model.NewProductSkuModel(conn, c.Cache),
		CategoryModel:   model.NewCategoryModel(conn, c.Cache),
	}
}
