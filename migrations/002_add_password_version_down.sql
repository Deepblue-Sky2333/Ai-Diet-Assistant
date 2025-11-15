-- 回滚密码版本字段迁移

USE ai_diet_assistant;

-- 删除索引
DROP INDEX idx_password_version ON users;

-- 删除 password_version 字段
ALTER TABLE users 
DROP COLUMN password_version;
