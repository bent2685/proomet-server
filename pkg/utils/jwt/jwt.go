package jwt

import (
	"errors"
	"proomet/config"
	"proomet/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义声明结构体
type Claims struct {
	models.JwtUser
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID uint, username, role string) (string, error) {
	// 设置Token过期时间
	expirationTime := time.Now().Add(time.Duration(config.AppConfig.JWT.Expired) * time.Second) // Token 24小时后过期
	// 创建声明
	claims := &Claims{
		JwtUser: models.JwtUser{
			UserID:   userID,
			Username: username,
			Role:     role,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "proomet",
		},
	}

	// 创建Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名Token
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证Token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
