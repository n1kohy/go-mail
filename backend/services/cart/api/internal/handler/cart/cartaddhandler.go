package cart

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/services/cart/api/internal/logic/cart"
	"go-mail/services/cart/api/internal/svc"
	"go-mail/services/cart/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CartAddHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CartAddReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := cart.NewCartAddLogic(r.Context(), svcCtx)
		resp, err := l.CartAdd(&req)
		response.Response(w, resp, err)
	}
}
