package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/yourusername/ai-diet-assistant/internal/model"
	"github.com/yourusername/ai-diet-assistant/internal/repository"
)

// MealService handles meal business logic
type MealService struct {
	mealRepo         *repository.MealRepository
	nutritionService *NutritionService
	validate         *validator.Validate
}

// NewMealService creates a new MealService instance
func NewMealService(mealRepo *repository.MealRepository, nutritionService *NutritionService) *MealService {
	return &MealService{
		mealRepo:         mealRepo,
		nutritionService: nutritionService,
		validate:         validator.New(),
	}
}

// CreateMeal creates a new meal record with nutrition calculation
func (s *MealService) CreateMeal(userID int64, meal *model.Meal) error {
	// Validate input
	if err := s.validate.Struct(meal); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Force user_id to the authenticated user
	meal.UserID = userID

	// Calculate nutrition from foods
	nutrition, err := s.nutritionService.CalculateNutrition(userID, meal.Foods)
	if err != nil {
		return fmt.Errorf("failed to calculate nutrition: %w", err)
	}

	meal.Nutrition = *nutrition

	// Create meal record
	return s.mealRepo.CreateMeal(meal)
}

// UpdateMeal updates an existing meal record
func (s *MealService) UpdateMeal(userID, mealID int64, meal *model.Meal) error {
	// Validate input
	if err := s.validate.Struct(meal); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Verify the meal exists and belongs to the user
	existing, err := s.mealRepo.GetMealByID(userID, mealID)
	if err != nil {
		return err
	}

	if existing == nil {
		return fmt.Errorf("meal not found")
	}

	// Recalculate nutrition from foods
	nutrition, err := s.nutritionService.CalculateNutrition(userID, meal.Foods)
	if err != nil {
		return fmt.Errorf("failed to calculate nutrition: %w", err)
	}

	meal.Nutrition = *nutrition

	return s.mealRepo.UpdateMeal(userID, mealID, meal)
}

// DeleteMeal deletes a meal record
func (s *MealService) DeleteMeal(userID, mealID int64) error {
	return s.mealRepo.DeleteMeal(userID, mealID)
}

// GetMeal retrieves a meal record by ID
func (s *MealService) GetMeal(userID, mealID int64) (*model.Meal, error) {
	return s.mealRepo.GetMealByID(userID, mealID)
}

// ListMeals retrieves a list of meals with filtering and pagination
func (s *MealService) ListMeals(userID int64, filter *model.MealFilter) ([]*model.Meal, int, error) {
	// Set default pagination values
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	return s.mealRepo.ListMeals(userID, filter)
}

// GetMonthlyStats retrieves monthly meal statistics
func (s *MealService) GetMonthlyStats(userID int64, year, month int) (*model.MonthlyStats, error) {
	// Get all meals for the month
	meals, err := s.mealRepo.GetMonthlyMeals(userID, year, month)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly meals: %w", err)
	}

	// Get daily trend
	dailyStats, err := s.nutritionService.GetMonthlyTrend(userID, year, month)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly trend: %w", err)
	}

	// Calculate total and average nutrition
	totalNutrition := model.NutritionData{
		Protein:  0,
		Carbs:    0,
		Fat:      0,
		Fiber:    0,
		Calories: 0,
	}

	for _, meal := range meals {
		totalNutrition.Protein += meal.Nutrition.Protein
		totalNutrition.Carbs += meal.Nutrition.Carbs
		totalNutrition.Fat += meal.Nutrition.Fat
		totalNutrition.Fiber += meal.Nutrition.Fiber
		totalNutrition.Calories += meal.Nutrition.Calories
	}

	// Calculate average per day (considering all days in the month)
	daysInMonth := len(dailyStats)
	avgNutrition := model.NutritionData{
		Protein:  0,
		Carbs:    0,
		Fat:      0,
		Fiber:    0,
		Calories: 0,
	}

	if daysInMonth > 0 {
		avgNutrition.Protein = totalNutrition.Protein / float64(daysInMonth)
		avgNutrition.Carbs = totalNutrition.Carbs / float64(daysInMonth)
		avgNutrition.Fat = totalNutrition.Fat / float64(daysInMonth)
		avgNutrition.Fiber = totalNutrition.Fiber / float64(daysInMonth)
		avgNutrition.Calories = totalNutrition.Calories / float64(daysInMonth)
	}

	stats := &model.MonthlyStats{
		Year:       year,
		Month:      month,
		TotalMeals: len(meals),
		DailyStats: dailyStats,
		AvgDaily:   avgNutrition,
		Total:      totalNutrition,
	}

	return stats, nil
}
