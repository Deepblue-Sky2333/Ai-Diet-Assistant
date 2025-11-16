# 餐饮记录模块

## 概述

餐饮记录模块提供用户日常饮食记录的完整管理功能，包括餐饮记录的创建、查询、更新和删除。用户可以记录每日的早餐、午餐、晚餐和零食，系统会自动计算每餐的营养数据，为营养分析和饮食计划提供基础数据。

**核心功能**：
- 创建餐饮记录（支持多种食材组合）
- 查询餐饮记录列表（支持日期范围和餐次类型过滤）
- 查询单个餐饮记录详情
- 更新餐饮记录
- 删除餐饮记录

**数据特性**：
- 每条记录包含餐次日期和类型（早餐、午餐、晚餐、零食）
- 支持一餐包含多种食材，每种食材可指定用量
- 系统自动计算整餐的营养数据（蛋白质、碳水化合物、脂肪、纤维、热量）
- 支持添加备注信息

---

## 接口列表

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/meals` | 创建餐饮记录 | 是 |
| GET | `/api/v1/meals` | 获取餐饮记录列表 | 是 |
| GET | `/api/v1/meals/:id` | 获取单个餐饮记录 | 是 |
| PUT | `/api/v1/meals/:id` | 更新餐饮记录 | 是 |
| DELETE | `/api/v1/meals/:id` | 删除餐饮记录 | 是 |

---

## 接口详情

### 创建餐饮记录

**接口**: `POST /api/v1/meals`

**说明**: 创建一条新的餐饮记录。系统会根据食材和用量自动计算整餐的营养数据。

**认证**: 是

#### 请求参数

##### 请求体

```json
{
  "meal_date": "2024-11-16T12:00:00Z",
  "meal_type": "lunch",
  "foods": [
    {
      "food_id": 1,
      "name": "鸡胸肉",
      "amount": 150,
      "unit": "g"
    },
    {
      "food_id": 2,
      "name": "西兰花",
      "amount": 200,
      "unit": "g"
    }
  ],
  "notes": "健康午餐"
}
```

| 字段 | 类型 | 必填 | 说明 | 验证规则 |
|------|------|------|------|----------|
| meal_date | string | 是 | 餐次日期时间 | ISO 8601 格式 |
| meal_type | string | 是 | 餐次类型 | 枚举值：breakfast, lunch, dinner, snack |
| foods | array | 是 | 食材列表 | 最少 1 项，最多 50 项 |
| foods[].food_id | int64 | 是 | 食材 ID | 必须 > 0 |
| foods[].name | string | 否 | 食材名称 | 长度 0-100 字符（可选，用于显示） |
| foods[].amount | number | 是 | 食材用量 | > 0，≤ 10000 |
| foods[].unit | string | 是 | 用量单位 | 长度 1-20 字符 |
| notes | string | 否 | 备注 | 最大 500 字符 |

#### 请求示例

```bash
curl -X POST http://localhost:9090/api/v1/meals \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "meal_date": "2024-11-16T12:00:00Z",
    "meal_type": "lunch",
    "foods": [
      {
        "food_id": 1,
        "name": "鸡胸肉",
        "amount": 150,
        "unit": "g"
      },
      {
        "food_id": 2,
        "name": "西兰花",
        "amount": 200,
        "unit": "g"
      }
    ],
    "notes": "健康午餐"
  }'
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
    "meal_date": "2024-11-16T12:00:00Z",
    "meal_type": "lunch",
    "foods": [
      {
        "food_id": 1,
        "name": "鸡胸肉",
        "amount": 150,
        "unit": "g"
      },
      {
        "food_id": 2,
        "name": "西兰花",
        "amount": 200,
        "unit": "g"
      }
    ],
    "nutrition": {
      "protein": 40.1,
      "carbs": 13.2,
      "fat": 2.6,
      "fiber": 5.2,
      "calories": 233.0
    },
    "notes": "健康午餐",
    "created_at": "2024-11-16T12:05:00Z",
    "updated_at": "2024-11-16T12:05:00Z"
  },
  "timestamp": 1699999999
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 餐饮记录唯一标识符 |
| user_id | int64 | 所属用户 ID |
| meal_date | string | 餐次日期时间（ISO 8601 格式） |
| meal_type | string | 餐次类型 |
| foods | array | 食材列表 |
| nutrition | object | 营养数据（自动计算） |
| nutrition.protein | number | 蛋白质总量（克） |
| nutrition.carbs | number | 碳水化合物总量（克） |
| nutrition.fat | number | 脂肪总量（克） |
| nutrition.fiber | number | 纤维总量（克） |
| nutrition.calories | number | 热量总量（千卡） |
| notes | string | 备注 |
| created_at | string | 创建时间（ISO 8601 格式） |
| updated_at | string | 更新时间（ISO 8601 格式） |

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid request parameters: Key: 'CreateMealRequest.MealType' Error:Field validation for 'MealType' failed on the 'oneof' tag",
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
| 40001 | 参数错误 | 缺少必填字段、参数类型不匹配、参数值超出范围 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | 服务器内部错误、营养计算失败 |

#### 注意事项

1. **营养自动计算**：系统会根据食材的营养数据和用量自动计算整餐的营养总量
2. **食材验证**：food_id 必须是用户市场面板中存在的食材
3. **用量单位**：用量单位应与食材定义的单位一致，系统会按比例计算营养
4. **餐次类型**：meal_type 必须是以下值之一：breakfast（早餐）、lunch（午餐）、dinner（晚餐）、snack（零食）
5. **日期格式**：meal_date 使用 ISO 8601 格式，包含日期和时间

---

### 获取餐饮记录列表

**接口**: `GET /api/v1/meals`

**说明**: 获取用户的餐饮记录列表，支持按日期范围和餐次类型过滤，支持分页查询。

**认证**: 是

#### 请求参数

##### 查询参数

| 参数 | 类型 | 必填 | 说明 | 默认值 | 示例 |
|------|------|------|------|--------|------|
| start_date | string | 否 | 开始日期 | - | 2024-11-01 |
| end_date | string | 否 | 结束日期 | - | 2024-11-30 |
| meal_type | string | 否 | 餐次类型过滤 | - | breakfast, lunch, dinner, snack |
| page | int | 否 | 页码（从 1 开始） | 1 | 1, 2, 3 |
| page_size | int | 否 | 每页数据量 | 20 | 10, 20, 50（最大 100） |

#### 请求示例

```bash
# 获取所有餐饮记录（默认分页）
curl -X GET "http://localhost:9090/api/v1/meals" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取指定日期范围的记录
curl -X GET "http://localhost:9090/api/v1/meals?start_date=2024-11-01&end_date=2024-11-30" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取午餐记录
curl -X GET "http://localhost:9090/api/v1/meals?meal_type=lunch" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取第 2 页，每页 50 条
curl -X GET "http://localhost:9090/api/v1/meals?page=2&page_size=50" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 组合过滤：获取 11 月的早餐记录
curl -X GET "http://localhost:9090/api/v1/meals?start_date=2024-11-01&end_date=2024-11-30&meal_type=breakfast" \
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
      "meal_date": "2024-11-16T12:00:00Z",
      "meal_type": "lunch",
      "foods": [
        {
          "food_id": 1,
          "name": "鸡胸肉",
          "amount": 150,
          "unit": "g"
        },
        {
          "food_id": 2,
          "name": "西兰花",
          "amount": 200,
          "unit": "g"
        }
      ],
      "nutrition": {
        "protein": 40.1,
        "carbs": 13.2,
        "fat": 2.6,
        "fiber": 5.2,
        "calories": 233.0
      },
      "notes": "健康午餐",
      "created_at": "2024-11-16T12:05:00Z",
      "updated_at": "2024-11-16T12:05:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 45,
    "total_pages": 3
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
| 40001 | 参数错误 | 日期格式不正确、meal_type 值无效 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | 服务器内部错误 |

#### 注意事项

1. **日期格式**：start_date 和 end_date 使用 YYYY-MM-DD 格式
2. **日期范围**：可以只指定 start_date 或 end_date，也可以同时指定
3. **分页默认值**：不指定分页参数时，默认返回第 1 页，每页 20 条数据
4. **分页限制**：page_size 最大值为 100，超过会自动调整为 100
5. **过滤组合**：可以同时使用多个过滤条件
6. **空结果**：当没有符合条件的数据时，返回空数组，不会报错

---

### 获取单个餐饮记录

**接口**: `GET /api/v1/meals/:id`

**说明**: 根据餐饮记录 ID 获取单个餐饮记录的详细信息。

**认证**: 是

#### 请求参数

##### 路径参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | int64 | 是 | 餐饮记录 ID | 1 |

#### 请求示例

```bash
curl -X GET http://localhost:9090/api/v1/meals/1 \
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
    "meal_date": "2024-11-16T12:00:00Z",
    "meal_type": "lunch",
    "foods": [
      {
        "food_id": 1,
        "name": "鸡胸肉",
        "amount": 150,
        "unit": "g"
      },
      {
        "food_id": 2,
        "name": "西兰花",
        "amount": 200,
        "unit": "g"
      }
    ],
    "nutrition": {
      "protein": 40.1,
      "carbs": 13.2,
      "fat": 2.6,
      "fiber": 5.2,
      "calories": 233.0
    },
    "notes": "健康午餐",
    "created_at": "2024-11-16T12:05:00Z",
    "updated_at": "2024-11-16T12:05:00Z"
  },
  "timestamp": 1699999999
}
```

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid meal id",
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
  "error": "meal not found",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | 餐饮记录 ID 格式不正确 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 40401 | 资源不存在 | 餐饮记录不存在或不属于当前用户 |
| 50001 | 内部错误 | 服务器内部错误 |

#### 注意事项

1. **权限验证**：只能查询属于当前用户的餐饮记录
2. **ID 格式**：ID 必须是有效的整数
3. **不存在处理**：如果餐饮记录不存在或不属于当前用户，返回 404 错误

---

### 更新餐饮记录

**接口**: `PUT /api/v1/meals/:id`

**说明**: 更新指定餐饮记录的信息。需要提供完整的餐饮数据，系统会重新计算营养数据。

**认证**: 是

#### 请求参数

##### 路径参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | int64 | 是 | 餐饮记录 ID | 1 |

##### 请求体

```json
{
  "meal_date": "2024-11-16T12:00:00Z",
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
      "food_id": 3,
      "name": "糙米饭",
      "amount": 100,
      "unit": "g"
    }
  ],
  "notes": "增加了主食"
}
```

| 字段 | 类型 | 必填 | 说明 | 验证规则 |
|------|------|------|------|----------|
| meal_date | string | 是 | 餐次日期时间 | ISO 8601 格式 |
| meal_type | string | 是 | 餐次类型 | 枚举值：breakfast, lunch, dinner, snack |
| foods | array | 是 | 食材列表 | 最少 1 项，最多 50 项 |
| foods[].food_id | int64 | 是 | 食材 ID | 必须 > 0 |
| foods[].name | string | 否 | 食材名称 | 长度 0-100 字符 |
| foods[].amount | number | 是 | 食材用量 | > 0，≤ 10000 |
| foods[].unit | string | 是 | 用量单位 | 长度 1-20 字符 |
| notes | string | 否 | 备注 | 最大 500 字符 |

#### 请求示例

```bash
curl -X PUT http://localhost:9090/api/v1/meals/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "meal_date": "2024-11-16T12:00:00Z",
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
        "food_id": 3,
        "name": "糙米饭",
        "amount": 100,
        "unit": "g"
      }
    ],
    "notes": "增加了主食"
  }'
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "meal updated successfully",
  "data": null,
  "timestamp": 1699999999
}
```

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid request parameters: Key: 'UpdateMealRequest.Foods' Error:Field validation for 'Foods' failed on the 'min' tag",
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
  "error": "meal not found",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | 缺少必填字段、参数类型不匹配、参数值超出范围 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 40401 | 资源不存在 | 餐饮记录不存在或不属于当前用户 |
| 50001 | 内部错误 | 服务器内部错误、营养计算失败 |

#### 注意事项

1. **完整更新**：需要提供所有必填字段，不支持部分更新
2. **权限验证**：只能更新属于当前用户的餐饮记录
3. **营养重算**：更新后系统会重新计算营养数据
4. **ID 不可变**：餐饮记录 ID 和用户 ID 不会被更新
5. **时间戳自动更新**：updated_at 字段会自动更新为当前时间

---

### 删除餐饮记录

**接口**: `DELETE /api/v1/meals/:id`

**说明**: 删除指定的餐饮记录。删除操作是物理删除，数据将无法恢复。

**认证**: 是

#### 请求参数

##### 路径参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | int64 | 是 | 餐饮记录 ID | 1 |

#### 请求示例

```bash
curl -X DELETE http://localhost:9090/api/v1/meals/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "meal deleted successfully",
  "data": null,
  "timestamp": 1699999999
}
```

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid meal id",
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
  "error": "meal not found",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | 餐饮记录 ID 格式不正确 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 40401 | 资源不存在 | 餐饮记录不存在或不属于当前用户 |
| 50001 | 内部错误 | 服务器内部错误 |

#### 注意事项

1. **物理删除**：删除操作是永久性的，数据无法恢复
2. **权限验证**：只能删除属于当前用户的餐饮记录
3. **营养统计影响**：删除记录会影响相关日期的营养统计数据
4. **谨慎操作**：建议在删除前向用户确认

---

## 数据模型

### Meal 模型

完整的 Meal 数据模型定义请参考 [数据模型文档](./data-models.md#meal-餐饮记录)。

**核心字段**：
- **id**: 餐饮记录唯一标识符
- **user_id**: 所属用户 ID
- **meal_date**: 餐次日期时间
- **meal_type**: 餐次类型（breakfast, lunch, dinner, snack）
- **foods**: 食材列表
- **nutrition**: 营养数据（自动计算）
- **notes**: 备注
- **created_at**: 创建时间
- **updated_at**: 更新时间

### MealFood 模型

**字段说明**：
- **food_id**: 食材 ID（引用 Food 表）
- **name**: 食材名称（可选，用于显示）
- **amount**: 食材用量
- **unit**: 用量单位

### 餐次类型说明

| 类型值 | 中文名称 | 说明 | 典型时间 |
|--------|---------|------|----------|
| breakfast | 早餐 | 早晨的第一餐 | 06:00 - 10:00 |
| lunch | 午餐 | 中午的正餐 | 11:00 - 14:00 |
| dinner | 晚餐 | 晚上的正餐 | 17:00 - 21:00 |
| snack | 零食 | 正餐之间的加餐 | 任意时间 |

---

## 使用场景

### 场景 1：记录早餐

用户在早餐后记录今天的早餐内容：

```bash
curl -X POST http://localhost:9090/api/v1/meals \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "meal_date": "2024-11-16T08:00:00Z",
    "meal_type": "breakfast",
    "foods": [
      {
        "food_id": 10,
        "name": "燕麦片",
        "amount": 50,
        "unit": "g"
      },
      {
        "food_id": 11,
        "name": "牛奶",
        "amount": 250,
        "unit": "ml"
      },
      {
        "food_id": 12,
        "name": "香蕉",
        "amount": 1,
        "unit": "个"
      }
    ],
    "notes": "健康早餐"
  }'
```

### 场景 2：查询本周的饮食记录

查询本周的所有餐饮记录，用于周总结：

```bash
curl -X GET "http://localhost:9090/api/v1/meals?start_date=2024-11-10&end_date=2024-11-16" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 场景 3：修改记录的食材用量

用户发现记录的食材用量不准确，需要修改：

```bash
curl -X PUT http://localhost:9090/api/v1/meals/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "meal_date": "2024-11-16T12:00:00Z",
    "meal_type": "lunch",
    "foods": [
      {
        "food_id": 1,
        "name": "鸡胸肉",
        "amount": 180,
        "unit": "g"
      },
      {
        "food_id": 2,
        "name": "西兰花",
        "amount": 200,
        "unit": "g"
      }
    ],
    "notes": "修正了鸡胸肉的用量"
  }'
```

### 场景 4：查询今天的午餐记录

查询今天是否已经记录了午餐：

```bash
curl -X GET "http://localhost:9090/api/v1/meals?start_date=2024-11-16&end_date=2024-11-16&meal_type=lunch" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## 最佳实践

### 1. 记录时机

- **及时记录**：建议在用餐后立即记录，避免遗忘
- **准确用量**：使用厨房秤等工具准确测量食材用量
- **完整记录**：记录所有食材，包括调味料（如果有营养价值）

### 2. 食材选择

- **使用市场面板**：从用户的市场面板中选择食材，确保营养数据准确
- **统一单位**：使用与食材定义一致的单位，便于营养计算
- **添加名称**：在 foods 中添加 name 字段，便于显示和理解

### 3. 餐次分类

- **合理分类**：根据用餐时间和性质选择合适的餐次类型
- **零食记录**：不要忽略零食和加餐，它们也会影响营养摄入
- **时间准确**：meal_date 应该反映实际用餐时间

### 4. 备注使用

- **记录特殊情况**：如外出就餐、特殊烹饪方式等
- **记录感受**：如饱腹感、口味等，便于后续分析
- **简洁明了**：备注应简洁，避免过长

### 5. 数据修正

- **及时修正**：发现错误及时修正，保证数据准确性
- **完整更新**：更新时提供完整数据，避免遗漏
- **谨慎删除**：删除前确认，避免误删重要数据

### 6. 查询优化

- **使用日期过滤**：查询时指定日期范围，减少数据量
- **合理分页**：根据实际需求设置 page_size
- **缓存数据**：对于当天的数据，可以在客户端缓存

### 7. 营养计算

- **信任自动计算**：系统会自动计算营养数据，无需手动计算
- **检查异常值**：如果营养数据异常，检查食材数据和用量是否正确
- **定期校验**：定期检查食材库的营养数据是否准确

---

## 常见问题

### Q: 如何计算整餐的营养数据？

A: 
- 系统会自动计算，无需手动计算
- 计算方式：根据每种食材的营养数据和用量，按比例计算后求和
- 例如：鸡胸肉 100g 含蛋白质 23g，使用 150g 则蛋白质为 23 * 1.5 = 34.5g

### Q: 可以修改餐次的日期吗？

A: 
- 可以，使用更新接口可以修改 meal_date
- 修改日期会影响营养统计，相关日期的统计数据会自动更新
- 建议谨慎修改，确保日期准确

### Q: 删除餐饮记录会影响营养统计吗？

A: 
- 会影响，删除记录后相关日期的营养统计会自动更新
- 删除是永久性的，无法恢复
- 如果不确定，建议先备份数据

### Q: 一餐可以包含多少种食材？

A: 
- 最少 1 种，最多 50 种
- 建议记录主要食材即可，不必过于详细
- 如果食材种类过多，可以合并同类食材

### Q: 如何处理外出就餐的记录？

A: 
- 尽量估算食材和用量，记录主要成分
- 在备注中说明是外出就餐，便于后续分析
- 可以参考餐厅提供的营养信息
- 如果无法准确估算，可以记录大致情况

### Q: 食材用量的单位必须与食材定义一致吗？

A: 
- 建议一致，便于营养计算
- 系统会按比例计算，但单位不一致可能导致计算错误
- 例如：食材定义单位是 "100g"，记录时也应使用 "g" 作为单位

### Q: 可以记录昨天或更早的餐饮吗？

A: 
- 可以，meal_date 可以是任意日期
- 建议及时记录，避免遗忘
- 补录历史数据时，确保日期和内容准确

### Q: 如何查询某个月的所有餐饮记录？

A: 
- 使用 start_date 和 end_date 参数指定月份范围
- 例如：start_date=2024-11-01&end_date=2024-11-30
- 注意处理分页，可能需要多次请求

### Q: 营养数据为什么与预期不符？

A: 
- 检查食材库中的营养数据是否准确
- 检查食材用量和单位是否正确
- 检查食材的单位定义（如 "100g" vs "1个"）
- 如果食材数据有误，更新食材后重新计算

### Q: 可以批量创建餐饮记录吗？

A: 
- 当前版本不支持批量创建
- 需要逐条创建餐饮记录
- 建议在客户端实现批量创建功能，循环调用创建接口

---

## 相关文档

- [数据模型](./data-models.md) - 查看 Meal 模型的完整定义
- [食材管理模块](./02-foods.md) - 了解如何管理食材库
- [营养分析模块](./06-nutrition.md) - 了解如何分析营养数据
- [Dashboard 模块](./07-dashboard.md) - 了解如何查看饮食概览
- [通用概念](./common-concepts.md) - 了解认证、分页等通用概念
- [错误码说明](./error-codes.md) - 查看所有错误码的详细说明
- [API 文档总览](./README.md) - 返回 API 文档首页
