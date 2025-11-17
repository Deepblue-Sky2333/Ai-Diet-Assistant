-- 添加用户角色和系统设置功能
-- 为用户表添加 role 字段，创建系统设置表

USE ai_diet_assistant;

-- 为 users 表添加 role 字段
ALTER TABLE users 
ADD COLUMN role ENUM('admin', 'user') DEFAULT 'user' COMMENT '用户角色' 
AFTER email;

-- 创建系统设置表
CREATE TABLE IF NOT EXISTS system_settings (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    setting_key VARCHAR(100) UNIQUE NOT NULL COMMENT '设置键',
    setting_value TEXT NOT NULL COMMENT '设置值',
    description VARCHAR(255) COMMENT '设置描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_setting_key (setting_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入默认的注册开关设置
INSERT INTO system_settings (setting_key, setting_value, description) 
VALUES ('registration_enabled', 'true', '是否允许新用户注册')
ON DUPLICATE KEY UPDATE setting_key = setting_key;

-- 将第一个用户设置为管理员（如果存在用户）
UPDATE users SET role = 'admin' WHERE id = (SELECT MIN(id) FROM (SELECT id FROM users) AS temp);
