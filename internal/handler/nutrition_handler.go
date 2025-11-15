package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/service"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
)

// NutritionHandler handles nutrition analysis and statistics HTTP requests
type NutritionHandler struct {
	nutritionService *service.NutritionService
	prefsRepo        repository.UserPreferencesRepository
}

// NewNutritionHandler creates a new NutritionHandler instance
func NewNutritionHandler(nutritionService *service.NutritionService, prefsRepo repository.UserPreferencesRepository) *NutritionHandler {
	return &NutritionHandler{
		nutritionService: nutritionService,
		prefsRepo:        prefsRepo,
	}
}

// GetDailyNutrition handles GET /api/v1/nutrition/daily/:date
// @Summary Get daily nutrition statistics
// @Description Get nutrition statistics for a specific date
// @Tags nutrition
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date path string true "Date in YYYY-MM-DD format"
// @Success 200 {object} utils.Response{data=model.DailyNutritionStats}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/nutrition/daily/{date} [get]
func (h *NutritionHandler) GetDailyNutrition(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Parse date parameter
	dateStr := c.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid date format, expected YYYY-MM-DD", err))
		return
	}

	// Get daily nutrition statistics
	stats, err := h.nutritionService.GetDailyStats(userID.(int64), date)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to get daily nutrition statistics", err))
		return
	}

	utils.Success(c, stats)
}

// GetMonthlyNutrition handles GET /api/v1/nutrition/monthly
// @Summary Get monthly nutrition trend
// @Description Get daily nutrition statistics for an entire month
// @Tags nutrition
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year query int true "Year (e.g., 2024)"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} utils.Response{data=[]model.DailyNutritionStats}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/nutrition/monthly [get]
func (h *NutritionHandler) GetMonthlyNutrition(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Parse year and month parameters
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr == "" || monthStr == "" {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "year and month parameters are required", nil))
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1900 || year > 2100 {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid year parameter", err))
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid month parameter, must be between 1 and 12", err))
		return
	}

	// Get monthly nutrition trend
	dailyStats, err := h.nutritionService.GetMonthlyTrend(userID.(int64), year, month)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to get monthly nutrition trend", err))
		return
	}

	utils.Success(c, dailyStats)
}

// CompareNutrition handles GET /api/v1/nutrition/compare
// @Summary Compare actual nutrition with target
// @Description Compare actual nutrition intake with user's target values for a date range
// @Tags nutrition
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string true "Start date in YYYY-MM-DD format"
// @Param end_date query string true "End date in YYYY-MM-DD format"
// @Success 200 {object} utils.Response{data=model.NutritionComparison}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/nutrition/compare [get]
func (h *NutritionHandler) CompareNutrition(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Parse date parameters - support both single date and date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	singleDateStr := c.Query("date") // Backward compatibility

	var startDate, endDate time.Time
	var err error

	// If start_date and end_date are provided, use date range
	if startDateStr != "" && endDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid start_date format, expected YYYY-MM-DD", err))
			return
		}

		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid end_date format, expected YYYY-MM-DD", err))
			return
		}

		// Validate date range
		if endDate.Before(startDate) {
			utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "end_date must be after start_date", nil))
			return
		}
	} else if singleDateStr != "" {
		// Backward compatibility: single date parameter
		startDate, err = time.Parse("2006-01-02", singleDateStr)
		if err != nil {
			utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid date format, expected YYYY-MM-DD", err))
			return
		}
		endDate = startDate
	} else {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "start_date and end_date parameters are required", nil))
		return
	}

	// Get actual nutrition for the date range
	// For now, we'll aggregate the stats across the date range
	var totalNutrition model.NutritionData
	dayCount := 0

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		stats, err := h.nutritionService.GetDailyStats(userID.(int64), d)
		if err != nil {
			// Skip days with errors (e.g., no data)
			continue
		}

		totalNutrition.Calories += stats.Nutrition.Calories
		totalNutrition.Protein += stats.Nutrition.Protein
		totalNutrition.Carbs += stats.Nutrition.Carbs
		totalNutrition.Fat += stats.Nutrition.Fat
		totalNutrition.Fiber += stats.Nutrition.Fiber
		dayCount++
	}

	// Calculate average if multiple days
	var avgNutrition model.NutritionData
	if dayCount > 0 {
		avgNutrition = model.NutritionData{
			Calories: totalNutrition.Calories / float64(dayCount),
			Protein:  totalNutrition.Protein / float64(dayCount),
			Carbs:    totalNutrition.Carbs / float64(dayCount),
			Fat:      totalNutrition.Fat / float64(dayCount),
			Fiber:    totalNutrition.Fiber / float64(dayCount),
		}
	}

	// Get user preferences to get target nutrition values
	prefs, err := h.prefsRepo.GetPreferences(userID.(int64))
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to get user preferences", err))
		return
	}

	// If no preferences set, use default values
	targetCalories := 2000.0
	targetProtein := 150.0
	targetCarbs := 250.0
	targetFat := 70.0
	targetFiber := 25.0

	if prefs != nil {
		if prefs.DailyCaloriesGoal > 0 {
			targetCalories = float64(prefs.DailyCaloriesGoal)
		}
		if prefs.DailyProteinGoal > 0 {
			targetProtein = float64(prefs.DailyProteinGoal)
		}
		if prefs.DailyCarbsGoal > 0 {
			targetCarbs = float64(prefs.DailyCarbsGoal)
		}
		if prefs.DailyFatGoal > 0 {
			targetFat = float64(prefs.DailyFatGoal)
		}
		if prefs.DailyFiberGoal > 0 {
			targetFiber = float64(prefs.DailyFiberGoal)
		}
	}

	target := &model.NutritionData{
		Calories: targetCalories,
		Protein:  targetProtein,
		Carbs:    targetCarbs,
		Fat:      targetFat,
		Fiber:    targetFiber,
	}

	// Compare actual with target
	comparison, err := h.nutritionService.CompareWithTarget(&avgNutrition, target)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to compare nutrition", err))
		return
	}

	utils.Success(c, comparison)
}

// RegisterRoutes registers nutrition-related routes
func (h *NutritionHandler) RegisterRoutes(router *gin.RouterGroup) {
	nutrition := router.Group("/nutrition")
	{
		nutrition.GET("/daily/:date", h.GetDailyNutrition)
		nutrition.GET("/monthly", h.GetMonthlyNutrition)
		nutrition.GET("/compare", h.CompareNutrition)
	}
}
