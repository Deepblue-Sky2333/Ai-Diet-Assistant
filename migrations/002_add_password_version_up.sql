-- 添加密码版本字段到用户表
-- 用于在密码修改后使旧的JWT令牌失效

USE ai_diet_assistant;

-- 添加 password_version 字段
ALTER TABLE users 
ADD COLUMN password_version BIGINT DEFAULT 0 COMMENT '密码版本（时间戳），用于令牌失效控制';

-- 初始化现有用户的密码版本为当前时间戳（毫秒）
UPDATE users 
SET password_version = UNIX_TIMESTAMP(NOW()) * 1000 
WHERE password_version = 0;

-- 添加索引以优化查询性能
CREATE INDEX idx_password_version ON users(password_version);
