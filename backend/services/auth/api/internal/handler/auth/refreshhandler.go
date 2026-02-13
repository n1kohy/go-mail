// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/services/auth/api/internal/logic/auth"
	"go-mail/services/auth/api/internal/svc"
	"go-mail/services/auth/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 刷新 Token
func RefreshHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := auth.NewRefreshLogic(r.Context(), svcCtx)
		resp, err := l.Refresh(&req)
		response.Response(w, resp, err)
	}
}
