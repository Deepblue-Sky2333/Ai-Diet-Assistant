-- 回滚数据库迁移
USE ai_diet_assistant;

-- 删除表（按照依赖关系的逆序）
DROP TABLE IF EXISTS login_attempts;
DROP TABLE IF EXISTS api_logs;
DROP TABLE IF EXISTS chat_history;
DROP TABLE IF EXISTS ai_settings;
DROP TABLE IF EXISTS plans;
DROP TABLE IF EXISTS meals;
DROP TABLE IF EXISTS foods;
DROP TABLE IF EXISTS user_preferences;
DROP TABLE IF EXISTS users;

-- 可选：删除数据库
-- DROP DATABASE IF EXISTS ai_diet_assistant;
