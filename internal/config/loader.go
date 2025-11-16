package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var globalConfig *Config

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// 设置配置文件路径
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./configs")
		v.AddConfigPath(".")
	}

	// 支持环境变量
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 验证配置
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	globalConfig = &cfg
	return &cfg, nil
}

// Get 获取全局配置
func Get() *Config {
	if globalConfig == nil {
		panic("config not loaded, call Load() first")
	}
	return globalConfig
}

// validateConfig 验证配置
func validateConfig(cfg *Config) error {
	// 验证服务器配置
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	// 验证数据库配置
	if cfg.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if cfg.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if cfg.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	// 验证 JWT 配置
	if cfg.JWT.Secret == "" {
		return fmt.Errorf("jwt secret is required")
	}
	if len(cfg.JWT.Secret) < 32 {
		return fmt.Errorf("jwt secret must be at least 32 characters")
	}

	// 验证加密配置
	if cfg.Encryption.AESKey == "" {
		return fmt.Errorf("aes key is required")
	}
	if len(cfg.Encryption.AESKey) != 32 {
		return fmt.Errorf("aes key must be exactly 32 bytes")
	}

	// 验证 AI 配置（在测试环境下可选）
	// AI API key 可以稍后通过 Web UI 配置
	// if cfg.AI.APIKey == "" {
	// 	return fmt.Errorf("ai api key is required")
	// }

	// 创建必要的目录
	if err := ensureDirectories(cfg); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	return nil
}

// ensureDirectories 确保必要的目录存在
func ensureDirectories(cfg *Config) error {
	dirs := []string{
		cfg.Upload.UploadPath,
	}

	// 如果日志输出到文件，确保日志目录存在
	if cfg.Log.Output != "" && cfg.Log.Output != "stdout" && cfg.Log.Output != "stderr" {
		// 提取目录路径
		logDir := cfg.Log.Output
		if idx := strings.LastIndex(logDir, "/"); idx != -1 {
			logDir = logDir[:idx]
			dirs = append(dirs, logDir)
		}
	}

	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Charset,
		c.ParseTime,
		c.Loc,
	)
}
