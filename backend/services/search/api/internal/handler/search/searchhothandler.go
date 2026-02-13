package search

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/services/search/api/internal/logic/search"
	"go-mail/services/search/api/internal/svc"
)

func SearchHotHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := search.NewSearchHotLogic(r.Context(), svcCtx)
		resp, err := l.SearchHot()
		response.Response(w, resp, err)
	}
}
