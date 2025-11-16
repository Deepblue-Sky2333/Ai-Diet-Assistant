# 认证模块

## 概述

认证模块提供用户身份验证和授权功能，包括用户登录、Token 刷新、登出和密码修改。系统使用 JWT (JSON Web Token) 进行无状态认证，采用双 Token 机制（Access Token 和 Refresh Token）以提高安全性和用户体验。

**核心功能**：
- 用户登录并获取 Token
- 使用 Refresh Token 刷新 Access Token
- 用户登出并使 Token 失效
- 修改登录密码

**安全特性**：
- 密码哈希存储
- 登录失败限流（防暴力破解）
- Token 黑名单机制
- 密码版本控制（密码修改后旧 Token 自动失效）

---

## 接口列表

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/auth/login` | 用户登录 | 否 |
| POST | `/api/v1/auth/refresh` | 刷新 Token | 否 |
| POST | `/api/v1/auth/logout` | 用户登出 | 是 |
| PUT | `/api/v1/auth/password` | 修改密码 | 是 |

---

## 接口详情

### 用户登录

**接口**: `POST /api/v1/auth/login`

**说明**: 用户使用用户名和密码登录系统，成功后返回 Access Token 和 Refresh Token。系统会记录登录尝试，连续失败超过限制次数（默认 5 次）后账户将被临时锁定。

**认证**: 否

#### 请求参数

##### 请求体

```json
{
  "username": "testuser",
  "password": "password123"
}
```

| 字段 | 类型 | 必填 | 说明 | 验证规则 |
|------|------|------|------|----------|
| username | string | 是 | 用户名 | 长度 3-50 字符，仅允许字母和数字 |
| password | string | 是 | 密码 | 长度 8-128 字符 |

#### 请求示例

```bash
curl -X POST http://localhost:9090/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwicHdkX3ZlciI6MTY5OTk5OTk5OSwiZXhwIjoxNzAwMDg2Mzk5LCJpYXQiOjE2OTk5OTk5OTksIm5iZiI6MTY5OTk5OTk5OX0.xxx",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwicHdkX3ZlciI6MTY5OTk5OTk5OSwiZXhwIjoxNzAwNjA0Nzk5LCJpYXQiOjE2OTk5OTk5OTksIm5iZiI6MTY5OTk5OTk5OX0.yyy",
    "expires_in": 86400
  },
  "timestamp": 1699999999
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| access_token | string | 访问令牌，用于访问受保护的 API，有效期 24 小时 |
| refresh_token | string | 刷新令牌，用于获取新的访问令牌，有效期 7 天 |
| expires_in | int64 | Access Token 有效期（秒），默认 86400（24 小时） |

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid request parameters: Key: 'LoginRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag",
  "timestamp": 1699999999
}
```

**错误响应 (401)**:

```json
{
  "code": 40101,
  "message": "unauthorized",
  "error": "invalid username or password",
  "timestamp": 1699999999
}
```

**错误响应 (429)**:

```json
{
  "code": 42901,
  "message": "too many requests",
  "error": "account locked due to too many failed login attempts",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | 用户名或密码格式不正确、缺少必填字段 |
| 40101 | 未授权 | 用户名或密码错误 |
| 42901 | 请求过于频繁 | 登录失败次数过多，账户被临时锁定 |
| 50001 | 内部错误 | 服务器内部错误 |

#### 注意事项

1. **登录限流**：连续登录失败 5 次后，账户将被锁定 15 分钟
2. **密码安全**：密码使用 bcrypt 算法加密存储
3. **Token 存储**：客户端应安全存储 Token，建议使用 localStorage 或 sessionStorage
4. **IP 记录**：系统会记录登录 IP 地址用于安全审计

---

### 刷新 Token

**接口**: `POST /api/v1/auth/refresh`

**说明**: 使用 Refresh Token 获取新的 Access Token。当 Access Token 过期时，客户端应使用此接口获取新的 Access Token，而无需用户重新登录。

**认证**: 否（但需要提供有效的 Refresh Token）

#### 请求参数

##### 请求体

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

| 字段 | 类型 | 必填 | 说明 | 验证规则 |
|------|------|------|------|----------|
| refresh_token | string | 是 | 刷新令牌 | 最小长度 20 字符 |

#### 请求示例

```bash
curl -X POST http://localhost:9090/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwicHdkX3ZlciI6MTY5OTk5OTk5OSwiZXhwIjoxNzAwNjA0Nzk5LCJpYXQiOjE2OTk5OTk5OTksIm5iZiI6MTY5OTk5OTk5OX0.yyy"
  }'
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwicHdkX3ZlciI6MTY5OTk5OTk5OSwiZXhwIjoxNzAwMDg2Mzk5LCJpYXQiOjE2OTk5OTk5OTksIm5iZiI6MTY5OTk5OTk5OX0.zzz"
  },
  "timestamp": 1699999999
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| access_token | string | 新的访问令牌，有效期 24 小时 |

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid request parameters: Key: 'RefreshTokenRequest.RefreshToken' Error:Field validation for 'RefreshToken' failed on the 'required' tag",
  "timestamp": 1699999999
}
```

**错误响应 (401)**:

```json
{
  "code": 40101,
  "message": "unauthorized",
  "error": "invalid or expired refresh token",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | Refresh Token 格式不正确、缺少必填字段 |
| 40101 | 未授权 | Refresh Token 无效、已过期或用户密码已修改 |
| 50001 | 内部错误 | 服务器内部错误 |

#### 注意事项

1. **Token 有效期**：Refresh Token 有效期为 7 天，过期后需要重新登录
2. **密码修改**：用户修改密码后，所有旧的 Refresh Token 将失效
3. **自动刷新**：建议在 Access Token 即将过期时自动调用此接口
4. **错误处理**：如果刷新失败，应引导用户重新登录

---

### 用户登出

**接口**: `POST /api/v1/auth/logout`

**说明**: 用户登出系统，将当前 Access Token 加入黑名单。登出后该 Token 将无法再使用，即使未过期。此接口具有幂等性，重复调用不会报错。

**认证**: 是

#### 请求参数

无请求参数，Token 从 Authorization 请求头中获取。

#### 请求示例

```bash
curl -X POST http://localhost:9090/api/v1/auth/logout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "logout successful",
  "data": null,
  "timestamp": 1699999999
}
```

**注意**：即使 Token 无效或已过期，此接口也会返回成功响应（幂等性设计）。

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 50001 | 内部错误 | 服务器内部错误（极少发生，且不影响登出结果） |

#### 注意事项

1. **幂等性**：重复调用此接口不会报错，始终返回成功
2. **Token 清理**：客户端应在调用此接口后清除本地存储的 Token
3. **黑名单机制**：Token 被加入 Redis 黑名单，有效期为 Token 的剩余有效时间
4. **Refresh Token**：登出只会使当前 Access Token 失效，Refresh Token 不受影响（但建议客户端一并清除）

---

### 修改密码

**接口**: `PUT /api/v1/auth/password`

**说明**: 用户修改登录密码。修改成功后，密码版本会更新，所有旧的 Token（包括 Access Token 和 Refresh Token）将立即失效，用户需要重新登录。

**认证**: 是

#### 请求参数

##### 请求体

```json
{
  "old_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

| 字段 | 类型 | 必填 | 说明 | 验证规则 |
|------|------|------|------|----------|
| old_password | string | 是 | 旧密码 | 长度 8-128 字符 |
| new_password | string | 是 | 新密码 | 长度 8-128 字符 |

#### 请求示例

```bash
curl -X PUT http://localhost:9090/api/v1/auth/password \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "oldpassword123",
    "new_password": "newpassword456"
  }'
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "password changed successfully, please login again",
  "data": null,
  "timestamp": 1699999999
}
```

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid request parameters: Key: 'ChangePasswordRequest.NewPassword' Error:Field validation for 'NewPassword' failed on the 'min' tag",
  "timestamp": 1699999999
}
```

**错误响应 (401) - 旧密码错误**:

```json
{
  "code": 40101,
  "message": "unauthorized",
  "error": "old password is incorrect",
  "timestamp": 1699999999
}
```

**错误响应 (401) - 未认证**:

```json
{
  "code": 40101,
  "message": "unauthorized",
  "error": "user not authenticated",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | 密码格式不正确、缺少必填字段 |
| 40101 | 未授权 | 旧密码错误或用户未认证 |
| 50001 | 内部错误 | 服务器内部错误 |

#### 注意事项

1. **Token 失效**：密码修改成功后，所有旧 Token 立即失效
2. **重新登录**：用户需要使用新密码重新登录
3. **密码强度**：建议新密码包含大小写字母、数字和特殊字符
4. **密码版本**：系统使用时间戳作为密码版本，确保密码修改后旧 Token 无法使用

---

## 认证流程

### 完整认证流程

```
1. 用户登录
   ↓
2. 获取 Access Token 和 Refresh Token
   ↓
3. 使用 Access Token 访问 API
   ↓
4. Access Token 过期
   ↓
5. 使用 Refresh Token 获取新的 Access Token
   ↓
6. 继续使用新的 Access Token 访问 API
   ↓
7. Refresh Token 过期或用户登出
   ↓
8. 重新登录
```

### Token 刷新流程

```javascript
// 示例：自动刷新 Token 的请求拦截器
async function apiRequest(url, options) {
  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        ...options.headers,
        'Authorization': `Bearer ${getAccessToken()}`
      }
    });
    
    const data = await response.json();
    
    // 检查 Token 是否过期
    if (data.code === 40101 && data.error.includes('expired')) {
      // 尝试刷新 Token
      const newAccessToken = await refreshAccessToken();
      
      // 使用新 Token 重试请求
      return apiRequest(url, {
        ...options,
        headers: {
          ...options.headers,
          'Authorization': `Bearer ${newAccessToken}`
        }
      });
    }
    
    return data;
  } catch (error) {
    console.error('API request failed:', error);
    throw error;
  }
}

async function refreshAccessToken() {
  const refreshToken = getRefreshToken();
  
  const response = await fetch('http://localhost:9090/api/v1/auth/refresh', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ refresh_token: refreshToken })
  });
  
  const data = await response.json();
  
  if (data.code === 0) {
    // 保存新的 Access Token
    saveAccessToken(data.data.access_token);
    return data.data.access_token;
  } else {
    // 刷新失败，跳转到登录页
    redirectToLogin();
    throw new Error('Token refresh failed');
  }
}
```

### 登出流程

```javascript
// 示例：登出流程
async function logout() {
  const accessToken = getAccessToken();
  
  try {
    // 调用登出接口
    await fetch('http://localhost:9090/api/v1/auth/logout', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${accessToken}`
      }
    });
  } catch (error) {
    console.error('Logout request failed:', error);
    // 即使请求失败，也继续清理本地 Token
  } finally {
    // 清除本地存储的 Token
    clearAccessToken();
    clearRefreshToken();
    
    // 跳转到登录页
    redirectToLogin();
  }
}
```

---

## 安全建议

### Token 存储

1. **使用安全存储**：
   - Web 应用：使用 localStorage 或 sessionStorage
   - 移动应用：使用安全的密钥存储（如 iOS Keychain、Android Keystore）
   - 避免在 Cookie 中存储（除非设置了 HttpOnly 和 Secure 标志）

2. **Token 加密**：
   - 对于敏感应用，可以在客户端对 Token 进行加密存储
   - 使用设备特定的密钥进行加密

### 密码安全

1. **密码强度**：
   - 最小长度 8 字符
   - 建议包含大小写字母、数字和特殊字符
   - 避免使用常见密码和个人信息

2. **密码传输**：
   - 始终使用 HTTPS 传输密码
   - 不要在 URL 中包含密码
   - 不要在日志中记录密码

### 防暴力破解

1. **登录限流**：
   - 系统默认限制：5 次失败后锁定 15 分钟
   - 建议客户端也实现限流机制

2. **验证码**：
   - 对于多次失败的登录尝试，可以要求输入验证码
   - 使用图形验证码或短信验证码

### Token 管理

1. **及时刷新**：
   - 在 Access Token 即将过期时自动刷新
   - 不要等到 Token 完全过期才刷新

2. **及时清理**：
   - 用户登出时清除所有 Token
   - 应用关闭时清除 sessionStorage 中的 Token

3. **异常处理**：
   - 捕获所有认证相关的错误
   - 对于 401 错误，尝试刷新 Token 或引导用户重新登录

---

## 常见问题

### Q: Access Token 和 Refresh Token 有什么区别？

A: 
- **Access Token**：用于访问 API，有效期短（24 小时），频繁使用
- **Refresh Token**：用于获取新的 Access Token，有效期长（7 天），较少使用
- 这种设计平衡了安全性和用户体验

### Q: 为什么修改密码后所有 Token 都会失效？

A: 
- 这是一种安全机制，防止密码泄露后攻击者继续使用旧 Token
- 系统使用密码版本（时间戳）来验证 Token 的有效性
- 密码修改后，密码版本更新，所有旧 Token 的密码版本不匹配，因此失效

### Q: 登录失败多次后账户被锁定怎么办？

A: 
- 等待 15 分钟后自动解锁
- 或联系管理员手动解锁
- 这是防暴力破解的安全机制

### Q: 如何判断 Token 是否过期？

A: 
- 方法 1：解析 JWT Token，检查 `exp` 字段
- 方法 2：调用 API，如果返回 40101 错误且消息包含 "expired"，则 Token 已过期
- 建议使用方法 2，因为更可靠（考虑了服务器时间）

### Q: Refresh Token 过期后怎么办？

A: 
- Refresh Token 过期后无法刷新，用户需要重新登录
- 建议在 Refresh Token 即将过期时提醒用户
- 或者实现"记住我"功能，自动重新登录

### Q: 可以同时在多个设备上登录吗？

A: 
- 可以，系统支持多设备同时登录
- 每个设备有独立的 Token
- 在一个设备上登出不会影响其他设备

### Q: 如何实现"记住我"功能？

A: 
- 将 Refresh Token 存储在持久化存储中（如 localStorage）
- 应用启动时检查 Refresh Token 是否有效
- 如果有效，自动获取新的 Access Token
- 如果无效，引导用户登录

---

## 相关文档

- [通用概念](./common-concepts.md) - 详细了解 JWT Token 认证机制
- [错误码说明](./error-codes.md) - 查看所有错误码的详细说明
- [API 文档总览](./README.md) - 返回 API 文档首页

