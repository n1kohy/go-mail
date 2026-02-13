package logic

import (
	"context"

	"go-mail/services/inventory/model"
	"go-mail/services/inventory/rpc/internal/svc"
	"go-mail/services/inventory/rpc/inventory"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStockLogic {
	return &GetStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetStock 查询 SKU 库存
func (l *GetStockLogic) GetStock(in *inventory.GetStockReq) (*inventory.GetStockResp, error) {
	stock, err := l.svcCtx.StockModel.FindOneBySkuId(l.ctx, in.SkuId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(codes.NotFound, "库存记录不存在")
		}
		logx.Errorf("查询库存失败: skuId=%d, err=%v", in.SkuId, err)
		return nil, status.Error(codes.Internal, "查询库存失败")
	}

	return &inventory.GetStockResp{
		SkuId:     stock.SkuId,
		Total:     stock.Total,
		Available: stock.Available,
		Locked:    stock.Locked,
	}, nil
}
