# 通用概念

本文档介绍 AI Diet Assistant API 中的通用概念和约定，包括认证机制、响应格式、分页机制、日期时间格式和数据验证规则。

---

## 认证机制

### JWT Token 认证

系统使用 JWT (JSON Web Token) 进行用户认证。JWT 是一种无状态的认证方式，token 中包含用户身份信息和权限信息。

#### Token 类型

系统使用双 Token 机制：

1. **Access Token (访问令牌)**
   - 用于访问受保护的 API 接口
   - 有效期较短（默认 24 小时）
   - 包含用户 ID、用户名和密码版本信息

2. **Refresh Token (刷新令牌)**
   - 用于获取新的 Access Token
   - 有效期较长（默认 168 小时，即 7 天）
   - 当 Access Token 过期时使用

#### Token 结构

JWT Token 包含以下声明（Claims）：

```json
{
  "user_id": 1,
  "username": "testuser",
  "pwd_ver": 1699999999,
  "exp": 1700086399,
  "iat": 1699999999,
  "nbf": 1699999999
}
```

| 字段 | 说明 |
|------|------|
| `user_id` | 用户 ID |
| `username` | 用户名 |
| `pwd_ver` | 密码版本（密码修改时间戳），用于密码修改后使旧 token 失效 |
| `exp` | 过期时间（Unix 时间戳） |
| `iat` | 签发时间（Unix 时间戳） |
| `nbf` | 生效时间（Unix 时间戳） |

### 获取 Token

通过登录接口获取 Token：

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

### 使用 Token

在请求头中携带 Access Token：

```bash
curl -X GET http://localhost:9090/api/v1/foods \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**格式要求**：
- 使用 `Authorization` 请求头
- 格式为 `Bearer <token>`
- `Bearer` 和 token 之间有一个空格

### Token 刷新流程

当 Access Token 过期时，使用 Refresh Token 获取新的 Access Token：

```bash
curl -X POST http://localhost:9090/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

**响应示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "timestamp": 1699999999
}
```

**刷新流程**：

1. 客户端检测到 Access Token 过期（收到 401 错误且错误消息为 "token has expired"）
2. 使用 Refresh Token 调用刷新接口
3. 获取新的 Access Token
4. 使用新的 Access Token 重试原请求

### Token 失效场景

Token 会在以下情况下失效：

1. **Token 过期**：超过有效期
   - 错误码：`40101`
   - 错误消息：`"token has expired"`

2. **密码已修改**：用户修改密码后，所有旧 token 失效
   - 错误码：`40101`
   - 错误消息：`"password has been changed, please login again"`

3. **Token 格式错误**：token 格式不正确或签名验证失败
   - 错误码：`40101`
   - 错误消息：`"invalid token"`

4. **Token 已登出**：用户主动登出后，token 被加入黑名单
   - 错误码：`40101`
   - 错误消息：`"invalid token"`

### 登出

登出时，当前 token 会被加入黑名单：

```bash
curl -X POST http://localhost:9090/api/v1/auth/logout \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## 用户角色和权限

### 角色类型

系统支持两种用户角色：

1. **管理员（admin）**
   - 拥有系统管理权限
   - 可以访问系统设置接口
   - 可以修改系统配置（如注册开关）
   - 第一个注册的用户自动成为管理员

2. **普通用户（user）**
   - 拥有基本功能权限
   - 可以管理自己的数据（食材、餐饮、计划等）
   - 可以使用 AI 服务
   - 无法访问系统设置接口

### 角色分配

#### 自动分配

- **第一个用户**：系统中第一个注册的用户自动成为管理员
- **后续用户**：后续注册的用户默认为普通用户

#### 角色识别

用户角色在以下场景中返回：

1. **注册响应**：注册成功后返回用户角色
   ```json
   {
     "code": 0,
     "message": "user registered successfully",
     "data": {
       "id": 1,
       "username": "admin",
       "role": "admin",
       "created_at": "2024-01-01T00:00:00Z"
     }
   }
   ```

2. **登录响应**：登录成功后可以通过用户信息接口获取角色
3. **用户信息接口**：查询用户信息时包含角色字段

### 权限控制

#### 管理员专属接口

以下接口仅限管理员访问，普通用户调用会返回 `40302` 错误：

| 接口 | 说明 |
|------|------|
| `GET /api/v1/settings/system` | 获取系统设置 |
| `PUT /api/v1/settings/system` | 更新系统设置 |

**错误响应示例**：

```json
{
  "code": 40302,
  "message": "forbidden",
  "error": "admin access required",
  "timestamp": 1699999999
}
```

#### 公开接口

以下接口无需认证，任何人都可以访问：

| 接口 | 说明 |
|------|------|
| `POST /api/v1/auth/register` | 用户注册 |
| `POST /api/v1/auth/login` | 用户登录 |
| `POST /api/v1/auth/refresh` | 刷新 Token |
| `GET /api/v1/system/info` | 获取公开系统信息 |

#### 认证接口

除管理员专属接口和公开接口外，其他接口需要认证但不限制角色。

### 前端显示

#### 管理员标识

前端应在用户界面中显示管理员身份，推荐使用以下方式：

1. **Emoji 标识**：在用户名旁显示 👑 或 🔑 emoji
   ```javascript
   function displayUserName(user) {
     if (user.role === 'admin') {
       return `👑 ${user.username}`;
     }
     return user.username;
   }
   ```

2. **徽章标识**：使用徽章或标签
   ```html
   <span class="username">admin</span>
   <span class="badge badge-admin">管理员</span>
   ```

3. **颜色区分**：使用不同颜色区分管理员和普通用户
   ```css
   .user-admin {
     color: #ff6b6b;
     font-weight: bold;
   }
   ```

#### 功能显示

根据用户角色显示或隐藏功能：

```javascript
// 检查用户角色
function isAdmin(user) {
  return user.role === 'admin';
}

// 根据角色显示系统设置菜单
if (isAdmin(currentUser)) {
  showSystemSettingsMenu();
} else {
  hideSystemSettingsMenu();
}
```

### 角色验证流程

```
1. 用户登录
   ↓
2. 获取用户信息（包含角色）
   ↓
3. 前端根据角色显示功能
   ↓
4. 用户访问接口
   ↓
5. 后端验证用户角色
   ↓
6. 返回结果或权限错误
```

### 最佳实践

1. **前后端双重验证**：
   - 前端根据角色控制界面显示
   - 后端验证每个请求的权限
   - 不要仅依赖前端验证

2. **角色缓存**：
   - 在客户端缓存用户角色信息
   - 登录后立即获取并缓存
   - 登出时清除缓存

3. **权限提示**：
   - 当用户尝试访问无权限的功能时，显示友好提示
   - 说明需要管理员权限
   - 提供联系管理员的方式

4. **安全考虑**：
   - 不要在客户端存储敏感的权限信息
   - 所有权限验证必须在服务器端进行
   - 定期刷新用户角色信息

### 常见问题

#### Q: 如何成为管理员？

A: 
- 第一个注册的用户自动成为管理员
- 后续用户需要现有管理员通过数据库或 CLI 工具提升权限
- 系统不提供自助提升权限的功能

#### Q: 可以有多个管理员吗？

A: 
- 可以，系统支持多个管理员
- 管理员可以通过 CLI 工具创建新的管理员账户
- 命令：`./bin/create-user -username <用户名> -password <密码> -role admin`

#### Q: 管理员可以降级为普通用户吗？

A: 
- 当前版本不支持通过 API 修改用户角色
- 需要通过数据库直接修改 users 表的 role 字段
- 建议至少保留一个管理员账户

#### Q: 普通用户可以看到其他用户的数据吗？

A: 
- 不可以，用户只能访问自己的数据
- 系统通过 Token 中的 user_id 进行数据隔离
- 即使知道其他用户的 ID，也无法访问其数据

---

## 响应格式

所有 API 接口使用统一的 JSON 响应格式。

### 成功响应

**结构**：

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": 1699999999
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | int | 响应码，0 表示成功 |
| `message` | string | 响应消息 |
| `data` | object/array | 响应数据，可以是对象或数组 |
| `timestamp` | int64 | 响应时间戳（Unix 时间戳，秒） |

**示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "鸡胸肉",
    "category": "meat",
    "price": 15.5,
    "unit": "100g"
  },
  "timestamp": 1699999999
}
```

### 错误响应

**结构**：

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "详细错误信息",
  "timestamp": 1699999999
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | int | 错误码，非 0 表示错误 |
| `message` | string | 错误消息 |
| `error` | string | 详细错误信息（可选，开发环境返回，生产环境返回通用消息） |
| `timestamp` | int64 | 响应时间戳（Unix 时间戳，秒） |

**示例**：

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid request parameters: Key: 'CreateFoodRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag",
  "timestamp": 1699999999
}
```

### 分页响应

**结构**：

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

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | int | 响应码，0 表示成功 |
| `message` | string | 响应消息 |
| `data` | array | 数据列表 |
| `pagination` | object | 分页信息 |
| `timestamp` | int64 | 响应时间戳（Unix 时间戳，秒） |

**分页信息字段**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `page` | int | 当前页码（从 1 开始） |
| `page_size` | int | 每页数据量 |
| `total` | int | 总数据量 |
| `total_pages` | int | 总页数 |

**示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "鸡胸肉",
      "category": "meat"
    },
    {
      "id": 2,
      "name": "西兰花",
      "category": "vegetable"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  },
  "timestamp": 1699999999
}
```

---

## 分页机制

### 分页参数

支持分页的接口接受以下查询参数：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `page` | int | 否 | 1 | 页码，从 1 开始 |
| `page_size` | int | 否 | 20 | 每页数据量，最大 100 |

**示例**：

```bash
# 获取第 1 页，每页 20 条数据（默认）
curl -X GET "http://localhost:9090/api/v1/foods" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 获取第 2 页，每页 50 条数据
curl -X GET "http://localhost:9090/api/v1/foods?page=2&page_size=50" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 分页计算

- 如果 `page` 小于 1，自动设置为 1
- 如果 `page_size` 小于 1，自动设置为 10
- `total_pages` 计算公式：`(total + page_size - 1) / page_size`

### 空结果处理

当请求的页码超出范围时，返回空数组：

```json
{
  "code": 0,
  "message": "success",
  "data": [],
  "pagination": {
    "page": 10,
    "page_size": 20,
    "total": 50,
    "total_pages": 3
  },
  "timestamp": 1699999999
}
```

---

## 日期时间格式

### 时间戳格式

系统中的时间戳统一使用 **Unix 时间戳（秒）**：

- 响应中的 `timestamp` 字段
- 数据模型中的 `created_at`、`updated_at` 字段

**示例**：

```json
{
  "timestamp": 1699999999,
  "created_at": 1699999999,
  "updated_at": 1699999999
}
```

### 日期格式

日期字段使用 **ISO 8601 格式**：

- 日期：`YYYY-MM-DD`
- 日期时间：`YYYY-MM-DDTHH:MM:SSZ` 或 `YYYY-MM-DDTHH:MM:SS+08:00`

**示例**：

```json
{
  "meal_date": "2024-11-16T12:00:00Z",
  "start_date": "2024-11-01",
  "end_date": "2024-11-30"
}
```

### 日期查询参数

查询参数中的日期使用 `YYYY-MM-DD` 格式：

```bash
# 查询指定日期范围的餐饮记录
curl -X GET "http://localhost:9090/api/v1/meals?start_date=2024-11-01&end_date=2024-11-30" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 时区处理

- 服务器使用 UTC 时区
- 客户端应根据本地时区进行转换
- 日期时间字段建议使用 ISO 8601 格式，包含时区信息

---

## 数据验证规则

### 通用验证规则

所有接口的请求参数都会进行验证，不符合规则的请求会返回 `40001` 错误码。

#### 字符串验证

| 规则 | 说明 | 示例 |
|------|------|------|
| `required` | 必填字段 | `"name" binding:"required"` |
| `min` | 最小长度 | `"name" binding:"min=3"` |
| `max` | 最大长度 | `"name" binding:"max=100"` |
| `alphanum` | 仅允许字母和数字 | `"username" binding:"alphanum"` |
| `oneof` | 枚举值 | `"category" binding:"oneof=meat vegetable"` |

#### 数值验证

| 规则 | 说明 | 示例 |
|------|------|------|
| `required` | 必填字段 | `"price" binding:"required"` |
| `gte` | 大于等于 | `"price" binding:"gte=0"` |
| `lte` | 小于等于 | `"price" binding:"lte=100000"` |
| `gt` | 大于 | `"quantity" binding:"gt=0"` |
| `lt` | 小于 | `"quantity" binding:"lt=1000"` |

#### 数组验证

| 规则 | 说明 | 示例 |
|------|------|------|
| `required` | 必填字段 | `"foods" binding:"required"` |
| `min` | 最小元素数量 | `"foods" binding:"min=1"` |
| `max` | 最大元素数量 | `"foods" binding:"max=50"` |
| `dive` | 验证数组元素 | `"foods" binding:"dive"` |

### 常见字段验证规则

#### 用户名

- 必填
- 长度：3-50 字符
- 仅允许字母和数字

```json
{
  "username": "testuser"
}
```

#### 密码

- 必填
- 长度：8-128 字符

```json
{
  "password": "password123"
}
```

#### 食材分类

- 必填
- 枚举值：`meat`（肉类）、`vegetable`（蔬菜）、`fruit`（水果）、`grain`（谷物）、`other`（其他）

```json
{
  "category": "meat"
}
```

#### 餐饮类型

- 必填
- 枚举值：`breakfast`（早餐）、`lunch`（午餐）、`dinner`（晚餐）、`snack`（加餐）

```json
{
  "meal_type": "breakfast"
}
```

#### 价格和营养数据

- 必填
- 大于等于 0
- 小于等于合理上限

```json
{
  "price": 15.5,
  "protein": 25.0,
  "carbs": 0.5,
  "fat": 3.0,
  "calories": 150.0
}
```

### 验证错误响应

当请求参数不符合验证规则时，返回详细的错误信息：

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid request parameters: Key: 'CreateFoodRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag",
  "timestamp": 1699999999
}
```

**错误信息格式**：
- 包含字段名称
- 包含验证规则
- 包含失败原因

### 批量操作验证

批量操作（如批量导入食材）会验证每个元素：

```json
{
  "foods": [
    {
      "name": "鸡胸肉",
      "category": "meat",
      "price": 15.5
    }
  ]
}
```

如果任何一个元素验证失败，整个请求都会失败。

---

## HTTP 状态码映射

系统的业务错误码会映射到相应的 HTTP 状态码：

| 业务错误码范围 | HTTP 状态码 | 说明 |
|---------------|-------------|------|
| 0 | 200 OK | 成功 |
| 40000-40999 | 400 Bad Request | 请求参数错误 |
| 40100-40199 | 401 Unauthorized | 未认证 |
| 40300-40399 | 403 Forbidden | 无权限 |
| 40400-40499 | 404 Not Found | 资源不存在 |
| 40900-40999 | 409 Conflict | 资源冲突 |
| 42900-42999 | 429 Too Many Requests | 请求过多 |
| 50000+ | 500 Internal Server Error | 服务器错误 |

**注意**：客户端应优先使用响应体中的 `code` 字段判断业务结果，而不是 HTTP 状态码。

---

## 安全性

### 敏感信息保护

- 生产环境不返回详细的错误堆栈信息
- 数据库连接字符串中的密码会被自动脱敏
- 文件路径会被清理，只保留文件名
- SQL 语句会被替换为通用消息

### 请求限制

- 批量操作限制元素数量（如批量导入最多 100 条）
- 分页查询限制每页最大数量（最大 100 条）
- 字符串字段限制最大长度

### 认证要求

大部分接口需要认证，未认证的请求会返回 `40101` 错误码。

**无需认证的接口**：
- `POST /api/v1/auth/login` - 登录
- `POST /api/v1/auth/refresh` - 刷新 Token

**需要认证的接口**：
- 所有其他接口

---

## 最佳实践

### Token 管理

1. **安全存储**：将 token 存储在安全的位置（如 localStorage 或 sessionStorage）
2. **自动刷新**：在 Access Token 即将过期时自动刷新
3. **错误处理**：捕获 401 错误，自动刷新 token 或跳转到登录页
4. **及时清理**：登出时清除本地存储的 token

### 错误处理

1. **检查错误码**：优先使用 `code` 字段判断错误类型
2. **显示友好消息**：根据错误码显示用户友好的错误消息
3. **记录详细错误**：在开发环境记录 `error` 字段的详细信息
4. **重试机制**：对于网络错误或服务器错误，实现合理的重试机制

### 分页查询

1. **合理设置页大小**：根据实际需求设置 `page_size`，避免一次请求过多数据
2. **缓存结果**：对于不常变化的数据，可以缓存分页结果
3. **显示总数**：使用 `total` 字段显示总数据量
4. **处理空结果**：当 `data` 为空数组时，显示"暂无数据"提示

### 日期时间处理

1. **使用标准库**：使用标准的日期时间库（如 JavaScript 的 Date、Moment.js 等）
2. **时区转换**：根据用户时区进行转换
3. **格式化显示**：根据用户习惯格式化日期时间显示
4. **验证格式**：提交前验证日期格式是否正确
