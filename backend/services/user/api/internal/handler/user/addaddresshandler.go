// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/services/user/api/internal/logic/user"
	"go-mail/services/user/api/internal/svc"
	"go-mail/services/user/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 新增收货地址
func AddAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddAddressReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewAddAddressLogic(r.Context(), svcCtx)
		resp, err := l.AddAddress(&req)
		response.Response(w, resp, err)
	}
}
