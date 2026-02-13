package order

import (
	"context"
	"encoding/json"

	"go-mail/common/xerr"
	"go-mail/services/order/api/internal/svc"
	"go-mail/services/order/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderListLogic {
	return &OrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// OrderList 订单列表
func (l *OrderListLogic) OrderList(req *types.OrderListReq) (resp *types.OrderListResp, err error) {
	userId, err := l.getUserId()
	if err != nil {
		return nil, err
	}

	orders, err := l.svcCtx.OrderMasterModel.FindByUserId(l.ctx, userId, req.Status, req.Page, req.Size)
	if err != nil {
		logx.Errorf("查询订单列表失败: userId=%d, err=%v", userId, err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	total, _ := l.svcCtx.OrderMasterModel.CountByUserId(l.ctx, userId, req.Status)

	var items []types.OrderListItem
	for _, o := range orders {
		// 统计订单商品数
		itemCount, _ := l.svcCtx.OrderItemModel.CountByOrderId(l.ctx, o.Id)

		items = append(items, types.OrderListItem{
			OrderId:    o.Id,
			Status:     int(o.Status),
			PayAmount:  o.PayAmount,
			ItemCount:  int(itemCount),
			CreateTime: o.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	if items == nil {
		items = []types.OrderListItem{}
	}

	return &types.OrderListResp{List: items, Total: total}, nil
}

func (l *OrderListLogic) getUserId() (int64, error) {
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
