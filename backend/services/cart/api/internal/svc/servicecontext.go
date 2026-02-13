package svc

import (
	"go-mail/services/cart/api/internal/config"
	"go-mail/services/cart/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	CartItemModel model.CartItemModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewSqlConn("postgres", c.DataSource)
	return &ServiceContext{
		Config:        c,
		CartItemModel: model.NewCartItemModel(conn, c.Cache),
	}
}
