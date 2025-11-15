package repository

import (
	"database/sql"
	"fmt"

	"github.com/yourusername/ai-diet-assistant/internal/model"
)

// UserPreferencesRepository 用户偏好仓储接口
type UserPreferencesRepository interface {
	CreatePreferences(prefs *model.UserPreferences) error
	UpdatePreferences(prefs *model.UserPreferences) error
	GetPreferences(userID int64) (*model.UserPreferences, error)
}

type userPreferencesRepository struct {
	db *sql.DB
}

// NewUserPreferencesRepository 创建用户偏好仓储实例
func NewUserPreferencesRepository(db *sql.DB) UserPreferencesRepository {
	return &userPreferencesRepository{db: db}
}

// CreatePreferences 创建用户偏好（使用扁平化结构）
func (r *userPreferencesRepository) CreatePreferences(prefs *model.UserPreferences) error {
	query := `
		INSERT INTO user_preferences (
			user_id, taste_preferences, dietary_restrictions, 
			daily_calories_goal, daily_protein_goal, daily_carbs_goal,
			daily_fat_goal, daily_fiber_goal
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		prefs.UserID,
		prefs.TastePreferences,
		prefs.DietaryRestrictions,
		prefs.DailyCaloriesGoal,
		prefs.DailyProteinGoal,
		prefs.DailyCarbsGoal,
		prefs.DailyFatGoal,
		prefs.DailyFiberGoal,
	)
	if err != nil {
		return fmt.Errorf("failed to create preferences: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	prefs.ID = id
	return nil
}

// UpdatePreferences 更新用户偏好（使用扁平化结构）
func (r *userPreferencesRepository) UpdatePreferences(prefs *model.UserPreferences) error {
	query := `
		UPDATE user_preferences 
		SET taste_preferences = ?,
		    dietary_restrictions = ?,
		    daily_calories_goal = ?,
		    daily_protein_goal = ?,
		    daily_carbs_goal = ?,
		    daily_fat_goal = ?,
		    daily_fiber_goal = ?,
		    updated_at = CURRENT_TIMESTAMP
		WHERE user_id = ?
	`

	result, err := r.db.Exec(
		query,
		prefs.TastePreferences,
		prefs.DietaryRestrictions,
		prefs.DailyCaloriesGoal,
		prefs.DailyProteinGoal,
		prefs.DailyCarbsGoal,
		prefs.DailyFatGoal,
		prefs.DailyFiberGoal,
		prefs.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update preferences: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("preferences not found for user_id: %d", prefs.UserID)
	}

	return nil
}

// GetPreferences 获取用户偏好（使用扁平化结构）
func (r *userPreferencesRepository) GetPreferences(userID int64) (*model.UserPreferences, error) {
	query := `
		SELECT id, user_id, taste_preferences, dietary_restrictions,
		       daily_calories_goal, daily_protein_goal, daily_carbs_goal,
		       daily_fat_goal, daily_fiber_goal, created_at, updated_at
		FROM user_preferences
		WHERE user_id = ?
	`

	var prefs model.UserPreferences

	err := r.db.QueryRow(query, userID).Scan(
		&prefs.ID,
		&prefs.UserID,
		&prefs.TastePreferences,
		&prefs.DietaryRestrictions,
		&prefs.DailyCaloriesGoal,
		&prefs.DailyProteinGoal,
		&prefs.DailyCarbsGoal,
		&prefs.DailyFatGoal,
		&prefs.DailyFiberGoal,
		&prefs.CreatedAt,
		&prefs.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // 返回 nil 表示未找到，不是错误
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get preferences: %w", err)
	}

	return &prefs, nil
}
