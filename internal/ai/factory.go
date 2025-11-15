package ai

import (
	"fmt"
	"strings"
)

// NewAIProvider creates a new AI provider based on the configuration
func NewAIProvider(config *ProviderConfig) (AIProvider, error) {
	if config == nil {
		return nil, fmt.Errorf("provider config is required")
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	// Set default values
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 1000
	}
	if config.Timeout == 0 {
		config.Timeout = 30
	}

	provider := strings.ToLower(config.Provider)
	switch provider {
	case "openai":
		if config.APIEndpoint == "" {
			config.APIEndpoint = "https://api.openai.com/v1"
		}
		if config.Model == "" {
			config.Model = "gpt-3.5-turbo"
		}
		return NewOpenAIProvider(config), nil

	case "deepseek":
		if config.APIEndpoint == "" {
			config.APIEndpoint = "https://api.deepseek.com/v1"
		}
		if config.Model == "" {
			config.Model = "deepseek-chat"
		}
		return NewDeepSeekProvider(config), nil

	case "custom":
		if config.APIEndpoint == "" {
			return nil, fmt.Errorf("custom provider requires API endpoint")
		}
		if config.Model == "" {
			config.Model = "default"
		}
		return NewCustomProvider(config), nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s (supported: openai, deepseek, custom)", config.Provider)
	}
}
