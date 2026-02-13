// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"go-mail/services/user/model"
	"go-mail/services/user/rpc/internal/svc"
	"go-mail/services/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUserByUsernameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByUsernameLogic {
	return &GetUserByUsernameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserByUsername 根据用户名查询用户
func (l *GetUserByUsernameLogic) GetUserByUsername(in *user.GetUserByUsernameReq) (*user.GetUserResp, error) {
	u, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		logx.Errorf("根据用户名查询用户失败: username=%s, err=%v", in.Username, err)
		return nil, status.Error(codes.Internal, "查询用户失败")
	}

	return &user.GetUserResp{
		Id:          u.Id,
		Username:    u.Username,
		Password:    u.Password,
		Mobile:      u.Mobile,
		Avatar:      u.Avatar,
		Gender:      int32(u.Gender),
		Role:        int32(u.Role),
		MemberLevel: int32(u.MemberLevel),
	}, nil
}
