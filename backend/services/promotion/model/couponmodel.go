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
	CouponModel interface {
		FindAvailable(ctx context.Context) ([]*Coupon, error)
		FindOne(ctx context.Context, id int64) (*Coupon, error)
		DecrRemainCount(ctx context.Context, id int64) error
	}

	Coupon struct {
		Id          int64     `db:"id"`
		Name        string    `db:"name"`
		Type        int64     `db:"type"`
		Threshold   float64   `db:"threshold"`
		Discount    float64   `db:"discount"`
		TotalCount  int       `db:"total_count"`
		RemainCount int       `db:"remain_count"`
		StartTime   time.Time `db:"start_time"`
		EndTime     time.Time `db:"end_time"`
		Status      int64     `db:"status"`
		CreateTime  time.Time `db:"create_time"`
	}

	defaultCouponModel struct {
		sqlc.CachedConn
		table string
	}
)

func NewCouponModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CouponModel {
	return &defaultCouponModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "coupon",
	}
}

func (m *defaultCouponModel) FindAvailable(ctx context.Context) ([]*Coupon, error) {
	var list []*Coupon
	query := fmt.Sprintf("SELECT id, name, type, threshold, discount, total_count, remain_count, start_time, end_time, status, create_time FROM %s WHERE status = 1 AND remain_count > 0 AND start_time <= NOW() AND end_time >= NOW() ORDER BY id DESC", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &list, query)
	return list, err
}

var cacheCouponIdPrefix = "cache:coupon:id:"

func (m *defaultCouponModel) FindOne(ctx context.Context, id int64) (*Coupon, error) {
	key := fmt.Sprintf("%s%v", cacheCouponIdPrefix, id)
	var resp Coupon
	err := m.QueryRowCtx(ctx, &resp, key, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("SELECT id, name, type, threshold, discount, total_count, remain_count, start_time, end_time, status, create_time FROM %s WHERE id = $1 LIMIT 1", m.table)
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

func (m *defaultCouponModel) DecrRemainCount(ctx context.Context, id int64) error {
	query := fmt.Sprintf("UPDATE %s SET remain_count = remain_count - 1 WHERE id = $1 AND remain_count > 0", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, id)
	if err != nil {
		return err
	}
	// 手动删除缓存
	key := fmt.Sprintf("%s%v", cacheCouponIdPrefix, id)
	_ = m.DelCacheCtx(ctx, key)
	return nil
}
