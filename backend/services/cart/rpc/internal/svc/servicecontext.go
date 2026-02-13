package svc

import (
	"go-mail/services/cart/model"
	"go-mail/services/cart/rpc/internal/config"

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
