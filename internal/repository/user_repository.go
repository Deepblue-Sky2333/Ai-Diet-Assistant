package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
)

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")
	// ErrUserAlreadyExists 用户已存在
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrInvalidPassword 密码错误
	ErrInvalidPassword = errors.New("invalid password")
)

// UserRepository 用户仓储接口
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User, password string) error
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByID(ctx context.Context, userID int64) (*model.User, error)
	UpdatePassword(ctx context.Context, userID int64, newPasswordHash string) error
	UpdatePasswordWithVersion(ctx context.Context, userID int64, newPasswordHash string, passwordVersion int64) error
}

// userRepository 用户仓储实现
type userRepository struct {
	db *sql.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// CreateUser 创建用户
func (r *userRepository) CreateUser(ctx context.Context, user *model.User, password string) error {
	// 哈希密码
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 使用预编译语句防止 SQL 注入
	query := `
		INSERT INTO users (username, password_hash, password_version, email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	passwordVersion := now.Unix()
	result, err := r.db.ExecContext(ctx, query,
		user.Username,
		passwordHash,
		passwordVersion,
		user.Email,
		now,
		now,
	)

	if err != nil {
		// 检查是否是唯一键冲突
		if isDuplicateKeyError(err) {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	// 获取插入的 ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id
	user.PasswordHash = passwordHash
	user.PasswordVersion = passwordVersion
	user.CreatedAt = now
	user.UpdatedAt = now

	return nil
}

// GetUserByUsername 根据用户名获取用户
func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	// 使用预编译语句防止 SQL 注入
	query := `
		SELECT id, username, password_hash, password_version, email, created_at, updated_at
		FROM users
		WHERE username = ?
	`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.PasswordVersion,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

// GetUserByID 根据用户 ID 获取用户
func (r *userRepository) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	// 使用预编译语句防止 SQL 注入
	query := `
		SELECT id, username, password_hash, password_version, email, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.PasswordVersion,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

// UpdatePassword 更新用户密码
func (r *userRepository) UpdatePassword(ctx context.Context, userID int64, newPasswordHash string) error {
	// 使用预编译语句防止 SQL 注入
	query := `
		UPDATE users
		SET password_hash = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, newPasswordHash, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdatePasswordWithVersion 更新用户密码和密码版本
func (r *userRepository) UpdatePasswordWithVersion(ctx context.Context, userID int64, newPasswordHash string, passwordVersion int64) error {
	// 使用预编译语句防止 SQL 注入
	query := `
		UPDATE users
		SET password_hash = ?, password_version = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, newPasswordHash, passwordVersion, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// isDuplicateKeyError 检查是否是唯一键冲突错误
func isDuplicateKeyError(err error) bool {
	// MySQL 错误码 1062 表示唯一键冲突
	return err != nil && (
		err.Error() == "Error 1062" ||
		contains(err.Error(), "Duplicate entry") ||
		contains(err.Error(), "duplicate key"))
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		containsMiddle(s, substr)))
}

// containsMiddle 检查字符串中间是否包含子串
func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
