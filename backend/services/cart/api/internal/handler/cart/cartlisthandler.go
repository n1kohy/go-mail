package cart

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/services/cart/api/internal/logic/cart"
	"go-mail/services/cart/api/internal/svc"
)

func CartListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := cart.NewCartListLogic(r.Context(), svcCtx)
		resp, err := l.CartList()
		response.Response(w, resp, err)
	}
}
