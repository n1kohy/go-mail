package search

import (
	"context"
	"fmt"

	"go-mail/common/xerr"
	"go-mail/services/search/api/internal/svc"
	"go-mail/services/search/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SearchProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchProductLogic {
	return &SearchProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SearchProduct 商品搜索（初版：使用 PostgreSQL LIKE 模糊查询）
func (l *SearchProductLogic) SearchProduct(req *types.SearchProductReq) (resp *types.SearchProductResp, err error) {
	// 初版实现：PostgreSQL LIKE 搜索，后续可替换为 Elasticsearch
	keyword := "%" + req.Keyword + "%"
	offset := (req.Page - 1) * req.Size

	type searchResult struct {
		Id        int64   `db:"id"`
		Name      string  `db:"name"`
		Price     float64 `db:"price"`
		MainImage string  `db:"main_image"`
	}

	var list []searchResult
	conn := sqlx.NewSqlConn("postgres", l.svcCtx.Config.DataSource)
	query := "SELECT p.id, p.name, COALESCE(s.price, 0) AS price, p.main_image FROM product p LEFT JOIN product_sku s ON s.product_id = p.id AND s.id = (SELECT MIN(id) FROM product_sku WHERE product_id = p.id) WHERE p.status = 1 AND p.name ILIKE $1 ORDER BY p.id DESC LIMIT $2 OFFSET $3"
	err = conn.QueryRowsCtx(l.ctx, &list, query, keyword, req.Size, offset)
	if err != nil {
		logx.Errorf("搜索商品失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ServerError)
	}

	// 统计总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM product WHERE status = 1 AND name ILIKE $1"
	_ = conn.QueryRowCtx(l.ctx, &total, countQuery, keyword)

	var items []types.SearchProductItem
	for _, p := range list {
		items = append(items, types.SearchProductItem{
			Id:    p.Id,
			Name:  p.Name,
			Price: p.Price,
			Image: p.MainImage,
		})
	}

	// 防止 nil slice
	if items == nil {
		items = []types.SearchProductItem{}
	}

	return &types.SearchProductResp{List: items, Total: total}, nil
}

// placeholder 避免 fmt 未使用
var _ = fmt.Sprintf
