package vo

// LoginVO 用户登陆
type AuthLoginVO struct {
	Token string `json:"token"`
}

// RegisterVO 用户注册
type AuthRegisterVO struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
