package promotion

import (
	"context"

	"go-mail/common/xerr"
	"go-mail/services/promotion/api/internal/svc"
	"go-mail/services/promotion/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCouponListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponListLogic {
	return &CouponListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CouponList 可领优惠券列表
func (l *CouponListLogic) CouponList() (resp *types.CouponListResp, err error) {
	coupons, err := l.svcCtx.CouponModel.FindAvailable(l.ctx)
	if err != nil {
		logx.Errorf("查询优惠券失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	var items []types.CouponItem
	for _, c := range coupons {
		items = append(items, types.CouponItem{
			Id:        c.Id,
			Name:      c.Name,
			Type:      int(c.Type),
			Threshold: c.Threshold,
			Discount:  c.Discount,
			EndTime:   c.EndTime.Format("2006-01-02"),
		})
	}

	if items == nil {
		items = []types.CouponItem{}
	}

	return &types.CouponListResp{List: items}, nil
}
