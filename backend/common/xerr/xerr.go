// Package xerr 定义统一的业务错误码和错误类型
package xerr

import "fmt"

// 业务错误码定义
const (
	// 通用错误码
	OK            uint32 = 200  // 成功
	BadRequest    uint32 = 400  // 参数错误
	Unauthorized  uint32 = 401  // 未授权
	Forbidden     uint32 = 403  // 权限不足
	NotFound      uint32 = 404  // 资源不存在
	ServerError   uint32 = 500  // 服务器内部错误
	BusinessError uint32 = 1001 // 通用业务错误

	// 秒杀专用错误码
	OutOfStock uint32 = 2001 // 库存不足
	SystemBusy uint32 = 2002 // 排队中
)

// 错误码描述映射
var codeMsg = map[uint32]string{
	OK:            "成功",
	BadRequest:    "参数错误",
	Unauthorized:  "未授权",
	Forbidden:     "权限不足",
	NotFound:      "资源不存在",
	ServerError:   "服务器内部错误",
	BusinessError: "业务错误",
	OutOfStock:    "库存不足",
	SystemBusy:    "系统繁忙，请稍后再试",
}

// CodeError 自定义业务错误类型
type CodeError struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

// NewCodeError 创建一个带错误码的错误
func NewCodeError(code uint32) *CodeError {
	return &CodeError{
		Code: code,
		Msg:  codeMsg[code],
	}
}

// NewCodeErrorMsg 创建一个自定义消息的错误
func NewCodeErrorMsg(code uint32, msg string) *CodeError {
	return &CodeError{
		Code: code,
		Msg:  msg,
	}
}

// NewBusinessError 快捷创建业务错误
func NewBusinessError(msg string) *CodeError {
	return &CodeError{
		Code: BusinessError,
		Msg:  msg,
	}
}

// Error 实现 error 接口
func (e *CodeError) Error() string {
	return fmt.Sprintf("错误码: %d, 消息: %s", e.Code, e.Msg)
}

// Data 返回错误的 JSON 响应结构
func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: e.Code,
		Msg:  e.Msg,
	}
}

// CodeErrorResponse 错误响应结构体
type CodeErrorResponse struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}
