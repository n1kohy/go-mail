package search

import (
	"context"

	"go-mail/services/search/api/internal/svc"
	"go-mail/services/search/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SearchSuggestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchSuggestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchSuggestLogic {
	return &SearchSuggestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SearchSuggest 搜索联想（返回前 10 个匹配的商品名称）
func (l *SearchSuggestLogic) SearchSuggest(req *types.SearchSuggestReq) (resp *types.SearchSuggestResp, err error) {
	keyword := req.Keyword + "%"

	type nameResult struct {
		Name string `db:"name"`
	}

	var list []nameResult
	conn := sqlx.NewSqlConn("postgres", l.svcCtx.Config.DataSource)
	query := "SELECT name FROM product WHERE status = 1 AND name ILIKE $1 LIMIT 10"
	_ = conn.QueryRowsCtx(l.ctx, &list, query, keyword)

	suggestions := make([]string, 0, len(list))
	for _, item := range list {
		suggestions = append(suggestions, item.Name)
	}

	return &types.SearchSuggestResp{Suggestions: suggestions}, nil
}
