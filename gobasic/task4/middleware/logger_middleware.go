package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/utils"
	"go.uber.org/zap"
)

// LoggerMiddleware HTTP 请求日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算请求耗时
		duration := time.Since(startTime)

		// 获取请求信息
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 构建日志字段
		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.Duration("duration", duration),
		}

		// 根据状态码选择日志级别
		if statusCode >= 500 {
			// 服务器错误
			if errorMessage != "" {
				fields = append(fields, zap.String("error", errorMessage))
			}
			utils.Error("服务器错误", fields...)
		} else if statusCode >= 400 {
			// 客户端错误
			if errorMessage != "" {
				fields = append(fields, zap.String("error", errorMessage))
			}
			utils.Warn("客户端错误", fields...)
		} else {
			// 正常请求
			utils.Info("HTTP 请求", fields...)
		}
	}
}
