package model

import "time"

// AIProxyConfig represents configuration for the AI proxy client
type AIProxyConfig struct {
	APIEndpoint string
	APIKey      string
	Model       string
	Timeout     time.Duration
	MaxRetries  int
}

// AIProxyRequest represents a request to the external AI service
type AIProxyRequest struct {
	Messages []AIProxyMessage `json:"messages"`
	Model    string           `json:"model,omitempty"`
}

// AIProxyMessage represents a single message in the conversation
type AIProxyMessage struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content string `json:"content"`
}

// AIProxyResponse represents a response from the external AI service
type AIProxyResponse struct {
	Content     string
	RawResponse string // Complete raw response JSON
}

// AIProxyError represents an error from the external AI service
type AIProxyError struct {
	StatusCode int
	Message    string
	Details    string
}

func (e *AIProxyError) Error() string {
	return e.Message
}
