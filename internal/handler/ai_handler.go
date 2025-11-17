package handler

import (
	"strconv"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/service"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
	"github.com/gin-gonic/gin"
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

	// Parse pagination parameters with validation
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Get chat history
	history, total, err := h.aiService.GetChatHistory(c.Request.Context(), userID.(int64), page, pageSize)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to get chat history", err))
		return
	}

	// Calculate pagination
	pagination := utils.CalculatePagination(page, pageSize, total)

	utils.SuccessWithPagination(c, history, pagination)
}

// RegisterRoutes registers AI-related routes
func (h *AIHandler) RegisterRoutes(router *gin.RouterGroup) {
	ai := router.Group("/ai")
	{
		ai.GET("/history", h.GetChatHistory)
	}
}
