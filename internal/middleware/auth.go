package middleware

import (
	"proomet/internal/domain/models"
	"proomet/internal/infra/auth"
	"proomet/pkg/utils/jwt"
	"proomet/pkg/utils/res"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// 1. 无 Header，注入 Guest 身份
		if authHeader == "" {
			c.Set("currentUser", models.JwtUser{Role: "guest", Username: "guest"})
			c.Next()
			return
		}

		// 2. 格式校验
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			res.ErrUnauthorized.Msg("Authorization 格式错误").Throw(c)
			c.Abort()
			return
		}

		// 3. 解析 Token
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			// Token 过期或非法，必须明确报错，而不是降级为 Guest
			res.ErrUnauthorized.Msg("登录已过期，请重新登录").Throw(c)
			c.Abort()
			return
		}

		c.Set("currentUser", claims.JwtUser)
		c.Next()
	}
}

// 权限校验中间件
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取已登录用户的 Role
		userRaw, exists := c.Get("currentUser")
		if !exists {
			res.ErrUnauthorized.Throw(c)
		}
		user := userRaw.(models.JwtUser)
		// 获取请求的资源(Object)和动作(Action)
		obj := c.Request.URL.Path
		act := c.Request.Method
		e := auth.GetEnforcer()

		// 调用 Casbin 进行决策
		ok, err := e.Enforce(user.Role, obj, act)
		if err != nil {
			res.ErrInternalServer.Msg("权限系统发生错误").Throw(c)
		}

		if !ok {
			res.ErrForbidden.Throw(c)
		}

		c.Next()
	}
}
