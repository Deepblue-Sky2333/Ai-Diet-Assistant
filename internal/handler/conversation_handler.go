package handler

import (
	"errors"
	"strconv"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/service"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
	"github.com/gin-gonic/gin"
)

// ConversationHandler handles conversation-related HTTP requests
type ConversationHandler struct {
	conversationService service.ConversationService
}

// NewConversationHandler creates a new ConversationHandler instance
func NewConversationHandler(conversationService service.ConversationService) *ConversationHandler {
	return &ConversationHandler{
		conversationService: conversationService,
	}
}

// CreateConversation handles POST /api/v1/conversations
func (h *ConversationHandler) CreateConversation(c *gin.Context) {
	var req model.CreateConversationRequest
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

	// Create conversation
	conv, err := h.conversationService.CreateConversation(c.Request.Context(), userID.(int64), req.Title)
	if err != nil {
		if errors.Is(err, service.ErrTitleTooLong) {
			utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "title too long: maximum 200 characters", err))
			return
		}
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to create conversation", err))
		return
	}

	utils.SuccessWithMessage(c, "conversation created successfully", conv)
}

// ListConversations handles GET /api/v1/conversations
func (h *ConversationHandler) ListConversations(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort_by", "updated_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// Parse is_favorited filter
	var isFavorited *bool
	if favStr := c.Query("is_favorited"); favStr != "" {
		if favStr == "true" {
			val := true
			isFavorited = &val
		} else if favStr == "false" {
			val := false
			isFavorited = &val
		}
	}

	filter := &model.ConversationFilter{
		IsFavorited: isFavorited,
		Page:        page,
		PageSize:    pageSize,
		SortBy:      sortBy,
		SortOrder:   sortOrder,
	}

	// List conversations
	conversations, total, err := h.conversationService.ListConversations(c.Request.Context(), userID.(int64), filter)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to list conversations", err))
		return
	}

	// Calculate pagination
	pagination := utils.CalculatePagination(page, pageSize, total)

	utils.SuccessWithPagination(c, conversations, pagination)
}

// GetConversation handles GET /api/v1/conversations/:id
func (h *ConversationHandler) GetConversation(c *gin.Context) {
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

	// Get conversation
	conv, err := h.conversationService.GetConversation(c.Request.Context(), userID.(int64), convID)
	if err != nil {
		if errors.Is(err, service.ErrConversationNotFound) {
			utils.Error(c, utils.NewAppError(utils.CodeNotFound, "conversation not found", err))
			return
		}
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to get conversation", err))
		return
	}

	utils.Success(c, conv)
}

// UpdateConversationTitle handles PUT /api/v1/conversations/:id
func (h *ConversationHandler) UpdateConversationTitle(c *gin.Context) {
	var req model.UpdateConversationTitleRequest
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

	// Update conversation title
	err = h.conversationService.UpdateConversationTitle(c.Request.Context(), userID.(int64), convID, req.Title)
	if err != nil {
		if errors.Is(err, service.ErrConversationNotFound) {
			utils.Error(c, utils.NewAppError(utils.CodeNotFound, "conversation not found", err))
			return
		}
		if errors.Is(err, service.ErrTitleTooLong) {
			utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "title too long: maximum 200 characters", err))
			return
		}
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to update conversation title", err))
		return
	}

	utils.SuccessWithMessage(c, "conversation title updated successfully", nil)
}

// DeleteConversation handles DELETE /api/v1/conversations/:id
func (h *ConversationHandler) DeleteConversation(c *gin.Context) {
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

	// Delete conversation
	err = h.conversationService.DeleteConversation(c.Request.Context(), userID.(int64), convID)
	if err != nil {
		if errors.Is(err, service.ErrConversationNotFound) {
			utils.Error(c, utils.NewAppError(utils.CodeNotFound, "conversation not found", err))
			return
		}
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to delete conversation", err))
		return
	}

	utils.SuccessWithMessage(c, "conversation deleted successfully", nil)
}

// FavoriteConversation handles POST /api/v1/conversations/:id/favorite
func (h *ConversationHandler) FavoriteConversation(c *gin.Context) {
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

	// Favorite conversation
	err = h.conversationService.FavoriteConversation(c.Request.Context(), userID.(int64), convID)
	if err != nil {
		if errors.Is(err, service.ErrConversationNotFound) {
			utils.Error(c, utils.NewAppError(utils.CodeNotFound, "conversation not found", err))
			return
		}
		if errors.Is(err, service.ErrFavoriteLimitReached) {
			utils.Error(c, utils.NewAppError(utils.CodeForbidden, "favorite limit reached: maximum 100 favorited conversations", err))
			return
		}
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to favorite conversation", err))
		return
	}

	utils.SuccessWithMessage(c, "conversation favorited successfully", nil)
}

// UnfavoriteConversation handles DELETE /api/v1/conversations/:id/favorite
func (h *ConversationHandler) UnfavoriteConversation(c *gin.Context) {
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

	// Unfavorite conversation
	err = h.conversationService.UnfavoriteConversation(c.Request.Context(), userID.(int64), convID)
	if err != nil {
		if errors.Is(err, service.ErrConversationNotFound) {
			utils.Error(c, utils.NewAppError(utils.CodeNotFound, "conversation not found", err))
			return
		}
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to unfavorite conversation", err))
		return
	}

	utils.SuccessWithMessage(c, "conversation unfavorited successfully", nil)
}

// SearchConversations handles GET /api/v1/conversations/search
func (h *ConversationHandler) SearchConversations(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, utils.NewAppError(utils.CodeUnauthorized, "user not authenticated", nil))
		return
	}

	// Parse query parameters
	keyword := c.Query("keyword")
	if keyword == "" {
		utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "keyword is required", nil))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort_by", "updated_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// Parse is_favorited filter
	var isFavorited *bool
	if favStr := c.Query("is_favorited"); favStr != "" {
		if favStr == "true" {
			val := true
			isFavorited = &val
		} else if favStr == "false" {
			val := false
			isFavorited = &val
		}
	}

	filter := &model.ConversationFilter{
		IsFavorited: isFavorited,
		Page:        page,
		PageSize:    pageSize,
		SortBy:      sortBy,
		SortOrder:   sortOrder,
	}

	// Search conversations
	conversations, total, err := h.conversationService.SearchConversations(c.Request.Context(), userID.(int64), keyword, filter)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to search conversations", err))
		return
	}

	// Calculate pagination
	pagination := utils.CalculatePagination(page, pageSize, total)

	utils.SuccessWithPagination(c, conversations, pagination)
}

// ExportConversation handles GET /api/v1/conversations/:id/export
func (h *ConversationHandler) ExportConversation(c *gin.Context) {
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

	// Export conversation
	jsonData, err := h.conversationService.ExportConversation(c.Request.Context(), userID.(int64), convID)
	if err != nil {
		if errors.Is(err, service.ErrConversationNotFound) {
			utils.Error(c, utils.NewAppError(utils.CodeNotFound, "conversation not found", err))
			return
		}
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to export conversation", err))
		return
	}

	// Return JSON file
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=conversation_"+c.Param("id")+".json")
	c.Data(200, "application/json", jsonData)
}

// ExportConversations handles POST /api/v1/conversations/export
func (h *ConversationHandler) ExportConversations(c *gin.Context) {
	var req model.ExportConversationsRequest
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

	// Export conversations
	jsonData, err := h.conversationService.ExportConversations(c.Request.Context(), userID.(int64), req.ConversationIDs)
	if err != nil {
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to export conversations", err))
		return
	}

	// Return JSON file
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=conversations_export.json")
	c.Data(200, "application/json", jsonData)
}

// GetMessages handles GET /api/v1/conversations/:id/messages
func (h *ConversationHandler) GetMessages(c *gin.Context) {
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

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	// Get messages
	messages, total, err := h.conversationService.GetMessages(c.Request.Context(), userID.(int64), convID, page, pageSize)
	if err != nil {
		if errors.Is(err, service.ErrConversationNotFound) {
			utils.Error(c, utils.NewAppError(utils.CodeNotFound, "conversation not found", err))
			return
		}
		utils.Error(c, utils.NewAppError(utils.CodeInternalError, "failed to get messages", err))
		return
	}

	// Calculate pagination
	pagination := utils.CalculatePagination(page, pageSize, total)

	utils.SuccessWithPagination(c, messages, pagination)
}

// RegisterRoutes registers conversation-related routes
func (h *ConversationHandler) RegisterRoutes(router *gin.RouterGroup) {
	conversations := router.Group("/conversations")
	{
		conversations.POST("", h.CreateConversation)
		conversations.GET("", h.ListConversations)
		conversations.GET("/search", h.SearchConversations)
		conversations.POST("/export", h.ExportConversations)
		conversations.GET("/:id", h.GetConversation)
		conversations.PUT("/:id", h.UpdateConversationTitle)
		conversations.DELETE("/:id", h.DeleteConversation)
		conversations.POST("/:id/favorite", h.FavoriteConversation)
		conversations.DELETE("/:id/favorite", h.UnfavoriteConversation)
		conversations.GET("/:id/export", h.ExportConversation)
		conversations.GET("/:id/messages", h.GetMessages)
	}
}
