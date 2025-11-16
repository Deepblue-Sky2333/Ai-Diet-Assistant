package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var redisClient *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis(cfg *config.RedisConfig, logger *zap.Logger) error {
	if !cfg.Enabled {
		if logger != nil {
			logger.Info("Redis is disabled, using in-memory storage for token blacklist")
		}
		return nil
	}

	// 创建 Redis 客户端
	redisClient = redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:   cfg.Password,
		DB:         cfg.DB,
		MaxRetries: cfg.MaxRetries,
		PoolSize:   cfg.PoolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		// 不直接包装错误，避免暴露连接信息
		return fmt.Errorf("failed to connect to redis: connection test failed")
	}

	if logger != nil {
		logger.Info("Redis connection established successfully")
	}
	return nil
}

// GetRedisClient 获取 Redis 客户端
func GetRedisClient() *redis.Client {
	return redisClient
}

// CloseRedis 关闭 Redis 连接
func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}

// RedisHealthCheck Redis 健康检查
func RedisHealthCheck() error {
	if redisClient == nil {
		return fmt.Errorf("redis not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		// 不直接包装错误，避免暴露连接信息
		return fmt.Errorf("redis health check failed")
	}

	return nil
}

// IsRedisEnabled 检查 Redis 是否已启用
func IsRedisEnabled() bool {
	return redisClient != nil
}
