package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
)

// APILogRepository API 日志仓储接口
type APILogRepository interface {
	CreateAPILog(ctx context.Context, log *model.APILog) error
	GetAPILogs(ctx context.Context, filter *model.APILogFilter) ([]*model.APILog, *model.Pagination, error)
	CleanupOldLogs(ctx context.Context, retentionDays int) (int64, error)
}

// apiLogRepository API 日志仓储实现
type apiLogRepository struct {
	db *sql.DB
}

// NewAPILogRepository 创建 API 日志仓储实例
func NewAPILogRepository(db *sql.DB) APILogRepository {
	return &apiLogRepository{
		db: db,
	}
}

// CreateAPILog 创建 API 日志
func (r *apiLogRepository) CreateAPILog(ctx context.Context, log *model.APILog) error {
	// 使用预编译语句防止 SQL 注入
	query := `
		INSERT INTO api_logs (user_id, method, path, status_code, ip_address, user_agent, response_time_ms, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		log.UserID,
		log.Method,
		log.Path,
		log.StatusCode,
		log.IPAddress,
		log.UserAgent,
		log.ResponseTimeMs,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create api log: %w", err)
	}

	// 获取插入的 ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	log.ID = id
	log.CreatedAt = now

	return nil
}

// GetAPILogs 获取 API 日志列表（支持分页和筛选）
func (r *apiLogRepository) GetAPILogs(ctx context.Context, filter *model.APILogFilter) ([]*model.APILog, *model.Pagination, error) {
	// 设置默认分页参数
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	// 构建查询条件
	var conditions []string
	var args []interface{}

	if filter.UserID != nil {
		conditions = append(conditions, "user_id = ?")
		args = append(args, *filter.UserID)
	}

	if filter.StartDate != nil {
		conditions = append(conditions, "created_at >= ?")
		args = append(args, *filter.StartDate)
	}

	if filter.EndDate != nil {
		conditions = append(conditions, "created_at <= ?")
		args = append(args, *filter.EndDate)
	}

	if filter.Method != "" {
		conditions = append(conditions, "method = ?")
		args = append(args, filter.Method)
	}

	if filter.Path != "" {
		conditions = append(conditions, "path LIKE ?")
		args = append(args, "%"+filter.Path+"%")
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM api_logs %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count api logs: %w", err)
	}

	// 计算分页信息
	totalPages := (total + filter.PageSize - 1) / filter.PageSize
	offset := (filter.Page - 1) * filter.PageSize

	pagination := &model.Pagination{
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	// 如果没有数据，直接返回
	if total == 0 {
		return []*model.APILog{}, pagination, nil
	}

	// 查询数据
	dataQuery := fmt.Sprintf(`
		SELECT id, user_id, method, path, status_code, ip_address, user_agent, response_time_ms, created_at
		FROM api_logs
		%s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, filter.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query api logs: %w", err)
	}
	defer rows.Close()

	// 解析结果
	logs := make([]*model.APILog, 0)
	for rows.Next() {
		log := &model.APILog{}
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Method,
			&log.Path,
			&log.StatusCode,
			&log.IPAddress,
			&log.UserAgent,
			&log.ResponseTimeMs,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan api log: %w", err)
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating api logs: %w", err)
	}

	return logs, pagination, nil
}

// CleanupOldLogs 清理旧日志（保留指定天数）
func (r *apiLogRepository) CleanupOldLogs(ctx context.Context, retentionDays int) (int64, error) {
	// 计算截止日期
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)

	// 使用预编译语句防止 SQL 注入
	query := `
		DELETE FROM api_logs
		WHERE created_at < ?
	`

	result, err := r.db.ExecContext(ctx, query, cutoffDate)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old logs: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}
