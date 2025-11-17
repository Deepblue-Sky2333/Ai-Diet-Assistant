-- 数据迁移脚本：将 chat_history 数据迁移到新的 conversation_flows 和 messages 表
-- 此脚本为每个用户创建默认对话流，并将历史消息迁移到新表结构

-- 步骤1: 为每个有聊天历史的用户创建默认对话流
INSERT INTO conversation_flows (user_id, title, is_favorited, message_count, created_at, updated_at)
SELECT 
    user_id,
    '历史对话' AS title,
    FALSE AS is_favorited,
    0 AS message_count,
    MIN(created_at) AS created_at,
    MAX(created_at) AS updated_at
FROM chat_history
GROUP BY user_id;

-- 步骤2: 将 chat_history 中的消息迁移到 messages 表
-- 为每条聊天历史创建两条消息记录：用户消息和AI响应
INSERT INTO messages (conversation_id, role, content, raw_request, raw_response, created_at)
SELECT 
    cf.id AS conversation_id,
    'user' AS role,
    ch.user_input AS content,
    NULL AS raw_request,
    NULL AS raw_response,
    ch.created_at
FROM chat_history ch
INNER JOIN conversation_flows cf ON ch.user_id = cf.user_id AND cf.title = '历史对话'
ORDER BY ch.created_at;

INSERT INTO messages (conversation_id, role, content, raw_request, raw_response, created_at)
SELECT 
    cf.id AS conversation_id,
    'assistant' AS role,
    ch.ai_response AS content,
    NULL AS raw_request,
    NULL AS raw_response,
    DATE_ADD(ch.created_at, INTERVAL 1 SECOND) AS created_at
FROM chat_history ch
INNER JOIN conversation_flows cf ON ch.user_id = cf.user_id AND cf.title = '历史对话'
ORDER BY ch.created_at;

-- 步骤3: 更新对话流的消息计数
UPDATE conversation_flows cf
SET message_count = (
    SELECT COUNT(*)
    FROM messages m
    WHERE m.conversation_id = cf.id
)
WHERE cf.title = '历史对话';

-- 步骤4: 验证迁移结果
-- 检查是否所有用户的聊天历史都已迁移
SELECT 
    'Migration Validation' AS check_type,
    COUNT(DISTINCT ch.user_id) AS users_in_chat_history,
    COUNT(DISTINCT cf.user_id) AS users_in_conversation_flows,
    (SELECT COUNT(*) FROM chat_history) AS total_chat_history_records,
    (SELECT COUNT(*) FROM messages WHERE role = 'user') AS total_user_messages,
    (SELECT COUNT(*) FROM messages WHERE role = 'assistant') AS total_assistant_messages
FROM chat_history ch
LEFT JOIN conversation_flows cf ON ch.user_id = cf.user_id AND cf.title = '历史对话';

