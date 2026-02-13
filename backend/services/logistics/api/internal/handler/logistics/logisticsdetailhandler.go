package logistics

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/common/xerr"
	"go-mail/services/logistics/api/internal/logic/logistics"
	"go-mail/services/logistics/api/internal/svc"
)

// 物流详情（路径参数 :id 即订单号）
func LogisticsDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderId := r.PathValue("id")
		if orderId == "" {
			response.Response(w, nil, xerr.NewCodeErrorMsg(xerr.BadRequest, "无效的订单 ID"))
			return
		}

		l := logistics.NewLogisticsDetailLogic(r.Context(), svcCtx)
		resp, err := l.LogisticsDetail(orderId)
		response.Response(w, resp, err)
	}
}
