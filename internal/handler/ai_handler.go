package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ai-diet-assistant/internal/service"
	"github.com/yourusername/ai-diet-assistant/internal/utils"
)

// AIHandler handles AI-related HTTP requests
type AIHandler struct {
	aiService *service.AIService
}

// NewAIHandler creates a new AIHandler instance
func NewAIHandler(aiService *service.AIService) *AIHandler {
	return &AIHandler{
		aiService: aiService,
	}
}

// ChatRequest represents the request body for chat
type ChatRequest struct {
	Message string            `json:"message" binding:"required,min=1,max=2000"`
	Context map[string]string `json:"context,omitempty"`
}

// SuggestMealPlanRequest represents the request body for meal plan suggestions
type SuggestMealPlanRequest struct {
	Days           int `json:"days" binding:"required,min=1,max=30"`
	TargetCalories int `json:"target_calories" binding:"omitempty,min=800,max=10000"`
}

// Chat handles POST /api/v1/ai/chat
// @Summary Chat with AI assistant
// @Description Send a message to the AI diet assistant and get a response
// @Tags ai
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChatRequest true "Chat request"
// @Success 200 {object} utils.Response{data=ai.ChatResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/ai/chat [post]
func (h *AIHandler) Chat(c *gin.Context) {
	var req ChatRequest
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

	// Call AI service
	response, err := h.aiService.Chat(c.Request.Context(), userID.(int64), req.Message, req.Context)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "AI chat failed", err))
		return
	}

	// Transform response to match frontend expectations
	// Frontend may expect "response" field, so provide both "message" and "response"
	utils.Success(c, gin.H{
		"message":     response.Message,
		"response":    response.Message,     // Alias for frontend compatibility
		"message_id":  response.MessageID,
		"tokens_used": response.TokensUsed,
	})
}

// SuggestMealPlan handles POST /api/v1/ai/suggest
// @Summary Generate meal plan suggestions
// @Description Generate AI-powered meal plan suggestions based on available foods
// @Tags ai
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SuggestMealPlanRequest true "Meal plan suggestion request"
// @Success 200 {object} utils.Response{data=ai.MealPlanResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/ai/suggest [post]
func (h *AIHandler) SuggestMealPlan(c *gin.Context) {
	var req SuggestMealPlanRequest
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

	// Generate meal plan
	response, err := h.aiService.GenerateMealPlan(c.Request.Context(), userID.(int64), req.Days, req.TargetCalories)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to generate meal plan", err))
		return
	}

	utils.Success(c, response)
}

// GetChatHistory handles GET /api/v1/ai/history
// @Summary Get chat history
// @Description Retrieve chat history with pagination
// @Tags ai
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} utils.PaginatedResponse{data=[]model.ChatHistory}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/ai/history [get]
func (h *AIHandler) GetChatHistory(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Get chat history
	history, total, err := h.aiService.GetChatHistory(c.Request.Context(), userID.(int64), page, pageSize)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to get chat history", err))
		return
	}

	// Calculate pagination
	totalPages := (total + pageSize - 1) / pageSize
	pagination := &utils.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	utils.SuccessWithPagination(c, history, pagination)
}

// RegisterRoutes registers AI-related routes
func (h *AIHandler) RegisterRoutes(router *gin.RouterGroup) {
	ai := router.Group("/ai")
	{
		ai.POST("/chat", h.Chat)
		ai.POST("/suggest", h.SuggestMealPlan)
		ai.GET("/history", h.GetChatHistory)
	}
}
