// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"go-mail/services/auth/api/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 注销登录
func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Logout 注销登录
// 将当前 token 加入黑名单（需要 Redis 支持）
func (l *LogoutLogic) Logout() error {
	// TODO: 将当前 token 加入 Redis 黑名单
	logx.Info("用户注销成功")
	return nil
}
