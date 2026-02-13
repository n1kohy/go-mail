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

type GetUserByMobileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByMobileLogic {
	return &GetUserByMobileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserByMobile 根据手机号查询用户
func (l *GetUserByMobileLogic) GetUserByMobile(in *user.GetUserByMobileReq) (*user.GetUserResp, error) {
	u, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		logx.Errorf("根据手机号查询用户失败: mobile=%s, err=%v", in.Mobile, err)
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
