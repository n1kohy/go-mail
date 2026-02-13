package seckill

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/services/seckill/api/internal/logic/seckill"
	"go-mail/services/seckill/api/internal/svc"
	"go-mail/services/seckill/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SeckillActionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SeckillActionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := seckill.NewSeckillActionLogic(r.Context(), svcCtx)
		resp, err := l.SeckillAction(&req)
		response.Response(w, resp, err)
	}
}
