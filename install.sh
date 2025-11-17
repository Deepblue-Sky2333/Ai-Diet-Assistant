#!/bin/bash

# ============================================
# AI Diet Assistant 优化安装脚本
# ============================================
# 智能检测已安装系统，提供重新安装、修复数据库、直接启动等选项
# Intelligent installation script with smart detection and multiple options

set -e

# ============================================
# 调试模式配置 / Debug mode configuration
# ============================================
# 使用方法: DEBUG=true ./install.sh
# Usage: DEBUG=true ./install.sh
#
# 调试模式功能 / Debug mode features:
# - 启用 set -x 显示所有执行的命令
# - 自定义 PS4 提示符显示文件名、行号和函数名
# - log_debug() 函数输出详细的调试信息到终端和日志文件
# - 记录所有关键步骤的执行跟踪信息
#
# 示例 / Examples:
#   DEBUG=true ./install.sh              # 启用调试模式运行安装
#   DEBUG=true ./install.sh 2>&1 | tee debug.log  # 保存调试输出到文件

DEBUG=${DEBUG:-false}

if [ "$DEBUG" = "true" ]; then
    # 启用详细的执行跟踪 / Enable detailed execution tracing
    set -x
    
    # 自定义 PS4 提示符，显示更详细的调试信息
    # Custom PS4 prompt for more detailed debug information
    # 格式: +(文件名:行号): 函数名(): 
    PS4='+(${BASH_SOURCE}:${LINENO}): ${FUNCNAME[0]:+${FUNCNAME[0]}(): }'
    
    echo "========================================" >&2
    echo "调试模式已启用 / Debug mode enabled" >&2
    echo "所有命令执行将被记录 / All command executions will be logged" >&2
    echo "========================================" >&2
    echo "" >&2
fi

# 颜色和图标 / Colors and Icons
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m'

# 图标 / Icons
ICON_CHECK="✓"
ICON_CROSS="✗"
ICON_ARROW="→"
ICON_STAR="★"
ICON_WARNING="⚠"
ICON_INFO="ℹ"

# 全局变量 / Global variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="${SCRIPT_DIR}/install.log"
CONFIG_FILE="${SCRIPT_DIR}/configs/config.yaml"
ENV_FILE="${SCRIPT_DIR}/.env"
INSTALL_FAILED=false

# 默认配置值 / Default configuration values
DEFAULT_DB_HOST="localhost"
DEFAULT_DB_PORT="3306"
DEFAULT_DB_USER="diet_user"
DEFAULT_DB_NAME="ai_diet_assistant"
DEFAULT_SERVER_PORT="9090"
DEFAULT_SERVER_MODE="release"
DEFAULT_REDIS_HOST="localhost"
DEFAULT_REDIS_PORT="6379"
DEFAULT_REDIS_DB="0"

# ============================================
# 日志和输出函数 / Logging and output functions
# ============================================

log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*" >> "$LOG_FILE"
}

log_debug() {
    # 仅在调试模式下输出详细信息 / Output detailed info only in debug mode
    if [ "$DEBUG" = "true" ]; then
        echo "[$(date +'%Y-%m-%d %H:%M:%S')] [DEBUG] $*" | tee -a "$LOG_FILE" >&2
    else
        echo "[$(date +'%Y-%m-%d %H:%M:%S')] [DEBUG] $*" >> "$LOG_FILE"
    fi
}

log_info() {
    echo -e "${BLUE}${ICON_INFO}${NC} $*" | tee -a "$LOG_FILE"
}

log_success() {
    echo -e "${GREEN}${ICON_CHECK}${NC} $*" | tee -a "$LOG_FILE"
}

log_warning() {
    echo -e "${YELLOW}${ICON_WARNING}${NC} $*" | tee -a "$LOG_FILE"
}

log_error() {
    echo -e "${RED}${ICON_CROSS}${NC} $*" | tee -a "$LOG_FILE"
}

print_header() {
    echo ""
    echo -e "${CYAN}╔════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║${NC}  $1"
    echo -e "${CYAN}╚════════════════════════════════════════╝${NC}"
    echo ""
}

print_step() {
    echo ""
    echo -e "${MAGENTA}${ICON_ARROW} [$1] $2${NC}"
    echo ""
}

print_progress() {
    echo -ne "${BLUE}${ICON_ARROW}${NC} $1...\r"
}

print_menu() {
    echo -e "${CYAN}请选择操作 / Please select an option:${NC}"
    echo ""
    echo "  1) 重新安装 (覆盖所有配置)"
    echo -e "  2) 修复数据库 (清空并重新初始化) ${RED}[危险操作]${NC}"
    echo "  3) 直接启动服务"
    echo "  4) 退出"
    echo ""
}

# ============================================
# 工具函数 / Utility functions
# ============================================

check_terminal() {
    if [ -t 1 ]; then
        log_success "检测到交互式终端"
        return 0
    else
        log_warning "非交互式终端，禁用颜色输出"
        # 禁用颜色
        RED=''
        GREEN=''
        YELLOW=''
        BLUE=''
        CYAN=''
        MAGENTA=''
        NC=''
        # 禁用图标
        ICON_CHECK="[OK]"
        ICON_CROSS="[X]"
        ICON_ARROW="->"
        ICON_STAR="*"
        ICON_WARNING="!"
        ICON_INFO="i"
        return 1
    fi
}

command_exists() {
    command -v "$1" &> /dev/null
}

version_ge() {
    printf '%s\n%s\n' "$2" "$1" | sort -V -C
}

is_installed() {
    [ -f "$CONFIG_FILE" ] && [ -f "$ENV_FILE" ]
}

mysql_cmd() {
    local host="$1"
    local port="$2"
    local user="$3"
    local password="$4"
    shift 4
    
    local error_file=$(mktemp)
    local result
    
    log "执行 MySQL 命令: mysql -h${host} -P${port} -u${user} [密码已隐藏] $*"
    log_debug "MySQL 连接参数 - 主机: $host, 端口: $port, 用户: $user, 密码长度: ${#password}"
    log_debug "MySQL 命令参数: $*"
    log_debug "错误输出文件: $error_file"
    
    if [ -z "$password" ]; then
        log_debug "使用无密码模式连接 MySQL"
        mysql -h"$host" -P"$port" -u"$user" "$@" 2>"$error_file"
        result=$?
    else
        log_debug "使用密码模式连接 MySQL"
        mysql -h"$host" -P"$port" -u"$user" -p"$password" "$@" 2>"$error_file"
        result=$?
    fi
    
    log_debug "MySQL 命令退出码: $result"
    
    if [ $result -ne 0 ]; then
        local error_msg=$(cat "$error_file" | grep -v "Warning: Using a password" | head -n 5)
        if [ -n "$error_msg" ]; then
            log "MySQL 命令执行失败 (退出码: $result)"
            log "错误详情: $error_msg"
            log_debug "完整错误输出: $(cat "$error_file")"
        else
            log "MySQL 命令执行失败 (退出码: $result)，无详细错误信息"
            log_debug "错误文件内容: $(cat "$error_file")"
        fi
    else
        log "MySQL 命令执行成功"
        log_debug "MySQL 命令成功执行，无错误"
    fi
    
    rm -f "$error_file"
    log_debug "已清理临时错误文件"
    return $result
}

# ============================================
# 依赖检查函数 / Dependency check functions
# ============================================

check_go() {
    print_progress "检查 Go 安装"
    
    if ! command_exists go; then
        log_error "Go 未安装"
        echo ""
        echo -e "${YELLOW}请安装 Go 1.21 或更高版本:${NC}"
        echo "  • 官方网站: https://golang.org/dl/"
        echo -e "  • macOS: ${CYAN}brew install go${NC}"
        echo -e "  • Ubuntu/Debian: ${CYAN}sudo apt-get install golang-go${NC}"
        echo -e "  • CentOS/RHEL: ${CYAN}sudo yum install golang${NC}"
        return 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_success "Go ${GO_VERSION} 已安装"
    
    if ! version_ge "$GO_VERSION" "1.21"; then
        log_warning "Go 版本过低 (需要 >= 1.21)，建议升级"
    fi
    
    return 0
}

check_mysql() {
    print_progress "检查 MySQL 安装"
    
    if ! command_exists mysql; then
        log_error "MySQL 客户端未安装"
        echo ""
        echo -e "${YELLOW}请安装 MySQL:${NC}"
        echo -e "  • macOS: ${CYAN}brew install mysql${NC}"
        echo -e "  • Ubuntu/Debian: ${CYAN}sudo apt-get install mysql-server mysql-client${NC}"
        echo -e "  • CentOS/RHEL: ${CYAN}sudo yum install mysql-server${NC}"
        return 1
    fi
    
    MYSQL_VERSION=$(mysql --version | awk '{print $5}' | sed 's/,//')
    log_success "MySQL ${MYSQL_VERSION} 已安装"
    
    if mysqladmin ping &> /dev/null; then
        log_success "MySQL 服务正在运行"
    else
        log_warning "MySQL 服务未运行或需要密码验证"
        echo ""
        echo -e "${YELLOW}请启动 MySQL 服务:${NC}"
        echo -e "  • macOS: ${CYAN}brew services start mysql${NC}"
        echo -e "  • Ubuntu/Debian: ${CYAN}sudo systemctl start mysql${NC}"
        echo -e "  • CentOS/RHEL: ${CYAN}sudo systemctl start mysqld${NC}"
        echo ""
        echo -n "是否继续? (y/n) "
        read -n 1 -r
        echo
        [[ ! $REPLY =~ ^[Yy]$ ]] && return 1
    fi
    
    return 0
}

check_redis() {
    print_progress "检查 Redis 安装 (可选)"
    
    if ! command_exists redis-cli; then
        log_warning "Redis 未安装 (可选功能)"
        echo ""
        echo -e "${CYAN}Redis 用于令牌黑名单管理。如果不安装，系统将使用内存存储。${NC}"
        echo ""
        echo -e "${YELLOW}安装 Redis:${NC}"
        echo -e "  • macOS: ${CYAN}brew install redis${NC}"
        echo -e "  • Ubuntu/Debian: ${CYAN}sudo apt-get install redis-server${NC}"
        echo -e "  • CentOS/RHEL: ${CYAN}sudo yum install redis${NC}"
        return 0
    fi
    
    REDIS_VERSION=$(redis-cli --version | awk '{print $2}')
    log_success "Redis ${REDIS_VERSION} 已安装"
    
    if redis-cli ping &> /dev/null; then
        log_success "Redis 服务正在运行"
    else
        log_warning "Redis 服务未运行"
    fi
    
    return 0
}

check_openssl() {
    print_progress "检查 OpenSSL 安装"
    
    if ! command_exists openssl; then
        log_error "OpenSSL 未安装"
        echo ""
        echo -e "${YELLOW}请安装 OpenSSL:${NC}"
        echo -e "  • macOS: ${CYAN}brew install openssl${NC}"
        echo -e "  • Ubuntu/Debian: ${CYAN}sudo apt-get install openssl${NC}"
        echo -e "  • CentOS/RHEL: ${CYAN}sudo yum install openssl${NC}"
        return 1
    fi
    
    OPENSSL_VERSION=$(openssl version | awk '{print $2}')
    log_success "OpenSSL ${OPENSSL_VERSION} 已安装"
    
    return 0
}

check_all_dependencies() {
    print_step "1/6" "检查系统依赖"
    
    local failed=false
    
    check_go || failed=true
    check_mysql || failed=true
    check_redis
    check_openssl || failed=true
    
    if [ "$failed" = true ]; then
        echo ""
        log_error "依赖检查失败，请安装缺失的依赖后重试"
        return 1
    fi
    
    echo ""
    log_success "所有必需依赖已安装"
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

read_with_default() {
    local prompt="$1"
    local default="$2"
    local value
    
    log_debug "read_with_default 调用 - 提示: '$prompt', 默认值: '$default'"
    log_debug "终端状态: $([ -t 0 ] && echo '交互式' || echo '非交互式')"
    
    echo -ne "${prompt} [${CYAN}${default}${NC}]: "
    
    log_debug "等待用户输入..."
    read -r value
    log_debug "用户输入完成，值: '${value}' (长度: ${#value})"
    
    local result="${value:-$default}"
    log_debug "返回值: '$result'"
    
    echo "$result"
}

read_password() {
    local prompt="$1"
    local password
    local password_confirm
    
    while true; do
        read -s -p "${prompt}: " password
        echo
        
        if [ -z "$password" ]; then
            log_error "密码不能为空"
            continue
        fi
        
        if [ ${#password} -lt 8 ]; then
            log_error "密码必须至少8个字符"
            continue
        fi
        
        read -s -p "确认密码: " password_confirm
        echo
        
        if [ "$password" != "$password_confirm" ]; then
            log_error "密码不匹配，请重试"
            continue
        fi
        
        echo "$password"
        return 0
    done
}

create_config_files() {
    print_step "2/6" "生成配置文件"
    log "开始生成配置文件"
    log_debug "进入 create_config_files 函数"
    
    # 备份现有配置 / Backup existing config
    if [ -f "$ENV_FILE" ]; then
        log "检测到现有配置文件，准备备份"
        log_debug "现有配置文件路径: $ENV_FILE"
        local backup_file="${ENV_FILE}.backup.$(date +%Y%m%d_%H%M%S)"
        cp "$ENV_FILE" "$backup_file"
        log_success "已备份现有配置到: $(basename "$backup_file")"
        log_debug "备份文件路径: $backup_file"
    fi
    
    # 生成安全密钥 / Generate security keys
    print_progress "生成安全密钥"
    log "开始生成 JWT 密钥和加密密钥"
    log_debug "调用 generate_jwt_secret 函数"
    JWT_SECRET=$(generate_jwt_secret)
    log_debug "调用 generate_aes_key 函数"
    ENCRYPTION_KEY=$(generate_aes_key)
    log "JWT_SECRET 长度: ${#JWT_SECRET}"
    log "ENCRYPTION_KEY 长度: ${#ENCRYPTION_KEY}"
    log_debug "JWT_SECRET 前10字符: ${JWT_SECRET:0:10}..."
    log_debug "ENCRYPTION_KEY 前10字符: ${ENCRYPTION_KEY:0:10}..."
    log_success "安全密钥已生成"
    
    # 数据库配置 / Database configuration
    echo ""
    echo -e "${CYAN}━━━ 数据库配置 ━━━${NC}"
    echo ""
    log "开始数据库配置输入阶段"
    log_debug "准备读取数据库配置参数"
    
    log "等待用户输入: MySQL 主机"
    DB_HOST=$(read_with_default "MySQL 主机" "$DEFAULT_DB_HOST")
    log "MySQL 主机已设置为: $DB_HOST"
    
    log "等待用户输入: MySQL 端口"
    DB_PORT=$(read_with_default "MySQL 端口" "$DEFAULT_DB_PORT")
    log "MySQL 端口已设置为: $DB_PORT"
    
    log "等待用户输入: 数据库用户"
    DB_USER=$(read_with_default "数据库用户" "$DEFAULT_DB_USER")
    log "数据库用户已设置为: $DB_USER"
    
    log "等待用户输入: 数据库名称"
    DB_NAME=$(read_with_default "数据库名称" "$DEFAULT_DB_NAME")
    log "数据库名称已设置为: $DB_NAME"
    
    echo ""
    log "等待用户输入: 数据库密码"
    DB_PASSWORD=$(read_password "数据库密码")
    log "数据库密码已设置 (长度: ${#DB_PASSWORD})"
    
    log_success "数据库配置已设置"

    # 服务器配置 / Server configuration
    echo ""
    echo -e "${CYAN}━━━ 服务器配置 ━━━${NC}"
    echo ""
    log "开始服务器配置输入阶段"
    
    log "等待用户输入: 服务器端口"
    SERVER_PORT=$(read_with_default "服务器端口" "$DEFAULT_SERVER_PORT")
    log "服务器端口已设置为: $SERVER_PORT"
    
    log "等待用户输入: 服务器模式"
    SERVER_MODE=$(read_with_default "服务器模式 (debug/release)" "$DEFAULT_SERVER_MODE")
    log "服务器模式已设置为: $SERVER_MODE"
    
    log_success "服务器配置已设置"
    
    # Redis 配置 / Redis configuration
    echo ""
    echo -e "${CYAN}━━━ Redis 配置 ━━━${NC}"
    echo ""
    log "开始 Redis 配置输入阶段"
    
    if command_exists redis-cli && redis-cli ping &> /dev/null; then
        log "检测到 Redis 正在运行"
        log "等待用户输入: 是否启用 Redis"
        printf "%b" "启用 Redis? (y/n) [${CYAN}y${NC}]: "
        read -r ENABLE_REDIS
        ENABLE_REDIS=${ENABLE_REDIS:-y}
        log "用户选择: $ENABLE_REDIS"
    else
        log "Redis 未运行，自动禁用 Redis 功能"
        log_warning "Redis 未运行，将禁用 Redis 功能"
        ENABLE_REDIS=n
    fi
    
    if [[ $ENABLE_REDIS =~ ^[Yy]$ ]]; then
        log "用户选择启用 Redis，开始配置"
        REDIS_ENABLED=true
        
        log "等待用户输入: Redis 主机"
        REDIS_HOST=$(read_with_default "Redis 主机" "$DEFAULT_REDIS_HOST")
        log "Redis 主机已设置为: $REDIS_HOST"
        
        log "等待用户输入: Redis 端口"
        REDIS_PORT=$(read_with_default "Redis 端口" "$DEFAULT_REDIS_PORT")
        log "Redis 端口已设置为: $REDIS_PORT"
        
        log "等待用户输入: Redis 密码"
        read -p "Redis 密码 (如无则留空): " REDIS_PASSWORD
        log "Redis 密码已设置 (长度: ${#REDIS_PASSWORD})"
        
        log "等待用户输入: Redis 数据库"
        REDIS_DB=$(read_with_default "Redis 数据库" "$DEFAULT_REDIS_DB")
        log "Redis 数据库已设置为: $REDIS_DB"
        
        log_success "Redis 配置已设置"
    else
        log "用户选择禁用 Redis，使用默认值"
        REDIS_ENABLED=false
        REDIS_HOST=$DEFAULT_REDIS_HOST
        REDIS_PORT=$DEFAULT_REDIS_PORT
        REDIS_PASSWORD=""
        REDIS_DB=$DEFAULT_REDIS_DB
        log_warning "Redis 已禁用 - 使用内存存储"
    fi
    
    # 创建 .env 文件 / Create .env file
    print_progress "创建配置文件"
    log "开始创建 .env 文件: $ENV_FILE"
    
    cat > "$ENV_FILE" << EOF
# ============================================
# AI Diet Assistant Configuration
# ============================================
# 由 install.sh 自动生成 / Auto-generated by install.sh
# 生成时间: $(date)

# Server Configuration
SERVER_PORT=${SERVER_PORT}
SERVER_MODE=${SERVER_MODE}

# Database Configuration
DB_HOST=${DB_HOST}
DB_PORT=${DB_PORT}
DB_USER=${DB_USER}
DB_PASSWORD=${DB_PASSWORD}
DB_NAME=${DB_NAME}

# JWT Configuration
JWT_SECRET=${JWT_SECRET}

# Encryption Key
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
    
    log ".env 文件写入完成，设置权限为 600"
    chmod 600 "$ENV_FILE"
    log_success ".env 文件已创建 (权限: 600)"

    # 创建 config.yaml 文件 / Create config.yaml file
    log "创建 configs 目录"
    mkdir -p "${SCRIPT_DIR}/configs"
    log "开始创建 config.yaml 文件: $CONFIG_FILE"
    
    cat > "$CONFIG_FILE" << EOF
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
    
    log "config.yaml 文件写入完成，设置权限为 600"
    chmod 600 "$CONFIG_FILE"
    log_success "config.yaml 文件已创建 (权限: 600)"
    
    # 导出环境变量 / Export environment variables
    log "导出环境变量到当前 shell"
    export DB_HOST DB_PORT DB_USER DB_PASSWORD DB_NAME
    export SERVER_PORT SERVER_MODE
    export REDIS_ENABLED REDIS_HOST REDIS_PORT REDIS_PASSWORD REDIS_DB
    log "环境变量导出完成"
    
    echo ""
    log_success "配置文件生成完成"
    log "create_config_files 函数执行完成"
    return 0
}

# ============================================
# 数据库初始化函数 / Database initialization functions
# ============================================

load_config_vars() {
    # 从配置文件加载变量 / Load variables from config files
    if [ -f "$ENV_FILE" ]; then
        export $(grep -v '^#' "$ENV_FILE" | xargs)
    fi
}

init_database() {
    print_step "3/6" "初始化数据库"
    log "开始初始化数据库"
    
    # 测试数据库连接 / Test database connection
    print_progress "测试数据库连接"
    log "测试数据库连接: $DB_HOST:$DB_PORT (用户: $DB_USER)"
    
    local test_output=$(mktemp)
    if mysql_cmd "$DB_HOST" "$DB_PORT" "$DB_USER" "$DB_PASSWORD" -e "SELECT 1;" > "$test_output" 2>&1; then
        log "数据库连接测试成功"
        log_success "数据库连接成功"
        rm -f "$test_output"
    else
        local exit_code=$?
        local error_details=$(cat "$test_output" | grep -v "Warning: Using a password" | grep -i "error" | head -n 3)
        rm -f "$test_output"
        
        log "数据库连接测试失败 (退出码: $exit_code)"
        log "错误详情: $error_details"
        log_error "无法连接到数据库"
        
        echo ""
        echo -e "${RED}${ICON_CROSS} 数据库连接失败${NC}"
        echo ""
        
        if echo "$error_details" | grep -qi "access denied"; then
            echo -e "${YELLOW}可能的原因: 用户名或密码错误${NC}"
            echo ""
            echo "请检查:"
            echo "  • 数据库用户名是否正确: ${DB_USER}"
            echo "  • 数据库密码是否正确"
            echo "  • 用户是否有权限连接到 MySQL"
        elif echo "$error_details" | grep -qi "can't connect"; then
            echo -e "${YELLOW}可能的原因: 无法连接到 MySQL 服务器${NC}"
            echo ""
            echo "请检查:"
            echo "  • MySQL 服务是否正在运行"
            echo "  • 主机地址是否正确: ${DB_HOST}"
            echo "  • 端口是否正确: ${DB_PORT}"
            echo "  • 防火墙是否阻止连接"
        else
            echo -e "${YELLOW}请检查:${NC}"
            echo "  1. MySQL 服务是否运行"
            echo "  2. 数据库用户是否存在"
            echo "  3. 用户名和密码是否正确"
            echo "  4. 主机和端口配置是否正确"
        fi
        
        echo ""
        echo -e "${CYAN}创建数据库用户:${NC}"
        echo "  mysql -u root -p"
        echo "  CREATE USER '${DB_USER}'@'localhost' IDENTIFIED BY 'your_password';"
        echo "  GRANT ALL PRIVILEGES ON ${DB_NAME}.* TO '${DB_USER}'@'localhost';"
        echo "  FLUSH PRIVILEGES;"
        echo ""
        echo -e "${CYAN}启动 MySQL 服务:${NC}"
        echo "  • macOS: brew services start mysql"
        echo "  • Ubuntu/Debian: sudo systemctl start mysql"
        echo "  • CentOS/RHEL: sudo systemctl start mysqld"
        
        return 1
    fi
    
    # 创建数据库 / Create database
    print_progress "创建数据库 ${DB_NAME}"
    log "执行数据库创建命令: CREATE DATABASE IF NOT EXISTS ${DB_NAME}"
    
    local create_output=$(mktemp)
    if mysql_cmd "$DB_HOST" "$DB_PORT" "$DB_USER" "$DB_PASSWORD" \
        -e "CREATE DATABASE IF NOT EXISTS ${DB_NAME} CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" > "$create_output" 2>&1; then
        log "数据库创建成功"
        log_success "数据库 ${DB_NAME} 已就绪"
        rm -f "$create_output"
    else
        local exit_code=$?
        local error_details=$(cat "$create_output" | grep -v "Warning: Using a password" | grep -i "error")
        rm -f "$create_output"
        
        log "数据库创建失败 (退出码: $exit_code)"
        log "错误详情: $error_details"
        log_error "数据库创建失败"
        
        echo ""
        echo -e "${RED}${ICON_CROSS} 无法创建数据库 ${DB_NAME}${NC}"
        echo ""
        
        if echo "$error_details" | grep -qi "access denied"; then
            echo -e "${YELLOW}可能的原因: 用户 ${DB_USER} 没有创建数据库的权限${NC}"
            echo ""
            echo "解决方案:"
            echo "  1. 使用 root 用户授予权限:"
            echo "     mysql -u root -p"
            echo "     GRANT CREATE ON *.* TO '${DB_USER}'@'localhost';"
            echo "     FLUSH PRIVILEGES;"
            echo ""
            echo "  2. 或者手动创建数据库后重试:"
            echo "     mysql -u root -p"
            echo "     CREATE DATABASE ${DB_NAME} CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
            echo "     GRANT ALL PRIVILEGES ON ${DB_NAME}.* TO '${DB_USER}'@'localhost';"
        else
            echo -e "${YELLOW}请检查用户 ${DB_USER} 是否有创建数据库的权限${NC}"
        fi
        
        return 1
    fi
    
    # 运行迁移脚本 / Run migration scripts
    print_progress "运行数据库迁移"
    log "检查 migrations 目录"
    
    if [ ! -d "${SCRIPT_DIR}/migrations" ]; then
        log_error "migrations 目录不存在: ${SCRIPT_DIR}/migrations"
        echo ""
        echo -e "${RED}${ICON_CROSS} migrations 目录不存在${NC}"
        echo ""
        echo "请确保项目结构完整，migrations 目录应包含数据库迁移脚本"
        return 1
    fi
    
    log "开始执行数据库迁移脚本"
    local migration_count=0
    local migration_failed=0
    
    for migration in "${SCRIPT_DIR}"/migrations/*_up.sql; do
        if [ -f "$migration" ]; then
            local migration_name=$(basename "$migration")
            log "执行迁移: $migration_name"
            
            local migration_output=$(mktemp)
            if mysql_cmd "$DB_HOST" "$DB_PORT" "$DB_USER" "$DB_PASSWORD" "$DB_NAME" \
                < "$migration" > "$migration_output" 2>&1; then
                ((migration_count++))
                log "迁移成功: $migration_name"
            else
                local exit_code=$?
                local error_details=$(cat "$migration_output" | grep -v "Warning: Using a password" | grep -i "error" | head -n 3)
                
                # 检查是否是因为表已存在（这是正常的）
                if echo "$error_details" | grep -qi "table.*already exists"; then
                    log "迁移跳过（表已存在）: $migration_name"
                else
                    ((migration_failed++))
                    log "迁移失败: $migration_name (退出码: $exit_code)"
                    log "错误详情: $error_details"
                    
                    echo ""
                    log_error "迁移失败: $migration_name"
                    echo -e "${YELLOW}错误详情:${NC}"
                    echo "$error_details" | head -n 5
                fi
            fi
            rm -f "$migration_output"
        fi
    done
    
    log "数据库迁移完成，成功: ${migration_count}，失败: ${migration_failed}"
    
    if [ $migration_failed -gt 0 ]; then
        log_warning "数据库迁移完成，但有 ${migration_failed} 个迁移失败"
        echo ""
        echo -e "${YELLOW}${ICON_WARNING} 部分迁移失败，请检查日志: ${LOG_FILE}${NC}"
        echo ""
        read -p "是否继续安装? (y/n) " -n 1 -r
        echo
        [[ ! $REPLY =~ ^[Yy]$ ]] && return 1
    else
        log_success "数据库迁移完成 (执行了 ${migration_count} 个迁移)"
    fi
    
    echo ""
    log_success "数据库初始化完成"
    log "init_database 函数执行完成"
    return 0
}

repair_database() {
    print_header "修复数据库"
    log "========== 开始数据库修复流程 =========="
    
    echo -e "${RED}${ICON_WARNING} 警告: 此操作将清空数据库中的所有数据！${NC}"
    echo ""
    echo "这将会:"
    echo "  • 删除数据库 ${DB_NAME}"
    echo "  • 重新创建数据库"
    echo "  • 运行所有迁移脚本"
    echo "  • 所有现有数据将永久丢失"
    echo ""
    
    log "等待用户确认数据库修复操作"
    read -p "确定要继续吗? 输入 'yes' 确认: " confirm
    log "用户输入: $confirm"
    
    if [ "$confirm" != "yes" ]; then
        log "用户取消了数据库修复操作"
        log_info "操作已取消"
        return 1
    fi
    
    log "用户确认修复数据库，加载配置变量"
    load_config_vars
    
    # 删除数据库 / Drop database
    print_progress "删除现有数据库"
    log "执行数据库删除命令: DROP DATABASE IF EXISTS ${DB_NAME}"
    
    local drop_output=$(mktemp)
    if mysql_cmd "$DB_HOST" "$DB_PORT" "$DB_USER" "$DB_PASSWORD" \
        -e "DROP DATABASE IF EXISTS ${DB_NAME};" > "$drop_output" 2>&1; then
        log "数据库删除成功"
        log_success "数据库已删除"
        rm -f "$drop_output"
    else
        local exit_code=$?
        local error_details=$(cat "$drop_output" | grep -v "Warning: Using a password" | grep -i "error")
        rm -f "$drop_output"
        
        log "数据库删除失败 (退出码: $exit_code)"
        log "错误详情: $error_details"
        log_error "数据库删除失败"
        
        echo ""
        echo -e "${RED}${ICON_CROSS} 无法删除数据库 ${DB_NAME}${NC}"
        echo ""
        
        if echo "$error_details" | grep -qi "access denied"; then
            echo -e "${YELLOW}可能的原因: 用户 ${DB_USER} 没有删除数据库的权限${NC}"
            echo ""
            echo "解决方案:"
            echo "  使用 root 用户手动删除数据库:"
            echo "    mysql -u root -p"
            echo "    DROP DATABASE IF EXISTS ${DB_NAME};"
        else
            echo -e "${YELLOW}请检查用户 ${DB_USER} 是否有删除数据库的权限${NC}"
        fi
        
        return 1
    fi
    
    # 重新初始化 / Reinitialize
    log "开始重新初始化数据库"
    init_database
    
    log "========== 数据库修复流程完成 =========="
    return $?
}

# ============================================
# 构建应用函数 / Build application functions
# ============================================

build_application() {
    print_step "4/6" "构建应用"
    log "开始构建应用"
    
    # 创建必要的目录 / Create necessary directories
    print_progress "创建必要的目录"
    log "创建必要的目录: bin, logs, uploads"
    mkdir -p "${SCRIPT_DIR}/bin"
    mkdir -p "${SCRIPT_DIR}/logs"
    mkdir -p "${SCRIPT_DIR}/uploads"
    log "目录创建完成"
    log_success "目录创建完成"
    
    # 下载 Go 依赖 / Download Go dependencies
    print_progress "下载 Go 依赖"
    log "切换到项目目录: ${SCRIPT_DIR}"
    cd "${SCRIPT_DIR}"
    
    log "执行 go mod download"
    local download_output=$(mktemp)
    if go mod download > "$download_output" 2>&1; then
        log "Go 依赖下载成功"
        cat "$download_output" >> "$LOG_FILE"
        rm -f "$download_output"
        log_success "Go 依赖下载成功"
    else
        local exit_code=$?
        local error_details=$(cat "$download_output" | head -n 10)
        cat "$download_output" >> "$LOG_FILE"
        rm -f "$download_output"
        
        log "Go 依赖下载失败 (退出码: $exit_code)"
        log "错误详情: $error_details"
        log_error "Go 依赖下载失败"
        
        echo ""
        echo -e "${RED}${ICON_CROSS} Go 依赖下载失败${NC}"
        echo ""
        echo -e "${YELLOW}可能的原因:${NC}"
        echo "  • 网络连接问题"
        echo "  • Go 模块代理不可访问"
        echo "  • go.mod 文件配置错误"
        echo ""
        echo -e "${CYAN}建议:${NC}"
        echo "  • 检查网络连接"
        echo "  • 设置 Go 代理: export GOPROXY=https://goproxy.cn,direct"
        echo "  • 查看详细日志: ${LOG_FILE}"
        
        return 1
    fi
    
    log "执行 go mod tidy"
    local tidy_output=$(mktemp)
    if go mod tidy > "$tidy_output" 2>&1; then
        log "Go 依赖整理完成"
        cat "$tidy_output" >> "$LOG_FILE"
        rm -f "$tidy_output"
        log_success "Go 依赖整理完成"
    else
        local error_details=$(cat "$tidy_output")
        cat "$tidy_output" >> "$LOG_FILE"
        rm -f "$tidy_output"
        
        log "Go 依赖整理有警告或失败"
        log "详情: $error_details"
        log_warning "Go 依赖整理有警告"
    fi
    
    # 构建应用 / Build application
    print_progress "编译应用"
    log "开始编译应用: ${SCRIPT_DIR}/cmd/server/main.go -> ${SCRIPT_DIR}/bin/diet-assistant"
    
    local build_output=$(mktemp)
    if go build -o "${SCRIPT_DIR}/bin/diet-assistant" "${SCRIPT_DIR}/cmd/server/main.go" > "$build_output" 2>&1; then
        log "应用编译成功"
        cat "$build_output" >> "$LOG_FILE"
        rm -f "$build_output"
        log_success "应用编译成功"
    else
        local exit_code=$?
        local error_details=$(cat "$build_output" | head -n 20)
        cat "$build_output" >> "$LOG_FILE"
        rm -f "$build_output"
        
        log "应用编译失败 (退出码: $exit_code)"
        log "错误详情: $error_details"
        log_error "应用编译失败"
        
        echo ""
        echo -e "${RED}${ICON_CROSS} 应用编译失败${NC}"
        echo ""
        echo -e "${YELLOW}编译错误详情:${NC}"
        echo "$error_details" | head -n 10
        echo ""
        echo -e "${CYAN}建议:${NC}"
        echo "  • 检查 Go 版本是否满足要求 (>= 1.21)"
        echo "  • 查看完整编译日志: ${LOG_FILE}"
        echo "  • 确保所有依赖已正确下载"
        
        return 1
    fi
    
    log "设置可执行权限"
    chmod +x "${SCRIPT_DIR}/bin/diet-assistant"
    
    echo ""
    log_success "应用构建完成"
    log "build_application 函数执行完成"
    return 0
}

# ============================================
# 初始用户创建函数 / Initial user creation functions
# ============================================

create_initial_user() {
    print_step "5/6" "创建初始用户"
    log "开始创建初始用户流程"
    
    load_config_vars
    log "配置变量已加载"
    
    # 检查是否已有用户 / Check if users exist
    log "检查数据库中是否已存在用户"
    
    local count_output=$(mktemp)
    local user_count
    
    if mysql_cmd "$DB_HOST" "$DB_PORT" "$DB_USER" "$DB_PASSWORD" "$DB_NAME" \
        -se "SELECT COUNT(*) FROM users" > "$count_output" 2>&1; then
        user_count=$(cat "$count_output" | grep -v "Warning: Using a password" | head -n 1)
        log "数据库中现有用户数: $user_count"
        rm -f "$count_output"
    else
        local error_details=$(cat "$count_output" | grep -v "Warning: Using a password" | grep -i "error")
        rm -f "$count_output"
        
        log "查询用户数失败，可能是 users 表不存在"
        log "错误详情: $error_details"
        
        if echo "$error_details" | grep -qi "doesn't exist"; then
            log_warning "users 表不存在，将继续创建初始用户"
            user_count=0
        else
            log_error "无法查询用户表"
            echo ""
            echo -e "${RED}${ICON_CROSS} 无法查询用户表${NC}"
            echo ""
            echo "请确保数据库迁移已成功执行"
            return 1
        fi
    fi
    
    if [ "$user_count" -gt 0 ]; then
        log_success "数据库中已存在 ${user_count} 个用户，跳过初始用户创建"
        return 0
    fi
    
    echo ""
    echo -e "${CYAN}━━━ 初始管理员用户设置 ━━━${NC}"
    echo ""
    echo "系统需要至少一个管理员用户才能登录。"
    echo "第一个创建的用户将自动成为管理员。"
    echo ""
    log "开始初始管理员用户输入阶段"
    
    # 获取用户名 / Get username
    log "等待用户输入: 用户名"
    while true; do
        read -p "用户名 (3-50个字符): " INIT_USERNAME
        log "用户输入的用户名: $INIT_USERNAME (长度: ${#INIT_USERNAME})"
        
        if [ -z "$INIT_USERNAME" ]; then
            log "用户名验证失败: 用户名为空"
            log_error "用户名不能为空"
            continue
        fi
        
        if [ ${#INIT_USERNAME} -lt 3 ] || [ ${#INIT_USERNAME} -gt 50 ]; then
            log "用户名验证失败: 长度不符合要求 (${#INIT_USERNAME})"
            log_error "用户名必须是3-50个字符"
            continue
        fi
        
        if ! echo "$INIT_USERNAME" | grep -qE '^[a-zA-Z0-9]+$'; then
            log "用户名验证失败: 包含非法字符"
            log_error "用户名只能包含字母和数字"
            continue
        fi
        
        log "用户名验证通过: $INIT_USERNAME"
        break
    done
    
    # 获取密码 / Get password
    echo ""
    log "等待用户输入: 管理员密码"
    INIT_PASSWORD=$(read_password "管理员密码 (至少8个字符)")
    log "管理员密码已设置 (长度: ${#INIT_PASSWORD})"
    
    # 获取邮箱 / Get email
    echo ""
    log "等待用户输入: 电子邮件 (可选)"
    read -p "电子邮件 (可选，按回车跳过): " INIT_EMAIL
    log "电子邮件已设置: ${INIT_EMAIL:-<未提供>}"

    # 构建 create-user 工具 / Build create-user tool
    print_progress "构建用户创建工具"
    log "检查 create-user 工具源码"
    
    if [ ! -f "${SCRIPT_DIR}/cmd/create-user/main.go" ]; then
        log_error "create-user 工具源码不存在: ${SCRIPT_DIR}/cmd/create-user/main.go"
        echo ""
        echo -e "${RED}${ICON_CROSS} create-user 工具源码不存在${NC}"
        echo ""
        echo "请确保项目结构完整"
        return 1
    fi
    
    log "开始构建 create-user 工具"
    local build_user_output=$(mktemp)
    
    if go build -o "${SCRIPT_DIR}/bin/create-user" "${SCRIPT_DIR}/cmd/create-user/main.go" > "$build_user_output" 2>&1; then
        chmod +x "${SCRIPT_DIR}/bin/create-user"
        log "create-user 工具构建成功，已设置可执行权限"
        cat "$build_user_output" >> "$LOG_FILE"
        rm -f "$build_user_output"
        log_success "用户创建工具构建成功"
    else
        local exit_code=$?
        local error_details=$(cat "$build_user_output" | head -n 15)
        cat "$build_user_output" >> "$LOG_FILE"
        rm -f "$build_user_output"
        
        log "create-user 工具构建失败 (退出码: $exit_code)"
        log "错误详情: $error_details"
        log_error "用户创建工具构建失败"
        
        echo ""
        echo -e "${RED}${ICON_CROSS} 用户创建工具构建失败${NC}"
        echo ""
        echo -e "${YELLOW}编译错误详情:${NC}"
        echo "$error_details" | head -n 10
        echo ""
        echo "查看完整日志: ${LOG_FILE}"
        
        return 1
    fi
    
    # 创建用户 / Create user
    print_progress "创建初始管理员用户"
    log "准备创建初始管理员用户"
    
    local create_cmd="${SCRIPT_DIR}/bin/create-user -username \"${INIT_USERNAME}\" -password \"${INIT_PASSWORD}\""
    
    if [ -n "$INIT_EMAIL" ]; then
        create_cmd="${create_cmd} -email \"${INIT_EMAIL}\""
        log "包含邮箱参数"
    fi
    
    create_cmd="${create_cmd} -config \"${CONFIG_FILE}\""
    log "执行命令: ${SCRIPT_DIR}/bin/create-user -username \"${INIT_USERNAME}\" -password \"***\" -config \"${CONFIG_FILE}\""
    
    local create_user_output=$(mktemp)
    if eval "$create_cmd" > "$create_user_output" 2>&1; then
        log "用户创建命令执行成功"
        cat "$create_user_output" >> "$LOG_FILE"
        rm -f "$create_user_output"
        log_success "初始管理员用户创建成功"
        
        # 初始化系统设置 / Initialize system settings
        log "初始化系统设置"
        local settings_output=$(mktemp)
        
        if mysql_cmd "$DB_HOST" "$DB_PORT" "$DB_USER" "$DB_PASSWORD" "$DB_NAME" <<EOF > "$settings_output" 2>&1
INSERT INTO system_settings (setting_key, setting_value, description) 
VALUES ('registration_enabled', 'true', '是否允许新用户注册')
ON DUPLICATE KEY UPDATE setting_key = setting_key;
EOF
        then
            log "系统设置初始化完成"
            log_success "系统设置初始化完成"
            rm -f "$settings_output"
        else
            local error_details=$(cat "$settings_output" | grep -v "Warning: Using a password" | grep -i "error")
            log "系统设置初始化失败"
            log "错误详情: $error_details"
            rm -f "$settings_output"
            log_warning "系统设置初始化失败，但可以继续"
        fi
    else
        local exit_code=$?
        local error_details=$(cat "$create_user_output")
        cat "$create_user_output" >> "$LOG_FILE"
        rm -f "$create_user_output"
        
        log "用户创建命令执行失败 (退出码: $exit_code)"
        log "错误详情: $error_details"
        log_error "初始用户创建失败"
        
        echo ""
        echo -e "${RED}${ICON_CROSS} 初始用户创建失败${NC}"
        echo ""
        
        if echo "$error_details" | grep -qi "duplicate"; then
            echo -e "${YELLOW}可能的原因: 用户名已存在${NC}"
            echo ""
            echo "用户 ${INIT_USERNAME} 可能已经存在，您可以:"
            echo "  1. 使用现有用户登录"
            echo "  2. 选择不同的用户名"
        elif echo "$error_details" | grep -qi "connection"; then
            echo -e "${YELLOW}可能的原因: 数据库连接问题${NC}"
            echo ""
            echo "请检查数据库连接是否正常"
        else
            echo -e "${YELLOW}详细错误信息已记录到日志文件: ${LOG_FILE}${NC}"
        fi
        
        echo ""
        echo "您可以在启动服务器后使用以下命令手动创建用户:"
        echo -e "  ${CYAN}./bin/create-user -username <用户名> -password <密码> [-email <邮箱>]${NC}"
        echo ""
        
        log "等待用户确认是否继续"
        read -p "是否继续? (y/n) " -n 1 -r
        echo
        log "用户选择: $REPLY"
        [[ ! $REPLY =~ ^[Yy]$ ]] && return 1
    fi
    
    echo ""
    log_success "初始用户设置完成"
    log "create_initial_user 函数执行完成"
    return 0
}

# ============================================
# 验证安装函数 / Verify installation functions
# ============================================

verify_installation() {
    print_step "6/6" "验证安装"
    
    local failed=false
    
    # 检查配置文件 / Check configuration files
    print_progress "检查配置文件"
    
    [ -f "$ENV_FILE" ] && log_success ".env 文件存在" || { log_error ".env 文件不存在"; failed=true; }
    [ -f "$CONFIG_FILE" ] && log_success "config.yaml 文件存在" || { log_error "config.yaml 文件不存在"; failed=true; }
    
    # 检查可执行文件 / Check executable
    print_progress "检查可执行文件"
    
    if [ -f "${SCRIPT_DIR}/bin/diet-assistant" ] && [ -x "${SCRIPT_DIR}/bin/diet-assistant" ]; then
        log_success "应用可执行文件存在"
    else
        log_error "应用可执行文件不存在或不可执行"
        failed=true
    fi
    
    # 检查数据库连接 / Check database connection
    print_progress "检查数据库连接"
    
    load_config_vars
    
    local verify_output=$(mktemp)
    if mysql_cmd "$DB_HOST" "$DB_PORT" "$DB_USER" "$DB_PASSWORD" "$DB_NAME" -e "SELECT 1;" > "$verify_output" 2>&1; then
        log_success "数据库连接正常"
        rm -f "$verify_output"
    else
        local error_details=$(cat "$verify_output" | grep -v "Warning: Using a password" | grep -i "error")
        rm -f "$verify_output"
        
        log "数据库连接验证失败"
        log "错误详情: $error_details"
        log_error "数据库连接失败"
        
        echo ""
        echo -e "${YELLOW}数据库连接问题，请检查配置${NC}"
        
        failed=true
    fi
    
    # 检查必要的目录 / Check necessary directories
    print_progress "检查必要的目录"
    
    for dir in logs uploads bin; do
        if [ -d "${SCRIPT_DIR}/${dir}" ]; then
            log_success "${dir}/ 目录存在"
        else
            log_error "${dir}/ 目录不存在"
            failed=true
        fi
    done
    
    echo ""
    
    if [ "$failed" = true ]; then
        log_error "安装验证失败"
        return 1
    fi
    
    log_success "安装验证通过"
    return 0
}

# ============================================
# 启动服务函数 / Start service functions
# ============================================

start_service() {
    print_header "启动服务"
    
    load_config_vars
    
    # 检查服务是否已在运行 / Check if service is already running
    if pgrep -f "diet-assistant" > /dev/null; then
        log_warning "服务似乎已在运行"
        echo ""
        read -p "是否停止现有服务并重新启动? (y/n) " -n 1 -r
        echo
        
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            pkill -f "diet-assistant"
            sleep 2
        else
            return 0
        fi
    fi
    
    # 启动服务 / Start service
    print_progress "启动服务"
    
    if [ -f "${SCRIPT_DIR}/scripts/start.sh" ]; then
        "${SCRIPT_DIR}/scripts/start.sh"
    else
        nohup "${SCRIPT_DIR}/bin/diet-assistant" > "${SCRIPT_DIR}/logs/app.log" 2>&1 &
        sleep 2
        
        if pgrep -f "diet-assistant" > /dev/null; then
            log_success "服务已启动"
        else
            log_error "服务启动失败"
            return 1
        fi
    fi
    
    echo ""
    log_success "服务启动完成"
    
    # 显示访问信息 / Show access information
    echo ""
    echo -e "${CYAN}━━━ 服务访问信息 ━━━${NC}"
    echo ""
    echo -e "  ${ICON_STAR} 访问地址: ${GREEN}http://localhost:${SERVER_PORT}${NC}"
    echo -e "  ${ICON_STAR} 健康检查: ${GREEN}http://localhost:${SERVER_PORT}/health${NC}"
    echo ""
    
    return 0
}

# ============================================
# 安装报告函数 / Installation report functions
# ============================================

generate_install_report() {
    local report_file="${SCRIPT_DIR}/install_report.txt"
    
    cat > "$report_file" << EOF
╔════════════════════════════════════════════════════════════════╗
║           AI Diet Assistant 安装报告                           ║
╚════════════════════════════════════════════════════════════════╝

安装时间: $(date)
安装目录: ${SCRIPT_DIR}

━━━ 配置信息 ━━━

服务器:
  • 端口: ${SERVER_PORT}
  • 模式: ${SERVER_MODE}
  • 访问地址: http://localhost:${SERVER_PORT}

数据库:
  • 主机: ${DB_HOST}:${DB_PORT}
  • 数据库: ${DB_NAME}
  • 用户: ${DB_USER}

Redis:
  • 状态: $([ "$REDIS_ENABLED" = "true" ] && echo "已启用" || echo "已禁用")
$([ "$REDIS_ENABLED" = "true" ] && echo "  • 主机: ${REDIS_HOST}:${REDIS_PORT}")

━━━ 文件位置 ━━━

  • 可执行文件: ./bin/diet-assistant
  • 配置文件: ./configs/config.yaml, ./.env
  • 日志目录: ./logs/
  • 上传目录: ./uploads/
  • 安装日志: ./install.log

━━━ 常用命令 ━━━

启动服务:
  ./scripts/start.sh
  或: ./bin/diet-assistant

停止服务:
  ./scripts/stop.sh

查看状态:
  ./scripts/status.sh

查看日志:
  tail -f logs/app.log

创建用户:
  ./bin/create-user -username <用户名> -password <密码>

━━━ 下一步操作 ━━━

1. 启动服务: ./scripts/start.sh
2. 访问健康检查: curl http://localhost:${SERVER_PORT}/health
3. 查看 API 文档: docs/api/README.md
4. 配置 Nginx (生产环境): docs/NGINX_CONFIGURATION.md

━━━ 重要提示 ━━━

⚠ 请妥善保管 .env 文件，不要提交到版本控制
⚠ 建议配置 Nginx 反向代理处理 HTTPS
⚠ 定期备份数据库
⚠ 生产环境建议启用 Redis

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

EOF
    
    echo "$report_file"
}

print_installation_summary() {
    load_config_vars
    
    echo ""
    print_header "安装成功完成！"
    
    echo -e "${CYAN}━━━ 服务信息 ━━━${NC}"
    echo ""
    echo -e "  ${ICON_STAR} 访问地址: ${GREEN}http://localhost:${SERVER_PORT}${NC}"
    echo -e "  ${ICON_STAR} 健康检查: ${GREEN}http://localhost:${SERVER_PORT}/health${NC}"
    echo -e "  ${ICON_STAR} 配置文件: ${CYAN}.env, configs/config.yaml${NC}"
    echo ""
    
    echo -e "${CYAN}━━━ 启动应用 ━━━${NC}"
    echo ""
    echo -e "  ${CYAN}./scripts/start.sh${NC}      # 使用脚本启动"
    echo -e "  ${CYAN}./bin/diet-assistant${NC}    # 直接运行"
    echo ""
    
    echo -e "${CYAN}━━━ 管理命令 ━━━${NC}"
    echo ""
    echo -e "  ${CYAN}./scripts/status.sh${NC}     # 查看状态"
    echo -e "  ${CYAN}./scripts/stop.sh${NC}       # 停止服务"
    echo -e "  ${CYAN}tail -f logs/app.log${NC}    # 查看日志"
    echo ""
    
    echo -e "${CYAN}━━━ 下一步操作 ━━━${NC}"
    echo ""
    echo -e "  1. 启动服务: ${CYAN}./scripts/start.sh${NC}"
    echo -e "  2. 测试健康检查: ${CYAN}curl http://localhost:${SERVER_PORT}/health${NC}"
    echo -e "  3. 查看 API 文档: ${CYAN}docs/api/README.md${NC}"
    echo -e "  4. 配置 Nginx (生产): ${CYAN}docs/NGINX_CONFIGURATION.md${NC}"
    echo ""
    
    local report_file=$(generate_install_report)
    log_success "安装报告已保存到: $(basename "$report_file")"
    log_success "安装日志已保存到: $(basename "$LOG_FILE")"
    
    echo ""
}

# ============================================
# 主安装流程 / Main installation process
# ============================================

perform_full_installation() {
    print_header "开始全新安装"
    log "========== 开始全新安装流程 =========="
    log_debug "进入 perform_full_installation 函数"
    
    # 检查依赖 / Check dependencies
    log "步骤 1/6: 检查系统依赖"
    log_debug "调用 check_all_dependencies"
    check_all_dependencies || exit 1
    log_debug "check_all_dependencies 完成"
    
    # 生成配置 / Generate configuration
    log "步骤 2/6: 生成配置文件"
    log_debug "调用 create_config_files"
    create_config_files || exit 1
    log_debug "create_config_files 完成"
    
    # 初始化数据库 / Initialize database
    log "步骤 3/6: 初始化数据库"
    log_debug "调用 init_database"
    init_database || exit 1
    log_debug "init_database 完成"
    
    # 构建应用 / Build application
    log "步骤 4/6: 构建应用"
    log_debug "调用 build_application"
    build_application || exit 1
    log_debug "build_application 完成"
    
    # 创建初始用户 / Create initial user
    log "步骤 5/6: 创建初始用户"
    log_debug "调用 create_initial_user"
    create_initial_user || exit 1
    log_debug "create_initial_user 完成"
    
    # 验证安装 / Verify installation
    log "步骤 6/6: 验证安装"
    log_debug "调用 verify_installation"
    verify_installation || exit 1
    log_debug "verify_installation 完成"
    
    # 显示安装总结 / Show installation summary
    log "========== 全新安装流程完成 =========="
    log_debug "perform_full_installation 函数执行完成"
    print_installation_summary
}

perform_reinstallation() {
    print_header "重新安装"
    
    echo -e "${YELLOW}${ICON_WARNING} 这将覆盖现有配置文件${NC}"
    echo ""
    read -p "确定要继续吗? (y/n) " -n 1 -r
    echo
    
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "操作已取消"
        return 1
    fi
    
    perform_full_installation
}

show_installed_menu() {
    print_header "系统已安装"
    log "检测到已安装的系统，显示菜单"
    
    load_config_vars
    
    echo -e "${GREEN}${ICON_CHECK}${NC} 检测到已安装的系统"
    echo ""
    echo "配置文件:"
    echo "  • ${CONFIG_FILE}"
    echo "  • ${ENV_FILE}"
    echo ""
    echo "数据库: ${DB_NAME} @ ${DB_HOST}:${DB_PORT}"
    echo "服务器端口: ${SERVER_PORT}"
    echo ""
    
    while true; do
        print_menu
        
        log "等待用户选择操作"
        read -p "请选择 [1-4]: " choice
        echo ""
        log "用户选择: $choice"
        
        case $choice in
            1)
                log "用户选择: 重新安装"
                perform_reinstallation
                break
                ;;
            2)
                log "用户选择: 修复数据库"
                repair_database
                echo ""
                read -p "按回车键继续..." -r
                ;;
            3)
                log "用户选择: 直接启动服务"
                start_service
                break
                ;;
            4)
                log "用户选择: 退出"
                log_info "退出安装程序"
                exit 0
                ;;
            *)
                log "无效选择: $choice"
                log_error "无效选择，请输入 1-4"
                echo ""
                ;;
        esac
    done
}

main() {
    # 清空日志文件 / Clear log file
    > "$LOG_FILE"
    
    log_debug "脚本启动参数: $*"
    log_debug "当前工作目录: $(pwd)"
    log_debug "脚本目录: ${SCRIPT_DIR}"
    log_debug "Shell 版本: ${BASH_VERSION}"
    log_debug "用户: $(whoami)"
    log_debug "调试模式: ${DEBUG}"
    
    # 检查终端兼容性 / Check terminal compatibility
    check_terminal
    
    print_header "AI Diet Assistant 安装程序"
    log "安装开始于: $(date)"
    log "安装目录: ${SCRIPT_DIR}"
    
    log_debug "检查配置文件是否存在"
    log_debug "CONFIG_FILE: ${CONFIG_FILE} (存在: $([ -f "$CONFIG_FILE" ] && echo 'yes' || echo 'no'))"
    log_debug "ENV_FILE: ${ENV_FILE} (存在: $([ -f "$ENV_FILE" ] && echo 'yes' || echo 'no'))"
    
    # 检查是否已安装 / Check if already installed
    if is_installed; then
        log_debug "检测到已安装的系统"
        show_installed_menu
    else
        log_debug "未检测到已安装的系统，开始全新安装"
        perform_full_installation
    fi
}

# 运行主函数 / Run main function
main "$@"
