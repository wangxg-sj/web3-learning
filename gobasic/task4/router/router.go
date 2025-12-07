package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/handler"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/middleware"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/service"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// 注册日志中间件
	router.Use(middleware.LoggerMiddleware())

	// 初始化服务
	userService := service.NewUserService(db)
	userHandler := handler.NewUserHandler(db, userService)

	postService := service.NewPostService(db)
	postHandler := handler.NewPostHandler(db, postService)

	commentService := service.NewCommentService(db)
	commentHandler := handler.NewCommentHandler(db, commentService)

	utils.Info("服务初始化完成", zap.Int("services", 3))

	// 公开路由
	group := router.Group("/api/v1")
	{
		group.POST("/register", userHandler.Register)
		group.POST("/login", userHandler.Login)
	}

	// 需要验证的路由
	privateGroup := router.Group("/api/v1")
	privateGroup.Use(middleware.AuthMiddleware())
	{
		privateGroup.POST("/posts", postHandler.CreatePost)
		privateGroup.GET("/posts", postHandler.GetPosts)
		privateGroup.GET("/posts/:id", postHandler.GetPost)
		privateGroup.DELETE("/posts/:id", postHandler.DeletePost)
		privateGroup.PUT("/posts/:id", postHandler.UpdatePost)

		privateGroup.POST("/comments", commentHandler.CreateComment)
		privateGroup.GET("/comments/:post_id", commentHandler.GetComments)

	}

	return router
}

func StartRouter(r *gin.Engine) {
	r.Run(":8080")
}
