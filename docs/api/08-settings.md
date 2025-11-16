# 设置管理模块

## 概述

设置管理模块提供用户配置管理功能，包括 AI 服务配置和用户个人偏好设置。通过这个模块，用户可以配置 AI 服务提供商（OpenAI、DeepSeek 或自定义）、设置 API 密钥和模型参数，以及管理个人的口味偏好、饮食限制和营养目标。这些设置直接影响 AI 服务的行为和饮食建议的生成。

**核心功能**：
- 获取所有设置（AI 配置和用户偏好）
- 更新 AI 服务配置
- 测试 AI 连接
- 获取用户资料
- 更新用户偏好

**数据特性**：
- AI 配置支持多种提供商
- API 密钥加密存储
- 用户偏好扁平化结构
- 自动填充默认值
- 安全的密钥掩码显示

**应用场景**：
- 配置 AI 服务提供商
- 设置营养目标
- 管理饮食偏好和限制
- 测试 AI 连接状态
- 查看和更新个人资料

---

## 接口列表

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/settings` | 获取所有设置 | 是 |
| PUT | `/api/v1/settings/ai` | 更新 AI 设置 | 是 |
| GET | `/api/v1/settings/ai/test` | 测试 AI 连接 | 是 |
| GET | `/api/v1/user/profile` | 获取用户资料 | 是 |
| PUT | `/api/v1/user/preferences` | 更新用户偏好 | 是 |

---

## 接口详情



### 获取所有设置

**接口**: `GET /api/v1/settings`

**说明**: 获取用户的所有设置，包括 AI 配置和用户偏好。AI 配置中的 API 密钥会被掩码处理（只显示前 4 位和后 4 位），确保安全性。如果用户还没有配置 AI 服务，ai_config 字段将为 null。用户偏好包含口味偏好、饮食限制和每日营养目标。

**认证**: 是


#### 请求参数

无需任何参数，系统会自动获取当前用户的设置。

#### 请求示例

```bash
# 获取所有设置
curl -X GET "http://localhost:9090/api/v1/settings" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 使用环境变量存储 Token
export TOKEN="YOUR_ACCESS_TOKEN"
curl -X GET "http://localhost:9090/api/v1/settings" \
  -H "Authorization: Bearer $TOKEN"
```

#### 响应示例

**成功响应 (200) - 已配置 AI**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "ai_config": {
      "provider": "openai",
      "api_endpoint": "https://api.openai.com/v1",
      "api_key_masked": "sk-p****xyz1",
      "model": "gpt-3.5-turbo",
      "temperature": 0.7,
      "max_tokens": 1000
    },
    "user_preferences": {
      "id": 1,
      "user_id": 1,
      "taste_preferences": "喜欢清淡口味，偏好蒸煮烹饪方式",
      "dietary_restrictions": "对海鲜过敏，不吃辛辣食物",
      "daily_calories_goal": 2000,
      "daily_protein_goal": 150,
      "daily_carbs_goal": 250,
      "daily_fat_goal": 70,
      "daily_fiber_goal": 30,
      "created_at": "2024-11-01T10:00:00Z",
      "updated_at": "2024-11-15T14:30:00Z"
    }
  },
  "timestamp": 1699999999
}
```

**成功响应 (200) - 未配置 AI**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "ai_config": null,
    "user_preferences": {
      "id": 1,
      "user_id": 1,
      "taste_preferences": "",
      "dietary_restrictions": "",
      "daily_calories_goal": 2000,
      "daily_protein_goal": 150,
      "daily_carbs_goal": 250,
      "daily_fat_goal": 70,
      "daily_fiber_goal": 30,
      "created_at": "2024-11-01T10:00:00Z",
      "updated_at": "2024-11-01T10:00:00Z"
    }
  },
  "timestamp": 1699999999
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| ai_config | object/null | AI 配置信息，未配置时为 null |
| ai_config.provider | string | AI 提供商（openai/deepseek/custom） |
| ai_config.api_endpoint | string | API 端点 URL |
| ai_config.api_key_masked | string | 掩码后的 API 密钥（前 4 位 + **** + 后 4 位） |
| ai_config.model | string | 使用的模型名称 |
| ai_config.temperature | float | 温度参数（0-2，控制随机性） |
| ai_config.max_tokens | int | 最大 Token 数量 |
| user_preferences | object | 用户偏好设置 |
| user_preferences.id | int | 偏好设置 ID |
| user_preferences.user_id | int | 用户 ID |
| user_preferences.taste_preferences | string | 口味偏好描述 |
| user_preferences.dietary_restrictions | string | 饮食限制描述 |
| user_preferences.daily_calories_goal | int | 每日热量目标（千卡） |
| user_preferences.daily_protein_goal | int | 每日蛋白质目标（克） |
| user_preferences.daily_carbs_goal | int | 每日碳水化合物目标（克） |
| user_preferences.daily_fat_goal | int | 每日脂肪目标（克） |
| user_preferences.daily_fiber_goal | int | 每日纤维目标（克） |


**错误响应 (401)**:

```json
{
  "code": 40101,
  "message": "unauthorized",
  "error": "user not authenticated",
  "timestamp": 1699999999
}
```

**错误响应 (500)**:

```json
{
  "code": 50001,
  "message": "internal error",
  "error": "failed to get user preferences",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | 服务器内部错误、数据库查询失败 |

---


### 更新 AI 设置

**接口**: `PUT /api/v1/settings/ai`

**说明**: 更新用户的 AI 服务配置，包括提供商、API 端点、API 密钥、模型和参数。支持 OpenAI、DeepSeek 和自定义提供商。API 密钥会被加密存储。如果不提供某些可选字段，系统会保留现有值或使用默认值。更新后的配置会立即生效，影响后续的 AI 服务调用。

**认证**: 是


#### 请求参数

##### 请求体

```json
{
  "provider": "openai",
  "api_endpoint": "https://api.openai.com/v1",
  "api_key": "sk-proj1234567890abcdefghijklmnopqrstuvwxyz1",
  "model": "gpt-3.5-turbo",
  "temperature": 0.7,
  "max_tokens": 1000,
  "is_active": true
}
```

| 字段 | 类型 | 必填 | 说明 | 验证规则 |
|------|------|------|------|----------|
| provider | string | 是 | AI 提供商 | 必须是 openai、deepseek 或 custom |
| api_endpoint | string | 否 | API 端点 URL | 必须是有效的 URL，最大 500 字符 |
| api_key | string | 条件 | API 密钥 | 新配置必填，更新时可选；最小 10 字符，最大 500 字符 |
| model | string | 否 | 模型名称 | 最大 100 字符，不提供时使用默认值 |
| temperature | float | 否 | 温度参数 | 0-2 之间，控制随机性，默认 0.7 |
| max_tokens | int | 否 | 最大 Token 数 | 1-32000 之间，默认 1000 |
| is_active | bool | 否 | 是否激活 | 总是设置为 true |

**默认模型**：
- openai: gpt-3.5-turbo
- deepseek: deepseek-chat
- custom: default

#### 请求示例

```bash
# 配置 OpenAI
curl -X PUT "http://localhost:9090/api/v1/settings/ai" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "openai",
    "api_endpoint": "https://api.openai.com/v1",
    "api_key": "sk-proj1234567890abcdefghijklmnopqrstuvwxyz1",
    "model": "gpt-3.5-turbo",
    "temperature": 0.7,
    "max_tokens": 1000
  }'

# 配置 DeepSeek
curl -X PUT "http://localhost:9090/api/v1/settings/ai" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "deepseek",
    "api_endpoint": "https://api.deepseek.com/v1",
    "api_key": "sk-1234567890abcdefghijklmnopqrstuvwxyz",
    "model": "deepseek-chat",
    "temperature": 0.8,
    "max_tokens": 2000
  }'

# 更新现有配置（只更新部分字段）
curl -X PUT "http://localhost:9090/api/v1/settings/ai" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "openai",
    "model": "gpt-4",
    "temperature": 0.5
  }'

# 配置自定义提供商
curl -X PUT "http://localhost:9090/api/v1/settings/ai" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "custom",
    "api_endpoint": "https://custom-ai.example.com/v1",
    "api_key": "custom-key-1234567890",
    "model": "custom-model",
    "temperature": 0.7,
    "max_tokens": 1500
  }'
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "AI settings updated successfully",
  "data": null,
  "timestamp": 1699999999
}
```

**错误响应 (400) - 参数错误**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "Key: 'UpdateAISettingsRequest.Provider' Error:Field validation for 'Provider' failed on the 'oneof' tag",
  "timestamp": 1699999999
}
```

**错误响应 (400) - 缺少 API 密钥**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "API key is required for new configuration",
  "timestamp": 1699999999
}
```

**错误响应 (401)**:

```json
{
  "code": 40101,
  "message": "unauthorized",
  "error": "user not authenticated",
  "timestamp": 1699999999
}
```

**错误响应 (500)**:

```json
{
  "code": 50001,
  "message": "internal error",
  "error": "failed to update AI settings",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | 请求参数不符合要求、provider 不在允许列表中 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | 服务器内部错误、数据库更新失败 |

#### 注意事项

1. **API 密钥安全**：API 密钥会被加密存储，不会以明文形式保存
2. **首次配置**：首次配置时必须提供 api_key
3. **更新配置**：更新时如果不提供 api_key，会保留现有的密钥
4. **字段保留**：不提供的可选字段会保留现有值或使用默认值
5. **默认值**：model、temperature、max_tokens 都有默认值
6. **立即生效**：配置更新后立即生效，影响后续的 AI 调用
7. **提供商切换**：可以随时切换不同的 AI 提供商
8. **参数范围**：temperature 必须在 0-2 之间，max_tokens 必须在 1-32000 之间
9. **URL 验证**：api_endpoint 必须是有效的 URL 格式
10. **测试建议**：更新配置后建议使用测试接口验证连接

---


### 测试 AI 连接

**接口**: `GET /api/v1/settings/ai/test`

**说明**: 测试当前配置的 AI 服务是否可用。系统会使用用户配置的 AI 设置发送一个测试请求，验证 API 端点、密钥和配置是否正确。如果测试失败，会返回详细的错误信息，帮助用户排查问题。建议在配置或更新 AI 设置后使用此接口进行验证。

**认证**: 是


#### 请求参数

无需任何参数，系统会自动使用当前用户的 AI 配置进行测试。

#### 请求示例

```bash
# 测试 AI 连接
curl -X GET "http://localhost:9090/api/v1/settings/ai/test" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 使用环境变量存储 Token
export TOKEN="YOUR_ACCESS_TOKEN"
curl -X GET "http://localhost:9090/api/v1/settings/ai/test" \
  -H "Authorization: Bearer $TOKEN"
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "AI connection test successful",
  "data": null,
  "timestamp": 1699999999
}
```

**错误响应 (401)**:

```json
{
  "code": 40101,
  "message": "unauthorized",
  "error": "user not authenticated",
  "timestamp": 1699999999
}
```

**错误响应 (500) - 未配置 AI**:

```json
{
  "code": 50001,
  "message": "internal error",
  "error": "AI connection test failed: no AI configuration found",
  "timestamp": 1699999999
}
```

**错误响应 (500) - API 密钥无效**:

```json
{
  "code": 50001,
  "message": "internal error",
  "error": "AI connection test failed: invalid API key",
  "timestamp": 1699999999
}
```

**错误响应 (500) - 网络错误**:

```json
{
  "code": 50001,
  "message": "internal error",
  "error": "AI connection test failed: connection timeout",
  "timestamp": 1699999999
}
```

**错误响应 (500) - 端点错误**:

```json
{
  "code": 50001,
  "message": "internal error",
  "error": "AI connection test failed: invalid endpoint URL",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | AI 连接失败、配置错误、网络问题 |

#### 注意事项

1. **测试时机**：建议在配置或更新 AI 设置后立即测试
2. **错误信息**：失败时会返回详细的错误信息，帮助排查问题
3. **网络要求**：需要服务器能够访问配置的 AI 端点
4. **超时设置**：测试请求有超时限制，避免长时间等待
5. **配置检查**：测试前确保已经配置了 AI 设置
6. **API 配额**：测试会消耗 AI 服务的 API 配额（通常很少）
7. **频率限制**：避免频繁测试，以免触发 AI 服务的频率限制
8. **故障排查**：
   - 检查 API 密钥是否正确
   - 检查 API 端点 URL 是否正确
   - 检查网络连接是否正常
   - 检查 AI 服务是否可用
9. **成功标准**：只要能成功调用 AI 服务并获得响应，就算测试成功
10. **失败处理**：测试失败不会影响现有配置，可以继续使用

---


### 获取用户资料

**接口**: `GET /api/v1/user/profile`

**说明**: 获取当前用户的资料信息，包括用户偏好设置。返回的数据与 GET /api/v1/settings 接口中的 user_preferences 字段相同，但这个接口专注于用户资料，不包含 AI 配置信息。适用于用户资料页面或个人设置页面。

**认证**: 是


#### 请求参数

无需任何参数，系统会自动获取当前用户的资料。

#### 请求示例

```bash
# 获取用户资料
curl -X GET "http://localhost:9090/api/v1/user/profile" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 使用环境变量存储 Token
export TOKEN="YOUR_ACCESS_TOKEN"
curl -X GET "http://localhost:9090/api/v1/user/profile" \
  -H "Authorization: Bearer $TOKEN"
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "user_id": 1,
    "taste_preferences": "喜欢清淡口味，偏好蒸煮烹饪方式",
    "dietary_restrictions": "对海鲜过敏，不吃辛辣食物",
    "daily_calories_goal": 2000,
    "daily_protein_goal": 150,
    "daily_carbs_goal": 250,
    "daily_fat_goal": 70,
    "daily_fiber_goal": 30,
    "created_at": "2024-11-01T10:00:00Z",
    "updated_at": "2024-11-15T14:30:00Z"
  },
  "timestamp": 1699999999
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 偏好设置 ID |
| user_id | int | 用户 ID |
| taste_preferences | string | 口味偏好描述 |
| dietary_restrictions | string | 饮食限制描述 |
| daily_calories_goal | int | 每日热量目标（千卡） |
| daily_protein_goal | int | 每日蛋白质目标（克） |
| daily_carbs_goal | int | 每日碳水化合物目标（克） |
| daily_fat_goal | int | 每日脂肪目标（克） |
| daily_fiber_goal | int | 每日纤维目标（克） |
| created_at | string | 创建时间（ISO 8601 格式） |
| updated_at | string | 更新时间（ISO 8601 格式） |


**错误响应 (401)**:

```json
{
  "code": 40101,
  "message": "unauthorized",
  "error": "user not authenticated",
  "timestamp": 1699999999
}
```

**错误响应 (500)**:

```json
{
  "code": 50001,
  "message": "internal error",
  "error": "failed to get user profile",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | 服务器内部错误、数据库查询失败 |

---


### 更新用户偏好

**接口**: `PUT /api/v1/user/preferences`

**说明**: 更新用户的个人偏好设置，包括口味偏好、饮食限制和每日营养目标。这些设置会影响 AI 生成的饮食建议和计划。所有字段都是可选的，不提供的字段会保留现有值或使用默认值。营养目标会在 Dashboard 和营养分析模块中使用。

**认证**: 是


#### 请求参数

##### 请求体

```json
{
  "taste_preferences": "喜欢清淡口味，偏好蒸煮烹饪方式",
  "dietary_restrictions": "对海鲜过敏，不吃辛辣食物",
  "daily_calories_goal": 2000,
  "daily_protein_goal": 150,
  "daily_carbs_goal": 250,
  "daily_fat_goal": 70,
  "daily_fiber_goal": 30
}
```

| 字段 | 类型 | 必填 | 说明 | 验证规则 |
|------|------|------|------|----------|
| taste_preferences | string | 否 | 口味偏好描述 | 最大 500 字符 |
| dietary_restrictions | string | 否 | 饮食限制描述 | 最大 500 字符 |
| daily_calories_goal | int | 否 | 每日热量目标 | 800-10000 千卡，默认 2000 |
| daily_protein_goal | int | 否 | 每日蛋白质目标 | 0-500 克，默认 150 |
| daily_carbs_goal | int | 否 | 每日碳水化合物目标 | 0-1000 克，默认 250 |
| daily_fat_goal | int | 否 | 每日脂肪目标 | 0-500 克，默认 70 |
| daily_fiber_goal | int | 否 | 每日纤维目标 | 0-200 克，默认 30 |

#### 请求示例

```bash
# 更新所有偏好
curl -X PUT "http://localhost:9090/api/v1/user/preferences" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "taste_preferences": "喜欢清淡口味，偏好蒸煮烹饪方式",
    "dietary_restrictions": "对海鲜过敏，不吃辛辣食物",
    "daily_calories_goal": 2000,
    "daily_protein_goal": 150,
    "daily_carbs_goal": 250,
    "daily_fat_goal": 70,
    "daily_fiber_goal": 30
  }'

# 只更新营养目标
curl -X PUT "http://localhost:9090/api/v1/user/preferences" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "daily_calories_goal": 1800,
    "daily_protein_goal": 120,
    "daily_carbs_goal": 200,
    "daily_fat_goal": 60,
    "daily_fiber_goal": 25
  }'

# 只更新口味偏好和饮食限制
curl -X PUT "http://localhost:9090/api/v1/user/preferences" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "taste_preferences": "喜欢重口味，偏好煎炒烹饪方式",
    "dietary_restrictions": "素食主义者，不吃任何肉类"
  }'

# 清空口味偏好和饮食限制
curl -X PUT "http://localhost:9090/api/v1/user/preferences" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "taste_preferences": "",
    "dietary_restrictions": ""
  }'
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "user preferences updated successfully",
  "data": null,
  "timestamp": 1699999999
}
```

**错误响应 (400) - 参数错误**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "Key: 'UpdateUserPreferencesRequest.DailyCaloriesGoal' Error:Field validation for 'DailyCaloriesGoal' failed on the 'min' tag",
  "timestamp": 1699999999
}
```

**错误响应 (401)**:

```json
{
  "code": 40101,
  "message": "unauthorized",
  "error": "user not authenticated",
  "timestamp": 1699999999
}
```

**错误响应 (500)**:

```json
{
  "code": 50001,
  "message": "internal error",
  "error": "failed to update user preferences",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | 请求参数不符合要求、数值超出范围 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | 服务器内部错误、数据库更新失败 |

#### 注意事项

1. **可选字段**：所有字段都是可选的，可以只更新部分字段
2. **字段保留**：不提供的字段会保留现有值
3. **默认值**：如果字段为 0 或空字符串，会使用默认值
4. **立即生效**：更新后立即生效，影响后续的 AI 建议和营养分析
5. **营养目标**：营养目标会在 Dashboard 和营养分析模块中使用
6. **AI 建议**：口味偏好和饮食限制会影响 AI 生成的饮食建议
7. **合理范围**：营养目标应该设置在合理范围内
8. **清空字段**：可以通过传递空字符串清空文本字段
9. **数值验证**：营养目标的数值会被验证，确保在合理范围内
10. **建议设置**：
    - 热量：根据个人情况设置，一般 1500-3000 千卡
    - 蛋白质：体重（kg）× 1.5-2.0 克
    - 碳水化合物：总热量的 45-65%
    - 脂肪：总热量的 20-35%
    - 纤维：25-35 克

---


## 数据模型

### AISettings 模型

AI 设置模型：

```typescript
interface AISettings {
  provider: 'openai' | 'deepseek' | 'custom';  // AI 提供商
  api_endpoint: string;                         // API 端点 URL
  api_key_masked: string;                       // 掩码后的 API 密钥
  model: string;                                // 模型名称
  temperature: number;                          // 温度参数（0-2）
  max_tokens: number;                           // 最大 Token 数量
}
```

### UserPreferences 模型

用户偏好模型：

```typescript
interface UserPreferences {
  id: number;                      // 偏好设置 ID
  user_id: number;                 // 用户 ID
  taste_preferences: string;       // 口味偏好描述
  dietary_restrictions: string;    // 饮食限制描述
  daily_calories_goal: number;     // 每日热量目标（千卡）
  daily_protein_goal: number;      // 每日蛋白质目标（克）
  daily_carbs_goal: number;        // 每日碳水化合物目标（克）
  daily_fat_goal: number;          // 每日脂肪目标（克）
  daily_fiber_goal: number;        // 每日纤维目标（克）
  created_at: string;              // 创建时间（ISO 8601）
  updated_at: string;              // 更新时间（ISO 8601）
}
```

### UpdateAISettingsRequest 模型

更新 AI 设置请求模型：

```typescript
interface UpdateAISettingsRequest {
  provider: 'openai' | 'deepseek' | 'custom';  // AI 提供商（必填）
  api_endpoint?: string;                        // API 端点 URL（可选）
  api_key?: string;                             // API 密钥（条件必填）
  model?: string;                               // 模型名称（可选）
  temperature?: number;                         // 温度参数（可选，0-2）
  max_tokens?: number;                          // 最大 Token 数量（可选，1-32000）
  is_active?: boolean;                          // 是否激活（可选）
}
```

### UpdateUserPreferencesRequest 模型

更新用户偏好请求模型：

```typescript
interface UpdateUserPreferencesRequest {
  taste_preferences?: string;       // 口味偏好描述（可选，最大 500 字符）
  dietary_restrictions?: string;    // 饮食限制描述（可选，最大 500 字符）
  daily_calories_goal?: number;     // 每日热量目标（可选，800-10000）
  daily_protein_goal?: number;      // 每日蛋白质目标（可选，0-500）
  daily_carbs_goal?: number;        // 每日碳水化合物目标（可选，0-1000）
  daily_fat_goal?: number;          // 每日脂肪目标（可选，0-500）
  daily_fiber_goal?: number;        // 每日纤维目标（可选，0-200）
}
```

---

## 使用场景

### 场景 1：首次配置 AI 服务

用户首次使用系统，需要配置 AI 服务：

```bash
# 1. 配置 OpenAI
curl -X PUT "http://localhost:9090/api/v1/settings/ai" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "openai",
    "api_endpoint": "https://api.openai.com/v1",
    "api_key": "sk-proj1234567890abcdefghijklmnopqrstuvwxyz1",
    "model": "gpt-3.5-turbo",
    "temperature": 0.7,
    "max_tokens": 1000
  }'

# 2. 测试连接
curl -X GET "http://localhost:9090/api/v1/settings/ai/test" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 3. 获取所有设置确认
curl -X GET "http://localhost:9090/api/v1/settings" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 场景 2：设置个人偏好

用户设置个人的口味偏好和营养目标：

```bash
curl -X PUT "http://localhost:9090/api/v1/user/preferences" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "taste_preferences": "喜欢清淡口味，偏好蒸煮烹饪方式",
    "dietary_restrictions": "对海鲜过敏，不吃辛辣食物",
    "daily_calories_goal": 2000,
    "daily_protein_goal": 150,
    "daily_carbs_goal": 250,
    "daily_fat_goal": 70,
    "daily_fiber_goal": 30
  }'
```

### 场景 3：切换 AI 提供商

用户想从 OpenAI 切换到 DeepSeek：

```bash
# 1. 更新 AI 配置
curl -X PUT "http://localhost:9090/api/v1/settings/ai" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "deepseek",
    "api_endpoint": "https://api.deepseek.com/v1",
    "api_key": "sk-1234567890abcdefghijklmnopqrstuvwxyz",
    "model": "deepseek-chat"
  }'

# 2. 测试新配置
curl -X GET "http://localhost:9090/api/v1/settings/ai/test" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 场景 4：调整营养目标

用户想调整每日营养目标：

```bash
curl -X PUT "http://localhost:9090/api/v1/user/preferences" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "daily_calories_goal": 1800,
    "daily_protein_goal": 120,
    "daily_carbs_goal": 200,
    "daily_fat_goal": 60,
    "daily_fiber_goal": 25
  }'
```

### 场景 5：查看当前设置

用户想查看当前的所有设置：

```bash
curl -X GET "http://localhost:9090/api/v1/settings" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 场景 6：更新 AI 模型参数

用户想调整 AI 模型的参数以获得不同的响应风格：

```bash
# 更保守的响应（降低 temperature）
curl -X PUT "http://localhost:9090/api/v1/settings/ai" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "openai",
    "temperature": 0.3,
    "max_tokens": 1500
  }'

# 更有创意的响应（提高 temperature）
curl -X PUT "http://localhost:9090/api/v1/settings/ai" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "openai",
    "temperature": 1.2,
    "max_tokens": 2000
  }'
```

---

## 最佳实践

### 1. AI 配置管理

- **安全存储**：不要在客户端存储 API 密钥明文
- **测试验证**：配置后立即测试连接
- **定期检查**：定期检查 AI 服务状态
- **备用方案**：准备多个 AI 提供商作为备用
- **参数调优**：根据实际效果调整 temperature 和 max_tokens

### 2. 用户偏好设置

- **引导设置**：新用户注册后引导设置偏好
- **合理目标**：帮助用户设置合理的营养目标
- **定期更新**：建议用户定期更新偏好设置
- **详细描述**：鼓励用户详细描述口味偏好和饮食限制
- **专业建议**：提供营养目标的参考建议

### 3. 安全性

- **密钥保护**：API 密钥加密存储，传输使用 HTTPS
- **权限控制**：只允许用户修改自己的设置
- **输入验证**：严格验证所有输入参数
- **错误处理**：不在错误信息中泄露敏感信息
- **审计日志**：记录设置变更日志

### 4. 用户体验

- **即时反馈**：设置更新后立即反馈结果
- **错误提示**：提供清晰的错误提示和解决建议
- **默认值**：提供合理的默认值
- **表单验证**：在客户端进行表单验证
- **保存提示**：提示用户保存设置

### 5. 性能优化

- **缓存设置**：在客户端缓存设置数据
- **批量更新**：尽量批量更新多个字段
- **减少请求**：避免频繁请求设置接口
- **异步处理**：使用异步方式更新设置
- **加载状态**：显示加载和保存状态

### 6. AI 服务选择

- **OpenAI**：功能强大，响应质量高，但成本较高
- **DeepSeek**：性价比高，中文支持好，适合国内用户
- **Custom**：灵活性高，可以使用自己部署的模型
- **选择建议**：根据预算、需求和地区选择合适的提供商

### 7. 营养目标设置

- **个性化**：根据个人情况（年龄、性别、活动量）设置
- **参考标准**：提供营养学标准作为参考
- **渐进调整**：不要一次性大幅调整目标
- **监测效果**：定期监测目标达成情况
- **专业咨询**：建议咨询营养师或医生

### 8. 错误处理

- **友好提示**：将技术错误转换为用户友好的提示
- **重试机制**：对于网络错误提供重试选项
- **回滚支持**：更新失败时保留原有设置
- **日志记录**：记录错误日志便于排查
- **用户支持**：提供联系支持的渠道

---


## 常见问题

### Q: 如何获取 OpenAI API 密钥？

A: 
- 访问 OpenAI 官网 (https://platform.openai.com)
- 注册并登录账号
- 进入 API Keys 页面
- 点击 "Create new secret key" 创建新密钥
- 复制密钥并妥善保管（只显示一次）
- 在系统中配置该密钥

### Q: API 密钥是否安全？

A: 
- API 密钥在服务器端加密存储
- 传输过程使用 HTTPS 加密
- 查询时只返回掩码后的密钥
- 不会在日志中记录明文密钥
- 建议定期更换 API 密钥

### Q: 为什么 AI 连接测试失败？

A: 
- 检查 API 密钥是否正确
- 检查 API 端点 URL 是否正确
- 检查网络连接是否正常
- 检查 AI 服务是否可用
- 检查是否有足够的 API 配额
- 查看错误信息获取详细原因

### Q: 可以同时配置多个 AI 提供商吗？

A: 
- 当前版本只支持配置一个活跃的 AI 提供商
- 可以随时切换不同的提供商
- 切换时会覆盖之前的配置
- 建议在切换前记录原有配置

### Q: temperature 参数如何设置？

A: 
- temperature 控制 AI 响应的随机性
- 0.0-0.3：更保守、更确定的响应
- 0.4-0.7：平衡的响应（推荐）
- 0.8-1.5：更有创意、更多样的响应
- 1.6-2.0：非常随机的响应
- 建议从 0.7 开始，根据效果调整

### Q: max_tokens 参数如何设置？

A: 
- max_tokens 控制 AI 响应的最大长度
- 1 token ≈ 0.75 个英文单词或 0.5 个中文字符
- 500-1000：简短回答
- 1000-2000：中等长度回答（推荐）
- 2000-4000：详细回答
- 注意：更大的值会消耗更多 API 配额

### Q: 如何设置合理的营养目标？

A: 
- **热量**：根据基础代谢率和活动量计算
  - 久坐：体重（kg）× 25-30 千卡
  - 轻度活动：体重（kg）× 30-35 千卡
  - 中度活动：体重（kg）× 35-40 千卡
- **蛋白质**：体重（kg）× 1.5-2.0 克
- **碳水化合物**：总热量的 45-65%
- **脂肪**：总热量的 20-35%
- **纤维**：25-35 克
- 建议咨询营养师或医生

### Q: 口味偏好和饮食限制如何描述？

A: 
- **口味偏好**：
  - 描述喜欢的口味（清淡、重口味、甜、咸等）
  - 描述喜欢的烹饪方式（蒸、煮、炒、烤等）
  - 描述喜欢的食材类型
- **饮食限制**：
  - 列出过敏食材
  - 列出不能吃的食物（宗教、健康原因等）
  - 列出特殊饮食要求（素食、低糖、低盐等）
- 描述越详细，AI 建议越准确

### Q: 更新设置后多久生效？

A: 
- 设置更新后立即生效
- AI 配置影响后续的 AI 服务调用
- 用户偏好影响后续的饮食建议和计划生成
- Dashboard 数据会在下次请求时使用新的营养目标
- 不影响已经生成的历史数据

### Q: 可以恢复默认设置吗？

A: 
- 当前版本不支持一键恢复默认设置
- 可以手动设置为默认值：
  - 热量：2000 千卡
  - 蛋白质：150 克
  - 碳水化合物：250 克
  - 脂肪：70 克
  - 纤维：30 克
- 可以清空口味偏好和饮食限制（传递空字符串）

### Q: 如何切换到自定义 AI 提供商？

A: 
- 设置 provider 为 "custom"
- 提供自定义的 API 端点 URL
- 提供相应的 API 密钥
- 设置模型名称
- 确保自定义服务兼容 OpenAI API 格式
- 测试连接确保配置正确

### Q: 为什么获取设置时 ai_config 为 null？

A: 
- 用户还没有配置 AI 服务
- 需要先调用 PUT /api/v1/settings/ai 配置 AI
- 配置后 ai_config 会返回配置信息
- 这是正常现象，不是错误

### Q: 可以导出和导入设置吗？

A: 
- 当前版本不支持导出和导入设置
- 可以通过 API 获取设置数据并保存
- 可以通过 API 重新设置保存的数据
- 未来版本可能会添加导出导入功能

### Q: 设置数据会被其他用户看到吗？

A: 
- 不会，出于隐私保护
- 每个用户只能看到和修改自己的设置
- 数据通过认证 Token 进行隔离
- API 密钥加密存储，即使管理员也看不到明文

### Q: 如何处理 API 配额不足的问题？

A: 
- 检查 AI 服务账户的配额使用情况
- 升级 AI 服务套餐
- 优化 max_tokens 参数减少消耗
- 减少不必要的 AI 调用
- 考虑切换到其他 AI 提供商

---

## 相关文档

- [认证模块](./01-authentication.md) - 了解如何获取认证 Token
- [AI 服务模块](./05-ai-services.md) - 了解 AI 服务的使用
- [Dashboard 模块](./07-dashboard.md) - 了解营养目标在 Dashboard 中的使用
- [营养分析模块](./06-nutrition.md) - 了解营养目标在营养分析中的使用
- [饮食计划模块](./04-plans.md) - 了解用户偏好如何影响计划生成
- [数据模型](./data-models.md) - 查看数据模型的完整定义
- [通用概念](./common-concepts.md) - 了解认证、响应格式等通用概念
- [错误码说明](./error-codes.md) - 查看所有错误码的详细说明
- [API 文档总览](./README.md) - 返回 API 文档首页

