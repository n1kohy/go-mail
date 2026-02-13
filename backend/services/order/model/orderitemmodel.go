package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	OrderItemModel interface {
		Insert(ctx context.Context, data *OrderItem) (sql.Result, error)
		FindByOrderId(ctx context.Context, orderId string) ([]*OrderItem, error)
		CountByOrderId(ctx context.Context, orderId string) (int64, error)
	}

	OrderItem struct {
		Id          int64     `db:"id"`
		OrderId     string    `db:"order_id"`
		ProductId   int64     `db:"product_id"`
		ProductName string    `db:"product_name"`
		SkuId       int64     `db:"sku_id"`
		SkuSpecs    string    `db:"sku_specs"`
		Price       float64   `db:"price"`
		Quantity    int       `db:"quantity"`
		CreateTime  time.Time `db:"create_time"`
	}

	defaultOrderItemModel struct {
		sqlc.CachedConn
		table string
	}
)

func NewOrderItemModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderItemModel {
	return &defaultOrderItemModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "order_item",
	}
}

func (m *defaultOrderItemModel) Insert(ctx context.Context, data *OrderItem) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (order_id, product_id, product_name, sku_id, sku_specs, price, quantity) VALUES ($1, $2, $3, $4, $5, $6, $7)", m.table)
	return m.ExecNoCacheCtx(ctx, query, data.OrderId, data.ProductId, data.ProductName, data.SkuId, data.SkuSpecs, data.Price, data.Quantity)
}

func (m *defaultOrderItemModel) FindByOrderId(ctx context.Context, orderId string) ([]*OrderItem, error) {
	var list []*OrderItem
	query := fmt.Sprintf("SELECT id, order_id, product_id, product_name, sku_id, sku_specs, price, quantity, create_time FROM %s WHERE order_id = $1 ORDER BY id", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, orderId)
	return list, err
}

func (m *defaultOrderItemModel) CountByOrderId(ctx context.Context, orderId string) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE order_id = $1", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &count, query, orderId)
	return count, err
}
