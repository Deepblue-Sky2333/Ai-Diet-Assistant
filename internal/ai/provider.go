package ai

import (
	"context"
	"time"

	"github.com/yourusername/ai-diet-assistant/internal/model"
)

// AIProvider defines the interface for AI service providers
type AIProvider interface {
	// GenerateMealPlan generates meal plans based on available foods and preferences
	GenerateMealPlan(ctx context.Context, request *MealPlanRequest) (*MealPlanResponse, error)
	
	// Chat handles conversational interactions with the AI
	Chat(ctx context.Context, request *ChatRequest) (*ChatResponse, error)
	
	// TestConnection tests the connection to the AI provider
	TestConnection(ctx context.Context) error
}

// MealPlanRequest represents a request to generate meal plans
type MealPlanRequest struct {
	AvailableFoods []model.Food        `json:"available_foods"`
	Preferences    *UserPreferences    `json:"preferences"`
	Days           int                 `json:"days"`
	TargetCalories int                 `json:"target_calories"`
}

// UserPreferences represents user dietary preferences
type UserPreferences struct {
	TastePreferences    []string          `json:"taste_preferences"`
	DietaryRestrictions []string          `json:"dietary_restrictions"`
	DailyCalorieTarget  int               `json:"daily_calorie_target"`
	PreferredMealTimes  map[string]string `json:"preferred_meal_times"`
}

// MealPlanResponse represents the AI's meal plan response
type MealPlanResponse struct {
	Plans []PlannedMeal `json:"plans"`
}

// PlannedMeal represents a single meal in the plan
type PlannedMeal struct {
	Date      string          `json:"date"`
	MealType  string          `json:"meal_type"`
	Foods     []MealFood      `json:"foods"`
	Reasoning string          `json:"reasoning"`
	Nutrition NutritionData   `json:"nutrition"`
}

// MealFood represents a food item in a meal with quantity
type MealFood struct {
	FoodID int64   `json:"food_id"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

// NutritionData represents nutritional information
type NutritionData struct {
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Fiber    float64 `json:"fiber"`
	Calories float64 `json:"calories"`
}

// ChatRequest represents a chat message to the AI
type ChatRequest struct {
	Message     string            `json:"message"`
	Context     map[string]string `json:"context,omitempty"`
	UserID      int64             `json:"user_id"`
	Preferences *UserPreferences  `json:"preferences,omitempty"`
}

// ChatResponse represents the AI's chat response
type ChatResponse struct {
	Message    string    `json:"message"`
	MessageID  int64     `json:"message_id"`
	TokensUsed int       `json:"tokens_used"`
	Timestamp  time.Time `json:"timestamp"`
}

// ProviderConfig represents configuration for an AI provider
type ProviderConfig struct {
	Provider    string  `json:"provider"`
	APIEndpoint string  `json:"api_endpoint"`
	APIKey      string  `json:"api_key"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	Timeout     int     `json:"timeout"` // in seconds
}
