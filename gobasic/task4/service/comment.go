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

type CommentService struct {
	db *gorm.DB
}

func (s CommentService) CreateComment(context *gin.Context, content string, postID uint) {
	userID, exists := context.Get("userID")
	if !exists {
		utils.Error("创建评论失败:未授权")
		utils.ErrorResponse(context, http.StatusUnauthorized, "未授权")
		return
	}

	uid := uint(reflect.ValueOf(userID).Float())
	utils.Debug("开始创建评论", zap.Uint("user_id", uid), zap.Uint("post_id", postID))

	comment := model.Comment{
		Content: content,
		PostID:  postID,
		UserID:  uid,
	}
	if err := s.db.Create(&comment).Error; err != nil {
		utils.Error("创建评论失败:数据库错误", zap.Error(err), zap.Uint("post_id", postID))
		utils.ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Info("创建评论成功(Service)", zap.Uint("user_id", uid), zap.Uint("post_id", postID), zap.Uint("comment_id", comment.ID))
	utils.SuccessResponse(context, http.StatusOK, "评论成功", nil)
}

func (s CommentService) GetComments(c *gin.Context, u uint) {
	utils.Debug("查询评论列表", zap.Uint("post_id", u))

	comments, err := gorm.G[model.Comment](s.db).Where("post_id = ?", u).Find(context.Background())
	if err != nil {
		utils.Error("查询评论列表失败", zap.Error(err), zap.Uint("post_id", u))
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Debug("查询评论列表成功", zap.Uint("post_id", u), zap.Int("count", len(comments)))
	utils.SuccessResponse(c, http.StatusOK, "获取成功", comments)
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}
