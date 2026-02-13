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

var cacheProductIdPrefix = "cache:product:id:"

type (
	ProductModel interface {
		FindOne(ctx context.Context, id int64) (*Product, error)
		FindListByCategoryId(ctx context.Context, categoryId int64, page, size int) ([]*Product, error)
		CountByCategoryId(ctx context.Context, categoryId int64) (int64, error)
		FindAll(ctx context.Context, page, size int) ([]*Product, error)
		CountAll(ctx context.Context) (int64, error)
		FindByIds(ctx context.Context, ids []int64) ([]*Product, error)
		Insert(ctx context.Context, data *Product) (sql.Result, error)
	}

	Product struct {
		Id         int64     `db:"id"`
		CategoryId int64     `db:"category_id"`
		Name       string    `db:"name"`
		SubTitle   string    `db:"sub_title"`
		MainImage  string    `db:"main_image"`
		DetailHtml string    `db:"detail_html"`
		Status     int64     `db:"status"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}

	defaultProductModel struct {
		sqlc.CachedConn
		table string
	}
)

func NewProductModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductModel {
	return &defaultProductModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "product",
	}
}

func (m *defaultProductModel) FindOne(ctx context.Context, id int64) (*Product, error) {
	key := fmt.Sprintf("%s%v", cacheProductIdPrefix, id)
	var resp Product
	err := m.QueryRowCtx(ctx, &resp, key, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("SELECT id, category_id, name, sub_title, main_image, detail_html, status, create_time, update_time FROM %s WHERE id = $1 LIMIT 1", m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultProductModel) FindListByCategoryId(ctx context.Context, categoryId int64, page, size int) ([]*Product, error) {
	var list []*Product
	offset := (page - 1) * size
	query := fmt.Sprintf("SELECT id, category_id, name, sub_title, main_image, status, create_time, update_time FROM %s WHERE category_id = $1 AND status = 1 ORDER BY id DESC LIMIT $2 OFFSET $3", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, categoryId, size, offset)
	return list, err
}

func (m *defaultProductModel) CountByCategoryId(ctx context.Context, categoryId int64) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE category_id = $1 AND status = 1", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &count, query, categoryId)
	return count, err
}

func (m *defaultProductModel) FindAll(ctx context.Context, page, size int) ([]*Product, error) {
	var list []*Product
	offset := (page - 1) * size
	query := fmt.Sprintf("SELECT id, category_id, name, sub_title, main_image, status, create_time, update_time FROM %s WHERE status = 1 ORDER BY id DESC LIMIT $1 OFFSET $2", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, size, offset)
	return list, err
}

func (m *defaultProductModel) CountAll(ctx context.Context) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE status = 1", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &count, query)
	return count, err
}

func (m *defaultProductModel) FindByIds(ctx context.Context, ids []int64) ([]*Product, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	// 构造 IN 子句: $1, $2, $3...
	args := make([]any, len(ids))
	holders := ""
	for i, id := range ids {
		if i > 0 {
			holders += ", "
		}
		holders += fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	var list []*Product
	query := fmt.Sprintf("SELECT id, category_id, name, sub_title, main_image, status, create_time, update_time FROM %s WHERE id IN (%s)", m.table, holders)
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...)
	return list, err
}

func (m *defaultProductModel) Insert(ctx context.Context, data *Product) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (category_id, name, sub_title, main_image, detail_html, status) VALUES ($1, $2, $3, $4, $5, $6)", m.table)
	return m.ExecNoCacheCtx(ctx, query, data.CategoryId, data.Name, data.SubTitle, data.MainImage, data.DetailHtml, data.Status)
}
