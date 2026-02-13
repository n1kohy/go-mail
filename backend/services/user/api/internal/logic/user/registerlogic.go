// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"go-mail/common/utils"
	"go-mail/common/xerr"
	"go-mail/services/user/api/internal/svc"
	"go-mail/services/user/api/internal/types"
	"go-mail/services/user/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Register 处理用户注册
// 1. 校验用户名唯一
// 2. 校验手机号唯一
// 3. bcrypt 加密密码
// 4. 写入数据库
func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 检查用户名是否已被注册
	_, err = l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err == nil {
		return nil, xerr.NewCodeErrorMsg(xerr.BadRequest, "用户名已被注册")
	}

	// 检查手机号是否已被注册
	_, err = l.svcCtx.UserModel.FindOneByMobile(l.ctx, req.Mobile)
	if err == nil {
		return nil, xerr.NewCodeErrorMsg(xerr.BadRequest, "手机号已被注册")
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logx.Errorf("密码加密失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	// 插入用户
	result, err := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Mobile:   req.Mobile,
	})
	if err != nil {
		logx.Errorf("插入用户失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	userId, err := result.LastInsertId()
	if err != nil {
		logx.Errorf("获取用户ID失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	logx.Infof("用户注册成功: username=%s, userId=%d", req.Username, userId)

	return &types.RegisterResp{
		UserId: userId,
	}, nil
}
