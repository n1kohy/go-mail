package svc

import (
	"go-mail/services/inventory/model"
	"go-mail/services/inventory/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	StockModel model.StockModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewSqlConn("postgres", c.DataSource)
	return &ServiceContext{
		Config:     c,
		StockModel: model.NewStockModel(conn, c.Cache),
	}
}
