package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/service"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
)

// PlanHandler handles plan-related HTTP requests
type PlanHandler struct {
	planService *service.PlanService
}

// NewPlanHandler creates a new PlanHandler instance
func NewPlanHandler(planService *service.PlanService) *PlanHandler {
	return &PlanHandler{
		planService: planService,
	}
}

// GeneratePlanRequest represents the request body for generating plans
type GeneratePlanRequest struct {
	Days        int    `json:"days" binding:"required,min=1,max=7"`
	Preferences string `json:"preferences" binding:"omitempty,max=500"`
}

// UpdatePlanRequest represents the request body for updating a plan
type UpdatePlanRequest struct {
	PlanDate    time.Time        `json:"plan_date" binding:"required"`
	MealType    string           `json:"meal_type" binding:"required,oneof=breakfast lunch dinner snack"`
	Foods       []model.MealFood `json:"foods" binding:"required,min=1,max=50,dive"`
	Status      string           `json:"status" binding:"omitempty,oneof=pending completed skipped"`
	AIReasoning string           `json:"ai_reasoning" binding:"omitempty,max=1000"`
}

// GeneratePlan handles POST /api/v1/plans/generate
// @Summary Generate meal plans using AI
// @Description Generate meal plans for future days based on available foods and preferences
// @Tags plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body GeneratePlanRequest true "Plan generation request"
// @Success 200 {object} utils.Response{data=[]model.Plan}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/plans/generate [post]
func (h *PlanHandler) GeneratePlan(c *gin.Context) {
	var req GeneratePlanRequest
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
	planRequest := &model.GeneratePlanRequest{
		Days:        req.Days,
		Preferences: req.Preferences,
	}

	// Generate plans using AI
	plans, err := h.planService.GeneratePlan(c.Request.Context(), userID.(int64), planRequest)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to generate plans", err))
		return
	}

	utils.Success(c, plans)
}

// GetPlan handles GET /api/v1/plans/:id
// @Summary Get a plan record
// @Description Get a plan record by ID
// @Tags plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Plan ID"
// @Success 200 {object} utils.Response{data=model.Plan}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/plans/{id} [get]
func (h *PlanHandler) GetPlan(c *gin.Context) {
	// Parse plan ID
	planID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid plan id", err))
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Get plan
	plan, err := h.planService.GetPlan(userID.(int64), planID)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeNotFound, "plan not found", err))
		return
	}

	utils.Success(c, plan)
}

// ListPlans handles GET /api/v1/plans
// @Summary List plan records
// @Description List plan records with filtering and pagination
// @Tags plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Param status query string false "Filter by status (pending, completed, skipped)"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} utils.PaginatedResponse{data=[]model.Plan}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/plans [get]
func (h *PlanHandler) ListPlans(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Parse query parameters
	filter := &model.PlanFilter{
		Status: c.Query("status"),
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

	// List plans
	plans, total, err := h.planService.ListPlans(userID.(int64), filter)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to list plans", err))
		return
	}

	// Calculate pagination
	pagination := utils.CalculatePagination(filter.Page, filter.PageSize, total)

	utils.SuccessWithPagination(c, plans, pagination)
}

// UpdatePlan handles PUT /api/v1/plans/:id
// @Summary Update a plan record
// @Description Update an existing plan record
// @Tags plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Plan ID"
// @Param request body UpdatePlanRequest true "Plan update request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/plans/{id} [put]
func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	// Parse plan ID
	planID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid plan id", err))
		return
	}

	var req UpdatePlanRequest
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
	plan := &model.Plan{
		PlanDate:    req.PlanDate,
		MealType:    req.MealType,
		Foods:       req.Foods,
		Status:      req.Status,
		AIReasoning: req.AIReasoning,
	}

	// Update plan (nutrition will be recalculated automatically)
	if err := h.planService.UpdatePlan(userID.(int64), planID, plan); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to update plan", err))
		return
	}

	utils.SuccessWithMessage(c, "plan updated successfully", nil)
}

// DeletePlan handles DELETE /api/v1/plans/:id
// @Summary Delete a plan record
// @Description Delete a plan record
// @Tags plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Plan ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/plans/{id} [delete]
func (h *PlanHandler) DeletePlan(c *gin.Context) {
	// Parse plan ID
	planID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid plan id", err))
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Delete plan
	if err := h.planService.DeletePlan(userID.(int64), planID); err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to delete plan", err))
		return
	}

	utils.SuccessWithMessage(c, "plan deleted successfully", nil)
}

// CompletePlan handles POST /api/v1/plans/:id/complete
// @Summary Complete a plan and convert to meal record
// @Description Mark a plan as completed and create a corresponding meal record
// @Tags plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Plan ID"
// @Success 200 {object} utils.Response{data=model.Meal}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/plans/{id}/complete [post]
func (h *PlanHandler) CompletePlan(c *gin.Context) {
	// Parse plan ID
	planID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid plan id", err))
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Complete plan and create meal record
	meal, err := h.planService.CompletePlan(userID.(int64), planID)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to complete plan", err))
		return
	}

	utils.SuccessWithMessage(c, "plan completed and meal record created", meal)
}

// RegisterRoutes registers plan-related routes
func (h *PlanHandler) RegisterRoutes(router *gin.RouterGroup) {
	plans := router.Group("/plans")
	{
		plans.POST("/generate", h.GeneratePlan)
		plans.GET("/:id", h.GetPlan)
		plans.GET("", h.ListPlans)
		plans.PUT("/:id", h.UpdatePlan)
		plans.DELETE("/:id", h.DeletePlan)
		plans.POST("/:id/complete", h.CompletePlan)
	}
}
