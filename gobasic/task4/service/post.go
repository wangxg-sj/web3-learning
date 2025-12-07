package service

import (
	"context"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/model"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db}
}
func (s *PostService) CreatePost(c *gin.Context, title string, content string) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error("创建文章失败:未授权")
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	uid := uint(reflect.ValueOf(userID).Float())
	utils.Debug("开始创建文章", zap.Uint("user_id", uid), zap.String("title", title))

	// 创建帖子
	post := model.Post{
		Title:   title,
		Content: content,
		UserID:  uid,
	}
	err := gorm.G[model.Post](s.db).Create(context.Background(), &post)
	if err != nil {
		utils.Error("创建文章失败:数据库错误", zap.Error(err), zap.Uint("user_id", uid))
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Info("创建文章成功(Service)", zap.Uint("user_id", uid), zap.Uint("post_id", post.ID), zap.String("title", title))
	utils.SuccessResponse(c, http.StatusOK, "创建成功", nil)
}

func (s *PostService) GetPost(c *gin.Context, id uint) {
	utils.Debug("查询文章详情", zap.Uint("post_id", id))

	post, err := gorm.G[model.Post](s.db).Preload("Comments", nil).
		Where("id = ?", id).First(context.Background())
	if err != nil {
		utils.Error("查询文章失败", zap.Error(err), zap.Uint("post_id", id))
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Debug("查询文章成功", zap.Uint("post_id", id), zap.Int("comments_count", len(post.Comments)))
	utils.SuccessResponse(c, http.StatusOK, "获取成功", post)
}

func (s *PostService) GetPosts(c *gin.Context) {
	utils.Debug("查询文章列表")

	posts, err := gorm.G[model.Post](s.db).Preload("Comments", nil).
		Find(context.Background())
	if err != nil {
		utils.Error("查询文章列表失败", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Debug("查询文章列表成功", zap.Int("count", len(posts)))
	utils.SuccessResponse(c, http.StatusOK, "获取成功", posts)
}

func (s *PostService) DeletePost(c *gin.Context, id uint) {
	userID, _ := c.Get("userID")
	uid := uint(reflect.ValueOf(userID).Float())

	utils.Debug("开始删除文章", zap.Uint("post_id", id), zap.Uint("user_id", uid))

	post, err := gorm.G[model.Post](s.db).Where("id = ?", id).First(context.Background())
	if err != nil {
		utils.Warn("删除文章失败:文章不存在", zap.Uint("post_id", id))
		utils.ErrorResponse(c, http.StatusNotFound, "帖子不存在")
		return
	}

	// 检查用户是否有权限删除
	if post.UserID != uid {
		utils.Warn("删除文章失败:权限不足", zap.Uint("post_id", id), zap.Uint("user_id", uid), zap.Uint("owner_id", post.UserID))
		utils.ErrorResponse(c, http.StatusForbidden, "没有权限删除该帖子")
		return
	}

	_, err = gorm.G[model.Post](s.db).Where("id = ?", id).Delete(context.Background())
	if err != nil {
		utils.Error("删除文章失败:数据库错误", zap.Error(err), zap.Uint("post_id", id))
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Info("删除文章成功(Service)", zap.Uint("post_id", id), zap.Uint("user_id", uid))
	utils.SuccessResponse(c, http.StatusOK, "删除成功", nil)
}

func (s *PostService) UpdatePost(c *gin.Context, u uint, title string, content string) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error("更新文章失败:未授权")
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	uid := uint(reflect.ValueOf(userID).Float())
	utils.Debug("开始更新文章", zap.Uint("post_id", u), zap.Uint("user_id", uid), zap.String("title", title))

	post := model.Post{}
	// 检查用户是否有权限更新
	post, err := gorm.G[model.Post](s.db).Where("id = ?", u).First(context.Background())
	if err != nil {
		utils.Warn("更新文章失败:文章不存在", zap.Uint("post_id", u))
		utils.ErrorResponse(c, http.StatusNotFound, "帖子不存在")
		return
	}

	if post.UserID != uid {
		utils.Warn("更新文章失败:权限不足", zap.Uint("post_id", u), zap.Uint("user_id", uid), zap.Uint("owner_id", post.UserID))
		utils.ErrorResponse(c, http.StatusForbidden, "没有权限更新该帖子")
		return
	}

	// 更新帖子
	post.Title = title
	post.Content = content
	_, err = gorm.G[model.Post](s.db).Where("id = ?", u).Updates(context.Background(), post)
	if err != nil {
		utils.Error("更新文章失败:数据库错误", zap.Error(err), zap.Uint("post_id", u))
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Info("更新文章成功(Service)", zap.Uint("post_id", u), zap.Uint("user_id", uid), zap.String("title", title))
	utils.SuccessResponse(c, http.StatusOK, "更新成功", nil)
}
