package svc

import (
	"go-mail/services/order/model"
	"go-mail/services/order/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config           config.Config
	OrderMasterModel model.OrderMasterModel
	OrderItemModel   model.OrderItemModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewSqlConn("postgres", c.DataSource)
	return &ServiceContext{
		Config:           c,
		OrderMasterModel: model.NewOrderMasterModel(conn, c.Cache),
		OrderItemModel:   model.NewOrderItemModel(conn, c.Cache),
	}
}
