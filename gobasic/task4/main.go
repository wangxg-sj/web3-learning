package main

import (
	"github.com/wangxg-sj/web3-learning/gobasic/task4/repository"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/router"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/utils"
	"go.uber.org/zap"
)

func main() {
	// 初始化日志系统
	if err := utils.InitLogger(); err != nil {
		panic("初始化日志系统失败: " + err.Error())
	}
	defer utils.Sync() // 程序退出时刷新日志缓冲区

	utils.Info("应用程序启动", zap.String("app", "task4"))

	// 初始化数据库
	db, err := repository.InitDB()
	if err != nil {
		utils.Fatal("数据库初始化失败", zap.Error(err))
	}
	utils.Info("数据库连接成功")

	// 数据库迁移
	if err := repository.AutoMigrate(db); err != nil {
		utils.Fatal("数据库迁移失败", zap.Error(err))
	}
	utils.Info("数据库迁移完成")

	// 初始化路由
	r := router.InitRouter(db)
	utils.Info("路由初始化完成")

	// 启动服务器
	utils.Info("服务器启动", zap.String("port", "8080"))
	router.StartRouter(r)
}
