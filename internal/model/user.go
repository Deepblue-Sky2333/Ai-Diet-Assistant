package model

import "time"

// User 用户模型
type User struct {
	ID              int64     `json:"id" db:"id"`
	Username        string    `json:"username" db:"username" binding:"required,min=3,max=50"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	PasswordVersion int64     `json:"-" db:"password_version"` // 密码版本（最后修改时间戳）
	Email           string    `json:"email,omitempty" db:"email" binding:"omitempty,email,max=100"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email,omitempty" binding:"omitempty,email,max=100"`
}

// UpdatePasswordRequest 更新密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
