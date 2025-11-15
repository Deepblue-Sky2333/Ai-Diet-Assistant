package model

import "time"

// Food represents a food item in the user's market panel
type Food struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name" binding:"required,max=100"`
	Category  string    `json:"category" db:"category" binding:"required,oneof=meat vegetable fruit grain other"`
	Price     float64   `json:"price" db:"price" binding:"gte=0"`
	Unit      string    `json:"unit" db:"unit" binding:"required,max=20"`
	Protein   float64   `json:"protein" db:"protein" binding:"gte=0"`
	Carbs     float64   `json:"carbs" db:"carbs" binding:"gte=0"`
	Fat       float64   `json:"fat" db:"fat" binding:"gte=0"`
	Fiber     float64   `json:"fiber" db:"fiber" binding:"gte=0"`
	Calories  float64   `json:"calories" db:"calories" binding:"gte=0"`
	Available bool      `json:"available" db:"available"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// FoodFilter represents filter criteria for listing foods
type FoodFilter struct {
	Category  string
	Available *bool
	Page      int
	PageSize  int
}

// BatchResult represents the result of a batch import operation
type BatchResult struct {
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors,omitempty"`
}
