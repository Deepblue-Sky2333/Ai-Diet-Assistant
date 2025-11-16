# AI Diet Assistant - 快速开始指南

本指南将帮助你在 5 分钟内启动并运行 AI Diet Assistant。

## 📋 前置要求

- **Go**: 1.21 或更高版本
- **MySQL**: 8.0 或更高版本
- **Redis**: 6 或更高版本（可选，用于 Token 黑名单）
- **Nginx**: 最新稳定版（生产环境推荐）

## 🚀 一键安装（5 分钟）

### 安装步骤

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

### 使用 Nginx 反向代理（推荐）

开发环境也可以配置 Nginx：

```nginx
server {
    listen 80;
    server_name localhost;

    # CORS 配置
    add_header 'Access-Control-Allow-Origin' '*' always;
    add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS' always;
    add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type' always;

    if ($request_method = 'OPTIONS') {
        return 204;
    }

    # 代理到后端
    location / {
        proxy_pass http://localhost:9090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

详细配置请参考：[Nginx 配置指南](docs/NGINX_CONFIGURATION.md)

## 📚 下一步

- 阅读完整文档：[README.md](README.md)
- 查看 API 文档：[docs/api/README.md](docs/api/README.md)
- 配置 Nginx：[docs/NGINX_CONFIGURATION.md](docs/NGINX_CONFIGURATION.md)
- 了解安装详情：[INSTALLATION_GUIDE.md](INSTALLATION_GUIDE.md)
- 查看安全指南：[docs/SECURITY.md](docs/SECURITY.md)

## 🆘 获取帮助

如果遇到问题：

1. 查看日志文件
2. 检查配置文件
3. 阅读完整文档
4. 提交 Issue

---

**祝你使用愉快！** 🎉
