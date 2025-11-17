-- 回滚用户角色和系统设置功能迁移

USE ai_diet_assistant;

-- 删除系统设置表
DROP TABLE IF EXISTS system_settings;

-- 删除 users 表的 role 字段
ALTER TABLE users DROP COLUMN role;
