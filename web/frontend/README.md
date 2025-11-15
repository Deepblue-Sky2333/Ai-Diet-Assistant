# AI Diet Assistant - Frontend

这是 AI Diet Assistant 的前端应用，使用 Next.js 15 和 React 19 构建。

## 🚀 快速开始

### 安装依赖

```bash
npm install
```

### 配置环境变量

创建 `.env.local` 文件：

```bash
# Demo 模式（无需后端，使用模拟数据）
NEXT_PUBLIC_DEMO_MODE=true

# 后端 API 地址（Demo 模式下仍可连接 localhost 后端）
NEXT_PUBLIC_API_URL=http://localhost:9090/api/v1
```

**注意**：如果你运行了 `scripts/install.sh`，这个文件已经自动创建。

### 启动开发服务器

```bash
npm run dev
```

应用将在 http://localhost:3000 启动。

### 构建生产版本

```bash
npm run build
npm start
```

## 📁 项目结构

```
web/frontend/
├── app/                    # Next.js App Router 页面
│   ├── dashboard/         # Dashboard 及子页面
│   │   ├── chat/         # AI 对话
│   │   ├── meals/        # 餐饮记录
│   │   ├── nutrition/    # 营养分析
│   │   ├── plans/        # 饮食计划
│   │   ├── settings/     # 设置
│   │   └── supermarket/  # 食材管理
│   ├── login/            # 登录页面
│   ├── layout.tsx        # 根布局
│   └── page.tsx          # 首页
├── components/            # React 组件
│   ├── ui/               # shadcn/ui 组件
│   ├── icons.tsx         # 图标组件
│   ├── sidebar.tsx       # 侧边栏
│   └── theme-provider.tsx # 主题提供者
├── lib/                   # 工具函数
│   ├── api.ts            # API 客户端
│   ├── security.ts       # 安全工具
│   └── utils.ts          # 通用工具
├── hooks/                 # React Hooks
│   ├── use-mobile.ts     # 移动端检测
│   └── use-toast.ts      # Toast 通知
├── public/                # 静态资源
├── styles/                # 全局样式
└── .env.local            # 环境变量（不提交到 Git）
```

## 🎨 技术栈

- **框架**: Next.js 15 (App Router)
- **UI 库**: React 19
- **语言**: TypeScript
- **样式**: Tailwind CSS v4
- **组件**: shadcn/ui
- **图表**: Recharts
- **日期**: date-fns
- **图标**: Lucide React

## 🔧 配置说明

### Demo 模式 vs 生产模式

#### Demo 模式（`NEXT_PUBLIC_DEMO_MODE=true`）

- ✅ 无需后端即可运行
- ✅ 使用模拟数据
- ✅ 所有功能可用
- ✅ 适合快速演示和测试
- ✅ 仍可连接 localhost 后端（如果可用）

**使用场景**：
- 前端开发和测试
- UI/UX 演示
- 无后端环境的快速预览

#### 生产模式（`NEXT_PUBLIC_DEMO_MODE=false`）

- ✅ 连接真实后端 API
- ✅ 真实数据存储
- ✅ 完整功能
- ⚠️ 需要后端服务运行

**使用场景**：
- 生产环境部署
- 完整功能测试
- 真实数据操作

### 环境变量

| 变量名 | 说明 | 默认值 | 必需 |
|--------|------|--------|------|
| `NEXT_PUBLIC_DEMO_MODE` | 是否启用 Demo 模式 | `true` | 否 |
| `NEXT_PUBLIC_API_URL` | 后端 API 地址 | `http://localhost:9090/api/v1` | 是 |

## 🎯 主要功能

### 1. Dashboard（仪表盘）
- 今日营养摄入统计
- 营养目标进度
- 未来饮食计划预览
- 快速操作入口

### 2. Supermarket（超市）
- 食材库管理
- 添加/编辑/删除食材
- 按分类筛选
- 搜索功能
- 营养信息录入

### 3. Meals（餐饮记录）
- 记录每日餐饮
- 自动计算营养
- 历史记录查询
- 按日期筛选

### 4. Plans（饮食计划）
- AI 生成饮食计划
- 基于食材库
- 营养平衡优化
- 计划管理

### 5. Nutrition（营养分析）
- 每日营养统计
- 月度趋势图表
- 实际与目标对比
- 多维度分析

### 6. Chat（AI 对话）
- 与 AI 助手对话
- 获取饮食建议
- 对话历史记录

### 7. Settings（设置）
- AI Provider 配置
- 营养目标设置
- 用户偏好管理
- 密码修改

## 🔐 安全特性

- ✅ JWT 认证
- ✅ Token 自动刷新
- ✅ XSS 防护（输入清理）
- ✅ 密码强度验证
- ✅ API 密钥掩码显示
- ✅ 安全的本地存储

## 🎨 UI/UX 特性

- ✅ 响应式设计（支持移动端）
- ✅ 深色模式支持
- ✅ 流畅的动画效果
- ✅ 友好的错误提示
- ✅ 加载状态指示
- ✅ Toast 通知

## 📝 开发指南

### 添加新页面

1. 在 `app/dashboard/` 下创建新目录
2. 创建 `page.tsx` 文件
3. 实现页面组件
4. 在侧边栏添加导航链接

### 添加新 API 调用

在 `lib/api.ts` 中添加新方法：

```typescript
async myNewApi(param: string) {
  return this.request('/my-endpoint', {
    method: 'POST',
    body: JSON.stringify({ param }),
  });
}
```

### 添加新组件

1. 在 `components/` 下创建组件文件
2. 使用 TypeScript 定义 Props
3. 导出组件

### 使用 shadcn/ui 组件

```bash
# 添加新组件
npx shadcn-ui@latest add button
npx shadcn-ui@latest add card
```

## 🧪 测试

```bash
# 运行测试（如果配置了）
npm test

# 类型检查
npm run type-check

# Lint 检查
npm run lint
```

## 📦 构建和部署

### 开发环境

```bash
npm run dev
```

### 生产构建

```bash
npm run build
npm start
```

### 静态导出（可选）

```bash
npm run build
# 输出到 out/ 目录
```

### 部署到 Vercel

1. 推送代码到 Git 仓库
2. 在 Vercel 导入项目
3. 配置环境变量
4. 部署

### 部署到其他平台

构建后，可以部署到：
- Netlify
- AWS Amplify
- Azure Static Web Apps
- 自托管（使用 nginx）

## 🐛 故障排查

### 无法连接后端

**问题**：显示 "Cannot connect to backend API"

**解决方案**：
1. 确认后端正在运行
2. 检查 `NEXT_PUBLIC_API_URL` 配置
3. 检查 CORS 配置
4. 或启用 Demo 模式

### 依赖安装失败

**问题**：`npm install` 失败

**解决方案**：
```bash
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
```

### 构建失败

**问题**：`npm run build` 失败

**解决方案**：
1. 检查 TypeScript 错误
2. 运行 `npm run lint` 检查代码
3. 确保所有依赖已安装

## 📚 相关文档

- [Next.js 文档](https://nextjs.org/docs)
- [React 文档](https://react.dev)
- [Tailwind CSS 文档](https://tailwindcss.com/docs)
- [shadcn/ui 文档](https://ui.shadcn.com)

## 🤝 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 开启 Pull Request

## 📄 许可证

MIT License - 详见根目录 LICENSE 文件
