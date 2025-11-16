# AI 服务模块

## 概述

AI 服务模块提供智能对话和餐饮建议功能，是系统的核心智能服务。用户可以与 AI 助手进行自然语言对话，获取饮食建议、营养知识和个性化推荐。AI 会根据用户的食材库、饮食记录和营养目标，提供专业的饮食指导。系统还会保存对话历史，方便用户回顾和追踪。

**核心功能**：
- 与 AI 助手进行自然语言对话
- 获取 AI 生成的餐饮建议
- 查询对话历史记录（支持分页）

**AI 特性**：
- 支持多种 AI 提供商（OpenAI、DeepSeek、自定义）
- 上下文感知对话（可传递用户数据作为上下文）
- Token 使用统计
- 对话历史持久化
- 智能重试机制

**数据特性**：
- 对话历史包含用户输入、AI 响应和上下文信息
- 支持传递结构化上下文数据（如当前日期、营养目标等）
- 自动记录 Token 使用量
- 支持分页查询历史记录

---

## 接口列表

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/ai/chat` | AI 对话 | 是 |
| POST | `/api/v1/ai/suggest` | AI 生成餐饮建议 | 是 |
| GET | `/api/v1/ai/history` | 获取对话历史 | 是 |

---

## 接口详情


### AI 对话

**接口**: `POST /api/v1/ai/chat`

**说明**: 与 AI 饮食助手进行自然语言对话。用户可以询问饮食建议、营养知识、食材搭配等问题，AI 会根据用户的个人情况提供专业的回答。支持传递上下文信息，使对话更加智能和个性化。

**认证**: 是


#### 请求参数

##### 请求体

```json
{
  "message": "我今天应该吃什么？",
  "context": {
    "current_date": "2024-11-17",
    "target_calories": "2000"
  }
}
```

| 字段 | 类型 | 必填 | 说明 | 验证规则 |
|------|------|------|------|----------|
| message | string | 是 | 用户消息内容 | 最小 1 字符，最大 2000 字符 |
| context | object | 否 | 上下文信息（键值对） | 可选，用于提供额外的上下文 |

#### 请求示例

```bash
# 简单对话
curl -X POST http://localhost:9090/api/v1/ai/chat \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "我今天应该吃什么？"
  }'

# 带上下文的对话
curl -X POST http://localhost:9090/api/v1/ai/chat \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "根据我的营养目标，推荐今天的晚餐",
    "context": {
      "current_date": "2024-11-17",
      "target_calories": "2000",
      "target_protein": "150"
    }
  }'
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "根据您的营养目标（每日 2000 千卡，蛋白质 150 克），我建议今天的晚餐可以选择：\n\n1. 主食：糙米饭 100 克（约 350 千卡）\n2. 蛋白质：鸡胸肉 200 克（约 330 千卡，蛋白质 62 克）\n3. 蔬菜：西兰花 150 克 + 胡萝卜 100 克（约 80 千卡）\n4. 健康脂肪：橄榄油 10 克（约 90 千卡）\n\n这份晚餐总计约 850 千卡，蛋白质约 65 克，营养均衡且符合您的目标。建议搭配清淡的烹饪方式，如蒸、煮或少油炒。",
    "response": "根据您的营养目标（每日 2000 千卡，蛋白质 150 克），我建议今天的晚餐可以选择：\n\n1. 主食：糙米饭 100 克（约 350 千卡）\n2. 蛋白质：鸡胸肉 200 克（约 330 千卡，蛋白质 62 克）\n3. 蔬菜：西兰花 150 克 + 胡萝卜 100 克（约 80 千卡）\n4. 健康脂肪：橄榄油 10 克（约 90 千卡）\n\n这份晚餐总计约 850 千卡，蛋白质约 65 克，营养均衡且符合您的目标。建议搭配清淡的烹饪方式，如蒸、煮或少油炒。",
    "message_id": 123,
    "tokens_used": 256
  },
  "timestamp": 1699999999
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| message | string | AI 的回复内容 |
| response | string | AI 的回复内容（与 message 相同，用于前端兼容） |
| message_id | int64 | 对话记录 ID（用于查询历史） |
| tokens_used | int | 本次对话使用的 Token 数量 |


**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid request parameters: Key: 'ChatRequest.Message' Error:Field validation for 'Message' failed on the 'required' tag",
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
  "error": "AI chat failed",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | message 为空、message 超过 2000 字符 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | AI 服务调用失败、AI 设置未配置、网络错误 |

#### 注意事项

1. **消息长度**：message 字段最少 1 字符，最多 2000 字符
2. **上下文信息**：context 是可选的，可以传递任何键值对信息
3. **上下文作用**：传递上下文可以让 AI 更好地理解用户的情况，提供更个性化的建议
4. **常用上下文**：current_date（当前日期）、target_calories（目标热量）、target_protein（目标蛋白质）等
5. **Token 统计**：tokens_used 字段记录了本次对话消耗的 Token 数量
6. **对话历史**：每次对话都会自动保存到历史记录中
7. **重试机制**：系统内置重试机制，如果 AI 服务暂时不可用会自动重试
8. **响应格式**：message 和 response 字段内容相同，提供两个字段是为了兼容不同的前端实现
9. **AI 配置**：需要在设置中配置 AI 服务才能使用此接口
10. **响应时间**：AI 响应可能需要几秒钟，建议前端显示加载状态

---

### AI 生成餐饮建议

**接口**: `POST /api/v1/ai/suggest`

**说明**: 使用 AI 生成未来几天的餐饮建议。系统会根据用户的食材库、营养目标和个人偏好，智能推荐合理的餐饮搭配。与生成饮食计划不同，此接口返回的是建议性的餐饮方案，不会直接创建计划记录。

**认证**: 是

#### 请求参数

##### 请求体

```json
{
  "days": 3,
  "target_calories": 2000
}
```

| 字段 | 类型 | 必填 | 说明 | 验证规则 |
|------|------|------|------|----------|
| days | int | 是 | 生成建议的天数 | 最小 1，最大 30 |
| target_calories | int | 否 | 目标热量（千卡/天） | 最小 800，最大 10000 |

#### 请求示例

```bash
# 生成 3 天的餐饮建议
curl -X POST http://localhost:9090/api/v1/ai/suggest \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "days": 3,
    "target_calories": 2000
  }'

# 只生成 1 天的建议
curl -X POST http://localhost:9090/api/v1/ai/suggest \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "days": 1
  }'
```


#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "days": [
      {
        "date": "2024-11-17",
        "meals": [
          {
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
              },
              {
                "food_id": 25,
                "name": "牛奶",
                "amount": 250,
                "unit": "ml"
              }
            ],
            "nutrition": {
              "protein": 25.5,
              "carbs": 35.3,
              "fat": 12.2,
              "fiber": 4.5,
              "calories": 350.0
            },
            "reasoning": "早餐提供优质蛋白质和复合碳水化合物，能够提供持久的能量，适合开启新的一天。"
          },
          {
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
                "amount": 100,
                "unit": "g"
              }
            ],
            "nutrition": {
              "protein": 55.3,
              "carbs": 45.8,
              "fat": 5.5,
              "fiber": 7.2,
              "calories": 450.0
            },
            "reasoning": "午餐以高蛋白鸡胸肉为主，搭配蔬菜和适量主食，营养均衡且热量适中。"
          },
          {
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
              },
              {
                "food_id": 12,
                "name": "红薯",
                "amount": 100,
                "unit": "g"
              }
            ],
            "nutrition": {
              "protein": 48.2,
              "carbs": 28.5,
              "fat": 15.3,
              "fiber": 6.8,
              "calories": 450.0
            },
            "reasoning": "晚餐选择牛肉提供优质蛋白质和铁质，搭配蔬菜和红薯，营养丰富且易消化。"
          }
        ],
        "daily_nutrition": {
          "protein": 129.0,
          "carbs": 109.6,
          "fat": 33.0,
          "fiber": 18.5,
          "calories": 1250.0
        }
      },
      {
        "date": "2024-11-18",
        "meals": [
          {
            "meal_type": "breakfast",
            "foods": [
              {
                "food_id": 30,
                "name": "燕麦片",
                "amount": 50,
                "unit": "g"
              },
              {
                "food_id": 25,
                "name": "牛奶",
                "amount": 250,
                "unit": "ml"
              },
              {
                "food_id": 35,
                "name": "香蕉",
                "amount": 1,
                "unit": "根"
              }
            ],
            "nutrition": {
              "protein": 18.5,
              "carbs": 55.3,
              "fat": 8.2,
              "fiber": 6.5,
              "calories": 360.0
            },
            "reasoning": "燕麦早餐提供丰富的膳食纤维和缓释能量，搭配香蕉补充钾元素。"
          }
        ],
        "daily_nutrition": {
          "protein": 18.5,
          "carbs": 55.3,
          "fat": 8.2,
          "fiber": 6.5,
          "calories": 360.0
        }
      }
    ],
    "total_nutrition": {
      "protein": 147.5,
      "carbs": 164.9,
      "fat": 41.2,
      "fiber": 25.0,
      "calories": 1610.0
    },
    "suggestions": [
      "建议每天保持充足的水分摄入（至少 2 升）",
      "可以在两餐之间添加健康零食，如坚果或水果",
      "注意烹饪方式，建议采用蒸、煮、烤等低油烹饪方法"
    ]
  },
  "timestamp": 1699999999
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| days | array | 每天的餐饮建议列表 |
| days[].date | string | 日期（YYYY-MM-DD 格式） |
| days[].meals | array | 当天的餐次列表 |
| days[].meals[].meal_type | string | 餐次类型（breakfast, lunch, dinner, snack） |
| days[].meals[].foods | array | 食材列表 |
| days[].meals[].nutrition | object | 该餐次的营养数据 |
| days[].meals[].reasoning | string | AI 推荐理由 |
| days[].daily_nutrition | object | 当天的总营养数据 |
| total_nutrition | object | 所有天数的总营养数据 |
| suggestions | array | AI 的额外建议 |


**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid request parameters: Key: 'SuggestMealPlanRequest.Days' Error:Field validation for 'Days' failed on the 'max' tag",
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
  "error": "failed to generate meal plan",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | days 超出范围（1-30）、target_calories 超出范围（800-10000） |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | AI 服务调用失败、食材库为空、AI 设置未配置 |

#### 注意事项

1. **天数范围**：可以生成 1-30 天的建议，建议生成 3-7 天
2. **热量目标**：target_calories 是可选的，如果不提供，AI 会根据用户的偏好设置自动确定
3. **食材来源**：AI 只会使用用户食材库中的食材生成建议
4. **建议性质**：此接口返回的是建议，不会创建实际的计划记录
5. **营养计算**：系统会自动计算每餐、每天和总的营养数据
6. **额外建议**：suggestions 字段包含 AI 的额外饮食建议
7. **与计划的区别**：此接口用于预览和参考，如需创建实际计划，请使用饮食计划模块的生成接口
8. **重试机制**：系统内置重试机制，如果 AI 服务暂时不可用会自动重试
9. **响应时间**：生成建议可能需要较长时间（5-15 秒），建议前端显示加载状态
10. **食材不足**：如果食材库中的食材不足以生成指定天数的建议，会返回错误

---

### 获取对话历史

**接口**: `GET /api/v1/ai/history`

**说明**: 获取用户与 AI 助手的对话历史记录，支持分页查询。对话历史按时间倒序排列（最新的在前）。

**认证**: 是

#### 请求参数

##### 查询参数

| 参数 | 类型 | 必填 | 说明 | 默认值 | 示例 |
|------|------|------|------|--------|------|
| page | int | 否 | 页码（从 1 开始） | 1 | 1, 2, 3 |
| page_size | int | 否 | 每页数据量 | 20 | 10, 20, 50（最大 100） |

#### 请求示例

```bash
# 获取第一页对话历史（默认每页 20 条）
curl -X GET "http://localhost:9090/api/v1/ai/history" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取第 2 页，每页 50 条
curl -X GET "http://localhost:9090/api/v1/ai/history?page=2&page_size=50" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取最近 10 条对话
curl -X GET "http://localhost:9090/api/v1/ai/history?page=1&page_size=10" \
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
      "id": 123,
      "user_id": 1,
      "user_input": "我今天应该吃什么？",
      "ai_response": "根据您的营养目标，我建议今天可以选择：\n\n早餐：鸡蛋 2 个 + 全麦面包 50g + 牛奶 250ml\n午餐：鸡胸肉 200g + 西兰花 150g + 糙米饭 100g\n晚餐：牛肉 150g + 菠菜 200g + 红薯 100g\n\n这样的搭配营养均衡，符合您的目标。",
      "context": {
        "current_date": "2024-11-17",
        "target_calories": "2000"
      },
      "tokens_used": 256,
      "created_at": "2024-11-17T10:30:00Z"
    },
    {
      "id": 122,
      "user_id": 1,
      "user_input": "鸡胸肉和牛肉哪个更适合减脂？",
      "ai_response": "从减脂角度来看，鸡胸肉更适合：\n\n1. 热量更低：鸡胸肉每 100g 约 165 千卡，牛肉约 250 千卡\n2. 脂肪含量更低：鸡胸肉脂肪含量约 3.6g/100g，牛肉约 15g/100g\n3. 蛋白质含量相近：两者都是优质蛋白质来源\n\n但牛肉含有更多的铁和维生素 B12，建议交替食用，保持营养均衡。",
      "context": {},
      "tokens_used": 198,
      "created_at": "2024-11-17T09:15:00Z"
    },
    {
      "id": 121,
      "user_id": 1,
      "user_input": "如何计算每日所需热量？",
      "ai_response": "每日所需热量（TDEE）的计算方法：\n\n1. 计算基础代谢率（BMR）：\n   - 男性：BMR = 10 × 体重(kg) + 6.25 × 身高(cm) - 5 × 年龄 + 5\n   - 女性：BMR = 10 × 体重(kg) + 6.25 × 身高(cm) - 5 × 年龄 - 161\n\n2. 根据活动水平调整：\n   - 久坐：BMR × 1.2\n   - 轻度活动：BMR × 1.375\n   - 中度活动：BMR × 1.55\n   - 高度活动：BMR × 1.725\n\n3. 根据目标调整：\n   - 减脂：TDEE - 300~500 千卡\n   - 增肌：TDEE + 300~500 千卡\n   - 维持：TDEE\n\n建议在设置中配置您的个人信息，系统会自动计算。",
      "context": {},
      "tokens_used": 312,
      "created_at": "2024-11-16T15:20:00Z"
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

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 对话记录唯一标识符 |
| user_id | int64 | 所属用户 ID |
| user_input | string | 用户的输入消息 |
| ai_response | string | AI 的回复内容 |
| context | object | 对话时传递的上下文信息 |
| tokens_used | int | 本次对话使用的 Token 数量 |
| created_at | string | 对话时间（ISO 8601 格式） |


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
  "error": "invalid page parameter",
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
  "error": "failed to get chat history",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | page 或 page_size 参数无效 |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | 服务器内部错误 |

#### 注意事项

1. **排序方式**：对话历史按时间倒序排列，最新的对话在最前面
2. **分页默认值**：不指定分页参数时，默认返回第 1 页，每页 20 条数据
3. **分页限制**：page_size 最大值为 100，超过会自动调整为 100
4. **空结果**：当没有对话历史时，返回空数组，不会报错
5. **上下文信息**：context 字段保存了对话时传递的上下文信息
6. **Token 统计**：tokens_used 字段记录了每次对话消耗的 Token 数量
7. **隐私保护**：只能查询属于当前用户的对话历史
8. **历史保留**：对话历史会永久保存，除非用户主动删除（当前版本不支持删除）

---

## 数据模型

### ChatHistory 模型

对话历史数据模型：

```typescript
interface ChatHistory {
  id: number;              // 对话记录 ID
  user_id: number;         // 用户 ID
  user_input: string;      // 用户输入
  ai_response: string;     // AI 响应
  context: object;         // 上下文信息（JSON 对象）
  tokens_used: number;     // Token 使用量
  created_at: string;      // 创建时间（ISO 8601 格式）
}
```

### MealPlanSuggestion 模型

餐饮建议数据模型：

```typescript
interface MealPlanSuggestion {
  days: DaySuggestion[];           // 每天的建议
  total_nutrition: NutritionData;  // 总营养数据
  suggestions: string[];           // 额外建议
}

interface DaySuggestion {
  date: string;                    // 日期（YYYY-MM-DD）
  meals: MealSuggestion[];         // 餐次列表
  daily_nutrition: NutritionData;  // 当天总营养
}

interface MealSuggestion {
  meal_type: string;               // 餐次类型
  foods: FoodItem[];               // 食材列表
  nutrition: NutritionData;        // 营养数据
  reasoning: string;               // 推荐理由
}

interface FoodItem {
  food_id: number;                 // 食材 ID
  name: string;                    // 食材名称
  amount: number;                  // 用量
  unit: string;                    // 单位
}

interface NutritionData {
  protein: number;                 // 蛋白质（克）
  carbs: number;                   // 碳水化合物（克）
  fat: number;                     // 脂肪（克）
  fiber: number;                   // 纤维（克）
  calories: number;                // 热量（千卡）
}
```

---

## 使用场景

### 场景 1：询问饮食建议

用户想知道今天应该吃什么：

```bash
curl -X POST http://localhost:9090/api/v1/ai/chat \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "我今天应该吃什么？",
    "context": {
      "current_date": "2024-11-17",
      "target_calories": "2000"
    }
  }'
```

### 场景 2：询问营养知识

用户想了解某种食材的营养价值：

```bash
curl -X POST http://localhost:9090/api/v1/ai/chat \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "鸡胸肉有什么营养价值？适合减脂吗？"
  }'
```

### 场景 3：生成餐饮建议

用户想预览未来几天的餐饮方案：

```bash
curl -X POST http://localhost:9090/api/v1/ai/suggest \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "days": 3,
    "target_calories": 2000
  }'
```

### 场景 4：查看对话历史

用户想回顾之前与 AI 的对话：

```bash
curl -X GET "http://localhost:9090/api/v1/ai/history?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 场景 5：带上下文的智能对话

用户想让 AI 根据当前情况提供建议：

```bash
curl -X POST http://localhost:9090/api/v1/ai/chat \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "我今天已经摄入了 1200 千卡，晚餐应该吃什么？",
    "context": {
      "current_date": "2024-11-17",
      "target_calories": "2000",
      "consumed_calories": "1200",
      "remaining_calories": "800"
    }
  }'
```

---

## 最佳实践

### 1. 对话交互

- **清晰提问**：提问时尽量清晰具体，避免模糊的问题
- **提供上下文**：传递相关的上下文信息，让 AI 更好地理解你的情况
- **分步询问**：复杂问题可以分步询问，逐步深入
- **参考历史**：查看对话历史，避免重复询问相同的问题

### 2. 上下文使用

- **常用上下文**：current_date、target_calories、target_protein、consumed_calories 等
- **结构化数据**：上下文信息使用键值对格式，便于 AI 理解
- **相关信息**：只传递与问题相关的上下文，避免信息过载
- **数据格式**：上下文值使用字符串格式

### 3. 餐饮建议

- **合理天数**：建议生成 3-7 天的建议，不要一次生成太多
- **设置目标**：提供 target_calories 参数，让建议更符合你的需求
- **参考性质**：建议仅供参考，可以根据实际情况调整
- **转换计划**：如果建议合适，可以使用饮食计划模块创建实际计划

### 4. 历史管理

- **定期回顾**：定期查看对话历史，回顾 AI 的建议
- **学习知识**：从历史对话中学习营养知识和饮食技巧
- **分页查询**：使用合理的 page_size，避免一次加载过多数据
- **Token 统计**：关注 tokens_used，了解 AI 使用成本

### 5. AI 配置

- **配置检查**：使用前确保在设置中配置了 AI 服务
- **提供商选择**：根据需求选择合适的 AI 提供商（OpenAI、DeepSeek 等）
- **模型选择**：选择合适的模型，平衡性能和成本
- **连接测试**：定期测试 AI 连接，确保服务可用

### 6. 错误处理

- **重试机制**：系统内置重试，但如果持续失败，检查 AI 配置
- **超时处理**：AI 响应可能较慢，前端应设置合理的超时时间
- **错误提示**：向用户清晰地展示错误信息
- **降级方案**：AI 服务不可用时，提供备用方案

### 7. 性能优化

- **加载状态**：AI 响应需要时间，前端应显示加载状态
- **缓存策略**：对话历史可以在客户端缓存
- **分页加载**：历史记录使用分页加载，避免一次加载过多
- **异步处理**：AI 请求使用异步处理，不阻塞用户操作

---

## 常见问题

### Q: AI 对话和餐饮建议有什么区别？

A: 
- **AI 对话**：自由的自然语言交互，可以询问任何饮食相关的问题，获取知识和建议
- **餐饮建议**：结构化的餐饮方案生成，返回具体的食材搭配和营养数据
- **使用场景**：对话用于咨询和学习，建议用于规划餐饮
- **数据格式**：对话返回文本，建议返回结构化数据

### Q: 如何配置 AI 服务？

A: 
- 在设置管理模块中配置 AI 服务
- 需要提供：AI 提供商、API 端点、API 密钥、模型名称等
- 配置后可以使用测试接口验证连接
- 详见 [设置管理模块文档](./08-settings.md)

### Q: 支持哪些 AI 提供商？

A: 
- **OpenAI**：GPT-3.5、GPT-4 等模型
- **DeepSeek**：DeepSeek Chat 模型
- **自定义**：支持兼容 OpenAI API 格式的自定义提供商
- 可以在设置中切换不同的提供商

### Q: Token 使用量如何计算？

A: 
- Token 使用量由 AI 提供商计算
- 包括用户输入和 AI 响应的 Token 数量
- 不同模型的 Token 计费标准不同
- 可以在对话历史中查看每次对话的 Token 使用量

### Q: 对话历史会保存多久？

A: 
- 对话历史会永久保存
- 当前版本不支持删除对话历史
- 建议定期导出重要的对话记录
- 未来版本可能会添加删除和导出功能

### Q: 如何让 AI 提供更准确的建议？

A: 
- **提供上下文**：传递相关的上下文信息
- **清晰提问**：问题尽量具体明确
- **完善资料**：在设置中完善个人信息和营养目标
- **丰富食材库**：食材库中的食材越多，AI 的建议越丰富

### Q: AI 响应很慢怎么办？

A: 
- AI 响应时间取决于提供商和网络状况
- 通常需要 2-10 秒
- 如果持续很慢，检查网络连接和 AI 配置
- 可以尝试切换到其他 AI 提供商

### Q: AI 建议不合理怎么办？

A: 
- AI 建议仅供参考，可以根据实际情况调整
- 可以在对话中向 AI 反馈，获取改进建议
- 可以在设置中调整偏好设置
- 可以手动修改建议的内容

### Q: 餐饮建议和饮食计划有什么区别？

A: 
- **餐饮建议**：AI 服务模块提供的预览性建议，不会创建实际记录
- **饮食计划**：饮食计划模块创建的实际计划记录，可以执行和跟踪
- **使用流程**：先用建议预览，满意后再创建计划
- **数据独立**：建议和计划是独立的，互不影响

### Q: 可以同时使用多个 AI 提供商吗？

A: 
- 当前版本只能配置一个活跃的 AI 提供商
- 可以在设置中切换不同的提供商
- 切换后新的对话会使用新的提供商
- 历史对话记录会保留使用的提供商信息

### Q: AI 服务需要额外付费吗？

A: 
- 系统本身不收费
- 但 AI 提供商（如 OpenAI）可能需要付费
- 需要自己申请 AI 提供商的 API 密钥
- Token 使用量会影响 AI 提供商的费用

### Q: 如何控制 AI 使用成本？

A: 
- **选择模型**：使用成本较低的模型（如 GPT-3.5）
- **精简提问**：避免过长的问题和上下文
- **合理使用**：只在需要时使用 AI 服务
- **监控使用**：定期查看 Token 使用量

---

## 相关文档

- [数据模型](./data-models.md) - 查看 ChatHistory 模型的完整定义
- [饮食计划模块](./04-plans.md) - 了解如何创建和管理饮食计划
- [设置管理模块](./08-settings.md) - 了解如何配置 AI 服务
- [通用概念](./common-concepts.md) - 了解认证、分页等通用概念
- [错误码说明](./error-codes.md) - 查看所有错误码的详细说明
- [API 文档总览](./README.md) - 返回 API 文档首页

