package svc

import (
	"go-mail/services/seckill/api/internal/config"
)

type ServiceContext struct {
	Config config.Config
	// 后续注入：InventoryRpc / OrderRpc / Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
