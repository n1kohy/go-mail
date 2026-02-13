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
	PaymentFlowModel interface {
		Insert(ctx context.Context, data *PaymentFlow) (sql.Result, error)
		FindByOrderId(ctx context.Context, orderId string) (*PaymentFlow, error)
		UpdateStatus(ctx context.Context, orderId string, status int, tradeNo string) error
		UpdateRefund(ctx context.Context, orderId string, refundAmount float64) error
	}

	PaymentFlow struct {
		Id           int64     `db:"id"`
		OrderId      string    `db:"order_id"`
		TradeNo      string    `db:"trade_no"`
		Amount       float64   `db:"amount"`
		Channel      int64     `db:"channel"`
		Status       int64     `db:"status"`
		CallbackTime time.Time `db:"callback_time"`
		RefundAmount float64   `db:"refund_amount"`
		CreateTime   time.Time `db:"create_time"`
		UpdateTime   time.Time `db:"update_time"`
	}

	defaultPaymentFlowModel struct {
		sqlc.CachedConn
		table string
	}
)

func NewPaymentFlowModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PaymentFlowModel {
	return &defaultPaymentFlowModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "payment_flow",
	}
}

func (m *defaultPaymentFlowModel) Insert(ctx context.Context, data *PaymentFlow) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (order_id, amount, channel, status) VALUES ($1, $2, $3, $4)", m.table)
	return m.ExecNoCacheCtx(ctx, query, data.OrderId, data.Amount, data.Channel, data.Status)
}

func (m *defaultPaymentFlowModel) FindByOrderId(ctx context.Context, orderId string) (*PaymentFlow, error) {
	var resp PaymentFlow
	query := fmt.Sprintf("SELECT id, order_id, trade_no, amount, channel, status, callback_time, refund_amount, create_time, update_time FROM %s WHERE order_id = $1 ORDER BY id DESC LIMIT 1", m.table)
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

func (m *defaultPaymentFlowModel) UpdateStatus(ctx context.Context, orderId string, status int, tradeNo string) error {
	query := fmt.Sprintf("UPDATE %s SET status = $1, trade_no = $2, callback_time = NOW(), update_time = NOW() WHERE order_id = $3 AND status = 0", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, status, tradeNo, orderId)
	return err
}

func (m *defaultPaymentFlowModel) UpdateRefund(ctx context.Context, orderId string, refundAmount float64) error {
	query := fmt.Sprintf("UPDATE %s SET status = 3, refund_amount = $1, update_time = NOW() WHERE order_id = $2 AND status = 1", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, refundAmount, orderId)
	return err
}
