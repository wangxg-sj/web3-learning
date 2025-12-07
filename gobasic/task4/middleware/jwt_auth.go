package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/utils"
	"go.uber.org/zap"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			utils.Warn("JWT认证失败:缺少Token",
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			utils.Warn("JWT认证失败:Token无效",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		expirationTime, err := claims.GetExpirationTime()
		if err != nil {
			utils.Warn("JWT认证失败:Token过期时间错误",
				zap.Error(err),
				zap.String("ip", c.ClientIP()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token过期时间错误"})
			c.Abort()
			return
		}

		if expirationTime.Time.Before(time.Now()) {
			utils.Warn("JWT认证失败:Token已过期",
				zap.Time("expired_at", expirationTime.Time),
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token已过期"})
			c.Abort()
			return
		}

		userID, ok := (*claims)["user_id"]
		if !ok {
			utils.Error("JWT认证失败:Token缺少user_id",
				zap.String("ip", c.ClientIP()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token中未包含user_id"})
			c.Abort()
			return
		}

		// 认证成功,记录 Debug 日志
		utils.Debug("JWT认证成功",
			zap.Any("user_id", userID),
			zap.String("path", c.Request.URL.Path))

		c.Set("userID", userID)
		//c.Next()
	}
}
