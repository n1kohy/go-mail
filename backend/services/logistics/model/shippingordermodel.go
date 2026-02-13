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
	ShippingOrderModel interface {
		Insert(ctx context.Context, data *ShippingOrder) (sql.Result, error)
		FindByOrderId(ctx context.Context, orderId string) (*ShippingOrder, error)
		UpdateStatus(ctx context.Context, orderId string, status int) error
	}

	ShippingOrder struct {
		Id         int64     `db:"id"`
		OrderId    string    `db:"order_id"`
		TrackingNo string    `db:"tracking_no"`
		Carrier    string    `db:"carrier"`
		Status     int64     `db:"status"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}

	defaultShippingOrderModel struct {
		sqlc.CachedConn
		table string
	}
)

func NewShippingOrderModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ShippingOrderModel {
	return &defaultShippingOrderModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "shipping_order",
	}
}

func (m *defaultShippingOrderModel) Insert(ctx context.Context, data *ShippingOrder) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (order_id, tracking_no, carrier, status) VALUES ($1, $2, $3, $4)", m.table)
	return m.ExecNoCacheCtx(ctx, query, data.OrderId, data.TrackingNo, data.Carrier, data.Status)
}

func (m *defaultShippingOrderModel) FindByOrderId(ctx context.Context, orderId string) (*ShippingOrder, error) {
	var resp ShippingOrder
	query := fmt.Sprintf("SELECT id, order_id, tracking_no, carrier, status, create_time, update_time FROM %s WHERE order_id = $1 ORDER BY id DESC LIMIT 1", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, orderId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultShippingOrderModel) UpdateStatus(ctx context.Context, orderId string, status int) error {
	query := fmt.Sprintf("UPDATE %s SET status = $1, update_time = NOW() WHERE order_id = $2", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, status, orderId)
	return err
}
