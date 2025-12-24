package handlers

import (
	"proomet/internal/application/services"
	"proomet/internal/interfaces/dto"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证endpoint
type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.AuthService{},
	}
}

// LoginWithPwd godoc
// @Summary 使用用户名密码登陆
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.LoginWithPwdDto true "登录请求"
// @Success 200 {object} res.Response{data=vo.AuthLoginVO} "登录成功"
// @Router /auth/sign/with-pwd [post]
func (h *AuthHandler) LoginWithPwd(c *gin.Context) {
	var req dto.LoginWithPwdDto
	if err := Bind(c, &req); err != nil {
		return
	}
	vo, err := h.authService.LoginWithPwd(&req)
	if err != nil {
		Error(c, err)
	}
	Success(c, vo)
}

// Register godoc
// @Summary 用户注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.RegisterDto true "注册请求"
// @Success 200 {object} res.Response{data=vo.AuthRegisterVO} "注册成功"
// @Router /auth/sign/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterDto
	if err := Bind(c, &req); err != nil {
		return
	}
	vo, err := h.authService.Register(&req)
	if err != nil {
		Error(c, err)
		return
	}
	Success(c, vo)
}
