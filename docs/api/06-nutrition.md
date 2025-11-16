# 营养分析模块

## 概述

营养分析模块提供全面的营养数据统计和分析功能，帮助用户了解自己的饮食营养摄入情况。通过每日统计、月度趋势和目标对比，用户可以清晰地看到自己的营养摄入是否达标，从而调整饮食计划。系统会自动汇总所有餐饮记录的营养数据，提供直观的数据展示和分析。

**核心功能**：
- 获取指定日期的营养统计数据
- 查看整月的营养趋势变化
- 对比实际摄入与目标营养值

**数据特性**：
- 自动汇总餐饮记录的营养数据
- 支持单日和日期范围查询
- 提供差值和百分比对比
- 包含餐次数量统计

**分析维度**：
- 蛋白质（Protein）
- 碳水化合物（Carbs）
- 脂肪（Fat）
- 纤维（Fiber）
- 热量（Calories）

---

## 接口列表

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/nutrition/daily/:date` | 获取每日营养统计 | 是 |
| GET | `/api/v1/nutrition/monthly` | 获取月度营养趋势 | 是 |
| GET | `/api/v1/nutrition/compare` | 对比实际与目标营养 | 是 |

---

## 接口详情


### 获取每日营养统计

**接口**: `GET /api/v1/nutrition/daily/:date`

**说明**: 获取指定日期的营养统计数据。系统会自动汇总该日期所有餐饮记录的营养数据，包括蛋白质、碳水化合物、脂肪、纤维和热量。同时提供餐次数量统计，帮助用户了解当天的饮食情况。

**认证**: 是


#### 请求参数

##### 路径参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| date | string | 是 | 日期（YYYY-MM-DD 格式） | 2024-11-17 |

#### 请求示例

```bash
# 获取 2024 年 11 月 17 日的营养统计
curl -X GET "http://localhost:9090/api/v1/nutrition/daily/2024-11-17" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取今天的营养统计
curl -X GET "http://localhost:9090/api/v1/nutrition/daily/$(date +%Y-%m-%d)" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取昨天的营养统计
curl -X GET "http://localhost:9090/api/v1/nutrition/daily/2024-11-16" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "date": "2024-11-17T00:00:00Z",
    "nutrition": {
      "protein": 145.5,
      "carbs": 220.3,
      "fat": 65.8,
      "fiber": 28.5,
      "calories": 2050.0
    },
    "meal_count": 4
  },
  "timestamp": 1699999999
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| date | string | 日期（ISO 8601 格式） |
| nutrition | object | 营养数据汇总 |
| nutrition.protein | float | 蛋白质总量（克） |
| nutrition.carbs | float | 碳水化合物总量（克） |
| nutrition.fat | float | 脂肪总量（克） |
| nutrition.fiber | float | 纤维总量（克） |
| nutrition.calories | float | 热量总量（千卡） |
| meal_count | int | 餐次数量 |


**无数据响应 (200)**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "date": "2024-11-17T00:00:00Z",
    "nutrition": {
      "protein": 0,
      "carbs": 0,
      "fat": 0,
      "fiber": 0,
      "calories": 0
    },
    "meal_count": 0
  },
  "timestamp": 1699999999
}
```

**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid date format, expected YYYY-MM-DD",
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
  "error": "failed to get daily nutrition statistics",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | 日期格式不正确（非 YYYY-MM-DD 格式） |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | 服务器内部错误、数据库查询失败 |

#### 注意事项

1. **日期格式**：必须使用 YYYY-MM-DD 格式，如 2024-11-17
2. **时区处理**：日期会被转换为当天的 00:00:00 到 23:59:59
3. **无数据情况**：如果当天没有餐饮记录，返回全 0 的营养数据，不会报错
4. **数据来源**：营养数据来自该日期所有餐饮记录的汇总
5. **实时更新**：添加或修改餐饮记录后，统计数据会实时更新
6. **餐次统计**：meal_count 表示当天记录的餐次数量
7. **历史数据**：可以查询任意历史日期的数据
8. **未来日期**：可以查询未来日期，但通常没有数据
9. **数据精度**：营养数据保留小数点后一位
10. **隐私保护**：只能查询属于当前用户的数据

---

### 获取月度营养趋势

**接口**: `GET /api/v1/nutrition/monthly`

**说明**: 获取指定月份每一天的营养统计数据，用于展示月度营养趋势。系统会返回该月每一天的营养数据，即使某天没有餐饮记录也会返回（营养值为 0）。这样可以完整展示整月的营养摄入趋势，方便用户进行月度分析和对比。

**认证**: 是


#### 请求参数

##### 查询参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| year | int | 是 | 年份 | 2024 |
| month | int | 是 | 月份（1-12） | 11 |

#### 请求示例

```bash
# 获取 2024 年 11 月的营养趋势
curl -X GET "http://localhost:9090/api/v1/nutrition/monthly?year=2024&month=11" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取 2024 年 1 月的营养趋势
curl -X GET "http://localhost:9090/api/v1/nutrition/monthly?year=2024&month=1" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取当前月份的营养趋势
curl -X GET "http://localhost:9090/api/v1/nutrition/monthly?year=$(date +%Y)&month=$(date +%-m)" \
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
      "date": "2024-11-01T00:00:00Z",
      "nutrition": {
        "protein": 135.2,
        "carbs": 210.5,
        "fat": 58.3,
        "fiber": 25.0,
        "calories": 1950.0
      },
      "meal_count": 3
    },
    {
      "date": "2024-11-02T00:00:00Z",
      "nutrition": {
        "protein": 142.8,
        "carbs": 225.3,
        "fat": 62.5,
        "fiber": 27.5,
        "calories": 2080.0
      },
      "meal_count": 4
    },
    {
      "date": "2024-11-03T00:00:00Z",
      "nutrition": {
        "protein": 0,
        "carbs": 0,
        "fat": 0,
        "fiber": 0,
        "calories": 0
      },
      "meal_count": 0
    }
  ],
  "timestamp": 1699999999
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| data | array | 每日营养统计数组（按日期升序） |
| data[].date | string | 日期（ISO 8601 格式） |
| data[].nutrition | object | 当天的营养数据汇总 |
| data[].nutrition.protein | float | 蛋白质总量（克） |
| data[].nutrition.carbs | float | 碳水化合物总量（克） |
| data[].nutrition.fat | float | 脂肪总量（克） |
| data[].nutrition.fiber | float | 纤维总量（克） |
| data[].nutrition.calories | float | 热量总量（千卡） |
| data[].meal_count | int | 当天的餐次数量 |


**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "year and month parameters are required",
  "timestamp": 1699999999
}
```

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid month parameter, must be between 1 and 12",
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
  "error": "failed to get monthly nutrition trend",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | year 或 month 参数缺失、year 超出范围（1900-2100）、month 超出范围（1-12） |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | 服务器内部错误、数据库查询失败 |

#### 注意事项

1. **必填参数**：year 和 month 参数都是必填的
2. **年份范围**：year 必须在 1900-2100 之间
3. **月份范围**：month 必须在 1-12 之间
4. **完整月份**：返回该月所有天的数据（28-31 天，根据月份而定）
5. **无数据天数**：没有餐饮记录的天数也会返回，营养值为 0
6. **数据排序**：数据按日期升序排列（从 1 号到月末）
7. **趋势分析**：可以用于绘制营养趋势图表
8. **数据完整性**：即使某天没有数据，也会在数组中占位
9. **性能考虑**：一次查询返回整月数据（最多 31 天）
10. **隐私保护**：只能查询属于当前用户的数据

---

### 对比实际与目标营养

**接口**: `GET /api/v1/nutrition/compare`

**说明**: 对比指定日期范围内的实际营养摄入与用户设置的目标营养值。系统会计算日期范围内的平均营养摄入，然后与用户在偏好设置中配置的目标值进行对比，提供差值和百分比数据。这有助于用户了解自己的营养摄入是否达标，从而调整饮食计划。

**认证**: 是


#### 请求参数

##### 查询参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| start_date | string | 是 | 开始日期（YYYY-MM-DD 格式） | 2024-11-01 |
| end_date | string | 是 | 结束日期（YYYY-MM-DD 格式） | 2024-11-17 |

#### 请求示例

```bash
# 对比 11 月 1 日到 17 日的营养摄入
curl -X GET "http://localhost:9090/api/v1/nutrition/compare?start_date=2024-11-01&end_date=2024-11-17" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 对比最近 7 天的营养摄入
curl -X GET "http://localhost:9090/api/v1/nutrition/compare?start_date=2024-11-10&end_date=2024-11-17" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 对比单日的营养摄入（开始和结束日期相同）
curl -X GET "http://localhost:9090/api/v1/nutrition/compare?start_date=2024-11-17&end_date=2024-11-17" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### 响应示例

**成功响应 (200)**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "target": {
      "protein": 150.0,
      "carbs": 250.0,
      "fat": 70.0,
      "fiber": 25.0,
      "calories": 2000.0
    },
    "actual": {
      "protein": 142.5,
      "carbs": 235.8,
      "fat": 65.3,
      "fiber": 23.5,
      "calories": 1950.0
    },
    "difference": {
      "protein": -7.5,
      "carbs": -14.2,
      "fat": -4.7,
      "fiber": -1.5,
      "calories": -50.0
    },
    "percentage": {
      "protein": 95.0,
      "carbs": 94.32,
      "fat": 93.29,
      "fiber": 94.0,
      "calories": 97.5
    }
  },
  "timestamp": 1699999999
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| target | object | 目标营养值（来自用户偏好设置） |
| target.protein | float | 目标蛋白质（克/天） |
| target.carbs | float | 目标碳水化合物（克/天） |
| target.fat | float | 目标脂肪（克/天） |
| target.fiber | float | 目标纤维（克/天） |
| target.calories | float | 目标热量（千卡/天） |
| actual | object | 实际平均营养摄入 |
| actual.protein | float | 实际平均蛋白质（克/天） |
| actual.carbs | float | 实际平均碳水化合物（克/天） |
| actual.fat | float | 实际平均脂肪（克/天） |
| actual.fiber | float | 实际平均纤维（克/天） |
| actual.calories | float | 实际平均热量（千卡/天） |
| difference | object | 差值（实际 - 目标） |
| difference.protein | float | 蛋白质差值（克） |
| difference.carbs | float | 碳水化合物差值（克） |
| difference.fat | float | 脂肪差值（克） |
| difference.fiber | float | 纤维差值（克） |
| difference.calories | float | 热量差值（千卡） |
| percentage | object | 完成百分比（实际 / 目标 × 100） |
| percentage.protein | float | 蛋白质完成百分比 |
| percentage.carbs | float | 碳水化合物完成百分比 |
| percentage.fat | float | 脂肪完成百分比 |
| percentage.fiber | float | 纤维完成百分比 |
| percentage.calories | float | 热量完成百分比 |


**错误响应 (400)**:

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "start_date and end_date parameters are required",
  "timestamp": 1699999999
}
```

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "invalid start_date format, expected YYYY-MM-DD",
  "timestamp": 1699999999
}
```

```json
{
  "code": 40001,
  "message": "invalid parameters",
  "error": "end_date must be after start_date",
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
  "error": "failed to get user preferences",
  "timestamp": 1699999999
}
```

```json
{
  "code": 50001,
  "message": "internal error",
  "error": "failed to compare nutrition",
  "timestamp": 1699999999
}
```

#### 错误码

| 错误码 | 说明 | 场景 |
|--------|------|------|
| 40001 | 参数错误 | start_date 或 end_date 缺失、日期格式不正确、end_date 早于 start_date |
| 40101 | 未授权 | 用户未认证或 Token 无效 |
| 50001 | 内部错误 | 服务器内部错误、数据库查询失败、获取用户偏好失败 |

#### 注意事项

1. **必填参数**：start_date 和 end_date 都是必填的
2. **日期格式**：必须使用 YYYY-MM-DD 格式
3. **日期范围**：end_date 必须大于或等于 start_date
4. **单日对比**：可以设置相同的开始和结束日期来对比单日数据
5. **平均计算**：如果日期范围包含多天，actual 是这些天的平均值
6. **目标来源**：target 值来自用户偏好设置（设置管理模块）
7. **默认目标**：如果用户未设置目标，使用系统默认值（热量 2000 千卡，蛋白质 150 克等）
8. **差值含义**：正值表示超过目标，负值表示未达到目标
9. **百分比含义**：100% 表示刚好达标，大于 100% 表示超标，小于 100% 表示未达标
10. **无数据处理**：如果日期范围内没有餐饮记录，actual 全为 0
11. **隐私保护**：只能查询属于当前用户的数据
12. **目标调整**：可以在设置管理模块中调整目标营养值

---

## 数据模型

### DailyNutritionStats 模型

每日营养统计数据模型：

```typescript
interface DailyNutritionStats {
  date: string;              // 日期（ISO 8601 格式）
  nutrition: NutritionData;  // 营养数据汇总
  meal_count: number;        // 餐次数量
}
```

### NutritionComparison 模型

营养对比数据模型：

```typescript
interface NutritionComparison {
  target: NutritionData;              // 目标营养值
  actual: NutritionData;              // 实际营养值
  difference: NutritionData;          // 差值（实际 - 目标）
  percentage: {                       // 完成百分比
    protein: number;                  // 蛋白质百分比
    carbs: number;                    // 碳水化合物百分比
    fat: number;                      // 脂肪百分比
    fiber: number;                    // 纤维百分比
    calories: number;                 // 热量百分比
  };
}
```

### NutritionData 模型

营养数据模型：

```typescript
interface NutritionData {
  protein: number;   // 蛋白质（克）
  carbs: number;     // 碳水化合物（克）
  fat: number;       // 脂肪（克）
  fiber: number;     // 纤维（克）
  calories: number;  // 热量（千卡）
}
```

---

## 使用场景

### 场景 1：查看今天的营养摄入

用户想知道今天吃了多少营养：

```bash
curl -X GET "http://localhost:9090/api/v1/nutrition/daily/2024-11-17" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 场景 2：分析本月的营养趋势

用户想查看本月每天的营养摄入变化：

```bash
curl -X GET "http://localhost:9090/api/v1/nutrition/monthly?year=2024&month=11" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 场景 3：对比最近一周的营养达标情况

用户想知道最近一周的营养摄入是否达标：

```bash
curl -X GET "http://localhost:9090/api/v1/nutrition/compare?start_date=2024-11-10&end_date=2024-11-17" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 场景 4：检查单日营养是否达标

用户想检查某一天的营养摄入是否达标：

```bash
curl -X GET "http://localhost:9090/api/v1/nutrition/compare?start_date=2024-11-17&end_date=2024-11-17" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 场景 5：对比不同月份的营养情况

用户想对比 10 月和 11 月的营养摄入：

```bash
# 获取 10 月数据
curl -X GET "http://localhost:9090/api/v1/nutrition/monthly?year=2024&month=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取 11 月数据
curl -X GET "http://localhost:9090/api/v1/nutrition/monthly?year=2024&month=11" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## 最佳实践

### 1. 数据查询

- **合理频率**：不要过于频繁地查询数据，建议使用缓存
- **日期范围**：对比接口建议使用 7-30 天的日期范围
- **时区处理**：注意客户端和服务器的时区差异
- **数据刷新**：添加或修改餐饮记录后，重新查询统计数据

### 2. 趋势分析

- **图表展示**：使用月度数据绘制营养趋势图表
- **异常识别**：关注营养摄入异常高或低的日期
- **规律发现**：分析营养摄入的周期性规律
- **对比分析**：对比不同月份或不同时期的数据

### 3. 目标管理

- **设置目标**：在设置管理模块中配置合理的营养目标
- **定期调整**：根据身体状况和目标变化调整营养目标
- **参考标准**：参考专业营养师的建议设置目标
- **个性化**：根据个人情况（年龄、性别、活动量等）设置目标

### 4. 数据解读

- **差值分析**：关注 difference 字段，了解与目标的差距
- **百分比参考**：percentage 在 90-110% 之间通常是合理的
- **综合评估**：不要只关注单一营养素，要综合评估
- **长期趋势**：关注长期趋势，而不是单日波动

### 5. 前端展示

- **可视化**：使用图表（折线图、柱状图、雷达图等）展示数据
- **颜色标识**：用颜色区分达标、未达标和超标
- **数据对比**：并排展示目标值和实际值
- **趋势指示**：显示营养摄入的上升或下降趋势

### 6. 性能优化

- **缓存策略**：对历史数据进行客户端缓存
- **按需加载**：使用分页或懒加载处理大量数据
- **数据聚合**：在客户端进行数据聚合和计算
- **异步加载**：使用异步请求，不阻塞用户操作

### 7. 用户体验

- **加载状态**：显示数据加载状态
- **空状态提示**：当没有数据时，提示用户添加餐饮记录
- **错误处理**：友好地展示错误信息
- **数据导出**：提供数据导出功能（CSV、PDF 等）

---

## 常见问题

### Q: 营养统计数据多久更新一次？

A: 
- 营养统计数据是实时计算的
- 每次查询都会重新汇总餐饮记录
- 添加、修改或删除餐饮记录后，统计数据会立即更新
- 不需要手动刷新或等待

### Q: 为什么我的营养数据全是 0？

A: 
- 可能是该日期没有餐饮记录
- 检查是否添加了餐饮记录
- 确认日期参数是否正确
- 确认餐饮记录的日期是否匹配

### Q: 如何设置营养目标？

A: 
- 在设置管理模块中配置用户偏好
- 可以设置每日热量、蛋白质、碳水化合物、脂肪和纤维的目标值
- 如果不设置，系统会使用默认值
- 详见 [设置管理模块文档](./08-settings.md)

### Q: 对比接口的平均值是如何计算的？

A: 
- 系统会汇总日期范围内所有天的营养数据
- 然后除以天数得到平均值
- 没有餐饮记录的天数也会计入天数
- 建议只对比有数据的日期范围

### Q: 百分比超过 100% 是什么意思？

A: 
- 百分比 = (实际值 / 目标值) × 100
- 超过 100% 表示实际摄入超过了目标值
- 例如：目标热量 2000 千卡，实际 2200 千卡，百分比为 110%
- 适度超标（100-110%）通常是可以接受的

### Q: 如何判断营养摄入是否合理？

A: 
- 参考 percentage 字段，90-110% 通常是合理的
- 关注 difference 字段，了解具体差距
- 综合评估所有营养素，不要只关注单一指标
- 建议咨询专业营养师

### Q: 月度趋势数据量很大怎么办？

A: 
- 月度数据最多 31 天，数据量不会太大
- 可以在客户端进行缓存
- 可以使用图表库进行可视化展示
- 可以按周或按旬进行数据聚合

### Q: 可以查询多个月的数据吗？

A: 
- 月度接口一次只能查询一个月
- 如需查询多个月，需要多次调用接口
- 可以使用对比接口查询跨月的日期范围
- 建议在客户端进行多月数据的聚合

### Q: 营养数据的精度是多少？

A: 
- 营养数据保留小数点后一位
- 例如：145.5 克、2050.0 千卡
- 百分比保留小数点后两位
- 精度足够满足日常使用需求

### Q: 如何导出营养数据？

A: 
- 当前版本不支持直接导出
- 可以在客户端获取数据后导出为 CSV 或 PDF
- 建议使用第三方库（如 Papa Parse、jsPDF）
- 未来版本可能会添加服务端导出功能

### Q: 营养统计包含哪些餐次？

A: 
- 包含所有类型的餐次：breakfast、lunch、dinner、snack
- meal_count 字段统计了餐次数量
- 不区分餐次类型，全部汇总
- 如需按餐次类型统计，需要在客户端处理

### Q: 可以对比不同用户的营养数据吗？

A: 
- 不可以，出于隐私保护
- 每个用户只能查询自己的数据
- 如需对比，可以在客户端手动对比
- 建议使用匿名化的数据进行对比

---

## 相关文档

- [数据模型](./data-models.md) - 查看 NutritionData 模型的完整定义
- [餐饮记录模块](./03-meals.md) - 了解如何添加和管理餐饮记录
- [设置管理模块](./08-settings.md) - 了解如何设置营养目标
- [Dashboard 模块](./07-dashboard.md) - 查看营养数据的综合展示
- [通用概念](./common-concepts.md) - 了解认证、响应格式等通用概念
- [错误码说明](./error-codes.md) - 查看所有错误码的详细说明
- [API 文档总览](./README.md) - 返回 API 文档首页
