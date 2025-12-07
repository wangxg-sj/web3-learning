package repository

import (
	"fmt"

	"github.com/wangxg-sj/web3-learning/gobasic/task4/model"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/utils"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := "root:1qaz@WSX@tcp(127.0.0.1:3306)/task4?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		utils.Error("数据库连接失败", zap.Error(err))
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}
	utils.Info("数据库连接成功", zap.String("database", "task4"))
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	utils.Info("开始数据库迁移")
	err := db.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Comment{},
	)
	if err != nil {
		utils.Error("数据库迁移失败", zap.Error(err))
		return err
	}
	utils.Info("数据库迁移成功", zap.Int("models", 3))
	return nil
}
