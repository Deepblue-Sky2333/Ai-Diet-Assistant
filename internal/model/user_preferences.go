package model

import "time"

// UserPreferences 用户偏好设置
// 采用扁平化结构以匹配前端期望
type UserPreferences struct {
	ID                  int64     `json:"id" db:"id"`
	UserID              int64     `json:"user_id" db:"user_id"`
	TastePreferences    string    `json:"taste_preferences" db:"taste_preferences"`
	DietaryRestrictions string    `json:"dietary_restrictions" db:"dietary_restrictions"`
	DailyCaloriesGoal   int       `json:"daily_calories_goal" db:"daily_calories_goal"`
	DailyProteinGoal    int       `json:"daily_protein_goal" db:"daily_protein_goal"`
	DailyCarbsGoal      int       `json:"daily_carbs_goal" db:"daily_carbs_goal"`
	DailyFatGoal        int       `json:"daily_fat_goal" db:"daily_fat_goal"`
	DailyFiberGoal      int       `json:"daily_fiber_goal" db:"daily_fiber_goal"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// UpdateUserPreferencesRequest 更新用户偏好请求
// 采用扁平化结构，字段类型统一为简单类型
type UpdateUserPreferencesRequest struct {
	TastePreferences    string `json:"taste_preferences" binding:"omitempty,max=500"`
	DietaryRestrictions string `json:"dietary_restrictions" binding:"omitempty,max=500"`
	DailyCaloriesGoal   int    `json:"daily_calories_goal" binding:"omitempty,min=800,max=10000"`
	DailyProteinGoal    int    `json:"daily_protein_goal" binding:"omitempty,min=0,max=500"`
	DailyCarbsGoal      int    `json:"daily_carbs_goal" binding:"omitempty,min=0,max=1000"`
	DailyFatGoal        int    `json:"daily_fat_goal" binding:"omitempty,min=0,max=500"`
	DailyFiberGoal      int    `json:"daily_fiber_goal" binding:"omitempty,min=0,max=200"`
}
