package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIProvider implements AIProvider for OpenAI API
type OpenAIProvider struct {
	config     *ProviderConfig
	httpClient *http.Client
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(config *ProviderConfig) *OpenAIProvider {
	return &OpenAIProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}
}

// openAIRequest represents a request to OpenAI API
type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
}

// openAIMessage represents a message in the conversation
type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openAIResponse represents a response from OpenAI API
type openAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// GenerateMealPlan generates meal plans using OpenAI
func (p *OpenAIProvider) GenerateMealPlan(ctx context.Context, request *MealPlanRequest) (*MealPlanResponse, error) {
	prompt := p.buildMealPlanPrompt(request)

	messages := []openAIMessage{
		{
			Role:    "system",
			Content: "You are a professional nutritionist and meal planning assistant. Generate meal plans in JSON format based on available foods and user preferences.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := p.callAPI(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	mealPlanResponse, err := p.parseMealPlanResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse meal plan response: %w", err)
	}

	return mealPlanResponse, nil
}

// Chat handles conversational interactions
func (p *OpenAIProvider) Chat(ctx context.Context, request *ChatRequest) (*ChatResponse, error) {
	messages := []openAIMessage{
		{
			Role:    "system",
			Content: "You are a helpful AI diet assistant. Provide friendly and informative responses about nutrition, meal planning, and healthy eating habits.",
		},
		{
			Role:    "user",
			Content: request.Message,
		},
	}

	response, err := p.callAPI(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	return &ChatResponse{
		Message:    response.Choices[0].Message.Content,
		TokensUsed: response.Usage.TotalTokens,
		Timestamp:  time.Now(),
	}, nil
}

// TestConnection tests the connection to OpenAI API
func (p *OpenAIProvider) TestConnection(ctx context.Context) error {
	messages := []openAIMessage{
		{
			Role:    "user",
			Content: "Hello",
		},
	}

	_, err := p.callAPI(ctx, messages)
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	return nil
}

// callAPI makes a request to OpenAI API
func (p *OpenAIProvider) callAPI(ctx context.Context, messages []openAIMessage) (*openAIResponse, error) {
	reqBody := openAIRequest{
		Model:       p.config.Model,
		Messages:    messages,
		Temperature: p.config.Temperature,
		MaxTokens:   p.config.MaxTokens,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/chat/completions", p.config.APIEndpoint)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.APIKey))

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var apiResponse openAIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if apiResponse.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s (type: %s, code: %s)", 
			apiResponse.Error.Message, apiResponse.Error.Type, apiResponse.Error.Code)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return &apiResponse, nil
}

// buildMealPlanPrompt constructs the prompt for meal plan generation
func (p *OpenAIProvider) buildMealPlanPrompt(request *MealPlanRequest) string {
	prompt := fmt.Sprintf(`Generate a %d-day meal plan based on the following information:

Available Foods:
`, request.Days)

	for _, food := range request.AvailableFoods {
		prompt += fmt.Sprintf("- %s (%s): Protein: %.1fg, Carbs: %.1fg, Fat: %.1fg, Calories: %.1f kcal per %s\n",
			food.Name, food.Category, food.Protein, food.Carbs, food.Fat, food.Calories, food.Unit)
	}

	if request.Preferences != nil {
		prompt += "\nUser Preferences:\n"
		if len(request.Preferences.TastePreferences) > 0 {
			prompt += fmt.Sprintf("- Taste preferences: %v\n", request.Preferences.TastePreferences)
		}
		if len(request.Preferences.DietaryRestrictions) > 0 {
			prompt += fmt.Sprintf("- Dietary restrictions: %v\n", request.Preferences.DietaryRestrictions)
		}
		if request.Preferences.DailyCalorieTarget > 0 {
			prompt += fmt.Sprintf("- Daily calorie target: %d kcal\n", request.Preferences.DailyCalorieTarget)
		}
	}

	if request.TargetCalories > 0 {
		prompt += fmt.Sprintf("\nTarget Calories: %d kcal per day\n", request.TargetCalories)
	}

	prompt += `
Please generate a meal plan in the following JSON format:
{
  "plans": [
    {
      "date": "YYYY-MM-DD",
      "meal_type": "breakfast|lunch|dinner|snack",
      "foods": [
        {
          "food_id": 0,
          "name": "food name",
          "amount": 100.0,
          "unit": "g"
        }
      ],
      "reasoning": "Brief explanation of why this meal was chosen",
      "nutrition": {
        "protein": 0.0,
        "carbs": 0.0,
        "fat": 0.0,
        "fiber": 0.0,
        "calories": 0.0
      }
    }
  ]
}

Generate meals for breakfast, lunch, and dinner for each day. Ensure nutritional balance and variety.`

	return prompt
}

// parseMealPlanResponse parses the AI response into a MealPlanResponse
func (p *OpenAIProvider) parseMealPlanResponse(response *openAIResponse) (*MealPlanResponse, error) {
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response choices available")
	}

	content := response.Choices[0].Message.Content

	var mealPlanResponse MealPlanResponse
	if err := json.Unmarshal([]byte(content), &mealPlanResponse); err != nil {
		// Try to extract JSON from markdown code blocks
		content = extractJSONFromMarkdown(content)
		if err := json.Unmarshal([]byte(content), &mealPlanResponse); err != nil {
			return nil, fmt.Errorf("failed to parse meal plan JSON: %w", err)
		}
	}

	return &mealPlanResponse, nil
}

// extractJSONFromMarkdown extracts JSON content from markdown code blocks
func extractJSONFromMarkdown(content string) string {
	// Simple extraction: look for ```json ... ``` or ``` ... ```
	start := -1
	end := -1

	// Find ```json or ```
	if idx := bytes.Index([]byte(content), []byte("```json")); idx != -1 {
		start = idx + 7
	} else if idx := bytes.Index([]byte(content), []byte("```")); idx != -1 {
		start = idx + 3
	}

	if start != -1 {
		// Find closing ```
		if idx := bytes.Index([]byte(content[start:]), []byte("```")); idx != -1 {
			end = start + idx
			return content[start:end]
		}
	}

	return content
}
