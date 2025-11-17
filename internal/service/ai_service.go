package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
)

// AIService handles AI-related business logic
type AIService struct {
	aiSettingsRepo  *repository.AISettingsRepository
	chatHistoryRepo *repository.ChatHistoryRepository
}

// NewAIService creates a new AIService instance
func NewAIService(
	aiSettingsRepo *repository.AISettingsRepository,
	chatHistoryRepo *repository.ChatHistoryRepository,
) *AIService {
	return &AIService{
		aiSettingsRepo:  aiSettingsRepo,
		chatHistoryRepo: chatHistoryRepo,
	}
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
