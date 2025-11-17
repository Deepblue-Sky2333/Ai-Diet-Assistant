package router

import (
	"context"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/config"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/handler"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/middleware"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handlers 包含所有处理器的结构
type Handlers struct {
	Auth         *handler.AuthHandler
	Food         *handler.FoodHandler
	Meal         *handler.MealHandler
	Plan         *handler.PlanHandler
	AI           *handler.AIHandler
	Nutrition    *handler.NutritionHandler
	Dashboard    *handler.DashboardHandler
	Settings     *handler.SettingsHandler
	Conversation *handler.ConversationHandler
	Message      *handler.MessageHandler
}

// SetupRouter 设置路由
//
// 文件上传配置说明：
//   - 全局请求体大小限制为 10MB
//   - 如需添加文件上传端点，使用 middleware.FileValidationMiddleware 进行文件类型验证
//   - 示例：
//     uploadConfig := middleware.FileValidationConfig{
//     MaxFileSize: 5 * 1024 * 1024, // 5MB
//     ValidateContent: true,
//     }
//     router.POST("/upload", middleware.FileValidationMiddleware(uploadConfig, logger), handler)
func SetupRouter(cfg *config.Config, logger *zap.Logger, jwtService *utils.JWTService, authService interface {
	ValidateToken(ctx context.Context, token string) (*utils.Claims, error)
}, handlers *Handlers, userRepo repository.UserRepository) *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由器
	router := gin.New()

	// 注册全局中间件
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.LoggerMiddleware(logger))
	router.Use(middleware.RateLimitMiddleware(&cfg.RateLimit, &cfg.Redis, logger))
	// 全局请求体大小限制（10MB）
	router.Use(middleware.RequestSizeLimitMiddleware(10*1024*1024, logger))
	// 输入清理中间件（防止 XSS 和 SQL 注入）
	router.Use(middleware.SanitizeMiddleware())

	// 健康检查端点（不需要认证）
	router.GET("/health", func(c *gin.Context) {
		utils.Success(c, gin.H{
			"status":  "ok",
			"service": "ai-diet-assistant",
		})
	})

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		// 认证路由（不需要认证中间件）
		handlers.Auth.RegisterRoutes(v1)

		// 公开的系统信息路由（不需要认证）
		system := v1.Group("/system")
		{
			system.GET("/info", handlers.Settings.GetSystemInfo)
		}

		// 需要认证的路由
		authenticated := v1.Group("")
		authenticated.Use(middleware.AuthMiddleware(jwtService, authService))
		{
			// 食材管理路由
			handlers.Food.RegisterRoutes(authenticated)

			// 餐饮记录路由
			handlers.Meal.RegisterRoutes(authenticated)

			// 饮食计划路由
			handlers.Plan.RegisterRoutes(authenticated)

			// AI 服务路由
			handlers.AI.RegisterRoutes(authenticated)

			// 营养分析路由
			handlers.Nutrition.RegisterRoutes(authenticated)

			// Dashboard 路由
			handlers.Dashboard.RegisterRoutes(authenticated)

			// 设置管理路由
			handlers.Settings.RegisterRoutes(authenticated, userRepo)

			// 对话流管理路由
			handlers.Conversation.RegisterRoutes(authenticated)

			// 消息代理路由
			handlers.Message.RegisterRoutes(authenticated)
		}
	}

	// 404 处理器
	router.NoRoute(func(c *gin.Context) {
		utils.ErrorWithMessage(c, utils.CodeNotFound, "endpoint not found")
	})

	return router
}
