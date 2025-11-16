# 饮食计划模块

## 概述

饮食计划模块提供基于 AI 的智能饮食计划生成和管理功能。用户可以使用 AI 生成未来几天的饮食计划，系统会根据用户的食材库、营养目标和个人偏好，智能推荐合理的餐饮搭配。生成的计划可以查询、修改、删除，也可以标记为完成并自动转换为餐饮记录。

**核心功能**：
- 使用 AI 生成未来几天的饮食计划（1-7 天）
- 查询饮食计划列表（支持日期范围和状态过滤）
- 查询单个饮食计划详情
- 更新饮食计划
- 删除饮食计划
- 完成计划并自动创建餐饮记录

**数据特性**：
- 每个计划包含计划日期、餐次类型、食材列表和营养数据
- AI 会提供推荐理由（ai_reasoning），说明为什么推荐这个搭配
- 计划状态包括：pending（待执行）、completed（已完成）、skipped（已跳过）
- 完成计划时会自动创建对应的餐饮记录

---

## 接口列表

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/plans/generate` | 生成 AI 饮食计划 | 是 |
| GET | `/api/v1/plans` | 获取饮食计划列表 | 是 |
| GET | `/api/v1/plans/:id` | 获取单个饮食计划 | 是 |
| PUT | `/api/v1/plans/:id` | 更新饮食计划 | 是 |
| DELETE | `/api/v1/plans/:id` | 删除饮食计划 | 是 |
| POST | `/api/v1/plans/:id/complete` | 完成计划并创建餐饮记录 | 是 |

---

## 接口详情


### 生成 AI 饮食计划

**接口**: `POST /api/v1/plans/generate`

**说明**: 使用 AI 生成未来几天的饮食计划。系统会根据用户的食材库、营养目标和个人偏好，智能推荐合理的餐饮搭配，并为每个计划提供推荐理由。

**认证**: 是

#### 请求参数

##### 请求体

```json
{
  "days": 3,
  "preferences": "低碳水，高蛋白，避免海鲜"
}
```

| 字段 | 类型 | 必填 | 说明 | 验证规则 |
|------|------|------|------|----------|
| days | int | 是 | 生成计划的天数 | 最小 1，最大 7 |
| preferences | string | 否 | 用户偏好说明 | 最大 500 字符 |

#### 请求示例

```bash
curl -X POST http://localhost:9090/api/v1/plans/generate \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "days": 3,
    "preferences": "低碳水，高蛋白，避免海鲜"
  }'
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "plan_date": "2024-11-17T00:00:00Z",
      "meal_type": "breakfast",
      "foods": [
        {
          "food_id": 10,
          "name": "鸡蛋",
          "amount": 2,
          "unit": "个"
        },
        {
          "food_id": 15,
          "name": "全麦面包",
          "amount": 50,
          "unit": "g"
        }
      ],
      "nutrition": {
        "protein": 18.5,
        "carbs": 25.3,
        "fat": 8.2,
        "fiber": 3.5,
        "calories": 245.0
      },
      "status": "pending",
      "ai_reasoning": "早餐选择鸡蛋和全麦面包，提供优质蛋白质和复合碳水化合物，符合低碳水高蛋白的要求，能够提供持久的能量。",
      "created_at": "2024-11-16T15:30:00Z",
      "updated_at": "2024-11-16T15:30:00Z"
    },
    {
      "id": 2,
      "user_id": 1,
      "plan_date": "2024-11-17T00:00:00Z",
      "meal_type": "lunch",
      "foods": [
        {
          "food_id": 1,
          "name": "鸡胸肉",
          "amount": 200,
          "unit": "g"
        },
        {
          "food_id": 2,
          "name": "西兰花",
          "amount": 150,
          "unit": "g"
        },
        {
          "food_id": 20,
          "name": "糙米饭",
          "amount": 80,
          "unit": "g"
        }
      ],
      "nutrition": {
        "protein": 52.3,
        "carbs": 35.8,
        "fat": 4.5,
        "fiber": 6.2,
        "calories": 385.0
      },
      "status": "pending",
      "ai_reasoning": "午餐以鸡胸肉为主要蛋白质来源，搭配西兰花和少量糙米饭，营养均衡且符合低碳水高蛋白的饮食偏好。",
      "created_at": "2024-11-16T15:30:00Z",
      "updated_at": "2024-11-16T15:30:00Z"
    },
    {
      "id": 3,
      "user_id": 1,
      "plan_date": "2024-11-17T00:00:00Z",
      "meal_type": "dinner",
      "foods": [
        {
          "food_id": 5,
          "name": "牛肉",
          "amount": 150,
          "unit": "g"
        },
        {
          "food_id": 8,
          "name": "菠菜",
          "amount": 200,
          "unit": "g"
        }
      ],
      "nutrition": {
        "protein": 45.2,
        "carbs": 8.5,
        "fat": 12.3,
        "fiber": 4.8,
        "calories": 325.0
      },
      "status": "pending",
      "ai_reasoning": "晚餐选择牛肉和菠菜，提供丰富的蛋白质和铁质，碳水化合物含量低，适合晚餐食用。",
      "created_at": "2024-11-16T15:30:00Z",
      "updated_at": "2024-11-16T15:30:00Z"
    }
  ],
  "timestamp": 1699999999
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 计划唯一标识符 |
| user_id | int64 | 所属用户 ID |
| plan_date | string | 计划日期（ISO 8601 格式） |
| meal_type | string | 餐次类型（breakfast, lunch, dinner, snack） |
| foods | array | 食材列表 |
| nutrition | object | 营养数据（自动计算） |
| nutrition.protein | number | 蛋白质总量（克） |
| nutrition.carbs | number | 碳水化合物总量（克） |
| nutrition.fat | number | 脂肪总量（克） |
| nutrition.fiber | number | 纤维总量（克） |
| nutrition.calories | number | 热量总量（千卡） |
| status | string | 计划状态（pending, completed, skipped） |
| ai_reasoning | string | AI 推荐理由 |
| created_at | string | 创建时间（ISO 8601 格式） |
| updated_at | string | 更新时间（ISO 8601 格式） |

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid request parameters: Key: 'GeneratePlanRequest.Days' Error:Field validation for 'Days' failed on the 'max' tag",
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
  "error": "failed to generate plans",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | days 超出范围（1-7）、preferences 过长 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | AI 服务调用失败、数据库错误 |

#### 注意事项

1. **生成范围**：一次最多生成 7 天的计划，建议生成 2-3 天
2. **AI 智能**：AI 会根据用户的食材库、营养目标和偏好生成计划
3. **偏好说明**：preferences 字段可以描述饮食偏好、过敏信息、口味要求等
4. **计划数量**：每天通常生成 3-4 个计划（早餐、午餐、晚餐，可能包含零食）
5. **营养计算**：系统会自动计算每个计划的营养数据
6. **推荐理由**：ai_reasoning 字段说明了 AI 为什么推荐这个搭配
7. **默认状态**：生成的计划默认状态为 pending（待执行）
8. **食材来源**：AI 只会使用用户食材库中的食材生成计划

---

### 获取饮食计划列表

**接口**: `GET /api/v1/plans`

**说明**: 获取用户的饮食计划列表，支持按日期范围和状态过滤，支持分页查询。

**认证**: 是

#### 请求参数

##### 查询参数

| 参数 | 类型 | 必填 | 说明 | 默认值 | 示例 |
|------|------|------|------|--------|------|
| start_date | string | 否 | 开始日期 | - | 2024-11-17 |
| end_date | string | 否 | 结束日期 | - | 2024-11-23 |
| status | string | 否 | 状态过滤 | - | pending, completed, skipped |
| page | int | 否 | 页码（从 1 开始） | 1 | 1, 2, 3 |
| page_size | int | 否 | 每页数据量 | 20 | 10, 20, 50（最大 100） |

#### 请求示例

```bash
# 获取所有饮食计划（默认分页）
curl -X GET "http://localhost:9090/api/v1/plans" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取指定日期范围的计划
curl -X GET "http://localhost:9090/api/v1/plans?start_date=2024-11-17&end_date=2024-11-23" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取待执行的计划
curl -X GET "http://localhost:9090/api/v1/plans?status=pending" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取第 2 页，每页 50 条
curl -X GET "http://localhost:9090/api/v1/plans?page=2&page_size=50" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 组合过滤：获取下周待执行的计划
curl -X GET "http://localhost:9090/api/v1/plans?start_date=2024-11-17&end_date=2024-11-23&status=pending" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "plan_date": "2024-11-17T00:00:00Z",
      "meal_type": "breakfast",
      "foods": [
        {
          "food_id": 10,
          "name": "鸡蛋",
          "amount": 2,
          "unit": "个"
        },
        {
          "food_id": 15,
          "name": "全麦面包",
          "amount": 50,
          "unit": "g"
        }
      ],
      "nutrition": {
        "protein": 18.5,
        "carbs": 25.3,
        "fat": 8.2,
        "fiber": 3.5,
        "calories": 245.0
      },
      "status": "pending",
      "ai_reasoning": "早餐选择鸡蛋和全麦面包，提供优质蛋白质和复合碳水化合物。",
      "created_at": "2024-11-16T15:30:00Z",
      "updated_at": "2024-11-16T15:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 15,
    "total_pages": 1
  },
  "timestamp": 1699999999
}
```

**空结果响应 (200)**:

```json
{
  "code": 0,
  "message": "success",
  "data": [],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 0,
    "total_pages": 0
  },
  "timestamp": 1699999999
}
```

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid start_date format",
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

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | 日期格式不正确、status 值无效 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | 服务器内部错误 |

#### 注意事项

1. **日期格式**：start_date 和 end_date 使用 YYYY-MM-DD 格式
2. **日期范围**：可以只指定 start_date 或 end_date，也可以同时指定
3. **状态过滤**：status 可选值为 pending、completed、skipped
4. **分页默认值**：不指定分页参数时，默认返回第 1 页，每页 20 条数据
5. **分页限制**：page_size 最大值为 100，超过会自动调整为 100
6. **过滤组合**：可以同时使用多个过滤条件
7. **空结果**：当没有符合条件的数据时，返回空数组，不会报错

---

### 获取单个饮食计划

**接口**: `GET /api/v1/plans/:id`

**说明**: 根据计划 ID 获取单个饮食计划的详细信息。

**认证**: 是

#### 请求参数

##### 路径参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | int64 | 是 | 计划 ID | 1 |

#### 请求示例

```bash
curl -X GET http://localhost:9090/api/v1/plans/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
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
    "plan_date": "2024-11-17T00:00:00Z",
    "meal_type": "breakfast",
    "foods": [
      {
        "food_id": 10,
        "name": "鸡蛋",
        "amount": 2,
        "unit": "个"
      },
      {
        "food_id": 15,
        "name": "全麦面包",
        "amount": 50,
        "unit": "g"
      }
    ],
    "nutrition": {
      "protein": 18.5,
      "carbs": 25.3,
      "fat": 8.2,
      "fiber": 3.5,
      "calories": 245.0
    },
    "status": "pending",
    "ai_reasoning": "早餐选择鸡蛋和全麦面包，提供优质蛋白质和复合碳水化合物。",
    "created_at": "2024-11-16T15:30:00Z",
    "updated_at": "2024-11-16T15:30:00Z"
  },
  "timestamp": 1699999999
}
```

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid plan id",
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

**错误响应 (404)**:

```json
{
  "code": 40401,
  "message": "resource not found",
  "error": "plan not found",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | 计划 ID 格式不正确 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 40401 | 资源不存在 | 计划不存在或不属于当前用户 |
| 50001 | 内部错误 | 服务器内部错误 |

#### 注意事项

1. **权限验证**：只能查询属于当前用户的饮食计划
2. **ID 格式**：ID 必须是有效的整数
3. **不存在处理**：如果计划不存在或不属于当前用户，返回 404 错误

---

### 更新饮食计划

**接口**: `PUT /api/v1/plans/:id`

**说明**: 更新指定饮食计划的信息。需要提供完整的计划数据，系统会重新计算营养数据。

**认证**: 是

#### 请求参数

##### 路径参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | int64 | 是 | 计划 ID | 1 |

##### 请求体

```json
{
  "plan_date": "2024-11-17T00:00:00Z",
  "meal_type": "breakfast",
  "foods": [
    {
      "food_id": 10,
      "name": "鸡蛋",
      "amount": 3,
      "unit": "个"
    },
    {
      "food_id": 15,
      "name": "全麦面包",
      "amount": 60,
      "unit": "g"
    },
    {
      "food_id": 25,
      "name": "牛奶",
      "amount": 200,
      "unit": "ml"
    }
  ],
  "status": "pending",
  "ai_reasoning": "增加了牛奶，提供更多钙质和蛋白质。"
}
```

| 字段 | 类型 | 必填 | 说明 | 验证规则 |
|------|------|------|------|----------|
| plan_date | string | 是 | 计划日期时间 | ISO 8601 格式 |
| meal_type | string | 是 | 餐次类型 | 枚举值：breakfast, lunch, dinner, snack |
| foods | array | 是 | 食材列表 | 最少 1 项，最多 50 项 |
| foods[].food_id | int64 | 是 | 食材 ID | 必须 > 0 |
| foods[].name | string | 否 | 食材名称 | 长度 0-100 字符 |
| foods[].amount | number | 是 | 食材用量 | > 0，≤ 10000 |
| foods[].unit | string | 是 | 用量单位 | 长度 1-20 字符 |
| status | string | 否 | 计划状态 | 枚举值：pending, completed, skipped |
| ai_reasoning | string | 否 | 推荐理由 | 最大 1000 字符 |

#### 请求示例

```bash
curl -X PUT http://localhost:9090/api/v1/plans/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "plan_date": "2024-11-17T00:00:00Z",
    "meal_type": "breakfast",
    "foods": [
      {
        "food_id": 10,
        "name": "鸡蛋",
        "amount": 3,
        "unit": "个"
      },
      {
        "food_id": 15,
        "name": "全麦面包",
        "amount": 60,
        "unit": "g"
      },
      {
        "food_id": 25,
        "name": "牛奶",
        "amount": 200,
        "unit": "ml"
      }
    ],
    "status": "pending",
    "ai_reasoning": "增加了牛奶，提供更多钙质和蛋白质。"
  }'
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "plan updated successfully",
  "data": null,
  "timestamp": 1699999999
}
```

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid request parameters: Key: 'UpdatePlanRequest.Foods' Error:Field validation for 'Foods' failed on the 'min' tag",
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

**错误响应 (404)**:

```json
{
  "code": 40401,
  "message": "resource not found",
  "error": "plan not found",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | 缺少必填字段、参数类型不匹配、参数值超出范围 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 40401 | 资源不存在 | 计划不存在或不属于当前用户 |
| 50001 | 内部错误 | 服务器内部错误、营养计算失败 |

#### 注意事项

1. **完整更新**：需要提供所有必填字段，不支持部分更新
2. **权限验证**：只能更新属于当前用户的饮食计划
3. **营养重算**：更新后系统会重新计算营养数据
4. **ID 不可变**：计划 ID 和用户 ID 不会被更新
5. **时间戳自动更新**：updated_at 字段会自动更新为当前时间
6. **状态修改**：可以通过此接口修改计划状态

---

### 删除饮食计划

**接口**: `DELETE /api/v1/plans/:id`

**说明**: 删除指定的饮食计划。删除操作是物理删除，数据将无法恢复。

**认证**: 是

#### 请求参数

##### 路径参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | int64 | 是 | 计划 ID | 1 |

#### 请求示例

```bash
curl -X DELETE http://localhost:9090/api/v1/plans/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "plan deleted successfully",
  "data": null,
  "timestamp": 1699999999
}
```

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid plan id",
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

**错误响应 (404)**:

```json
{
  "code": 40401,
  "message": "resource not found",
  "error": "plan not found",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | 计划 ID 格式不正确 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 40401 | 资源不存在 | 计划不存在或不属于当前用户 |
| 50001 | 内部错误 | 服务器内部错误 |

#### 注意事项

1. **物理删除**：删除操作是永久性的，数据无法恢复
2. **权限验证**：只能删除属于当前用户的饮食计划
3. **谨慎操作**：建议在删除前向用户确认
4. **已完成计划**：即使计划已完成，也可以删除（不会影响已创建的餐饮记录）

---

### 完成计划并创建餐饮记录

**接口**: `POST /api/v1/plans/:id/complete`

**说明**: 标记计划为已完成，并自动创建对应的餐饮记录。这是将计划转换为实际饮食记录的便捷方式。

**认证**: 是

#### 请求参数

##### 路径参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | int64 | 是 | 计划 ID | 1 |

#### 请求示例

```bash
curl -X POST http://localhost:9090/api/v1/plans/1/complete \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "plan completed and meal record created",
  "data": {
    "id": 15,
    "user_id": 1,
    "meal_date": "2024-11-17T00:00:00Z",
    "meal_type": "breakfast",
    "foods": [
      {
        "food_id": 10,
        "name": "鸡蛋",
        "amount": 2,
        "unit": "个"
      },
      {
        "food_id": 15,
        "name": "全麦面包",
        "amount": 50,
        "unit": "g"
      }
    ],
    "nutrition": {
      "protein": 18.5,
      "carbs": 25.3,
      "fat": 8.2,
      "fiber": 3.5,
      "calories": 245.0
    },
    "notes": "Completed from plan #1",
    "created_at": "2024-11-17T08:30:00Z",
    "updated_at": "2024-11-17T08:30:00Z"
  },
  "timestamp": 1699999999
}
```

**字段说明**：

返回的是新创建的餐饮记录（Meal）对象，包含以下字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 餐饮记录唯一标识符（新创建的） |
| user_id | int64 | 所属用户 ID |
| meal_date | string | 餐次日期时间（来自计划） |
| meal_type | string | 餐次类型（来自计划） |
| foods | array | 食材列表（来自计划） |
| nutrition | object | 营养数据（来自计划） |
| notes | string | 备注（自动生成，包含计划 ID） |
| created_at | string | 创建时间（当前时间） |
| updated_at | string | 更新时间（当前时间） |

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid plan id",
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

**错误响应 (404)**:

```json
{
  "code": 40401,
  "message": "resource not found",
  "error": "plan not found",
  "timestamp": 1699999999
}
```

**错误响应 (500)**:

```json
{
  "code": 50001,
  "message": "internal error",
  "error": "failed to complete plan",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | 计划 ID 格式不正确、计划已完成 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 40401 | 资源不存在 | 计划不存在或不属于当前用户 |
| 50001 | 内部错误 | 创建餐饮记录失败、更新计划状态失败 |

#### 注意事项

1. **自动转换**：此接口会自动将计划转换为餐饮记录
2. **状态更新**：计划状态会自动更新为 completed
3. **数据复制**：餐饮记录会复制计划的所有数据（日期、餐次、食材、营养）
4. **备注生成**：餐饮记录的备注会自动生成，格式为 "Completed from plan #[计划ID]"
5. **重复完成**：已完成的计划不能再次完成，会返回错误
6. **权限验证**：只能完成属于当前用户的计划
7. **原子操作**：创建餐饮记录和更新计划状态是原子操作，要么都成功，要么都失败
8. **返回数据**：返回的是新创建的餐饮记录，不是计划本身

---

## 数据模型

### Plan 模型

完整的 Plan 数据模型定义请参考 [数据模型文档](./data-models.md#plan-饮食计划)。

**核心字段**：
- **id**: 计划唯一标识符
- **user_id**: 所属用户 ID
- **plan_date**: 计划日期时间
- **meal_type**: 餐次类型（breakfast, lunch, dinner, snack）
- **foods**: 食材列表
- **nutrition**: 营养数据（自动计算）
- **status**: 计划状态（pending, completed, skipped）
- **ai_reasoning**: AI 推荐理由
- **created_at**: 创建时间
- **updated_at**: 更新时间

### 计划状态说明

| 状态值 | 中文名称 | 说明 | 使用场景 |
|--------|---------|------|----------|
| pending | 待执行 | 计划已生成，等待执行 | 默认状态，AI 生成后的初始状态 |
| completed | 已完成 | 计划已执行，已创建餐饮记录 | 通过完成接口自动设置 |
| skipped | 已跳过 | 计划未执行，用户选择跳过 | 用户手动设置或不想执行该计划 |

---

## 使用场景

### 场景 1：生成一周的饮食计划

用户想要为下周生成饮食计划：

```bash
curl -X POST http://localhost:9090/api/v1/plans/generate \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "days": 7,
    "preferences": "低脂肪，高蛋白，适合健身增肌"
  }'
```

### 场景 2：查看今天的饮食计划

查询今天的所有待执行计划：

```bash
curl -X GET "http://localhost:9090/api/v1/plans?start_date=2024-11-17&end_date=2024-11-17&status=pending" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 场景 3：修改计划的食材

用户觉得 AI 推荐的食材不合适，想要修改：

```bash
curl -X PUT http://localhost:9090/api/v1/plans/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "plan_date": "2024-11-17T00:00:00Z",
    "meal_type": "breakfast",
    "foods": [
      {
        "food_id": 12,
        "name": "燕麦片",
        "amount": 50,
        "unit": "g"
      },
      {
        "food_id": 25,
        "name": "牛奶",
        "amount": 250,
        "unit": "ml"
      }
    ],
    "status": "pending",
    "ai_reasoning": "用户偏好燕麦早餐"
  }'
```

### 场景 4：执行计划并记录

用户按照计划吃完早餐后，标记计划为完成：

```bash
curl -X POST http://localhost:9090/api/v1/plans/1/complete \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 场景 5：跳过不想执行的计划

用户决定不执行某个计划，标记为跳过：

```bash
curl -X PUT http://localhost:9090/api/v1/plans/2 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "plan_date": "2024-11-17T00:00:00Z",
    "meal_type": "lunch",
    "foods": [...],
    "status": "skipped",
    "ai_reasoning": "用户外出就餐"
  }'
```

### 场景 6：删除过期的计划

删除已经过期且未执行的计划：

```bash
curl -X DELETE http://localhost:9090/api/v1/plans/5 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## 最佳实践

### 1. 计划生成

- **合理天数**：建议生成 2-3 天的计划，不要一次生成太多
- **明确偏好**：在 preferences 中清楚说明饮食偏好、过敏信息、口味要求
- **定期生成**：建议每 2-3 天生成一次新计划，保持计划的新鲜度
- **参考营养目标**：AI 会参考用户的营养目标，确保目标设置合理

### 2. 计划管理

- **及时查看**：每天查看当天的计划，提前准备食材
- **灵活调整**：如果计划不合适，及时修改或跳过
- **状态管理**：合理使用 pending、completed、skipped 状态
- **定期清理**：删除过期且未执行的计划，保持数据整洁

### 3. 计划执行

- **按计划执行**：尽量按照计划执行，有助于养成良好的饮食习惯
- **完成记录**：执行后使用完成接口，自动创建餐饮记录
- **记录偏差**：如果实际执行与计划有偏差，可以在餐饮记录中修改
- **总结反馈**：定期总结计划执行情况，调整偏好设置

### 4. AI 推荐理由

- **理解理由**：阅读 ai_reasoning 字段，了解 AI 的推荐逻辑
- **学习知识**：从推荐理由中学习营养搭配知识
- **反馈优化**：如果推荐不合理，通过偏好设置优化

### 5. 与餐饮记录的关系

- **计划先行**：先生成计划，再执行并记录
- **灵活记录**：也可以不按计划，直接创建餐饮记录
- **数据一致**：完成计划后，餐饮记录会自动创建，数据保持一致
- **独立管理**：计划和餐饮记录是独立的，删除计划不会影响已创建的餐饮记录

### 6. 查询优化

- **使用日期过滤**：查询时指定日期范围，减少数据量
- **状态过滤**：使用状态过滤快速找到待执行或已完成的计划
- **合理分页**：根据实际需求设置 page_size
- **缓存数据**：对于当天的计划，可以在客户端缓存

### 7. 错误处理

- **检查状态**：完成计划前检查状态，避免重复完成
- **验证数据**：更新计划时验证食材和用量的合理性
- **处理失败**：如果完成计划失败，检查错误信息并重试

---

## 常见问题

### Q: AI 生成计划的依据是什么？

A: 
- AI 会综合考虑以下因素：
  - 用户的食材库（只使用用户拥有的食材）
  - 用户的营养目标（热量、蛋白质、碳水化合物等）
  - 用户的饮食偏好（通过 preferences 参数传递）
  - 营养均衡原则（确保营养搭配合理）
  - 餐次特点（早餐、午餐、晚餐的不同需求）

### Q: 可以一次生成多少天的计划？

A: 
- 最少 1 天，最多 7 天
- 建议生成 2-3 天的计划
- 生成太多天可能导致食材过期或计划不够灵活

### Q: 生成的计划可以修改吗？

A: 
- 可以，使用更新接口可以修改计划的任何内容
- 修改后营养数据会自动重新计算
- 可以修改食材、用量、餐次类型、状态等

### Q: 完成计划后可以修改餐饮记录吗？

A: 
- 可以，完成计划后会创建餐饮记录
- 餐饮记录可以通过餐饮记录模块的接口修改
- 修改餐饮记录不会影响原计划

### Q: 已完成的计划可以删除吗？

A: 
- 可以删除，删除计划不会影响已创建的餐饮记录
- 餐饮记录和计划是独立的数据
- 建议保留已完成的计划，便于后续分析

### Q: 可以重复完成同一个计划吗？

A: 
- 不可以，每个计划只能完成一次
- 如果需要重复执行，可以创建新的计划或直接创建餐饮记录
- 已完成的计划状态为 completed，不能再次完成

### Q: 跳过的计划可以恢复吗？

A: 
- 可以，使用更新接口将状态改回 pending
- 或者删除跳过的计划，重新生成
- 跳过只是标记状态，不影响数据

### Q: 如何处理计划与实际执行的偏差？

A: 
- 如果偏差较小，可以直接完成计划
- 如果偏差较大，可以先修改计划，再完成
- 也可以跳过计划，直接创建餐饮记录

### Q: AI 推荐的食材不喜欢怎么办？

A: 
- 可以修改计划，替换为喜欢的食材
- 可以在 preferences 中说明不喜欢的食材
- 可以删除计划，重新生成（调整 preferences）

### Q: 计划的日期可以是过去的日期吗？

A: 
- 技术上可以，但不建议
- 计划是为未来准备的，过去的应该直接创建餐饮记录
- AI 生成的计划默认是未来的日期

### Q: 如何查看计划的执行率？

A: 
- 可以通过查询接口统计不同状态的计划数量
- 计算 completed / (completed + skipped + pending) 得到执行率
- 建议在客户端实现统计功能

### Q: 计划生成失败怎么办？

A: 
- 检查食材库是否有足够的食材
- 检查 preferences 是否过于严格
- 检查 AI 服务是否正常
- 查看错误信息，根据提示调整

---

## 相关文档

- [数据模型](./data-models.md) - 查看 Plan 模型的完整定义
- [餐饮记录模块](./03-meals.md) - 了解如何管理餐饮记录
- [AI 服务模块](./05-ai-services.md) - 了解 AI 服务的其他功能
- [食材管理模块](./02-foods.md) - 了解如何管理食材库
- [营养分析模块](./06-nutrition.md) - 了解如何分析营养数据
- [通用概念](./common-concepts.md) - 了解认证、分页等通用概念
- [错误码说明](./error-codes.md) - 查看所有错误码的详细说明
- [API 文档总览](./README.md) - 返回 API 文档首页
