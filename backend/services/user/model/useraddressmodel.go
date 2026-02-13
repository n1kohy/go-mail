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

var cacheUserAddressIdPrefix = "cache:userAddress:id:"

type (
	// UserAddressModel 自定义接口
	UserAddressModel interface {
		Insert(ctx context.Context, data *UserAddress) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*UserAddress, error)
		FindByUserId(ctx context.Context, userId int64) ([]*UserAddress, error)
		Update(ctx context.Context, data *UserAddress) error
		Delete(ctx context.Context, id int64) error
	}

	// UserAddress 对应数据库 user_address 表
	UserAddress struct {
		Id         int64     `db:"id"`
		UserId     int64     `db:"user_id"`
		Receiver   string    `db:"receiver"`
		Phone      string    `db:"phone"`
		Province   string    `db:"province"`
		City       string    `db:"city"`
		District   string    `db:"district"`
		Detail     string    `db:"detail"`
		IsDefault  bool      `db:"is_default"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}

	defaultUserAddressModel struct {
		sqlc.CachedConn
		table string
	}
)

// NewUserAddressModel 创建地址 Model 实例
func NewUserAddressModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserAddressModel {
	return &defaultUserAddressModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "user_address",
	}
}

// Insert 新增地址
func (m *defaultUserAddressModel) Insert(ctx context.Context, data *UserAddress) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, receiver, phone, province, city, district, detail, is_default) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", m.table)
	ret, err := m.ExecNoCacheCtx(ctx, query, data.UserId, data.Receiver, data.Phone, data.Province, data.City, data.District, data.Detail, data.IsDefault)
	return ret, err
}

// FindOne 根据 ID 查询地址（带缓存）
func (m *defaultUserAddressModel) FindOne(ctx context.Context, id int64) (*UserAddress, error) {
	key := fmt.Sprintf("%s%v", cacheUserAddressIdPrefix, id)
	var resp UserAddress
	err := m.QueryRowCtx(ctx, &resp, key, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("SELECT id, user_id, receiver, phone, province, city, district, detail, is_default, create_time, update_time FROM %s WHERE id = $1 LIMIT 1", m.table)
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

// FindByUserId 查询用户的所有地址（默认地址优先，无缓存直查）
func (m *defaultUserAddressModel) FindByUserId(ctx context.Context, userId int64) ([]*UserAddress, error) {
	var list []*UserAddress
	query := fmt.Sprintf("SELECT id, user_id, receiver, phone, province, city, district, detail, is_default, create_time, update_time FROM %s WHERE user_id = $1 ORDER BY is_default DESC, id DESC", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, userId)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Update 更新地址
func (m *defaultUserAddressModel) Update(ctx context.Context, data *UserAddress) error {
	key := fmt.Sprintf("%s%v", cacheUserAddressIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("UPDATE %s SET receiver = $1, phone = $2, province = $3, city = $4, district = $5, detail = $6, is_default = $7 WHERE id = $8", m.table)
		return conn.ExecCtx(ctx, query, data.Receiver, data.Phone, data.Province, data.City, data.District, data.Detail, data.IsDefault, data.Id)
	}, key)
	return err
}

// Delete 删除地址
func (m *defaultUserAddressModel) Delete(ctx context.Context, id int64) error {
	key := fmt.Sprintf("%s%v", cacheUserAddressIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, key)
	return err
}
