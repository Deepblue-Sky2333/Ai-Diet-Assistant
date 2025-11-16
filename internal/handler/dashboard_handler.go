package handler

import (
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/service"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard-related HTTP requests
type DashboardHandler struct {
	dashboardService *service.DashboardService
	prefsRepo        repository.UserPreferencesRepository
}

// NewDashboardHandler creates a new DashboardHandler instance
func NewDashboardHandler(dashboardService *service.DashboardService, prefsRepo repository.UserPreferencesRepository) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
		prefsRepo:        prefsRepo,
	}
}

// GetDashboard handles GET /api/v1/dashboard
// @Summary Get dashboard data
// @Description Get aggregated dashboard data including today's nutrition, nutrition goals, and upcoming plans
// @Tags dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=map[string]interface{}}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/dashboard [get]
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	// Get user ID from context (injected by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Get dashboard data
	dashboardData, err := h.dashboardService.GetDashboardData(userID.(int64))
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to get dashboard data", err))
		return
	}

	// Get user nutrition goals from preferences
	userPrefs, err := h.prefsRepo.GetPreferences(userID.(int64))
	if err != nil {
		// If preferences not found, use default values
		userPrefs = nil
	}

	// Set default nutrition goals if preferences not found
	caloriesGoal := 2000
	proteinGoal := 150
	carbsGoal := 250
	fatGoal := 70

	if userPrefs != nil {
		if userPrefs.DailyCaloriesGoal > 0 {
			caloriesGoal = userPrefs.DailyCaloriesGoal
		}
		if userPrefs.DailyProteinGoal > 0 {
			proteinGoal = userPrefs.DailyProteinGoal
		}
		if userPrefs.DailyCarbsGoal > 0 {
			carbsGoal = userPrefs.DailyCarbsGoal
		}
		if userPrefs.DailyFatGoal > 0 {
			fatGoal = userPrefs.DailyFatGoal
		}
	}

	// Transform plans to match frontend expectations
	transformedPlans := make([]gin.H, len(dashboardData.FuturePlans))
	for i, plan := range dashboardData.FuturePlans {
		transformedPlans[i] = gin.H{
			"id":        plan.ID,
			"date":      plan.PlanDate.Format("2006-01-02"), // Frontend expects "date" field
			"meal_type": plan.MealType,
			"reason":    plan.AIReasoning, // Frontend expects "reason" field
		}
	}

	// Build frontend-expected response structure
	response := gin.H{
		"today_nutrition": gin.H{
			"calories": dashboardData.TodayStats.Nutrition.Calories,
			"protein":  dashboardData.TodayStats.Nutrition.Protein,
			"carbs":    dashboardData.TodayStats.Nutrition.Carbs,
			"fat":      dashboardData.TodayStats.Nutrition.Fat,
		},
		"nutrition_goal": gin.H{
			"calories": caloriesGoal,
			"protein":  proteinGoal,
			"carbs":    carbsGoal,
			"fat":      fatGoal,
		},
		"upcoming_plans": transformedPlans,
	}

	utils.Success(c, response)
}

// RegisterRoutes registers dashboard-related routes
func (h *DashboardHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/dashboard", h.GetDashboard)
}
