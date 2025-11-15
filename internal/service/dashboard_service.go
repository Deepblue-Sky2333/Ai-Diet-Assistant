package service

import (
	"fmt"
	"time"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
)

// DashboardService handles dashboard data aggregation
type DashboardService struct {
	mealService      *MealService
	planService      *PlanService
	nutritionService *NutritionService
}

// NewDashboardService creates a new DashboardService instance
func NewDashboardService(
	mealService *MealService,
	planService *PlanService,
	nutritionService *NutritionService,
) *DashboardService {
	return &DashboardService{
		mealService:      mealService,
		planService:      planService,
		nutritionService: nutritionService,
	}
}

// GetDashboardData aggregates data for the dashboard view
func (s *DashboardService) GetDashboardData(userID int64) (*model.DashboardData, error) {
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	// Get current month meal records
	monthlyStats, err := s.mealService.GetMonthlyStats(userID, currentYear, currentMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly stats: %w", err)
	}

	// Get future plans (next 2 days)
	tomorrow := now.AddDate(0, 0, 1)
	dayAfterTomorrow := tomorrow.AddDate(0, 0, 1)
	endOfDayAfterTomorrow := time.Date(
		dayAfterTomorrow.Year(),
		dayAfterTomorrow.Month(),
		dayAfterTomorrow.Day(),
		23, 59, 59, 0,
		dayAfterTomorrow.Location(),
	)

	planFilter := &model.PlanFilter{
		StartDate: &tomorrow,
		EndDate:   &endOfDayAfterTomorrow,
		Status:    "pending",
		Page:      1,
		PageSize:  100,
	}

	futurePlans, _, err := s.planService.ListPlans(userID, planFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get future plans: %w", err)
	}

	// Get today's nutrition stats
	todayStats, err := s.nutritionService.GetDailyStats(userID, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's stats: %w", err)
	}

	// Assemble dashboard data
	dashboardData := &model.DashboardData{
		MonthlyStats:  monthlyStats,
		FuturePlans:   futurePlans,
		TodayStats:    todayStats,
		CurrentMonth:  currentMonth,
		CurrentYear:   currentYear,
		GeneratedAt:   now,
	}

	return dashboardData, nil
}
