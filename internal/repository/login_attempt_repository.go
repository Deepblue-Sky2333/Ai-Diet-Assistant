package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
)

// LoginAttemptRepository 登录尝试仓储接口
type LoginAttemptRepository interface {
	RecordLoginAttempt(ctx context.Context, attempt *model.LoginAttempt) error
	GetRecentFailedAttempts(ctx context.Context, username string, duration time.Duration) (int, error)
	CleanupOldAttempts(ctx context.Context, olderThan time.Duration) error
}

// loginAttemptRepository 登录尝试仓储实现
type loginAttemptRepository struct {
	db *sql.DB
}

// NewLoginAttemptRepository 创建登录尝试仓储实例
func NewLoginAttemptRepository(db *sql.DB) LoginAttemptRepository {
	return &loginAttemptRepository{
		db: db,
	}
}

// RecordLoginAttempt 记录登录尝试
func (r *loginAttemptRepository) RecordLoginAttempt(ctx context.Context, attempt *model.LoginAttempt) error {
	// 使用预编译语句防止 SQL 注入
	query := `
		INSERT INTO login_attempts (username, ip_address, success, attempted_at)
		VALUES (?, ?, ?, ?)
	`

	if attempt.AttemptedAt.IsZero() {
		attempt.AttemptedAt = time.Now()
	}

	result, err := r.db.ExecContext(ctx, query,
		attempt.Username,
		attempt.IPAddress,
		attempt.Success,
		attempt.AttemptedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to record login attempt: %w", err)
	}

	// 获取插入的 ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	attempt.ID = id
	return nil
}

// GetRecentFailedAttempts 获取最近失败的登录尝试次数（默认15分钟内）
func (r *loginAttemptRepository) GetRecentFailedAttempts(ctx context.Context, username string, duration time.Duration) (int, error) {
	// 使用预编译语句防止 SQL 注入
	query := `
		SELECT COUNT(*)
		FROM login_attempts
		WHERE username = ? 
		  AND success = false 
		  AND attempted_at >= ?
	`

	cutoffTime := time.Now().Add(-duration)
	var count int

	err := r.db.QueryRowContext(ctx, query, username, cutoffTime).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get recent failed attempts: %w", err)
	}

	return count, nil
}

// CleanupOldAttempts 清理过期的登录尝试记录
func (r *loginAttemptRepository) CleanupOldAttempts(ctx context.Context, olderThan time.Duration) error {
	// 使用预编译语句防止 SQL 注入
	query := `
		DELETE FROM login_attempts
		WHERE attempted_at < ?
	`

	cutoffTime := time.Now().Add(-olderThan)
	_, err := r.db.ExecContext(ctx, query, cutoffTime)
	if err != nil {
		return fmt.Errorf("failed to cleanup old attempts: %w", err)
	}

	return nil
}
