-- 回滚脚本：删除迁移创建的对话流和消息
-- 注意：此脚本不会恢复 chat_history 表的数据，因为原始数据应该保持不变

-- 步骤1: 删除所有标题为 '历史对话' 的对话流的消息
-- (由于外键级联删除，这一步实际上会自动执行，但为了明确性还是写出来)
DELETE m FROM messages m
INNER JOIN conversation_flows cf ON m.conversation_id = cf.id
WHERE cf.title = '历史对话';

-- 步骤2: 删除所有标题为 '历史对话' 的对话流
DELETE FROM conversation_flows
WHERE title = '历史对话';

-- 验证回滚结果
SELECT 
    'Rollback Validation' AS check_type,
    (SELECT COUNT(*) FROM conversation_flows WHERE title = '历史对话') AS remaining_migrated_conversations,
    (SELECT COUNT(*) FROM messages m 
     INNER JOIN conversation_flows cf ON m.conversation_id = cf.id 
     WHERE cf.title = '历史对话') AS remaining_migrated_messages,
    (SELECT COUNT(*) FROM chat_history) AS chat_history_records_preserved;

