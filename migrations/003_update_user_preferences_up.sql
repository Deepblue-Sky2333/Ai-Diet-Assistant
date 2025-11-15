-- 更新用户偏好表结构以匹配前端期望的扁平化数据结构
-- 修改日期: 2024-11-15
-- 目的: 统一前后端数据结构，解决字段类型和命名不匹配问题

USE ai_diet_assistant;

-- 1. 修改 taste_preferences 从 JSON 改为 TEXT 类型
ALTER TABLE user_preferences 
    MODIFY COLUMN taste_preferences TEXT COMMENT '口味偏好（逗号分隔的字符串）';

-- 2. 修改 dietary_restrictions 从 JSON 改为 TEXT 类型
ALTER TABLE user_preferences 
    MODIFY COLUMN dietary_restrictions TEXT COMMENT '饮食限制（逗号分隔的字符串）';

-- 3. 重命名 daily_calorie_target 为 daily_calories_goal
ALTER TABLE user_preferences 
    CHANGE COLUMN daily_calorie_target daily_calories_goal INT COMMENT '每日卡路里目标';

-- 4. 新增营养目标字段
ALTER TABLE user_preferences 
    ADD COLUMN daily_protein_goal INT DEFAULT 0 COMMENT '每日蛋白质目标(g)' AFTER daily_calories_goal,
    ADD COLUMN daily_carbs_goal INT DEFAULT 0 COMMENT '每日碳水化合物目标(g)' AFTER daily_protein_goal,
    ADD COLUMN daily_fat_goal INT DEFAULT 0 COMMENT '每日脂肪目标(g)' AFTER daily_carbs_goal,
    ADD COLUMN daily_fiber_goal INT DEFAULT 0 COMMENT '每日纤维目标(g)' AFTER daily_fat_goal;

-- 5. 删除不再使用的 preferred_meal_times 字段（如果前端不需要）
-- ALTER TABLE user_preferences DROP COLUMN preferred_meal_times;
-- 注释掉以保持向后兼容，如果确认不需要可以取消注释
