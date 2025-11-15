# AI Diet Assistant - 快速开始指南

本指南将帮助你在 5 分钟内启动并运行 AI Diet Assistant。

## 📋 前置要求

- **Go**: 1.21 或更高版本
- **MySQL**: 8.0 或更高版本
- **Node.js**: 18 或更高版本
- **npm**: 9 或更高版本

## 🚀 一键安装（5 分钟）

### 安装步骤

```bash
# 1. 克隆项目
git clone https://github.com/Deepblue-Sky2333/Ai-Diet-Assistant/
cd ai-diet-assistant

# 2. 运行一键安装脚本
./install.sh
```

### 安装脚本会自动完成

1. **检测系统依赖**
   - 检测 Go、Node.js、MySQL、openssl
   - 如果缺少，提示自动安装

2. **配置系统**
   - 自动生成安全密钥（JWT 和 AES）
   - 配置数据库连接
   - 配置 CORS 允许的域名
   - 配置 Go 模块路径
   - 配置 Redis（可选）
   - 配置前端（Demo 模式或生产模式）

3. **创建数据库**
   - 自动创建数据库
   - 运行数据库迁移

4. **构建应用**
   - 下载 Go 依赖
   - 构建后端
   - 安装前端依赖
   - 构建前端

5. **配置服务**（Linux 系统可选）
   - 配置 systemd 服务
   - 设置开机自启
   - 提供服务管理命令

### 安装完成

安装完成后，应用会自动启动（如果配置了服务）。

访问：**http://localhost:9090**

### 服务管理（Linux）

如果配置了系统服务：

```bash
# 启动服务
sudo systemctl start diet-assistant

# 停止服务
sudo systemctl stop diet-assistant

# 重启服务
sudo systemctl restart diet-assistant

# 查看状态
sudo systemctl status diet-assistant

# 查看日志
sudo journalctl -u diet-assistant -f
```

### 手动启动（未配置服务）

```bash
# 直接运行
./bin/diet-assistant

# 或使用脚本
./scripts/start.sh
```

### 开发模式

如果需要前后端分离开发（热重载）：

```bash
# 启动后端
./scripts/start.sh

# 启动前端（新终端）
cd web/frontend && npm run dev
```

访问：
- 前端：http://localhost:3000
- 后端：http://localhost:9090

### 登录

#### Demo 模式

点击 "Login with Test Account" 按钮即可登录。

#### 生产模式

使用安装时配置的用户名和密码登录。

## 📱 功能概览

登录后，你可以：

1. **Dashboard（仪表盘）**
   - 查看今日营养摄入
   - 查看未来饮食计划
   - 快速访问各功能模块

2. **Supermarket（超市）**
   - 添加食材到你的食材库
   - 记录价格和营养信息
   - 按分类管理食材

3. **Meals（餐饮记录）**
   - 记录每日餐饮
   - 自动计算营养摄入
   - 查看历史记录

4. **Plans（饮食计划）**
   - AI 生成个性化饮食计划
   - 基于你的食材库
   - 考虑营养平衡

5. **Nutrition（营养分析）**
   - 每日营养统计
   - 月度趋势分析
   - 实际与目标对比

6. **Chat（AI 对话）**
   - 与 AI 助手对话
   - 获取饮食建议
   - 查看对话历史

7. **Settings（设置）**
   - 配置 AI Provider
   - 设置营养目标
   - 管理用户偏好

## 🔧 配置 AI Provider

1. 访问 **Settings** 页面
2. 选择 AI Provider（OpenAI、DeepSeek 或自定义）
3. 输入 API Key
4. 点击 "Test Connection" 验证
5. 保存配置

## 📝 常见问题

### 前端无法连接后端

**问题**：前端显示 "Cannot connect to backend API"

**解决方案**：
1. 确保后端正在运行：`curl http://localhost:9090/health`
2. 检查前端配置：`cat web/frontend/.env.local`
3. 确认 `NEXT_PUBLIC_API_URL` 设置正确
4. 或启用 Demo 模式：设置 `NEXT_PUBLIC_DEMO_MODE=true`

### 数据库连接失败

**问题**：后端启动时报数据库连接错误

**解决方案**：
1. 确认 MySQL 正在运行：`mysql -u root -p`
2. 检查数据库配置：`cat .env | grep DB_`
3. 确认数据库已创建：`SHOW DATABASES;`
4. 检查用户权限

### 端口已被占用

**问题**：启动时提示端口 9090 或 3000 已被占用

**解决方案**：
```bash
# 查找占用端口的进程
lsof -i :9090  # 后端
lsof -i :3000  # 前端

# 停止进程
kill -9 <PID>

# 或修改端口配置
# 后端：编辑 .env 中的 SERVER_PORT
# 前端：Next.js 会自动使用下一个可用端口
```

### 前端依赖安装失败

**问题**：`npm install` 失败

**解决方案**：
```bash
# 清理缓存
cd web/frontend
rm -rf node_modules package-lock.json
npm cache clean --force

# 重新安装
npm install

# 或使用 yarn/pnpm
yarn install
# 或
pnpm install
```

## 🛠️ 开发模式

### 热重载开发

**后端**（使用 air）：
```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 启动热重载
air
```

**前端**（Next.js 自带）：
```bash
cd web/frontend
npm run dev
```

### 查看日志

**后端日志**：
```bash
tail -f logs/app.log
# 或
tail -f logs/backend.log  # 使用 start-all.sh 时
```

**前端日志**：
在浏览器控制台查看（F12 → Console）

## 📚 下一步

- 阅读完整文档：[README.md](README.md)
- 查看 API 文档：[docs/api.md](docs/api.md)
- 了解部署指南：[README.md#生产部署](README.md#生产部署)
- 配置 AI Provider：在 Settings 页面配置

## 🆘 获取帮助

如果遇到问题：

1. 查看日志文件
2. 检查配置文件
3. 阅读完整文档
4. 提交 Issue

---

**祝你使用愉快！** 🎉
