package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/model"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db}
}

func (s *UserService) Register(name, password, email string) error {
	utils.Debug("开始用户注册", zap.String("email", email), zap.String("name", name))

	// 检查用户是否已存在
	var existingUser model.User
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		utils.Warn("用户注册失败:邮箱已存在", zap.String("email", email))
		return errors.New("用户已存在")
	}

	// 创建新用户
	user := model.User{
		Name:     name,
		Email:    email,
		Password: utils.CreatePassword(password), // 密码加密存储
	}
	if err := s.db.Create(&user).Error; err != nil {
		utils.Error("用户注册失败:数据库错误", zap.Error(err), zap.String("email", email))
		return err
	}

	utils.Info("用户注册成功(Service)", zap.String("email", email), zap.Uint("user_id", user.ID))
	return nil
}

func (s *UserService) Login(email string, password string) (string, error) {
	utils.Debug("开始用户登录", zap.String("email", email))

	var user model.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		utils.Warn("用户登录失败:用户不存在", zap.String("email", email))
		return "", errors.New("用户不存在")
	}

	if !utils.CheckPassword(password, user.Password) {
		utils.Warn("用户登录失败:密码错误", zap.String("email", email), zap.Uint("user_id", user.ID))
		return "", errors.New("密码错误")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token, err := utils.GenerateToken(claims)
	if err != nil {
		utils.Error("生成Token失败", zap.Error(err), zap.String("email", email))
		return "", err
	}

	utils.Info("用户登录成功(Service)", zap.String("email", email), zap.Uint("user_id", user.ID))
	return token, nil
}
