package svc

import (
	"go-mail/services/logistics/api/internal/config"
	"go-mail/services/logistics/model"
	"go-mail/services/order/rpc/orderrpc"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config             config.Config
	ShippingOrderModel model.ShippingOrderModel
	OrderRpc           orderrpc.OrderRpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewSqlConn("postgres", c.DataSource)
	return &ServiceContext{
		Config:             c,
		ShippingOrderModel: model.NewShippingOrderModel(conn, c.Cache),
		OrderRpc:           orderrpc.NewOrderRpc(zrpc.MustNewClient(c.OrderRpcConf)),
	}
}
