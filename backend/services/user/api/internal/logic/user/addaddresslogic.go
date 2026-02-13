// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"encoding/json"

	"go-mail/common/xerr"
	"go-mail/services/user/api/internal/svc"
	"go-mail/services/user/api/internal/types"
	"go-mail/services/user/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增地址
func NewAddAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAddressLogic {
	return &AddAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AddAddress 新增收货地址
func (l *AddAddressLogic) AddAddress(req *types.AddAddressReq) (resp *types.AddAddressResp, err error) {
	// 从 JWT context 提取 userId
	userId, ok := l.ctx.Value("userId").(json.Number)
	if !ok {
		return nil, xerr.NewCodeError(xerr.Unauthorized)
	}
	uid, err := userId.Int64()
	if err != nil {
		return nil, xerr.NewCodeError(xerr.Unauthorized)
	}

	// 插入地址
	result, err := l.svcCtx.UserAddressModel.Insert(l.ctx, &model.UserAddress{
		UserId:    uid,
		Receiver:  req.Receiver,
		Phone:     req.Phone,
		Province:  req.Province,
		City:      req.City,
		District:  req.District,
		Detail:    req.Detail,
		IsDefault: req.IsDefault == 1,
	})
	if err != nil {
		logx.Errorf("新增地址失败: userId=%d, err=%v", uid, err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	addressId, err := result.LastInsertId()
	if err != nil {
		logx.Errorf("获取地址ID失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	return &types.AddAddressResp{Id: addressId}, nil
}
