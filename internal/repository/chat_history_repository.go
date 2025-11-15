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
	// ErrChatHistoryNotFound chat history not found
	ErrChatHistoryNotFound = errors.New("chat history not found")
)

// ChatHistoryRepository handles chat history data operations
type ChatHistoryRepository struct {
	db *sql.DB
}

// NewChatHistoryRepository creates a new ChatHistoryRepository instance
func NewChatHistoryRepository(db *sql.DB) *ChatHistoryRepository {
	return &ChatHistoryRepository{
		db: db,
	}
}

// CreateChatHistory creates a new chat history record
func (r *ChatHistoryRepository) CreateChatHistory(ctx context.Context, history *model.ChatHistory) error {
	query := `
		INSERT INTO chat_history (user_id, user_input, ai_response, context, tokens_used, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		history.UserID,
		history.UserInput,
		history.AIResponse,
		history.Context,
		history.TokensUsed,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create chat history: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	history.ID = id
	history.CreatedAt = now

	return nil
}

// GetChatHistory retrieves chat history for a user with pagination
func (r *ChatHistoryRepository) GetChatHistory(ctx context.Context, userID int64, page, pageSize int) ([]*model.ChatHistory, int, error) {
	// Set default pagination values
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM chat_history
		WHERE user_id = ?
	`
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count chat history: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	query := `
		SELECT id, user_id, user_input, ai_response, context, tokens_used, created_at
		FROM chat_history
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get chat history: %w", err)
	}
	defer rows.Close()

	var historyList []*model.ChatHistory
	for rows.Next() {
		history := &model.ChatHistory{}
		err := rows.Scan(
			&history.ID,
			&history.UserID,
			&history.UserInput,
			&history.AIResponse,
			&history.Context,
			&history.TokensUsed,
			&history.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan chat history: %w", err)
		}
		historyList = append(historyList, history)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating chat history: %w", err)
	}

	return historyList, total, nil
}

// GetChatHistoryByID retrieves a specific chat history record
func (r *ChatHistoryRepository) GetChatHistoryByID(ctx context.Context, userID, historyID int64) (*model.ChatHistory, error) {
	query := `
		SELECT id, user_id, user_input, ai_response, context, tokens_used, created_at
		FROM chat_history
		WHERE id = ? AND user_id = ?
	`

	history := &model.ChatHistory{}
	err := r.db.QueryRowContext(ctx, query, historyID, userID).Scan(
		&history.ID,
		&history.UserID,
		&history.UserInput,
		&history.AIResponse,
		&history.Context,
		&history.TokensUsed,
		&history.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrChatHistoryNotFound
		}
		return nil, fmt.Errorf("failed to get chat history: %w", err)
	}

	return history, nil
}

// DeleteChatHistory deletes a specific chat history record
func (r *ChatHistoryRepository) DeleteChatHistory(ctx context.Context, userID, historyID int64) error {
	query := `
		DELETE FROM chat_history
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query, historyID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete chat history: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrChatHistoryNotFound
	}

	return nil
}

// CleanupOldRecords deletes chat history records older than the specified days
func (r *ChatHistoryRepository) CleanupOldRecords(ctx context.Context, daysToKeep int) (int64, error) {
	if daysToKeep <= 0 {
		daysToKeep = 30 // Default to 30 days
	}

	cutoffDate := time.Now().AddDate(0, 0, -daysToKeep)

	query := `
		DELETE FROM chat_history
		WHERE created_at < ?
	`

	result, err := r.db.ExecContext(ctx, query, cutoffDate)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old records: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// DeleteAllChatHistory deletes all chat history for a user
func (r *ChatHistoryRepository) DeleteAllChatHistory(ctx context.Context, userID int64) error {
	query := `
		DELETE FROM chat_history
		WHERE user_id = ?
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete all chat history: %w", err)
	}

	return nil
}
