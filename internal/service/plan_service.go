package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
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
func (s *PlanService) GeneratePlan(ctx context.Context, userID int64, request *model.GeneratePlanRequest) ([]*model.Plan, error) {
	// Validate input
	if err := s.validate.Struct(request); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Set default days if not specified
	days := request.Days
	if days <= 0 {
		days = 2 // Default to 2 days as per requirements
	}

	// Default target calories (can be retrieved from user preferences in the future)
	targetCalories := 2000

	// Call AI service to generate meal plan
	aiResponse, err := s.aiService.GenerateMealPlan(ctx, userID, days, targetCalories)
	if err != nil {
		return nil, fmt.Errorf("failed to generate meal plan with AI: %w", err)
	}

	// Convert AI response to plan records
	plans := make([]*model.Plan, 0)
	now := time.Now()

	for _, aiPlan := range aiResponse.Plans {
		// Parse the date from AI response
		planDate, err := time.Parse("2006-01-02", aiPlan.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to parse plan date %s: %w", aiPlan.Date, err)
		}

		// Convert ai.MealFood to model.MealFood
		foods := make([]model.MealFood, len(aiPlan.Foods))
		for i, f := range aiPlan.Foods {
			foods[i] = model.MealFood{
				FoodID: f.FoodID,
				Name:   f.Name,
				Amount: f.Amount,
				Unit:   f.Unit,
			}
		}

		// Calculate nutrition for this plan
		nutrition, err := s.nutritionService.CalculateNutrition(userID, foods)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate nutrition: %w", err)
		}

		plan := &model.Plan{
			UserID:      userID,
			PlanDate:    planDate,
			MealType:    aiPlan.MealType,
			Foods:       foods,
			Nutrition:   *nutrition,
			Status:      "pending",
			AIReasoning: aiPlan.Reasoning,
		}

		// Create plan in database
		if err := s.planRepo.CreatePlan(plan); err != nil {
			return nil, fmt.Errorf("failed to create plan: %w", err)
		}

		plans = append(plans, plan)
	}

	// Log generation time
	_ = now

	return plans, nil
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
