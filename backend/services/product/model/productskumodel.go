package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	ProductSkuModel interface {
		FindByProductId(ctx context.Context, productId int64) ([]*ProductSku, error)
	}

	ProductSku struct {
		Id         int64     `db:"id"`
		ProductId  int64     `db:"product_id"`
		Specs      string    `db:"specs"`
		Price      float64   `db:"price"`
		Image      string    `db:"image"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}

	defaultProductSkuModel struct {
		sqlc.CachedConn
		table string
	}
)

func NewProductSkuModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductSkuModel {
	return &defaultProductSkuModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "product_sku",
	}
}

// FindByProductId 查询商品的所有 SKU
func (m *defaultProductSkuModel) FindByProductId(ctx context.Context, productId int64) ([]*ProductSku, error) {
	var list []*ProductSku
	query := fmt.Sprintf("SELECT id, product_id, specs, price, image, create_time, update_time FROM %s WHERE product_id = $1 ORDER BY id", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, productId)
	return list, err
}
