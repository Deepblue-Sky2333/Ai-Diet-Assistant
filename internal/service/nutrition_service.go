package service

import (
	"fmt"
	"time"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
)

// NutritionService handles nutrition calculation and analysis
type NutritionService struct {
	foodRepo *repository.FoodRepository
	mealRepo *repository.MealRepository
}

// NewNutritionService creates a new NutritionService instance
func NewNutritionService(foodRepo *repository.FoodRepository, mealRepo *repository.MealRepository) *NutritionService {
	return &NutritionService{
		foodRepo: foodRepo,
		mealRepo: mealRepo,
	}
}

// CalculateNutrition calculates total nutrition from a list of meal foods
func (s *NutritionService) CalculateNutrition(userID int64, foods []model.MealFood) (*model.NutritionData, error) {
	nutrition := &model.NutritionData{
		Protein:  0,
		Carbs:    0,
		Fat:      0,
		Fiber:    0,
		Calories: 0,
	}

	for _, mealFood := range foods {
		// Get food details from database
		food, err := s.foodRepo.GetFoodByID(userID, mealFood.FoodID)
		if err != nil {
			return nil, fmt.Errorf("failed to get food %d: %w", mealFood.FoodID, err)
		}

		// Calculate nutrition based on amount
		// Assuming food nutrition is per 100g and amount is in grams
		ratio := mealFood.Amount / 100.0

		nutrition.Protein += food.Protein * ratio
		nutrition.Carbs += food.Carbs * ratio
		nutrition.Fat += food.Fat * ratio
		nutrition.Fiber += food.Fiber * ratio
		nutrition.Calories += food.Calories * ratio
	}

	return nutrition, nil
}

// GetDailyStats calculates nutrition statistics for a specific day
func (s *NutritionService) GetDailyStats(userID int64, date time.Time) (*model.DailyNutritionStats, error) {
	// Set date to start of day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24*time.Hour - time.Second)

	// Get all meals for the day
	filter := &model.MealFilter{
		StartDate: &startOfDay,
		EndDate:   &endOfDay,
		Page:      1,
		PageSize:  100, // Assume max 100 meals per day
	}

	meals, _, err := s.mealRepo.ListMeals(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily meals: %w", err)
	}

	// Sum up nutrition from all meals
	stats := &model.DailyNutritionStats{
		Date:      startOfDay,
		MealCount: len(meals),
		Nutrition: model.NutritionData{
			Protein:  0,
			Carbs:    0,
			Fat:      0,
			Fiber:    0,
			Calories: 0,
		},
	}

	for _, meal := range meals {
		stats.Nutrition.Protein += meal.Nutrition.Protein
		stats.Nutrition.Carbs += meal.Nutrition.Carbs
		stats.Nutrition.Fat += meal.Nutrition.Fat
		stats.Nutrition.Fiber += meal.Nutrition.Fiber
		stats.Nutrition.Calories += meal.Nutrition.Calories
	}

	return stats, nil
}

// GetMonthlyTrend calculates daily nutrition statistics for an entire month
func (s *NutritionService) GetMonthlyTrend(userID int64, year, month int) ([]*model.DailyNutritionStats, error) {
	// Get all meals for the month
	meals, err := s.mealRepo.GetMonthlyMeals(userID, year, month)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly meals: %w", err)
	}

	// Group meals by date
	dailyMeals := make(map[string][]*model.Meal)
	for _, meal := range meals {
		dateKey := meal.MealDate.Format("2006-01-02")
		dailyMeals[dateKey] = append(dailyMeals[dateKey], meal)
	}

	// Calculate stats for each day
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	var dailyStats []*model.DailyNutritionStats

	for date := startDate; date.Before(endDate); date = date.AddDate(0, 0, 1) {
		dateKey := date.Format("2006-01-02")
		meals := dailyMeals[dateKey]

		stats := &model.DailyNutritionStats{
			Date:      date,
			MealCount: len(meals),
			Nutrition: model.NutritionData{
				Protein:  0,
				Carbs:    0,
				Fat:      0,
				Fiber:    0,
				Calories: 0,
			},
		}

		for _, meal := range meals {
			stats.Nutrition.Protein += meal.Nutrition.Protein
			stats.Nutrition.Carbs += meal.Nutrition.Carbs
			stats.Nutrition.Fat += meal.Nutrition.Fat
			stats.Nutrition.Fiber += meal.Nutrition.Fiber
			stats.Nutrition.Calories += meal.Nutrition.Calories
		}

		dailyStats = append(dailyStats, stats)
	}

	return dailyStats, nil
}

// CompareWithTarget compares actual nutrition with target values
func (s *NutritionService) CompareWithTarget(actual *model.NutritionData, target *model.NutritionData) (*model.NutritionComparison, error) {
	if target == nil {
		return nil, fmt.Errorf("target nutrition data is required")
	}

	comparison := &model.NutritionComparison{
		Target: *target,
		Actual: *actual,
		Difference: model.NutritionData{
			Protein:  actual.Protein - target.Protein,
			Carbs:    actual.Carbs - target.Carbs,
			Fat:      actual.Fat - target.Fat,
			Fiber:    actual.Fiber - target.Fiber,
			Calories: actual.Calories - target.Calories,
		},
		Percentage: make(map[string]float64),
	}

	// Calculate percentage differences
	if target.Protein > 0 {
		comparison.Percentage["protein"] = (actual.Protein / target.Protein) * 100
	}
	if target.Carbs > 0 {
		comparison.Percentage["carbs"] = (actual.Carbs / target.Carbs) * 100
	}
	if target.Fat > 0 {
		comparison.Percentage["fat"] = (actual.Fat / target.Fat) * 100
	}
	if target.Fiber > 0 {
		comparison.Percentage["fiber"] = (actual.Fiber / target.Fiber) * 100
	}
	if target.Calories > 0 {
		comparison.Percentage["calories"] = (actual.Calories / target.Calories) * 100
	}

	return comparison, nil
}
