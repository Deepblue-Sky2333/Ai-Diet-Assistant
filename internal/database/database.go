package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yourusername/ai-diet-assistant/internal/config"
)

var db *sql.DB

// Init 初始化数据库连接
func Init(cfg *config.DatabaseConfig) error {
	var err error
	
	// 创建数据库连接
	db, err = sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		// 不直接包装错误，避免暴露DSN
		return fmt.Errorf("failed to open database connection")
	}

	// 配置连接池
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 测试连接
	if err := db.Ping(); err != nil {
		// 不直接包装错误，避免暴露DSN
		return fmt.Errorf("failed to connect to database: connection test failed")
	}

	return nil
}

// GetDB 获取数据库连接
func GetDB() *sql.DB {
	if db == nil {
		panic("database not initialized, call Init() first")
	}
	return db
}

// Close 关闭数据库连接
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// HealthCheck 数据库健康检查
func HealthCheck() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	// 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping 数据库
	if err := db.PingContext(ctx); err != nil {
		// 不直接包装错误，避免暴露DSN
		return fmt.Errorf("database health check failed")
	}

	return nil
}

// GetStats 获取数据库连接池统计信息
func GetStats() sql.DBStats {
	if db == nil {
		return sql.DBStats{}
	}
	return db.Stats()
}
