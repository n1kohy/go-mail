package logic

import (
	"context"

	"go-mail/services/promotion/rpc/internal/svc"
	"go-mail/services/promotion/rpc/promotion"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CalculateDiscountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCalculateDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CalculateDiscountLogic {
	return &CalculateDiscountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CalculateDiscount 计算优惠价格
func (l *CalculateDiscountLogic) CalculateDiscount(in *promotion.CalculateDiscountReq) (*promotion.CalculateDiscountResp, error) {
	if in.CouponId == 0 {
		// 无优惠券
		return &promotion.CalculateDiscountResp{
			DiscountAmount: 0,
			PayAmount:      in.TotalAmount,
			Valid:          true,
		}, nil
	}

	coupon, err := l.svcCtx.CouponModel.FindOne(l.ctx, in.CouponId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "优惠券不存在")
	}

	// 校验门槛
	if in.TotalAmount < coupon.Threshold {
		return &promotion.CalculateDiscountResp{
			DiscountAmount: 0,
			PayAmount:      in.TotalAmount,
			Valid:          false,
		}, nil
	}

	// 计算优惠金额
	discountAmount := coupon.Discount
	payAmount := in.TotalAmount - discountAmount
	if payAmount < 0 {
		payAmount = 0
	}

	return &promotion.CalculateDiscountResp{
		DiscountAmount: discountAmount,
		PayAmount:      payAmount,
		Valid:          true,
	}, nil
}
