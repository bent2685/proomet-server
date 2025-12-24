package models

import (
	"gorm.io/gorm"
)

type UserRole string

// 角色常量
const (
	RoleAdmin  = "admin"  // 超级管理员
	RoleMember = "member" // 普通成员
	RoleGuest  = "guest"  // 访客
)

// User 用户模型
type User struct {
	gorm.Model
	Username     string `gorm:"type:varchar(64);uniqueIndex;not null;comment:登录名" json:"username"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Nickname     string `gorm:"type:varchar(64);comment:昵称" json:"nickname"`
	Email        string `gorm:"type:varchar(128);index;comment:邮箱" json:"email"`
	Role         string `gorm:"type:varchar(20);default:'member';index;comment:角色标识" json:"role"`
	Status       int    `gorm:"type:smallint;default:1;comment:状态(1:正常, 2:禁用)" json:"status"`
}

type JwtUser struct {
	UserID   uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
