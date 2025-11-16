#!/bin/bash

# AI Diet Assistant 安装脚本
# 此脚本生成安全密钥并创建配置文件

set -e

# 输出颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # 无颜色

echo -e "${BLUE}==================================="
echo "AI Diet Assistant 安装程序"
echo "===================================${NC}"
echo ""

# 检查是否安装了 openssl
if ! command -v openssl &> /dev/null; then
    echo -e "${RED}错误：未安装 openssl！${NC}"
    echo "请先安装 openssl："
    echo "  - macOS: brew install openssl"
    echo "  - Ubuntu/Debian: sudo apt-get install openssl"
    exit 1
fi

# 生成强随机 JWT 密钥（48 字节 base64）
generate_jwt_secret() {
    openssl rand -base64 48 | tr -d '\n'
}

# 生成 AES 加密密钥（32 字节）
generate_aes_key() {
    openssl rand -base64 32 | head -c 32
}

# 验证密码（仅基本检查）
validate_password() {
    local password="$1"
    
    if [ -z "$password" ]; then
        echo -e "${RED}密码不能为空${NC}"
        return 1
    fi
    
    return 0
}

# 验证域名格式
validate_domain() {
    local domain="$1"
    
    # 去除空格
    domain=$(echo "$domain" | xargs)
    
    # 空值检查
    if [ -z "$domain" ]; then
        echo -e "${RED}域名不能为空${NC}"
        return 1
    fi
    
    # 检查是否为有效的 URL 格式
    if [[ ! "$domain" =~ ^https?:// ]]; then
        echo -e "${RED}域名必须以 http:// 或 https:// 开头${NC}"
        return 1
    fi
    
    # 检查通配符（但允许 *.domain.com 格式的子域名）
    if [[ "$domain" == "*" ]] || [[ "$domain" == "http://*" ]] || [[ "$domain" == "https://*" ]]; then
        echo -e "${RED}出于安全原因，不允许使用纯通配符域名${NC}"
        return 1
    fi
    
    return 0
}

# 验证模块路径
validate_module_path() {
    local path="$1"
    
    if [[ "$path" == *"yourusername"* ]] || [[ "$path" == *"example"* ]]; then
        echo -e "${RED}模块路径不能包含 'yourusername' 或 'example'${NC}"
        return 1
    fi
    
    if [[ ! "$path" =~ ^[a-zA-Z0-9._/-]+$ ]]; then
        echo -e "${RED}模块路径包含无效字符${NC}"
        return 1
    fi
    
    return 0
}

echo -e "${YELLOW}此脚本将生成安全密钥并创建配置文件。${NC}"
echo ""

# 检查 .env 是否已存在
if [ -f .env ]; then
    echo -e "${YELLOW}警告：.env 文件已存在！${NC}"
    read -p "是否要覆盖它？(y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "安装已取消。"
        exit 0
    fi
    # 备份现有的 .env
    cp .env .env.backup.$(date +%Y%m%d_%H%M%S)
    echo -e "${GREEN}已备份现有的 .env${NC}"
fi

# 生成安全密钥
echo -e "${BLUE}正在生成安全密钥...${NC}"
JWT_SECRET=$(generate_jwt_secret)
ENCRYPTION_KEY=$(generate_aes_key)
echo -e "${GREEN}✓ JWT 密钥已生成（48 字节 base64）${NC}"
echo -e "${GREEN}✓ AES 加密密钥已生成（32 字节）${NC}"
echo ""

# 获取数据库密码
echo -e "${BLUE}数据库配置${NC}"
echo "请输入数据库密码："
while true; do
    read -s -p "密码: " DB_PASSWORD
    echo
    if validate_password "$DB_PASSWORD"; then
        read -s -p "确认密码: " DB_PASSWORD_CONFIRM
        echo
        if [ "$DB_PASSWORD" == "$DB_PASSWORD_CONFIRM" ]; then
            echo -e "${GREEN}✓ 数据库密码已设置${NC}"
            break
        else
            echo -e "${RED}密码不匹配。请重试。${NC}"
        fi
    fi
done
echo ""

# 获取数据库配置
read -p "数据库主机 [localhost]: " DB_HOST
DB_HOST=${DB_HOST:-localhost}

read -p "数据库端口 [3306]: " DB_PORT
DB_PORT=${DB_PORT:-3306}

read -p "数据库用户 [diet_user]: " DB_USER
DB_USER=${DB_USER:-diet_user}

read -p "数据库名称 [ai_diet_assistant]: " DB_NAME
DB_NAME=${DB_NAME:-ai_diet_assistant}

echo -e "${GREEN}✓ 数据库配置已设置${NC}"
echo ""

# 获取 CORS 允许的来源
echo -e "${BLUE}CORS 配置${NC}"
echo "请输入允许的来源（逗号分隔，例如：http://example.com,http://app.example.com）："
echo "注意："
echo "  - 出于安全原因，不允许使用通配符（*）"
echo "  - HTTPS 由 Nginx 反向代理处理，此处配置内网地址即可"
while true; do
    read -p "允许的来源: " CORS_ORIGINS
    
    # 验证每个域名
    IFS=',' read -ra DOMAINS <<< "$CORS_ORIGINS"
    all_valid=true
    for d in "${DOMAINS[@]}"; do
        # 去除空格
        d=$(echo "$d" | xargs)
        if ! validate_domain "$d"; then
            all_valid=false
            break
        fi
    done
    
    if [ "$all_valid" = true ]; then
        echo -e "${GREEN}✓ CORS 来源已设置${NC}"
        break
    else
        echo "请输入有效的域名。"
    fi
done
echo ""

# 模块路径已固定
MODULE_PATH="github.com/Deepblue-Sky2333/Ai-Diet-Assistant"
echo -e "${BLUE}模块路径配置${NC}"
echo -e "模块路径: ${GREEN}${MODULE_PATH}${NC}"
echo ""

# 获取服务器配置
echo -e "${BLUE}服务器配置${NC}"
read -p "服务器端口 [9090]: " SERVER_PORT
SERVER_PORT=${SERVER_PORT:-9090}

read -p "服务器模式 (debug/release) [release]: " SERVER_MODE
SERVER_MODE=${SERVER_MODE:-release}

echo -e "${GREEN}✓ 服务器配置已设置${NC}"
echo ""

# 获取 Redis 配置（可选）
echo -e "${BLUE}Redis 配置（可选）${NC}"
echo "Redis 用于令牌黑名单管理（登出功能）。"
echo "如果禁用，系统将使用内存存储（不推荐用于生产环境）。"
read -p "启用 Redis？(y/n) [n]: " ENABLE_REDIS
ENABLE_REDIS=${ENABLE_REDIS:-n}

if [[ $ENABLE_REDIS =~ ^[Yy]$ ]]; then
    REDIS_ENABLED=true
    
    read -p "Redis 主机 [localhost]: " REDIS_HOST
    REDIS_HOST=${REDIS_HOST:-localhost}
    
    read -p "Redis 端口 [6379]: " REDIS_PORT
    REDIS_PORT=${REDIS_PORT:-6379}
    
    read -p "Redis 密码（如无则留空）: " REDIS_PASSWORD
    
    read -p "Redis 数据库 [0]: " REDIS_DB
    REDIS_DB=${REDIS_DB:-0}
    
    echo -e "${GREEN}✓ Redis 配置已设置${NC}"
else
    REDIS_ENABLED=false
    REDIS_HOST=localhost
    REDIS_PORT=6379
    REDIS_PASSWORD=
    REDIS_DB=0
    echo -e "${YELLOW}⚠ Redis 已禁用 - 使用内存存储（不推荐用于生产环境）${NC}"
fi
echo ""

# 创建 .env 文件
echo -e "${BLUE}正在创建 .env 文件...${NC}"
cat > .env << EOF
# Server Configuration
SERVER_PORT=${SERVER_PORT}
SERVER_MODE=${SERVER_MODE}

# Database Configuration
DB_HOST=${DB_HOST}
DB_PORT=${DB_PORT}
DB_USER=${DB_USER}
DB_PASSWORD=${DB_PASSWORD}
DB_NAME=${DB_NAME}

# JWT Configuration (auto-generated 48-byte base64 key)
JWT_SECRET=${JWT_SECRET}

# Encryption Key (auto-generated 32-byte key for AES-256)
ENCRYPTION_KEY=${ENCRYPTION_KEY}

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_MINUTE=100

# CORS Configuration
CORS_ALLOWED_ORIGINS=${CORS_ORIGINS}

# Redis Configuration (optional)
REDIS_ENABLED=false
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Security
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION=15m
PASSWORD_MIN_LENGTH=8

EOF

# 设置 .env 文件的安全权限
chmod 600 .env
echo -e "${GREEN}✓ .env 文件已创建，权限为 600${NC}"
echo ""

# 如果 config.yaml 存在则更新
if [ -f configs/config.yaml ]; then
    echo -e "${YELLOW}警告：configs/config.yaml 已存在！${NC}"
    read -p "是否要更新它？(y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # 备份现有配置
        cp configs/config.yaml configs/config.yaml.backup.$(date +%Y%m%d_%H%M%S)
        echo -e "${GREEN}已备份现有的 config.yaml${NC}"
        UPDATE_CONFIG=true
    else
        UPDATE_CONFIG=false
    fi
else
    UPDATE_CONFIG=true
fi

if [ "$UPDATE_CONFIG" = true ]; then
    echo -e "${BLUE}正在创建 configs/config.yaml...${NC}"
    
    # Convert comma-separated origins to YAML array
    IFS=',' read -ra DOMAINS <<< "$CORS_ORIGINS"
    YAML_ORIGINS=""
    for domain in "${DOMAINS[@]}"; do
        domain=$(echo "$domain" | xargs)
        YAML_ORIGINS="${YAML_ORIGINS}    - ${domain}\n"
    done
    
    cat > configs/config.yaml << EOF
server:
  port: ${SERVER_PORT}
  mode: ${SERVER_MODE}
  read_timeout: 30s
  write_timeout: 30s

database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  user: ${DB_USER}
  password: ${DB_PASSWORD}
  dbname: ${DB_NAME}
  charset: utf8mb4
  parse_time: true
  loc: Local
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600s

jwt:
  secret: ${JWT_SECRET}
  expire_hours: 1
  refresh_expire_hours: 168

encryption:
  aes_key: ${ENCRYPTION_KEY}

redis:
  enabled: ${REDIS_ENABLED}
  host: ${REDIS_HOST}
  port: ${REDIS_PORT}
  password: ${REDIS_PASSWORD}
  db: ${REDIS_DB}
  max_retries: 3
  pool_size: 10

rate_limit:
  enabled: true
  requests_per_minute: 100
  burst: 10

cors:
  allowed_origins:
$(echo -e "$YAML_ORIGINS")  allowed_methods:
    - GET
    - POST
    - PUT
    - DELETE
    - OPTIONS
  allowed_headers:
    - Origin
    - Content-Type
    - Authorization
  expose_headers:
    - Content-Length
  allow_credentials: true
  max_age: 12h

log:
  level: info
  format: json
  output: logs/app.log
  max_size: 100
  max_backups: 3
  max_age: 28
  compress: true

security:
  max_login_attempts: 5
  lockout_duration: 15m
  password_min_length: 8
  require_special_char: true
  require_number: true
  require_uppercase: true

upload:
  max_file_size: 10485760
  allowed_types:
    - image/jpeg
    - image/png
    - image/gif
  upload_path: uploads/

EOF
    
    chmod 600 configs/config.yaml
    echo -e "${GREEN}✓ configs/config.yaml 已创建，权限为 600${NC}"
fi
echo ""

# 验证所有必需的配置是否已设置
echo -e "${BLUE}正在验证配置...${NC}"

VERIFICATION_FAILED=false

# 检查 JWT 密钥长度
if [ ${#JWT_SECRET} -lt 32 ]; then
    echo -e "${RED}✗ JWT 密钥太短（< 32 字符）${NC}"
    VERIFICATION_FAILED=true
else
    echo -e "${GREEN}✓ JWT 密钥长度足够${NC}"
fi

# 检查 AES 密钥长度
if [ ${#ENCRYPTION_KEY} -ne 32 ]; then
    echo -e "${RED}✗ AES 加密密钥不是 32 字节${NC}"
    VERIFICATION_FAILED=true
else
    echo -e "${GREEN}✓ AES 加密密钥为 32 字节${NC}"
fi

# 检查是否使用示例值
if [[ "$JWT_SECRET" == *"example"* ]] || [[ "$JWT_SECRET" == *"change"* ]]; then
    echo -e "${RED}✗ JWT 密钥似乎是示例值${NC}"
    VERIFICATION_FAILED=true
else
    echo -e "${GREEN}✓ JWT 密钥不是示例值${NC}"
fi

if [[ "$ENCRYPTION_KEY" == *"example"* ]] || [[ "$ENCRYPTION_KEY" == *"change"* ]]; then
    echo -e "${RED}✗ 加密密钥似乎是示例值${NC}"
    VERIFICATION_FAILED=true
else
    echo -e "${GREEN}✓ 加密密钥不是示例值${NC}"
fi

# 检查 CORS 配置
if [[ "$CORS_ORIGINS" == *"*"* ]]; then
    echo -e "${RED}✗ CORS 配置包含通配符（*）${NC}"
    VERIFICATION_FAILED=true
else
    echo -e "${GREEN}✓ CORS 配置未使用通配符${NC}"
fi

# 检查模块路径
if [[ "$MODULE_PATH" == *"yourusername"* ]] || [[ "$MODULE_PATH" == *"example"* ]]; then
    echo -e "${RED}✗ 模块路径包含占位符值${NC}"
    VERIFICATION_FAILED=true
else
    echo -e "${GREEN}✓ 模块路径配置正确${NC}"
fi

echo ""

if [ "$VERIFICATION_FAILED" = true ]; then
    echo -e "${RED}配置验证失败！${NC}"
    echo "请检查上述错误并重新运行脚本。"
    exit 1
fi

# 更新模块路径配置
echo -e "${BLUE}正在更新模块路径配置...${NC}"

# 更新 configs/module.conf
MODULE_CONF="configs/module.conf"
if [ -f "$MODULE_CONF" ]; then
    # 备份现有的 module.conf
    cp "$MODULE_CONF" "${MODULE_CONF}.backup.$(date +%Y%m%d_%H%M%S)"
    echo -e "${GREEN}已备份现有的 configs/module.conf${NC}"
fi

# 在 module.conf 中更新 MODULE_PATH
# 使用更兼容的方式更新配置文件
if grep -q "^MODULE_PATH=" "$MODULE_CONF" 2>/dev/null; then
    # 如果存在，则替换
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS 需要提供备份扩展名
        sed -i '' "s|^MODULE_PATH=.*|MODULE_PATH=${MODULE_PATH}|" "$MODULE_CONF"
    else
        # Linux
        sed -i "s|^MODULE_PATH=.*|MODULE_PATH=${MODULE_PATH}|" "$MODULE_CONF"
    fi
else
    # 如果不存在，则添加
    echo "MODULE_PATH=${MODULE_PATH}" >> "$MODULE_CONF"
fi
echo -e "${GREEN}✓ 模块路径配置已更新为: ${MODULE_PATH}${NC}"

# 运行模块路径更新脚本
if [ -f "./scripts/update-module-path-auto.sh" ]; then
    echo -e "${BLUE}正在运行模块路径更新脚本...${NC}"
    echo ""
    
    # 使用自动更新脚本（无需交互）
    if bash ./scripts/update-module-path-auto.sh "$MODULE_PATH"; then
        echo ""
        echo -e "${GREEN}✓ 模块路径更新成功${NC}"
    else
        echo ""
        echo -e "${YELLOW}⚠ 模块路径更新失败${NC}"
        echo "您可以手动运行："
        echo "  ./scripts/update-module-path-auto.sh $MODULE_PATH"
        echo "或者："
        echo "  ./scripts/update-module-path.sh"
        VERIFICATION_FAILED=true
    fi
elif [ -f "./scripts/update-module-path.sh" ]; then
    echo -e "${BLUE}正在运行模块路径更新脚本（交互模式）...${NC}"
    echo ""
    echo -e "${YELLOW}请按照提示确认更新${NC}"
    
    if bash ./scripts/update-module-path.sh; then
        echo ""
        echo -e "${GREEN}✓ 模块路径更新成功${NC}"
    else
        echo ""
        echo -e "${YELLOW}⚠ 模块路径更新失败或被取消${NC}"
        echo "您可以稍后手动运行：./scripts/update-module-path.sh"
        VERIFICATION_FAILED=true
    fi
else
    echo -e "${YELLOW}警告：未找到模块路径更新脚本${NC}"
    echo "模块路径配置已保存到 configs/module.conf"
    echo "但未应用到 Go 文件。"
    echo ""
    echo "请手动更新："
    echo "  1. go.mod 中的 module 声明"
    echo "  2. 所有 .go 文件中的导入路径"
    echo "  3. 运行 go mod tidy"
fi
echo ""

if [ "$VERIFICATION_FAILED" = true ]; then
    echo -e "${YELLOW}==================================="
    echo "安装完成但有警告"
    echo "===================================${NC}"
    echo ""
    echo "请在部署到生产环境之前解决上述警告。"
else
    echo -e "${GREEN}==================================="
    echo "安装成功完成！"
    echo "===================================${NC}"
fi
# 配置前端
echo -e "${BLUE}正在配置前端...${NC}"

FRONTEND_DIR="web/frontend"
if [ -d "$FRONTEND_DIR" ]; then
    # 获取前端配置
    echo ""
    echo -e "${BLUE}前端配置${NC}"
    echo "前端可以运行在两种模式："
    echo "  1. 演示模式（无需后端，使用模拟数据）"
    echo "  2. 生产模式（连接到后端 API）"
    echo ""
    read -p "启用演示模式？(y/n) [n]: " ENABLE_DEMO
    ENABLE_DEMO=${ENABLE_DEMO:-n}
    
    if [[ $ENABLE_DEMO =~ ^[Yy]$ ]]; then
        DEMO_MODE="true"
        # 在演示模式下，仍然使用 localhost 连接后端
        FRONTEND_API_URL="http://localhost:${SERVER_PORT}/api/v1"
        echo -e "${GREEN}✓ 演示模式已启用（可用时将使用 localhost 后端）${NC}"
    else
        DEMO_MODE="false"
        echo ""
        echo "请输入生产环境的后端 API URL："
        echo "  - 本地开发：http://localhost:${SERVER_PORT}/api/v1"
        echo "  - 生产环境：https://api.yourdomain.com/api/v1"
        read -p "后端 API URL [http://localhost:${SERVER_PORT}/api/v1]: " FRONTEND_API_URL
        FRONTEND_API_URL=${FRONTEND_API_URL:-http://localhost:${SERVER_PORT}/api/v1}
        echo -e "${GREEN}✓ 生产模式已启用${NC}"
    fi
    
    # 创建前端 .env.local 文件
    cat > "${FRONTEND_DIR}/.env.local" << EOF
# AI Diet Assistant 前端配置
# 由 install.sh 自动生成

# 演示模式配置
# 设置为 'true' 以在无后端情况下运行（使用模拟数据）
# 设置为 'false' 以连接到真实后端 API
NEXT_PUBLIC_DEMO_MODE=${DEMO_MODE}

# 后端 API 配置
# 后端 API 服务器的 URL
NEXT_PUBLIC_API_URL=${FRONTEND_API_URL}

EOF
    
    chmod 600 "${FRONTEND_DIR}/.env.local"
    echo -e "${GREEN}✓ 前端配置已创建：${FRONTEND_DIR}/.env.local${NC}"
    echo ""
    
    # 检查 node_modules 是否存在
    if [ ! -d "${FRONTEND_DIR}/node_modules" ]; then
        echo -e "${YELLOW}前端依赖未安装。${NC}"
        echo "要安装前端依赖，请运行："
        echo "  cd ${FRONTEND_DIR} && npm install"
        echo ""
    fi
else
    echo -e "${YELLOW}在 ${FRONTEND_DIR} 未找到前端目录${NC}"
    echo "跳过前端配置。"
    echo ""
fi

echo ""
echo "已创建的配置文件："
echo "  - .env（权限：600）"
if [ "$UPDATE_CONFIG" = true ]; then
    echo "  - configs/config.yaml（权限：600）"
fi
echo "  - configs/module.conf（模块路径配置）"
if [ -f "${FRONTEND_DIR}/.env.local" ]; then
    echo "  - ${FRONTEND_DIR}/.env.local（权限：600）"
fi
echo ""
echo "下一步："
echo "1. 检查生成的配置文件"
echo "2. 如需要，更新 AI 提供商 API 密钥（在 .env 或 Web UI 中）"
echo "3. 运行数据库迁移：make migrate-up"
echo "4. 构建应用程序：make build"
echo "5. 启动后端：./scripts/start.sh"
if [ -d "$FRONTEND_DIR" ]; then
    echo "6. 安装前端依赖：cd ${FRONTEND_DIR} && npm install"
    echo "7. 启动前端：cd ${FRONTEND_DIR} && npm run dev"
    echo "8. 访问应用程序：http://localhost:3000"
fi
echo ""
echo -e "${YELLOW}重要提示：请保护好您的 .env 文件，切勿将其提交到版本控制系统！${NC}"
echo ""
