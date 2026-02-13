package promotion

import (
	"context"
	"encoding/json"

	"go-mail/common/xerr"
	"go-mail/services/promotion/api/internal/svc"
	"go-mail/services/promotion/api/internal/types"
	"go-mail/services/promotion/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponClaimLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCouponClaimLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponClaimLogic {
	return &CouponClaimLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CouponClaim 领取优惠券
func (l *CouponClaimLogic) CouponClaim(req *types.CouponClaimReq) (resp *types.CouponClaimResp, err error) {
	userId, err := l.getUserId()
	if err != nil {
		return nil, err
	}

	// 检查优惠券是否存在
	coupon, err := l.svcCtx.CouponModel.FindOne(l.ctx, req.CouponId)
	if err != nil {
		return nil, xerr.NewCodeErrorMsg(xerr.NotFound, "优惠券不存在")
	}

	// 检查是否已领取
	_, err = l.svcCtx.CouponRecordModel.FindByUserIdAndCouponId(l.ctx, userId, coupon.Id)
	if err == nil {
		return nil, xerr.NewCodeErrorMsg(xerr.BadRequest, "已领取过该优惠券")
	}

	// 扣减库存
	err = l.svcCtx.CouponModel.DecrRemainCount(l.ctx, coupon.Id)
	if err != nil {
		return nil, xerr.NewCodeErrorMsg(xerr.BadRequest, "优惠券已被领完")
	}

	// 创建领取记录
	_, err = l.svcCtx.CouponRecordModel.Insert(l.ctx, &model.CouponRecord{
		UserId:   userId,
		CouponId: coupon.Id,
		Status:   0,
	})
	if err != nil {
		logx.Errorf("创建领券记录失败: userId=%d, couponId=%d, err=%v", userId, coupon.Id, err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	return &types.CouponClaimResp{Msg: "领取成功"}, nil
}

func (l *CouponClaimLogic) getUserId() (int64, error) {
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
