// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"time"

	"go-mail/common/utils"
	"go-mail/common/xerr"
	"go-mail/services/auth/api/internal/svc"
	"go-mail/services/auth/api/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Login 处理用户登录
// 1. 根据用户名查询用户
// 2. bcrypt 校验密码
// 3. 签发 JWT Token
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 查询用户
	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil {
		logx.Errorf("查询用户失败: username=%s, err=%v", req.Username, err)
		return nil, xerr.NewCodeErrorMsg(xerr.Unauthorized, "用户名或密码错误")
	}

	// 校验密码
	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, xerr.NewCodeErrorMsg(xerr.Unauthorized, "用户名或密码错误")
	}

	// 签发 JWT Token
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessSecret := l.svcCtx.Config.Auth.AccessSecret

	claims := jwt.MapClaims{
		"userId": user.Id,
		"exp":    now + accessExpire,
		"iat":    now,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(accessSecret))
	if err != nil {
		logx.Errorf("生成 Token 失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	logx.Infof("用户登录成功: userId=%d, username=%s", user.Id, user.Username)

	return &types.LoginResp{
		Token:        token,
		RefreshToken: "", // 后续实现
		ExpireIn:     accessExpire,
	}, nil
}
