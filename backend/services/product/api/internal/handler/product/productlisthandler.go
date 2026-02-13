package product

import (
	"net/http"

	"go-mail/common/response"
	"go-mail/services/product/api/internal/logic/product"
	"go-mail/services/product/api/internal/svc"
	"go-mail/services/product/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 商品列表（支持分类筛选+分页）
func ProductListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProductListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := product.NewProductListLogic(r.Context(), svcCtx)
		resp, err := l.ProductList(&req)
		response.Response(w, resp, err)
	}
}
