package promotion

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/services/promotion/api/internal/logic/promotion"
	"go-mail/services/promotion/api/internal/svc"
	"go-mail/services/promotion/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CouponClaimHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CouponClaimReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := promotion.NewCouponClaimLogic(r.Context(), svcCtx)
		resp, err := l.CouponClaim(&req)
		response.Response(w, resp, err)
	}
}
