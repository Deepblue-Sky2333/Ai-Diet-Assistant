package service

import (
	"fmt"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
	"github.com/go-playground/validator/v10"
)

// FoodService handles food business logic
type FoodService struct {
	foodRepo *repository.FoodRepository
	validate *validator.Validate
}

// NewFoodService creates a new FoodService instance
func NewFoodService(foodRepo *repository.FoodRepository) *FoodService {
	return &FoodService{
		foodRepo: foodRepo,
		validate: validator.New(),
	}
}

// CreateFood creates a new food item with validation
func (s *FoodService) CreateFood(userID int64, food *model.Food) error {
	// Validate input
	if err := s.validate.Struct(food); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Force user_id to the authenticated user
	food.UserID = userID

	// Set default values
	if food.Unit == "" {
		food.Unit = "g"
	}
	food.Available = true

	return s.foodRepo.CreateFood(food)
}

// UpdateFood updates an existing food item
func (s *FoodService) UpdateFood(userID, foodID int64, food *model.Food) error {
	// Validate input
	if err := s.validate.Struct(food); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Verify the food exists and belongs to the user
	existing, err := s.foodRepo.GetFoodByID(userID, foodID)
	if err != nil {
		return err
	}

	if existing == nil {
		return fmt.Errorf("food not found")
	}

	return s.foodRepo.UpdateFood(userID, foodID, food)
}

// DeleteFood deletes a food item
func (s *FoodService) DeleteFood(userID, foodID int64) error {
	return s.foodRepo.DeleteFood(userID, foodID)
}

// GetFood retrieves a food item by ID
func (s *FoodService) GetFood(userID, foodID int64) (*model.Food, error) {
	return s.foodRepo.GetFoodByID(userID, foodID)
}

// ListFoods retrieves a list of foods with filtering and pagination
func (s *FoodService) ListFoods(userID int64, filter *model.FoodFilter) ([]*model.Food, int, error) {
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

	return s.foodRepo.ListFoods(userID, filter)
}

// BatchImport imports multiple food items with validation
func (s *FoodService) BatchImport(userID int64, foods []*model.Food) (*model.BatchResult, error) {
	result := &model.BatchResult{
		Success: 0,
		Failed:  0,
		Errors:  make([]string, 0),
	}

	if len(foods) == 0 {
		return result, nil
	}

	// Validate each food item
	validFoods := make([]*model.Food, 0)
	for i, food := range foods {
		if err := s.validate.Struct(food); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: %v", i+1, err))
			continue
		}

		// Set default values
		if food.Unit == "" {
			food.Unit = "g"
		}
		food.Available = true
		food.UserID = userID

		validFoods = append(validFoods, food)
	}

	// Batch insert valid foods
	if len(validFoods) > 0 {
		if err := s.foodRepo.BatchInsertFoods(userID, validFoods); err != nil {
			return nil, fmt.Errorf("failed to batch insert foods: %w", err)
		}
		result.Success = len(validFoods)
	}

	return result, nil
}
