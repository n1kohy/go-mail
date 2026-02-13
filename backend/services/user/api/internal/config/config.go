// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	// PostgreSQL 数据源
	DataSource string
	// Redis 缓存配置
	Cache cache.CacheConf
	// JWT 认证配置（goctl 从 .api 的 jwt: Auth 自动生成）
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
}
