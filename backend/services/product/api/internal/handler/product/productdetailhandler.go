package product

import (
	"net/http"
	"strconv"

	"go-mail/common/response"
	"go-mail/common/xerr"
	"go-mail/services/product/api/internal/logic/product"
	"go-mail/services/product/api/internal/svc"
)

// 商品详情（SPU+SKU 联查）
func ProductDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.Response(w, nil, xerr.NewCodeErrorMsg(xerr.BadRequest, "无效的商品 ID"))
			return
		}

		l := product.NewProductDetailLogic(r.Context(), svcCtx)
		resp, err := l.ProductDetail(id)
		response.Response(w, resp, err)
	}
}
