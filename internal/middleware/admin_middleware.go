package middleware

import (
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
	"github.com/gin-gonic/gin"
)

const (
	// ContextKeyUserRole 用户角色的 context key
	ContextKeyUserRole = "user_role"
)

// AdminMiddleware 管理员权限检查中间件
// 此中间件必须在 AuthMiddleware 之后使用，因为它依赖于 user_id
func AdminMiddleware(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 context 获取用户 ID（由 AuthMiddleware 设置）
		userID, exists := c.Get(ContextKeyUserID)
		if !exists {
			utils.Error(c, utils.NewAppError(
				utils.CodeUnauthorized,
				"user not authenticated",
				nil,
			))
			c.Abort()
			return
		}

		// 类型断言
		uid, ok := userID.(int64)
		if !ok {
			utils.Error(c, utils.NewAppError(
				utils.CodeUnauthorized,
				"invalid user id in context",
				nil,
			))
			c.Abort()
			return
		}

		// 从数据库获取用户信息以验证角色
		user, err := userRepo.GetUserByID(c.Request.Context(), uid)
		if err != nil {
			if err == repository.ErrUserNotFound {
				utils.Error(c, utils.NewAppError(
					utils.CodeUnauthorized,
					"user not found",
					err,
				))
			} else {
				utils.Error(c, utils.NewAppError(
					utils.CodeInternalError,
					"failed to get user information",
					err,
				))
			}
			c.Abort()
			return
		}

		// 检查用户角色是否为管理员
		if user.Role != model.RoleAdmin {
			utils.Error(c, utils.NewAppError(
				utils.CodeForbidden,
				"admin access required",
				nil,
			))
			c.Abort()
			return
		}

		// 将用户角色存入上下文，供后续处理器使用
		c.Set(ContextKeyUserRole, user.Role)

		c.Next()
	}
}

// GetUserRole 从 context 获取用户角色
func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get(ContextKeyUserRole)
	if !exists {
		return "", false
	}
	r, ok := role.(string)
	return r, ok
}

// IsAdmin 检查当前用户是否为管理员
func IsAdmin(c *gin.Context) bool {
	role, exists := GetUserRole(c)
	if !exists {
		return false
	}
	return role == model.RoleAdmin
}
