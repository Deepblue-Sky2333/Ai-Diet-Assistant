# AI Diet Assistant API 接口文档

## 基础信息

- **Base URL**: `http://localhost:9090/api/v1`（开发环境）
- **认证方式**: JWT Bearer Token
- **请求格式**: `application/json`
- **响应格式**: `application/json`
- **字符编码**: UTF-8

## 认证说明

### JWT Token 机制

系统使用双 Token 认证机制：
- **Access Token**: 有效期 15 分钟，用于 API 请求认证
- **Refresh Token**: 有效期 7 天，用于刷新 Access Token

### 使用方式

在所有需要认证的 API 请求中，添加 Authorization Header：

```
Authorization: Bearer {access_token}
```

### Token 刷新流程

1. Access Token 过期时，API 返回 401 错误
2. 使用 Refresh Token 调用 `/auth/refresh` 获取新的 Access Token
3. 如果 Refresh Token 也过期，需要重新登录

## 响应格式

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
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
  }
}
```

### 错误响应

```json
{
  "code": 40001,
  "message": "invalid request parameters",
  "error": "详细错误信息"
}
```


## 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 40001 | 参数错误 |
| 40101 | 未认证 |
| 40301 | 禁止访问 |
| 40401 | 资源不存在 |
| 42901 | 请求过于频繁 |
| 50001 | 服务器内部错误 |

---

## 1. 认证接口

### 1.1 用户登录

**接口**: `POST /auth/login`

**说明**: 用户使用用户名和密码登录系统

**请求参数**:

```json
{
  "username": "user",
  "password": "password123"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名，3-50字符，仅字母数字 |
| password | string | 是 | 密码，8-128字符 |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900
  }
}
```

**错误码**:
- `40001`: 参数错误
- `40101`: 用户名或密码错误
- `42901`: 账户被锁定（登录失败次数过多）

---

### 1.2 刷新 Token

**接口**: `POST /auth/refresh`

**说明**: 使用 Refresh Token 获取新的 Access Token

**请求参数**:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

### 1.3 用户登出

**接口**: `POST /auth/logout`

**说明**: 用户登出系统，将当前 Token 加入黑名单

**认证**: 需要

**响应示例**:

```json
{
  "code": 0,
  "message": "logout successful"
}
```

---

### 1.4 修改密码

**接口**: `PUT /auth/password`

**说明**: 用户修改登录密码

**认证**: 需要

**请求参数**:

```json
{
  "old_password": "oldpassword123",
  "new_password": "newpassword123"
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "password changed successfully, please login again"
}
```

---


## 2. 食材管理接口

### 2.1 创建食材

**接口**: `POST /foods`

**说明**: 在用户的食材库中创建新食材

**认证**: 需要

**请求参数**:

```json
{
  "name": "鸡胸肉",
  "category": "meat",
  "price": 15.5,
  "unit": "100g",
  "protein": 23.0,
  "carbs": 0.0,
  "fat": 1.2,
  "fiber": 0.0,
  "calories": 105.0,
  "available": true
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 食材名称，1-100字符 |
| category | string | 是 | 分类：meat, vegetable, fruit, grain, other |
| price | number | 是 | 价格，≥0 |
| unit | string | 是 | 单位，1-20字符 |
| protein | number | 是 | 蛋白质（克），≥0 |
| carbs | number | 是 | 碳水化合物（克），≥0 |
| fat | number | 是 | 脂肪（克），≥0 |
| fiber | number | 是 | 纤维（克），≥0 |
| calories | number | 是 | 卡路里，≥0 |
| available | boolean | 否 | 是否可用，默认 true |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "鸡胸肉",
    "category": "meat",
    "price": 15.5,
    "unit": "100g",
    "protein": 23.0,
    "carbs": 0.0,
    "fat": 1.2,
    "fiber": 0.0,
    "calories": 105.0,
    "available": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

---

### 2.2 获取食材列表

**接口**: `GET /foods`

**说明**: 获取用户的食材列表，支持分页和筛选

**认证**: 需要

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| category | string | 否 | 按分类筛选：meat, vegetable, fruit, grain, other |
| available | boolean | 否 | 按可用性筛选：true, false |
| page | integer | 否 | 页码，默认 1 |
| page_size | integer | 否 | 每页数量，默认 20，最大 100 |

**请求示例**:

```
GET /foods?category=meat&available=true&page=1&page_size=20
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "name": "鸡胸肉",
      "category": "meat",
      "price": 15.5,
      "unit": "100g",
      "protein": 23.0,
      "carbs": 0.0,
      "fat": 1.2,
      "fiber": 0.0,
      "calories": 105.0,
      "available": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 50,
    "total_pages": 3
  }
}
```

---

### 2.3 获取单个食材

**接口**: `GET /foods/:id`

**说明**: 获取指定 ID 的食材详情

**认证**: 需要

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 食材 ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "鸡胸肉",
    "category": "meat",
    "price": 15.5,
    "unit": "100g",
    "protein": 23.0,
    "carbs": 0.0,
    "fat": 1.2,
    "fiber": 0.0,
    "calories": 105.0,
    "available": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

---

### 2.4 更新食材

**接口**: `PUT /foods/:id`

**说明**: 更新指定 ID 的食材信息

**认证**: 需要

**请求参数**: 同创建食材

**响应示例**:

```json
{
  "code": 0,
  "message": "food updated successfully"
}
```

---

### 2.5 删除食材

**接口**: `DELETE /foods/:id`

**说明**: 删除指定 ID 的食材

**认证**: 需要

**响应示例**:

```json
{
  "code": 0,
  "message": "food deleted successfully"
}
```

---

### 2.6 批量导入食材

**接口**: `POST /foods/batch`

**说明**: 批量导入多个食材

**认证**: 需要

**请求参数**:

```json
{
  "foods": [
    {
      "name": "鸡胸肉",
      "category": "meat",
      "price": 15.5,
      "unit": "100g",
      "protein": 23.0,
      "carbs": 0.0,
      "fat": 1.2,
      "fiber": 0.0,
      "calories": 105.0,
      "available": true
    },
    {
      "name": "西兰花",
      "category": "vegetable",
      "price": 8.0,
      "unit": "100g",
      "protein": 2.8,
      "carbs": 6.6,
      "fat": 0.4,
      "fiber": 2.6,
      "calories": 34.0,
      "available": true
    }
  ]
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "success_count": 2,
    "failed_count": 0,
    "total_count": 2
  }
}
```

---


## 3. 餐饮记录接口

### 3.1 创建餐饮记录

**接口**: `POST /meals`

**说明**: 创建新的餐饮记录，系统自动计算营养摄入

**认证**: 需要

**请求参数**:

```json
{
  "meal_date": "2024-01-15T12:00:00Z",
  "meal_type": "lunch",
  "foods": [
    {
      "food_id": 1,
      "quantity": 2.0
    },
    {
      "food_id": 2,
      "quantity": 1.5
    }
  ],
  "notes": "午餐很美味"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| meal_date | string | 是 | 餐饮日期时间（ISO 8601 格式） |
| meal_type | string | 是 | 餐型：breakfast, lunch, dinner, snack |
| foods | array | 是 | 食物列表，1-50 项 |
| foods[].food_id | integer | 是 | 食材 ID |
| foods[].quantity | number | 是 | 份量（倍数） |
| notes | string | 否 | 备注，最多 500 字符 |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "user_id": 1,
    "meal_date": "2024-01-15T12:00:00Z",
    "meal_type": "lunch",
    "foods": [
      {
        "food_id": 1,
        "food_name": "鸡胸肉",
        "quantity": 2.0,
        "unit": "100g"
      }
    ],
    "nutrition": {
      "protein": 46.0,
      "carbs": 0.0,
      "fat": 2.4,
      "fiber": 0.0,
      "calories": 210.0
    },
    "notes": "午餐很美味",
    "created_at": "2024-01-15T12:00:00Z",
    "updated_at": "2024-01-15T12:00:00Z"
  }
}
```

---

### 3.2 获取餐饮记录列表

**接口**: `GET /meals`

**说明**: 获取餐饮记录列表，支持日期筛选和分页

**认证**: 需要

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_date | string | 否 | 开始日期（YYYY-MM-DD） |
| end_date | string | 否 | 结束日期（YYYY-MM-DD） |
| meal_type | string | 否 | 餐型筛选：breakfast, lunch, dinner, snack |
| page | integer | 否 | 页码，默认 1 |
| page_size | integer | 否 | 每页数量，默认 20，最大 100 |

**请求示例**:

```
GET /meals?start_date=2024-01-01&end_date=2024-01-31&meal_type=lunch&page=1
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "meal_date": "2024-01-15T12:00:00Z",
      "meal_type": "lunch",
      "foods": [...],
      "nutrition": {...},
      "notes": "午餐很美味",
      "created_at": "2024-01-15T12:00:00Z",
      "updated_at": "2024-01-15T12:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 30,
    "total_pages": 2
  }
}
```

---

### 3.3 获取单个餐饮记录

**接口**: `GET /meals/:id`

**说明**: 获取指定 ID 的餐饮记录详情

**认证**: 需要

**响应示例**: 同创建餐饮记录响应

---

### 3.4 更新餐饮记录

**接口**: `PUT /meals/:id`

**说明**: 更新指定 ID 的餐饮记录

**认证**: 需要

**请求参数**: 同创建餐饮记录

**响应示例**:

```json
{
  "code": 0,
  "message": "meal updated successfully"
}
```

---

### 3.5 删除餐饮记录

**接口**: `DELETE /meals/:id`

**说明**: 删除指定 ID 的餐饮记录

**认证**: 需要

**响应示例**:

```json
{
  "code": 0,
  "message": "meal deleted successfully"
}
```

---


## 4. 饮食计划接口

### 4.1 生成 AI 饮食计划

**接口**: `POST /plans/generate`

**说明**: 使用 AI 生成未来几天的饮食计划

**认证**: 需要

**请求参数**:

```json
{
  "days": 2,
  "preferences": "我喜欢高蛋白低碳水的饮食，不吃辣"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| days | integer | 是 | 生成天数，1-7 天 |
| preferences | string | 否 | 用户偏好，最多 500 字符 |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "plan_date": "2024-01-16T08:00:00Z",
      "meal_type": "breakfast",
      "foods": [
        {
          "food_id": 1,
          "food_name": "鸡胸肉",
          "quantity": 2.0,
          "unit": "100g"
        }
      ],
      "nutrition": {
        "protein": 46.0,
        "carbs": 0.0,
        "fat": 2.4,
        "fiber": 0.0,
        "calories": 210.0
      },
      "status": "pending",
      "ai_reasoning": "高蛋白低碳水，符合您的偏好",
      "created_at": "2024-01-15T12:00:00Z",
      "updated_at": "2024-01-15T12:00:00Z"
    }
  ]
}
```

---

### 4.2 获取饮食计划列表

**接口**: `GET /plans`

**说明**: 获取饮食计划列表，支持日期和状态筛选

**认证**: 需要

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_date | string | 否 | 开始日期（YYYY-MM-DD） |
| end_date | string | 否 | 结束日期（YYYY-MM-DD） |
| status | string | 否 | 状态筛选：pending, completed, skipped |
| page | integer | 否 | 页码，默认 1 |
| page_size | integer | 否 | 每页数量，默认 20，最大 100 |

**请求示例**:

```
GET /plans?start_date=2024-01-16&status=pending&page=1
```

**响应示例**: 同生成计划响应（带分页）

---

### 4.3 获取单个饮食计划

**接口**: `GET /plans/:id`

**说明**: 获取指定 ID 的饮食计划详情

**认证**: 需要

**响应示例**: 同生成计划响应（单个对象）

---

### 4.4 更新饮食计划

**接口**: `PUT /plans/:id`

**说明**: 更新指定 ID 的饮食计划

**认证**: 需要

**请求参数**:

```json
{
  "plan_date": "2024-01-16T08:00:00Z",
  "meal_type": "breakfast",
  "foods": [
    {
      "food_id": 1,
      "quantity": 2.0
    }
  ],
  "status": "pending",
  "ai_reasoning": "高蛋白低碳水"
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "plan updated successfully"
}
```

---

### 4.5 删除饮食计划

**接口**: `DELETE /plans/:id`

**说明**: 删除指定 ID 的饮食计划

**认证**: 需要

**响应示例**:

```json
{
  "code": 0,
  "message": "plan deleted successfully"
}
```

---

### 4.6 完成饮食计划

**接口**: `POST /plans/:id/complete`

**说明**: 标记计划为已完成，并自动创建对应的餐饮记录

**认证**: 需要

**响应示例**:

```json
{
  "code": 0,
  "message": "plan completed and meal record created",
  "data": {
    "id": 10,
    "user_id": 1,
    "meal_date": "2024-01-16T08:00:00Z",
    "meal_type": "breakfast",
    "foods": [...],
    "nutrition": {...},
    "notes": "来自饮食计划",
    "created_at": "2024-01-16T08:00:00Z",
    "updated_at": "2024-01-16T08:00:00Z"
  }
}
```

---


## 5. AI 服务接口

### 5.1 AI 对话

**接口**: `POST /ai/chat`

**说明**: 与 AI 助手进行对话，获取饮食建议

**认证**: 需要

**请求参数**:

```json
{
  "message": "我今天应该吃什么？",
  "context": {
    "current_calories": "1200",
    "target_calories": "2000"
  }
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| message | string | 是 | 用户消息，1-2000 字符 |
| context | object | 否 | 上下文信息（键值对） |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "response": "根据您今天的摄入情况，建议晚餐选择高蛋白低脂的食物...",
    "conversation_id": "conv_123456"
  }
}
```

---

### 5.2 AI 生成餐饮建议

**接口**: `POST /ai/suggest`

**说明**: 生成 AI 餐饮计划建议

**认证**: 需要

**请求参数**:

```json
{
  "days": 3,
  "target_calories": 2000
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| days | integer | 是 | 建议天数，1-30 天 |
| target_calories | integer | 否 | 目标卡路里，800-10000 |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "suggestions": [
      {
        "date": "2024-01-16",
        "meals": [
          {
            "meal_type": "breakfast",
            "foods": ["鸡蛋", "全麦面包", "牛奶"],
            "calories": 450,
            "reasoning": "高蛋白早餐，提供充足能量"
          }
        ]
      }
    ]
  }
}
```

---

### 5.3 获取对话历史

**接口**: `GET /ai/history`

**说明**: 获取用户的 AI 对话历史记录

**认证**: 需要

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | integer | 否 | 页码，默认 1 |
| page_size | integer | 否 | 每页数量，默认 20，最大 100 |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "message": "我今天应该吃什么？",
      "response": "根据您今天的摄入情况...",
      "created_at": "2024-01-15T12:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 50,
    "total_pages": 3
  }
}
```

---


## 6. 营养分析接口

### 6.1 获取每日营养统计

**接口**: `GET /nutrition/daily/:date`

**说明**: 获取指定日期的营养统计数据

**认证**: 需要

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| date | string | 日期（YYYY-MM-DD） |

**请求示例**:

```
GET /nutrition/daily/2024-01-15
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "date": "2024-01-15",
    "nutrition": {
      "protein": 120.5,
      "carbs": 250.0,
      "fat": 60.0,
      "fiber": 30.0,
      "calories": 2100.0
    },
    "meals": [
      {
        "meal_type": "breakfast",
        "nutrition": {
          "protein": 30.0,
          "carbs": 50.0,
          "fat": 15.0,
          "fiber": 8.0,
          "calories": 450.0
        }
      },
      {
        "meal_type": "lunch",
        "nutrition": {
          "protein": 45.0,
          "carbs": 100.0,
          "fat": 20.0,
          "fiber": 12.0,
          "calories": 800.0
        }
      },
      {
        "meal_type": "dinner",
        "nutrition": {
          "protein": 40.0,
          "carbs": 80.0,
          "fat": 20.0,
          "fiber": 8.0,
          "calories": 700.0
        }
      },
      {
        "meal_type": "snack",
        "nutrition": {
          "protein": 5.5,
          "carbs": 20.0,
          "fat": 5.0,
          "fiber": 2.0,
          "calories": 150.0
        }
      }
    ]
  }
}
```

---

### 6.2 获取月度营养趋势

**接口**: `GET /nutrition/monthly`

**说明**: 获取指定月份的每日营养统计数据

**认证**: 需要

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| year | integer | 是 | 年份（如 2024） |
| month | integer | 是 | 月份（1-12） |

**请求示例**:

```
GET /nutrition/monthly?year=2024&month=1
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "date": "2024-01-01",
      "nutrition": {
        "protein": 100.0,
        "carbs": 200.0,
        "fat": 50.0,
        "fiber": 25.0,
        "calories": 1800.0
      }
    },
    {
      "date": "2024-01-02",
      "nutrition": {
        "protein": 120.0,
        "carbs": 250.0,
        "fat": 60.0,
        "fiber": 30.0,
        "calories": 2100.0
      }
    }
  ]
}
```

---

### 6.3 对比实际与目标营养

**接口**: `GET /nutrition/compare`

**说明**: 对比指定日期的实际营养摄入与用户目标

**认证**: 需要

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| date | string | 是 | 日期（YYYY-MM-DD） |

**请求示例**:

```
GET /nutrition/compare?date=2024-01-15
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "date": "2024-01-15",
    "actual": {
      "protein": 120.5,
      "carbs": 250.0,
      "fat": 60.0,
      "fiber": 30.0,
      "calories": 2100.0
    },
    "target": {
      "protein": 150.0,
      "carbs": 200.0,
      "fat": 66.7,
      "fiber": 25.0,
      "calories": 2000.0
    },
    "difference": {
      "protein": -29.5,
      "carbs": 50.0,
      "fat": -6.7,
      "fiber": 5.0,
      "calories": 100.0
    },
    "percentage": {
      "protein": 80.3,
      "carbs": 125.0,
      "fat": 90.0,
      "fiber": 120.0,
      "calories": 105.0
    }
  }
}
```

---


## 7. Dashboard 接口

### 7.1 获取 Dashboard 数据

**接口**: `GET /dashboard`

**说明**: 获取综合面板数据，包括今日营养、月度统计、未来计划

**认证**: 需要

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "today": {
      "date": "2024-01-15",
      "nutrition": {
        "protein": 120.5,
        "carbs": 250.0,
        "fat": 60.0,
        "fiber": 30.0,
        "calories": 2100.0
      },
      "target": {
        "protein": 150.0,
        "carbs": 200.0,
        "fat": 66.7,
        "fiber": 25.0,
        "calories": 2000.0
      },
      "percentage": {
        "protein": 80.3,
        "carbs": 125.0,
        "fat": 90.0,
        "fiber": 120.0,
        "calories": 105.0
      }
    },
    "monthly_stats": {
      "year": 2024,
      "month": 1,
      "average_calories": 1950.0,
      "total_meals": 90,
      "days_tracked": 15
    },
    "upcoming_plans": [
      {
        "id": 1,
        "plan_date": "2024-01-16T08:00:00Z",
        "meal_type": "breakfast",
        "foods": [...],
        "nutrition": {...},
        "status": "pending"
      }
    ],
    "recent_meals": [
      {
        "id": 10,
        "meal_date": "2024-01-15T18:00:00Z",
        "meal_type": "dinner",
        "foods": [...],
        "nutrition": {...}
      }
    ]
  }
}
```

---


## 8. 设置管理接口

### 8.1 获取所有设置

**接口**: `GET /settings`

**说明**: 获取用户的 AI 设置和偏好设置

**认证**: 需要

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "ai_settings": [
      {
        "id": 1,
        "user_id": 1,
        "provider": "openai",
        "api_endpoint": "https://api.openai.com/v1",
        "model": "gpt-3.5-turbo",
        "temperature": 0.7,
        "max_tokens": 1000,
        "is_active": true,
        "created_at": "2024-01-15T10:00:00Z",
        "updated_at": "2024-01-15T10:00:00Z"
      }
    ],
    "user_preferences": {
      "id": 1,
      "user_id": 1,
      "taste_preferences": ["高蛋白", "低碳水"],
      "dietary_restrictions": ["不吃辣", "不吃海鲜"],
      "daily_calorie_target": 2000,
      "preferred_meal_times": {
        "breakfast": "08:00",
        "lunch": "12:00",
        "dinner": "18:00"
      },
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  }
}
```

---

### 8.2 更新 AI 设置

**接口**: `PUT /settings/ai`

**说明**: 更新用户的 AI 配置

**认证**: 需要

**请求参数**:

```json
{
  "provider": "openai",
  "api_endpoint": "https://api.openai.com/v1",
  "api_key": "sk-xxxxxxxxxxxxxxxx",
  "model": "gpt-3.5-turbo",
  "temperature": 0.7,
  "max_tokens": 1000,
  "is_active": true
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| provider | string | 是 | AI 提供商：openai, deepseek, custom |
| api_endpoint | string | 否 | API 端点 URL |
| api_key | string | 是 | API 密钥，10-500 字符 |
| model | string | 否 | 模型名称，最多 100 字符 |
| temperature | number | 否 | 温度参数，0-2，默认 0.7 |
| max_tokens | integer | 否 | 最大 Token 数，1-32000，默认 1000 |
| is_active | boolean | 否 | 是否激活，默认 true |

**响应示例**:

```json
{
  "code": 0,
  "message": "AI settings updated successfully"
}
```

---

### 8.3 测试 AI 连接

**接口**: `GET /settings/ai/test`

**说明**: 测试当前配置的 AI 服务是否可用

**认证**: 需要

**响应示例**:

```json
{
  "code": 0,
  "message": "AI connection test successful"
}
```

**错误示例**:

```json
{
  "code": 50001,
  "message": "AI connection test failed: invalid API key",
  "error": "401 Unauthorized"
}
```

---

### 8.4 获取用户资料

**接口**: `GET /user/profile`

**说明**: 获取用户基本信息和偏好设置

**认证**: 需要

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "user_id": 1,
    "taste_preferences": ["高蛋白", "低碳水"],
    "dietary_restrictions": ["不吃辣", "不吃海鲜"],
    "daily_calorie_target": 2000,
    "preferred_meal_times": {
      "breakfast": "08:00",
      "lunch": "12:00",
      "dinner": "18:00"
    },
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
}
```

---

### 8.5 更新用户偏好

**接口**: `PUT /user/preferences`

**说明**: 更新用户的口味偏好、饮食限制和营养目标

**认证**: 需要

**请求参数**:

```json
{
  "taste_preferences": ["高蛋白", "低碳水", "清淡"],
  "dietary_restrictions": ["不吃辣", "不吃海鲜"],
  "daily_calorie_target": 2000,
  "preferred_meal_times": {
    "breakfast": "08:00",
    "lunch": "12:00",
    "dinner": "18:00",
    "snack": "15:00"
  }
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| taste_preferences | array | 否 | 口味偏好列表 |
| dietary_restrictions | array | 否 | 饮食限制列表 |
| daily_calorie_target | integer | 否 | 每日卡路里目标，默认 2000 |
| preferred_meal_times | object | 否 | 偏好用餐时间（键值对） |

**响应示例**:

```json
{
  "code": 0,
  "message": "user preferences updated successfully"
}
```

---


## 9. 健康检查接口

### 9.1 健康检查

**接口**: `GET /health`

**说明**: 检查服务是否正常运行（无需认证）

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "ok",
    "service": "ai-diet-assistant"
  }
}
```

---

## 附录

### A. 数据类型说明

#### Food（食材）

```typescript
interface Food {
  id: number;
  user_id: number;
  name: string;
  category: "meat" | "vegetable" | "fruit" | "grain" | "other";
  price: number;
  unit: string;
  protein: number;
  carbs: number;
  fat: number;
  fiber: number;
  calories: number;
  available: boolean;
  created_at: string;
  updated_at: string;
}
```

#### Meal（餐饮记录）

```typescript
interface Meal {
  id: number;
  user_id: number;
  meal_date: string;
  meal_type: "breakfast" | "lunch" | "dinner" | "snack";
  foods: MealFood[];
  nutrition: NutritionData;
  notes: string;
  created_at: string;
  updated_at: string;
}

interface MealFood {
  food_id: number;
  food_name: string;
  quantity: number;
  unit: string;
}
```

#### Plan（饮食计划）

```typescript
interface Plan {
  id: number;
  user_id: number;
  plan_date: string;
  meal_type: "breakfast" | "lunch" | "dinner" | "snack";
  foods: MealFood[];
  nutrition: NutritionData;
  status: "pending" | "completed" | "skipped";
  ai_reasoning: string;
  created_at: string;
  updated_at: string;
}
```

#### NutritionData（营养数据）

```typescript
interface NutritionData {
  protein: number;   // 蛋白质（克）
  carbs: number;     // 碳水化合物（克）
  fat: number;       // 脂肪（克）
  fiber: number;     // 纤维（克）
  calories: number;  // 卡路里
}
```

#### UserPreferences（用户偏好）

```typescript
interface UserPreferences {
  id: number;
  user_id: number;
  taste_preferences: string[];
  dietary_restrictions: string[];
  daily_calorie_target: number;
  preferred_meal_times: Record<string, string>;
  created_at: string;
  updated_at: string;
}
```

#### AISettings（AI 设置）

```typescript
interface AISettings {
  id: number;
  user_id: number;
  provider: "openai" | "deepseek" | "custom";
  api_endpoint: string;
  model: string;
  temperature: number;
  max_tokens: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
```

---

### B. 常见问题

#### 1. 如何处理 Token 过期？

当 Access Token 过期时，API 会返回 401 错误。此时应该：
1. 使用 Refresh Token 调用 `/auth/refresh` 获取新的 Access Token
2. 使用新的 Access Token 重试原请求
3. 如果 Refresh Token 也过期，跳转到登录页

#### 2. 如何实现分页？

所有列表接口都支持分页，使用 `page` 和 `page_size` 参数：
- `page`: 页码，从 1 开始
- `page_size`: 每页数量，默认 20，最大 100

响应中包含 `pagination` 对象，包含总数和总页数。

#### 3. 日期格式要求

- 日期时间：ISO 8601 格式（`2024-01-15T12:00:00Z`）
- 仅日期：`YYYY-MM-DD` 格式（`2024-01-15`）

#### 4. 如何处理错误？

所有错误响应都包含 `code` 和 `message` 字段：
- 检查 `code` 判断错误类型
- 显示 `message` 给用户
- 根据错误码采取相应措施（如 401 跳转登录）

#### 5. 营养数据如何计算？

营养数据由后端自动计算：
- 创建/更新餐饮记录时，根据食材和份量自动计算
- 生成饮食计划时，AI 会考虑营养平衡
- 营养统计接口会汇总指定时间范围的数据

---

### C. 开发建议

1. **使用 TypeScript**: 定义接口类型，提高代码可维护性
2. **封装 API 请求**: 创建统一的 API 请求函数，处理认证和错误
3. **实现请求拦截器**: 自动添加 Token，自动刷新过期 Token
4. **错误处理**: 统一处理 API 错误，显示友好提示
5. **加载状态**: 所有 API 请求显示加载动画
6. **数据缓存**: 缓存不常变化的数据（如食材列表）
7. **防抖节流**: 搜索、滚动等操作使用防抖或节流

---

### D. 测试建议

#### 使用 cURL 测试

```bash
# 登录
curl -X POST http://localhost:9090/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"password123"}'

# 获取食材列表
curl -X GET "http://localhost:9090/api/v1/foods?page=1&page_size=20" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 创建食材
curl -X POST http://localhost:9090/api/v1/foods \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "鸡胸肉",
    "category": "meat",
    "price": 15.5,
    "unit": "100g",
    "protein": 23.0,
    "carbs": 0.0,
    "fat": 1.2,
    "fiber": 0.0,
    "calories": 105.0,
    "available": true
  }'
```

#### 使用 Postman

1. 导入 API 集合（可以根据此文档创建）
2. 设置环境变量（base_url, access_token）
3. 测试所有端点
4. 验证响应格式和数据

---

## 更新日志

### v1.0.0 (2024-01-15)
- 初始版本
- 实现所有核心功能接口
- 支持 JWT 认证
- 支持分页和筛选

---

**文档版本**: v1.0.0  
**最后更新**: 2024-01-15  
**维护者**: AI Diet Assistant Team
