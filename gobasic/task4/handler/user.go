package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/service"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserHandler struct {
	db          *gorm.DB
	userService *service.UserService
}

func NewUserHandler(db *gorm.DB, userService *service.UserService) *UserHandler {
	return &UserHandler{db, userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Error("用户注册参数验证失败", zap.Error(err), zap.String("email", req.Email))
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.userService.Register(req.Name, req.Password, req.Email)
	if err != nil {
		utils.Error("用户注册失败", zap.Error(err), zap.String("email", req.Email))
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Info("用户注册成功", zap.String("email", req.Email), zap.String("name", req.Name))
	utils.SuccessResponse(c, http.StatusCreated, "注册成功", nil)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error("用户登录参数验证失败", zap.Error(err), zap.String("email", req.Email))
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		utils.Error("用户登录失败", zap.Error(err), zap.String("email", req.Email))
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Info("用户登录成功", zap.String("email", req.Email))
	utils.SuccessResponse(c, http.StatusOK, "登录成功", map[string]string{"token": token})
}
