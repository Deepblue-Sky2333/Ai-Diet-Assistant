package handler

import (
	"errors"
	"strconv"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/service"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
	"github.com/gin-gonic/gin"
)

// MessageHandler handles message-related HTTP requests
type MessageHandler struct {
	messageProxyService service.MessageProxyService
}

// NewMessageHandler creates a new MessageHandler instance
func NewMessageHandler(messageProxyService service.MessageProxyService) *MessageHandler {
	return &MessageHandler{
		messageProxyService: messageProxyService,
	}
}

// SendMessage handles POST /api/v1/conversations/:id/messages
func (h *MessageHandler) SendMessage(c *gin.Context) {
	var req model.SendMessageRequest
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

	// Parse conversation ID
	convID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "invalid conversation ID", err))
		return
	}

	// Validate message size (10MB limit)
	if len(req.Content) > service.MaxMessageSize {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "message too large: maximum 10MB", nil))
		return
	}

	// Send message
	response, err := h.messageProxyService.SendMessage(c.Request.Context(), userID.(int64), convID, req.Content)
	if err != nil {
		if errors.Is(err, service.ErrConversationNotFound) {
			utils.Error(c, utils.NewAppError(utils.CodeNotFound, "conversation not found", err))
			return
		}
		if errors.Is(err, service.ErrMessageTooLarge) {
			utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "message too large: maximum 10MB", err))
			return
		}
		if errors.Is(err, service.ErrAIServiceUnavailable) {
			utils.Error(c, utils.NewAppError(utils.CodeAIServiceError, "AI service unavailable", err))
			return
		}
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to send message", err))
		return
	}

	utils.Success(c, response)
}

// RegisterRoutes registers message-related routes
func (h *MessageHandler) RegisterRoutes(router *gin.RouterGroup) {
	conversations := router.Group("/conversations")
	{
		conversations.POST("/:id/messages", h.SendMessage)
	}
}
