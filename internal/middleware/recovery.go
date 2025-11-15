package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ai-diet-assistant/internal/utils"
	"go.uber.org/zap"
)

// RecoveryMiddleware 恢复中间件，捕获 panic 并记录错误
func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				stack := debug.Stack()
				
				// 获取用户信息
				userID, _ := GetUserID(c)
				username, _ := GetUsername(c)
				
				// 记录错误日志
				fields := []zap.Field{
					zap.Any("error", err),
					zap.String("stack", string(stack)),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("ip", c.ClientIP()),
				}
				
				if userID > 0 {
					fields = append(fields, zap.Int64("user_id", userID))
				}
				if username != "" {
					fields = append(fields, zap.String("username", username))
				}
				
				logger.Error("Panic recovered", fields...)
				
				// 返回 500 错误响应
				utils.Error(c, utils.NewAppError(
					utils.CodeInternalError,
					"internal server error",
					fmt.Errorf("%v", err),
				))
				
				// 中止请求处理
				c.Abort()
			}
		}()
		
		c.Next()
	}
}
