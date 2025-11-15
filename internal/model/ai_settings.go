package model

import "time"

// AISettings represents AI provider configuration for a user
type AISettings struct {
	ID              int64     `json:"id" db:"id"`
	UserID          int64     `json:"user_id" db:"user_id"`
	Provider        string    `json:"provider" db:"provider" binding:"required,oneof=openai deepseek custom"`
	APIEndpoint     string    `json:"api_endpoint" db:"api_endpoint" binding:"omitempty,url"`
	APIKeyEncrypted string    `json:"-" db:"api_key_encrypted"`
	APIKey          string    `json:"api_key,omitempty" db:"-" binding:"required"`
	Model           string    `json:"model" db:"model" binding:"required"`
	Temperature     float64   `json:"temperature" db:"temperature" binding:"gte=0,lte=2"`
	MaxTokens       int       `json:"max_tokens" db:"max_tokens" binding:"gte=1,lte=4096"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// MaskAPIKey returns a masked version of the API key for safe display
// Shows first 4 and last 4 characters, masks the middle with ****
func (s *AISettings) MaskAPIKey() string {
	if s.APIKey == "" {
		return ""
	}
	if len(s.APIKey) < 8 {
		return "****"
	}
	return s.APIKey[:4] + "****" + s.APIKey[len(s.APIKey)-4:]
}

// ChatHistory represents a chat conversation record
type ChatHistory struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	UserInput  string    `json:"user_input" db:"user_input"`
	AIResponse string    `json:"ai_response" db:"ai_response"`
	Context    string    `json:"context,omitempty" db:"context"` // JSON string
	TokensUsed int       `json:"tokens_used" db:"tokens_used"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
