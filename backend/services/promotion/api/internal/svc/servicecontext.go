package svc

import (
	"go-mail/services/promotion/api/internal/config"
	"go-mail/services/promotion/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	CouponModel       model.CouponModel
	CouponRecordModel model.CouponRecordModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewSqlConn("postgres", c.DataSource)
	return &ServiceContext{
		Config:            c,
		CouponModel:       model.NewCouponModel(conn, c.Cache),
		CouponRecordModel: model.NewCouponRecordModel(conn, c.Cache),
	}
}
