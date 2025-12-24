package services

import (
	"proomet/internal/domain/models"
	"proomet/internal/infra/database"
	"proomet/internal/interfaces/dto"
	"proomet/internal/interfaces/vo"
	"proomet/pkg/utils/jwt"
	"proomet/pkg/utils/res"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

// LoginWithPwd 登录
func (s *AuthService) LoginWithPwd(dto *dto.LoginWithPwdDto) (*vo.AuthLoginVO, error) {
	db := database.GetDB()

	pwd := dto.Password
	if pwd == "" {
		return nil, res.ErrInvalidParam.Msg("密码不能为空")
	}
	// 判断email/username 是否存在
	var user models.User
	if err := db.Where("email = ?", dto.Account).Or("username = ?", dto.Account).First(&user).Error; err != nil {
		return nil, res.ErrNotFound.Msg("账号不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(pwd)); err != nil {
		return nil, res.ErrInvalidCredentials.Msg("密码错误")
	}
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, res.ErrInternalServer.Msg("生成Token失败")
	}

	return &vo.AuthLoginVO{
		Token: token,
	}, nil
}

// Register 用户注册
func (s *AuthService) Register(dto *dto.RegisterDto) (*vo.AuthRegisterVO, error) {
	db := database.GetDB()

	// 检查用户名是否已存在
	var existingUser models.User
	if err := db.Where("username = ?", dto.Username).First(&existingUser).Error; err == nil {
		return nil, res.ErrUsernameTaken
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, res.ErrInternalServer.Msg("密码加密失败")
	}

	// 创建用户
	user := models.User{
		Username:     dto.Username,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleMember,
		Status:       1,
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, res.ErrInternalServer.Msg("创建用户失败")
	}

	// 生成Token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, res.ErrInternalServer.Msg("生成Token失败")
	}

	return &vo.AuthRegisterVO{
		UserID:   user.ID,
		Username: user.Username,
		Token:    token,
	}, nil
}
