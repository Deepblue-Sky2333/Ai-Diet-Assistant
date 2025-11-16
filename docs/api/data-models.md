# 数据模型

本文档定义了 AI Diet Assistant 系统中使用的所有核心数据模型。

## 目录

- [Food (食材)](#food-食材)
- [Meal (餐饮记录)](#meal-餐饮记录)
- [Plan (饮食计划)](#plan-饮食计划)
- [NutritionData (营养数据)](#nutritiondata-营养数据)
- [UserPreferences (用户偏好)](#userpreferences-用户偏好)
- [AISettings (AI 设置)](#aisettings-ai-设置)
- [辅助模型](#辅助模型)

---

## Food (食材)

食材模型表示用户市场面板中的食材项目。

### 字段定义

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | integer | 食材唯一标识符 | 主键，自动生成 |
| user_id | integer | 所属用户 ID | 必填，外键 |
| name | string | 食材名称 | 必填，最大 100 字符 |
| category | string | 食材分类 | 必填，枚举值：meat, vegetable, fruit, grain, other |
| price | number | 价格 | 必填，≥ 0 |
| unit | string | 单位 | 必填，最大 20 字符 |
| protein | number | 蛋白质含量（克/单位） | 必填，≥ 0 |
| carbs | number | 碳水化合物含量（克/单位） | 必填，≥ 0 |
| fat | number | 脂肪含量（克/单位） | 必填，≥ 0 |
| fiber | number | 纤维含量（克/单位） | 必填，≥ 0 |
| calories | number | 热量（千卡/单位） | 必填，≥ 0 |
| available | boolean | 是否可用 | 默认 true |
| created_at | string | 创建时间 | ISO 8601 格式 |
| updated_at | string | 更新时间 | ISO 8601 格式 |

### TypeScript 接口

```typescript
interface Food {
  id: number;
  user_id: number;
  name: string;
  category: 'meat' | 'vegetable' | 'fruit' | 'grain' | 'other';
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

### 示例数据

```json
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
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

---

## Meal (餐饮记录)

餐饮记录模型表示用户的实际用餐记录。

### 字段定义

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | integer | 记录唯一标识符 | 主键，自动生成 |
| user_id | integer | 所属用户 ID | 必填，外键 |
| meal_date | string | 用餐日期 | 必填，ISO 8601 格式 |
| meal_type | string | 餐次类型 | 必填，枚举值：breakfast, lunch, dinner, snack |
| foods | array | 食材列表 | 必填，至少包含 1 项，参见 MealFood |
| nutrition | object | 营养汇总 | 自动计算，参见 NutritionData |
| notes | string | 备注 | 可选，最大 500 字符 |
| created_at | string | 创建时间 | ISO 8601 格式 |
| updated_at | string | 更新时间 | ISO 8601 格式 |

### TypeScript 接口

```typescript
interface Meal {
  id: number;
  user_id: number;
  meal_date: string;
  meal_type: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  foods: MealFood[];
  nutrition: NutritionData;
  notes?: string;
  created_at: string;
  updated_at: string;
}

interface MealFood {
  food_id: number;
  name: string;
  amount: number;
  unit: string;
}
```

### 示例数据

```json
{
  "id": 1,
  "user_id": 1,
  "meal_date": "2024-01-15T12:00:00Z",
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
    "protein": 40.5,
    "carbs": 12.0,
    "fat": 3.5,
    "fiber": 5.2,
    "calories": 245.0
  },
  "notes": "午餐，健身后",
  "created_at": "2024-01-15T12:30:00Z",
  "updated_at": "2024-01-15T12:30:00Z"
}
```

---

## Plan (饮食计划)

饮食计划模型表示未来日期的用餐计划。

### 字段定义

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | integer | 计划唯一标识符 | 主键，自动生成 |
| user_id | integer | 所属用户 ID | 必填，外键 |
| plan_date | string | 计划日期 | 必填，ISO 8601 格式 |
| meal_type | string | 餐次类型 | 必填，枚举值：breakfast, lunch, dinner, snack |
| foods | array | 食材列表 | 必填，至少包含 1 项，参见 MealFood |
| nutrition | object | 营养汇总 | 自动计算，参见 NutritionData |
| status | string | 计划状态 | 枚举值：pending, completed, skipped，默认 pending |
| ai_reasoning | string | AI 推荐理由 | 可选，最大 1000 字符 |
| created_at | string | 创建时间 | ISO 8601 格式 |
| updated_at | string | 更新时间 | ISO 8601 格式 |

### TypeScript 接口

```typescript
interface Plan {
  id: number;
  user_id: number;
  plan_date: string;
  meal_type: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  foods: MealFood[];
  nutrition: NutritionData;
  status: 'pending' | 'completed' | 'skipped';
  ai_reasoning?: string;
  created_at: string;
  updated_at: string;
}
```

### 示例数据

```json
{
  "id": 1,
  "user_id": 1,
  "plan_date": "2024-01-16T08:00:00Z",
  "meal_type": "breakfast",
  "foods": [
    {
      "food_id": 5,
      "name": "燕麦",
      "amount": 50,
      "unit": "g"
    },
    {
      "food_id": 6,
      "name": "牛奶",
      "amount": 250,
      "unit": "ml"
    },
    {
      "food_id": 7,
      "name": "香蕉",
      "amount": 1,
      "unit": "个"
    }
  ],
  "nutrition": {
    "protein": 15.0,
    "carbs": 65.0,
    "fat": 8.0,
    "fiber": 6.0,
    "calories": 380.0
  },
  "status": "pending",
  "ai_reasoning": "根据您的健身目标，早餐需要充足的碳水化合物和蛋白质。燕麦提供复合碳水，牛奶提供优质蛋白，香蕉补充快速能量和钾元素。",
  "created_at": "2024-01-15T20:00:00Z",
  "updated_at": "2024-01-15T20:00:00Z"
}
```

---

## NutritionData (营养数据)

营养数据模型表示食物或餐饮的营养信息汇总。

### 字段定义

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| protein | number | 蛋白质（克） | ≥ 0 |
| carbs | number | 碳水化合物（克） | ≥ 0 |
| fat | number | 脂肪（克） | ≥ 0 |
| fiber | number | 纤维（克） | ≥ 0 |
| calories | number | 热量（千卡） | ≥ 0 |

### TypeScript 接口

```typescript
interface NutritionData {
  protein: number;
  carbs: number;
  fat: number;
  fiber: number;
  calories: number;
}
```

### 示例数据

```json
{
  "protein": 25.5,
  "carbs": 45.0,
  "fat": 12.0,
  "fiber": 8.5,
  "calories": 385.0
}
```

---

## UserPreferences (用户偏好)

用户偏好模型存储用户的饮食偏好和营养目标。

### 字段定义

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | integer | 偏好设置唯一标识符 | 主键，自动生成 |
| user_id | integer | 所属用户 ID | 必填，外键 |
| taste_preferences | string | 口味偏好 | 可选，最大 500 字符 |
| dietary_restrictions | string | 饮食限制 | 可选，最大 500 字符 |
| daily_calories_goal | integer | 每日热量目标（千卡） | 800-10000 |
| daily_protein_goal | integer | 每日蛋白质目标（克） | 0-500 |
| daily_carbs_goal | integer | 每日碳水化合物目标（克） | 0-1000 |
| daily_fat_goal | integer | 每日脂肪目标（克） | 0-500 |
| daily_fiber_goal | integer | 每日纤维目标（克） | 0-200 |
| created_at | string | 创建时间 | ISO 8601 格式 |
| updated_at | string | 更新时间 | ISO 8601 格式 |

### TypeScript 接口

```typescript
interface UserPreferences {
  id: number;
  user_id: number;
  taste_preferences: string;
  dietary_restrictions: string;
  daily_calories_goal: number;
  daily_protein_goal: number;
  daily_carbs_goal: number;
  daily_fat_goal: number;
  daily_fiber_goal: number;
  created_at: string;
  updated_at: string;
}
```

### 示例数据

```json
{
  "id": 1,
  "user_id": 1,
  "taste_preferences": "喜欢清淡口味，偏好蒸煮烹饪方式",
  "dietary_restrictions": "不吃海鲜，对花生过敏",
  "daily_calories_goal": 2000,
  "daily_protein_goal": 120,
  "daily_carbs_goal": 250,
  "daily_fat_goal": 60,
  "daily_fiber_goal": 30,
  "created_at": "2024-01-10T08:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

---

## AISettings (AI 设置)

AI 设置模型存储用户的 AI 服务提供商配置。

### 字段定义

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | integer | 设置唯一标识符 | 主键，自动生成 |
| user_id | integer | 所属用户 ID | 必填，外键 |
| provider | string | AI 提供商 | 必填，枚举值：openai, deepseek, custom |
| api_endpoint | string | API 端点 URL | 可选，必须是有效 URL |
| api_key | string | API 密钥 | 必填，存储时加密 |
| model | string | 模型名称 | 必填 |
| temperature | number | 温度参数 | 0.0-2.0，控制输出随机性 |
| max_tokens | integer | 最大 Token 数 | 1-4096 |
| is_active | boolean | 是否激活 | 默认 true |
| created_at | string | 创建时间 | ISO 8601 格式 |
| updated_at | string | 更新时间 | ISO 8601 格式 |

### TypeScript 接口

```typescript
interface AISettings {
  id: number;
  user_id: number;
  provider: 'openai' | 'deepseek' | 'custom';
  api_endpoint?: string;
  api_key: string;
  model: string;
  temperature: number;
  max_tokens: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
```

### 示例数据

```json
{
  "id": 1,
  "user_id": 1,
  "provider": "openai",
  "api_endpoint": "https://api.openai.com/v1",
  "api_key": "sk-****",
  "model": "gpt-4",
  "temperature": 0.7,
  "max_tokens": 2048,
  "is_active": true,
  "created_at": "2024-01-10T08:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

**注意**: 
- API 密钥在数据库中加密存储
- 返回给客户端时，API 密钥会被掩码处理（仅显示前 4 位和后 4 位）

---

## 辅助模型

### MealFood (餐饮食材)

表示餐饮记录或计划中的单个食材项。

```typescript
interface MealFood {
  food_id: number;      // 食材 ID
  name: string;         // 食材名称
  amount: number;       // 数量，> 0，≤ 10000
  unit: string;         // 单位，1-20 字符
}
```

### FoodFilter (食材过滤器)

用于食材列表查询的过滤条件。

```typescript
interface FoodFilter {
  category?: string;    // 分类过滤
  available?: boolean;  // 可用性过滤
  page?: number;        // 页码
  page_size?: number;   // 每页数量
}
```

### MealFilter (餐饮记录过滤器)

用于餐饮记录列表查询的过滤条件。

```typescript
interface MealFilter {
  start_date?: string;  // 开始日期
  end_date?: string;    // 结束日期
  meal_type?: string;   // 餐次类型
  page?: number;        // 页码
  page_size?: number;   // 每页数量
}
```

### PlanFilter (计划过滤器)

用于饮食计划列表查询的过滤条件。

```typescript
interface PlanFilter {
  start_date?: string;  // 开始日期
  end_date?: string;    // 结束日期
  status?: string;      // 计划状态
  page?: number;        // 页码
  page_size?: number;   // 每页数量
}
```

### DailyNutritionStats (每日营养统计)

表示某一天的营养统计数据。

```typescript
interface DailyNutritionStats {
  date: string;              // 日期
  nutrition: NutritionData;  // 营养汇总
  meal_count: number;        // 餐次数量
}
```

### MonthlyStats (月度统计)

表示某个月的营养统计数据。

```typescript
interface MonthlyStats {
  year: number;                        // 年份
  month: number;                       // 月份
  total_meals: number;                 // 总餐次数
  daily_stats: DailyNutritionStats[];  // 每日统计
  avg_daily: NutritionData;            // 日均营养
  total: NutritionData;                // 月度总营养
}
```

### NutritionComparison (营养对比)

表示实际营养与目标营养的对比。

```typescript
interface NutritionComparison {
  target: NutritionData;              // 目标营养
  actual: NutritionData;              // 实际营养
  difference: NutritionData;          // 差值
  percentage: {                       // 完成百分比
    [key: string]: number;
  };
}
```

### DashboardData (仪表盘数据)

表示仪表盘聚合数据。

```typescript
interface DashboardData {
  monthly_stats: MonthlyStats;         // 月度统计
  future_plans: Plan[];                // 未来计划
  today_stats: DailyNutritionStats;    // 今日统计
  current_month: number;               // 当前月份
  current_year: number;                // 当前年份
  generated_at: string;                // 生成时间
}
```

### BatchResult (批量操作结果)

表示批量导入操作的结果。

```typescript
interface BatchResult {
  success: number;      // 成功数量
  failed: number;       // 失败数量
  errors?: string[];    // 错误信息列表
}
```

### ChatHistory (对话历史)

表示 AI 对话历史记录。

```typescript
interface ChatHistory {
  id: number;           // 记录 ID
  user_id: number;      // 用户 ID
  user_input: string;   // 用户输入
  ai_response: string;  // AI 响应
  context?: string;     // 上下文（JSON 字符串）
  tokens_used: number;  // 使用的 Token 数
  created_at: string;   // 创建时间
}
```

---

## 数据关系

### 关系图

```
User (用户)
  ├── Food (食材) [1:N]
  ├── Meal (餐饮记录) [1:N]
  │     └── MealFood (餐饮食材) [1:N]
  ├── Plan (饮食计划) [1:N]
  │     └── MealFood (计划食材) [1:N]
  ├── UserPreferences (用户偏好) [1:1]
  ├── AISettings (AI 设置) [1:N]
  └── ChatHistory (对话历史) [1:N]
```

### 关系说明

1. **User → Food**: 一个用户可以创建多个食材
2. **User → Meal**: 一个用户可以创建多个餐饮记录
3. **User → Plan**: 一个用户可以创建多个饮食计划
4. **User → UserPreferences**: 一个用户有一个偏好设置
5. **User → AISettings**: 一个用户可以配置多个 AI 设置（但只有一个激活）
6. **User → ChatHistory**: 一个用户可以有多条对话历史
7. **Meal/Plan → MealFood**: 一个餐饮记录或计划包含多个食材项
8. **MealFood → Food**: 每个食材项引用一个食材

---

## 验证规则

### 通用规则

1. **日期格式**: 所有日期时间字段使用 ISO 8601 格式（如 `2024-01-15T10:30:00Z`）
2. **字符串长度**: 所有字符串字段都有最大长度限制
3. **数值范围**: 所有数值字段都有合理的范围限制
4. **枚举值**: 枚举类型字段只接受预定义的值

### 特定规则

1. **Food.category**: 必须是 `meat`, `vegetable`, `fruit`, `grain`, `other` 之一
2. **Meal.meal_type / Plan.meal_type**: 必须是 `breakfast`, `lunch`, `dinner`, `snack` 之一
3. **Plan.status**: 必须是 `pending`, `completed`, `skipped` 之一
4. **AISettings.provider**: 必须是 `openai`, `deepseek`, `custom` 之一
5. **MealFood.amount**: 必须 > 0 且 ≤ 10000
6. **UserPreferences.daily_calories_goal**: 必须在 800-10000 之间
7. **AISettings.temperature**: 必须在 0.0-2.0 之间
8. **AISettings.max_tokens**: 必须在 1-4096 之间

---

## 计算规则

### 营养数据计算

餐饮记录和计划的营养数据通过以下公式计算：

```
总营养 = Σ (食材营养 × 数量 / 食材单位量)
```

例如：
- 食材：鸡胸肉，蛋白质 23g/100g
- 用量：150g
- 计算：23 × 150 / 100 = 34.5g 蛋白质

### 营养对比计算

```typescript
// 差值
difference = actual - target

// 完成百分比
percentage = (actual / target) × 100
```

---

## 注意事项

1. **时区处理**: 所有时间戳使用 UTC 时区，客户端需要根据用户时区进行转换
2. **精度**: 营养数据保留 1 位小数，价格保留 2 位小数
3. **安全性**: API 密钥在数据库中加密存储，返回时进行掩码处理
4. **软删除**: 系统不使用软删除，删除操作为物理删除
5. **并发**: 更新操作使用乐观锁（通过 updated_at 字段）
6. **默认值**: 布尔字段默认值在数据库层面设置
