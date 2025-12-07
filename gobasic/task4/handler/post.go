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

type PostHandler struct {
	db          *gorm.DB
	postService *service.PostService
}

func NewPostHandler(db *gorm.DB, postService *service.PostService) *PostHandler {
	return &PostHandler{db, postService}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error("创建文章参数验证失败", zap.Error(err))
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	utils.Debug("创建文章请求", zap.String("title", req.Title))
	h.postService.CreatePost(c, req.Title, req.Content)
}

func (h *PostHandler) GetPost(context *gin.Context) {
	postId, err := strconv.ParseUint(context.Param("id"), 10, 32)
	if err != nil {
		utils.Error("获取文章参数错误", zap.Error(err))
		utils.ErrorResponse(context, http.StatusBadRequest, "请求参数错误")
		return
	}

	utils.Debug("获取文章详情", zap.Uint64("post_id", postId))
	h.postService.GetPost(context, uint(postId))

}

func (h *PostHandler) GetPosts(context *gin.Context) {
	utils.Debug("获取文章列表")
	h.postService.GetPosts(context)
}

func (h *PostHandler) DeletePost(context *gin.Context) {
	postId, err := strconv.ParseUint(context.Param("id"), 10, 32)
	if err != nil {
		utils.Error("删除文章参数错误", zap.Error(err))
		utils.ErrorResponse(context, http.StatusBadRequest, "请求参数错误")
		return
	}

	utils.Info("删除文章请求", zap.Uint64("post_id", postId))
	h.postService.DeletePost(context, uint(postId))
}

func (h *PostHandler) UpdatePost(context *gin.Context) {
	postId, err := strconv.ParseUint(context.Param("id"), 10, 32)
	if err != nil {
		utils.Error("更新文章参数错误", zap.Error(err))
		utils.ErrorResponse(context, http.StatusBadRequest, "请求参数错误")
		return
	}
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := context.ShouldBindJSON(&req); err != nil {
		utils.Error("更新文章参数验证失败", zap.Error(err))
		utils.ErrorResponse(context, http.StatusBadRequest, "请求参数错误")
		return
	}

	utils.Info("更新文章请求", zap.Uint64("post_id", postId), zap.String("title", req.Title))
	h.postService.UpdatePost(context, uint(postId), req.Title, req.Content)
}
