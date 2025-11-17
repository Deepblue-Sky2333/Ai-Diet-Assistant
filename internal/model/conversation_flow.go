package model

import "time"

// ConversationFlow 对话流模型
type ConversationFlow struct {
	ID           int64     `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	Title        string    `json:"title" db:"title"`
	IsFavorited  bool      `json:"is_favorited" db:"is_favorited"`
	MessageCount int       `json:"message_count" db:"message_count"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ConversationFilter 对话流过滤器
type ConversationFilter struct {
	IsFavorited *bool  `json:"is_favorited,omitempty"`
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
	SortBy      string `json:"sort_by"`    // "created_at", "updated_at"
	SortOrder   string `json:"sort_order"` // "asc", "desc"
}

// CreateConversationRequest 创建对话流请求
type CreateConversationRequest struct {
	Title string `json:"title" binding:"omitempty,max=200"`
}

// UpdateConversationTitleRequest 更新对话流标题请求
type UpdateConversationTitleRequest struct {
	Title string `json:"title" binding:"required,max=200"`
}

// ConversationResponse 对话流响应
type ConversationResponse struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Title        string    `json:"title"`
	IsFavorited  bool      `json:"is_favorited"`
	MessageCount int       `json:"message_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ConversationListResponse 对话流列表响应
type ConversationListResponse struct {
	Conversations []*ConversationResponse `json:"conversations"`
	Total         int                     `json:"total"`
	Page          int                     `json:"page"`
	PageSize      int                     `json:"page_size"`
}

// SearchConversationsRequest 搜索对话流请求
type SearchConversationsRequest struct {
	Keyword     string `json:"keyword" binding:"required"`
	IsFavorited *bool  `json:"is_favorited,omitempty"`
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
}

// ExportConversationsRequest 批量导出对话流请求
type ExportConversationsRequest struct {
	ConversationIDs []int64 `json:"conversation_ids" binding:"required,min=1"`
}
