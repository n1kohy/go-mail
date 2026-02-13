package promotion

import (
	"context"
	"encoding/json"

	"go-mail/common/xerr"
	"go-mail/services/promotion/api/internal/svc"
	"go-mail/services/promotion/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponMineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCouponMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponMineLogic {
	return &CouponMineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CouponMine 我的优惠券
func (l *CouponMineLogic) CouponMine(req *types.CouponMineReq) (resp *types.CouponMineResp, err error) {
	userId, err := l.getUserId()
	if err != nil {
		return nil, err
	}

	records, err := l.svcCtx.CouponRecordModel.FindByUserId(l.ctx, userId, req.Status)
	if err != nil {
		logx.Errorf("查询我的优惠券失败: userId=%d, err=%v", userId, err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	var items []types.MyCouponItem
	for _, r := range records {
		// 查优惠券详情
		coupon, err := l.svcCtx.CouponModel.FindOne(l.ctx, r.CouponId)
		if err != nil {
			continue
		}
		items = append(items, types.MyCouponItem{
			Id:        r.Id,
			CouponId:  r.CouponId,
			Name:      coupon.Name,
			Threshold: coupon.Threshold,
			Discount:  coupon.Discount,
			Status:    int(r.Status),
			EndTime:   coupon.EndTime.Format("2006-01-02"),
		})
	}

	if items == nil {
		items = []types.MyCouponItem{}
	}

	return &types.CouponMineResp{List: items}, nil
}

func (l *CouponMineLogic) getUserId() (int64, error) {
	uid := l.ctx.Value("userId")
	if uid == nil {
		return 0, xerr.NewCodeError(xerr.Unauthorized)
	}
	userId, err := uid.(json.Number).Int64()
	if err != nil {
		return 0, xerr.NewCodeError(xerr.Unauthorized)
	}
	return userId, nil
}
