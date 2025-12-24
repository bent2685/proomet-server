package dto

// LoginWithPwdDto 使用密码登陆
type LoginWithPwdDto struct {
	Account  string `json:"account" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// RegisterDto 用户注册
type RegisterDto struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}
