package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/config"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/database"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RateLimiterInterface 限流器接口
type RateLimiterInterface interface {
	Allow(key string) bool
}

// rateLimitEntry 限流条目（内存实现使用）
type rateLimitEntry struct {
	count     int
	resetTime time.Time
	mu        sync.Mutex
}

// MemoryRateLimiter 内存限流器
type MemoryRateLimiter struct {
	entries         map[string]*rateLimitEntry
	mu              sync.RWMutex
	limit           int
	window          time.Duration
	cleanupInterval time.Duration
}

// RedisRateLimiter Redis限流器（滑动窗口算法）
type RedisRateLimiter struct {
	client      *redis.Client
	limit       int
	window      time.Duration
	fallback    *MemoryRateLimiter
	logger      *zap.Logger
	useFallback bool
	fallbackMu  sync.RWMutex
}

// NewMemoryRateLimiter 创建内存限流器
func NewMemoryRateLimiter(requestsPerMinute int) *MemoryRateLimiter {
	rl := &MemoryRateLimiter{
		entries:         make(map[string]*rateLimitEntry),
		limit:           requestsPerMinute,
		window:          time.Minute,
		cleanupInterval: 5 * time.Minute,
	}

	// 启动清理过期条目的 goroutine
	go rl.cleanup()

	return rl
}

// NewRedisRateLimiter 创建Redis限流器
func NewRedisRateLimiter(client *redis.Client, requestsPerMinute int, logger *zap.Logger) *RedisRateLimiter {
	rl := &RedisRateLimiter{
		client:      client,
		limit:       requestsPerMinute,
		window:      time.Minute,
		fallback:    NewMemoryRateLimiter(requestsPerMinute),
		logger:      logger,
		useFallback: false,
	}

	// 测试Redis连接
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		if logger != nil {
			logger.Warn("Redis connection failed, using memory fallback for rate limiting", zap.Error(err))
		}
		rl.useFallback = true
	}

	return rl
}

// cleanup 定期清理过期的限流条目
func (rl *MemoryRateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		rl.mu.Lock()
		for key, entry := range rl.entries {
			entry.mu.Lock()
			if now.After(entry.resetTime) {
				delete(rl.entries, key)
			}
			entry.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// Allow 检查是否允许请求（内存实现）
func (rl *MemoryRateLimiter) Allow(key string) bool {
	now := time.Now()

	rl.mu.Lock()
	entry, exists := rl.entries[key]
	if !exists {
		entry = &rateLimitEntry{
			count:     0,
			resetTime: now.Add(rl.window),
		}
		rl.entries[key] = entry
	}
	rl.mu.Unlock()

	entry.mu.Lock()
	defer entry.mu.Unlock()

	// 如果窗口已过期，重置计数器
	if now.After(entry.resetTime) {
		entry.count = 0
		entry.resetTime = now.Add(rl.window)
	}

	// 检查是否超过限制
	if entry.count >= rl.limit {
		return false
	}

	entry.count++
	return true
}

// Allow 检查是否允许请求（Redis实现，使用滑动窗口算法）
func (rl *RedisRateLimiter) Allow(key string) bool {
	// 检查是否需要使用降级方案
	rl.fallbackMu.RLock()
	useFallback := rl.useFallback
	rl.fallbackMu.RUnlock()

	if useFallback {
		return rl.fallback.Allow(key)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	now := time.Now()

	// Redis key
	redisKey := fmt.Sprintf("ratelimit:%s", key)

	// 使用 Lua 脚本实现原子操作的滑动窗口算法
	script := redis.NewScript(`
		local key = KEYS[1]
		local now = tonumber(ARGV[1])
		local window = tonumber(ARGV[2])
		local limit = tonumber(ARGV[3])
		local window_start = now - window
		
		-- 删除窗口外的旧记录
		redis.call('ZREMRANGEBYSCORE', key, 0, window_start)
		
		-- 获取当前窗口内的请求数
		local count = redis.call('ZCARD', key)
		
		if count < limit then
			-- 添加当前请求
			redis.call('ZADD', key, now, now)
			-- 设置过期时间（窗口大小 + 1秒）
			redis.call('EXPIRE', key, window + 1)
			return 1
		else
			return 0
		end
	`)

	result, err := script.Run(ctx, rl.client, []string{redisKey},
		now.UnixNano(),
		int64(rl.window.Seconds()),
		rl.limit).Result()

	if err != nil {
		// Redis 出错时降级到内存限流
		if rl.logger != nil {
			rl.logger.Warn("Redis rate limiter error, falling back to memory",
				zap.Error(err),
				zap.String("key", key))
		}

		// 标记使用降级方案
		rl.fallbackMu.Lock()
		rl.useFallback = true
		rl.fallbackMu.Unlock()

		// 启动恢复检查
		go rl.checkRedisRecovery()

		return rl.fallback.Allow(key)
	}

	return result.(int64) == 1
}

// checkRedisRecovery 检查Redis是否恢复
func (rl *RedisRateLimiter) checkRedisRecovery() {
	// 等待一段时间后尝试恢复
	time.Sleep(10 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rl.client.Ping(ctx).Err(); err == nil {
		// Redis 已恢复
		rl.fallbackMu.Lock()
		rl.useFallback = false
		rl.fallbackMu.Unlock()

		if rl.logger != nil {
			rl.logger.Info("Redis connection recovered, switching back from memory fallback")
		}
	}
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(cfg *config.RateLimitConfig, redisCfg *config.RedisConfig, logger *zap.Logger) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// 根据配置选择限流实现
	var limiter RateLimiterInterface
	var storageType string

	// 如果Redis已启用且可用，使用Redis限流器
	if redisCfg.Enabled && database.IsRedisEnabled() {
		redisClient := database.GetRedisClient()
		if redisClient != nil {
			limiter = NewRedisRateLimiter(redisClient, cfg.RequestsPerMinute, logger)
			storageType = "redis"
			if logger != nil {
				logger.Info("Rate limiter initialized with Redis storage")
			}
		} else {
			limiter = NewMemoryRateLimiter(cfg.RequestsPerMinute)
			storageType = "memory"
			if logger != nil {
				logger.Warn("Redis client not available, using memory storage for rate limiting")
			}
		}
	} else {
		limiter = NewMemoryRateLimiter(cfg.RequestsPerMinute)
		storageType = "memory"
		if logger != nil {
			logger.Info("Rate limiter initialized with memory storage")
		}
	}

	return func(c *gin.Context) {
		startTime := time.Now()

		// 获取限流键（优先使用用户 ID，否则使用 IP）
		var key string
		if userID, exists := GetUserID(c); exists {
			key = fmt.Sprintf("user:%d", userID)
		} else {
			key = fmt.Sprintf("ip:%s", c.ClientIP())
		}

		// 检查是否允许请求
		allowed := limiter.Allow(key)

		// 记录性能监控日志
		duration := time.Since(startTime)
		if logger != nil && duration > 100*time.Millisecond {
			logger.Warn("Rate limit check took too long",
				zap.Duration("duration", duration),
				zap.String("storage", storageType),
				zap.String("key", key),
				zap.Bool("allowed", allowed))
		}

		if !allowed {
			if logger != nil {
				logger.Info("Rate limit exceeded",
					zap.String("key", key),
					zap.String("storage", storageType),
					zap.String("ip", c.ClientIP()))
			}

			utils.Error(c, utils.NewAppError(
				utils.CodeTooManyRequests,
				"too many requests, please try again later",
				nil,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}
