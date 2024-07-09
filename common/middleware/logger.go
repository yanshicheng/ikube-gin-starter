package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikube-gin-starter/global"
	"go.uber.org/zap"
	"time"
)

// GinLogger 自定义 Gin 日志中间件
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 继续处理请求
		c.Next()

		// 排除 /healthz 路径的日志记录
		if path == "/healthz" {
			return
		}

		// 计算请求处理耗时
		cost := time.Since(start)

		// 构建日志字段
		logger := global.L.Named("gin-web")
		if global.C.Logger.Format == "json" {
			logger.Info(
				"Request handled",
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Duration("cost", cost),
			)
		} else {
			// 构建日志消息
			logMessage := fmt.Sprintf("Request handled status=%d method=%s path=%s query=%s ip=%s user-agent=%s errors=%s cost=%s",
				c.Writer.Status(),
				c.Request.Method,
				path,
				query,
				c.ClientIP(),
				c.Request.UserAgent(),
				c.Errors.ByType(gin.ErrorTypePrivate).String(),
				cost,
			)
			logger.Info(logMessage)
		}

	}
}
