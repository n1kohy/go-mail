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
	StockModel interface {
		FindOneBySkuId(ctx context.Context, skuId int64) (*Stock, error)
		DeductStock(ctx context.Context, skuId int64, quantity int32) error
		RollbackStock(ctx context.Context, skuId int64, quantity int32) error
	}

	Stock struct {
		Id         int64     `db:"id"`
		SkuId      int64     `db:"sku_id"`
		Total      int32     `db:"total"`
		Available  int32     `db:"available"`
		Locked     int32     `db:"locked"`
		Version    int32     `db:"version"`
		UpdateTime time.Time `db:"update_time"`
	}

	defaultStockModel struct {
		sqlc.CachedConn
		table string
	}
)

func NewStockModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) StockModel {
	return &defaultStockModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "stock",
	}
}

// FindOneBySkuId 查询库存
func (m *defaultStockModel) FindOneBySkuId(ctx context.Context, skuId int64) (*Stock, error) {
	var resp Stock
	query := fmt.Sprintf("SELECT id, sku_id, total, available, locked, version, update_time FROM %s WHERE sku_id = $1 LIMIT 1", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, skuId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// DeductStock 乐观锁扣减库存
// UPDATE stock SET available = available - qty, locked = locked + qty, version = version + 1
// WHERE sku_id = ? AND available >= qty AND version = current_version
func (m *defaultStockModel) DeductStock(ctx context.Context, skuId int64, quantity int32) error {
	// 先获取当前版本
	stock, err := m.FindOneBySkuId(ctx, skuId)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET available = available - $1, locked = locked + $1, version = version + 1, update_time = NOW() WHERE sku_id = $2 AND available >= $1 AND version = $3", m.table)
	result, err := m.ExecNoCacheCtx(ctx, query, quantity, skuId, stock.Version)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows // 库存不足或版本冲突
	}
	return nil
}

// RollbackStock 回滚库存
func (m *defaultStockModel) RollbackStock(ctx context.Context, skuId int64, quantity int32) error {
	query := fmt.Sprintf("UPDATE %s SET available = available + $1, locked = locked - $1, version = version + 1, update_time = NOW() WHERE sku_id = $2 AND locked >= $1", m.table)
	result, err := m.ExecNoCacheCtx(ctx, query, quantity, skuId)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
