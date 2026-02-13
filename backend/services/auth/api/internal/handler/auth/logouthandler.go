// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/services/auth/api/internal/logic/auth"
	"go-mail/services/auth/api/internal/svc"
)

// 注销登录
func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewLogoutLogic(r.Context(), svcCtx)
		err := l.Logout()
		response.Response(w, nil, err)
	}
}
