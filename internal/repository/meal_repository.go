package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/yourusername/ai-diet-assistant/internal/model"
)

// MealRepository handles meal data access operations
type MealRepository struct {
	db *sql.DB
}

// NewMealRepository creates a new MealRepository instance
func NewMealRepository(db *sql.DB) *MealRepository {
	return &MealRepository{db: db}
}

// CreateMeal creates a new meal record
func (r *MealRepository) CreateMeal(meal *model.Meal) error {
	// Marshal foods and nutrition to JSON
	foodsJSON, err := json.Marshal(meal.Foods)
	if err != nil {
		return fmt.Errorf("failed to marshal foods: %w", err)
	}

	nutritionJSON, err := json.Marshal(meal.Nutrition)
	if err != nil {
		return fmt.Errorf("failed to marshal nutrition: %w", err)
	}

	query := `
		INSERT INTO meals (user_id, meal_date, meal_type, foods, nutrition, notes)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		meal.UserID,
		meal.MealDate,
		meal.MealType,
		foodsJSON,
		nutritionJSON,
		meal.Notes,
	)
	if err != nil {
		return fmt.Errorf("failed to create meal: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	meal.ID = id
	return nil
}

// UpdateMeal updates an existing meal record (with ownership verification)
func (r *MealRepository) UpdateMeal(userID, mealID int64, meal *model.Meal) error {
	// Marshal foods and nutrition to JSON
	foodsJSON, err := json.Marshal(meal.Foods)
	if err != nil {
		return fmt.Errorf("failed to marshal foods: %w", err)
	}

	nutritionJSON, err := json.Marshal(meal.Nutrition)
	if err != nil {
		return fmt.Errorf("failed to marshal nutrition: %w", err)
	}

	query := `
		UPDATE meals 
		SET meal_date = ?, meal_type = ?, foods = ?, nutrition = ?, notes = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.Exec(
		query,
		meal.MealDate,
		meal.MealType,
		foodsJSON,
		nutritionJSON,
		meal.Notes,
		mealID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update meal: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("meal not found or access denied")
	}

	return nil
}

// DeleteMeal deletes a meal record (with ownership verification)
func (r *MealRepository) DeleteMeal(userID, mealID int64) error {
	query := `DELETE FROM meals WHERE id = ? AND user_id = ?`

	result, err := r.db.Exec(query, mealID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete meal: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("meal not found or access denied")
	}

	return nil
}

// GetMealByID retrieves a meal record by ID (with ownership verification)
func (r *MealRepository) GetMealByID(userID, mealID int64) (*model.Meal, error) {
	query := `
		SELECT id, user_id, meal_date, meal_type, foods, nutrition, notes, created_at, updated_at
		FROM meals
		WHERE id = ? AND user_id = ?
	`

	meal := &model.Meal{}
	var foodsJSON, nutritionJSON []byte

	err := r.db.QueryRow(query, mealID, userID).Scan(
		&meal.ID,
		&meal.UserID,
		&meal.MealDate,
		&meal.MealType,
		&foodsJSON,
		&nutritionJSON,
		&meal.Notes,
		&meal.CreatedAt,
		&meal.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("meal not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get meal: %w", err)
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(foodsJSON, &meal.Foods); err != nil {
		return nil, fmt.Errorf("failed to unmarshal foods: %w", err)
	}

	if err := json.Unmarshal(nutritionJSON, &meal.Nutrition); err != nil {
		return nil, fmt.Errorf("failed to unmarshal nutrition: %w", err)
	}

	return meal, nil
}

// ListMeals retrieves a list of meals with filtering and pagination
func (r *MealRepository) ListMeals(userID int64, filter *model.MealFilter) ([]*model.Meal, int, error) {
	// Build the WHERE clause
	whereClauses := []string{"user_id = ?"}
	args := []interface{}{userID}

	if filter.StartDate != nil {
		whereClauses = append(whereClauses, "meal_date >= ?")
		args = append(args, *filter.StartDate)
	}

	if filter.EndDate != nil {
		whereClauses = append(whereClauses, "meal_date <= ?")
		args = append(args, *filter.EndDate)
	}

	if filter.MealType != "" {
		whereClauses = append(whereClauses, "meal_type = ?")
		args = append(args, filter.MealType)
	}

	whereClause := strings.Join(whereClauses, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM meals WHERE %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count meals: %w", err)
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT id, user_id, meal_date, meal_type, foods, nutrition, notes, created_at, updated_at
		FROM meals
		WHERE %s
		ORDER BY meal_date DESC, created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list meals: %w", err)
	}
	defer rows.Close()

	meals := make([]*model.Meal, 0)
	for rows.Next() {
		meal := &model.Meal{}
		var foodsJSON, nutritionJSON []byte

		err := rows.Scan(
			&meal.ID,
			&meal.UserID,
			&meal.MealDate,
			&meal.MealType,
			&foodsJSON,
			&nutritionJSON,
			&meal.Notes,
			&meal.CreatedAt,
			&meal.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan meal: %w", err)
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(foodsJSON, &meal.Foods); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal foods: %w", err)
		}

		if err := json.Unmarshal(nutritionJSON, &meal.Nutrition); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal nutrition: %w", err)
		}

		meals = append(meals, meal)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating meals: %w", err)
	}

	return meals, total, nil
}

// GetMonthlyMeals retrieves all meals for a specific month
func (r *MealRepository) GetMonthlyMeals(userID int64, year, month int) ([]*model.Meal, error) {
	// Calculate start and end dates for the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	query := `
		SELECT id, user_id, meal_date, meal_type, foods, nutrition, notes, created_at, updated_at
		FROM meals
		WHERE user_id = ? AND meal_date >= ? AND meal_date <= ?
		ORDER BY meal_date ASC, created_at ASC
	`

	rows, err := r.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly meals: %w", err)
	}
	defer rows.Close()

	meals := make([]*model.Meal, 0)
	for rows.Next() {
		meal := &model.Meal{}
		var foodsJSON, nutritionJSON []byte

		err := rows.Scan(
			&meal.ID,
			&meal.UserID,
			&meal.MealDate,
			&meal.MealType,
			&foodsJSON,
			&nutritionJSON,
			&meal.Notes,
			&meal.CreatedAt,
			&meal.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan meal: %w", err)
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(foodsJSON, &meal.Foods); err != nil {
			return nil, fmt.Errorf("failed to unmarshal foods: %w", err)
		}

		if err := json.Unmarshal(nutritionJSON, &meal.Nutrition); err != nil {
			return nil, fmt.Errorf("failed to unmarshal nutrition: %w", err)
		}

		meals = append(meals, meal)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating meals: %w", err)
	}

	return meals, nil
}
