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

var cacheUserIdPrefix = "cache:user:id:"

type (
	// UserModel 自定义接口，可在此扩展方法
	UserModel interface {
		Insert(ctx context.Context, data *User) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*User, error)
		FindOneByUsername(ctx context.Context, username string) (*User, error)
		FindOneByMobile(ctx context.Context, mobile string) (*User, error)
		Update(ctx context.Context, data *User) error
		Delete(ctx context.Context, id int64) error
	}

	// User 对应数据库 "user" 表
	User struct {
		Id          int64     `db:"id"`
		Username    string    `db:"username"`
		Password    string    `db:"password"`
		Mobile      string    `db:"mobile"`
		Avatar      string    `db:"avatar"`
		Gender      int64     `db:"gender"`
		Role        int64     `db:"role"`
		MemberLevel int64     `db:"member_level"`
		CreateTime  time.Time `db:"create_time"`
		UpdateTime  time.Time `db:"update_time"`
	}

	defaultUserModel struct {
		sqlc.CachedConn
		table string
	}
)

// NewUserModel 创建用户 Model 实例
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &defaultUserModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      `"user"`,
	}
}

// Insert 插入用户
func (m *defaultUserModel) Insert(ctx context.Context, data *User) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (username, password, mobile, avatar, gender, role, member_level) VALUES ($1, $2, $3, $4, $5, $6, $7)", m.table)
	ret, err := m.ExecNoCacheCtx(ctx, query, data.Username, data.Password, data.Mobile, data.Avatar, data.Gender, data.Role, data.MemberLevel)
	return ret, err
}

// FindOne 根据 ID 查询用户（带缓存）
func (m *defaultUserModel) FindOne(ctx context.Context, id int64) (*User, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	var resp User
	err := m.QueryRowCtx(ctx, &resp, userIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("SELECT id, username, password, mobile, avatar, gender, role, member_level, create_time, update_time FROM %s WHERE id = $1 LIMIT 1", m.table)
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

// FindOneByUsername 根据用户名查询（无缓存直查）
func (m *defaultUserModel) FindOneByUsername(ctx context.Context, username string) (*User, error) {
	var resp User
	query := fmt.Sprintf("SELECT id, username, password, mobile, avatar, gender, role, member_level, create_time, update_time FROM %s WHERE username = $1 LIMIT 1", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, username)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindOneByMobile 根据手机号查询（无缓存直查）
func (m *defaultUserModel) FindOneByMobile(ctx context.Context, mobile string) (*User, error) {
	var resp User
	query := fmt.Sprintf("SELECT id, username, password, mobile, avatar, gender, role, member_level, create_time, update_time FROM %s WHERE mobile = $1 LIMIT 1", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, mobile)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Update 更新用户
func (m *defaultUserModel) Update(ctx context.Context, data *User) error {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("UPDATE %s SET username = $1, password = $2, mobile = $3, avatar = $4, gender = $5, role = $6, member_level = $7 WHERE id = $8", m.table)
		return conn.ExecCtx(ctx, query, data.Username, data.Password, data.Mobile, data.Avatar, data.Gender, data.Role, data.MemberLevel, data.Id)
	}, userIdKey)
	return err
}

// Delete 删除用户
func (m *defaultUserModel) Delete(ctx context.Context, id int64) error {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, userIdKey)
	return err
}
