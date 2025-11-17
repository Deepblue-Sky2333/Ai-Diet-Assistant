package handler

import (
	"errors"
	"time"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/service"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=128"`
	Email    string `json:"email,omitempty" binding:"omitempty,email,max=100"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required,min=20"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=8,max=128"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=128"`
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户使用用户名和密码登录系统
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录请求"
// @Success 200 {object} utils.Response{data=utils.TokenPair}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid request parameters", err))
		return
	}

	// 获取客户端 IP 地址
	ipAddress := c.ClientIP()

	// 调用服务层登录
	tokenPair, err := h.authService.Login(c.Request.Context(), req.Username, req.Password, ipAddress)
	if err != nil {
		if errors.Is(err, service.ErrAccountLocked) {
			utils.Error(c, utils.NewAppError(utils.CodeTooManyRequests, "account locked due to too many failed login attempts", err))
			return
		}
		if errors.Is(err, service.ErrInvalidCredentials) {
			utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "invalid username or password", err))
			return
		}
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "login failed", err))
		return
	}

	utils.Success(c, tokenPair)
}

// RefreshToken 刷新访问令牌
// @Summary 刷新访问令牌
// @Description 使用 Refresh Token 获取新的 Access Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "刷新令牌请求"
// @Success 200 {object} utils.Response{data=map[string]string}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid request parameters", err))
		return
	}

	// 刷新 token
	accessToken, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "invalid or expired refresh token", err))
		return
	}

	utils.Success(c, gin.H{
		"access_token": accessToken,
	})
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出系统，将当前令牌加入黑名单
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从 Authorization 头获取 token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.SuccessWithMessage(c, "logout successful", nil)
		return
	}

	// 提取 token
	var token string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		utils.SuccessWithMessage(c, "logout successful", nil)
		return
	}

	// 调用服务层登出
	if err := h.authService.Logout(c.Request.Context(), token); err != nil {
		// 即使失败也返回成功（幂等性）
		// 日志记录可以在中间件层处理
	}

	utils.SuccessWithMessage(c, "logout successful", nil)
}

// Register 用户注册
// @Summary 用户注册
// @Description 新用户注册账户
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册请求"
// @Success 200 {object} utils.Response{data=RegisterResponse}
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid request parameters", err))
		return
	}

	// 调用服务层注册
	user, err := h.authService.Register(c.Request.Context(), req.Username, req.Password, req.Email)
	if err != nil {
		if errors.Is(err, service.ErrUsernameExists) {
			utils.Error(c, utils.NewAppError(utils.CodeConflict, "username already exists", err))
			return
		}
		if errors.Is(err, service.ErrRegistrationDisabled) {
			utils.Error(c, utils.NewAppError(utils.CodeForbidden, "registration is currently disabled", err))
			return
		}
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "registration failed", err))
		return
	}

	// 构建响应
	response := RegisterResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	utils.SuccessWithMessage(c, "user registered successfully", response)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 用户修改登录密码
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/auth/password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid request parameters", err))
		return
	}

	// 从上下文中获取用户 ID（由认证中间件注入）
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// 修改密码
	err := h.authService.ChangePassword(c.Request.Context(), userID.(int64), req.OldPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidPassword) {
			utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "old password is incorrect", err))
			return
		}
		if errors.Is(err, service.ErrInvalidCredentials) {
			utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not found", err))
			return
		}
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to change password", err))
		return
	}

	utils.SuccessWithMessage(c, "password changed successfully, please login again", nil)
}

// RegisterRoutes 注册认证相关路由
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register) // 注册路由，不需要认证
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/logout", h.Logout)
		auth.PUT("/password", h.ChangePassword) // 需要认证中间件
	}
}
