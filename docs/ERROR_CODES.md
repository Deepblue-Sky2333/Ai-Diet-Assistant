# 错误码规范文档

## 概述

本文档定义了 AI Diet Assistant 后端系统使用的标准错误码。所有 API 端点必须使用这些标准错误码，以确保错误处理的一致性和可维护性。

## 标准错误码

根据需求文档 3.2，系统使用以下 7 个标准错误码：

| 错误码 | HTTP 状态码 | 说明 | 使用场景 |
|--------|------------|------|---------|
| 0 | 200 | 成功 | 请求成功处理 |
| 40001 | 400 | 参数错误 | 请求参数验证失败、格式错误、缺少必需参数 |
| 40101 | 401 | 未授权 | 未提供认证信息、Token 无效或过期、密码错误 |
| 40401 | 404 | 资源不存在 | 请求的资源（食材、餐饮记录、计划等）不存在 |
| 42901 | 429 | 请求过于频繁 | 触发限流机制 |
| 50001 | 500 | 内部错误 | 服务器内部错误、加密错误等 |
| 50002 | 500 | 数据库错误 | 数据库连接失败、查询错误等 |
| 50003 | 500 | 外部服务错误 | AI 服务调用失败、第三方 API 错误等 |

## 错误响应格式

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": { ... },
  "timestamp": 1234567890
}
```

### 分页响应

```json
{
  "code": 0,
  "message": "success",
  "data": [ ... ],
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
  "error": "详细错误信息（仅开发环境）",
  "timestamp": 1234567890
}
```

## 使用指南

### 在 Handler 中使用

推荐使用便捷函数创建错误：

```go
// 参数错误
if err := c.ShouldBindJSON(&req); err != nil {
    utils.Error(c, utils.NewInvalidParamsError("invalid request parameters", err))
    return
}

// 未授权
if token == "" {
    utils.Error(c, utils.NewUnauthorizedError("missing authorization token", nil))
    return
}

// 资源不存在
food, err := h.foodService.GetFood(userID, foodID)
if err != nil {
    utils.Error(c, utils.NewNotFoundError("food not found", err))
    return
}

// 内部错误
if err := h.service.DoSomething(); err != nil {
    utils.Error(c, utils.NewInternalError("failed to process request", err))
    return
}

// 数据库错误
if err := h.repo.Save(data); err != nil {
    utils.Error(c, utils.NewDatabaseError("failed to save data", err))
    return
}

// AI 服务错误
response, err := h.aiService.Chat(ctx, message)
if err != nil {
    utils.Error(c, utils.NewAIServiceError("AI service unavailable", err))
    return
}
```

### 便捷函数列表

- `NewInvalidParamsError(message, err)` - 创建参数错误 (40001)
- `NewUnauthorizedError(message, err)` - 创建未授权错误 (40101)
- `NewNotFoundError(message, err)` - 创建资源不存在错误 (40401)
- `NewTooManyRequestsError(message, err)` - 创建限流错误 (42901)
- `NewInternalError(message, err)` - 创建内部错误 (50001)
- `NewDatabaseError(message, err)` - 创建数据库错误 (50002)
- `NewAIServiceError(message, err)` - 创建 AI 服务错误 (50003)

如果 message 参数为空字符串，将使用默认错误消息。

### 通用方式

也可以使用通用的 `NewAppError` 函数：

```go
utils.Error(c, utils.NewAppError(utils.CodeInvalidParams, "custom message", err))
```

## 错误日志

系统会自动记录所有错误到日志，包括：

- 错误码和消息
- HTTP 状态码
- 请求路径和方法
- 客户端 IP
- 用户 ID 和用户名（如果已认证）
- 底层错误详情

### 日志级别

- **Warn 级别**：客户端错误（40001, 40101, 40401, 42901）
- **Error 级别**：服务器错误（50001, 50002, 50003）

## 安全考虑

### 生产环境

在生产环境（`gin.ReleaseMode`）中：
- 不返回详细的错误信息
- 只返回通用的错误消息
- 敏感信息（如数据库连接字符串、文件路径）会被清理

### 开发环境

在开发环境中：
- 返回清理后的详细错误信息
- 帮助开发者快速定位问题
- 仍然会清理敏感信息（密码、密钥等）

## 错误信息清理

系统会自动清理错误消息中的敏感信息：

1. **数据库连接字符串**：密码部分替换为 `***`
   - `user:password@tcp(host:3306)/db` → `user:***@tcp(host:3306)/db`

2. **文件路径**：只保留文件名，移除完整路径

3. **SQL 语句**：替换为 `[SQL query]`

## 迁移指南

### 已移除的错误码

以下错误码已被移除，请使用标准错误码替代：

| 旧错误码 | 新错误码 | 说明 |
|---------|---------|------|
| CodeForbidden (40301) | CodeUnauthorized (40101) | 统一使用未授权错误 |
| CodeConflict (40901) | CodeInvalidParams (40001) | 资源冲突视为参数错误 |
| CodeValidationError (40002) | CodeInvalidParams (40001) | 验证错误统一为参数错误 |
| CodeEncryptionError (50004) | CodeInternalError (50001) | 加密错误视为内部错误 |

### 更新步骤

1. 搜索代码中使用旧错误码的地方
2. 根据上表替换为新的标准错误码
3. 更新错误消息以提供更清晰的说明
4. 测试确保错误处理正确

## 最佳实践

1. **始终提供有意义的错误消息**
   ```go
   // 好的做法
   utils.NewInvalidParamsError("email format is invalid", err)
   
   // 不好的做法
   utils.NewInvalidParamsError("", err)  // 使用默认消息
   ```

2. **传递底层错误**
   ```go
   // 好的做法
   if err := service.DoSomething(); err != nil {
       utils.Error(c, utils.NewInternalError("failed to process", err))
       return
   }
   
   // 不好的做法
   if err := service.DoSomething(); err != nil {
       utils.Error(c, utils.NewInternalError("failed to process", nil))  // 丢失了错误上下文
       return
   }
   ```

3. **选择合适的错误码**
   - 参数验证失败 → 40001
   - 认证失败 → 40101
   - 资源不存在 → 40401
   - 业务逻辑错误 → 50001
   - 数据库操作失败 → 50002
   - 外部服务调用失败 → 50003

4. **不要在 Service 层使用错误码**
   - Service 层返回标准的 Go error
   - Handler 层负责将 error 转换为 AppError

## 测试

确保所有 API 端点的错误处理都经过测试：

```go
func TestCreateFood_InvalidParams(t *testing.T) {
    // 测试参数验证
    // 期望返回 40001
}

func TestGetFood_NotFound(t *testing.T) {
    // 测试资源不存在
    // 期望返回 40401
}

func TestCreateFood_Unauthorized(t *testing.T) {
    // 测试未授权访问
    // 期望返回 40101
}
```

## 参考

- 需求文档：`.kiro/specs/backend-cleanup/requirements.md` - 需求 3.2
- 设计文档：`.kiro/specs/backend-cleanup/design.md` - 错误处理章节
- 实现代码：`internal/utils/response.go`
