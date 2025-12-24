package handlers

import (
	"proomet/internal/interfaces/validators"
	"proomet/internal/middleware"
	"proomet/pkg/utils/res"

	"github.com/gin-gonic/gin"
)

// Handler 处理器函数类型
type Handler func(c *gin.Context) error

// Handle 统一处理函数，自动处理错误
func Handle(h Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h(c); err != nil {
			res.ErrInternalServer.ThrowMsg(c, err.Error())
		}
	}
}

// BindHandler 绑定请求并处理的函数
func BindHandler(req any, processor func(*gin.Context, any) error) gin.HandlerFunc {
	return middleware.BindAndProcess(req, func(c *gin.Context, boundReq any) any {
		if err := processor(c, boundReq); err != nil {
			return res.Response{
				Code:    res.ErrInternalServer.Code,
				Message: err.Error(),
				Data:    nil,
			}
		}
		return res.Response{
			Code:    20000,
			Message: "操作成功",
			Data:    nil,
		}
	})
}

// Success 返回成功响应
func Success(c *gin.Context, data any) {
	res.Success(c, data)
}

// SuccessWithMessage 返回带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data any) {
	res.SuccessMsg(c, message, data)
}

// Error 返回错误响应并 panic
func Error(c *gin.Context, err error) {
	// 如果是 BusinessError，直接 panic 以保留原始错误码
	if businessErr, ok := err.(*res.BusinessError); ok {
		panic(businessErr)
	}
	// 其他错误统一使用内部服务器错误
	res.ErrInternalServer.ThrowMsg(c, err.Error())
}

// Bind 绑定并验证请求
func Bind(c *gin.Context, req any) error {
	if err := c.ShouldBindJSON(req); err != nil {
		res.ErrInvalidParam.ThrowMsg(c, validators.GetValidationError(err))
		return err
	}
	return nil
}
