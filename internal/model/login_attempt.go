package model

import "time"

// LoginAttempt 登录尝试记录
type LoginAttempt struct {
	ID          int64     `json:"id" db:"id"`
	Username    string    `json:"username" db:"username"`
	IPAddress   string    `json:"ip_address" db:"ip_address"`
	Success     bool      `json:"success" db:"success"`
	AttemptedAt time.Time `json:"attempted_at" db:"attempted_at"`
}
