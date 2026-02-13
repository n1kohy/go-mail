// Package response 提供统一的 HTTP 响应格式
// 所有 API 服务的 handler 层应调用此包返回响应
package response

import (
	"net/http"

	"go-mail/common/xerr"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Body 统一响应体结构
type Body struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Response 返回成功响应
func Response(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		// 判断是否为自定义业务错误
		codeErr, ok := err.(*xerr.CodeError)
		if ok {
			httpx.WriteJson(w, http.StatusOK, &Body{
				Code: codeErr.Code,
				Msg:  codeErr.Msg,
			})
		} else {
			// 未知错误 → 500
			httpx.WriteJson(w, http.StatusOK, &Body{
				Code: xerr.ServerError,
				Msg:  "服务器内部错误",
			})
		}
		return
	}

	// 成功响应
	httpx.WriteJson(w, http.StatusOK, &Body{
		Code: xerr.OK,
		Msg:  "成功",
		Data: data,
	})
}
