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
git clone https://github.com/Deepblue-Sky2333/Ai-Diet-Assistant/
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

**详细说明：** 查看 [快速开始指南](QUICKSTART.md)

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

### 🏗️ 技术栈

**后端：**
- Go 1.21+
- Gin Web Framework
- MySQL 8.0+
- Redis 6+ (可选)
- JWT 认证

**部署：**
- Nginx (反向代理和 CORS 处理)
- Systemd (服务管理)

### 📚 文档

- [快速开始](QUICKSTART.md) - 快速安装和部署
- [安装指南](INSTALLATION_GUIDE.md) - 详细安装步骤
- [API 文档](docs/api/README.md) - 完整 API 接口说明
- [Nginx 配置](docs/NGINX_CONFIGURATION.md) - Nginx 反向代理配置指南
- [安全最佳实践](docs/SECURITY.md) - 安全配置指南
- [错误码说明](docs/ERROR_CODES.md) - 错误码参考

### 🔧 开发

```bash
# 后端开发
make run

# 运行测试
go test ./...

# 代码检查
go vet ./...
golint ./...
```

### 🚀 生产部署

#### 系统要求

- Go 1.21+
- MySQL 8.0+
- Redis 6+ (可选，用于 Token 黑名单)
- Nginx (推荐，用于反向代理和 CORS 处理)

#### 部署步骤

1. **安装依赖**
   ```bash
   # 安装 Go, PostgreSQL, Redis
   # 参考 INSTALLATION_GUIDE.md
   ```

2. **配置应用**
   ```bash
   # 复制配置文件
   cp configs/config.yaml.example configs/config.yaml
   
   # 编辑配置文件
   vim configs/config.yaml
   ```

3. **初始化数据库**
   ```bash
   # 运行数据库迁移
   ./scripts/run-migrations.sh
   ```

4. **构建应用**
   ```bash
   # 构建二进制文件
   go build -o bin/diet-assistant cmd/server/main.go
   ```

5. **配置 Nginx**
   
   后端是纯 API 服务，需要 Nginx 处理：
   - 反向代理
   - CORS 跨域请求
   - SSL/TLS 终止
   - 负载均衡（可选）
   
   详细配置请参考：[Nginx 配置指南](docs/NGINX_CONFIGURATION.md)

6. **配置系统服务**
   ```bash
   # 复制服务文件
   sudo cp scripts/diet-assistant.service /etc/systemd/system/
   
   # 启动服务
   sudo systemctl enable diet-assistant
   sudo systemctl start diet-assistant
   ```

7. **验证部署**
   ```bash
   # 检查服务状态
   sudo systemctl status diet-assistant
   
   # 测试健康检查
   curl http://localhost:9090/health
   ```

#### Nginx 配置示例

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

完整配置和 HTTPS 设置请参考：[Nginx 配置指南](docs/NGINX_CONFIGURATION.md)

### 📝 项目结构

```
.
├── cmd/                    # 应用程序入口
│   └── server/            # 主服务器
├── internal/              # 私有应用代码
│   ├── handler/          # HTTP 处理器
│   ├── service/          # 业务逻辑层
│   ├── repository/       # 数据访问层
│   ├── middleware/       # 中间件
│   ├── model/            # 数据模型
│   ├── config/           # 配置管理
│   ├── database/         # 数据库连接
│   ├── ai/               # AI 服务集成
│   └── utils/            # 工具函数
├── configs/              # 配置文件
├── migrations/           # 数据库迁移
├── scripts/              # 部署和管理脚本
├── docs/                 # 文档
│   ├── api/             # API 文档
│   └── NGINX_CONFIGURATION.md  # Nginx 配置指南
└── bin/                  # 编译后的二进制文件
```

### 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

