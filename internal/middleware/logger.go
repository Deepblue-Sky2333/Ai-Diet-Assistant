package middleware

import (
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// sensitiveFields 敏感字段列表
var sensitiveFields = []string{
	"password",
	"api_key",
	"secret",
	"token",
	"authorization",
}

// dsnPattern 匹配数据库连接字符串的正则表达式
// 匹配格式: user:password@tcp(host:port)/dbname?params
var dsnPattern = regexp.MustCompile(`([^:]+):([^@]+)@tcp\(([^)]+)\)/([^?]+)(\?.*)?`)

// LoggerMiddleware 日志中间件
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()
		
		// 处理请求
		c.Next()
		
		// 计算响应时间
		duration := time.Since(startTime)
		
		// 获取用户信息
		userID, _ := GetUserID(c)
		username, _ := GetUsername(c)
		
		// 构建日志字段
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", sanitizeQuery(c.Request.URL.RawQuery)),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.Int64("duration_ms", duration.Milliseconds()),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		}
		
		// 添加用户信息（如果存在）
		if userID > 0 {
			fields = append(fields, zap.Int64("user_id", userID))
		}
		if username != "" {
			fields = append(fields, zap.String("username", username))
		}
		
		// 添加错误信息（如果存在）
		if len(c.Errors) > 0 {
			// 清理错误消息中的敏感信息（如DSN）
			sanitizedErrors := sanitizeDSN(c.Errors.String())
			fields = append(fields, zap.String("errors", sanitizedErrors))
		}
		
		// 根据状态码选择日志级别
		switch {
		case c.Writer.Status() >= 500:
			logger.Error("HTTP Request", fields...)
		case c.Writer.Status() >= 400:
			logger.Warn("HTTP Request", fields...)
		default:
			logger.Info("HTTP Request", fields...)
		}
	}
}

// sanitizeQuery 脱敏查询参数
func sanitizeQuery(query string) string {
	if query == "" {
		return ""
	}
	
	// 简单的脱敏处理：检查是否包含敏感字段
	lowerQuery := strings.ToLower(query)
	for _, field := range sensitiveFields {
		if strings.Contains(lowerQuery, field) {
			return "[REDACTED]"
		}
	}
	
	return query
}

// sanitizeValue 脱敏值
func sanitizeValue(key, value string) string {
	lowerKey := strings.ToLower(key)
	for _, field := range sensitiveFields {
		if strings.Contains(lowerKey, field) {
			return "***"
		}
	}
	return value
}

// sanitizeDSN 清理字符串中的数据库连接字符串密码
// 将 user:password@tcp(host:port)/dbname 格式中的密码替换为 ***
func sanitizeDSN(message string) string {
	// 使用正则表达式匹配并替换DSN中的密码部分
	// 匹配格式: user:password@tcp(host:port)/dbname?params
	// 替换为: user:***@tcp(host:port)/dbname?params
	return dsnPattern.ReplaceAllString(message, "$1:***@tcp($3)/$4$5")
}
