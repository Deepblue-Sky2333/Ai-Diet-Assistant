-- 回滚用户偏好表结构更新
-- 将表结构恢复到之前的状态

USE ai_diet_assistant;

-- 1. 删除新增的营养目标字段
ALTER TABLE user_preferences 
    DROP COLUMN daily_fiber_goal,
    DROP COLUMN daily_fat_goal,
    DROP COLUMN daily_carbs_goal,
    DROP COLUMN daily_protein_goal;

-- 2. 重命名 daily_calories_goal 回 daily_calorie_target
ALTER TABLE user_preferences 
    CHANGE COLUMN daily_calories_goal daily_calorie_target INT COMMENT '每日卡路里目标';

-- 3. 恢复 dietary_restrictions 为 JSON 类型
ALTER TABLE user_preferences 
    MODIFY COLUMN dietary_restrictions JSON COMMENT '饮食限制';

-- 4. 恢复 taste_preferences 为 JSON 类型
ALTER TABLE user_preferences 
    MODIFY COLUMN taste_preferences JSON COMMENT '口味偏好';
