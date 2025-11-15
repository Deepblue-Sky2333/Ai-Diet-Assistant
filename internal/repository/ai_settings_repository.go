package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/yourusername/ai-diet-assistant/internal/model"
	"github.com/yourusername/ai-diet-assistant/internal/utils"
)

var (
	// ErrAISettingsNotFound AI settings not found
	ErrAISettingsNotFound = errors.New("ai settings not found")
)

// AISettingsRepository handles AI settings data operations
type AISettingsRepository struct {
	db     *sql.DB
	crypto *utils.CryptoService
}

// NewAISettingsRepository creates a new AISettingsRepository instance
func NewAISettingsRepository(db *sql.DB, crypto *utils.CryptoService) *AISettingsRepository {
	return &AISettingsRepository{
		db:     db,
		crypto: crypto,
	}
}

// CreateAISettings creates new AI settings with encrypted API key
func (r *AISettingsRepository) CreateAISettings(ctx context.Context, settings *model.AISettings) error {
	// Encrypt API key
	encryptedKey, err := r.crypto.EncryptAES(settings.APIKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt API key: %w", err)
	}

	// Deactivate other settings if this one is active
	if settings.IsActive {
		if err := r.deactivateAllSettings(ctx, settings.UserID); err != nil {
			return fmt.Errorf("failed to deactivate other settings: %w", err)
		}
	}

	query := `
		INSERT INTO ai_settings (user_id, provider, api_endpoint, api_key_encrypted, model, temperature, max_tokens, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		settings.UserID,
		settings.Provider,
		settings.APIEndpoint,
		encryptedKey,
		settings.Model,
		settings.Temperature,
		settings.MaxTokens,
		settings.IsActive,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create AI settings: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	settings.ID = id
	settings.APIKeyEncrypted = encryptedKey
	settings.CreatedAt = now
	settings.UpdatedAt = now

	return nil
}

// UpdateAISettings updates existing AI settings with encrypted API key
func (r *AISettingsRepository) UpdateAISettings(ctx context.Context, settings *model.AISettings) error {
	// Encrypt API key if provided
	var encryptedKey string
	if settings.APIKey != "" {
		var err error
		encryptedKey, err = r.crypto.EncryptAES(settings.APIKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt API key: %w", err)
		}
	} else {
		// Keep existing encrypted key
		existing, err := r.GetAISettingsByID(ctx, settings.UserID, settings.ID)
		if err != nil {
			return err
		}
		encryptedKey = existing.APIKeyEncrypted
	}

	// Deactivate other settings if this one is being activated
	if settings.IsActive {
		if err := r.deactivateAllSettings(ctx, settings.UserID); err != nil {
			return fmt.Errorf("failed to deactivate other settings: %w", err)
		}
	}

	query := `
		UPDATE ai_settings
		SET provider = ?, api_endpoint = ?, api_key_encrypted = ?, model = ?, temperature = ?, max_tokens = ?, is_active = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		settings.Provider,
		settings.APIEndpoint,
		encryptedKey,
		settings.Model,
		settings.Temperature,
		settings.MaxTokens,
		settings.IsActive,
		time.Now(),
		settings.ID,
		settings.UserID,
	)

	if err != nil {
		return fmt.Errorf("failed to update AI settings: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrAISettingsNotFound
	}

	return nil
}

// GetAISettings retrieves AI settings by ID with decrypted API key
func (r *AISettingsRepository) GetAISettingsByID(ctx context.Context, userID, settingsID int64) (*model.AISettings, error) {
	query := `
		SELECT id, user_id, provider, api_endpoint, api_key_encrypted, model, temperature, max_tokens, is_active, created_at, updated_at
		FROM ai_settings
		WHERE id = ? AND user_id = ?
	`

	settings := &model.AISettings{}
	err := r.db.QueryRowContext(ctx, query, settingsID, userID).Scan(
		&settings.ID,
		&settings.UserID,
		&settings.Provider,
		&settings.APIEndpoint,
		&settings.APIKeyEncrypted,
		&settings.Model,
		&settings.Temperature,
		&settings.MaxTokens,
		&settings.IsActive,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAISettingsNotFound
		}
		return nil, fmt.Errorf("failed to get AI settings: %w", err)
	}

	// Decrypt API key
	decryptedKey, err := r.crypto.DecryptAES(settings.APIKeyEncrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt API key: %w", err)
	}
	settings.APIKey = decryptedKey

	return settings, nil
}

// GetActiveAISettings retrieves the active AI settings for a user
func (r *AISettingsRepository) GetActiveAISettings(ctx context.Context, userID int64) (*model.AISettings, error) {
	query := `
		SELECT id, user_id, provider, api_endpoint, api_key_encrypted, model, temperature, max_tokens, is_active, created_at, updated_at
		FROM ai_settings
		WHERE user_id = ? AND is_active = true
		LIMIT 1
	`

	settings := &model.AISettings{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&settings.ID,
		&settings.UserID,
		&settings.Provider,
		&settings.APIEndpoint,
		&settings.APIKeyEncrypted,
		&settings.Model,
		&settings.Temperature,
		&settings.MaxTokens,
		&settings.IsActive,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAISettingsNotFound
		}
		return nil, fmt.Errorf("failed to get active AI settings: %w", err)
	}

	// Decrypt API key
	decryptedKey, err := r.crypto.DecryptAES(settings.APIKeyEncrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt API key: %w", err)
	}
	settings.APIKey = decryptedKey

	return settings, nil
}

// ListAISettings retrieves all AI settings for a user
func (r *AISettingsRepository) ListAISettings(ctx context.Context, userID int64) ([]*model.AISettings, error) {
	query := `
		SELECT id, user_id, provider, api_endpoint, api_key_encrypted, model, temperature, max_tokens, is_active, created_at, updated_at
		FROM ai_settings
		WHERE user_id = ?
		ORDER BY is_active DESC, created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list AI settings: %w", err)
	}
	defer rows.Close()

	var settingsList []*model.AISettings
	for rows.Next() {
		settings := &model.AISettings{}
		err := rows.Scan(
			&settings.ID,
			&settings.UserID,
			&settings.Provider,
			&settings.APIEndpoint,
			&settings.APIKeyEncrypted,
			&settings.Model,
			&settings.Temperature,
			&settings.MaxTokens,
			&settings.IsActive,
			&settings.CreatedAt,
			&settings.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan AI settings: %w", err)
		}

		// Decrypt API key
		decryptedKey, err := r.crypto.DecryptAES(settings.APIKeyEncrypted)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt API key: %w", err)
		}
		settings.APIKey = decryptedKey

		settingsList = append(settingsList, settings)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating AI settings: %w", err)
	}

	return settingsList, nil
}

// DeleteAISettings deletes AI settings
func (r *AISettingsRepository) DeleteAISettings(ctx context.Context, userID, settingsID int64) error {
	query := `
		DELETE FROM ai_settings
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query, settingsID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete AI settings: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrAISettingsNotFound
	}

	return nil
}

// deactivateAllSettings deactivates all AI settings for a user
func (r *AISettingsRepository) deactivateAllSettings(ctx context.Context, userID int64) error {
	query := `
		UPDATE ai_settings
		SET is_active = false, updated_at = ?
		WHERE user_id = ? AND is_active = true
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate settings: %w", err)
	}

	return nil
}
