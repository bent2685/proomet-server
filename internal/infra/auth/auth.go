package auth

import (
	"proomet/pkg/utils"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var Enforcer *casbin.Enforcer

func InitCasbin(db *gorm.DB) {
	// 1. 初始化适配器（让 Casbin 使用现有的 Gorm 实例）
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		utils.Log.Fatalf("Casbin 初始化失败 %s", err)
	}

	// 2. 加载模型配置
	Enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		utils.Log.Fatalf("Casbin 初始化失败 %s", err)
	}

	// 3. 加载策略
	if err := Enforcer.LoadPolicy(); err != nil {
		utils.Log.Fatalf("Casbin 加载策略失败 %s", err)
	}
	utils.Log.Info("Casbin 初始化成功")
}

func GetEnforcer() *casbin.Enforcer {
	return Enforcer
}
