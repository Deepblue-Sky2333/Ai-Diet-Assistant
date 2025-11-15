-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS ai_diet_assistant CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE ai_diet_assistant;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 用户偏好设置表
CREATE TABLE IF NOT EXISTS user_preferences (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    taste_preferences JSON COMMENT '口味偏好',
    dietary_restrictions JSON COMMENT '饮食限制',
    daily_calorie_target INT COMMENT '每日卡路里目标',
    preferred_meal_times JSON COMMENT '偏好用餐时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uk_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 食材表
CREATE TABLE IF NOT EXISTS foods (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL COMMENT '食材名称',
    category ENUM('meat', 'vegetable', 'fruit', 'grain', 'other') NOT NULL COMMENT '食材分类',
    price DECIMAL(10,2) COMMENT '价格',
    unit VARCHAR(20) DEFAULT 'g' COMMENT '单位',
    protein DECIMAL(8,2) COMMENT '蛋白质(g)',
    carbs DECIMAL(8,2) COMMENT '碳水化合物(g)',
    fat DECIMAL(8,2) COMMENT '脂肪(g)',
    fiber DECIMAL(8,2) COMMENT '纤维(g)',
    calories DECIMAL(8,2) COMMENT '卡路里',
    available BOOLEAN DEFAULT TRUE COMMENT '是否可用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_category (user_id, category),
    INDEX idx_available (available),
    INDEX idx_user_available (user_id, available)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 餐饮记录表
CREATE TABLE IF NOT EXISTS meals (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    meal_date DATE NOT NULL COMMENT '用餐日期',
    meal_type ENUM('breakfast', 'lunch', 'dinner', 'snack') NOT NULL COMMENT '餐次类型',
    foods JSON NOT NULL COMMENT '食材列表',
    nutrition JSON NOT NULL COMMENT '营养数据',
    notes TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_date (user_id, meal_date),
    INDEX idx_date (meal_date),
    INDEX idx_user_date_type (user_id, meal_date, meal_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 饮食计划表
CREATE TABLE IF NOT EXISTS plans (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    plan_date DATE NOT NULL COMMENT '计划日期',
    meal_type ENUM('breakfast', 'lunch', 'dinner', 'snack') NOT NULL COMMENT '餐次类型',
    foods JSON NOT NULL COMMENT '食材列表',
    nutrition JSON NOT NULL COMMENT '营养数据',
    status ENUM('pending', 'completed', 'skipped') DEFAULT 'pending' COMMENT '状态',
    ai_reasoning TEXT COMMENT 'AI推荐理由',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_date_status (user_id, plan_date, status),
    INDEX idx_user_status (user_id, status),
    UNIQUE KEY uk_user_date_type (user_id, plan_date, meal_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- AI设置表
CREATE TABLE IF NOT EXISTS ai_settings (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    provider VARCHAR(50) NOT NULL COMMENT 'AI提供商',
    api_endpoint VARCHAR(255) COMMENT 'API端点',
    api_key_encrypted TEXT NOT NULL COMMENT '加密的API密钥',
    model VARCHAR(100) COMMENT '模型名称',
    temperature DECIMAL(3,2) DEFAULT 0.7 COMMENT '温度参数',
    max_tokens INT DEFAULT 1000 COMMENT '最大token数',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否激活',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_active (user_id, is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 聊天历史表
CREATE TABLE IF NOT EXISTS chat_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    user_input TEXT NOT NULL COMMENT '用户输入',
    ai_response TEXT NOT NULL COMMENT 'AI响应',
    context JSON COMMENT '上下文信息',
    tokens_used INT COMMENT '使用的token数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_created (user_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- API日志表
CREATE TABLE IF NOT EXISTS api_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT COMMENT '用户ID',
    method VARCHAR(10) COMMENT 'HTTP方法',
    path VARCHAR(255) COMMENT '请求路径',
    status_code INT COMMENT '状态码',
    ip_address VARCHAR(45) COMMENT 'IP地址',
    user_agent TEXT COMMENT '用户代理',
    response_time_ms INT COMMENT '响应时间(毫秒)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_created (user_id, created_at),
    INDEX idx_created (created_at),
    INDEX idx_status (status_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 登录尝试表（用于限流）
CREATE TABLE IF NOT EXISTS login_attempts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    ip_address VARCHAR(45) NOT NULL COMMENT 'IP地址',
    success BOOLEAN NOT NULL COMMENT '是否成功',
    attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '尝试时间',
    INDEX idx_username_time (username, attempted_at),
    INDEX idx_ip_time (ip_address, attempted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
