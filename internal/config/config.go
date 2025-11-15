package config

import (
	"time"
)

// Config 应用配置结构
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Encryption EncryptionConfig `mapstructure:"encryption"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	CORS      CORSConfig      `mapstructure:"cors"`
	Log       LogConfig       `mapstructure:"log"`
	AI        AIConfig        `mapstructure:"ai"`
	Security  SecurityConfig  `mapstructure:"security"`
	Upload    UploadConfig    `mapstructure:"upload"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	Charset         string        `mapstructure:"charset"`
	ParseTime       bool          `mapstructure:"parse_time"`
	Loc             string        `mapstructure:"loc"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Password   string `mapstructure:"password"`
	DB         int    `mapstructure:"db"`
	MaxRetries int    `mapstructure:"max_retries"`
	PoolSize   int    `mapstructure:"pool_size"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret            string `mapstructure:"secret"`
	ExpireHours       int    `mapstructure:"expire_hours"`
	RefreshExpireHours int   `mapstructure:"refresh_expire_hours"`
}

// EncryptionConfig 加密配置
type EncryptionConfig struct {
	AESKey string `mapstructure:"aes_key"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled           bool   `mapstructure:"enabled"`
	RequestsPerMinute int    `mapstructure:"requests_per_minute"`
	Burst             int    `mapstructure:"burst"`
	StorageType       string `mapstructure:"storage_type"` // "memory" or "redis"
}

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowedOrigins   []string      `mapstructure:"allowed_origins"`
	AllowedMethods   []string      `mapstructure:"allowed_methods"`
	AllowedHeaders   []string      `mapstructure:"allowed_headers"`
	ExposeHeaders    []string      `mapstructure:"expose_headers"`
	AllowCredentials bool          `mapstructure:"allow_credentials"`
	MaxAge           time.Duration `mapstructure:"max_age"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// AIConfig AI 配置
type AIConfig struct {
	Provider    string        `mapstructure:"provider"`
	APIKey      string        `mapstructure:"api_key"`
	Model       string        `mapstructure:"model"`
	MaxTokens   int           `mapstructure:"max_tokens"`
	Temperature float64       `mapstructure:"temperature"`
	Timeout     time.Duration `mapstructure:"timeout"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	MaxLoginAttempts   int           `mapstructure:"max_login_attempts"`
	LockoutDuration    time.Duration `mapstructure:"lockout_duration"`
	PasswordMinLength  int           `mapstructure:"password_min_length"`
	RequireSpecialChar bool          `mapstructure:"require_special_char"`
	RequireNumber      bool          `mapstructure:"require_number"`
	RequireUppercase   bool          `mapstructure:"require_uppercase"`
	DefaultUser        DefaultUserConfig `mapstructure:"default_user"`
}

// DefaultUserConfig 默认用户配置
type DefaultUserConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	MaxFileSize  int64    `mapstructure:"max_file_size"`
	AllowedTypes []string `mapstructure:"allowed_types"`
	UploadPath   string   `mapstructure:"upload_path"`
}
