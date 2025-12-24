package routes

import (
	"proomet/internal/middleware"
	"proomet/pkg/utils/res"

	"github.com/gin-gonic/gin"
)

type TestRouter struct{}

// NewUserRouter 创建用户路由实例
func NewTestRouter() *TestRouter {
	return &TestRouter{}
}

// RegisterRoutes 注册路由
func (tr *TestRouter) RegisterRoutes(router *gin.RouterGroup) {
	testGroup := router.Group("/test")
	{
		testGroup.GET("/health",
			middleware.Authenticate(),
			middleware.Authorize(),
			tr.Health)
	}
}

// PING godoc
// @Summary 测试路由
// @Tags 测试
// @Accept json
// @Produce json
// @Success 200 {object} res.Response{data=bool} "成功"
// @Router /test/health [get]
func (tr *TestRouter) Health(c *gin.Context) {
	res.SuccessMsg(c, "PONG", true)
}
