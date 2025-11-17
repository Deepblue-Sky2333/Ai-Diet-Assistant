package model

import "time"

// 消息角色常量
const (
	MessageRoleUser      = "user"
	MessageRoleAssistant = "assistant"
)

// Message 消息模型
type Message struct {
	ID             int64     `json:"id" db:"id"`
	ConversationID int64     `json:"conversation_id" db:"conversation_id"`
	Role           string    `json:"role" db:"role"` // "user" or "assistant"
	Content        string    `json:"content" db:"content"`
	RawRequest     string    `json:"raw_request,omitempty" db:"raw_request"`   // 原始请求JSON
	RawResponse    string    `json:"raw_response,omitempty" db:"raw_response"` // 原始响应JSON
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// MessageResponse 消息响应
type MessageResponse struct {
	ID             int64     `json:"id"`
	ConversationID int64     `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

// MessageListResponse 消息列表响应
type MessageListResponse struct {
	Messages []*MessageResponse `json:"messages"`
	Total    int                `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}
