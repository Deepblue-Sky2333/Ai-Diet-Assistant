package model

import "time"

// Plan represents a meal plan for future dates
type Plan struct {
	ID          int64         `json:"id" db:"id"`
	UserID      int64         `json:"user_id" db:"user_id"`
	PlanDate    time.Time     `json:"plan_date" db:"plan_date" binding:"required"`
	MealType    string        `json:"meal_type" db:"meal_type" binding:"required,oneof=breakfast lunch dinner snack"`
	Foods       []MealFood    `json:"foods" binding:"required,gte=1,dive"`
	Nutrition   NutritionData `json:"nutrition"`
	Status      string        `json:"status" db:"status" binding:"omitempty,oneof=pending completed skipped"`
	AIReasoning string        `json:"ai_reasoning,omitempty" db:"ai_reasoning" binding:"omitempty,max=1000"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
}

// PlanFilter represents filter criteria for listing plans
type PlanFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Status    string
	Page      int
	PageSize  int
}

// GeneratePlanRequest represents a request to generate meal plans
type GeneratePlanRequest struct {
	Days        int    `json:"days" binding:"omitempty,gte=1,lte=7"`
	Preferences string `json:"preferences,omitempty" binding:"omitempty,max=500"`
}
