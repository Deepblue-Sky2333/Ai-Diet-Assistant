package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
)

var (
	// ErrConversationNotFound 对话流不存在
	ErrConversationNotFound = errors.New("conversation not found")
	// ErrFavoriteLimitReached 收藏数量已达上限
	ErrFavoriteLimitReached = errors.New("favorite limit reached")
)

// ConversationRepository 对话流仓储接口
type ConversationRepository interface {
	// Create creates a new conversation flow
	Create(ctx context.Context, conv *model.ConversationFlow) error

	// GetByID retrieves a conversation by ID
	GetByID(ctx context.Context, userID, convID int64) (*model.ConversationFlow, error)

	// List retrieves conversations with pagination and filters
	List(ctx context.Context, userID int64, filter *model.ConversationFilter) ([]*model.ConversationFlow, int, error)

	// Update updates a conversation
	Update(ctx context.Context, conv *model.ConversationFlow) error

	// Delete deletes a conversation and all its messages
	Delete(ctx context.Context, userID, convID int64) error

	// SetFavorite sets the favorite status of a conversation
	SetFavorite(ctx context.Context, userID, convID int64, isFavorited bool) error

	// GetFavoriteCount gets the count of favorited conversations for a user
	GetFavoriteCount(ctx context.Context, userID int64) (int, error)

	// GetRecentCount gets the count of recent (non-favorited) conversations
	GetRecentCount(ctx context.Context, userID int64) (int, error)

	// DeleteOldestRecent deletes the oldest non-favorited conversation
	DeleteOldestRecent(ctx context.Context, userID int64) error

	// Search searches conversations by title
	Search(ctx context.Context, userID int64, keyword string, filter *model.ConversationFilter) ([]*model.ConversationFlow, int, error)

	// IncrementMessageCount increments the message count
	IncrementMessageCount(ctx context.Context, convID int64) error
}

// conversationRepository 对话流仓储实现
type conversationRepository struct {
	db *sql.DB
}

// NewConversationRepository 创建对话流仓储实例
func NewConversationRepository(db *sql.DB) ConversationRepository {
	return &conversationRepository{
		db: db,
	}
}

// Create creates a new conversation flow
func (r *conversationRepository) Create(ctx context.Context, conv *model.ConversationFlow) error {
	query := `
		INSERT INTO conversation_flows (user_id, title, is_favorited, message_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		conv.UserID,
		conv.Title,
		conv.IsFavorited,
		conv.MessageCount,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	conv.ID = id
	conv.CreatedAt = now
	conv.UpdatedAt = now

	return nil
}

// GetByID retrieves a conversation by ID
func (r *conversationRepository) GetByID(ctx context.Context, userID, convID int64) (*model.ConversationFlow, error) {
	query := `
		SELECT id, user_id, title, is_favorited, message_count, created_at, updated_at
		FROM conversation_flows
		WHERE id = ? AND user_id = ?
	`

	conv := &model.ConversationFlow{}
	err := r.db.QueryRowContext(ctx, query, convID, userID).Scan(
		&conv.ID,
		&conv.UserID,
		&conv.Title,
		&conv.IsFavorited,
		&conv.MessageCount,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrConversationNotFound
		}
		return nil, fmt.Errorf("failed to get conversation by id: %w", err)
	}

	return conv, nil
}

// List retrieves conversations with pagination and filters
func (r *conversationRepository) List(ctx context.Context, userID int64, filter *model.ConversationFilter) ([]*model.ConversationFlow, int, error) {
	// Build query with filters
	query := `
		SELECT id, user_id, title, is_favorited, message_count, created_at, updated_at
		FROM conversation_flows
		WHERE user_id = ?
	`
	args := []interface{}{userID}

	// Apply favorite filter if specified
	if filter.IsFavorited != nil {
		query += " AND is_favorited = ?"
		args = append(args, *filter.IsFavorited)
	}

	// Apply sorting
	sortBy := "created_at"
	if filter.SortBy == "updated_at" {
		sortBy = "updated_at"
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Apply pagination
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize
	query += " LIMIT ? OFFSET ?"
	args = append(args, filter.PageSize, offset)

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list conversations: %w", err)
	}
	defer rows.Close()

	// Scan results
	conversations := make([]*model.ConversationFlow, 0)
	for rows.Next() {
		conv := &model.ConversationFlow{}
		err := rows.Scan(
			&conv.ID,
			&conv.UserID,
			&conv.Title,
			&conv.IsFavorited,
			&conv.MessageCount,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan conversation: %w", err)
		}
		conversations = append(conversations, conv)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating conversations: %w", err)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM conversation_flows WHERE user_id = ?`
	countArgs := []interface{}{userID}
	if filter.IsFavorited != nil {
		countQuery += " AND is_favorited = ?"
		countArgs = append(countArgs, *filter.IsFavorited)
	}

	var total int
	err = r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	return conversations, total, nil
}

// Update updates a conversation
func (r *conversationRepository) Update(ctx context.Context, conv *model.ConversationFlow) error {
	query := `
		UPDATE conversation_flows
		SET title = ?, is_favorited = ?, message_count = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		conv.Title,
		conv.IsFavorited,
		conv.MessageCount,
		time.Now(),
		conv.ID,
		conv.UserID,
	)

	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrConversationNotFound
	}

	conv.UpdatedAt = time.Now()
	return nil
}

// Delete deletes a conversation and all its messages
func (r *conversationRepository) Delete(ctx context.Context, userID, convID int64) error {
	query := `DELETE FROM conversation_flows WHERE id = ? AND user_id = ?`

	result, err := r.db.ExecContext(ctx, query, convID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrConversationNotFound
	}

	return nil
}

// SetFavorite sets the favorite status of a conversation
func (r *conversationRepository) SetFavorite(ctx context.Context, userID, convID int64, isFavorited bool) error {
	query := `
		UPDATE conversation_flows
		SET is_favorited = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query, isFavorited, time.Now(), convID, userID)
	if err != nil {
		return fmt.Errorf("failed to set favorite status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrConversationNotFound
	}

	return nil
}

// GetFavoriteCount gets the count of favorited conversations for a user
func (r *conversationRepository) GetFavoriteCount(ctx context.Context, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM conversation_flows WHERE user_id = ? AND is_favorited = true`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get favorite count: %w", err)
	}

	return count, nil
}

// GetRecentCount gets the count of recent (non-favorited) conversations
func (r *conversationRepository) GetRecentCount(ctx context.Context, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM conversation_flows WHERE user_id = ? AND is_favorited = false`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get recent count: %w", err)
	}

	return count, nil
}

// DeleteOldestRecent deletes the oldest non-favorited conversation
func (r *conversationRepository) DeleteOldestRecent(ctx context.Context, userID int64) error {
	query := `
		DELETE FROM conversation_flows
		WHERE id = (
			SELECT id FROM (
				SELECT id FROM conversation_flows
				WHERE user_id = ? AND is_favorited = false
				ORDER BY created_at ASC
				LIMIT 1
			) AS oldest
		)
	`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete oldest recent conversation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrConversationNotFound
	}

	return nil
}

// Search searches conversations by title
func (r *conversationRepository) Search(ctx context.Context, userID int64, keyword string, filter *model.ConversationFilter) ([]*model.ConversationFlow, int, error) {
	// Build query with filters
	query := `
		SELECT id, user_id, title, is_favorited, message_count, created_at, updated_at
		FROM conversation_flows
		WHERE user_id = ? AND title LIKE ?
	`
	args := []interface{}{userID, "%" + keyword + "%"}

	// Apply favorite filter if specified
	if filter.IsFavorited != nil {
		query += " AND is_favorited = ?"
		args = append(args, *filter.IsFavorited)
	}

	// Apply sorting
	sortBy := "created_at"
	if filter.SortBy == "updated_at" {
		sortBy = "updated_at"
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Apply pagination
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize
	query += " LIMIT ? OFFSET ?"
	args = append(args, filter.PageSize, offset)

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search conversations: %w", err)
	}
	defer rows.Close()

	// Scan results
	conversations := make([]*model.ConversationFlow, 0)
	for rows.Next() {
		conv := &model.ConversationFlow{}
		err := rows.Scan(
			&conv.ID,
			&conv.UserID,
			&conv.Title,
			&conv.IsFavorited,
			&conv.MessageCount,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan conversation: %w", err)
		}
		conversations = append(conversations, conv)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating conversations: %w", err)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM conversation_flows WHERE user_id = ? AND title LIKE ?`
	countArgs := []interface{}{userID, "%" + keyword + "%"}
	if filter.IsFavorited != nil {
		countQuery += " AND is_favorited = ?"
		countArgs = append(countArgs, *filter.IsFavorited)
	}

	var total int
	err = r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	return conversations, total, nil
}

// IncrementMessageCount increments the message count
func (r *conversationRepository) IncrementMessageCount(ctx context.Context, convID int64) error {
	query := `
		UPDATE conversation_flows
		SET message_count = message_count + 1, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), convID)
	if err != nil {
		return fmt.Errorf("failed to increment message count: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrConversationNotFound
	}

	return nil
}
