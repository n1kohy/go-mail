package search

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/services/search/api/internal/logic/search"
	"go-mail/services/search/api/internal/svc"
	"go-mail/services/search/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SearchSuggestHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchSuggestReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := search.NewSearchSuggestLogic(r.Context(), svcCtx)
		resp, err := l.SearchSuggest(&req)
		response.Response(w, resp, err)
	}
}
