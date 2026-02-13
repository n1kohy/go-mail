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
	OrderMasterModel interface {
		Insert(ctx context.Context, data *OrderMaster) (sql.Result, error)
		FindOne(ctx context.Context, id string) (*OrderMaster, error)
		FindByUserId(ctx context.Context, userId int64, status, page, size int) ([]*OrderMaster, error)
		CountByUserId(ctx context.Context, userId int64, status int) (int64, error)
		UpdateStatus(ctx context.Context, id string, status int) error
	}

	OrderMaster struct {
		Id              string    `db:"id"`
		UserId          int64     `db:"user_id"`
		TotalAmount     float64   `db:"total_amount"`
		DiscountAmount  float64   `db:"discount_amount"`
		FreightAmount   float64   `db:"freight_amount"`
		PayAmount       float64   `db:"pay_amount"`
		CouponId        int64     `db:"coupon_id"`
		PayType         int64     `db:"pay_type"`
		Status          int64     `db:"status"`
		AddressSnapshot string    `db:"address_snapshot"`
		ExpireTime      time.Time `db:"expire_time"`
		CreateTime      time.Time `db:"create_time"`
		UpdateTime      time.Time `db:"update_time"`
	}

	defaultOrderMasterModel struct {
		sqlc.CachedConn
		table string
	}
)

func NewOrderMasterModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderMasterModel {
	return &defaultOrderMasterModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "order_master",
	}
}

var cacheOrderIdPrefix = "cache:order:id:"

func (m *defaultOrderMasterModel) Insert(ctx context.Context, data *OrderMaster) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (id, user_id, total_amount, discount_amount, freight_amount, pay_amount, coupon_id, status, address_snapshot, expire_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", m.table)
	return m.ExecNoCacheCtx(ctx, query, data.Id, data.UserId, data.TotalAmount, data.DiscountAmount, data.FreightAmount, data.PayAmount, data.CouponId, data.Status, data.AddressSnapshot, data.ExpireTime)
}

func (m *defaultOrderMasterModel) FindOne(ctx context.Context, id string) (*OrderMaster, error) {
	key := fmt.Sprintf("%s%v", cacheOrderIdPrefix, id)
	var resp OrderMaster
	err := m.QueryRowCtx(ctx, &resp, key, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("SELECT id, user_id, total_amount, discount_amount, freight_amount, pay_amount, coupon_id, pay_type, status, address_snapshot, expire_time, create_time, update_time FROM %s WHERE id = $1 LIMIT 1", m.table)
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

func (m *defaultOrderMasterModel) FindByUserId(ctx context.Context, userId int64, status, page, size int) ([]*OrderMaster, error) {
	var list []*OrderMaster
	offset := (page - 1) * size
	var query string
	if status >= 0 {
		query = fmt.Sprintf("SELECT id, user_id, total_amount, discount_amount, freight_amount, pay_amount, coupon_id, pay_type, status, expire_time, create_time, update_time FROM %s WHERE user_id = $1 AND status = $2 ORDER BY create_time DESC LIMIT $3 OFFSET $4", m.table)
		err := m.QueryRowsNoCacheCtx(ctx, &list, query, userId, status, size, offset)
		return list, err
	}
	query = fmt.Sprintf("SELECT id, user_id, total_amount, discount_amount, freight_amount, pay_amount, coupon_id, pay_type, status, expire_time, create_time, update_time FROM %s WHERE user_id = $1 ORDER BY create_time DESC LIMIT $2 OFFSET $3", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, userId, size, offset)
	return list, err
}

func (m *defaultOrderMasterModel) CountByUserId(ctx context.Context, userId int64, status int) (int64, error) {
	var count int64
	var query string
	if status >= 0 {
		query = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id = $1 AND status = $2", m.table)
		err := m.QueryRowNoCacheCtx(ctx, &count, query, userId, status)
		return count, err
	}
	query = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id = $1", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &count, query, userId)
	return count, err
}

func (m *defaultOrderMasterModel) UpdateStatus(ctx context.Context, id string, status int) error {
	query := fmt.Sprintf("UPDATE %s SET status = $1, update_time = NOW() WHERE id = $2", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, status, id)
	if err != nil {
		return err
	}
	// 手动删除缓存
	key := fmt.Sprintf("%s%v", cacheOrderIdPrefix, id)
	_ = m.DelCacheCtx(ctx, key)
	return nil
}
