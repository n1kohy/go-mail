// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	// PostgreSQL 数据源（用于查询用户进行验证）
	DataSource string
	// Redis 缓存配置
	Cache cache.CacheConf
	// JWT 认证配置
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
}
