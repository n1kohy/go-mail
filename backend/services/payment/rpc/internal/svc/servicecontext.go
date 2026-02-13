package svc

import (
	"go-mail/services/payment/model"
	"go-mail/services/payment/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config           config.Config
	PaymentFlowModel model.PaymentFlowModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewSqlConn("postgres", c.DataSource)
	return &ServiceContext{
		Config:           c,
		PaymentFlowModel: model.NewPaymentFlowModel(conn, c.Cache),
	}
}
