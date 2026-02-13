// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"encoding/json"

	"go-mail/common/xerr"
	"go-mail/services/user/api/internal/svc"
	"go-mail/services/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户信息
func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetUserInfo 获取当前登录用户信息
// 从 JWT context 中解析 userId，查库并返回（手机号脱敏）
func (l *GetUserInfoLogic) GetUserInfo() (resp *types.UserInfoResp, err error) {
	// 从 JWT context 提取 userId
	userId, ok := l.ctx.Value("userId").(json.Number)
	if !ok {
		return nil, xerr.NewCodeError(xerr.Unauthorized)
	}
	uid, err := userId.Int64()
	if err != nil {
		return nil, xerr.NewCodeError(xerr.Unauthorized)
	}

	// 查询用户
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, uid)
	if err != nil {
		logx.Errorf("查询用户失败: userId=%d, err=%v", uid, err)
		return nil, xerr.NewCodeError(xerr.NotFound)
	}

	// 手机号脱敏: 138****8000
	mobileMasked := user.Mobile
	if len(mobileMasked) >= 11 {
		mobileMasked = mobileMasked[:3] + "****" + mobileMasked[7:]
	}

	return &types.UserInfoResp{
		Id:          user.Id,
		Username:    user.Username,
		Mobile:      mobileMasked,
		Avatar:      user.Avatar,
		Gender:      int(user.Gender),
		MemberLevel: int(user.MemberLevel),
	}, nil
}
