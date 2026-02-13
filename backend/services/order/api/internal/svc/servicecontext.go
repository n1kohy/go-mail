package svc

import (
	"go-mail/services/cart/rpc/cartrpc"
	"go-mail/services/inventory/rpc/inventoryrpc"
	"go-mail/services/order/api/internal/config"
	"go-mail/services/order/model"
	"go-mail/services/product/rpc/productrpc"
	"go-mail/services/promotion/rpc/promotionrpc"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config           config.Config
	OrderMasterModel model.OrderMasterModel
	OrderItemModel   model.OrderItemModel
	CartRpc          cartrpc.CartRpc
	ProductRpc       productrpc.ProductRpc
	InventoryRpc     inventoryrpc.InventoryRpc
	PromotionRpc     promotionrpc.PromotionRpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewSqlConn("postgres", c.DataSource)
	return &ServiceContext{
		Config:           c,
		OrderMasterModel: model.NewOrderMasterModel(conn, c.Cache),
		OrderItemModel:   model.NewOrderItemModel(conn, c.Cache),
		CartRpc:          cartrpc.NewCartRpc(zrpc.MustNewClient(c.CartRpcConf)),
		ProductRpc:       productrpc.NewProductRpc(zrpc.MustNewClient(c.ProductRpcConf)),
		InventoryRpc:     inventoryrpc.NewInventoryRpc(zrpc.MustNewClient(c.InventoryRpcConf)),
		PromotionRpc:     promotionrpc.NewPromotionRpc(zrpc.MustNewClient(c.PromotionRpcConf)),
	}
}
