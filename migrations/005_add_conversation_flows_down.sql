-- 删除消息表（先删除，因为有外键依赖）
DROP TABLE IF EXISTS messages;

-- 删除对话流表
DROP TABLE IF EXISTS conversation_flows;
