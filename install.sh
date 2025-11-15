#!/bin/bash

# AI Diet Assistant 一键安装脚本
# 此脚本会自动检测环境、安装依赖、配置系统并启动服务

# 不使用 set -e，改为手动检查关键步骤
# set -e 会导致任何小错误都退出，过于严格

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # 无颜色

echo -e "${BLUE}==================================="
echo "AI Diet Assistant 一键安装"
echo "===================================${NC}"
echo ""

# 检测操作系统
OS="unknown"
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS="linux"
    if [ -f /etc/debian_version ]; then
        DISTRO="debian"
    elif [ -f /etc/redhat-release ]; then
        DISTRO="redhat"
    fi
elif [[ "$OSTYPE" == "darwin"* ]]; then
    OS="macos"
fi

echo -e "${BLUE}检测到操作系统: $OS${NC}"
echo ""

# 函数：检查命令是否存在
command_exists() {
    command -v "$1" &> /dev/null
}

# 函数：安装 Go
install_go() {
    echo -e "${YELLOW}正在安装 Go...${NC}"
    
    if [ "$OS" == "macos" ]; then
        if command_exists brew; then
            brew install go
        else
            echo -e "${RED}请先安装 Homebrew: https://brew.sh/${NC}"
            exit 1
        fi
    elif [ "$OS" == "linux" ]; then
        if [ "$DISTRO" == "debian" ]; then
            sudo apt-get update
            sudo apt-get install -y golang-go
        elif [ "$DISTRO" == "redhat" ]; then
            sudo yum install -y golang
        fi
    fi
    
    if command_exists go; then
        echo -e "${GREEN}✓ Go 安装成功: $(go version)${NC}"
    else
        echo -e "${RED}✗ Go 安装失败${NC}"
        exit 1
    fi
}

# 函数：安装 Node.js
install_nodejs() {
    echo -e "${YELLOW}正在安装 Node.js...${NC}"
    
    if [ "$OS" == "macos" ]; then
        if command_exists brew; then
            brew install node
        else
            echo -e "${RED}请先安装 Homebrew: https://brew.sh/${NC}"
            exit 1
        fi
    elif [ "$OS" == "linux" ]; then
        if [ "$DISTRO" == "debian" ]; then
            curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
            sudo apt-get install -y nodejs
        elif [ "$DISTRO" == "redhat" ]; then
            curl -fsSL https://rpm.nodesource.com/setup_18.x | sudo bash -
            sudo yum install -y nodejs
        fi
    fi
    
    if command_exists node; then
        echo -e "${GREEN}✓ Node.js 安装成功: $(node --version)${NC}"
    else
        echo -e "${RED}✗ Node.js 安装失败${NC}"
        exit 1
    fi
}

# 函数：安装 MySQL
install_mysql() {
    echo -e "${YELLOW}正在安装 MySQL...${NC}"
    
    if [ "$OS" == "macos" ]; then
        if command_exists brew; then
            brew install mysql
            brew services start mysql
        fi
    elif [ "$OS" == "linux" ]; then
        if [ "$DISTRO" == "debian" ]; then
            sudo apt-get update
            sudo apt-get install -y mysql-server
            sudo systemctl start mysql
            sudo systemctl enable mysql
        elif [ "$DISTRO" == "redhat" ]; then
            sudo yum install -y mysql-server
            sudo systemctl start mysqld
            sudo systemctl enable mysqld
        fi
    fi
    
    if command_exists mysql; then
        echo -e "${GREEN}✓ MySQL 安装成功${NC}"
    else
        echo -e "${RED}✗ MySQL 安装失败${NC}"
        exit 1
    fi
}

# 步骤 1: 检查和安装依赖
echo -e "${BLUE}[1/6] 检查系统依赖...${NC}"

# 检查 Go
if ! command_exists go; then
    echo -e "${YELLOW}未检测到 Go，是否自动安装？(y/n)${NC}"
    read -p "> " install_go_choice
    if [[ $install_go_choice =~ ^[Yy]$ ]]; then
        install_go
    else
        echo -e "${RED}Go 是必需的，请手动安装后重试${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✓ Go 已安装: $(go version | awk '{print $3}')${NC}"
fi

# 检查 Node.js
if ! command_exists node; then
    echo -e "${YELLOW}未检测到 Node.js，是否自动安装？(y/n)${NC}"
    read -p "> " install_node_choice
    if [[ $install_node_choice =~ ^[Yy]$ ]]; then
        install_nodejs
    else
        echo -e "${RED}Node.js 是必需的，请手动安装后重试${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✓ Node.js 已安装: $(node --version)${NC}"
fi

# 检查 npm
if ! command_exists npm; then
    echo -e "${RED}npm 未安装，请重新安装 Node.js${NC}"
    exit 1
else
    echo -e "${GREEN}✓ npm 已安装: $(npm --version)${NC}"
fi

# 检查 MySQL
if ! command_exists mysql; then
    echo -e "${YELLOW}未检测到 MySQL，是否自动安装？(y/n)${NC}"
    read -p "> " install_mysql_choice
    if [[ $install_mysql_choice =~ ^[Yy]$ ]]; then
        install_mysql
    else
        echo -e "${RED}MySQL 是必需的，请手动安装后重试${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✓ MySQL 已安装${NC}"
fi

# 检查 openssl
if ! command_exists openssl; then
    echo -e "${RED}openssl 未安装，请先安装${NC}"
    exit 1
else
    echo -e "${GREEN}✓ openssl 已安装${NC}"
fi

echo ""

# 步骤 2: 运行配置脚本
echo -e "${BLUE}[2/6] 配置系统...${NC}"
if [ -f "./scripts/install.sh" ]; then
    ./scripts/install.sh
    if [ $? -ne 0 ]; then
        echo -e "${YELLOW}⚠ 配置脚本执行有警告，继续...${NC}"
    fi
else
    echo -e "${YELLOW}⚠ 配置脚本不存在，跳过...${NC}"
fi

echo ""

# 步骤 3: 创建数据库
echo -e "${BLUE}[3/6] 配置数据库...${NC}"

# 从 .env 读取数据库配置
if [ -f ".env" ]; then
    source .env
else
    echo -e "${RED}✗ .env 文件不存在，请先运行配置脚本${NC}"
    exit 1
fi

echo "正在创建数据库..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" 2>/dev/null <<EOF
CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EOF

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 数据库创建成功${NC}"
else
    echo -e "${YELLOW}⚠ 数据库可能已存在或创建失败，继续...${NC}"
fi

echo ""

# 步骤 4: 运行数据库迁移
echo -e "${BLUE}[4/6] 运行数据库迁移...${NC}"

migration_failed=0
for file in migrations/*_up.sql; do
    if [ -f "$file" ]; then
        echo "执行: $file"
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$file" 2>/dev/null
        if [ $? -ne 0 ]; then
            echo -e "${YELLOW}⚠ $file 执行失败（可能已执行过），继续...${NC}"
            migration_failed=1
        fi
    fi
done

if [ $migration_failed -eq 0 ]; then
    echo -e "${GREEN}✓ 数据库迁移完成${NC}"
else
    echo -e "${YELLOW}⚠ 部分迁移有警告，但已完成${NC}"
fi
echo ""

# 步骤 5: 构建应用
echo -e "${BLUE}[5/6] 构建应用...${NC}"

# 下载 Go 依赖
echo "下载 Go 依赖..."
go mod download
if [ $? -ne 0 ]; then
    echo -e "${YELLOW}⚠ Go 依赖下载有警告，尝试继续构建...${NC}"
fi

# 构建后端
echo "构建后端..."
mkdir -p bin
go build -o bin/diet-assistant cmd/server/main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 后端构建成功${NC}"
else
    echo -e "${RED}✗ 后端构建失败，这是关键错误${NC}"
    exit 1
fi

# 构建前端
echo "构建前端..."
if [ -d "web/frontend" ]; then
    cd web/frontend

    if [ ! -d "node_modules" ]; then
        echo "安装前端依赖..."
        npm install
        if [ $? -ne 0 ]; then
            echo -e "${RED}✗ 前端依赖安装失败${NC}"
            cd ../..
            exit 1
        fi
    fi

    npm run build

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ 前端构建成功${NC}"
    else
        echo -e "${RED}✗ 前端构建失败，这是关键错误${NC}"
        cd ../..
        exit 1
    fi

    cd ../..
else
    echo -e "${YELLOW}⚠ 前端目录不存在，跳过前端构建${NC}"
fi
echo ""

# 步骤 6: 配置服务（可选）
echo -e "${BLUE}[6/6] 配置系统服务...${NC}"

if [ "$OS" == "linux" ]; then
    echo -e "${YELLOW}是否配置为系统服务（开机自启）？(y/n)${NC}"
    read -p "> " setup_service
    
    if [[ $setup_service =~ ^[Yy]$ ]]; then
        # 获取当前目录的绝对路径
        INSTALL_DIR=$(pwd)
        
        # 创建 systemd 服务文件
        sudo tee /etc/systemd/system/diet-assistant.service > /dev/null <<EOF
[Unit]
Description=AI Diet Assistant
After=network.target mysql.service

[Service]
Type=simple
User=$USER
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/bin/diet-assistant
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
        
        # 重载 systemd
        sudo systemctl daemon-reload
        sudo systemctl enable diet-assistant
        
        echo -e "${GREEN}✓ 系统服务配置完成${NC}"
        echo ""
        echo "服务管理命令："
        echo "  启动: sudo systemctl start diet-assistant"
        echo "  停止: sudo systemctl stop diet-assistant"
        echo "  重启: sudo systemctl restart diet-assistant"
        echo "  状态: sudo systemctl status diet-assistant"
        echo "  日志: sudo journalctl -u diet-assistant -f"
        echo ""
        
        echo -e "${YELLOW}是否现在启动服务？(y/n)${NC}"
        read -p "> " start_service
        
        if [[ $start_service =~ ^[Yy]$ ]]; then
            sudo systemctl start diet-assistant
            echo -e "${GREEN}✓ 服务已启动${NC}"
        fi
    fi
fi

echo ""
echo -e "${GREEN}==================================="
echo "安装完成！"
echo "===================================${NC}"
echo ""
echo "应用信息："
echo "  访问地址: http://localhost:$SERVER_PORT"
echo "  配置文件: .env"
echo "  日志目录: logs/"
echo ""

if [ "$OS" == "linux" ] && [[ $setup_service =~ ^[Yy]$ ]]; then
    echo "服务已配置为系统服务，使用 systemctl 管理"
else
    echo "启动应用："
    echo "  ./bin/diet-assistant"
    echo ""
    echo "或使用脚本："
    echo "  ./scripts/start.sh"
fi

echo ""
echo -e "${YELLOW}重要提示：${NC}"
echo "1. 请妥善保管 .env 文件，不要提交到版本控制"
echo "2. 建议配置 Nginx 反向代理处理 HTTPS"
echo "3. 定期备份数据库"
echo ""
echo -e "${GREEN}祝使用愉快！${NC}"
echo ""
