package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	DataSource string // 连接 db_inventory（扣减库存）
	Cache      cache.CacheConf
	Auth       struct {
		AccessSecret string
		AccessExpire int64
	}
}
