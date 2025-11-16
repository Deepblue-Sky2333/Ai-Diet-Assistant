# Dashboard 模块

## 概述

Dashboard 模块提供聚合的仪表盘数据，为用户展示饮食管理的全局视图。通过单个接口调用，用户可以获取今日营养摄入、营养目标、未来饮食计划等关键信息，无需多次调用不同的接口。这个模块是应用首页的核心数据源，帮助用户快速了解当前的饮食状况和即将到来的计划。

**核心功能**：
- 获取今日营养摄入统计
- 获取用户设置的营养目标
- 获取未来 2 天的饮食计划

**数据特性**：
- 聚合多个模块的数据
- 一次请求获取所有关键信息
- 实时计算和更新
- 自动处理默认值

**应用场景**：
- 应用首页数据展示
- 用户每日饮食概览
- 快速了解营养达标情况
- 查看即将到来的饮食计划

---

## 接口列表

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/dashboard` | 获取 Dashboard 数据 | 是 |

---

## 接口详情


### 获取 Dashboard 数据

**接口**: `GET /api/v1/dashboard`

**说明**: 获取聚合的仪表盘数据，包括今日营养摄入、营养目标和未来饮食计划。这是一个聚合接口，将多个模块的数据整合在一起，方便前端一次性获取首页所需的所有关键信息。系统会自动计算今日的营养摄入，从用户偏好中读取营养目标（如未设置则使用默认值），并获取未来 2 天的待执行饮食计划。

**认证**: 是


#### 请求参数

无需任何参数，系统会自动获取当前用户的数据。

#### 请求示例

```bash
# 获取 Dashboard 数据
curl -X GET "http://localhost:9090/api/v1/dashboard" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 使用环境变量存储 Token
export TOKEN="YOUR_ACCESS_TOKEN"
curl -X GET "http://localhost:9090/api/v1/dashboard" \
  -H "Authorization: Bearer $TOKEN"
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "today_nutrition": {
      "calories": 1850.5,
      "protein": 142.3,
      "carbs": 215.8,
      "fat": 62.4
    },
    "nutrition_goal": {
      "calories": 2000,
      "protein": 150,
      "carbs": 250,
      "fat": 70
    },
    "upcoming_plans": [
      {
        "id": 101,
        "date": "2024-11-18",
        "meal_type": "breakfast",
        "reason": "高蛋白早餐，提供充足能量开始新的一天"
      },
      {
        "id": 102,
        "date": "2024-11-18",
        "meal_type": "lunch",
        "reason": "均衡午餐，包含优质蛋白质和复合碳水化合物"
      },
      {
        "id": 103,
        "date": "2024-11-19",
        "meal_type": "breakfast",
        "reason": "轻食早餐，适合工作日快速准备"
      }
    ]
  },
  "timestamp": 1699999999
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| today_nutrition | object | 今日营养摄入统计 |
| today_nutrition.calories | float | 今日摄入热量（千卡） |
| today_nutrition.protein | float | 今日摄入蛋白质（克） |
| today_nutrition.carbs | float | 今日摄入碳水化合物（克） |
| today_nutrition.fat | float | 今日摄入脂肪（克） |
| nutrition_goal | object | 营养目标（来自用户偏好设置） |
| nutrition_goal.calories | int | 目标热量（千卡/天） |
| nutrition_goal.protein | int | 目标蛋白质（克/天） |
| nutrition_goal.carbs | int | 目标碳水化合物（克/天） |
| nutrition_goal.fat | int | 目标脂肪（克/天） |
| upcoming_plans | array | 未来 2 天的待执行饮食计划 |
| upcoming_plans[].id | int | 计划 ID |
| upcoming_plans[].date | string | 计划日期（YYYY-MM-DD 格式） |
| upcoming_plans[].meal_type | string | 餐次类型（breakfast/lunch/dinner/snack） |
| upcoming_plans[].reason | string | AI 生成的计划理由 |


**无今日数据响应 (200)**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "today_nutrition": {
      "calories": 0,
      "protein": 0,
      "carbs": 0,
      "fat": 0
    },
    "nutrition_goal": {
      "calories": 2000,
      "protein": 150,
      "carbs": 250,
      "fat": 70
    },
    "upcoming_plans": []
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

**错误响应 (500)**:

```json
{
  "code": 50001,
  "message": "internal error",
  "error": "failed to get dashboard data",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | 服务器内部错误、数据库查询失败 |

#### 注意事项

1. **无需参数**：接口不需要任何参数，自动获取当前用户的数据
2. **今日数据**：today_nutrition 是今天（服务器时间）的营养摄入统计
3. **实时更新**：添加或修改餐饮记录后，today_nutrition 会实时更新
4. **默认目标**：如果用户未设置营养目标，使用系统默认值（热量 2000 千卡，蛋白质 150 克，碳水化合物 250 克，脂肪 70 克）
5. **未来计划**：upcoming_plans 包含明天和后天的待执行计划（status 为 pending）
6. **计划排序**：未来计划按日期和餐次类型排序
7. **空数据处理**：如果今天没有餐饮记录，today_nutrition 全为 0
8. **无计划情况**：如果未来 2 天没有计划，upcoming_plans 为空数组
9. **数据来源**：
   - today_nutrition：来自营养分析模块的每日统计
   - nutrition_goal：来自用户偏好设置
   - upcoming_plans：来自饮食计划模块
10. **性能优化**：建议在客户端缓存数据，避免频繁请求
11. **刷新时机**：建议在以下情况刷新数据：
    - 用户打开应用时
    - 添加或修改餐饮记录后
    - 添加或修改饮食计划后
    - 修改用户偏好设置后
12. **时区处理**：日期基于服务器时区，注意客户端和服务器的时区差异

---

## 数据模型

### DashboardData 模型

Dashboard 数据模型：

```typescript
interface DashboardData {
  today_nutrition: TodayNutrition;     // 今日营养摄入
  nutrition_goal: NutritionGoal;       // 营养目标
  upcoming_plans: UpcomingPlan[];      // 未来计划
}
```

### TodayNutrition 模型

今日营养摄入模型：

```typescript
interface TodayNutrition {
  calories: number;   // 热量（千卡）
  protein: number;    // 蛋白质（克）
  carbs: number;      // 碳水化合物（克）
  fat: number;        // 脂肪（克）
}
```

### NutritionGoal 模型

营养目标模型：

```typescript
interface NutritionGoal {
  calories: number;   // 目标热量（千卡/天）
  protein: number;    // 目标蛋白质（克/天）
  carbs: number;      // 目标碳水化合物（克/天）
  fat: number;        // 目标脂肪（克/天）
}
```

### UpcomingPlan 模型

未来计划模型：

```typescript
interface UpcomingPlan {
  id: number;                                      // 计划 ID
  date: string;                                    // 计划日期（YYYY-MM-DD）
  meal_type: 'breakfast' | 'lunch' | 'dinner' | 'snack';  // 餐次类型
  reason: string;                                  // AI 生成的计划理由
}
```

---

## 使用场景

### 场景 1：应用首页加载

用户打开应用时，加载 Dashboard 数据：

```bash
curl -X GET "http://localhost:9090/api/v1/dashboard" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 场景 2：添加餐饮记录后刷新

用户添加餐饮记录后，刷新 Dashboard 数据以更新今日营养摄入：

```bash
# 1. 添加餐饮记录
curl -X POST "http://localhost:9090/api/v1/meals" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "meal_date": "2024-11-17T12:00:00Z",
    "meal_type": "lunch",
    "foods": [...]
  }'

# 2. 刷新 Dashboard 数据
curl -X GET "http://localhost:9090/api/v1/dashboard" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 场景 3：查看营养达标情况

用户想快速了解今天的营养摄入是否达标：

```bash
curl -X GET "http://localhost:9090/api/v1/dashboard" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

然后在客户端计算达标百分比：
```javascript
const response = await fetch('/api/v1/dashboard');
const data = await response.json();

const caloriesPercentage = (data.data.today_nutrition.calories / data.data.nutrition_goal.calories) * 100;
const proteinPercentage = (data.data.today_nutrition.protein / data.data.nutrition_goal.protein) * 100;
// ...
```

### 场景 4：查看未来计划

用户想查看明天和后天的饮食计划：

```bash
curl -X GET "http://localhost:9090/api/v1/dashboard" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 场景 5：定期刷新数据

在单页应用中，定期刷新 Dashboard 数据：

```javascript
// 每 5 分钟刷新一次
setInterval(async () => {
  const response = await fetch('/api/v1/dashboard', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  const data = await response.json();
  updateDashboard(data.data);
}, 5 * 60 * 1000);
```

---

## 最佳实践

### 1. 数据加载

- **首次加载**：应用启动时立即加载 Dashboard 数据
- **缓存策略**：在客户端缓存数据，避免频繁请求
- **刷新时机**：在关键操作后刷新数据（添加餐饮记录、修改计划等）
- **加载状态**：显示加载状态，提升用户体验
- **错误处理**：优雅地处理加载失败的情况

### 2. 数据展示

- **可视化**：使用图表展示营养摄入和目标的对比
- **进度条**：使用进度条展示营养达标百分比
- **颜色标识**：用颜色区分达标、未达标和超标
- **卡片布局**：使用卡片布局分别展示不同的数据模块
- **响应式设计**：适配不同屏幕尺寸

### 3. 营养对比

- **百分比计算**：计算实际摄入与目标的百分比
- **差值显示**：显示与目标的差值（正值或负值）
- **视觉反馈**：使用颜色和图标提供视觉反馈
- **建议提示**：根据达标情况提供饮食建议
- **趋势指示**：显示与昨天相比的变化趋势

### 4. 计划展示

- **时间排序**：按日期和餐次类型排序展示计划
- **快速操作**：提供快速查看计划详情的入口
- **状态标识**：标识计划的状态（待执行、已完成等）
- **提醒功能**：在计划时间临近时提醒用户
- **空状态**：当没有计划时，引导用户创建计划

### 5. 性能优化

- **数据缓存**：在客户端缓存 Dashboard 数据
- **增量更新**：只更新变化的部分，而不是整个页面
- **懒加载**：对于非关键数据，使用懒加载
- **防抖节流**：对频繁的刷新操作进行防抖或节流
- **离线支持**：支持离线查看缓存的数据

### 6. 用户体验

- **快速响应**：确保页面快速加载和响应
- **平滑过渡**：使用动画和过渡效果
- **即时反馈**：操作后立即更新 UI
- **错误提示**：友好地展示错误信息
- **引导提示**：为新用户提供引导和提示

### 7. 数据刷新

- **手动刷新**：提供下拉刷新功能
- **自动刷新**：在后台自动刷新数据
- **刷新指示**：显示数据刷新状态
- **刷新频率**：控制刷新频率，避免过于频繁
- **网络状态**：根据网络状态调整刷新策略

### 8. 交互设计

- **点击跳转**：点击数据卡片跳转到详情页
- **快捷操作**：提供快捷操作按钮（添加餐饮记录、创建计划等）
- **手势支持**：支持手势操作（滑动、长按等）
- **反馈动画**：操作时提供视觉反馈
- **无障碍**：支持无障碍访问

---

## 常见问题

### Q: Dashboard 数据多久更新一次？

A: 
- Dashboard 数据是实时计算的
- 每次请求都会重新获取最新数据
- 建议在关键操作后手动刷新
- 可以在客户端实现定期自动刷新（如每 5 分钟）
- 避免过于频繁的请求，建议使用缓存

### Q: 为什么今日营养数据全是 0？

A: 
- 可能是今天还没有添加餐饮记录
- 检查是否添加了今天的餐饮记录
- 确认餐饮记录的日期是否正确
- 确认服务器时区设置是否正确

### Q: 如何修改营养目标？

A: 
- 在设置管理模块中修改用户偏好
- 可以设置每日热量、蛋白质、碳水化合物和脂肪的目标值
- 修改后，Dashboard 会自动使用新的目标值
- 详见 [设置管理模块文档](./08-settings.md)

### Q: 未来计划为什么是空的？

A: 
- 可能是未来 2 天没有创建饮食计划
- 只显示状态为 pending（待执行）的计划
- 已完成的计划不会显示在这里
- 可以在饮食计划模块中创建新计划
- 详见 [饮食计划模块文档](./04-plans.md)

### Q: Dashboard 数据包含哪些营养素？

A: 
- today_nutrition 包含：热量、蛋白质、碳水化合物、脂肪
- 不包含纤维数据（纤维数据可以在营养分析模块中查看）
- nutrition_goal 也只包含这 4 个营养素
- 这是为了简化首页展示，突出核心营养指标

### Q: 如何计算营养达标百分比？

A: 
- 百分比 = (实际值 / 目标值) × 100
- 例如：实际热量 1850 千卡，目标 2000 千卡，百分比为 92.5%
- 可以在客户端进行计算
- 90-110% 通常是合理的范围

### Q: 未来计划最多显示多少条？

A: 
- 显示未来 2 天（明天和后天）的计划
- 每天可能有多个计划（早餐、午餐、晚餐、加餐）
- 最多可能显示 8 条（2 天 × 4 餐次）
- 实际数量取决于用户创建的计划数量

### Q: Dashboard 数据可以导出吗？

A: 
- 当前版本不支持直接导出
- 可以在客户端获取数据后导出
- 建议使用截图或 PDF 导出功能
- 未来版本可能会添加导出功能

### Q: 如何处理时区问题？

A: 
- 服务器使用服务器时区计算"今天"
- 客户端应该根据服务器时区显示数据
- 或者在客户端转换为本地时区
- 建议在 API 响应中包含时区信息

### Q: Dashboard 接口性能如何？

A: 
- 接口会聚合多个模块的数据，可能稍慢
- 建议在客户端缓存数据
- 避免频繁请求
- 可以使用加载状态提升用户体验
- 服务器端已经进行了优化

### Q: 可以自定义 Dashboard 显示的内容吗？

A: 
- 当前版本不支持自定义
- 显示的内容是固定的（今日营养、营养目标、未来计划）
- 未来版本可能会添加自定义功能
- 可以在客户端选择性地展示部分数据

### Q: Dashboard 数据会包含其他用户的数据吗？

A: 
- 不会，出于隐私保护
- 每个用户只能看到自己的数据
- 数据通过认证 Token 进行隔离
- 确保数据安全和隐私

---

## 相关文档

- [营养分析模块](./06-nutrition.md) - 了解营养统计的详细信息
- [饮食计划模块](./04-plans.md) - 了解如何创建和管理饮食计划
- [餐饮记录模块](./03-meals.md) - 了解如何添加和管理餐饮记录
- [设置管理模块](./08-settings.md) - 了解如何设置营养目标
- [数据模型](./data-models.md) - 查看数据模型的完整定义
- [通用概念](./common-concepts.md) - 了解认证、响应格式等通用概念
- [错误码说明](./error-codes.md) - 查看所有错误码的详细说明
- [API 文档总览](./README.md) - 返回 API 文档首页
