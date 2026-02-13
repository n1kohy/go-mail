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

type GetAddressListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取地址列表
func NewGetAddressListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAddressListLogic {
	return &GetAddressListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetAddressList 获取当前用户的所有收货地址
func (l *GetAddressListLogic) GetAddressList() (resp *types.AddressListResp, err error) {
	// 从 JWT context 提取 userId
	userId, ok := l.ctx.Value("userId").(json.Number)
	if !ok {
		return nil, xerr.NewCodeError(xerr.Unauthorized)
	}
	uid, err := userId.Int64()
	if err != nil {
		return nil, xerr.NewCodeError(xerr.Unauthorized)
	}

	// 查询该用户所有地址
	addresses, err := l.svcCtx.UserAddressModel.FindByUserId(l.ctx, uid)
	if err != nil {
		logx.Errorf("查询地址列表失败: userId=%d, err=%v", uid, err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	// 转换为响应类型
	var list []types.AddressItem
	for _, addr := range addresses {
		isDefault := 0
		if addr.IsDefault {
			isDefault = 1
		}
		list = append(list, types.AddressItem{
			Id:        addr.Id,
			Receiver:  addr.Receiver,
			Phone:     addr.Phone,
			Province:  addr.Province,
			City:      addr.City,
			District:  addr.District,
			Detail:    addr.Detail,
			IsDefault: isDefault,
		})
	}

	return &types.AddressListResp{List: list}, nil
}
