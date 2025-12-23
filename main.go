package main

import (
	"fmt"
	"log"
	"os"
	"proomet/config"
	_ "proomet/docs"
	"proomet/internal/infra/database"
	"proomet/internal/infra/ofs"
	"proomet/internal/interfaces/routes"
	"proomet/internal/interfaces/validators"
	"proomet/internal/middleware"
	"proomet/pkg/utils"
	"slices"

	"github.com/gin-gonic/gin"
)

// @title proomet api docs
// @version 1.0
// @description proomet api docs
// @host localhost:7071
// @BasePath /
func main() {
	// 检查是否有迁移参数
	args := os.Args[1:]
	runMigration := slices.Contains(args, "migrate")

	// 初始化配置
	config.Init("")

	// 设置Gin模式
	if config.AppConfig.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 注册自定义验证器
	validators.RegisterCustomValidators()

	// 初始化日志系统
	if err := utils.SetupLogging(config.AppConfig.Log); err != nil {
		utils.Log.Fatalf("日志系统初始化失败: %v", err)
	}

	// 记录启动日志
	utils.Log.Info("服务启动中...")

	// 初始化数据库
	defer database.Close()
	database.InitDatabase()
	// 初始化s3
	ofs.InitOfs()
	if runMigration {
		utils.Log.Info("执行数据库迁移...")
		database.AutoMigrate()

		utils.Log.Info("数据库迁移完成")
		return
	}

	r := gin.New()
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.RequestIDMiddleware())

	routerManager := routes.NewRouterManager()
	routerManager.RegisterRouter(routes.NewTestRouter())
	routerManager.SetupRoutes(r)

	addr := fmt.Sprintf("%s:%s", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	utils.Log.Success("🎉[proomet-server] 服务启动完成")

	// 启动服务器
	if err := r.Run(addr); err != nil {
		log.Fatal("服务器启动失败:", err)
	}

}
