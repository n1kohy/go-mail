package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	DataSource string
	Cache      cache.CacheConf
	Auth       struct {
		AccessSecret string
		AccessExpire int64
	}
	OrderRpcConf zrpc.RpcClientConf // 调用 Order RPC 更新订单状态
}
