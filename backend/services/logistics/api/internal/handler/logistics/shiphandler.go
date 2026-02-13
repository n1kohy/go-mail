package logistics

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/services/logistics/api/internal/logic/logistics"
	"go-mail/services/logistics/api/internal/svc"
	"go-mail/services/logistics/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ShipHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShipReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logistics.NewShipLogic(r.Context(), svcCtx)
		resp, err := l.Ship(&req)
		response.Response(w, resp, err)
	}
}
