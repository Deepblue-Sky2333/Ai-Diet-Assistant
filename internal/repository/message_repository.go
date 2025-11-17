package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
)

var (
	// ErrMessageNotFound 消息不存在
	ErrMessageNotFound = errors.New("message not found")
)

// MessageRepository 消息仓储接口
type MessageRepository interface {
	// Create creates a new message
	Create(ctx context.Context, msg *model.Message) error

	// GetByConversationID retrieves all messages for a conversation
	GetByConversationID(ctx context.Context, userID, convID int64, page, pageSize int) ([]*model.Message, int, error)

	// DeleteByConversationID deletes all messages for a conversation
	DeleteByConversationID(ctx context.Context, convID int64) error

	// GetByID retrieves a message by ID
	GetByID(ctx context.Context, userID, msgID int64) (*model.Message, error)
}

// messageRepository 消息仓储实现
type messageRepository struct {
	db *sql.DB
}

// NewMessageRepository 创建消息仓储实例
func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{
		db: db,
	}
}

// Create creates a new message
func (r *messageRepository) Create(ctx context.Context, msg *model.Message) error {
	query := `
		INSERT INTO messages (conversation_id, role, content, raw_request, raw_response, created_at)
		VALUES (?, ?, ?, ?, ?, NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		msg.ConversationID,
		msg.Role,
		msg.Content,
		msg.RawRequest,
		msg.RawResponse,
	)

	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	msg.ID = id

	// Get the created_at timestamp
	err = r.db.QueryRowContext(ctx, "SELECT created_at FROM messages WHERE id = ?", id).Scan(&msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to get created_at: %w", err)
	}

	return nil
}

// GetByConversationID retrieves all messages for a conversation
func (r *messageRepository) GetByConversationID(ctx context.Context, userID, convID int64, page, pageSize int) ([]*model.Message, int, error) {
	// First verify the conversation belongs to the user
	var conversationUserID int64
	err := r.db.QueryRowContext(ctx,
		"SELECT user_id FROM conversation_flows WHERE id = ?",
		convID,
	).Scan(&conversationUserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, ErrConversationNotFound
		}
		return nil, 0, fmt.Errorf("failed to verify conversation ownership: %w", err)
	}

	if conversationUserID != userID {
		return nil, 0, ErrConversationNotFound
	}

	// Apply pagination
	if pageSize <= 0 {
		pageSize = 50
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize

	// Get messages
	query := `
		SELECT id, conversation_id, role, content, raw_request, raw_response, created_at
		FROM messages
		WHERE conversation_id = ?
		ORDER BY created_at ASC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, convID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	messages := make([]*model.Message, 0)
	for rows.Next() {
		msg := &model.Message{}
		var rawRequest, rawResponse sql.NullString

		err := rows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.Role,
			&msg.Content,
			&rawRequest,
			&rawResponse,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan message: %w", err)
		}

		if rawRequest.Valid {
			msg.RawRequest = rawRequest.String
		}
		if rawResponse.Valid {
			msg.RawResponse = rawResponse.String
		}

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating messages: %w", err)
	}

	// Get total count
	var total int
	err = r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM messages WHERE conversation_id = ?",
		convID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	return messages, total, nil
}

// DeleteByConversationID deletes all messages for a conversation
func (r *messageRepository) DeleteByConversationID(ctx context.Context, convID int64) error {
	query := `DELETE FROM messages WHERE conversation_id = ?`

	_, err := r.db.ExecContext(ctx, query, convID)
	if err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}

	return nil
}

// GetByID retrieves a message by ID
func (r *messageRepository) GetByID(ctx context.Context, userID, msgID int64) (*model.Message, error) {
	query := `
		SELECT m.id, m.conversation_id, m.role, m.content, m.raw_request, m.raw_response, m.created_at
		FROM messages m
		INNER JOIN conversation_flows c ON m.conversation_id = c.id
		WHERE m.id = ? AND c.user_id = ?
	`

	msg := &model.Message{}
	var rawRequest, rawResponse sql.NullString

	err := r.db.QueryRowContext(ctx, query, msgID, userID).Scan(
		&msg.ID,
		&msg.ConversationID,
		&msg.Role,
		&msg.Content,
		&rawRequest,
		&rawResponse,
		&msg.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMessageNotFound
		}
		return nil, fmt.Errorf("failed to get message by id: %w", err)
	}

	if rawRequest.Valid {
		msg.RawRequest = rawRequest.String
	}
	if rawResponse.Valid {
		msg.RawResponse = rawResponse.String
	}

	return msg, nil
}
