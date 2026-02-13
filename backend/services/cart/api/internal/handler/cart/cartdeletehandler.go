package cart

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/services/cart/api/internal/logic/cart"
	"go-mail/services/cart/api/internal/svc"
	"go-mail/services/cart/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CartDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CartDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := cart.NewCartDeleteLogic(r.Context(), svcCtx)
		resp, err := l.CartDelete(&req)
		response.Response(w, resp, err)
	}
}
