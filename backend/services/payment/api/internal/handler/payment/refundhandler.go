package payment

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/services/payment/api/internal/logic/payment"
	"go-mail/services/payment/api/internal/svc"
	"go-mail/services/payment/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RefundHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefundReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := payment.NewRefundLogic(r.Context(), svcCtx)
		resp, err := l.Refund(&req)
		response.Response(w, resp, err)
	}
}
