package svc

import (
	"go-mail/services/order/rpc/orderrpc"
	"go-mail/services/payment/api/internal/config"
	"go-mail/services/payment/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config           config.Config
	PaymentFlowModel model.PaymentFlowModel
	OrderRpc         orderrpc.OrderRpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewSqlConn("postgres", c.DataSource)
	return &ServiceContext{
		Config:           c,
		PaymentFlowModel: model.NewPaymentFlowModel(conn, c.Cache),
		OrderRpc:         orderrpc.NewOrderRpc(zrpc.MustNewClient(c.OrderRpcConf)),
	}
}
