// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"go-mail/services/user/api/internal/config"
	"go-mail/services/user/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config           config.Config
	UserModel        model.UserModel
	UserAddressModel model.UserAddressModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewSqlConn("postgres", c.DataSource)
	return &ServiceContext{
		Config:           c,
		UserModel:        model.NewUserModel(conn, c.Cache),
		UserAddressModel: model.NewUserAddressModel(conn, c.Cache),
	}
}
