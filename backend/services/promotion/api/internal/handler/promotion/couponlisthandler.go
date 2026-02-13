package promotion

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/services/promotion/api/internal/logic/promotion"
	"go-mail/services/promotion/api/internal/svc"
)

func CouponListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := promotion.NewCouponListLogic(r.Context(), svcCtx)
		resp, err := l.CouponList()
		response.Response(w, resp, err)
	}
}
