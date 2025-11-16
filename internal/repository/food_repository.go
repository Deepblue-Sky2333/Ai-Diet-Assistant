package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
)

// FoodRepository handles food data access operations
type FoodRepository struct {
	db *sql.DB
}

// NewFoodRepository creates a new FoodRepository instance
func NewFoodRepository(db *sql.DB) *FoodRepository {
	return &FoodRepository{db: db}
}

// CreateFood creates a new food item for a user
func (r *FoodRepository) CreateFood(food *model.Food) error {
	query := `
		INSERT INTO foods (user_id, name, category, price, unit, protein, carbs, fat, fiber, calories, available)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		food.UserID,
		food.Name,
		food.Category,
		food.Price,
		food.Unit,
		food.Protein,
		food.Carbs,
		food.Fat,
		food.Fiber,
		food.Calories,
		food.Available,
	)
	if err != nil {
		return fmt.Errorf("failed to create food: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	food.ID = id
	return nil
}

// UpdateFood updates an existing food item (with ownership verification)
func (r *FoodRepository) UpdateFood(userID, foodID int64, food *model.Food) error {
	query := `
		UPDATE foods 
		SET name = ?, category = ?, price = ?, unit = ?, protein = ?, carbs = ?, 
		    fat = ?, fiber = ?, calories = ?, available = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.Exec(
		query,
		food.Name,
		food.Category,
		food.Price,
		food.Unit,
		food.Protein,
		food.Carbs,
		food.Fat,
		food.Fiber,
		food.Calories,
		food.Available,
		foodID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update food: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("food not found or access denied")
	}

	return nil
}

// DeleteFood deletes a food item (with ownership verification)
func (r *FoodRepository) DeleteFood(userID, foodID int64) error {
	query := `DELETE FROM foods WHERE id = ? AND user_id = ?`

	result, err := r.db.Exec(query, foodID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete food: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("food not found or access denied")
	}

	return nil
}

// GetFoodByID retrieves a food item by ID (with ownership verification)
func (r *FoodRepository) GetFoodByID(userID, foodID int64) (*model.Food, error) {
	query := `
		SELECT id, user_id, name, category, price, unit, protein, carbs, fat, fiber, 
		       calories, available, created_at, updated_at
		FROM foods
		WHERE id = ? AND user_id = ?
	`

	food := &model.Food{}
	err := r.db.QueryRow(query, foodID, userID).Scan(
		&food.ID,
		&food.UserID,
		&food.Name,
		&food.Category,
		&food.Price,
		&food.Unit,
		&food.Protein,
		&food.Carbs,
		&food.Fat,
		&food.Fiber,
		&food.Calories,
		&food.Available,
		&food.CreatedAt,
		&food.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("food not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get food: %w", err)
	}

	return food, nil
}

// ListFoods retrieves a list of foods with filtering and pagination
func (r *FoodRepository) ListFoods(userID int64, filter *model.FoodFilter) ([]*model.Food, int, error) {
	// Build the WHERE clause
	whereClauses := []string{"user_id = ?"}
	args := []interface{}{userID}

	if filter.Category != "" {
		whereClauses = append(whereClauses, "category = ?")
		args = append(args, filter.Category)
	}

	if filter.Available != nil {
		whereClauses = append(whereClauses, "available = ?")
		args = append(args, *filter.Available)
	}

	whereClause := strings.Join(whereClauses, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM foods WHERE %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count foods: %w", err)
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT id, user_id, name, category, price, unit, protein, carbs, fat, fiber, 
		       calories, available, created_at, updated_at
		FROM foods
		WHERE %s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list foods: %w", err)
	}
	defer rows.Close()

	foods := make([]*model.Food, 0)
	for rows.Next() {
		food := &model.Food{}
		err := rows.Scan(
			&food.ID,
			&food.UserID,
			&food.Name,
			&food.Category,
			&food.Price,
			&food.Unit,
			&food.Protein,
			&food.Carbs,
			&food.Fat,
			&food.Fiber,
			&food.Calories,
			&food.Available,
			&food.CreatedAt,
			&food.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan food: %w", err)
		}
		foods = append(foods, food)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating foods: %w", err)
	}

	return foods, total, nil
}

// BatchInsertFoods inserts multiple food items in a batch
func (r *FoodRepository) BatchInsertFoods(userID int64, foods []*model.Food) error {
	if len(foods) == 0 {
		return nil
	}

	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO foods (user_id, name, category, price, unit, protein, carbs, fat, fiber, calories, available)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, food := range foods {
		// Force user_id to the authenticated user
		food.UserID = userID

		_, err := stmt.Exec(
			food.UserID,
			food.Name,
			food.Category,
			food.Price,
			food.Unit,
			food.Protein,
			food.Carbs,
			food.Fat,
			food.Fiber,
			food.Calories,
			food.Available,
		)
		if err != nil {
			return fmt.Errorf("failed to insert food '%s': %w", food.Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
