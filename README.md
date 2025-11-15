# AI Diet Assistant

<div align="center">

🍎 AI 驱动的个性化饮食计划助手系统

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![MySQL](https://img.shields.io/badge/MySQL-8.0+-4479A1?style=flat&logo=mysql&logoColor=white)](https://www.mysql.com)
[![Next.js](https://img.shields.io/badge/Next.js-15-black?style=flat&logo=next.js)](https://nextjs.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>

---

## 中文

### 📖 简介

AI Diet Assistant 是一个智能饮食管理系统，通过集成 AI 服务帮助用户管理饮食和营养。

**核心功能：**

- 🥗 **智能饮食计划** - AI 自动生成个性化饮食计划
- 📊 **营养分析** - 实时追踪营养摄入，对比目标值
- 🛒 **食材管理** - 管理个人食材库，记录价格和营养信息
- 💬 **AI 对话** - 获取个性化饮食建议
- 📈 **数据可视化** - 直观展示饮食历史和营养趋势

### 🚀 快速开始

#### 一键安装

```bash
# 1. 克隆项目
git clone <repository-url>
cd ai-diet-assistant

# 2. 运行一键安装脚本
./install.sh
```

安装脚本会自动：
- ✅ 检测并安装依赖（Go、Node.js、MySQL）
- ✅ 配置系统（生成密钥、配置数据库、CORS 等）
- ✅ 创建数据库并运行迁移
- ✅ 构建前后端应用
- ✅ 配置系统服务（可选）

安装完成后，访问：**http://localhost:9090**

#### 开发模式

如果需要前后端分离开发：

```bash
# 启动后端
./scripts/start.sh

# 启动前端（新终端）
cd web/frontend && npm run dev
```

访问：
- 前端：http://localhost:3000
- 后端：http://localhost:9090

**详细说明：** 查看 [快速开始指南](QUICKSTART.md)

### 🏗️ 技术栈

**前端：**
- Next.js 15 + React 19
- TypeScript
- Tailwind CSS v4
- shadcn/ui

**后端：**
- Go 1.25.4
- Gin Web Framework
- MySQL 8.0+
- JWT 认证

### 📚 文档

- [快速开始](QUICKSTART.md) - 5分钟一键安装
- [API 文档](docs/API.md) - API 接口说明
- [安全最佳实践](docs/SECURITY.md) - 安全配置指南
- [前端文档](web/frontend/README.md) - 前端开发文档

### 🔧 开发

```bash
# 后端开发
make run

# 前端开发
cd web/frontend
npm run dev

# 运行测试
make test

# 代码检查
make lint
```

### 📝 项目结构

```
.
├── cmd/                    # 应用程序入口
├── internal/              # 私有应用代码
│   ├── handler/          # HTTP 处理器
│   ├── service/          # 业务逻辑层
│   ├── repository/       # 数据访问层
│   └── middleware/       # 中间件
├── web/frontend/         # Next.js 前端应用
├── configs/              # 配置文件
├── migrations/           # 数据库迁移
├── scripts/              # 部署和管理脚本
└── docs/                 # 文档
```

### 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

