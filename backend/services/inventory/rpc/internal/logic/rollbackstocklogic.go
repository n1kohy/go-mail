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

type RollbackStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRollbackStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RollbackStockLogic {
	return &RollbackStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RollbackStock 回滚库存（取消订单时调用）
func (l *RollbackStockLogic) RollbackStock(in *inventory.RollbackStockReq) (*inventory.RollbackStockResp, error) {
	err := l.svcCtx.StockModel.RollbackStock(l.ctx, in.SkuId, in.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			logx.Errorf("回滚库存失败（锁定量不足）: skuId=%d, qty=%d", in.SkuId, in.Quantity)
			return nil, status.Error(codes.FailedPrecondition, "回滚库存失败")
		}
		logx.Errorf("回滚库存失败: skuId=%d, err=%v", in.SkuId, err)
		return nil, status.Error(codes.Internal, "回滚库存失败")
	}

	logx.Infof("库存回滚成功: skuId=%d, qty=%d", in.SkuId, in.Quantity)
	return &inventory.RollbackStockResp{Success: true}, nil
}
