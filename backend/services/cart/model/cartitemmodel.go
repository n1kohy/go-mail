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
	CartItemModel interface {
		FindByUserId(ctx context.Context, userId int64) ([]*CartItem, error)
		FindByUserIdAndSkuId(ctx context.Context, userId, skuId int64) (*CartItem, error)
		FindSelectedByUserId(ctx context.Context, userId int64) ([]*CartItem, error)
		Insert(ctx context.Context, data *CartItem) (sql.Result, error)
		UpdateQuantity(ctx context.Context, userId, skuId int64, quantity int) error
		DeleteByUserIdAndSkuIds(ctx context.Context, userId int64, skuIds []int64) error
		Upsert(ctx context.Context, data *CartItem) error
	}

	CartItem struct {
		Id          int64     `db:"id"`
		UserId      int64     `db:"user_id"`
		ProductId   int64     `db:"product_id"`
		SkuId       int64     `db:"sku_id"`
		ProductName string    `db:"product_name"`
		SkuSpecs    string    `db:"sku_specs"`
		Price       float64   `db:"price"`
		Quantity    int       `db:"quantity"`
		Selected    bool      `db:"selected"`
		CreateTime  time.Time `db:"create_time"`
		UpdateTime  time.Time `db:"update_time"`
	}

	defaultCartItemModel struct {
		sqlc.CachedConn
		table string
	}
)

func NewCartItemModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CartItemModel {
	return &defaultCartItemModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "cart_item",
	}
}

// FindByUserId 查询用户全部购物车
func (m *defaultCartItemModel) FindByUserId(ctx context.Context, userId int64) ([]*CartItem, error) {
	var list []*CartItem
	query := fmt.Sprintf("SELECT id, user_id, product_id, sku_id, product_name, sku_specs, price, quantity, selected, create_time, update_time FROM %s WHERE user_id = $1 ORDER BY create_time DESC", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, userId)
	return list, err
}

// FindByUserIdAndSkuId 查询单条购物车记录
func (m *defaultCartItemModel) FindByUserIdAndSkuId(ctx context.Context, userId, skuId int64) (*CartItem, error) {
	var resp CartItem
	query := fmt.Sprintf("SELECT id, user_id, product_id, sku_id, product_name, sku_specs, price, quantity, selected, create_time, update_time FROM %s WHERE user_id = $1 AND sku_id = $2 LIMIT 1", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, userId, skuId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindSelectedByUserId 查询用户选中的购物车商品
func (m *defaultCartItemModel) FindSelectedByUserId(ctx context.Context, userId int64) ([]*CartItem, error) {
	var list []*CartItem
	query := fmt.Sprintf("SELECT id, user_id, product_id, sku_id, product_name, sku_specs, price, quantity, selected, create_time, update_time FROM %s WHERE user_id = $1 AND selected = TRUE ORDER BY create_time DESC", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, userId)
	return list, err
}

// Insert 新增购物车记录
func (m *defaultCartItemModel) Insert(ctx context.Context, data *CartItem) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, product_id, sku_id, product_name, sku_specs, price, quantity, selected) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", m.table)
	return m.ExecNoCacheCtx(ctx, query, data.UserId, data.ProductId, data.SkuId, data.ProductName, data.SkuSpecs, data.Price, data.Quantity, data.Selected)
}

// Upsert 加入购物车：存在则累加数量，不存在则新增
func (m *defaultCartItemModel) Upsert(ctx context.Context, data *CartItem) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, product_id, sku_id, product_name, sku_specs, price, quantity, selected) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (user_id, sku_id) DO UPDATE SET quantity = %s.quantity + EXCLUDED.quantity, update_time = NOW()", m.table, m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, data.UserId, data.ProductId, data.SkuId, data.ProductName, data.SkuSpecs, data.Price, data.Quantity, data.Selected)
	return err
}

// UpdateQuantity 更新数量
func (m *defaultCartItemModel) UpdateQuantity(ctx context.Context, userId, skuId int64, quantity int) error {
	query := fmt.Sprintf("UPDATE %s SET quantity = $1, update_time = NOW() WHERE user_id = $2 AND sku_id = $3", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, quantity, userId, skuId)
	return err
}

// DeleteByUserIdAndSkuIds 批量删除购物车
func (m *defaultCartItemModel) DeleteByUserIdAndSkuIds(ctx context.Context, userId int64, skuIds []int64) error {
	if len(skuIds) == 0 {
		return nil
	}
	args := make([]any, 0, len(skuIds)+1)
	args = append(args, userId)
	holders := ""
	for i, id := range skuIds {
		if i > 0 {
			holders += ", "
		}
		holders += fmt.Sprintf("$%d", i+2)
		args = append(args, id)
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND sku_id IN (%s)", m.table, holders)
	_, err := m.ExecNoCacheCtx(ctx, query, args...)
	return err
}
