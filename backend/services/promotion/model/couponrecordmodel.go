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
	CouponRecordModel interface {
		Insert(ctx context.Context, data *CouponRecord) (sql.Result, error)
		FindByUserIdAndCouponId(ctx context.Context, userId, couponId int64) (*CouponRecord, error)
		FindByUserId(ctx context.Context, userId int64, status int) ([]*CouponRecord, error)
		MarkUsed(ctx context.Context, userId, couponId int64, orderId string) error
	}

	CouponRecord struct {
		Id          int64     `db:"id"`
		UserId      int64     `db:"user_id"`
		CouponId    int64     `db:"coupon_id"`
		Status      int64     `db:"status"`
		UsedOrderId string    `db:"used_order_id"`
		CreateTime  time.Time `db:"create_time"`
	}

	defaultCouponRecordModel struct {
		sqlc.CachedConn
		table string
	}
)

func NewCouponRecordModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CouponRecordModel {
	return &defaultCouponRecordModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "coupon_record",
	}
}

func (m *defaultCouponRecordModel) Insert(ctx context.Context, data *CouponRecord) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, coupon_id, status) VALUES ($1, $2, $3)", m.table)
	return m.ExecNoCacheCtx(ctx, query, data.UserId, data.CouponId, data.Status)
}

func (m *defaultCouponRecordModel) FindByUserIdAndCouponId(ctx context.Context, userId, couponId int64) (*CouponRecord, error) {
	var resp CouponRecord
	query := fmt.Sprintf("SELECT id, user_id, coupon_id, status, used_order_id, create_time FROM %s WHERE user_id = $1 AND coupon_id = $2 LIMIT 1", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, userId, couponId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultCouponRecordModel) FindByUserId(ctx context.Context, userId int64, status int) ([]*CouponRecord, error) {
	var list []*CouponRecord
	query := fmt.Sprintf("SELECT id, user_id, coupon_id, status, used_order_id, create_time FROM %s WHERE user_id = $1 AND status = $2 ORDER BY create_time DESC", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, userId, status)
	return list, err
}

func (m *defaultCouponRecordModel) MarkUsed(ctx context.Context, userId, couponId int64, orderId string) error {
	query := fmt.Sprintf("UPDATE %s SET status = 1, used_order_id = $1 WHERE user_id = $2 AND coupon_id = $3 AND status = 0", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, orderId, userId, couponId)
	return err
}
