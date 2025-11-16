# 错误码说明

## 概述

本文档列出了 AI Diet Assistant API 中所有可能返回的错误码及其说明。所有错误响应都遵循统一的响应格式。

## 错误响应格式

### 标准错误响应结构

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "详细错误信息（开发环境）",
  "timestamp": 1699999999
}
```

### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 业务错误码，详见下方错误码列表 |
| message | string | 错误消息的简短描述 |
| error | string | 详细错误信息（可选，仅在开发环境返回） |
| timestamp | int64 | 响应时间戳（Unix 时间戳） |

### 环境差异

- **开发环境**: `error` 字段包含详细的错误信息，便于调试
- **生产环境**: `error` 字段返回通用的错误描述，不暴露内部实现细节

---

## 错误码分类

错误码采用 5 位数字格式，按范围分类：

- **0**: 成功
- **40xxx**: 客户端错误（4xx 系列）
- **50xxx**: 服务器错误（5xx 系列）

---

## 成功码

### 0 - Success

**说明**: 请求成功处理

**HTTP 状态码**: 200 OK

**触发场景**: 
- 所有成功的 API 请求

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "示例数据"
  },
  "timestamp": 1699999999
}
```

---

## 客户端错误码 (40xxx)

### 40001 - Invalid Parameters

**说明**: 请求参数无效或格式错误

**HTTP 状态码**: 400 Bad Request

**触发场景**:
- 请求参数缺失必填字段
- 参数类型不匹配
- 参数格式不正确（如日期格式错误）
- JSON 解析失败
- 参数值超出允许范围

**常见示例**:
- 创建食材时缺少 `name` 字段
- 日期参数格式不是 `YYYY-MM-DD`
- ID 参数不是有效的数字
- 请求体 JSON 格式错误

**响应示例**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid request parameters: name is required",
  "timestamp": 1699999999
}
```

**处理建议**:
- 检查请求参数是否完整
- 验证参数类型和格式
- 参考接口文档确认参数要求

---

### 40002 - Validation Error

**说明**: 数据验证失败

**HTTP 状态码**: 400 Bad Request

**触发场景**:
- 数据不符合业务规则
- 字段值不在允许的范围内
- 数据格式验证失败

**常见示例**:
- 营养数据为负数
- 邮箱格式不正确
- 密码强度不符合要求

**响应示例**:

```json
{
  "code": 40002,
  "message": "validation error",
  "error": "calories must be greater than 0",
  "timestamp": 1699999999
}
```

**处理建议**:
- 检查数据是否符合业务规则
- 验证字段值的有效性
- 参考数据模型文档了解验证规则

---

### 40101 - Unauthorized

**说明**: 未授权，需要身份认证

**HTTP 状态码**: 401 Unauthorized

**触发场景**:
- 未提供 JWT Token
- Token 格式错误
- Token 已过期
- Token 无效或被篡改
- Token 已被加入黑名单（已登出）
- 用户未认证

**常见示例**:
- 请求头缺少 `Authorization` 字段
- Token 过期需要刷新
- 使用已登出的 Token

**响应示例**:

```json
{
  "code": 40101,
  "message": "unauthorized",
  "error": "user not authenticated",
  "timestamp": 1699999999
}
```

**处理建议**:
- 检查是否已登录
- 验证 Token 是否有效
- 如果 Token 过期，使用刷新 Token 接口获取新 Token
- 如果刷新失败，重新登录

---

### 40301 - Forbidden

**说明**: 禁止访问，权限不足

**HTTP 状态码**: 403 Forbidden

**触发场景**:
- 用户已认证但无权访问该资源
- 尝试访问其他用户的私有数据
- 操作权限不足

**常见示例**:
- 尝试修改其他用户的食材
- 尝试删除其他用户的餐饮记录
- 访问无权限的管理功能

**响应示例**:

```json
{
  "code": 40301,
  "message": "forbidden",
  "error": "you do not have permission to access this resource",
  "timestamp": 1699999999
}
```

**处理建议**:
- 确认当前用户是否有权限执行该操作
- 检查是否访问了正确的资源 ID
- 联系管理员获取必要权限

---

### 40401 - Not Found

**说明**: 资源不存在

**HTTP 状态码**: 404 Not Found

**触发场景**:
- 请求的资源 ID 不存在
- 资源已被删除
- URL 路径错误

**常见示例**:
- 查询不存在的食材 ID
- 访问已删除的餐饮记录
- 获取不存在的饮食计划

**响应示例**:

```json
{
  "code": 40401,
  "message": "resource not found",
  "error": "plan not found",
  "timestamp": 1699999999
}
```

**处理建议**:
- 验证资源 ID 是否正确
- 检查资源是否已被删除
- 确认 API 路径是否正确

---

### 40901 - Conflict

**说明**: 资源冲突

**HTTP 状态码**: 409 Conflict

**触发场景**:
- 创建重复的资源
- 资源状态冲突
- 并发操作冲突

**常见示例**:
- 用户名已存在
- 同一时间段已有餐饮记录
- 资源正在被其他操作使用

**响应示例**:

```json
{
  "code": 40901,
  "message": "resource conflict",
  "error": "meal already exists for this time period",
  "timestamp": 1699999999
}
```

**处理建议**:
- 检查是否存在重复数据
- 刷新数据后重试
- 避免并发操作同一资源

---

### 42901 - Too Many Requests

**说明**: 请求过于频繁，触发限流

**HTTP 状态码**: 429 Too Many Requests

**触发场景**:
- 超过 API 调用频率限制
- 短时间内发送大量请求
- 触发防刷机制

**常见示例**:
- 1 分钟内登录尝试超过 5 次
- 短时间内大量调用 AI 接口
- 批量操作频率过高

**响应示例**:

```json
{
  "code": 42901,
  "message": "too many requests",
  "error": "rate limit exceeded, please try again later",
  "timestamp": 1699999999
}
```

**处理建议**:
- 等待一段时间后重试
- 实现请求限流和重试机制
- 优化请求频率
- 使用批量接口减少请求次数

---

## 服务器错误码 (50xxx)

### 50001 - Internal Server Error

**说明**: 服务器内部错误

**HTTP 状态码**: 500 Internal Server Error

**触发场景**:
- 服务器内部异常
- 未预期的错误
- 业务逻辑执行失败
- 配置错误

**常见示例**:
- 服务器配置错误
- 代码执行异常
- 资源不足

**响应示例**:

```json
{
  "code": 50001,
  "message": "internal server error",
  "error": "an internal error occurred, please try again later",
  "timestamp": 1699999999
}
```

**处理建议**:
- 稍后重试
- 如果问题持续，联系技术支持
- 检查请求是否触发了服务器异常

---

### 50002 - Database Error

**说明**: 数据库操作错误

**HTTP 状态码**: 500 Internal Server Error

**触发场景**:
- 数据库连接失败
- SQL 执行错误
- 数据库超时
- 事务失败

**常见示例**:
- 数据库连接池耗尽
- 查询超时
- 数据完整性约束违反
- 数据库服务不可用

**响应示例**:

```json
{
  "code": 50002,
  "message": "database error",
  "error": "a database error occurred, please try again later",
  "timestamp": 1699999999
}
```

**处理建议**:
- 稍后重试
- 检查数据是否符合约束条件
- 如果问题持续，联系技术支持

---

### 50003 - AI Service Error

**说明**: AI 服务调用失败

**HTTP 状态码**: 500 Internal Server Error

**触发场景**:
- AI 服务不可用
- AI API 调用超时
- AI 服务返回错误
- API Key 无效或配额不足

**常见示例**:
- OpenAI API 调用失败
- DeepSeek API 超时
- AI 服务配置错误
- API 配额已用完

**响应示例**:

```json
{
  "code": 50003,
  "message": "AI service error",
  "error": "an external service error occurred, please try again later",
  "timestamp": 1699999999
}
```

**处理建议**:
- 检查 AI 服务配置是否正确
- 验证 API Key 是否有效
- 检查 API 配额是否充足
- 稍后重试
- 如果问题持续，联系技术支持

---

### 50004 - Encryption Error

**说明**: 加密/解密操作失败

**HTTP 状态码**: 500 Internal Server Error

**触发场景**:
- 密码加密失败
- 数据解密失败
- 加密密钥错误
- 加密算法异常

**常见示例**:
- 密码哈希生成失败
- Token 签名失败
- 敏感数据加密失败

**响应示例**:

```json
{
  "code": 50004,
  "message": "encryption error",
  "error": "an internal error occurred, please try again later",
  "timestamp": 1699999999
}
```

**处理建议**:
- 稍后重试
- 如果问题持续，联系技术支持
- 检查系统配置

---

## 错误码映射表

### HTTP 状态码与业务错误码对应关系

| HTTP 状态码 | 业务错误码范围 | 说明 |
|------------|--------------|------|
| 200 OK | 0 | 成功 |
| 400 Bad Request | 40000-40999 | 请求参数错误 |
| 401 Unauthorized | 40100-40199 | 未授权 |
| 403 Forbidden | 40300-40399 | 禁止访问 |
| 404 Not Found | 40400-40499 | 资源不存在 |
| 409 Conflict | 40900-40999 | 资源冲突 |
| 429 Too Many Requests | 42900-42999 | 请求过于频繁 |
| 500 Internal Server Error | 50000+ | 服务器错误 |

### 完整错误码列表

| 错误码 | 错误消息 | HTTP 状态码 | 分类 |
|--------|---------|------------|------|
| 0 | success | 200 | 成功 |
| 40001 | invalid parameters | 400 | 参数错误 |
| 40002 | validation error | 400 | 验证错误 |
| 40101 | unauthorized | 401 | 未授权 |
| 40301 | forbidden | 403 | 禁止访问 |
| 40401 | resource not found | 404 | 资源不存在 |
| 40901 | resource conflict | 409 | 资源冲突 |
| 42901 | too many requests | 429 | 请求过于频繁 |
| 50001 | internal server error | 500 | 内部错误 |
| 50002 | database error | 500 | 数据库错误 |
| 50003 | AI service error | 500 | AI 服务错误 |
| 50004 | encryption error | 500 | 加密错误 |

---

## 错误处理最佳实践

### 1. 客户端错误处理流程

```javascript
async function handleAPIRequest(url, options) {
  try {
    const response = await fetch(url, options);
    const data = await response.json();
    
    // 检查业务错误码
    if (data.code !== 0) {
      switch (data.code) {
        case 40101: // Unauthorized
          // Token 过期，尝试刷新
          await refreshToken();
          // 重试原请求
          return handleAPIRequest(url, options);
          
        case 40401: // Not Found
          // 资源不存在
          showError('资源不存在');
          break;
          
        case 42901: // Too Many Requests
          // 请求过于频繁，等待后重试
          await sleep(5000);
          return handleAPIRequest(url, options);
          
        case 50001:
        case 50002:
        case 50003:
          // 服务器错误，稍后重试
          showError('服务器错误，请稍后重试');
          break;
          
        default:
          showError(data.message);
      }
      throw new Error(data.message);
    }
    
    return data.data;
  } catch (error) {
    console.error('API request failed:', error);
    throw error;
  }
}
```

### 2. 重试策略

对于以下错误码，建议实现自动重试机制：

- **42901** (Too Many Requests): 等待后重试
- **50001** (Internal Server Error): 指数退避重试
- **50002** (Database Error): 短暂延迟后重试
- **50003** (AI Service Error): 延迟后重试

### 3. 用户友好的错误提示

| 错误码 | 用户提示 |
|--------|---------|
| 40001 | 请检查输入的信息是否正确 |
| 40002 | 数据格式不正确，请重新输入 |
| 40101 | 登录已过期，请重新登录 |
| 40301 | 您没有权限执行此操作 |
| 40401 | 请求的内容不存在 |
| 40901 | 操作冲突，请刷新后重试 |
| 42901 | 操作过于频繁，请稍后再试 |
| 50001 | 服务器繁忙，请稍后重试 |
| 50002 | 数据保存失败，请稍后重试 |
| 50003 | AI 服务暂时不可用，请稍后重试 |
| 50004 | 系统错误，请联系技术支持 |

### 4. 日志记录

建议记录以下信息用于问题排查：

- 错误码和错误消息
- 请求 URL 和方法
- 请求参数（脱敏）
- 用户 ID（如果已登录）
- 时间戳
- 客户端信息（浏览器、设备等）

---

## 常见问题

### Q: 如何区分客户端错误和服务器错误？

A: 通过错误码的范围：
- 40xxx 开头的是客户端错误，通常是请求参数或权限问题
- 50xxx 开头的是服务器错误，通常是服务端异常

### Q: Token 过期后如何处理？

A: 收到 40101 错误时：
1. 使用 Refresh Token 调用刷新接口获取新的 Access Token
2. 如果刷新失败，引导用户重新登录
3. 获取新 Token 后重试原请求

### Q: 遇到 50xxx 错误应该如何处理？

A: 服务器错误的处理建议：
1. 实现重试机制（建议 2-3 次）
2. 使用指数退避策略
3. 如果持续失败，提示用户稍后重试
4. 记录错误日志便于排查

### Q: 如何避免触发限流（42901）？

A: 限流避免策略：
1. 实现客户端请求限流
2. 使用防抖和节流技术
3. 合理使用批量接口
4. 缓存不常变化的数据
5. 避免短时间内重复请求

### Q: 生产环境和开发环境的错误响应有什么区别？

A: 主要区别在 `error` 字段：
- **开发环境**: 返回详细的错误堆栈和调试信息
- **生产环境**: 返回通用的错误描述，不暴露内部实现细节

---

## 相关文档

- [通用概念](./common-concepts.md) - 了解认证、响应格式等通用概念
- [API 文档总览](./README.md) - 返回 API 文档首页
- [数据模型](./data-models.md) - 了解数据结构定义
