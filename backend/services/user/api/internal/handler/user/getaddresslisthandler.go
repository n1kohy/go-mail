// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/services/user/api/internal/logic/user"
	"go-mail/services/user/api/internal/svc"
)

// 获取地址列表
func GetAddressListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetAddressListLogic(r.Context(), svcCtx)
		resp, err := l.GetAddressList()
		response.Response(w, resp, err)
	}
}
