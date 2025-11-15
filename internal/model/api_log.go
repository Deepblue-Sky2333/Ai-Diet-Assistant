package model

import "time"

// APILog API 日志模型
type APILog struct {
	ID             int64     `json:"id" db:"id"`
	UserID         *int64    `json:"user_id,omitempty" db:"user_id"`
	Method         string    `json:"method" db:"method"`
	Path           string    `json:"path" db:"path"`
	StatusCode     int       `json:"status_code" db:"status_code"`
	IPAddress      string    `json:"ip_address" db:"ip_address"`
	UserAgent      string    `json:"user_agent" db:"user_agent"`
	ResponseTimeMs int       `json:"response_time_ms" db:"response_time_ms"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// APILogFilter API 日志筛选条件
type APILogFilter struct {
	UserID    *int64
	StartDate *time.Time
	EndDate   *time.Time
	Method    string
	Path      string
	Page      int
	PageSize  int
}

// Pagination 分页信息
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
