package svc

import (
	"go-mail/services/logistics/model"
	"go-mail/services/logistics/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config             config.Config
	ShippingOrderModel model.ShippingOrderModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewSqlConn("postgres", c.DataSource)
	return &ServiceContext{
		Config:             c,
		ShippingOrderModel: model.NewShippingOrderModel(conn, c.Cache),
	}
}
