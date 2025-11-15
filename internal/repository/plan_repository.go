package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yourusername/ai-diet-assistant/internal/model"
)

// PlanRepository handles plan data access operations
type PlanRepository struct {
	db *sql.DB
}

// NewPlanRepository creates a new PlanRepository instance
func NewPlanRepository(db *sql.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

// CreatePlan creates a new plan record
func (r *PlanRepository) CreatePlan(plan *model.Plan) error {
	// Marshal foods and nutrition to JSON
	foodsJSON, err := json.Marshal(plan.Foods)
	if err != nil {
		return fmt.Errorf("failed to marshal foods: %w", err)
	}

	nutritionJSON, err := json.Marshal(plan.Nutrition)
	if err != nil {
		return fmt.Errorf("failed to marshal nutrition: %w", err)
	}

	query := `
		INSERT INTO plans (user_id, plan_date, meal_type, foods, nutrition, status, ai_reasoning)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		plan.UserID,
		plan.PlanDate,
		plan.MealType,
		foodsJSON,
		nutritionJSON,
		plan.Status,
		plan.AIReasoning,
	)
	if err != nil {
		return fmt.Errorf("failed to create plan: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	plan.ID = id
	return nil
}

// UpdatePlan updates an existing plan record (with ownership verification)
func (r *PlanRepository) UpdatePlan(userID, planID int64, plan *model.Plan) error {
	// Marshal foods and nutrition to JSON
	foodsJSON, err := json.Marshal(plan.Foods)
	if err != nil {
		return fmt.Errorf("failed to marshal foods: %w", err)
	}

	nutritionJSON, err := json.Marshal(plan.Nutrition)
	if err != nil {
		return fmt.Errorf("failed to marshal nutrition: %w", err)
	}

	query := `
		UPDATE plans 
		SET plan_date = ?, meal_type = ?, foods = ?, nutrition = ?, status = ?, ai_reasoning = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.Exec(
		query,
		plan.PlanDate,
		plan.MealType,
		foodsJSON,
		nutritionJSON,
		plan.Status,
		plan.AIReasoning,
		planID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("plan not found or access denied")
	}

	return nil
}

// DeletePlan deletes a plan record (with ownership verification)
func (r *PlanRepository) DeletePlan(userID, planID int64) error {
	query := `DELETE FROM plans WHERE id = ? AND user_id = ?`

	result, err := r.db.Exec(query, planID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("plan not found or access denied")
	}

	return nil
}

// GetPlanByID retrieves a plan record by ID (with ownership verification)
func (r *PlanRepository) GetPlanByID(userID, planID int64) (*model.Plan, error) {
	query := `
		SELECT id, user_id, plan_date, meal_type, foods, nutrition, status, ai_reasoning, created_at, updated_at
		FROM plans
		WHERE id = ? AND user_id = ?
	`

	plan := &model.Plan{}
	var foodsJSON, nutritionJSON []byte

	err := r.db.QueryRow(query, planID, userID).Scan(
		&plan.ID,
		&plan.UserID,
		&plan.PlanDate,
		&plan.MealType,
		&foodsJSON,
		&nutritionJSON,
		&plan.Status,
		&plan.AIReasoning,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("plan not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(foodsJSON, &plan.Foods); err != nil {
		return nil, fmt.Errorf("failed to unmarshal foods: %w", err)
	}

	if err := json.Unmarshal(nutritionJSON, &plan.Nutrition); err != nil {
		return nil, fmt.Errorf("failed to unmarshal nutrition: %w", err)
	}

	return plan, nil
}

// ListPlans retrieves a list of plans with filtering and pagination
func (r *PlanRepository) ListPlans(userID int64, filter *model.PlanFilter) ([]*model.Plan, int, error) {
	// Build the WHERE clause
	whereClauses := []string{"user_id = ?"}
	args := []interface{}{userID}

	if filter.StartDate != nil {
		whereClauses = append(whereClauses, "plan_date >= ?")
		args = append(args, *filter.StartDate)
	}

	if filter.EndDate != nil {
		whereClauses = append(whereClauses, "plan_date <= ?")
		args = append(args, *filter.EndDate)
	}

	if filter.Status != "" {
		whereClauses = append(whereClauses, "status = ?")
		args = append(args, filter.Status)
	}

	whereClause := strings.Join(whereClauses, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM plans WHERE %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count plans: %w", err)
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT id, user_id, plan_date, meal_type, foods, nutrition, status, ai_reasoning, created_at, updated_at
		FROM plans
		WHERE %s
		ORDER BY plan_date ASC, created_at ASC
		LIMIT ? OFFSET ?
	`, whereClause)

	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list plans: %w", err)
	}
	defer rows.Close()

	plans := make([]*model.Plan, 0)
	for rows.Next() {
		plan := &model.Plan{}
		var foodsJSON, nutritionJSON []byte

		err := rows.Scan(
			&plan.ID,
			&plan.UserID,
			&plan.PlanDate,
			&plan.MealType,
			&foodsJSON,
			&nutritionJSON,
			&plan.Status,
			&plan.AIReasoning,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan plan: %w", err)
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(foodsJSON, &plan.Foods); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal foods: %w", err)
		}

		if err := json.Unmarshal(nutritionJSON, &plan.Nutrition); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal nutrition: %w", err)
		}

		plans = append(plans, plan)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating plans: %w", err)
	}

	return plans, total, nil
}

// UpdatePlanStatus updates the status of a plan (with ownership verification)
func (r *PlanRepository) UpdatePlanStatus(userID, planID int64, status string) error {
	query := `
		UPDATE plans 
		SET status = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.Exec(query, status, planID, userID)
	if err != nil {
		return fmt.Errorf("failed to update plan status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("plan not found or access denied")
	}

	return nil
}
