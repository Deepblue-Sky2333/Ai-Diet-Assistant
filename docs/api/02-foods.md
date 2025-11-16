# 食材管理模块

## 概述

食材管理模块提供用户市场面板中食材的完整管理功能，包括食材的创建、查询、更新、删除和批量导入。用户可以维护自己的食材库，记录食材的营养信息和价格，为餐饮记录和饮食计划提供基础数据。

**核心功能**：
- 创建单个食材项
- 查询食材列表（支持分类和可用性过滤）
- 查询单个食材详情
- 更新食材信息
- 删除食材
- 批量导入食材

**数据特性**：
- 每个食材包含完整的营养信息（蛋白质、碳水化合物、脂肪、纤维、热量）
- 支持食材分类（肉类、蔬菜、水果、谷物、其他）
- 支持自定义单位和价格
- 支持可用性标记

---

## 接口列表

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/foods` | 创建食材 | 是 |
| GET | `/api/v1/foods` | 获取食材列表 | 是 |
| GET | `/api/v1/foods/:id` | 获取单个食材 | 是 |
| PUT | `/api/v1/foods/:id` | 更新食材 | 是 |
| DELETE | `/api/v1/foods/:id` | 删除食材 | 是 |
| POST | `/api/v1/foods/batch` | 批量导入食材 | 是 |

---

## 接口详情

### 创建食材

**接口**: `POST /api/v1/foods`

**说明**: 在用户的市场面板中创建一个新的食材项。食材创建后可用于餐饮记录和饮食计划。

**认证**: 是

#### 请求参数

##### 请求体

```json
{
  "name": "鸡胸肉",
  "category": "meat",
  "price": 15.99,
  "unit": "100g",
  "protein": 23.0,
  "carbs": 0.0,
  "fat": 1.2,
  "fiber": 0.0,
  "calories": 110.0,
  "available": true
}
```

| 字段 | 类型 | 必填 | 说明 | 验证规则 |
|------|------|------|------|----------|
| name | string | 是 | 食材名称 | 长度 1-100 字符 |
| category | string | 是 | 食材分类 | 枚举值：meat, vegetable, fruit, grain, other |
| price | number | 是 | 价格 | ≥ 0，≤ 100000 |
| unit | string | 是 | 单位 | 长度 1-20 字符，如 "100g", "个", "ml" |
| protein | number | 是 | 蛋白质含量（克/单位） | ≥ 0，≤ 1000 |
| carbs | number | 是 | 碳水化合物含量（克/单位） | ≥ 0，≤ 1000 |
| fat | number | 是 | 脂肪含量（克/单位） | ≥ 0，≤ 1000 |
| fiber | number | 是 | 纤维含量（克/单位） | ≥ 0，≤ 1000 |
| calories | number | 是 | 热量（千卡/单位） | ≥ 0，≤ 10000 |
| available | boolean | 否 | 是否可用 | 默认 true |

#### 请求示例

```bash
curl -X POST http://localhost:9090/api/v1/foods \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "鸡胸肉",
    "category": "meat",
    "price": 15.99,
    "unit": "100g",
    "protein": 23.0,
    "carbs": 0.0,
    "fat": 1.2,
    "fiber": 0.0,
    "calories": 110.0,
    "available": true
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
    "name": "鸡胸肉",
    "category": "meat",
    "price": 15.99,
    "unit": "100g",
    "protein": 23.0,
    "carbs": 0.0,
    "fat": 1.2,
    "fiber": 0.0,
    "calories": 110.0,
    "available": true,
    "created_at": "2024-11-16T10:30:00Z",
    "updated_at": "2024-11-16T10:30:00Z"
  },
  "timestamp": 1699999999
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 食材唯一标识符 |
| user_id | int64 | 所属用户 ID |
| name | string | 食材名称 |
| category | string | 食材分类 |
| price | number | 价格 |
| unit | string | 单位 |
| protein | number | 蛋白质含量（克/单位） |
| carbs | number | 碳水化合物含量（克/单位） |
| fat | number | 脂肪含量（克/单位） |
| fiber | number | 纤维含量（克/单位） |
| calories | number | 热量（千卡/单位） |
| available | boolean | 是否可用 |
| created_at | string | 创建时间（ISO 8601 格式） |
| updated_at | string | 更新时间（ISO 8601 格式） |

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid request parameters: Key: 'CreateFoodRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag",
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
| 50001 | 内部错误 | 服务器内部错误 |

#### 注意事项

1. **营养数据**：所有营养数据都是基于指定单位的含量
2. **单位灵活性**：单位可以是重量（如 "100g"）、体积（如 "ml"）或数量（如 "个"）
3. **分类枚举**：category 必须是以下值之一：meat（肉类）、vegetable（蔬菜）、fruit（水果）、grain（谷物）、other（其他）
4. **可用性标记**：available 字段用于标记食材是否可用，不可用的食材仍保留在数据库中

---

### 获取食材列表

**接口**: `GET /api/v1/foods`

**说明**: 获取用户的食材列表，支持按分类和可用性过滤，支持分页查询。

**认证**: 是

#### 请求参数

##### 查询参数

| 参数 | 类型 | 必填 | 说明 | 默认值 | 示例 |
|------|------|------|------|--------|------|
| category | string | 否 | 按分类过滤 | - | meat, vegetable, fruit, grain, other |
| available | boolean | 否 | 按可用性过滤 | - | true, false |
| page | int | 否 | 页码（从 1 开始） | 1 | 1, 2, 3 |
| page_size | int | 否 | 每页数据量 | 20 | 10, 20, 50（最大 100） |

#### 请求示例

```bash
# 获取所有食材（默认分页）
curl -X GET "http://localhost:9090/api/v1/foods" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取肉类食材
curl -X GET "http://localhost:9090/api/v1/foods?category=meat" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取可用的食材
curl -X GET "http://localhost:9090/api/v1/foods?available=true" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取第 2 页，每页 50 条
curl -X GET "http://localhost:9090/api/v1/foods?page=2&page_size=50" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 组合过滤：获取可用的蔬菜类食材
curl -X GET "http://localhost:9090/api/v1/foods?category=vegetable&available=true" \
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
      "name": "鸡胸肉",
      "category": "meat",
      "price": 15.99,
      "unit": "100g",
      "protein": 23.0,
      "carbs": 0.0,
      "fat": 1.2,
      "fiber": 0.0,
      "calories": 110.0,
      "available": true,
      "created_at": "2024-11-16T10:30:00Z",
      "updated_at": "2024-11-16T10:30:00Z"
    },
    {
      "id": 2,
      "user_id": 1,
      "name": "西兰花",
      "category": "vegetable",
      "price": 8.50,
      "unit": "100g",
      "protein": 2.8,
      "carbs": 6.6,
      "fat": 0.4,
      "fiber": 2.6,
      "calories": 34.0,
      "available": true,
      "created_at": "2024-11-16T10:35:00Z",
      "updated_at": "2024-11-16T10:35:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 50,
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
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | 服务器内部错误 |

#### 注意事项

1. **分页默认值**：不指定分页参数时，默认返回第 1 页，每页 20 条数据
2. **分页限制**：page_size 最大值为 100，超过会自动调整为 100
3. **过滤组合**：可以同时使用多个过滤条件
4. **空结果**：当没有符合条件的数据时，返回空数组，不会报错

---

### 获取单个食材

**接口**: `GET /api/v1/foods/:id`

**说明**: 根据食材 ID 获取单个食材的详细信息。

**认证**: 是

#### 请求参数

##### 路径参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | int64 | 是 | 食材 ID | 1 |

#### 请求示例

```bash
curl -X GET http://localhost:9090/api/v1/foods/1 \
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
    "name": "鸡胸肉",
    "category": "meat",
    "price": 15.99,
    "unit": "100g",
    "protein": 23.0,
    "carbs": 0.0,
    "fat": 1.2,
    "fiber": 0.0,
    "calories": 110.0,
    "available": true,
    "created_at": "2024-11-16T10:30:00Z",
    "updated_at": "2024-11-16T10:30:00Z"
  },
  "timestamp": 1699999999
}
```

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid food id",
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
  "error": "food not found",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | 食材 ID 格式不正确 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 40401 | 资源不存在 | 食材不存在或不属于当前用户 |
| 50001 | 内部错误 | 服务器内部错误 |

#### 注意事项

1. **权限验证**：只能查询属于当前用户的食材
2. **ID 格式**：ID 必须是有效的整数
3. **不存在处理**：如果食材不存在或不属于当前用户，返回 404 错误

---

### 更新食材

**接口**: `PUT /api/v1/foods/:id`

**说明**: 更新指定食材的信息。需要提供完整的食材数据，部分更新不支持。

**认证**: 是

#### 请求参数

##### 路径参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | int64 | 是 | 食材 ID | 1 |

##### 请求体

```json
{
  "name": "鸡胸肉（有机）",
  "category": "meat",
  "price": 18.99,
  "unit": "100g",
  "protein": 24.0,
  "carbs": 0.0,
  "fat": 1.0,
  "fiber": 0.0,
  "calories": 115.0,
  "available": true
}
```

| 字段 | 类型 | 必填 | 说明 | 验证规则 |
|------|------|------|------|----------|
| name | string | 是 | 食材名称 | 长度 1-100 字符 |
| category | string | 是 | 食材分类 | 枚举值：meat, vegetable, fruit, grain, other |
| price | number | 是 | 价格 | ≥ 0，≤ 100000 |
| unit | string | 是 | 单位 | 长度 1-20 字符 |
| protein | number | 是 | 蛋白质含量（克/单位） | ≥ 0，≤ 1000 |
| carbs | number | 是 | 碳水化合物含量（克/单位） | ≥ 0，≤ 1000 |
| fat | number | 是 | 脂肪含量（克/单位） | ≥ 0，≤ 1000 |
| fiber | number | 是 | 纤维含量（克/单位） | ≥ 0，≤ 1000 |
| calories | number | 是 | 热量（千卡/单位） | ≥ 0，≤ 10000 |
| available | boolean | 否 | 是否可用 | 默认 true |

#### 请求示例

```bash
curl -X PUT http://localhost:9090/api/v1/foods/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "鸡胸肉（有机）",
    "category": "meat",
    "price": 18.99,
    "unit": "100g",
    "protein": 24.0,
    "carbs": 0.0,
    "fat": 1.0,
    "fiber": 0.0,
    "calories": 115.0,
    "available": true
  }'
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "food updated successfully",
  "data": null,
  "timestamp": 1699999999
}
```

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid request parameters: Key: 'UpdateFoodRequest.Price' Error:Field validation for 'Price' failed on the 'gte' tag",
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
  "error": "food not found",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | 缺少必填字段、参数类型不匹配、参数值超出范围 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 40401 | 资源不存在 | 食材不存在或不属于当前用户 |
| 50001 | 内部错误 | 服务器内部错误 |

#### 注意事项

1. **完整更新**：需要提供所有必填字段，不支持部分更新
2. **权限验证**：只能更新属于当前用户的食材
3. **ID 不可变**：食材 ID 和用户 ID 不会被更新
4. **时间戳自动更新**：updated_at 字段会自动更新为当前时间

---

### 删除食材

**接口**: `DELETE /api/v1/foods/:id`

**说明**: 从用户的市场面板中删除指定的食材。删除操作是物理删除，数据将无法恢复。

**认证**: 是

#### 请求参数

##### 路径参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | int64 | 是 | 食材 ID | 1 |

#### 请求示例

```bash
curl -X DELETE http://localhost:9090/api/v1/foods/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "food deleted successfully",
  "data": null,
  "timestamp": 1699999999
}
```

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid food id",
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
  "error": "food not found",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | 食材 ID 格式不正确 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 40401 | 资源不存在 | 食材不存在或不属于当前用户 |
| 50001 | 内部错误 | 服务器内部错误 |

#### 注意事项

1. **物理删除**：删除操作是永久性的，数据无法恢复
2. **权限验证**：只能删除属于当前用户的食材
3. **关联数据**：删除食材前，请确保没有餐饮记录或计划引用该食材
4. **建议做法**：如果不确定是否要永久删除，可以使用更新接口将 available 设置为 false

---

### 批量导入食材

**接口**: `POST /api/v1/foods/batch`

**说明**: 批量导入多个食材项。适用于初始化食材库或从其他系统迁移数据。系统会验证每个食材项，验证失败的项会被跳过，成功的项会被导入。

**认证**: 是

#### 请求参数

##### 请求体

```json
{
  "foods": [
    {
      "name": "鸡胸肉",
      "category": "meat",
      "price": 15.99,
      "unit": "100g",
      "protein": 23.0,
      "carbs": 0.0,
      "fat": 1.2,
      "fiber": 0.0,
      "calories": 110.0,
      "available": true
    },
    {
      "name": "西兰花",
      "category": "vegetable",
      "price": 8.50,
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

| 字段 | 类型 | 必填 | 说明 | 验证规则 |
|------|------|------|------|----------|
| foods | array | 是 | 食材列表 | 最少 1 项，最多 100 项 |
| foods[].name | string | 是 | 食材名称 | 长度 1-100 字符 |
| foods[].category | string | 是 | 食材分类 | 枚举值：meat, vegetable, fruit, grain, other |
| foods[].price | number | 是 | 价格 | ≥ 0，≤ 100000 |
| foods[].unit | string | 是 | 单位 | 长度 1-20 字符 |
| foods[].protein | number | 是 | 蛋白质含量（克/单位） | ≥ 0，≤ 1000 |
| foods[].carbs | number | 是 | 碳水化合物含量（克/单位） | ≥ 0，≤ 1000 |
| foods[].fat | number | 是 | 脂肪含量（克/单位） | ≥ 0，≤ 1000 |
| foods[].fiber | number | 是 | 纤维含量（克/单位） | ≥ 0，≤ 1000 |
| foods[].calories | number | 是 | 热量（千卡/单位） | ≥ 0，≤ 10000 |
| foods[].available | boolean | 否 | 是否可用 | 默认 true |

#### 请求示例

```bash
curl -X POST http://localhost:9090/api/v1/foods/batch \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "foods": [
      {
        "name": "鸡胸肉",
        "category": "meat",
        "price": 15.99,
        "unit": "100g",
        "protein": 23.0,
        "carbs": 0.0,
        "fat": 1.2,
        "fiber": 0.0,
        "calories": 110.0,
        "available": true
      },
      {
        "name": "西兰花",
        "category": "vegetable",
        "price": 8.50,
        "unit": "100g",
        "protein": 2.8,
        "carbs": 6.6,
        "fat": 0.4,
        "fiber": 2.6,
        "calories": 34.0,
        "available": true
      }
    ]
  }'
```

#### 响应示例

**成功响应 (200) - 全部成功**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "success": 2,
    "failed": 0,
    "errors": []
  },
  "timestamp": 1699999999
}
```

**成功响应 (200) - 部分成功**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "success": 1,
    "failed": 1,
    "errors": [
      "row 2: Key: 'Food.Name' Error:Field validation for 'Name' failed on the 'required' tag"
    ]
  },
  "timestamp": 1699999999
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| success | int | 成功导入的食材数量 |
| failed | int | 导入失败的食材数量 |
| errors | array | 错误信息列表（可选，仅在有失败项时返回） |

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid request parameters: Key: 'BatchImportRequest.Foods' Error:Field validation for 'Foods' failed on the 'required' tag",
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
| 40001 | 参数错误 | foods 数组为空、超过最大数量限制 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | 服务器内部错误 |

#### 注意事项

1. **数量限制**：单次最多导入 100 个食材项
2. **部分成功**：即使部分食材验证失败，成功的食材仍会被导入
3. **错误信息**：errors 数组包含每个失败项的行号和错误原因
4. **事务处理**：成功的食材会被批量插入，提高性能
5. **验证规则**：每个食材项都会进行完整的验证，与单个创建接口的验证规则相同

---

## 数据模型

### Food 模型

完整的 Food 数据模型定义请参考 [数据模型文档](./data-models.md#food-食材)。

**核心字段**：
- **id**: 食材唯一标识符
- **user_id**: 所属用户 ID
- **name**: 食材名称
- **category**: 食材分类（meat, vegetable, fruit, grain, other）
- **price**: 价格
- **unit**: 单位
- **protein**: 蛋白质含量（克/单位）
- **carbs**: 碳水化合物含量（克/单位）
- **fat**: 脂肪含量（克/单位）
- **fiber**: 纤维含量（克/单位）
- **calories**: 热量（千卡/单位）
- **available**: 是否可用
- **created_at**: 创建时间
- **updated_at**: 更新时间

### 食材分类说明

| 分类值 | 中文名称 | 说明 | 示例 |
|--------|---------|------|------|
| meat | 肉类 | 各种肉类和海鲜 | 鸡胸肉、牛肉、鱼肉、虾 |
| vegetable | 蔬菜 | 各种蔬菜 | 西兰花、菠菜、胡萝卜、番茄 |
| fruit | 水果 | 各种水果 | 苹果、香蕉、橙子、草莓 |
| grain | 谷物 | 谷物和主食 | 米饭、面条、燕麦、面包 |
| other | 其他 | 其他食材 | 牛奶、鸡蛋、坚果、调味料 |

---

## 使用场景

### 场景 1：初始化食材库

用户首次使用系统时，可以通过批量导入接口快速建立食材库：

```bash
# 准备食材数据（可以从 Excel 或其他系统导出）
curl -X POST http://localhost:9090/api/v1/foods/batch \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d @foods_data.json
```

### 场景 2：查询可用的肉类食材

在创建餐饮记录时，查询可用的肉类食材供用户选择：

```bash
curl -X GET "http://localhost:9090/api/v1/foods?category=meat&available=true" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 场景 3：更新食材价格

当市场价格变化时，更新食材的价格信息：

```bash
curl -X PUT http://localhost:9090/api/v1/foods/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "鸡胸肉",
    "category": "meat",
    "price": 16.99,
    "unit": "100g",
    "protein": 23.0,
    "carbs": 0.0,
    "fat": 1.2,
    "fiber": 0.0,
    "calories": 110.0,
    "available": true
  }'
```

### 场景 4：标记食材为不可用

当某个食材暂时缺货时，可以标记为不可用而不删除：

```bash
curl -X PUT http://localhost:9090/api/v1/foods/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "鸡胸肉",
    "category": "meat",
    "price": 15.99,
    "unit": "100g",
    "protein": 23.0,
    "carbs": 0.0,
    "fat": 1.2,
    "fiber": 0.0,
    "calories": 110.0,
    "available": false
  }'
```

---

## 最佳实践

### 1. 食材命名

- **清晰明确**：使用准确的食材名称，如 "鸡胸肉" 而不是 "鸡肉"
- **包含规格**：如果有特殊规格，可以在名称中注明，如 "鸡胸肉（有机）"
- **避免重复**：同一食材不要创建多个条目，使用更新接口修改信息

### 2. 单位设置

- **统一标准**：建议使用标准单位，如 "100g"、"ml"、"个"
- **营养对应**：确保营养数据与单位对应，如 "100g" 对应的是每 100 克的营养含量
- **便于计算**：使用便于计算的单位，如 "100g" 比 "1斤" 更适合营养计算

### 3. 营养数据

- **准确性**：尽量使用权威来源的营养数据（如食物成分表）
- **一致性**：同类食材使用相同的数据来源
- **定期更新**：营养数据可能会更新，建议定期检查和更新

### 4. 分类管理

- **合理分类**：根据食材的主要特征选择分类
- **一致性**：同类食材使用相同的分类
- **灵活使用 other**：对于难以归类的食材，使用 "other" 分类

### 5. 可用性管理

- **软删除**：对于暂时不用的食材，使用 available=false 而不是删除
- **定期清理**：定期检查不可用的食材，确认是否需要删除
- **保留历史**：如果食材已被餐饮记录引用，建议保留而不是删除

### 6. 批量导入

- **数据准备**：导入前仔细检查数据格式和内容
- **分批导入**：如果数据量大，建议分批导入（每批不超过 100 条）
- **错误处理**：检查返回的错误信息，修正失败的数据后重新导入
- **备份数据**：导入前备份原始数据，以便出错时恢复

### 7. 性能优化

- **使用过滤**：查询时使用 category 和 available 过滤，减少数据量
- **合理分页**：根据实际需求设置 page_size，避免一次加载过多数据
- **缓存结果**：对于不常变化的食材列表，可以在客户端缓存

---

## 常见问题

### Q: 如何计算食材的营养数据？

A: 
- 营养数据应基于指定的单位
- 例如：单位是 "100g"，则 protein=23.0 表示每 100 克含 23 克蛋白质
- 在餐饮记录中使用时，系统会根据实际用量自动计算营养

### Q: 可以修改食材的单位吗？

A: 
- 可以，使用更新接口可以修改单位
- 但要注意同时更新营养数据，确保数据与新单位对应
- 如果食材已被餐饮记录引用，修改单位可能影响历史数据的准确性

### Q: 删除食材会影响已有的餐饮记录吗？

A: 
- 删除食材不会删除餐饮记录中的数据
- 但餐饮记录中的 food_id 可能会失效
- 建议使用 available=false 标记不可用，而不是删除

### Q: 批量导入时部分失败怎么办？

A: 
- 系统会返回成功和失败的数量，以及失败的错误信息
- 成功的食材已经导入，不需要重新导入
- 根据错误信息修正失败的数据，然后重新导入失败的部分

### Q: 如何处理同名但不同规格的食材？

A: 
- 在名称中包含规格信息，如 "鸡胸肉（有机）" 和 "鸡胸肉（普通）"
- 或者使用不同的价格和营养数据区分
- 避免创建完全相同的食材

### Q: 食材列表支持搜索吗？

A: 
- 当前版本不支持名称搜索
- 可以使用分类过滤缩小范围
- 建议在客户端实现搜索功能

### Q: 如何导出食材数据？

A: 
- 使用列表接口获取所有食材数据
- 设置较大的 page_size（最大 100）
- 多次请求获取所有页的数据
- 在客户端处理和导出数据

---

## 相关文档

- [数据模型](./data-models.md) - 查看 Food 模型的完整定义
- [餐饮记录模块](./03-meals.md) - 了解如何使用食材创建餐饮记录
- [饮食计划模块](./04-plans.md) - 了解如何使用食材创建饮食计划
- [通用概念](./common-concepts.md) - 了解认证、分页等通用概念
- [错误码说明](./error-codes.md) - 查看所有错误码的详细说明
- [API 文档总览](./README.md) - 返回 API 文档首页
