package service

import (
	"context"
	"fmt"

	"github.com/yourusername/ai-diet-assistant/internal/ai"
	"github.com/yourusername/ai-diet-assistant/internal/model"
	"github.com/yourusername/ai-diet-assistant/internal/repository"
)

// SettingsService 设置服务接口
type SettingsService interface {
	GetAISettings(ctx context.Context, userID int64) (*model.AISettings, error)
	ListAISettings(ctx context.Context, userID int64) ([]*model.AISettings, error)
	UpdateAISettings(ctx context.Context, userID int64, settings *model.AISettings) error
	TestAIConnection(ctx context.Context, userID int64) error
	GetUserPreferences(ctx context.Context, userID int64) (*model.UserPreferences, error)
	UpdateUserPreferences(ctx context.Context, userID int64, prefs *model.UserPreferences) error
}

type settingsService struct {
	aiSettingsRepo *repository.AISettingsRepository
	userPrefsRepo  repository.UserPreferencesRepository
}

// NewSettingsService 创建设置服务实例
func NewSettingsService(
	aiSettingsRepo *repository.AISettingsRepository,
	userPrefsRepo repository.UserPreferencesRepository,
) SettingsService {
	return &settingsService{
		aiSettingsRepo: aiSettingsRepo,
		userPrefsRepo:  userPrefsRepo,
	}
}

// GetAISettings 获取活跃的 AI 设置（调用 AI Settings Repository）
func (s *settingsService) GetAISettings(ctx context.Context, userID int64) (*model.AISettings, error) {
	settings, err := s.aiSettingsRepo.GetActiveAISettings(ctx, userID)
	if err != nil {
		if err == repository.ErrAISettingsNotFound {
			return nil, fmt.Errorf("no active AI settings found for user")
		}
		return nil, fmt.Errorf("failed to get AI settings: %w", err)
	}

	return settings, nil
}

// ListAISettings 获取所有 AI 设置
func (s *settingsService) ListAISettings(ctx context.Context, userID int64) ([]*model.AISettings, error) {
	settings, err := s.aiSettingsRepo.ListAISettings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list AI settings: %w", err)
	}

	return settings, nil
}

// UpdateAISettings 更新 AI 设置（验证配置、加密密钥）
func (s *settingsService) UpdateAISettings(ctx context.Context, userID int64, settings *model.AISettings) error {
	// 验证配置
	if err := s.validateAISettings(settings); err != nil {
		return fmt.Errorf("invalid AI settings: %w", err)
	}

	// 设置用户 ID
	settings.UserID = userID

	// 检查是否存在设置
	existing, err := s.aiSettingsRepo.GetActiveAISettings(ctx, userID)
	if err != nil && err != repository.ErrAISettingsNotFound {
		return fmt.Errorf("failed to check existing settings: %w", err)
	}

	// 如果存在则更新，否则创建
	if existing != nil {
		settings.ID = existing.ID
		if err := s.aiSettingsRepo.UpdateAISettings(ctx, settings); err != nil {
			return fmt.Errorf("failed to update AI settings: %w", err)
		}
	} else {
		if err := s.aiSettingsRepo.CreateAISettings(ctx, settings); err != nil {
			return fmt.Errorf("failed to create AI settings: %w", err)
		}
	}

	return nil
}

// TestAIConnection 测试 AI 连接（调用 AI Provider）
func (s *settingsService) TestAIConnection(ctx context.Context, userID int64) error {
	// 获取活跃的 AI 设置
	settings, err := s.aiSettingsRepo.GetActiveAISettings(ctx, userID)
	if err != nil {
		if err == repository.ErrAISettingsNotFound {
			return fmt.Errorf("no active AI settings found")
		}
		return fmt.Errorf("failed to get AI settings: %w", err)
	}

	// 创建 AI Provider
	providerConfig := &ai.ProviderConfig{
		Provider:    settings.Provider,
		APIKey:      settings.APIKey,
		APIEndpoint: settings.APIEndpoint,
		Model:       settings.Model,
		Temperature: settings.Temperature,
		MaxTokens:   settings.MaxTokens,
		Timeout:     30,
	}

	provider, err := ai.NewAIProvider(providerConfig)
	if err != nil {
		return fmt.Errorf("failed to create AI provider: %w", err)
	}

	// 测试连接
	if err := provider.TestConnection(ctx); err != nil {
		return fmt.Errorf("AI connection test failed: %w", err)
	}

	return nil
}

// GetUserPreferences 获取用户偏好
func (s *settingsService) GetUserPreferences(ctx context.Context, userID int64) (*model.UserPreferences, error) {
	prefs, err := s.userPrefsRepo.GetPreferences(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	// 如果不存在，返回默认值
	if prefs == nil {
		prefs = &model.UserPreferences{
			UserID:            userID,
			TastePreferences:  "",
			DietaryRestrictions: "",
			DailyCaloriesGoal: 2000,
			DailyProteinGoal:  150,
			DailyCarbsGoal:    250,
			DailyFatGoal:      70,
			DailyFiberGoal:    30,
		}
	}

	return prefs, nil
}

// UpdateUserPreferences 更新用户偏好
func (s *settingsService) UpdateUserPreferences(ctx context.Context, userID int64, prefs *model.UserPreferences) error {
	// 验证偏好设置
	if err := s.validateUserPreferences(prefs); err != nil {
		return fmt.Errorf("invalid user preferences: %w", err)
	}

	// 设置用户 ID
	prefs.UserID = userID

	// 检查是否存在偏好设置
	existing, err := s.userPrefsRepo.GetPreferences(userID)
	if err != nil {
		return fmt.Errorf("failed to check existing preferences: %w", err)
	}

	// 如果存在则更新，否则创建
	if existing != nil {
		if err := s.userPrefsRepo.UpdatePreferences(prefs); err != nil {
			return fmt.Errorf("failed to update preferences: %w", err)
		}
	} else {
		if err := s.userPrefsRepo.CreatePreferences(prefs); err != nil {
			return fmt.Errorf("failed to create preferences: %w", err)
		}
	}

	return nil
}

// validateAISettings 验证 AI 设置
func (s *settingsService) validateAISettings(settings *model.AISettings) error {
	if settings.Provider == "" {
		return fmt.Errorf("provider is required")
	}

	// API Key 在 handler 层已经处理（新建时必需，更新时可选）
	// 这里只验证如果提供了 API Key，它不能为空
	if settings.APIKey == "" {
		return fmt.Errorf("API key is required")
	}

	// 验证 provider 类型
	validProviders := map[string]bool{
		"openai":   true,
		"deepseek": true,
		"custom":   true,
	}
	if !validProviders[settings.Provider] {
		return fmt.Errorf("invalid provider: %s (must be openai, deepseek, or custom)", settings.Provider)
	}

	// 自定义 provider 需要 API 端点
	if settings.Provider == "custom" && settings.APIEndpoint == "" {
		return fmt.Errorf("custom provider requires API endpoint")
	}

	// 验证温度参数
	if settings.Temperature < 0 || settings.Temperature > 2 {
		return fmt.Errorf("temperature must be between 0 and 2")
	}

	// 验证 max tokens
	if settings.MaxTokens < 1 || settings.MaxTokens > 4096 {
		return fmt.Errorf("max_tokens must be between 1 and 4096")
	}

	return nil
}

// validateUserPreferences 验证用户偏好
func (s *settingsService) validateUserPreferences(prefs *model.UserPreferences) error {
	// 验证每日卡路里目标
	if prefs.DailyCaloriesGoal < 0 {
		return fmt.Errorf("daily calories goal must be non-negative")
	}

	if prefs.DailyCaloriesGoal > 10000 {
		return fmt.Errorf("daily calories goal seems unreasonably high (max: 10000)")
	}

	// 验证蛋白质目标
	if prefs.DailyProteinGoal < 0 || prefs.DailyProteinGoal > 500 {
		return fmt.Errorf("daily protein goal must be between 0 and 500g")
	}

	// 验证碳水化合物目标
	if prefs.DailyCarbsGoal < 0 || prefs.DailyCarbsGoal > 1000 {
		return fmt.Errorf("daily carbs goal must be between 0 and 1000g")
	}

	// 验证脂肪目标
	if prefs.DailyFatGoal < 0 || prefs.DailyFatGoal > 500 {
		return fmt.Errorf("daily fat goal must be between 0 and 500g")
	}

	// 验证纤维目标
	if prefs.DailyFiberGoal < 0 || prefs.DailyFiberGoal > 200 {
		return fmt.Errorf("daily fiber goal must be between 0 and 200g")
	}

	return nil
}
