package model

import "time"

// Meal represents a meal record
type Meal struct {
	ID        int64         `json:"id" db:"id"`
	UserID    int64         `json:"user_id" db:"user_id"`
	MealDate  time.Time     `json:"meal_date" db:"meal_date" binding:"required"`
	MealType  string        `json:"meal_type" db:"meal_type" binding:"required,oneof=breakfast lunch dinner snack"`
	Foods     []MealFood    `json:"foods" binding:"required,gte=1,dive"`
	Nutrition NutritionData `json:"nutrition"`
	Notes     string        `json:"notes,omitempty" db:"notes" binding:"omitempty,max=500"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
}

// MealFood represents a food item in a meal
type MealFood struct {
	FoodID int64   `json:"food_id" binding:"required,gt=0"`
	Name   string  `json:"name" binding:"omitempty,min=1,max=100"`
	Amount float64 `json:"amount" binding:"required,gt=0,lte=10000"`
	Unit   string  `json:"unit" binding:"required,min=1,max=20"`
}

// NutritionData represents nutritional information
type NutritionData struct {
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Fiber    float64 `json:"fiber"`
	Calories float64 `json:"calories"`
}

// MealFilter represents filter criteria for listing meals
type MealFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	MealType  string
	Page      int
	PageSize  int
}

// DailyNutritionStats represents daily nutrition statistics
type DailyNutritionStats struct {
	Date      time.Time     `json:"date"`
	Nutrition NutritionData `json:"nutrition"`
	MealCount int           `json:"meal_count"`
}

// MonthlyStats represents monthly meal statistics
type MonthlyStats struct {
	Year       int                    `json:"year"`
	Month      int                    `json:"month"`
	TotalMeals int                    `json:"total_meals"`
	DailyStats []*DailyNutritionStats `json:"daily_stats"`
	AvgDaily   NutritionData          `json:"avg_daily"`
	Total      NutritionData          `json:"total"`
}

// NutritionComparison represents comparison between actual and target nutrition
type NutritionComparison struct {
	Target     NutritionData      `json:"target"`
	Actual     NutritionData      `json:"actual"`
	Difference NutritionData      `json:"difference"`
	Percentage map[string]float64 `json:"percentage"`
}

// DashboardData represents aggregated data for the dashboard view
type DashboardData struct {
	MonthlyStats *MonthlyStats        `json:"monthly_stats"`
	FuturePlans  []*Plan              `json:"future_plans"`
	TodayStats   *DailyNutritionStats `json:"today_stats"`
	CurrentMonth int                  `json:"current_month"`
	CurrentYear  int                  `json:"current_year"`
	GeneratedAt  time.Time            `json:"generated_at"`
}
