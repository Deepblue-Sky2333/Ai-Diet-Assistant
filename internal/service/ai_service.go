package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yourusername/ai-diet-assistant/internal/ai"
	"github.com/yourusername/ai-diet-assistant/internal/model"
	"github.com/yourusername/ai-diet-assistant/internal/repository"
)

// AIService handles AI-related business logic
type AIService struct {
	aiSettingsRepo  *repository.AISettingsRepository
	chatHistoryRepo *repository.ChatHistoryRepository
	foodRepo        *repository.FoodRepository
}

// NewAIService creates a new AIService instance
func NewAIService(
	aiSettingsRepo *repository.AISettingsRepository,
	chatHistoryRepo *repository.ChatHistoryRepository,
	foodRepo *repository.FoodRepository,
) *AIService {
	return &AIService{
		aiSettingsRepo:  aiSettingsRepo,
		chatHistoryRepo: chatHistoryRepo,
		foodRepo:        foodRepo,
	}
}

// GenerateMealPlan generates meal plans using AI based on available foods and preferences
func (s *AIService) GenerateMealPlan(ctx context.Context, userID int64, days int, targetCalories int) (*ai.MealPlanResponse, error) {
	// Get active AI settings
	settings, err := s.aiSettingsRepo.GetActiveAISettings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI settings: %w", err)
	}

	// Create AI provider
	provider, err := s.createProvider(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI provider: %w", err)
	}

	// Get available foods
	foods, err := s.getAvailableFoods(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get available foods: %w", err)
	}

	if len(foods) == 0 {
		return nil, fmt.Errorf("no available foods found")
	}

	// Build meal plan request
	request := &ai.MealPlanRequest{
		AvailableFoods: foods,
		Days:           days,
		TargetCalories: targetCalories,
		// Preferences can be added here when preference system is implemented
	}

	// Generate meal plan with retry logic
	var response *ai.MealPlanResponse
	maxRetries := 2
	for i := 0; i <= maxRetries; i++ {
		response, err = provider.GenerateMealPlan(ctx, request)
		if err == nil {
			break
		}
		if i < maxRetries {
			time.Sleep(time.Second * time.Duration(i+1))
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate meal plan after retries: %w", err)
	}

	return response, nil
}

// Chat handles conversational interactions with AI
func (s *AIService) Chat(ctx context.Context, userID int64, message string, contextData map[string]string) (*ai.ChatResponse, error) {
	// Get active AI settings
	settings, err := s.aiSettingsRepo.GetActiveAISettings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI settings: %w", err)
	}

	// Create AI provider
	provider, err := s.createProvider(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI provider: %w", err)
	}

	// Build chat request
	request := &ai.ChatRequest{
		Message: message,
		Context: contextData,
		UserID:  userID,
	}

	// Send chat request with retry logic
	var response *ai.ChatResponse
	maxRetries := 2
	for i := 0; i <= maxRetries; i++ {
		response, err = provider.Chat(ctx, request)
		if err == nil {
			break
		}
		if i < maxRetries {
			time.Sleep(time.Second * time.Duration(i+1))
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to chat after retries: %w", err)
	}

	// Save chat history and get the message ID
	messageID, err := s.SaveChatHistory(ctx, userID, message, response.Message, contextData, response.TokensUsed)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("warning: failed to save chat history: %v\n", err)
		// Set message ID to 0 if save failed
		response.MessageID = 0
	} else {
		// Set the message ID in the response
		response.MessageID = messageID
	}

	return response, nil
}

// SaveChatHistory saves a chat interaction to history and returns the message ID
func (s *AIService) SaveChatHistory(ctx context.Context, userID int64, userInput, aiResponse string, contextData map[string]string, tokensUsed int) (int64, error) {
	// Convert context to JSON string
	var contextJSON string
	if len(contextData) > 0 {
		contextBytes, err := json.Marshal(contextData)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal context: %w", err)
		}
		contextJSON = string(contextBytes)
	} else {
		// Use empty JSON object instead of empty string for MySQL JSON column
		contextJSON = "{}"
	}

	history := &model.ChatHistory{
		UserID:     userID,
		UserInput:  userInput,
		AIResponse: aiResponse,
		Context:    contextJSON,
		TokensUsed: tokensUsed,
	}

	if err := s.chatHistoryRepo.CreateChatHistory(ctx, history); err != nil {
		return 0, fmt.Errorf("failed to create chat history: %w", err)
	}

	return history.ID, nil
}

// GetChatHistory retrieves chat history for a user with pagination
func (s *AIService) GetChatHistory(ctx context.Context, userID int64, page, pageSize int) ([]*model.ChatHistory, int, error) {
	return s.chatHistoryRepo.GetChatHistory(ctx, userID, page, pageSize)
}

// TestAIConnection tests the connection to the AI provider
func (s *AIService) TestAIConnection(ctx context.Context, userID int64) error {
	// Get active AI settings
	settings, err := s.aiSettingsRepo.GetActiveAISettings(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get AI settings: %w", err)
	}

	// Create AI provider
	provider, err := s.createProvider(settings)
	if err != nil {
		return fmt.Errorf("failed to create AI provider: %w", err)
	}

	// Test connection
	if err := provider.TestConnection(ctx); err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	return nil
}

// createProvider creates an AI provider from settings
func (s *AIService) createProvider(settings *model.AISettings) (ai.AIProvider, error) {
	config := &ai.ProviderConfig{
		Provider:    settings.Provider,
		APIEndpoint: settings.APIEndpoint,
		APIKey:      settings.APIKey,
		Model:       settings.Model,
		Temperature: settings.Temperature,
		MaxTokens:   settings.MaxTokens,
		Timeout:     30, // Default timeout
	}

	return ai.NewAIProvider(config)
}

// getAvailableFoods retrieves all available foods for a user
func (s *AIService) getAvailableFoods(ctx context.Context, userID int64) ([]model.Food, error) {
	available := true
	filter := &model.FoodFilter{
		Available: &available,
		Page:      1,
		PageSize:  1000, // Get all available foods
	}

	foods, _, err := s.foodRepo.ListFoods(userID, filter)
	if err != nil {
		return nil, err
	}

	// Convert to non-pointer slice
	result := make([]model.Food, len(foods))
	for i, food := range foods {
		result[i] = *food
	}

	return result, nil
}
