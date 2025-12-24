package res

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 通用响应结构
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// BusinessError 业务异常结构
type BusinessError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *BusinessError) Error() string {
	return e.Message
}

// ============================================================
// 预设业务异常 - 支持直接抛出、自定义消息、链式调用
// ============================================================

// Throw 直接抛出预设异常（通过 panic，由 recovery 中间件统一处理）
func (e *BusinessError) Throw(c *gin.Context) {
	panic(e)
}

// Msg 返回带自定义消息的异常副本（不抛出）
func (e *BusinessError) Msg(message string) *BusinessError {
	return &BusinessError{
		Code:    e.Code,
		Message: message,
	}
}

// ThrowMsg 抛出带自定义消息的异常（通过 panic，由 recovery 中间件统一处理）
func (e *BusinessError) ThrowMsg(c *gin.Context, message string) {
	panic(&BusinessError{
		Code:    e.Code,
		Message: message,
	})
}

// Msgf 返回带格式化消息的异常副本（不抛出）
func (e *BusinessError) Msgf(format string, args ...any) *BusinessError {
	return &BusinessError{
		Code:    e.Code,
		Message: fmt.Sprintf(format, args...),
	}
}

// ThrowMsgf 抛出带格式化消息的异常（通过 panic，由 recovery 中间件统一处理）
func (e *BusinessError) ThrowMsgf(c *gin.Context, format string, args ...any) {
	panic(&BusinessError{
		Code:    e.Code,
		Message: fmt.Sprintf(format, args...),
	})
}

// ============================================================
// 自定义业务异常创建器
// ============================================================

// Err 创建自定义业务异常
func Err(code int, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

// Errf 创建带格式化消息的自定义业务异常
func Errf(code int, format string, args ...any) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Throw 直接抛出自定义异常（通过 panic，由 recovery 中间件统一处理）
func Throw(c *gin.Context, code int, message string) {
	panic(&BusinessError{
		Code:    code,
		Message: message,
	})
}

// Throwf 直接抛出带格式化消息的自定义异常（通过 panic，由 recovery 中间件统一处理）
func Throwf(c *gin.Context, code int, format string, args ...any) {
	panic(&BusinessError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	})
}

// ============================================================
// 成功响应
// ============================================================

// Success 成功响应
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessMsg 带自定义消息的成功响应
func SuccessMsg(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// SuccessMsgf 带格式化消息的成功响应
func SuccessMsgf(c *gin.Context, data any, format string, args ...any) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: fmt.Sprintf(format, args...),
		Data:    data,
	})
}

// ============================================================
// 预设业务异常常量
// ============================================================

var (
	// 通用错误
	ErrClientBussiness = &BusinessError{Code: 400000, Message: "客户端业务错误"}
	ErrInvalidParam    = &BusinessError{Code: 400010, Message: "参数错误"}
	ErrUnauthorized    = &BusinessError{Code: 400001, Message: "未授权"}
	ErrForbidden       = &BusinessError{Code: 400003, Message: "禁止访问"}
	ErrNotFound        = &BusinessError{Code: 400004, Message: "资源不存在"}
	ErrInternalServer  = &BusinessError{Code: 500001, Message: "服务器内部错误"}

	// 认证相关错误
	ErrInvalidCredentials = &BusinessError{Code: 400002, Message: "凭证错误"}

	// 用户相关错误
	ErrUserNotFound      = &BusinessError{Code: 400101, Message: "用户不存在"}
	ErrDataAlreadyExists = &BusinessError{Code: 400102, Message: "数据已存在"}
	ErrInvalidPassword   = &BusinessError{Code: 400103, Message: "密码错误"}
	ErrEmailAlreadyUsed  = &BusinessError{Code: 400104, Message: "邮箱已被使用"}
	ErrUsernameTaken     = &BusinessError{Code: 400105, Message: "用户名已存在"}

	// 权限相关错误
	ErrInsufficientPermissions = &BusinessError{Code: 400009, Message: "权限不足"}
)
