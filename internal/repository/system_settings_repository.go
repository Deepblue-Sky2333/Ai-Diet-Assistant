package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	// ErrSettingNotFound 设置不存在
	ErrSettingNotFound = errors.New("setting not found")
)

// SystemSettingsRepository 系统设置仓储接口
type SystemSettingsRepository interface {
	GetSetting(ctx context.Context, key string) (string, error)
	GetAllSettings(ctx context.Context) (map[string]string, error)
	UpdateSetting(ctx context.Context, key, value string) error
}

// systemSettingsRepository 系统设置仓储实现
type systemSettingsRepository struct {
	db *sql.DB
}

// NewSystemSettingsRepository 创建系统设置仓储实例
func NewSystemSettingsRepository(db *sql.DB) SystemSettingsRepository {
	return &systemSettingsRepository{
		db: db,
	}
}

// GetSetting 获取单个设置
func (r *systemSettingsRepository) GetSetting(ctx context.Context, key string) (string, error) {
	query := `SELECT setting_value FROM system_settings WHERE setting_key = ?`

	var value string
	err := r.db.QueryRowContext(ctx, query, key).Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrSettingNotFound
		}
		return "", fmt.Errorf("failed to get setting: %w", err)
	}

	return value, nil
}

// GetAllSettings 获取所有设置
func (r *systemSettingsRepository) GetAllSettings(ctx context.Context) (map[string]string, error) {
	query := `SELECT setting_key, setting_value FROM system_settings`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all settings: %w", err)
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("failed to scan setting: %w", err)
		}
		settings[key] = value
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating settings: %w", err)
	}

	return settings, nil
}

// UpdateSetting 更新设置（如果不存在则插入）
func (r *systemSettingsRepository) UpdateSetting(ctx context.Context, key, value string) error {
	query := `
		INSERT INTO system_settings (setting_key, setting_value, updated_at)
		VALUES (?, ?, NOW())
		ON DUPLICATE KEY UPDATE setting_value = ?, updated_at = NOW()
	`

	_, err := r.db.ExecContext(ctx, query, key, value, value)
	if err != nil {
		return fmt.Errorf("failed to update setting: %w", err)
	}

	return nil
}
