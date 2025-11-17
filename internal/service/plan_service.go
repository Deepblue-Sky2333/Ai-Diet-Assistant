package service

import (
	"context"
	"fmt"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
	"github.com/go-playground/validator/v10"
)

// PlanService handles plan business logic
type PlanService struct {
	planRepo         *repository.PlanRepository
	mealRepo         *repository.MealRepository
	aiService        *AIService
	nutritionService *NutritionService
	validate         *validator.Validate
}

// NewPlanService creates a new PlanService instance
func NewPlanService(
	planRepo *repository.PlanRepository,
	mealRepo *repository.MealRepository,
	aiService *AIService,
	nutritionService *NutritionService,
) *PlanService {
	return &PlanService{
		planRepo:         planRepo,
		mealRepo:         mealRepo,
		aiService:        aiService,
		nutritionService: nutritionService,
		validate:         validator.New(),
	}
}

// GeneratePlan generates meal plans for future days using AI
// NOTE: This method is deprecated and no longer functional after removing AI provider implementations.
// It is kept for reference but will return an error if called.
func (s *PlanService) GeneratePlan(ctx context.Context, userID int64, request *model.GeneratePlanRequest) ([]*model.Plan, error) {
	return nil, fmt.Errorf("AI-based meal plan generation is no longer supported - use the new message proxy service instead")
}

// GetPlan retrieves a plan record by ID
func (s *PlanService) GetPlan(userID, planID int64) (*model.Plan, error) {
	return s.planRepo.GetPlanByID(userID, planID)
}

// ListPlans retrieves a list of plans with filtering and pagination
func (s *PlanService) ListPlans(userID int64, filter *model.PlanFilter) ([]*model.Plan, int, error) {
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

	return s.planRepo.ListPlans(userID, filter)
}

// UpdatePlan updates an existing plan record
func (s *PlanService) UpdatePlan(userID, planID int64, plan *model.Plan) error {
	// Validate input
	if err := s.validate.Struct(plan); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Verify the plan exists and belongs to the user
	existing, err := s.planRepo.GetPlanByID(userID, planID)
	if err != nil {
		return err
	}

	if existing == nil {
		return fmt.Errorf("plan not found")
	}

	// Recalculate nutrition from foods
	nutrition, err := s.nutritionService.CalculateNutrition(userID, plan.Foods)
	if err != nil {
		return fmt.Errorf("failed to calculate nutrition: %w", err)
	}

	plan.Nutrition = *nutrition

	return s.planRepo.UpdatePlan(userID, planID, plan)
}

// DeletePlan deletes a plan record
func (s *PlanService) DeletePlan(userID, planID int64) error {
	return s.planRepo.DeletePlan(userID, planID)
}

// CompletePlan marks a plan as completed and converts it to a meal record
func (s *PlanService) CompletePlan(userID, planID int64) (*model.Meal, error) {
	// Get the plan
	plan, err := s.planRepo.GetPlanByID(userID, planID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}

	if plan.Status == "completed" {
		return nil, fmt.Errorf("plan is already completed")
	}

	// Create a meal record from the plan
	meal := &model.Meal{
		UserID:    userID,
		MealDate:  plan.PlanDate,
		MealType:  plan.MealType,
		Foods:     plan.Foods,
		Nutrition: plan.Nutrition,
		Notes:     fmt.Sprintf("Completed from plan #%d", plan.ID),
	}

	// Create the meal record
	if err := s.mealRepo.CreateMeal(meal); err != nil {
		return nil, fmt.Errorf("failed to create meal from plan: %w", err)
	}

	// Update plan status to completed
	if err := s.planRepo.UpdatePlanStatus(userID, planID, "completed"); err != nil {
		return nil, fmt.Errorf("failed to update plan status: %w", err)
	}

	return meal, nil
}
