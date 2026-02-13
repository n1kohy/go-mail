package logic

import (
	"context"
	"database/sql"

	"go-mail/services/inventory/rpc/internal/svc"
	"go-mail/services/inventory/rpc/inventory"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeductStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductStockLogic {
	return &DeductStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeductStock 扣减库存（乐观锁）
func (l *DeductStockLogic) DeductStock(in *inventory.DeductStockReq) (*inventory.DeductStockResp, error) {
	err := l.svcCtx.StockModel.DeductStock(l.ctx, in.SkuId, in.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			logx.Errorf("库存不足或版本冲突: skuId=%d, qty=%d", in.SkuId, in.Quantity)
			return nil, status.Error(codes.FailedPrecondition, "库存不足")
		}
		logx.Errorf("扣减库存失败: skuId=%d, err=%v", in.SkuId, err)
		return nil, status.Error(codes.Internal, "扣减库存失败")
	}

	logx.Infof("库存扣减成功: skuId=%d, qty=%d", in.SkuId, in.Quantity)
	return &inventory.DeductStockResp{Success: true}, nil
}
