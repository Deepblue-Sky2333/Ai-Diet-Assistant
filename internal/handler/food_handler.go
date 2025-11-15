package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/service"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
)

// FoodHandler handles food-related HTTP requests
type FoodHandler struct {
	foodService *service.FoodService
}

// NewFoodHandler creates a new FoodHandler instance
func NewFoodHandler(foodService *service.FoodService) *FoodHandler {
	return &FoodHandler{
		foodService: foodService,
	}
}

// CreateFoodRequest represents the request body for creating a food item
type CreateFoodRequest struct {
	Name      string  `json:"name" binding:"required,min=1,max=100"`
	Category  string  `json:"category" binding:"required,oneof=meat vegetable fruit grain other"`
	Price     float64 `json:"price" binding:"required,gte=0,lte=100000"`
	Unit      string  `json:"unit" binding:"required,min=1,max=20"`
	Protein   float64 `json:"protein" binding:"required,gte=0,lte=1000"`
	Carbs     float64 `json:"carbs" binding:"required,gte=0,lte=1000"`
	Fat       float64 `json:"fat" binding:"required,gte=0,lte=1000"`
	Fiber     float64 `json:"fiber" binding:"required,gte=0,lte=1000"`
	Calories  float64 `json:"calories" binding:"required,gte=0,lte=10000"`
	Available bool    `json:"available"`
}

// UpdateFoodRequest represents the request body for updating a food item
type UpdateFoodRequest struct {
	Name      string  `json:"name" binding:"required,min=1,max=100"`
	Category  string  `json:"category" binding:"required,oneof=meat vegetable fruit grain other"`
	Price     float64 `json:"price" binding:"required,gte=0,lte=100000"`
	Unit      string  `json:"unit" binding:"required,min=1,max=20"`
	Protein   float64 `json:"protein" binding:"required,gte=0,lte=1000"`
	Carbs     float64 `json:"carbs" binding:"required,gte=0,lte=1000"`
	Fat       float64 `json:"fat" binding:"required,gte=0,lte=1000"`
	Fiber     float64 `json:"fiber" binding:"required,gte=0,lte=1000"`
	Calories  float64 `json:"calories" binding:"required,gte=0,lte=10000"`
	Available bool    `json:"available"`
}

// BatchImportRequest represents the request body for batch importing foods
type BatchImportRequest struct {
	Foods []CreateFoodRequest `json:"foods" binding:"required,min=1,max=100,dive"`
}

// CreateFood handles POST /api/v1/foods
// @Summary Create a new food item
// @Description Create a new food item in the user's market panel
// @Tags foods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateFoodRequest true "Food creation request"
// @Success 200 {object} utils.Response{data=model.Food}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/foods [post]
func (h *FoodHandler) CreateFood(c *gin.Context) {
	var req CreateFoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 提供更详细的验证错误信息
		errorMsg := "invalid request parameters: " + err.Error()
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, errorMsg, err))
		return
	}

	// Get user ID from context (injected by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Convert request to model
	food := &model.Food{
		Name:      req.Name,
		Category:  req.Category,
		Price:     req.Price,
		Unit:      req.Unit,
		Protein:   req.Protein,
		Carbs:     req.Carbs,
		Fat:       req.Fat,
		Fiber:     req.Fiber,
		Calories:  req.Calories,
		Available: req.Available,
	}

	// Create food
	if err := h.foodService.CreateFood(userID.(int64), food); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to create food", err))
		return
	}

	utils.Success(c, food)
}

// UpdateFood handles PUT /api/v1/foods/:id
// @Summary Update a food item
// @Description Update an existing food item
// @Tags foods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Food ID"
// @Param request body UpdateFoodRequest true "Food update request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/foods/{id} [put]
func (h *FoodHandler) UpdateFood(c *gin.Context) {
	// Parse food ID
	foodID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid food id", err))
		return
	}

	var req UpdateFoodRequest
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
	food := &model.Food{
		Name:      req.Name,
		Category:  req.Category,
		Price:     req.Price,
		Unit:      req.Unit,
		Protein:   req.Protein,
		Carbs:     req.Carbs,
		Fat:       req.Fat,
		Fiber:     req.Fiber,
		Calories:  req.Calories,
		Available: req.Available,
	}

	// Update food
	if err := h.foodService.UpdateFood(userID.(int64), foodID, food); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to update food", err))
		return
	}

	utils.SuccessWithMessage(c, "food updated successfully", nil)
}

// DeleteFood handles DELETE /api/v1/foods/:id
// @Summary Delete a food item
// @Description Delete a food item from the user's market panel
// @Tags foods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Food ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/foods/{id} [delete]
func (h *FoodHandler) DeleteFood(c *gin.Context) {
	// Parse food ID
	foodID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid food id", err))
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Delete food
	if err := h.foodService.DeleteFood(userID.(int64), foodID); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to delete food", err))
		return
	}

	utils.SuccessWithMessage(c, "food deleted successfully", nil)
}

// GetFood handles GET /api/v1/foods/:id
// @Summary Get a food item
// @Description Get a food item by ID
// @Tags foods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Food ID"
// @Success 200 {object} utils.Response{data=model.Food}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/foods/{id} [get]
func (h *FoodHandler) GetFood(c *gin.Context) {
	// Parse food ID
	foodID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid food id", err))
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Get food
	food, err := h.foodService.GetFood(userID.(int64), foodID)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeNotFound, "food not found", err))
		return
	}

	utils.Success(c, food)
}

// ListFoods handles GET /api/v1/foods
// @Summary List food items
// @Description List food items with filtering and pagination
// @Tags foods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category query string false "Filter by category (meat, vegetable, fruit, grain, other)"
// @Param available query bool false "Filter by availability"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} utils.PaginatedResponse{data=[]model.Food}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/foods [get]
func (h *FoodHandler) ListFoods(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Parse query parameters
	filter := &model.FoodFilter{
		Category: c.Query("category"),
	}

	// Parse available filter
	if availableStr := c.Query("available"); availableStr != "" {
		available := availableStr == "true"
		filter.Available = &available
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	filter.Page = page
	filter.PageSize = pageSize

	// List foods
	foods, total, err := h.foodService.ListFoods(userID.(int64), filter)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to list foods", err))
		return
	}

	// Calculate pagination
	totalPages := (total + filter.PageSize - 1) / filter.PageSize
	pagination := &utils.Pagination{
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	utils.SuccessWithPagination(c, foods, pagination)
}

// BatchImport handles POST /api/v1/foods/batch
// @Summary Batch import food items
// @Description Import multiple food items at once
// @Tags foods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BatchImportRequest true "Batch import request"
// @Success 200 {object} utils.Response{data=model.BatchResult}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/foods/batch [post]
func (h *FoodHandler) BatchImport(c *gin.Context) {
	var req BatchImportRequest
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

	// Convert requests to models
	foods := make([]*model.Food, len(req.Foods))
	for i, foodReq := range req.Foods {
		foods[i] = &model.Food{
			Name:      foodReq.Name,
			Category:  foodReq.Category,
			Price:     foodReq.Price,
			Unit:      foodReq.Unit,
			Protein:   foodReq.Protein,
			Carbs:     foodReq.Carbs,
			Fat:       foodReq.Fat,
			Fiber:     foodReq.Fiber,
			Calories:  foodReq.Calories,
			Available: foodReq.Available,
		}
	}

	// Batch import
	result, err := h.foodService.BatchImport(userID.(int64), foods)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to batch import foods", err))
		return
	}

	utils.Success(c, result)
}

// RegisterRoutes registers food-related routes
func (h *FoodHandler) RegisterRoutes(router *gin.RouterGroup) {
	foods := router.Group("/foods")
	{
		foods.POST("", h.CreateFood)
		foods.PUT("/:id", h.UpdateFood)
		foods.DELETE("/:id", h.DeleteFood)
		foods.GET("/:id", h.GetFood)
		foods.GET("", h.ListFoods)
		foods.POST("/batch", h.BatchImport)
	}
}
