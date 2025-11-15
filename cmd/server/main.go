package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/app"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/config"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
	"go.uber.org/zap"
)

var (
	configPath = flag.String("config", "", "Path to config file (default: ./configs/config.yaml)")
	version    = "1.0.0"
	buildTime  = "unknown"
)

func main() {
	flag.Parse()

	// 打印版本信息
	fmt.Printf("AI Diet Assistant v%s (built at %s)\n", version, buildTime)

	// 初始化默认日志（用于启动阶段）
	logger := utils.InitDefaultLogger()
	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("Starting AI Diet Assistant",
		zap.String("version", version),
		zap.String("build_time", buildTime),
	)

	// 加载配置文件
	logger.Info("Loading configuration...")
	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}
	logger.Info("Configuration loaded successfully")

	// 初始化生产日志系统
	logger, err = utils.InitLogger(&cfg.Log)
	if err != nil {
		logger.Fatal("Failed to initialize logger", zap.Error(err))
	}
	logger.Info("Logger initialized successfully",
		zap.String("level", cfg.Log.Level),
		zap.String("format", cfg.Log.Format),
	)

	// 创建应用程序实例
	logger.Info("Initializing application...")
	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize application", zap.Error(err))
	}
	logger.Info("Application initialized successfully")

	// 运行应用程序
	logger.Info("Starting application server...")
	if err := application.Run(); err != nil {
		logger.Fatal("Application error", zap.Error(err))
	}

	logger.Info("Application stopped gracefully")
	os.Exit(0)
}
