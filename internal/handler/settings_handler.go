package handler

import (
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/middleware"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/service"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
	"github.com/gin-gonic/gin"
)

// SettingsHandler 设置处理器
type SettingsHandler struct {
	settingsService service.SettingsService
}

// NewSettingsHandler 创建设置处理器实例
func NewSettingsHandler(settingsService service.SettingsService) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
	}
}

// UpdateAISettingsRequest 更新 AI 设置请求
type UpdateAISettingsRequest struct {
	Provider    string  `json:"provider" binding:"required,oneof=openai deepseek custom"`
	APIEndpoint string  `json:"api_endpoint" binding:"omitempty,url,max=500"`
	APIKey      string  `json:"api_key" binding:"omitempty,min=10,max=500"` // 可选，更新时可以不提供
	Model       string  `json:"model" binding:"omitempty,min=1,max=100"`
	Temperature float64 `json:"temperature" binding:"omitempty,gte=0,lte=2"`
	MaxTokens   int     `json:"max_tokens" binding:"omitempty,gte=1,lte=32000"`
	IsActive    bool    `json:"is_active"`
}

// GetSettings 获取所有设置
// @Summary 获取所有设置
// @Description 获取用户的 AI 设置和偏好设置
// @Tags 设置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/settings [get]
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// 获取活跃的 AI 配置（单个对象）
	activeAIConfig, err := h.settingsService.GetAISettings(c.Request.Context(), userID.(int64))
	if err != nil {
		// 如果没有配置，不返回错误，而是返回 null
		activeAIConfig = nil
	}

	// 获取用户偏好
	userPrefs, err := h.settingsService.GetUserPreferences(c.Request.Context(), userID.(int64))
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to get user preferences", err))
		return
	}

	// 构建响应
	response := gin.H{
		"user_preferences": userPrefs,
	}

	// 如果有活跃的 AI 配置，返回安全的版本（掩码 API Key）
	if activeAIConfig != nil {
		response["ai_config"] = gin.H{
			"provider":       activeAIConfig.Provider,
			"api_endpoint":   activeAIConfig.APIEndpoint,
			"api_key_masked": activeAIConfig.MaskAPIKey(),
			"model":          activeAIConfig.Model,
			"temperature":    activeAIConfig.Temperature,
			"max_tokens":     activeAIConfig.MaxTokens,
		}
	} else {
		response["ai_config"] = nil
	}

	utils.Success(c, response)
}

// UpdateAISettings 更新 AI 设置
// @Summary 更新 AI 设置
// @Description 更新用户的 AI 配置（Provider、API Key 等）
// @Tags 设置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateAISettingsRequest true "AI 设置"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/settings/ai [put]
func (h *SettingsHandler) UpdateAISettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	var req UpdateAISettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid request parameters", err))
		return
	}

	// 获取现有配置（如果存在）
	existing, _ := h.settingsService.GetAISettings(c.Request.Context(), userID.(int64))

	// 构建 AI 设置对象
	settings := &model.AISettings{
		Provider:    req.Provider,
		APIEndpoint: req.APIEndpoint,
		IsActive:    true, // 总是设置为活跃
	}

	// 处理 API Key：如果提供了新的则使用新的，否则保留旧的
	if req.APIKey != "" {
		settings.APIKey = req.APIKey
	} else if existing != nil {
		// 保留现有的 API Key
		settings.APIKey = existing.APIKey
	} else {
		// 新配置必须提供 API Key
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "API key is required for new configuration", nil))
		return
	}

	// 处理 Model：如果提供了则使用提供的，否则使用现有的或默认值
	if req.Model != "" {
		settings.Model = req.Model
	} else if existing != nil && existing.Model != "" {
		// 保留现有的 Model
		settings.Model = existing.Model
	} else {
		// 根据 provider 设置默认模型
		switch settings.Provider {
		case "openai":
			settings.Model = "gpt-3.5-turbo"
		case "deepseek":
			settings.Model = "deepseek-chat"
		default:
			settings.Model = "default"
		}
	}

	// 处理 Temperature：如果提供了则使用提供的，否则使用现有的或默认值
	if req.Temperature != 0 {
		settings.Temperature = req.Temperature
	} else if existing != nil && existing.Temperature != 0 {
		// 保留现有的 Temperature
		settings.Temperature = existing.Temperature
	} else {
		// 使用默认值
		settings.Temperature = 0.7
	}

	// 处理 MaxTokens：如果提供了则使用提供的，否则使用现有的或默认值
	if req.MaxTokens != 0 {
		settings.MaxTokens = req.MaxTokens
	} else if existing != nil && existing.MaxTokens != 0 {
		// 保留现有的 MaxTokens
		settings.MaxTokens = existing.MaxTokens
	} else {
		// 使用默认值
		settings.MaxTokens = 1000
	}

	// 更新设置
	err := h.settingsService.UpdateAISettings(c.Request.Context(), userID.(int64), settings)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to update AI settings", err))
		return
	}

	utils.SuccessWithMessage(c, "AI settings updated successfully", nil)
}

// TestAIConnection 测试 AI 连接
// @Summary 测试 AI 连接
// @Description 测试当前配置的 AI 服务是否可用
// @Tags 设置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/settings/ai/test [get]
func (h *SettingsHandler) TestAIConnection(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// 测试连接
	err := h.settingsService.TestAIConnection(c.Request.Context(), userID.(int64))
	if err != nil {
		// 提供更友好的错误消息
		errorMsg := "AI connection test failed: " + err.Error()
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, errorMsg, err))
		return
	}

	utils.SuccessWithMessage(c, "AI connection test successful", nil)
}

// GetUserProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取用户基本信息和偏好设置
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/user/profile [get]
func (h *SettingsHandler) GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// 获取用户偏好
	userPrefs, err := h.settingsService.GetUserPreferences(c.Request.Context(), userID.(int64))
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to get user profile", err))
		return
	}

	utils.Success(c, userPrefs)
}

// UpdateUserPreferences 更新用户偏好
// @Summary 更新用户偏好
// @Description 更新用户的口味偏好、饮食限制和营养目标
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.UpdateUserPreferencesRequest true "用户偏好"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/user/preferences [put]
func (h *SettingsHandler) UpdateUserPreferences(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	var req model.UpdateUserPreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid request parameters", err))
		return
	}

	// 获取现有偏好（如果存在）以保留未提供的值
	existing, _ := h.settingsService.GetUserPreferences(c.Request.Context(), userID.(int64))

	// 构建用户偏好对象
	prefs := &model.UserPreferences{
		UserID: userID.(int64),
	}

	// 处理文本字段：使用提供的值，如果为空则保留现有值或使用空字符串
	prefs.TastePreferences = req.TastePreferences
	if prefs.TastePreferences == "" && existing != nil {
		prefs.TastePreferences = existing.TastePreferences
	}

	prefs.DietaryRestrictions = req.DietaryRestrictions
	if prefs.DietaryRestrictions == "" && existing != nil {
		prefs.DietaryRestrictions = existing.DietaryRestrictions
	}

	// 处理营养目标：使用提供的值，如果为0则使用现有值或默认值
	if req.DailyCaloriesGoal > 0 {
		prefs.DailyCaloriesGoal = req.DailyCaloriesGoal
	} else if existing != nil && existing.DailyCaloriesGoal > 0 {
		prefs.DailyCaloriesGoal = existing.DailyCaloriesGoal
	} else {
		prefs.DailyCaloriesGoal = 2000
	}

	if req.DailyProteinGoal > 0 {
		prefs.DailyProteinGoal = req.DailyProteinGoal
	} else if existing != nil && existing.DailyProteinGoal > 0 {
		prefs.DailyProteinGoal = existing.DailyProteinGoal
	} else {
		prefs.DailyProteinGoal = 150
	}

	if req.DailyCarbsGoal > 0 {
		prefs.DailyCarbsGoal = req.DailyCarbsGoal
	} else if existing != nil && existing.DailyCarbsGoal > 0 {
		prefs.DailyCarbsGoal = existing.DailyCarbsGoal
	} else {
		prefs.DailyCarbsGoal = 250
	}

	if req.DailyFatGoal > 0 {
		prefs.DailyFatGoal = req.DailyFatGoal
	} else if existing != nil && existing.DailyFatGoal > 0 {
		prefs.DailyFatGoal = existing.DailyFatGoal
	} else {
		prefs.DailyFatGoal = 70
	}

	if req.DailyFiberGoal > 0 {
		prefs.DailyFiberGoal = req.DailyFiberGoal
	} else if existing != nil && existing.DailyFiberGoal > 0 {
		prefs.DailyFiberGoal = existing.DailyFiberGoal
	} else {
		prefs.DailyFiberGoal = 30
	}

	// 更新偏好
	err := h.settingsService.UpdateUserPreferences(c.Request.Context(), userID.(int64), prefs)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to update user preferences", err))
		return
	}

	utils.SuccessWithMessage(c, "user preferences updated successfully", nil)
}

// GetSystemSettings 获取系统设置
// @Summary 获取系统设置
// @Description 获取系统级别的配置（需要管理员权限）
// @Tags 设置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/v1/settings/system [get]
func (h *SettingsHandler) GetSystemSettings(c *gin.Context) {
	// 获取系统设置
	settings, err := h.settingsService.GetSystemSettings(c.Request.Context())
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to get system settings", err))
		return
	}

	utils.Success(c, settings)
}

// UpdateSystemSettingsRequest 更新系统设置请求
type UpdateSystemSettingsRequest struct {
	RegistrationEnabled *bool `json:"registration_enabled,omitempty"`
}

// UpdateSystemSettings 更新系统设置
// @Summary 更新系统设置
// @Description 更新系统级别的配置（需要管理员权限）
// @Tags 设置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateSystemSettingsRequest true "系统设置"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/v1/settings/system [put]
func (h *SettingsHandler) UpdateSystemSettings(c *gin.Context) {
	var req UpdateSystemSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid request parameters", err))
		return
	}

	// 构建设置映射
	settings := make(map[string]interface{})
	if req.RegistrationEnabled != nil {
		settings["registration_enabled"] = *req.RegistrationEnabled
	}

	// 如果没有提供任何设置，返回错误
	if len(settings) == 0 {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "no settings provided", nil))
		return
	}

	// 更新系统设置
	err := h.settingsService.UpdateSystemSettings(c.Request.Context(), settings)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to update system settings", err))
		return
	}

	utils.SuccessWithMessage(c, "system settings updated successfully", nil)
}

// GetSystemInfo 获取公开的系统信息
// @Summary 获取公开的系统信息
// @Description 获取系统的公开信息，如注册是否开放、版本号等（无需认证）
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/system/info [get]
func (h *SettingsHandler) GetSystemInfo(c *gin.Context) {
	// 获取注册开关状态
	registrationEnabled, err := h.settingsService.IsRegistrationEnabled(c.Request.Context())
	if err != nil {
		// 如果获取失败，默认返回 true（允许注册）
		registrationEnabled = false
	}

	// 构建响应
	info := gin.H{
		"registration_enabled": registrationEnabled,
		"version":              "1.0.0", // 可以从配置或常量中读取
	}

	utils.Success(c, info)
}

// RegisterRoutes 注册设置相关路由
func (h *SettingsHandler) RegisterRoutes(router *gin.RouterGroup, userRepo repository.UserRepository) {
	settings := router.Group("/settings")
	{
		settings.GET("", h.GetSettings)
		settings.PUT("/ai", h.UpdateAISettings)
		settings.GET("/ai/test", h.TestAIConnection)

		// 系统设置路由（需要管理员权限）
		settings.GET("/system", middleware.AdminMiddleware(userRepo), h.GetSystemSettings)
		settings.PUT("/system", middleware.AdminMiddleware(userRepo), h.UpdateSystemSettings)
	}

	user := router.Group("/user")
	{
		user.GET("/profile", h.GetUserProfile)
		user.PUT("/preferences", h.UpdateUserPreferences)
	}
}
