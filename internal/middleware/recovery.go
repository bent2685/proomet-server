package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"proomet/pkg/utils"
	"proomet/pkg/utils/res"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RecoveryMiddleware 全局异常拦截中间件
func RecoveryMiddleware() gin.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 处理业务异常
				if businessErr, ok := err.(*res.BusinessError); ok {
					utils.Log.Error(logrus.Fields{
						"type":       "API_RESPONSE",
						"method":     c.Request.Method,
						"path":       c.Request.URL.Path,
						"client_ip":  c.ClientIP(),
						"error_code": businessErr.Code,
						"error_msg":  businessErr.Message,
					})

					c.JSON(http.StatusOK, businessErr)
					c.Abort()
					return
				}

				// 处理系统异常
				stack := string(debug.Stack())
				utils.Log.Error(logrus.Fields{
					"type":      "SYSTEM_PANIC",
					"method":    c.Request.Method,
					"path":      c.Request.URL.Path,
					"client_ip": c.ClientIP(),
					"error":     fmt.Sprintf("%v", err),
					"stack":     stack,
				})

				c.JSON(http.StatusOK, res.Response{
					Code:    res.ErrInternalServer.Code,
					Message: res.ErrInternalServer.Message,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RequestIDMiddleware 请求ID中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
