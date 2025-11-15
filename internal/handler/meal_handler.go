package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ai-diet-assistant/internal/model"
	"github.com/yourusername/ai-diet-assistant/internal/service"
	"github.com/yourusername/ai-diet-assistant/internal/utils"
)

// MealHandler handles meal-related HTTP requests
type MealHandler struct {
	mealService *service.MealService
}

// NewMealHandler creates a new MealHandler instance
func NewMealHandler(mealService *service.MealService) *MealHandler {
	return &MealHandler{
		mealService: mealService,
	}
}

// CreateMealRequest represents the request body for creating a meal
type CreateMealRequest struct {
	MealDate time.Time        `json:"meal_date" binding:"required"`
	MealType string           `json:"meal_type" binding:"required,oneof=breakfast lunch dinner snack"`
	Foods    []model.MealFood `json:"foods" binding:"required,min=1,max=50,dive"`
	Notes    string           `json:"notes" binding:"omitempty,max=500"`
}

// UpdateMealRequest represents the request body for updating a meal
type UpdateMealRequest struct {
	MealDate time.Time        `json:"meal_date" binding:"required"`
	MealType string           `json:"meal_type" binding:"required,oneof=breakfast lunch dinner snack"`
	Foods    []model.MealFood `json:"foods" binding:"required,min=1,max=50,dive"`
	Notes    string           `json:"notes" binding:"omitempty,max=500"`
}

// CreateMeal handles POST /api/v1/meals
// @Summary Create a new meal record
// @Description Create a new meal record with automatic nutrition calculation
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateMealRequest true "Meal creation request"
// @Success 200 {object} utils.Response{data=model.Meal}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/meals [post]
func (h *MealHandler) CreateMeal(c *gin.Context) {
	var req CreateMealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid request parameters", err))
		return
	}

	// Get user ID from context (injected by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Convert request to model
	meal := &model.Meal{
		MealDate: req.MealDate,
		MealType: req.MealType,
		Foods:    req.Foods,
		Notes:    req.Notes,
	}

	// Create meal (nutrition will be calculated automatically)
	if err := h.mealService.CreateMeal(userID.(int64), meal); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to create meal", err))
		return
	}

	utils.Success(c, meal)
}

// UpdateMeal handles PUT /api/v1/meals/:id
// @Summary Update a meal record
// @Description Update an existing meal record
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Meal ID"
// @Param request body UpdateMealRequest true "Meal update request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/meals/{id} [put]
func (h *MealHandler) UpdateMeal(c *gin.Context) {
	// Parse meal ID
	mealID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid meal id", err))
		return
	}

	var req UpdateMealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid request parameters", err))
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Convert request to model
	meal := &model.Meal{
		MealDate: req.MealDate,
		MealType: req.MealType,
		Foods:    req.Foods,
		Notes:    req.Notes,
	}

	// Update meal (nutrition will be recalculated automatically)
	if err := h.mealService.UpdateMeal(userID.(int64), mealID, meal); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to update meal", err))
		return
	}

	utils.SuccessWithMessage(c, "meal updated successfully", nil)
}

// DeleteMeal handles DELETE /api/v1/meals/:id
// @Summary Delete a meal record
// @Description Delete a meal record
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Meal ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/meals/{id} [delete]
func (h *MealHandler) DeleteMeal(c *gin.Context) {
	// Parse meal ID
	mealID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid meal id", err))
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Delete meal
	if err := h.mealService.DeleteMeal(userID.(int64), mealID); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to delete meal", err))
		return
	}

	utils.SuccessWithMessage(c, "meal deleted successfully", nil)
}

// GetMeal handles GET /api/v1/meals/:id
// @Summary Get a meal record
// @Description Get a meal record by ID
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Meal ID"
// @Success 200 {object} utils.Response{data=model.Meal}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/meals/{id} [get]
func (h *MealHandler) GetMeal(c *gin.Context) {
	// Parse meal ID
	mealID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid meal id", err))
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Get meal
	meal, err := h.mealService.GetMeal(userID.(int64), mealID)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeNotFound, "meal not found", err))
		return
	}

	utils.Success(c, meal)
}

// ListMeals handles GET /api/v1/meals
// @Summary List meal records
// @Description List meal records with filtering and pagination
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Param meal_type query string false "Filter by meal type (breakfast, lunch, dinner, snack)"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} utils.PaginatedResponse{data=[]model.Meal}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/meals [get]
func (h *MealHandler) ListMeals(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Parse query parameters
	filter := &model.MealFilter{
		MealType: c.Query("meal_type"),
	}

	// Parse date filters
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid start_date format", err))
			return
		}
		filter.StartDate = &startDate
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid end_date format", err))
			return
		}
		filter.EndDate = &endDate
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	filter.Page = page
	filter.PageSize = pageSize

	// List meals
	meals, total, err := h.mealService.ListMeals(userID.(int64), filter)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to list meals", err))
		return
	}

	// Calculate pagination
	pagination := utils.CalculatePagination(filter.Page, filter.PageSize, total)

	utils.SuccessWithPagination(c, meals, pagination)
}

// RegisterRoutes registers meal-related routes
func (h *MealHandler) RegisterRoutes(router *gin.RouterGroup) {
	meals := router.Group("/meals")
	{
		meals.POST("", h.CreateMeal)
		meals.PUT("/:id", h.UpdateMeal)
		meals.DELETE("/:id", h.DeleteMeal)
		meals.GET("/:id", h.GetMeal)
		meals.GET("", h.ListMeals)
	}
}
