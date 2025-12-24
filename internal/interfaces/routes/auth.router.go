package routes

import (
	"proomet/internal/interfaces/handlers"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	authHandler handlers.AuthHandler
}

// NewUserRouter 创建用户路由实例
func NewAuthRouter() *AuthRouter {
	return &AuthRouter{
		authHandler: *handlers.NewAuthHandler(),
	}
}

// RegisterRoutes 注册路由
func (tr *AuthRouter) RegisterRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	{
		signGroup := authGroup.Group("/sign")
		{
			signGroup.POST("/with-pwd",
				tr.authHandler.LoginWithPwd)
			signGroup.POST("/register",
				tr.authHandler.Register)
		}
	}
}
