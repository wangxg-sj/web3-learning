package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/service"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CommentHandler struct {
	db             *gorm.DB
	commentService *service.CommentService
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req struct {
		Content string `json:"content"`
		PostID  uint   `json:"post_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error("创建评论参数验证失败", zap.Error(err))
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	utils.Info("创建评论请求", zap.Uint("post_id", req.PostID))
	h.commentService.CreateComment(c, req.Content, req.PostID)
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	postId, err := strconv.ParseUint(c.Param("post_id"), 10, 32)
	if err != nil {
		utils.Error("获取评论参数错误", zap.Error(err))
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	utils.Debug("获取评论列表", zap.Uint64("post_id", postId))
	h.commentService.GetComments(c, uint(postId))
}

func NewCommentHandler(db *gorm.DB, commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{db, commentService}
}
