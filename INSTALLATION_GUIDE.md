# AI Diet Assistant 安装指南

## 快速开始

### 前置要求

在运行安装脚本之前，请确保系统满足以下要求：

1. **Go 语言** (>= 1.21)
2. **MySQL** (>= 8.0)
3. **Redis** (>= 6，可选，用于令牌黑名单管理)
4. **Nginx** (最新稳定版，生产环境必需)
5. **OpenSSL** (用于生成安全密钥)

### 一键安装

```bash
# 克隆项目
git clone https://github.com/Deepblue-Sky2333/Ai-Diet-Assistant.git
cd Ai-Diet-Assistant

# 运行安装脚本
./install.sh
```

安装脚本会自动：
- ✅ 检查系统依赖
- ✅ 生成安全配置文件
- ✅ 初始化数据库
- ✅ 构建应用
- ✅ 验证安装

### 安装过程

安装脚本会引导您完成以下配置：

#### 1. 数据库配置
```
MySQL 主机 [localhost]: 
MySQL 端口 [3306]: 
数据库用户 [diet_user]: 
数据库名称 [ai_diet_assistant]: 
密码: ********
确认密码: ********
```

#### 2. 服务器配置
```
服务器端口 [9090]: 
服务器模式 (debug/release) [release]: 
```

#### 3. Redis 配置（可选）
```
启用 Redis？(y/n) [y]: 
Redis 主机 [localhost]: 
Redis 端口 [6379]: 
Redis 密码（如无则留空）: 
Redis 数据库 [0]: 
```

## 手动安装

如果您希望手动安装，请按照以下步骤操作：

### 1. 安装依赖

#### macOS
```bash
brew install go mysql redis openssl nginx
brew services start mysql
brew services start redis
brew services start nginx
```

#### Ubuntu/Debian
```bash
sudo apt-get update
sudo apt-get install golang-go mysql-server redis-server openssl nginx
sudo systemctl start mysql
sudo systemctl start redis
sudo systemctl start nginx
```

#### CentOS/RHEL
```bash
sudo yum install golang mysql-server redis openssl nginx
sudo systemctl start mysqld
sudo systemctl start redis
sudo systemctl start nginx
```

### 2. 配置数据库

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

### 3. 配置应用

```bash
# 复制配置文件示例
cp configs/config.yaml.example configs/config.yaml

# 编辑配置文件
vim configs/config.yaml
```

### 4. 运行数据库迁移

```bash
# 使用迁移脚本
./scripts/run-migrations.sh

# 或手动运行
for file in migrations/*_up.sql; do
    mysql -u diet_user -p ai_diet_assistant < "$file"
done
```

### 5. 构建应用

```bash
# 下载依赖
go mod download
go mod tidy

# 创建目录
mkdir -p bin logs uploads

# 编译应用
go build -o bin/diet-assistant cmd/server/main.go
```

### 6. 启动应用

```bash
# 直接运行
./bin/diet-assistant

# 或使用脚本
./scripts/start.sh
```

## 验证安装

### 检查应用状态

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

### 检查日志

```bash
# 查看应用日志
tail -f logs/app.log

# 查看安装日志
cat install.log
```

## 常见问题

### 1. MySQL 连接失败

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

### 2. Go 版本过低

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

### 3. 端口已被占用

**问题**: `bind: address already in use`

**解决方案**:
```bash
# 查找占用端口的进程
lsof -i :9090

# 终止进程
kill -9 <PID>

# 或修改配置文件中的端口
vim .env
# 修改 SERVER_PORT=9091
```

### 4. 权限问题

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

### 5. Redis 未运行

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

## 配置说明

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

## 管理命令

```bash
# 启动服务
./scripts/start.sh

# 停止服务
./scripts/stop.sh

# 查看状态
./scripts/status.sh

# 查看日志
tail -f logs/app.log

# 重启服务
./scripts/stop.sh && ./scripts/start.sh
```

## 生产环境部署

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

**完整配置和高级功能**（负载均衡、限流、缓存等）请参考：
- [Nginx 配置指南](docs/NGINX_CONFIGURATION.md)

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

## 安全建议

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

## 获取帮助

- **文档**: 查看 `docs/` 目录
- **API 文档**: 查看 `docs/api/` 目录
- **问题反馈**: https://github.com/Deepblue-Sky2333/Ai-Diet-Assistant/issues

## 许可证

查看 LICENSE 文件了解详情。
