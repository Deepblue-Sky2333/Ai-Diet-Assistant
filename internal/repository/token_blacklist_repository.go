package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenBlacklistRepository 令牌黑名单仓库接口
type TokenBlacklistRepository interface {
	// Add 添加令牌到黑名单
	Add(ctx context.Context, token string, expiry time.Duration) error
	// IsBlacklisted 检查令牌是否在黑名单中
	IsBlacklisted(ctx context.Context, token string) (bool, error)
	// Clean 清理过期的令牌（主要用于内存实现）
	Clean(ctx context.Context) error
}

// redisTokenBlacklistRepository Redis 实现的令牌黑名单
type redisTokenBlacklistRepository struct {
	client *redis.Client
	prefix string
}

// NewRedisTokenBlacklistRepository 创建 Redis 令牌黑名单仓库
func NewRedisTokenBlacklistRepository(client *redis.Client) TokenBlacklistRepository {
	return &redisTokenBlacklistRepository{
		client: client,
		prefix: "token:blacklist:",
	}
}

// Add 添加令牌到黑名单
func (r *redisTokenBlacklistRepository) Add(ctx context.Context, token string, expiry time.Duration) error {
	key := r.prefix + token
	err := r.client.Set(ctx, key, "1", expiry).Err()
	if err != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}
	return nil
}

// IsBlacklisted 检查令牌是否在黑名单中
func (r *redisTokenBlacklistRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := r.prefix + token
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check token blacklist: %w", err)
	}
	return result > 0, nil
}

// Clean Redis 自动处理过期，无需手动清理
func (r *redisTokenBlacklistRepository) Clean(ctx context.Context) error {
	return nil
}

// memoryTokenBlacklistRepository 内存实现的令牌黑名单（开发环境备用）
type memoryTokenBlacklistRepository struct {
	mu     sync.RWMutex
	tokens map[string]time.Time // token -> expiry time
}

// NewMemoryTokenBlacklistRepository 创建内存令牌黑名单仓库
func NewMemoryTokenBlacklistRepository() TokenBlacklistRepository {
	repo := &memoryTokenBlacklistRepository{
		tokens: make(map[string]time.Time),
	}

	// 启动后台清理 goroutine
	go repo.cleanupLoop()

	return repo
}

// Add 添加令牌到黑名单
func (m *memoryTokenBlacklistRepository) Add(ctx context.Context, token string, expiry time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	expiryTime := time.Now().Add(expiry)
	m.tokens[token] = expiryTime

	return nil
}

// IsBlacklisted 检查令牌是否在黑名单中
func (m *memoryTokenBlacklistRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	expiryTime, exists := m.tokens[token]
	if !exists {
		return false, nil
	}

	// 检查是否已过期
	if time.Now().After(expiryTime) {
		return false, nil
	}

	return true, nil
}

// Clean 清理过期的令牌
func (m *memoryTokenBlacklistRepository) Clean(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for token, expiryTime := range m.tokens {
		if now.After(expiryTime) {
			delete(m.tokens, token)
		}
	}

	return nil
}

// cleanupLoop 后台清理循环
func (m *memoryTokenBlacklistRepository) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		_ = m.Clean(context.Background())
	}
}
