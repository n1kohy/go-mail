package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	CategoryModel interface {
		FindOne(ctx context.Context, id int64) (*Category, error)
	}

	Category struct {
		Id       int64  `db:"id"`
		ParentId int64  `db:"parent_id"`
		Name     string `db:"name"`
		Level    int64  `db:"level"`
		Sort     int64  `db:"sort"`
	}

	defaultCategoryModel struct {
		sqlc.CachedConn
		table string
	}
)

func NewCategoryModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CategoryModel {
	return &defaultCategoryModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "category",
	}
}

var cacheCategoryIdPrefix = "cache:category:id:"

func (m *defaultCategoryModel) FindOne(ctx context.Context, id int64) (*Category, error) {
	key := fmt.Sprintf("%s%v", cacheCategoryIdPrefix, id)
	var resp Category
	err := m.QueryRowCtx(ctx, &resp, key, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("SELECT id, parent_id, name, level, sort FROM %s WHERE id = $1 LIMIT 1", m.table)
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
