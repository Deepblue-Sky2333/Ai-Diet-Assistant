package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
)

// AIProxyClient defines the interface for AI proxy operations
type AIProxyClient interface {
	// SendMessage sends a message to external AI service
	SendMessage(ctx context.Context, request *model.AIProxyRequest) (*model.AIProxyResponse, error)

	// TestConnection tests the connection to external AI service
	TestConnection(ctx context.Context) error
}

// HTTPProxyClient implements AIProxyClient using HTTP
type HTTPProxyClient struct {
	config     *model.AIProxyConfig
	httpClient *http.Client
	logger     *utils.Logger
}

// NewHTTPProxyClient creates a new HTTP proxy client
func NewHTTPProxyClient(config *model.AIProxyConfig, logger *utils.Logger) *HTTPProxyClient {
	return &HTTPProxyClient{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		logger: logger,
	}
}

// SendMessage sends a message to the external AI service with retry logic
func (c *HTTPProxyClient) SendMessage(ctx context.Context, request *model.AIProxyRequest) (*model.AIProxyResponse, error) {
	var lastErr error

	// Set model from config if not specified in request
	if request.Model == "" {
		request.Model = c.config.Model
	}

	// Retry logic
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s, etc.
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			c.logger.Info(fmt.Sprintf("Retrying AI request (attempt %d/%d) after %v", attempt, c.config.MaxRetries, backoff))

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		response, err := c.sendMessageAttempt(ctx, request)
		if err == nil {
			return response, nil
		}

		lastErr = err
		c.logger.Error(fmt.Sprintf("AI request attempt %d failed: %v", attempt+1, err))

		// Don't retry on client errors (4xx)
		if proxyErr, ok := err.(*model.AIProxyError); ok {
			if proxyErr.StatusCode >= 400 && proxyErr.StatusCode < 500 {
				return nil, err
			}
		}
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", c.config.MaxRetries+1, lastErr)
}

// sendMessageAttempt performs a single attempt to send a message
func (c *HTTPProxyClient) sendMessageAttempt(ctx context.Context, request *model.AIProxyRequest) (*model.AIProxyResponse, error) {
	// Marshal request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.APIEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	// Send request
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, &model.AIProxyError{
			StatusCode: 0,
			Message:    "Failed to send request to AI service",
			Details:    err.Error(),
		}
	}
	defer httpResp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if httpResp.StatusCode != http.StatusOK {
		return nil, &model.AIProxyError{
			StatusCode: httpResp.StatusCode,
			Message:    fmt.Sprintf("AI service returned error status: %d", httpResp.StatusCode),
			Details:    string(responseBody),
		}
	}

	// Parse response
	var responseData map[string]interface{}
	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	// Extract content from response
	content, err := c.extractContent(responseData)
	if err != nil {
		return nil, fmt.Errorf("failed to extract content from response: %w", err)
	}

	return &model.AIProxyResponse{
		Content:     content,
		RawResponse: string(responseBody),
	}, nil
}

// extractContent extracts the message content from the AI response
// Supports common response formats from OpenAI, DeepSeek, and similar providers
func (c *HTTPProxyClient) extractContent(responseData map[string]interface{}) (string, error) {
	// Try OpenAI format: choices[0].message.content
	if choices, ok := responseData["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content, nil
				}
			}
		}
	}

	// Try alternative format: response or content field
	if content, ok := responseData["content"].(string); ok {
		return content, nil
	}

	if response, ok := responseData["response"].(string); ok {
		return response, nil
	}

	// Try text field
	if text, ok := responseData["text"].(string); ok {
		return text, nil
	}

	return "", fmt.Errorf("unable to extract content from response: unsupported format")
}

// TestConnection tests the connection to the external AI service
func (c *HTTPProxyClient) TestConnection(ctx context.Context) error {
	// Create a simple test request
	testRequest := &model.AIProxyRequest{
		Messages: []model.AIProxyMessage{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
		Model: c.config.Model,
	}

	// Try to send the test message
	_, err := c.sendMessageAttempt(ctx, testRequest)
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	return nil
}
