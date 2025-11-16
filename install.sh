#!/bin/bash

# ============================================
# AI Diet Assistant 优化安装脚本
# ============================================
# 此脚本会自动检查依赖、生成配置、初始化数据库并验证安装
# This script checks dependencies, generates config, initializes database and verifies installation

set -e

# 颜色输出 / Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # 无颜色 / No color

# 全局变量 / Global variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="${SCRIPT_DIR}/install.log"
INSTALL_FAILED=false

# ============================================
# 日志函数 / Logging functions
# ============================================

log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOG_FILE"
}

log_info() {
    echo -e "${BLUE}[INFO]${NC} $*" | tee -a "$LOG_FILE"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $*" | tee -a "$LOG_FILE"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $*" | tee -a "$LOG_FILE"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*" | tee -a "$LOG_FILE"
}

print_header() {
    echo -e "${CYAN}========================================${NC}"
    echo -e "${CYAN}$1${NC}"
    echo -e "${CYAN}========================================${NC}"
}

print_step() {
    echo ""
    echo -e "${BLUE}[$1] $2${NC}"
    echo ""
}

# ============================================
# 工具函数 / Utility functions
# ============================================

command_exists() {
    command -v "$1" &> /dev/null
}

version_ge() {
    # 比较版本号 / Compare version numbers
    # 返回 0 如果 $1 >= $2 / Returns 0 if $1 >= $2
    printf '%s\n%s\n' "$2" "$1" | sort -V -C
}

# ============================================
# 依赖检查函数 / Dependency check functions
# ============================================

check_go() {
    log_info "检查 Go 安装..."
    
    if ! command_exists go; then
        log_error "Go 未安装"
        echo ""
        echo "请安装 Go 1.21 或更高版本："
        echo "  - 官方网站: https://golang.org/dl/"
        echo "  - macOS: brew install go"
        echo "  - Ubuntu/Debian: sudo apt-get install golang-go"
        echo "  - CentOS/RHEL: sudo yum install golang"
        return 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_success "Go 已安装: ${GO_VERSION}"
    
    # 检查 Go 版本 / Check Go version
    if ! version_ge "$GO_VERSION" "1.21"; then
        log_warning "Go 版本过低 (需要 >= 1.21)，建议升级"
    fi
    
    return 0
}

check_mysql() {
    log_info "检查 MySQL 安装..."
    
    # 检查 mysql 命令 / Check mysql command
    if ! command_exists mysql; then
        log_error "MySQL 客户端未安装"
        echo ""
        echo "请安装 MySQL："
        echo "  - 官方网站: https://dev.mysql.com/downloads/"
        echo "  - macOS: brew install mysql"
        echo "  - Ubuntu/Debian: sudo apt-get install mysql-server mysql-client"
        echo "  - CentOS/RHEL: sudo yum install mysql-server"
        return 1
    fi
    
    MYSQL_VERSION=$(mysql --version | awk '{print $5}' | sed 's/,//')
    log_success "MySQL 客户端已安装: ${MYSQL_VERSION}"
    
    # 检查 MySQL 服务是否运行 / Check if MySQL service is running
    if mysqladmin ping &> /dev/null; then
        log_success "MySQL 服务正在运行"
    else
        log_warning "MySQL 服务未运行或需要密码验证"
        echo ""
        echo "请启动 MySQL 服务："
        echo "  - macOS: brew services start mysql"
        echo "  - Ubuntu/Debian: sudo systemctl start mysql"
        echo "  - CentOS/RHEL: sudo systemctl start mysqld"
        echo ""
        read -p "是否继续安装？(y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            return 1
        fi
    fi
    
    return 0
}

check_redis() {
    log_info "检查 Redis 安装（可选）..."
    
    if ! command_exists redis-cli; then
        log_warning "Redis 未安装（可选功能）"
        echo ""
        echo "Redis 用于令牌黑名单管理。如果不安装，系统将使用内存存储。"
        echo ""
        echo "安装 Redis："
        echo "  - macOS: brew install redis"
        echo "  - Ubuntu/Debian: sudo apt-get install redis-server"
        echo "  - CentOS/RHEL: sudo yum install redis"
        echo ""
        read -p "是否继续安装（不使用 Redis）？(y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            return 1
        fi
        return 0
    fi
    
    REDIS_VERSION=$(redis-cli --version | awk '{print $2}')
    log_success "Redis 已安装: ${REDIS_VERSION}"
    
    # 检查 Redis 服务是否运行 / Check if Redis service is running
    if redis-cli ping &> /dev/null; then
        log_success "Redis 服务正在运行"
    else
        log_warning "Redis 服务未运行"
        echo ""
        echo "请启动 Redis 服务："
        echo "  - macOS: brew services start redis"
        echo "  - Ubuntu/Debian: sudo systemctl start redis"
        echo "  - CentOS/RHEL: sudo systemctl start redis"
        echo ""
        read -p "是否继续安装（不使用 Redis）？(y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            return 1
        fi
    fi
    
    return 0
}

check_openssl() {
    log_info "检查 OpenSSL 安装..."
    
    if ! command_exists openssl; then
        log_error "OpenSSL 未安装"
        echo ""
        echo "请安装 OpenSSL："
        echo "  - macOS: brew install openssl"
        echo "  - Ubuntu/Debian: sudo apt-get install openssl"
        echo "  - CentOS/RHEL: sudo yum install openssl"
        return 1
    fi
    
    OPENSSL_VERSION=$(openssl version | awk '{print $2}')
    log_success "OpenSSL 已安装: ${OPENSSL_VERSION}"
    
    return 0
}

# ============================================
# 配置生成函数 / Configuration generation functions
# ============================================

generate_jwt_secret() {
    openssl rand -base64 48 | tr -d '\n'
}

generate_aes_key() {
    openssl rand -base64 32 | head -c 32
}

create_config_files() {
    log_info "生成配置文件..."
    
    # 检查 .env 是否已存在 / Check if .env already exists
    if [ -f "${SCRIPT_DIR}/.env" ]; then
        log_warning ".env 文件已存在"
        read -p "是否覆盖现有配置？(y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "跳过配置文件生成"
            return 0
        fi
        # 备份现有配置 / Backup existing config
        cp "${SCRIPT_DIR}/.env" "${SCRIPT_DIR}/.env.backup.$(date +%Y%m%d_%H%M%S)"
        log_success "已备份现有 .env 文件"
    fi
    
    # 生成安全密钥 / Generate security keys
    log_info "生成安全密钥..."
    JWT_SECRET=$(generate_jwt_secret)
    ENCRYPTION_KEY=$(generate_aes_key)
    log_success "安全密钥已生成"
    
    # 获取数据库配置 / Get database configuration
    echo ""
    echo -e "${CYAN}数据库配置 / Database Configuration${NC}"
    echo ""
    
    read -p "MySQL 主机 [localhost]: " DB_HOST
    DB_HOST=${DB_HOST:-localhost}
    
    read -p "MySQL 端口 [3306]: " DB_PORT
    DB_PORT=${DB_PORT:-3306}
    
    read -p "数据库用户 [diet_user]: " DB_USER
    DB_USER=${DB_USER:-diet_user}
    
    read -p "数据库名称 [ai_diet_assistant]: " DB_NAME
    DB_NAME=${DB_NAME:-ai_diet_assistant}
    
    echo ""
    echo "请输入数据库密码："
    read -s -p "密码: " DB_PASSWORD
    echo
    read -s -p "确认密码: " DB_PASSWORD_CONFIRM
    echo
    
    if [ "$DB_PASSWORD" != "$DB_PASSWORD_CONFIRM" ]; then
        log_error "密码不匹配"
        return 1
    fi
    
    if [ -z "$DB_PASSWORD" ]; then
        log_error "密码不能为空"
        return 1
    fi
    
    log_success "数据库配置已设置"
    
    # 获取服务器配置 / Get server configuration
    echo ""
    echo -e "${CYAN}服务器配置 / Server Configuration${NC}"
    echo ""
    
    read -p "服务器端口 [9090]: " SERVER_PORT
    SERVER_PORT=${SERVER_PORT:-9090}
    
    read -p "服务器模式 (debug/release) [release]: " SERVER_MODE
    SERVER_MODE=${SERVER_MODE:-release}
    
    log_success "服务器配置已设置"
    
    # 获取 Redis 配置 / Get Redis configuration
    echo ""
    echo -e "${CYAN}Redis 配置 / Redis Configuration${NC}"
    echo ""
    
    if command_exists redis-cli && redis-cli ping &> /dev/null; then
        read -p "启用 Redis？(y/n) [y]: " ENABLE_REDIS
        ENABLE_REDIS=${ENABLE_REDIS:-y}
    else
        log_warning "Redis 未运行，将禁用 Redis 功能"
        ENABLE_REDIS=n
    fi
    
    if [[ $ENABLE_REDIS =~ ^[Yy]$ ]]; then
        REDIS_ENABLED=true
        
        read -p "Redis 主机 [localhost]: " REDIS_HOST
        REDIS_HOST=${REDIS_HOST:-localhost}
        
        read -p "Redis 端口 [6379]: " REDIS_PORT
        REDIS_PORT=${REDIS_PORT:-6379}
        
        read -p "Redis 密码（如无则留空）: " REDIS_PASSWORD
        
        read -p "Redis 数据库 [0]: " REDIS_DB
        REDIS_DB=${REDIS_DB:-0}
        
        log_success "Redis 配置已设置"
    else
        REDIS_ENABLED=false
        REDIS_HOST=localhost
        REDIS_PORT=6379
        REDIS_PASSWORD=
        REDIS_DB=0
        log_warning "Redis 已禁用 - 使用内存存储"
    fi
    
    # 创建 .env 文件 / Create .env file
    log_info "创建 .env 文件..."
    
    cat > "${SCRIPT_DIR}/.env" << EOF
# ============================================
# AI Diet Assistant Configuration
# ============================================
# 由 install.sh 自动生成 / Auto-generated by install.sh
# 生成时间 / Generated at: $(date)

# Server Configuration
SERVER_PORT=${SERVER_PORT}
SERVER_MODE=${SERVER_MODE}

# Database Configuration (MySQL)
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

# Redis Configuration
REDIS_ENABLED=${REDIS_ENABLED}
REDIS_HOST=${REDIS_HOST}
REDIS_PORT=${REDIS_PORT}
REDIS_PASSWORD=${REDIS_PASSWORD}
REDIS_DB=${REDIS_DB}

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Security
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION=15m
PASSWORD_MIN_LENGTH=8

EOF
    
    chmod 600 "${SCRIPT_DIR}/.env"
    log_success ".env 文件已创建（权限：600）"
    
    # 创建 config.yaml 文件 / Create config.yaml file
    if [ ! -f "${SCRIPT_DIR}/configs/config.yaml" ] || [[ $REPLY =~ ^[Yy]$ ]]; then
        log_info "创建 configs/config.yaml..."
        
        mkdir -p "${SCRIPT_DIR}/configs"
        
        cat > "${SCRIPT_DIR}/configs/config.yaml" << EOF
# ============================================
# AI Diet Assistant Configuration
# ============================================
# 由 install.sh 自动生成 / Auto-generated by install.sh

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
        
        chmod 600 "${SCRIPT_DIR}/configs/config.yaml"
        log_success "configs/config.yaml 已创建（权限：600）"
    fi
    
    # 导出环境变量供后续使用 / Export environment variables for later use
    export DB_HOST DB_PORT DB_USER DB_PASSWORD DB_NAME
    export SERVER_PORT SERVER_MODE
    export REDIS_ENABLED REDIS_HOST REDIS_PORT REDIS_PASSWORD REDIS_DB
    
    return 0
}

# ============================================
# 数据库初始化函数 / Database initialization functions
# ============================================

init_database() {
    log_info "初始化数据库..."
    
    # 检查数据库连接 / Check database connection
    log_info "测试数据库连接..."
    
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1;" &> /dev/null; then
        log_success "数据库连接成功"
    else
        log_error "无法连接到数据库"
        echo ""
        echo "请检查："
        echo "  1. MySQL 服务是否运行"
        echo "  2. 数据库用户是否存在"
        echo "  3. 用户名和密码是否正确"
        echo "  4. 主机和端口是否正确"
        echo ""
        echo "创建数据库用户："
        echo "  mysql -u root -p"
        echo "  CREATE USER '${DB_USER}'@'localhost' IDENTIFIED BY 'your_password';"
        echo "  GRANT ALL PRIVILEGES ON ${DB_NAME}.* TO '${DB_USER}'@'localhost';"
        echo "  FLUSH PRIVILEGES;"
        return 1
    fi
    
    # 创建数据库 / Create database
    log_info "创建数据库 ${DB_NAME}..."
    
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "CREATE DATABASE IF NOT EXISTS ${DB_NAME} CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" &> /dev/null; then
        log_success "数据库 ${DB_NAME} 创建成功"
    else
        log_warning "数据库可能已存在"
    fi
    
    # 运行迁移脚本 / Run migration scripts
    log_info "运行数据库迁移..."
    
    if [ ! -d "${SCRIPT_DIR}/migrations" ]; then
        log_error "migrations 目录不存在"
        return 1
    fi
    
    migration_count=0
    for migration in "${SCRIPT_DIR}"/migrations/*_up.sql; do
        if [ -f "$migration" ]; then
            migration_name=$(basename "$migration")
            log_info "执行迁移: ${migration_name}"
            
            if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$migration" 2>&1 | grep -v "Warning: Using a password" > /dev/null; then
                log_success "迁移 ${migration_name} 执行成功"
                ((migration_count++))
            else
                log_warning "迁移 ${migration_name} 执行失败（可能已执行过）"
            fi
        fi
    done
    
    if [ $migration_count -eq 0 ]; then
        log_warning "没有执行任何迁移"
    else
        log_success "共执行 ${migration_count} 个迁移"
    fi
    
    return 0
}

# ============================================
# 构建应用函数 / Build application functions
# ============================================

build_application() {
    log_info "构建应用..."
    
    # 创建必要的目录 / Create necessary directories
    log_info "创建必要的目录..."
    mkdir -p "${SCRIPT_DIR}/bin"
    mkdir -p "${SCRIPT_DIR}/logs"
    mkdir -p "${SCRIPT_DIR}/uploads"
    log_success "目录创建完成"
    
    # 下载 Go 依赖 / Download Go dependencies
    log_info "下载 Go 依赖..."
    cd "${SCRIPT_DIR}"
    
    if go mod download; then
        log_success "Go 依赖下载成功"
    else
        log_error "Go 依赖下载失败"
        return 1
    fi
    
    if go mod tidy; then
        log_success "Go 依赖整理完成"
    else
        log_warning "Go 依赖整理有警告"
    fi
    
    # 构建应用 / Build application
    log_info "编译应用..."
    
    if go build -o "${SCRIPT_DIR}/bin/diet-assistant" "${SCRIPT_DIR}/cmd/server/main.go"; then
        log_success "应用编译成功"
    else
        log_error "应用编译失败"
        return 1
    fi
    
    # 设置可执行权限 / Set executable permission
    chmod +x "${SCRIPT_DIR}/bin/diet-assistant"
    
    return 0
}

# ============================================
# 验证安装函数 / Verify installation functions
# ============================================

verify_installation() {
    log_info "验证安装..."
    
    local verification_failed=false
    
    # 检查配置文件 / Check configuration files
    if [ -f "${SCRIPT_DIR}/.env" ]; then
        log_success ".env 文件存在"
    else
        log_error ".env 文件不存在"
        verification_failed=true
    fi
    
    if [ -f "${SCRIPT_DIR}/configs/config.yaml" ]; then
        log_success "config.yaml 文件存在"
    else
        log_error "config.yaml 文件不存在"
        verification_failed=true
    fi
    
    # 检查可执行文件 / Check executable
    if [ -f "${SCRIPT_DIR}/bin/diet-assistant" ] && [ -x "${SCRIPT_DIR}/bin/diet-assistant" ]; then
        log_success "应用可执行文件存在"
    else
        log_error "应用可执行文件不存在或不可执行"
        verification_failed=true
    fi
    
    # 检查数据库连接 / Check database connection
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "SELECT 1;" &> /dev/null; then
        log_success "数据库连接正常"
    else
        log_error "数据库连接失败"
        verification_failed=true
    fi
    
    # 检查必要的目录 / Check necessary directories
    for dir in logs uploads; do
        if [ -d "${SCRIPT_DIR}/${dir}" ]; then
            log_success "${dir} 目录存在"
        else
            log_error "${dir} 目录不存在"
            verification_failed=true
        fi
    done
    
    if [ "$verification_failed" = true ]; then
        log_error "安装验证失败"
        return 1
    fi
    
    log_success "安装验证通过"
    return 0
}

# ============================================
# 主安装流程 / Main installation process
# ============================================

main() {
    # 清空日志文件 / Clear log file
    > "$LOG_FILE"
    
    print_header "AI Diet Assistant 安装程序"
    log "安装开始于: $(date)"
    log "安装目录: ${SCRIPT_DIR}"
    echo ""
    
    # 步骤 1: 检查依赖 / Step 1: Check dependencies
    print_step "1/5" "检查系统依赖"
    
    if ! check_go; then
        log_error "Go 检查失败"
        INSTALL_FAILED=true
    fi
    
    if ! check_mysql; then
        log_error "MySQL 检查失败"
        INSTALL_FAILED=true
    fi
    
    if ! check_redis; then
        log_warning "Redis 检查失败（可选）"
    fi
    
    if ! check_openssl; then
        log_error "OpenSSL 检查失败"
        INSTALL_FAILED=true
    fi
    
    if [ "$INSTALL_FAILED" = true ]; then
        echo ""
        log_error "依赖检查失败，请安装缺失的依赖后重试"
        exit 1
    fi
    
    log_success "所有必需依赖已安装"
    
    # 步骤 2: 生成配置文件 / Step 2: Generate configuration files
    print_step "2/5" "生成配置文件"
    
    if ! create_config_files; then
        log_error "配置文件生成失败"
        exit 1
    fi
    
    # 步骤 3: 初始化数据库 / Step 3: Initialize database
    print_step "3/5" "初始化数据库"
    
    if ! init_database; then
        log_error "数据库初始化失败"
        exit 1
    fi
    
    # 步骤 4: 构建应用 / Step 4: Build application
    print_step "4/5" "构建应用"
    
    if ! build_application; then
        log_error "应用构建失败"
        exit 1
    fi
    
    # 步骤 5: 验证安装 / Step 5: Verify installation
    print_step "5/5" "验证安装"
    
    if ! verify_installation; then
        log_error "安装验证失败"
        exit 1
    fi
    
    # 安装完成 / Installation complete
    echo ""
    print_header "安装成功完成！"
    echo ""
    
    log_success "安装日志已保存到: ${LOG_FILE}"
    echo ""
    
    echo -e "${CYAN}应用信息：${NC}"
    echo "  访问地址: http://localhost:${SERVER_PORT}"
    echo "  健康检查: http://localhost:${SERVER_PORT}/health"
    echo "  配置文件: .env, configs/config.yaml"
    echo "  日志目录: logs/"
    echo "  上传目录: uploads/"
    echo ""
    
    echo -e "${CYAN}启动应用：${NC}"
    echo "  直接运行: ./bin/diet-assistant"
    echo "  使用脚本: ./scripts/start.sh"
    echo ""
    
    echo -e "${CYAN}管理命令：${NC}"
    echo "  查看状态: ./scripts/status.sh"
    echo "  停止服务: ./scripts/stop.sh"
    echo "  查看日志: tail -f logs/app.log"
    echo ""
    
    echo -e "${YELLOW}重要提示：${NC}"
    echo "  1. 请妥善保管 .env 文件，不要提交到版本控制"
    echo "  2. 建议配置 Nginx 反向代理处理 HTTPS 和 CORS"
    echo "  3. 定期备份数据库"
    echo "  4. 生产环境建议启用 Redis"
    echo ""
    
    log "安装完成于: $(date)"
}

# 运行主函数 / Run main function
main "$@"
