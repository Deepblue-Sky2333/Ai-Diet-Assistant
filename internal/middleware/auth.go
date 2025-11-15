package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
)

const (
	// ContextKeyUserID 用户 ID 的 context key
	ContextKeyUserID = "user_id"
	// ContextKeyUsername 用户名的 context key
	ContextKeyUsername = "username"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware(jwtService *utils.JWTService, authService interface {
	ValidateToken(ctx context.Context, token string) (*utils.Claims, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization 头获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, utils.NewAppError(
				utils.CodeUnauthorized,
				"missing authorization header",
				nil,
			))
			c.Abort()
			return
		}

		// 检查 Bearer 前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Error(c, utils.NewAppError(
				utils.CodeUnauthorized,
				"invalid authorization header format",
				nil,
			))
			c.Abort()
			return
		}

		token := parts[1]

		// 使用 authService 验证 token（包括黑名单和密码版本检查）
		claims, err := authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			var message string
			if err == utils.ErrExpiredToken {
				message = "token has expired"
			} else if err == utils.ErrPasswordChanged {
				message = "password has been changed, please login again"
			} else {
				message = "invalid token"
			}
			utils.Error(c, utils.NewAppError(
				utils.CodeUnauthorized,
				message,
				err,
			))
			c.Abort()
			return
		}

		// 将用户信息注入到 context
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)

		c.Next()
	}
}

// GetUserID 从 context 获取用户 ID
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		return 0, false
	}
	id, ok := userID.(int64)
	return id, ok
}

// GetUsername 从 context 获取用户名
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get(ContextKeyUsername)
	if !exists {
		return "", false
	}
	name, ok := username.(string)
	return name, ok
}

// MustGetUserID 从 context 获取用户 ID，如果不存在则 panic
func MustGetUserID(c *gin.Context) int64 {
	userID, ok := GetUserID(c)
	if !ok {
		panic("user_id not found in context")
	}
	return userID
}
