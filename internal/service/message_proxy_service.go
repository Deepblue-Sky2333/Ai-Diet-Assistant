package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/ai"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
	"go.uber.org/zap"
)

// getSimpleLogger creates a simple logger for AI proxy client
func getSimpleLogger() *utils.Logger {
	// Create a no-op zap logger for the AI proxy client
	zapLogger := zap.NewNop()
	return utils.NewLogger(zapLogger)
}

const (
	// MaxMessageSize 最大消息大小 (10MB)
	MaxMessageSize = 10 * 1024 * 1024
	// DefaultTimeout 默认超时时间 (30秒)
	DefaultTimeout = 30 * time.Second
	// DefaultMaxRetries 默认最大重试次数
	DefaultMaxRetries = 2
)

var (
	// ErrMessageTooLarge 消息过大
	ErrMessageTooLarge = errors.New("message too large: maximum 10MB")
	// ErrAIServiceUnavailable AI服务不可用
	ErrAIServiceUnavailable = errors.New("AI service unavailable")
)

// MessageProxyService 消息代理服务接口
type MessageProxyService interface {
	// SendMessage sends a message to external AI service and stores the conversation
	SendMessage(ctx context.Context, userID, convID int64, content string) (*model.MessageResponse, error)
}

// messageProxyService 消息代理服务实现
type messageProxyService struct {
	convRepo       repository.ConversationRepository
	msgRepo        repository.MessageRepository
	aiSettingsRepo *repository.AISettingsRepository
}

// NewMessageProxyService creates a new message proxy service
func NewMessageProxyService(
	convRepo repository.ConversationRepository,
	msgRepo repository.MessageRepository,
	aiSettingsRepo *repository.AISettingsRepository,
) MessageProxyService {
	return &messageProxyService{
		convRepo:       convRepo,
		msgRepo:        msgRepo,
		aiSettingsRepo: aiSettingsRepo,
	}
}

// SendMessage sends a message to external AI service and stores the conversation
func (s *messageProxyService) SendMessage(ctx context.Context, userID, convID int64, content string) (*model.MessageResponse, error) {
	// Validate message size
	if len(content) > MaxMessageSize {
		return nil, ErrMessageTooLarge
	}

	// Trim content
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}

	// Verify conversation exists and belongs to user
	conv, err := s.convRepo.GetByID(ctx, userID, convID)
	if err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			return nil, ErrConversationNotFound
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Load AI configuration for user
	aiConfig, err := s.loadAIConfig(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAIServiceUnavailable, err)
	}

	// Create AI proxy client with simple logger
	aiClient := ai.NewHTTPProxyClient(aiConfig, getSimpleLogger())

	// Get conversation history to build context
	messages, _, err := s.msgRepo.GetByConversationID(ctx, userID, convID, 1, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation history: %w", err)
	}

	// Build AI request with conversation history
	aiMessages := make([]model.AIProxyMessage, 0, len(messages)+1)
	for _, msg := range messages {
		aiMessages = append(aiMessages, model.AIProxyMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	// Add current user message
	aiMessages = append(aiMessages, model.AIProxyMessage{
		Role:    model.MessageRoleUser,
		Content: content,
	})

	aiRequest := &model.AIProxyRequest{
		Messages: aiMessages,
	}

	// Marshal request to JSON for storage
	rawRequestJSON, err := json.Marshal(aiRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal AI request: %w", err)
	}

	// Send message to AI service
	aiResponse, err := aiClient.SendMessage(ctx, aiRequest)
	if err != nil {
		// Check if it's an AI service error
		if _, ok := err.(*model.AIProxyError); ok {
			return nil, fmt.Errorf("%w: %v", ErrAIServiceUnavailable, err)
		}
		return nil, fmt.Errorf("failed to send message to AI service: %w", err)
	}

	// Store user message
	userMessage := &model.Message{
		ConversationID: convID,
		Role:           model.MessageRoleUser,
		Content:        content,
		RawRequest:     string(rawRequestJSON),
		RawResponse:    "",
	}

	if err := s.msgRepo.Create(ctx, userMessage); err != nil {
		return nil, fmt.Errorf("failed to store user message: %w", err)
	}

	// Store AI response
	assistantMessage := &model.Message{
		ConversationID: convID,
		Role:           model.MessageRoleAssistant,
		Content:        aiResponse.Content,
		RawRequest:     "",
		RawResponse:    aiResponse.RawResponse,
	}

	if err := s.msgRepo.Create(ctx, assistantMessage); err != nil {
		return nil, fmt.Errorf("failed to store AI response: %w", err)
	}

	// Update conversation message count and updated_at
	if err := s.convRepo.IncrementMessageCount(ctx, convID); err != nil {
		// Log error but don't fail the request
		// The messages are already stored successfully
		fmt.Printf("Warning: failed to increment message count for conversation %d: %v\n", convID, err)
	}
	// Increment again for the assistant message
	if err := s.convRepo.IncrementMessageCount(ctx, convID); err != nil {
		fmt.Printf("Warning: failed to increment message count for conversation %d: %v\n", convID, err)
	}

	// Update conversation title if this is the first message
	if conv.MessageCount == 0 && conv.Title == "New Conversation" {
		// Use first 30 characters of user message as title
		title := content
		if len(title) > DefaultTitleMaxLength {
			// Find a good breaking point (space, punctuation)
			title = title[:DefaultTitleMaxLength]
			if lastSpace := strings.LastIndexAny(title, " \t\n.,;!?"); lastSpace > 0 {
				title = title[:lastSpace]
			}
			title = title + "..."
		}
		conv.Title = title
		if err := s.convRepo.Update(ctx, conv); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: failed to update conversation title for conversation %d: %v\n", convID, err)
		}
	}

	// Return the assistant's response
	return &model.MessageResponse{
		ID:             assistantMessage.ID,
		ConversationID: convID,
		Role:           model.MessageRoleAssistant,
		Content:        aiResponse.Content,
		CreatedAt:      assistantMessage.CreatedAt,
	}, nil
}

// loadAIConfig loads AI configuration for a user
func (s *messageProxyService) loadAIConfig(ctx context.Context, userID int64) (*model.AIProxyConfig, error) {
	// Get active AI settings for user
	settings, err := s.aiSettingsRepo.GetActiveAISettings(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrAISettingsNotFound) {
			return nil, fmt.Errorf("no active AI settings found for user")
		}
		return nil, fmt.Errorf("failed to get AI settings: %w", err)
	}

	// Validate required fields
	if settings.APIEndpoint == "" {
		return nil, fmt.Errorf("AI API endpoint is not configured")
	}
	if settings.APIKey == "" {
		return nil, fmt.Errorf("AI API key is not configured")
	}

	// Create AI proxy config
	config := &model.AIProxyConfig{
		APIEndpoint: settings.APIEndpoint,
		APIKey:      settings.APIKey,
		Model:       settings.Model,
		Timeout:     DefaultTimeout,
		MaxRetries:  DefaultMaxRetries,
	}

	return config, nil
}
