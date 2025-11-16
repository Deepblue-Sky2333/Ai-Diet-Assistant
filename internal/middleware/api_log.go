package middleware

import (
	"context"
	"time"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// APILogMiddleware API 日志中间件
// 记录每个 API 请求的详细信息并异步写入数据库
func APILogMiddleware(repo repository.APILogRepository, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算响应时间
		duration := time.Since(startTime)

		// 获取用户 ID（如果已认证）
		var userID *int64
		if id, exists := GetUserID(c); exists && id > 0 {
			userID = &id
		}

		// 构建日志对象
		apiLog := &model.APILog{
			UserID:         userID,
			Method:         c.Request.Method,
			Path:           c.Request.URL.Path,
			StatusCode:     c.Writer.Status(),
			IPAddress:      c.ClientIP(),
			UserAgent:      c.Request.UserAgent(),
			ResponseTimeMs: int(duration.Milliseconds()),
		}

		// 异步写入数据库，避免阻塞请求
		go func() {
			// 创建新的 context，避免使用已取消的请求 context
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// 写入数据库
			if err := repo.CreateAPILog(ctx, apiLog); err != nil {
				// 记录错误但不影响请求处理
				logger.Error("Failed to create API log",
					zap.Error(err),
					zap.String("method", apiLog.Method),
					zap.String("path", apiLog.Path),
					zap.Int("status", apiLog.StatusCode),
				)
			}
		}()
	}
}
