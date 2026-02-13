package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	DataSource       string
	Cache            cache.CacheConf
	Auth             struct {
		AccessSecret string
		AccessExpire int64
	}
	CartRpcConf      zrpc.RpcClientConf // Cart RPC (9904)
	ProductRpcConf   zrpc.RpcClientConf // Product RPC (9902)
	InventoryRpcConf zrpc.RpcClientConf // Inventory RPC (9903)
	PromotionRpcConf zrpc.RpcClientConf // Promotion RPC (9905)
}
