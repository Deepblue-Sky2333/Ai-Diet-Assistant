package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ai-diet-assistant/internal/config"
	"github.com/yourusername/ai-diet-assistant/internal/handler"
	"github.com/yourusername/ai-diet-assistant/internal/middleware"
	"github.com/yourusername/ai-diet-assistant/internal/utils"
	"go.uber.org/zap"
)

// Handlers 包含所有处理器的结构
type Handlers struct {
	Auth       *handler.AuthHandler
	Food       *handler.FoodHandler
	Meal       *handler.MealHandler
	Plan       *handler.PlanHandler
	AI         *handler.AIHandler
	Nutrition  *handler.NutritionHandler
	Dashboard  *handler.DashboardHandler
	Settings   *handler.SettingsHandler
}

// SetupRouter 设置路由
// 
// 文件上传配置说明：
// - 全局请求体大小限制为 10MB
// - 如需添加文件上传端点，使用 middleware.FileValidationMiddleware 进行文件类型验证
// - 示例：
//   uploadConfig := middleware.FileValidationConfig{
//       MaxFileSize: 5 * 1024 * 1024, // 5MB
//       ValidateContent: true,
//   }
//   router.POST("/upload", middleware.FileValidationMiddleware(uploadConfig, logger), handler)
func SetupRouter(cfg *config.Config, logger *zap.Logger, jwtService *utils.JWTService, authService interface {
	ValidateToken(ctx context.Context, token string) (*utils.Claims, error)
}, handlers *Handlers) *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由器
	router := gin.New()

	// 注册全局中间件
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.LoggerMiddleware(logger))
	router.Use(middleware.CORSMiddleware(&cfg.CORS, logger))
	router.Use(middleware.RateLimitMiddleware(&cfg.RateLimit, &cfg.Redis, logger))
	// 全局请求体大小限制（10MB）
	router.Use(middleware.RequestSizeLimitMiddleware(10*1024*1024, logger))
	// 输入清理中间件（防止 XSS 和 SQL 注入）
	router.Use(middleware.SanitizeMiddleware())

	// 健康检查端点（不需要认证）
	router.GET("/health", func(c *gin.Context) {
		utils.Success(c, gin.H{
			"status": "ok",
			"service": "ai-diet-assistant",
		})
	})

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		// 认证路由（不需要认证中间件）
		handlers.Auth.RegisterRoutes(v1)

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
			handlers.Settings.RegisterRoutes(authenticated)
		}
	}

	// 静态文件服务（Next.js 前端）
	// 生产模式：服务构建后的静态文件
	// 开发模式：前端独立运行在 3000 端口
	setupFrontendRoutes(router, cfg, logger)

	return router
}

// setupFrontendRoutes 设置前端路由
func setupFrontendRoutes(router *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	// 检查是否存在构建后的前端文件
	frontendPath := "./web/frontend/.next"
	
	// 如果是生产模式且前端已构建，则服务静态文件
	if cfg.Server.Mode == "release" {
		// 服务 Next.js 静态资源
		router.Static("/_next/static", "./web/frontend/.next/static")
		router.StaticFile("/favicon.ico", "./web/frontend/public/favicon.ico")
		
		// 服务其他静态文件
		router.Static("/static", "./web/frontend/public")
		
		// 所有其他路由返回 index.html（SPA 路由）
		router.NoRoute(func(c *gin.Context) {
			// 如果是 API 请求，返回 404
			if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
				c.JSON(404, gin.H{
					"code":    40401,
					"message": "endpoint not found",
				})
				return
			}
			
			// 否则返回前端页面
			c.File("./web/frontend/out/index.html")
		})
		
		logger.Info("Frontend static files enabled",
			zap.String("path", frontendPath),
			zap.String("mode", "production"),
		)
	} else {
		// 开发模式提示
		logger.Info("Frontend running in development mode",
			zap.String("frontend_url", "http://localhost:3000"),
			zap.String("backend_url", "http://localhost:"+cfg.Server.Port),
		)
	}
}
