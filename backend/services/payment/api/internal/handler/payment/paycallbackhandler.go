package payment

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/services/payment/api/internal/logic/payment"
	"go-mail/services/payment/api/internal/svc"
	"go-mail/services/payment/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PayCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PayCallbackReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := payment.NewPayCallbackLogic(r.Context(), svcCtx)
		resp, err := l.PayCallback(&req)
		response.Response(w, resp, err)
	}
}
