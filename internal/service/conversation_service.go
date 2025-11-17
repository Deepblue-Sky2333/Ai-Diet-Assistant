package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
)

const (
	// MaxRecentConversations 最大最近对话流数量
	MaxRecentConversations = 10
	// MaxFavoriteConversations 最大收藏对话流数量
	MaxFavoriteConversations = 100
	// DefaultTitleMaxLength 默认标题最大长度（从第一条消息截取）
	DefaultTitleMaxLength = 30
	// MaxTitleLength 标题最大长度
	MaxTitleLength = 200
)

var (
	// ErrConversationNotFound 对话流不存在
	ErrConversationNotFound = errors.New("conversation not found")
	// ErrFavoriteLimitReached 收藏数量已达上限
	ErrFavoriteLimitReached = errors.New("favorite limit reached: maximum 100 favorited conversations")
	// ErrTitleTooLong 标题过长
	ErrTitleTooLong = errors.New("title too long: maximum 200 characters")
)

// ConversationService 对话流服务接口
type ConversationService interface {
	// CreateConversation creates a new conversation flow
	CreateConversation(ctx context.Context, userID int64, title string) (*model.ConversationFlow, error)

	// GetConversation retrieves a conversation by ID
	GetConversation(ctx context.Context, userID, convID int64) (*model.ConversationFlow, error)

	// ListConversations lists conversations with filters
	ListConversations(ctx context.Context, userID int64, filter *model.ConversationFilter) ([]*model.ConversationFlow, int, error)

	// UpdateConversationTitle updates the title of a conversation
	UpdateConversationTitle(ctx context.Context, userID, convID int64, title string) error

	// DeleteConversation deletes a conversation and all its messages
	DeleteConversation(ctx context.Context, userID, convID int64) error

	// FavoriteConversation marks a conversation as favorited
	FavoriteConversation(ctx context.Context, userID, convID int64) error

	// UnfavoriteConversation removes favorite status from a conversation
	UnfavoriteConversation(ctx context.Context, userID, convID int64) error

	// SearchConversations searches conversations by keyword
	SearchConversations(ctx context.Context, userID int64, keyword string, filter *model.ConversationFilter) ([]*model.ConversationFlow, int, error)

	// ExportConversation exports a conversation to JSON
	ExportConversation(ctx context.Context, userID, convID int64) ([]byte, error)

	// ExportConversations exports multiple conversations to JSON
	ExportConversations(ctx context.Context, userID int64, convIDs []int64) ([]byte, error)

	// GetMessages retrieves messages for a conversation
	GetMessages(ctx context.Context, userID, convID int64, page, pageSize int) ([]*model.Message, int, error)
}

// conversationService 对话流服务实现
type conversationService struct {
	convRepo repository.ConversationRepository
	msgRepo  repository.MessageRepository
}

// NewConversationService creates a new conversation service
func NewConversationService(
	convRepo repository.ConversationRepository,
	msgRepo repository.MessageRepository,
) ConversationService {
	return &conversationService{
		convRepo: convRepo,
		msgRepo:  msgRepo,
	}
}

// CreateConversation creates a new conversation flow with automatic cleanup
func (s *conversationService) CreateConversation(ctx context.Context, userID int64, title string) (*model.ConversationFlow, error) {
	// Validate title length
	if len(title) > MaxTitleLength {
		return nil, ErrTitleTooLong
	}

	// Use default title if empty
	if title == "" {
		title = "New Conversation"
	}

	// Check recent conversation count and cleanup if needed
	recentCount, err := s.convRepo.GetRecentCount(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent conversation count: %w", err)
	}

	// If we have 10 or more recent conversations, delete the oldest one
	if recentCount >= MaxRecentConversations {
		if err := s.convRepo.DeleteOldestRecent(ctx, userID); err != nil {
			// Log error but don't fail the creation
			// This could happen if all conversations are favorited
			if !errors.Is(err, repository.ErrConversationNotFound) {
				return nil, fmt.Errorf("failed to delete oldest recent conversation: %w", err)
			}
		}
	}

	// Create new conversation
	conv := &model.ConversationFlow{
		UserID:       userID,
		Title:        title,
		IsFavorited:  false,
		MessageCount: 0,
	}

	if err := s.convRepo.Create(ctx, conv); err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return conv, nil
}

// GetConversation retrieves a conversation by ID
func (s *conversationService) GetConversation(ctx context.Context, userID, convID int64) (*model.ConversationFlow, error) {
	conv, err := s.convRepo.GetByID(ctx, userID, convID)
	if err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			return nil, ErrConversationNotFound
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return conv, nil
}

// ListConversations lists conversations with filters
func (s *conversationService) ListConversations(ctx context.Context, userID int64, filter *model.ConversationFilter) ([]*model.ConversationFlow, int, error) {
	// Set default values
	if filter == nil {
		filter = &model.ConversationFilter{}
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.SortBy == "" {
		filter.SortBy = "updated_at"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	conversations, total, err := s.convRepo.List(ctx, userID, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list conversations: %w", err)
	}

	return conversations, total, nil
}

// UpdateConversationTitle updates the title of a conversation
func (s *conversationService) UpdateConversationTitle(ctx context.Context, userID, convID int64, title string) error {
	// Validate title length
	if len(title) > MaxTitleLength {
		return ErrTitleTooLong
	}

	if title == "" {
		return errors.New("title cannot be empty")
	}

	// Get existing conversation
	conv, err := s.convRepo.GetByID(ctx, userID, convID)
	if err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			return ErrConversationNotFound
		}
		return fmt.Errorf("failed to get conversation: %w", err)
	}

	// Update title
	conv.Title = title
	if err := s.convRepo.Update(ctx, conv); err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			return ErrConversationNotFound
		}
		return fmt.Errorf("failed to update conversation title: %w", err)
	}

	return nil
}

// DeleteConversation deletes a conversation and all its messages
func (s *conversationService) DeleteConversation(ctx context.Context, userID, convID int64) error {
	// Delete conversation (messages will be deleted by CASCADE)
	if err := s.convRepo.Delete(ctx, userID, convID); err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			return ErrConversationNotFound
		}
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	return nil
}

// FavoriteConversation marks a conversation as favorited
func (s *conversationService) FavoriteConversation(ctx context.Context, userID, convID int64) error {
	// Check if conversation exists and belongs to user
	conv, err := s.convRepo.GetByID(ctx, userID, convID)
	if err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			return ErrConversationNotFound
		}
		return fmt.Errorf("failed to get conversation: %w", err)
	}

	// If already favorited, nothing to do
	if conv.IsFavorited {
		return nil
	}

	// Check favorite count limit
	favoriteCount, err := s.convRepo.GetFavoriteCount(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get favorite count: %w", err)
	}

	if favoriteCount >= MaxFavoriteConversations {
		return ErrFavoriteLimitReached
	}

	// Set favorite status
	if err := s.convRepo.SetFavorite(ctx, userID, convID, true); err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			return ErrConversationNotFound
		}
		return fmt.Errorf("failed to favorite conversation: %w", err)
	}

	return nil
}

// UnfavoriteConversation removes favorite status from a conversation
func (s *conversationService) UnfavoriteConversation(ctx context.Context, userID, convID int64) error {
	// Check if conversation exists and belongs to user
	conv, err := s.convRepo.GetByID(ctx, userID, convID)
	if err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			return ErrConversationNotFound
		}
		return fmt.Errorf("failed to get conversation: %w", err)
	}

	// If not favorited, nothing to do
	if !conv.IsFavorited {
		return nil
	}

	// Remove favorite status
	if err := s.convRepo.SetFavorite(ctx, userID, convID, false); err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			return ErrConversationNotFound
		}
		return fmt.Errorf("failed to unfavorite conversation: %w", err)
	}

	return nil
}

// SearchConversations searches conversations by keyword
func (s *conversationService) SearchConversations(ctx context.Context, userID int64, keyword string, filter *model.ConversationFilter) ([]*model.ConversationFlow, int, error) {
	// Set default values
	if filter == nil {
		filter = &model.ConversationFilter{}
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.SortBy == "" {
		filter.SortBy = "updated_at"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	// Trim keyword
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []*model.ConversationFlow{}, 0, nil
	}

	conversations, total, err := s.convRepo.Search(ctx, userID, keyword, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search conversations: %w", err)
	}

	return conversations, total, nil
}

// ExportConversation exports a conversation to JSON
func (s *conversationService) ExportConversation(ctx context.Context, userID, convID int64) ([]byte, error) {
	// Get conversation
	conv, err := s.convRepo.GetByID(ctx, userID, convID)
	if err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			return nil, ErrConversationNotFound
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Get all messages (no pagination for export)
	messages, _, err := s.msgRepo.GetByConversationID(ctx, userID, convID, 1, 10000)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	// Build export structure
	exportData := map[string]interface{}{
		"version":     "1.0",
		"exported_at": time.Now().Format(time.RFC3339),
		"conversations": []map[string]interface{}{
			{
				"id":           conv.ID,
				"title":        conv.Title,
				"is_favorited": conv.IsFavorited,
				"created_at":   conv.CreatedAt.Format(time.RFC3339),
				"updated_at":   conv.UpdatedAt.Format(time.RFC3339),
				"messages":     s.formatMessagesForExport(messages),
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal export data: %w", err)
	}

	return jsonData, nil
}

// ExportConversations exports multiple conversations to JSON
func (s *conversationService) ExportConversations(ctx context.Context, userID int64, convIDs []int64) ([]byte, error) {
	if len(convIDs) == 0 {
		return nil, errors.New("no conversation IDs provided")
	}

	conversations := make([]map[string]interface{}, 0, len(convIDs))

	for _, convID := range convIDs {
		// Get conversation
		conv, err := s.convRepo.GetByID(ctx, userID, convID)
		if err != nil {
			if errors.Is(err, repository.ErrConversationNotFound) {
				// Skip conversations that don't exist or don't belong to user
				continue
			}
			return nil, fmt.Errorf("failed to get conversation %d: %w", convID, err)
		}

		// Get all messages
		messages, _, err := s.msgRepo.GetByConversationID(ctx, userID, convID, 1, 10000)
		if err != nil {
			return nil, fmt.Errorf("failed to get messages for conversation %d: %w", convID, err)
		}

		conversations = append(conversations, map[string]interface{}{
			"id":           conv.ID,
			"title":        conv.Title,
			"is_favorited": conv.IsFavorited,
			"created_at":   conv.CreatedAt.Format(time.RFC3339),
			"updated_at":   conv.UpdatedAt.Format(time.RFC3339),
			"messages":     s.formatMessagesForExport(messages),
		})
	}

	// Build export structure
	exportData := map[string]interface{}{
		"version":       "1.0",
		"exported_at":   time.Now().Format(time.RFC3339),
		"conversations": conversations,
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal export data: %w", err)
	}

	return jsonData, nil
}

// formatMessagesForExport formats messages for export (removes raw request/response)
func (s *conversationService) formatMessagesForExport(messages []*model.Message) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(messages))
	for _, msg := range messages {
		result = append(result, map[string]interface{}{
			"role":       msg.Role,
			"content":    msg.Content,
			"created_at": msg.CreatedAt.Format(time.RFC3339),
		})
	}
	return result
}

// GetMessages retrieves messages for a conversation
func (s *conversationService) GetMessages(ctx context.Context, userID, convID int64, page, pageSize int) ([]*model.Message, int, error) {
	// Verify conversation exists and belongs to user
	_, err := s.convRepo.GetByID(ctx, userID, convID)
	if err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			return nil, 0, ErrConversationNotFound
		}
		return nil, 0, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Set default pagination values
	if pageSize <= 0 {
		pageSize = 50
	}
	if page <= 0 {
		page = 1
	}

	// Get messages
	messages, total, err := s.msgRepo.GetByConversationID(ctx, userID, convID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get messages: %w", err)
	}

	return messages, total, nil
}
