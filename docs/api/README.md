# AI Diet Assistant API 文档

欢迎使用 AI Diet Assistant API 文档！本文档提供了完整的 API 接口说明，帮助您快速集成和使用我们的饮食管理系统。

## 项目简介

AI Diet Assistant 是一个智能饮食管理系统，帮助用户记录饮食、分析营养、制定饮食计划，并通过 AI 提供个性化的饮食建议。系统提供 RESTful API，支持食材管理、餐饮记录、营养分析、AI 对话等功能。

**核心功能**：
- 🍎 **食材管理**：创建、查询、更新和删除食材信息，支持批量导入
- 🍽️ **餐饮记录**：记录每日三餐和加餐，自动计算营养摄入
- 📊 **营养分析**：统计每日、每月营养数据，对比目标值
- 🤖 **AI 服务**：智能对话、餐饮建议、饮食计划生成
- 📅 **饮食计划**：创建和管理个性化饮食计划
- ⚙️ **设置管理**：配置 AI 服务、用户偏好等

---

## 基础信息

### Base URL

```
http://localhost:9090/api/v1
```

**说明**：
- 开发环境默认端口：`9090`
- 生产环境请替换为实际域名
- 所有接口路径都基于此 Base URL

### 认证方式

系统使用 **JWT (JSON Web Token)** 进行身份认证，采用双 Token 机制：

- **Access Token**：用于访问 API，有效期 24 小时
- **Refresh Token**：用于刷新 Access Token，有效期 7 天

**使用方式**：

在请求头中携带 Access Token：

```bash
Authorization: Bearer YOUR_ACCESS_TOKEN
```

**示例**：

```bash
curl -X GET http://localhost:9090/api/v1/foods \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

详细说明请参考：[通用概念 - 认证机制](./common-concepts.md#认证机制)

### 请求格式

- **Content-Type**: `application/json`
- **字符编码**: UTF-8
- **请求方法**: GET, POST, PUT, DELETE

**请求示例**：

```bash
curl -X POST http://localhost:9090/api/v1/foods \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "鸡胸肉",
    "category": "meat",
    "price": 15.5
  }'
```

### 响应格式

所有接口返回统一的 JSON 格式：

**成功响应**：

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": 1699999999
}
```

**错误响应**：

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "详细错误信息",
  "timestamp": 1699999999
}
```

**分页响应**：

```json
{
  "code": 0,
  "message": "success",
  "data": [],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  },
  "timestamp": 1699999999
}
```

详细说明请参考：[通用概念 - 响应格式](./common-concepts.md#响应格式)

---

## 快速开始

### 1. 环境准备

**系统要求**：
- Go 1.21+
- PostgreSQL 14+
- Redis 6+

**配置文件**：

复制配置文件模板并修改：

```bash
cp configs/config.yaml.example configs/config.yaml
```

编辑 `configs/config.yaml`，配置数据库连接、JWT 密钥等：

```yaml
server:
  port: 9090
  mode: debug

database:
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  dbname: diet_assistant

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your_jwt_secret_key
  access_token_expire: 86400   # 24 hours
  refresh_token_expire: 604800 # 7 days
```

**启动服务**：

```bash
# 运行数据库迁移
make migrate-up

# 启动服务
make run
```

服务启动后，访问 `http://localhost:9090` 验证是否正常运行。

### 2. 获取 Token

首先需要登录获取 Access Token 和 Refresh Token。

**步骤 1：登录**

```bash
curl -X POST http://localhost:9090/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

**响应示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400
  },
  "timestamp": 1699999999
}
```

**步骤 2：保存 Token**

将返回的 `access_token` 和 `refresh_token` 保存到安全的位置（如 localStorage）。

**步骤 3：使用 Token**

在后续的 API 请求中，在请求头中携带 Access Token：

```bash
Authorization: Bearer YOUR_ACCESS_TOKEN
```

### 3. 第一个 API 调用

使用获取的 Token 调用 API，例如获取食材列表：

```bash
curl -X GET http://localhost:9090/api/v1/foods \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "鸡胸肉",
      "category": "meat",
      "price": 15.5,
      "unit": "100g",
      "protein": 25.0,
      "carbs": 0.5,
      "fat": 3.0,
      "calories": 150.0,
      "created_at": 1699999999,
      "updated_at": 1699999999
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 1,
    "total_pages": 1
  },
  "timestamp": 1699999999
}
```

恭喜！您已经成功调用了第一个 API。

---

## 模块导航

API 按功能模块组织，每个模块提供一组相关的接口：

### 核心模块

| 模块 | 说明 | 文档链接 |
|------|------|---------|
| 🔐 认证模块 | 用户登录、Token 刷新、登出、密码修改 | [01-authentication.md](./01-authentication.md) |
| 🍎 食材管理 | 食材的增删改查、批量导入 | [02-foods.md](./02-foods.md) |
| 🍽️ 餐饮记录 | 餐饮记录的增删改查 | [03-meals.md](./03-meals.md) |
| 📅 饮食计划 | 生成和管理饮食计划 | [04-plans.md](./04-plans.md) |
| 🤖 AI 服务 | AI 对话、餐饮建议、对话历史 | [05-ai-services.md](./05-ai-services.md) |
| 📊 营养分析 | 每日统计、月度趋势、营养对比 | [06-nutrition.md](./06-nutrition.md) |
| 📈 Dashboard | 获取仪表盘数据 | [07-dashboard.md](./07-dashboard.md) |
| ⚙️ 设置管理 | AI 设置、用户偏好、用户资料 | [08-settings.md](./08-settings.md) |

### 参考文档

| 文档 | 说明 | 链接 |
|------|------|------|
| 📖 通用概念 | 认证机制、响应格式、分页、日期格式、数据验证 | [common-concepts.md](./common-concepts.md) |
| 📋 数据模型 | 所有数据模型的定义和字段说明 | [data-models.md](./data-models.md) |
| ⚠️ 错误码说明 | 所有错误码的详细说明和处理建议 | [error-codes.md](./error-codes.md) |

---

## 接口快速索引

### 认证模块 (4 个接口)

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| POST | `/auth/login` | 用户登录 | 否 |
| POST | `/auth/refresh` | 刷新 Token | 否 |
| POST | `/auth/logout` | 用户登出 | 是 |
| PUT | `/auth/password` | 修改密码 | 是 |

### 食材管理 (6 个接口)

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| POST | `/foods` | 创建食材 | 是 |
| GET | `/foods` | 获取食材列表 | 是 |
| GET | `/foods/:id` | 获取单个食材 | 是 |
| PUT | `/foods/:id` | 更新食材 | 是 |
| DELETE | `/foods/:id` | 删除食材 | 是 |
| POST | `/foods/batch` | 批量导入食材 | 是 |

### 餐饮记录 (5 个接口)

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| POST | `/meals` | 创建餐饮记录 | 是 |
| GET | `/meals` | 获取餐饮记录列表 | 是 |
| GET | `/meals/:id` | 获取单个餐饮记录 | 是 |
| PUT | `/meals/:id` | 更新餐饮记录 | 是 |
| DELETE | `/meals/:id` | 删除餐饮记录 | 是 |

### 饮食计划 (6 个接口)

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| POST | `/plans/generate` | 生成 AI 饮食计划 | 是 |
| GET | `/plans` | 获取计划列表 | 是 |
| GET | `/plans/:id` | 获取单个计划 | 是 |
| PUT | `/plans/:id` | 更新计划 | 是 |
| DELETE | `/plans/:id` | 删除计划 | 是 |
| POST | `/plans/:id/complete` | 完成计划 | 是 |

### AI 服务 (3 个接口)

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| POST | `/ai/chat` | AI 对话 | 是 |
| POST | `/ai/suggest` | AI 生成餐饮建议 | 是 |
| GET | `/ai/history` | 获取对话历史 | 是 |

### 营养分析 (3 个接口)

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| GET | `/nutrition/daily/:date` | 获取每日营养统计 | 是 |
| GET | `/nutrition/monthly` | 获取月度营养趋势 | 是 |
| GET | `/nutrition/compare` | 对比实际与目标营养 | 是 |

### Dashboard (1 个接口)

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| GET | `/dashboard` | 获取 Dashboard 数据 | 是 |

### 设置管理 (5 个接口)

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| GET | `/settings` | 获取所有设置 | 是 |
| PUT | `/settings/ai` | 更新 AI 设置 | 是 |
| GET | `/settings/ai/test` | 测试 AI 连接 | 是 |
| GET | `/user/profile` | 获取用户资料 | 是 |
| PUT | `/user/preferences` | 更新用户偏好 | 是 |

**总计**：33 个接口

---

## 常见问题

### Token 相关

#### Q: Token 过期后如何处理？

A: 当 Access Token 过期时，使用 Refresh Token 获取新的 Access Token：

```bash
curl -X POST http://localhost:9090/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

如果 Refresh Token 也过期，需要重新登录。

**自动刷新示例**：

```javascript
async function apiRequest(url, options) {
  let response = await fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      'Authorization': `Bearer ${getAccessToken()}`
    }
  });
  
  let data = await response.json();
  
  // 如果 Token 过期，自动刷新
  if (data.code === 40101 && data.error.includes('expired')) {
    const newToken = await refreshAccessToken();
    
    // 使用新 Token 重试
    response = await fetch(url, {
      ...options,
      headers: {
        ...options.headers,
        'Authorization': `Bearer ${newToken}`
      }
    });
    
    data = await response.json();
  }
  
  return data;
}
```

#### Q: 如何判断 Token 是否即将过期？

A: 解析 JWT Token 的 `exp` 字段（过期时间戳），与当前时间对比：

```javascript
function isTokenExpiringSoon(token, thresholdSeconds = 300) {
  const payload = JSON.parse(atob(token.split('.')[1]));
  const expiresAt = payload.exp * 1000; // 转换为毫秒
  const now = Date.now();
  const timeUntilExpiry = expiresAt - now;
  
  return timeUntilExpiry < thresholdSeconds * 1000;
}

// 如果 Token 在 5 分钟内过期，自动刷新
if (isTokenExpiringSoon(accessToken, 300)) {
  await refreshAccessToken();
}
```

### 分页相关

#### Q: 如何使用分页？

A: 在支持分页的接口中，使用 `page` 和 `page_size` 查询参数：

```bash
# 获取第 1 页，每页 20 条（默认）
curl -X GET "http://localhost:9090/api/v1/foods" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 获取第 2 页，每页 50 条
curl -X GET "http://localhost:9090/api/v1/foods?page=2&page_size=50" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": [...],
  "pagination": {
    "page": 2,
    "page_size": 50,
    "total": 150,
    "total_pages": 3
  },
  "timestamp": 1699999999
}
```

**分页参数说明**：
- `page`: 页码，从 1 开始，默认 1
- `page_size`: 每页数据量，默认 20，最大 100
- `total`: 总数据量
- `total_pages`: 总页数

详细说明请参考：[通用概念 - 分页机制](./common-concepts.md#分页机制)

### 日期格式

#### Q: 日期参数应该使用什么格式？

A: 系统使用以下日期格式：

- **日期参数**：`YYYY-MM-DD`（如 `2024-11-16`）
- **日期时间**：ISO 8601 格式（如 `2024-11-16T12:00:00Z`）
- **时间戳**：Unix 时间戳（秒）

**示例**：

```bash
# 查询指定日期范围的餐饮记录
curl -X GET "http://localhost:9090/api/v1/meals?start_date=2024-11-01&end_date=2024-11-30" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 获取指定日期的营养统计
curl -X GET "http://localhost:9090/api/v1/nutrition/daily/2024-11-16" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

详细说明请参考：[通用概念 - 日期时间格式](./common-concepts.md#日期时间格式)

### 错误处理

#### Q: 如何处理 API 错误？

A: 所有错误响应都包含 `code` 和 `message` 字段，根据错误码进行处理：

```javascript
async function handleAPIResponse(response) {
  const data = await response.json();
  
  if (data.code !== 0) {
    switch (data.code) {
      case 40001: // 参数错误
        alert('请检查输入的信息是否正确');
        break;
        
      case 40101: // 未授权
        if (data.error.includes('expired')) {
          // Token 过期，尝试刷新
          await refreshAccessToken();
        } else {
          // 其他认证错误，跳转登录
          redirectToLogin();
        }
        break;
        
      case 40401: // 资源不存在
        alert('请求的内容不存在');
        break;
        
      case 42901: // 请求过于频繁
        alert('操作过于频繁，请稍后再试');
        break;
        
      case 50001:
      case 50002:
      case 50003: // 服务器错误
        alert('服务器繁忙，请稍后重试');
        break;
        
      default:
        alert(data.message);
    }
    
    throw new Error(data.message);
  }
  
  return data.data;
}
```

**常见错误码**：

| 错误码 | 说明 | 处理建议 |
|--------|------|---------|
| 40001 | 参数错误 | 检查请求参数 |
| 40101 | 未授权 | 刷新 Token 或重新登录 |
| 40401 | 资源不存在 | 检查资源 ID |
| 42901 | 请求过于频繁 | 等待后重试 |
| 50001 | 服务器错误 | 稍后重试 |

详细说明请参考：[错误码说明](./error-codes.md)

### 营养数据计算

#### Q: 营养数据是如何计算的？

A: 营养数据的计算规则：

1. **食材营养数据**：
   - 每个食材定义了单位重量（默认 100g）的营养成分
   - 包括：蛋白质、碳水化合物、脂肪、热量

2. **餐饮记录营养计算**：
   - 根据食材数量和单位重量计算实际营养摄入
   - 公式：`实际营养 = 食材营养 × (数量 / 单位重量)`
   - 例如：200g 鸡胸肉的蛋白质 = 25g × (200 / 100) = 50g

3. **每日营养统计**：
   - 汇总当天所有餐饮记录的营养数据
   - 按餐饮类型（早餐、午餐、晚餐、加餐）分类统计

4. **营养对比**：
   - 将实际摄入与用户设置的目标值对比
   - 计算达成率：`达成率 = (实际值 / 目标值) × 100%`

**示例**：

```json
{
  "date": "2024-11-16",
  "total": {
    "protein": 120.5,
    "carbs": 250.0,
    "fat": 45.0,
    "calories": 1850.0
  },
  "target": {
    "protein": 150.0,
    "carbs": 300.0,
    "fat": 50.0,
    "calories": 2000.0
  },
  "achievement": {
    "protein": 80.3,
    "carbs": 83.3,
    "fat": 90.0,
    "calories": 92.5
  }
}
```

---

## 开发建议

### 1. 错误处理

- 始终检查响应的 `code` 字段
- 实现统一的错误处理机制
- 对于服务器错误（50xxx），实现重试机制
- 记录错误日志便于排查问题

### 2. Token 管理

- 安全存储 Token（使用 localStorage 或安全存储）
- 实现自动刷新机制
- Token 过期时自动刷新或引导用户登录
- 登出时清除所有 Token

### 3. 请求优化

- 使用分页避免一次请求过多数据
- 缓存不常变化的数据（如食材列表）
- 实现请求防抖和节流
- 避免短时间内重复请求

### 4. 安全性

- 始终使用 HTTPS（生产环境）
- 不要在 URL 中包含敏感信息
- 不要在日志中记录密码和 Token
- 实现 CSRF 保护

### 5. 用户体验

- 显示友好的错误提示
- 实现加载状态提示
- 处理网络异常情况
- 提供离线功能（如果适用）

---

## 版本信息

- **API 版本**: v1
- **文档版本**: 1.0.0
- **最后更新**: 2024-11-16

---

## 技术支持

如果您在使用 API 过程中遇到问题，请：

1. 查看相关模块的详细文档
2. 查看 [错误码说明](./error-codes.md) 了解错误原因
3. 查看 [通用概念](./common-concepts.md) 了解基础概念
4. 联系技术支持团队

---

## 相关资源

- [项目 README](../../README.md) - 项目概述和安装指南
- [快速开始指南](../../QUICKSTART.md) - 快速部署和使用
- [安全文档](../SECURITY.md) - 安全最佳实践
- [OpenAPI 规范](../openapi.yaml) - OpenAPI 格式的 API 定义

---

**祝您使用愉快！** 🎉
