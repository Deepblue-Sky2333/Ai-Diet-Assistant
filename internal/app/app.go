package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/config"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/database"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/handler"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/router"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/service"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// App 应用程序结构
type App struct {
	config     *config.Config
	logger     *zap.Logger
	db         *sql.DB
	httpServer *http.Server
	router     *gin.Engine
}

// New 创建应用程序实例
func New(cfg *config.Config, logger *zap.Logger) (*App, error) {
	app := &App{
		config: cfg,
		logger: logger,
	}

	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	app.db = database.GetDB()
	logger.Info("Database initialized successfully")

	// 初始化 Redis（可选）
	if err := database.InitRedis(&cfg.Redis, logger); err != nil {
		logger.Warn("Failed to initialize Redis, using in-memory storage", zap.Error(err))
	}

	// 初始化依赖
	if err := app.initDependencies(); err != nil {
		return nil, fmt.Errorf("failed to initialize dependencies: %w", err)
	}

	// 创建默认用户（如果启用）
	if err := app.ensureDefaultUser(); err != nil {
		logger.Warn("Failed to create default user", zap.Error(err))
	}

	return app, nil
}

// initDependencies 初始化所有依赖（Repository -> Service -> Handler -> Router）
func (a *App) initDependencies() error {
	// 创建 JWT 服务
	jwtService := utils.NewJWTService(
		a.config.JWT.Secret,
		a.config.JWT.ExpireHours,
		a.config.JWT.RefreshExpireHours,
	)

	// 创建加密服务
	cryptoService, err := utils.NewCryptoService(a.config.Encryption.AESKey)
	if err != nil {
		return fmt.Errorf("failed to create crypto service: %w", err)
	}

	// ========== 创建所有 Repository 实例 ==========
	userRepo := repository.NewUserRepository(a.db)
	userPrefsRepo := repository.NewUserPreferencesRepository(a.db)
	foodRepo := repository.NewFoodRepository(a.db)
	mealRepo := repository.NewMealRepository(a.db)
	planRepo := repository.NewPlanRepository(a.db)
	aiSettingsRepo := repository.NewAISettingsRepository(a.db, cryptoService)
	chatHistoryRepo := repository.NewChatHistoryRepository(a.db)
	loginAttemptRepo := repository.NewLoginAttemptRepository(a.db)

	// 创建令牌黑名单仓库（根据 Redis 是否可用选择实现）
	var tokenBlacklistRepo repository.TokenBlacklistRepository
	if database.IsRedisEnabled() {
		tokenBlacklistRepo = repository.NewRedisTokenBlacklistRepository(database.GetRedisClient())
		a.logger.Info("Using Redis for token blacklist")
	} else {
		tokenBlacklistRepo = repository.NewMemoryTokenBlacklistRepository()
		a.logger.Warn("Using in-memory storage for token blacklist (not recommended for production)")
	}

	a.logger.Info("All repositories initialized")

	// ========== 创建所有 Service 实例 ==========
	authService := service.NewAuthService(
		userRepo,
		loginAttemptRepo,
		tokenBlacklistRepo,
		jwtService,
		a.config.Security.MaxLoginAttempts,
		a.config.Security.LockoutDuration,
	)

	foodService := service.NewFoodService(foodRepo)

	nutritionService := service.NewNutritionService(foodRepo, mealRepo)

	mealService := service.NewMealService(mealRepo, nutritionService)

	aiService := service.NewAIService(
		aiSettingsRepo,
		chatHistoryRepo,
		foodRepo,
	)

	planService := service.NewPlanService(
		planRepo,
		mealRepo,
		aiService,
		nutritionService,
	)

	dashboardService := service.NewDashboardService(
		mealService,
		planService,
		nutritionService,
	)

	settingsService := service.NewSettingsService(
		aiSettingsRepo,
		userPrefsRepo,
	)

	a.logger.Info("All services initialized")

	// ========== 创建所有 Handler 实例 ==========
	authHandler := handler.NewAuthHandler(authService)
	foodHandler := handler.NewFoodHandler(foodService)
	mealHandler := handler.NewMealHandler(mealService)
	planHandler := handler.NewPlanHandler(planService)
	aiHandler := handler.NewAIHandler(aiService)
	nutritionHandler := handler.NewNutritionHandler(nutritionService, userPrefsRepo)
	dashboardHandler := handler.NewDashboardHandler(dashboardService, userPrefsRepo)
	settingsHandler := handler.NewSettingsHandler(settingsService)

	a.logger.Info("All handlers initialized")

	// ========== 创建 Handlers 结构 ==========
	handlers := &router.Handlers{
		Auth:      authHandler,
		Food:      foodHandler,
		Meal:      mealHandler,
		Plan:      planHandler,
		AI:        aiHandler,
		Nutrition: nutritionHandler,
		Dashboard: dashboardHandler,
		Settings:  settingsHandler,
	}

	// ========== 设置路由 ==========
	a.router = router.SetupRouter(a.config, a.logger, jwtService, authService, handlers)
	a.logger.Info("Router initialized successfully")

	return nil
}

// Run 启动应用程序
func (a *App) Run() error {
	// 创建 HTTP 服务器
	a.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", a.config.Server.Port),
		Handler:      a.router,
		ReadTimeout:  a.config.Server.ReadTimeout,
		WriteTimeout: a.config.Server.WriteTimeout,
	}

	// 启动服务器（在 goroutine 中）
	go func() {
		a.logger.Info("Starting HTTP server",
			zap.Int("port", a.config.Server.Port),
			zap.String("mode", a.config.Server.Mode),
		)

		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Shutting down server...")

	// 优雅关闭
	return a.Shutdown()
}

// Shutdown 优雅关闭应用程序
func (a *App) Shutdown() error {
	// 创建超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭 HTTP 服务器
	if a.httpServer != nil {
		a.logger.Info("Shutting down HTTP server...")
		if err := a.httpServer.Shutdown(ctx); err != nil {
			a.logger.Error("HTTP server shutdown error", zap.Error(err))
			return err
		}
		a.logger.Info("HTTP server stopped")
	}

	// 关闭数据库连接
	if a.db != nil {
		a.logger.Info("Closing database connection...")
		if err := database.Close(); err != nil {
			a.logger.Error("Database close error", zap.Error(err))
			return err
		}
		a.logger.Info("Database connection closed")
	}

	// 关闭 Redis 连接
	if database.IsRedisEnabled() {
		a.logger.Info("Closing Redis connection...")
		if err := database.CloseRedis(); err != nil {
			a.logger.Error("Redis close error", zap.Error(err))
			return err
		}
		a.logger.Info("Redis connection closed")
	}

	// 同步日志
	if err := a.logger.Sync(); err != nil {
		// 忽略 sync 错误（在某些系统上可能会失败）
		// a.logger.Error("Logger sync error", zap.Error(err))
	}

	a.logger.Info("Application shutdown complete")
	return nil
}

// ensureDefaultUser 确保默认用户存在
func (a *App) ensureDefaultUser() error {
	// 检查是否启用默认用户
	if !a.config.Security.DefaultUser.Enabled {
		a.logger.Info("Default user is disabled")
		return nil
	}

	username := a.config.Security.DefaultUser.Username
	password := a.config.Security.DefaultUser.Password

	if username == "" || password == "" {
		return fmt.Errorf("default user username or password is empty")
	}

	ctx := context.Background()
	userRepo := repository.NewUserRepository(a.db)

	// 检查用户是否已存在
	existingUser, err := userRepo.GetUserByUsername(ctx, username)
	if err == nil && existingUser != nil {
		// 用户已存在，更新密码
		passwordHash, err := utils.HashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// 更新密码和密码版本
		newPasswordVersion := time.Now().Unix()
		if err := userRepo.UpdatePasswordWithVersion(ctx, existingUser.ID, passwordHash, newPasswordVersion); err != nil {
			return fmt.Errorf("failed to update default user password: %w", err)
		}

		a.logger.Info("Default user password updated",
			zap.String("username", username),
		)
		return nil
	}

	// 用户不存在，创建新用户
	user := &model.User{
		Username: username,
		Email:    "",
	}

	if err := userRepo.CreateUser(ctx, user, password); err != nil {
		return fmt.Errorf("failed to create default user: %w", err)
	}

	a.logger.Info("Default user created successfully",
		zap.String("username", username),
	)

	return nil
}
