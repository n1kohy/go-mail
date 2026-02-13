package logic

import (
	"context"

	"go-mail/services/promotion/rpc/internal/svc"
	"go-mail/services/promotion/rpc/promotion"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UseCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUseCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UseCouponLogic {
	return &UseCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UseCoupon 核销优惠券
func (l *UseCouponLogic) UseCoupon(in *promotion.UseCouponReq) (*promotion.UseCouponResp, error) {
	err := l.svcCtx.CouponRecordModel.MarkUsed(l.ctx, in.UserId, in.CouponId, in.UsedOrderId)
	if err != nil {
		logx.Errorf("核销优惠券失败: userId=%d, couponId=%d, err=%v", in.UserId, in.CouponId, err)
		return nil, status.Error(codes.Internal, "核销优惠券失败")
	}

	return &promotion.UseCouponResp{Success: true}, nil
}
