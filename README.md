# AI Diet Assistant

<div align="center">

🍎 AI 驱动的个性化饮食计划助手系统

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>

---

## 📖 简介

AI Diet Assistant 是一个智能饮食管理系统，通过集成 AI 服务帮助用户管理饮食和营养。

**核心功能：**

- 🥗 **智能饮食计划** - AI 自动生成个性化饮食计划
- 📊 **营养分析** - 实时追踪营养摄入，对比目标值
- 🛒 **食材管理** - 管理个人食材库，记录价格和营养信息
- � **AI 对*话流管理** - 管理与 AI 的对话历史，支持收藏和搜索
- � **食消息代理** - 转发消息到外部 AI 服务，保存对话记录
- � **数据可视化*管* - 直观展示饮食历史和营养趋势

## � 快速开始

### 前置要求

在运行安装脚本之前，请确保系统满足以下要求：

- **Go**: 1.21 或更高版本
- **MySQL**: 8.0 或更高版本
- **Redis**: 6 或更高版本（可选，用于 Token 黑名单）
- **Nginx**: 最新稳定版（生产环境推荐）
- **OpenSSL**: 用于生成安全密钥

### 一键安装（5 分钟）

```bash
# 1. 克隆项目
git clone https://github.com/Deepblue-Sky2333/Ai-Diet-Assistant/
cd Ai-Diet-Assistant

# 2. 运行一键安装脚本
./install.sh
```

### 安装脚本会自动完成

1. **检测系统依赖**
   - 检测 Go、MySQL、Redis、openssl
   - 如果缺少，提示安装方法

2. **配置系统**
   - 自动生成安全密钥（JWT 和 AES）
   - 配置数据库连接
   - 配置 Redis（可选）
   - 配置 Go 模块路径

3. **创建数据库**
   - 自动创建数据库
   - 运行数据库迁移

4. **构建应用**
   - 下载 Go 依赖
   - 构建后端 API 服务

5. **配置服务**（Linux 系统可选）
   - 配置 systemd 服务
   - 设置开机自启
   - 提供服务管理命令

### 安装完成

安装完成后，后端 API 服务会自动启动（如果配置了服务）。

**后端 API 地址**：http://localhost:9090/api/v1

**健康检查**：http://localhost:9090/health

**注意**：本项目是纯后端 API 服务，不包含前端界面。如需访问 API，请：
- 使用 API 客户端（如 Postman、curl）
- 开发自己的前端应用
- 配置 Nginx 作为反向代理（生产环境推荐）

### 服务管理

#### Linux 系统（systemd）

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

#### 手动启动

```bash
# 直接运行
./bin/diet-assistant

# 或使用脚本
./scripts/start.sh

# 停止服务
./scripts/stop.sh

# 查看状态
./scripts/status.sh
```

### 测试 API

使用 curl 测试 API：

```bash
# 健康检查
curl http://localhost:9090/health

# 登录获取 Token
curl -X POST http://localhost:9090/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your_username",
    "password": "your_password"
  }'

# 使用 Token 访问 API
curl -X GET http://localhost:9090/api/v1/foods \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 📱 API 功能概览

本系统提供 33 个 RESTful API 接口，涵盖以下功能模块：

1. **认证模块** (4 个接口)
   - 用户登录、Token 刷新、登出、密码修改

2. **食材管理** (6 个接口)
   - 创建、查询、更新、删除食材，支持批量导入

3. **餐饮记录** (5 个接口)
   - 记录每日餐饮，自动计算营养摄入

4. **饮食计划** (6 个接口)
   - AI 生成个性化饮食计划，管理计划

5. **AI 服务** (3 个接口)
   - AI 对话、餐饮建议、对话历史

6. **营养分析** (3 个接口)
   - 每日统计、月度趋势、营养对比

7. **Dashboard** (1 个接口)
   - 获取仪表盘数据

8. **设置管理** (5 个接口)
   - AI 设置、用户偏好、用户资料

详细 API 文档请参考：[docs/api/README.md](docs/api/README.md)

## 💬 AI 对话流功能

系统提供完整的 AI 对话流管理功能：

### 对话流管理
- **创建对话流** - 开始新的 AI 对话会话
- **自动历史管理** - 自动保留最近 10 条对话流
- **收藏功能** - 收藏重要对话（最多 100 条）
- **搜索功能** - 按标题搜索对话流
- **导出功能** - 导出对话为 JSON 格式

### 消息代理
- **消息转发** - 将用户消息转发到外部 AI 服务
- **原始存储** - 保存原始请求和响应数据
- **历史记录** - 完整保存对话历史

### API 使用示例

```bash
# 1. 创建新对话流
curl -X POST http://localhost:9090/api/v1/conversations \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "饮食咨询"}'

# 2. 发送消息
curl -X POST http://localhost:9090/api/v1/conversations/1/messages \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": "今天晚餐吃什么好？"}'

# 3. 获取对话历史
curl -X GET http://localhost:9090/api/v1/conversations/1/messages \
  -H "Authorization: Bearer YOUR_TOKEN"

# 4. 收藏对话流
curl -X POST http://localhost:9090/api/v1/conversations/1/favorite \
  -H "Authorization: Bearer YOUR_TOKEN"

# 5. 搜索对话流
curl -X GET "http://localhost:9090/api/v1/conversations/search?keyword=饮食" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 6. 导出对话流
curl -X GET http://localhost:9090/api/v1/conversations/1/export \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 🔧 配置 AI Provider

通过 API 配置 AI Provider：

```bash
# 更新 AI 设置
curl -X PUT http://localhost:9090/api/v1/settings/ai \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "openai",
    "api_key": "your_api_key",
    "model": "gpt-4",
    "base_url": "https://api.openai.com/v1"
  }'

# 测试 AI 连接
curl -X GET http://localhost:9090/api/v1/settings/ai/test \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 🏗️ 技术栈

**后端：**
- Go 1.21+
- Gin Web Framework
- MySQL 8.0+
- Redis 6+ (可选)
- JWT 认证

**部署：**
- Nginx (反向代理和 CORS 处理)
- Systemd (服务管理)

## 📝 项目结构

```
.
├── cmd/                    # 应用程序入口
│   ├── server/            # 主服务器
│   └── create-user/       # 用户创建工具
├── internal/              # 私有应用代码
│   ├── handler/          # HTTP 处理器
│   │   ├── conversation_handler.go  # 对话流管理
│   │   └── message_handler.go       # 消息代理
│   ├── service/          # 业务逻辑层
│   │   ├── conversation_service.go  # 对话流服务
│   │   └── message_proxy_service.go # 消息代理服务
│   ├── repository/       # 数据访问层
│   │   ├── conversation_repository.go  # 对话流数据访问
│   │   └── message_repository.go       # 消息数据访问
│   ├── middleware/       # 中间件
│   ├── model/            # 数据模型
│   │   ├── conversation_flow.go  # 对话流模型
│   │   └── message.go            # 消息模型
│   ├── config/           # 配置管理
│   ├── database/         # 数据库连接
│   ├── ai/               # AI 代理客户端
│   │   └── proxy_client.go  # 外部 AI 服务客户端
│   └── utils/            # 工具函数
├── configs/              # 配置文件
├── migrations/           # 数据库迁移
├── scripts/              # 部署和管理脚本
├── docs/                 # 文档
│   ├── api/             # API 文档
│   └── NGINX_CONFIGURATION.md  # Nginx 配置指南
└── bin/                  # 编译后的二进制文件
```

## 📚 详细安装指南

### 手动安装

如果您希望手动安装，请按照以下步骤操作：

#### 1. 安装依赖

**macOS**
```bash
brew install go mysql redis openssl nginx
brew services start mysql
brew services start redis
brew services start nginx
```

**Ubuntu/Debian**
```bash
sudo apt-get update
sudo apt-get install golang-go mysql-server redis-server openssl nginx
sudo systemctl start mysql
sudo systemctl start redis
sudo systemctl start nginx
```

**CentOS/RHEL**
```bash
sudo yum install golang mysql-server redis openssl nginx
sudo systemctl start mysqld
sudo systemctl start redis
sudo systemctl start nginx
```

#### 2. 配置数据库

```bash
# 登录 MySQL
mysql -u root -p

# 创建数据库和用户
CREATE DATABASE ai_diet_assistant CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'diet_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON ai_diet_assistant.* TO 'diet_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

#### 3. 配置应用

```bash
# 复制配置文件示例
cp configs/config.yaml.example configs/config.yaml

# 编辑配置文件
vim configs/config.yaml
```

#### 4. 运行数据库迁移

```bash
# 使用迁移脚本
./scripts/run-migrations.sh

# 或手动运行
for file in migrations/*_up.sql; do
    mysql -u diet_user -p ai_diet_assistant < "$file"
done
```

#### 5. 构建应用

```bash
# 下载依赖
go mod download
go mod tidy

# 创建目录
mkdir -p bin logs uploads

# 编译应用
go build -o bin/diet-assistant cmd/server/main.go
```

#### 6. 启动应用

```bash
# 直接运行
./bin/diet-assistant

# 或使用脚本
./scripts/start.sh
```

### 验证安装

```bash
# 使用状态脚本
./scripts/status.sh

# 或手动检查
curl http://localhost:9090/health
```

预期响应：
```json
{
  "status": "ok",
  "timestamp": 1234567890
}
```

## 🚀 生产环境部署

### 重要说明

**本项目是纯后端 API 服务**，不包含前端代码。在生产环境中：

1. **必须使用 Nginx** 作为反向代理
2. **CORS 由 Nginx 处理**，后端不处理 CORS
3. **建议启用 HTTPS**，使用 Let's Encrypt 免费证书

### 1. 配置 Nginx 反向代理

Nginx 负责处理：
- 反向代理到后端 API
- CORS 跨域请求
- SSL/TLS 终止
- 负载均衡（可选）
- 请求限流（可选）

**基础配置示例**：

创建 Nginx 配置文件 `/etc/nginx/sites-available/diet-assistant`:

```nginx
server {
    listen 80;
    server_name api.yourdomain.com;

    # CORS 配置
    add_header 'Access-Control-Allow-Origin' '$http_origin' always;
    add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS' always;
    add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type' always;
    add_header 'Access-Control-Allow-Credentials' 'true' always;

    # 处理 OPTIONS 请求
    if ($request_method = 'OPTIONS') {
        return 204;
    }

    # 代理到后端
    location / {
        proxy_pass http://localhost:9090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

启用配置：
```bash
sudo ln -s /etc/nginx/sites-available/diet-assistant /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

**完整配置和高级功能**（负载均衡、限流、缓存等）请参考：[Nginx 配置指南](docs/NGINX_CONFIGURATION.md)

### 2. 配置 SSL (Let's Encrypt)

```bash
# 安装 Certbot
sudo apt-get install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d api.yourdomain.com

# 测试自动续期
sudo certbot renew --dry-run
```

### 3. 配置系统服务

创建 systemd 服务文件 `/etc/systemd/system/diet-assistant.service`:

```ini
[Unit]
Description=AI Diet Assistant API Server
After=network.target mysql.service redis.service

[Service]
Type=simple
User=your_user
WorkingDirectory=/path/to/Ai-Diet-Assistant
ExecStart=/path/to/Ai-Diet-Assistant/bin/diet-assistant
Restart=on-failure
RestartSec=10

# 环境变量（可选）
Environment="GIN_MODE=release"

[Install]
WantedBy=multi-user.target
```

启用服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable diet-assistant
sudo systemctl start diet-assistant
sudo systemctl status diet-assistant
```

## ⚙️ 配置说明

### 环境变量 (.env)

```bash
# 服务器配置
SERVER_PORT=9090              # 服务器端口
SERVER_MODE=release           # 运行模式: debug/release

# 数据库配置
DB_HOST=localhost             # MySQL 主机
DB_PORT=3306                  # MySQL 端口
DB_USER=diet_user             # 数据库用户
DB_PASSWORD=your_password     # 数据库密码
DB_NAME=ai_diet_assistant     # 数据库名称

# JWT 配置
JWT_SECRET=<auto_generated>   # JWT 密钥（自动生成）

# 加密配置
ENCRYPTION_KEY=<auto_generated> # AES 密钥（自动生成）

# 限流配置
RATE_LIMIT_ENABLED=true       # 启用限流
RATE_LIMIT_REQUESTS_PER_MINUTE=100 # 每分钟请求数

# Redis 配置
REDIS_ENABLED=true            # 启用 Redis
REDIS_HOST=localhost          # Redis 主机
REDIS_PORT=6379               # Redis 端口
REDIS_PASSWORD=               # Redis 密码
REDIS_DB=0                    # Redis 数据库

# 日志配置
LOG_LEVEL=info                # 日志级别
LOG_FORMAT=json               # 日志格式

# 安全配置
MAX_LOGIN_ATTEMPTS=5          # 最大登录尝试次数
LOCKOUT_DURATION=15m          # 锁定时长
PASSWORD_MIN_LENGTH=8         # 最小密码长度
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

### 查看日志

**后端日志**：
```bash
# 应用日志
tail -f logs/app.log

# 系统服务日志
sudo journalctl -u diet-assistant -f
```

## 📝 常见问题

### API 无法访问

**问题**：无法访问 API 接口

**解决方案**：
1. 确保后端正在运行：`curl http://localhost:9090/health`
2. 检查服务状态：`sudo systemctl status diet-assistant`
3. 查看日志：`tail -f logs/app.log`
4. 检查端口是否被占用：`lsof -i :9090`

### 数据库连接失败

**问题**：后端启动时报数据库连接错误

**解决方案**：
1. 确认 MySQL 正在运行：`mysql -u root -p`
2. 检查数据库配置：`cat .env | grep DB_`
3. 确认数据库已创建：`SHOW DATABASES;`
4. 检查用户权限

### 端口已被占用

**问题**：启动时提示端口 9090 已被占用

**解决方案**：
```bash
# 查找占用端口的进程
lsof -i :9090

# 停止进程
kill -9 <PID>

# 或修改端口配置
vim configs/config.yaml
# 修改 server.port 配置
```

### CORS 错误

**问题**：前端调用 API 时出现 CORS 错误

**解决方案**：
1. **开发环境**：后端已移除 CORS 中间件，需要配置 Nginx
2. **生产环境**：必须使用 Nginx 处理 CORS
3. 参考 [Nginx 配置指南](docs/NGINX_CONFIGURATION.md) 配置 CORS

### MySQL 连接失败

**问题**: `无法连接到数据库`

**解决方案**:
```bash
# 检查 MySQL 服务状态
# macOS
brew services list | grep mysql

# Linux
sudo systemctl status mysql

# 启动 MySQL
# macOS
brew services start mysql

# Linux
sudo systemctl start mysql
```

### Go 版本过低

**问题**: `Go 版本过低 (需要 >= 1.21)`

**解决方案**:
```bash
# macOS
brew upgrade go

# Ubuntu/Debian
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt-get update
sudo apt-get install golang-go

# 或从官网下载: https://golang.org/dl/
```

### 权限问题

**问题**: `permission denied`

**解决方案**:
```bash
# 设置可执行权限
chmod +x install.sh
chmod +x bin/diet-assistant
chmod +x scripts/*.sh

# 设置配置文件权限
chmod 600 .env
chmod 600 configs/config.yaml
```

### Redis 未运行

**问题**: Redis 连接失败

**解决方案**:
```bash
# 启动 Redis
# macOS
brew services start redis

# Linux
sudo systemctl start redis

# 或在配置中禁用 Redis
vim .env
# 修改 REDIS_ENABLED=false
```

## 🔒 安全建议

1. **保护配置文件**
   - 不要将 `.env` 和 `config.yaml` 提交到版本控制
   - 设置适当的文件权限（600）

2. **使用强密码**
   - 数据库密码至少 16 字符
   - 包含大小写字母、数字和特殊字符

3. **启用 HTTPS**
   - 在生产环境中始终使用 HTTPS
   - 使用 Let's Encrypt 免费证书

4. **定期备份**
   - 定期备份数据库
   - 备份配置文件

5. **监控日志**
   - 定期检查应用日志
   - 设置日志轮转

6. **更新依赖**
   - 定期更新 Go 依赖
   - 关注安全公告

详细安全指南请参考：[docs/SECURITY.md](docs/SECURITY.md)

## 📚 文档

- [API 文档](docs/api/README.md) - 完整 API 接口说明
- [Nginx 配置](docs/NGINX_CONFIGURATION.md) - Nginx 反向代理配置指南
- [安全最佳实践](docs/SECURITY.md) - 安全配置指南
- [错误码说明](docs/ERROR_CODES.md) - 错误码参考

## 🆘 获取帮助

如果遇到问题：

1. 查看日志文件
2. 检查配置文件
3. 阅读完整文档
4. 提交 Issue: https://github.com/Deepblue-Sky2333/Ai-Diet-Assistant/issues

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

---

**祝你使用愉快！** 🎉
