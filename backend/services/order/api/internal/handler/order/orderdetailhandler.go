package order

import (
	"net/http"
	"strconv"

	"go-mail/common/response"
	"go-mail/common/xerr"
	"go-mail/services/order/api/internal/logic/order"
	"go-mail/services/order/api/internal/svc"
)

// 订单详情（路径参数 :id）
func OrderDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderId := r.PathValue("id")
		if orderId == "" {
			response.Response(w, nil, xerr.NewCodeErrorMsg(xerr.BadRequest, "无效的订单 ID"))
			return
		}

		// 验证订单 ID 格式（以 SN 开头的字符串）
		if len(orderId) < 2 {
			// 尝试作为数字 ID 处理
			_, err := strconv.ParseInt(orderId, 10, 64)
			if err != nil {
				response.Response(w, nil, xerr.NewCodeErrorMsg(xerr.BadRequest, "无效的订单 ID"))
				return
			}
		}

		l := order.NewOrderDetailLogic(r.Context(), svcCtx)
		resp, err := l.OrderDetail(orderId)
		response.Response(w, resp, err)
	}
}
