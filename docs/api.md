# API 文档 / API Documentation

本文档提供 AI Diet Assistant 的 API 接口说明。

## 📋 基础信息

- **Base URL**: `http://localhost:9090/api/v1`
- **认证方式**: JWT Bearer Token
- **请求格式**: JSON
- **响应格式**: JSON

## 🔐 认证流程

### 1. 登录获取 Token

```bash
curl -X POST http://localhost:9090/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"password"}'
```

**响应：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "expires_in": 900
  },
  "timestamp": 1234567890
}
```

### 2. 使用 Token 访问 API

```bash
curl -X GET http://localhost:9090/api/v1/foods \
  -H "Authorization: Bearer eyJhbGc..."
```

### 3. 刷新 Token

```bash
curl -X POST http://localhost:9090/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"eyJhbGc..."}'
```

## 📚 API 端点

### 认证 (Authentication)

| 方法 | 端点 | 说明 |
|------|------|------|
| POST | `/auth/login` | 用户登录 |
| POST | `/auth/refresh` | 刷新访问令牌 |
| POST | `/auth/logout` | 用户登出 |
| PUT | `/auth/password` | 修改密码 |

### 食材管理 (Foods)

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/foods` | 获取食材列表（支持分页和分类筛选） |
| POST | `/foods` | 创建新食材 |
| GET | `/foods/:id` | 获取食材详情 |
| PUT | `/foods/:id` | 更新食材信息 |
| DELETE | `/foods/:id` | 删除食材 |
| POST | `/foods/batch` | 批量导入食材 |

### 餐饮记录 (Meals)

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/meals` | 获取餐饮记录列表（支持日期筛选） |
| POST | `/meals` | 创建餐饮记录 |
| GET | `/meals/:id` | 获取餐饮记录详情 |
| PUT | `/meals/:id` | 更新餐饮记录 |
| DELETE | `/meals/:id` | 删除餐饮记录 |

### 饮食计划 (Plans)

| 方法 | 端点 | 说明 |
|------|------|------|
| POST | `/plans/generate` | 生成 AI 饮食计划 |
| GET | `/plans` | 获取饮食计划列表 |
| GET | `/plans/:id` | 获取饮食计划详情 |
| PUT | `/plans/:id` | 更新饮食计划 |
| DELETE | `/plans/:id` | 删除饮食计划 |
| POST | `/plans/:id/complete` | 完成计划并转为餐饮记录 |

### AI 服务 (AI)

| 方法 | 端点 | 说明 |
|------|------|------|
| POST | `/ai/chat` | 与 AI 对话 |
| GET | `/ai/history` | 获取对话历史 |

### 营养分析 (Nutrition)

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/nutrition/daily/:date` | 获取指定日期的营养统计 |
| GET | `/nutrition/monthly` | 获取月度营养趋势 |
| GET | `/nutrition/compare` | 对比实际与目标营养摄入（支持日期范围） |

### Dashboard

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/dashboard` | 获取综合面板数据 |

### 设置 (Settings)

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/settings` | 获取所有设置 |
| PUT | `/settings/ai` | 更新 AI 配置 |
| GET | `/settings/ai/test` | 测试 AI 连接 |
| GET | `/user/profile` | 获取用户资料 |
| PUT | `/user/preferences` | 更新用户偏好 |

## 📝 请求示例

### 创建食材

```bash
curl -X POST http://localhost:9090/api/v1/foods \
  -H "Authorization: Bearer eyJhbGc..." \
  -H "Content-Type: application/json" \
  -d '{
    "name": "鸡胸肉",
    "category": "meat",
    "price": 15.99,
    "calories": 165,
    "protein": 31,
    "carbs": 0,
    "fat": 3.6,
    "fiber": 0
  }'
```

### 生成饮食计划

```bash
curl -X POST http://localhost:9090/api/v1/plans/generate \
  -H "Authorization: Bearer eyJhbGc..." \
  -H "Content-Type: application/json" \
  -d '{
    "days": 2,
    "preferences": "低碳水，高蛋白"
  }'
```

### AI 对话

```bash
curl -X POST http://localhost:9090/api/v1/ai/chat \
  -H "Authorization: Bearer eyJhbGc..." \
  -H "Content-Type: application/json" \
  -d '{
    "message": "我想减肥，应该怎么吃？"
  }'
```

**响应：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "AI 的回复内容",
    "response": "AI 的回复内容",
    "message_id": 123,
    "tokens_used": 150
  },
  "timestamp": 1234567890
}
```

## 📊 响应格式

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    // 响应数据
  },
  "timestamp": 1234567890
}
```

### 分页响应

```json
{
  "code": 0,
  "message": "success",
  "data": [
    // 数据列表
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  },
  "timestamp": 1234567890
}
```

### 错误响应

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "详细错误信息",
  "timestamp": 1234567890
}
```

## 🔢 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 40001 | 参数无效 |
| 40101 | 未授权 |
| 40301 | 禁止访问 |
| 40401 | 资源未找到 |
| 40901 | 资源冲突 |
| 42901 | 请求过多 |
| 50001 | 内部服务器错误 |
| 50002 | 数据库错误 |
| 50003 | AI 服务错误 |

## 📚 相关文档

- [快速开始](../QUICKSTART.md)
- [安全最佳实践](SECURITY.md)
- [OpenAPI 规范](openapi.yaml)
