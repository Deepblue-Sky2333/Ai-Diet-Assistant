#!/bin/bash

# AI Diet Assistant 启动脚本

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}==================================="
echo "AI Diet Assistant 启动脚本"
echo "===================================${NC}"
echo ""

# 检查 .env 文件
if [ -f .env ]; then
    echo -e "${GREEN}✓ 加载 .env 环境变量${NC}"
    source .env
else
    echo -e "${YELLOW}⚠ .env 文件不存在，使用默认配置${NC}"
fi

# 检查配置文件
if [ ! -f "configs/config.yaml" ]; then
    echo -e "${YELLOW}⚠ 配置文件不存在，从示例创建...${NC}"
    if [ -f "configs/config.yaml.example" ]; then
        cp configs/config.yaml.example configs/config.yaml
        echo -e "${GREEN}✓ 已创建 configs/config.yaml${NC}"
        echo -e "${YELLOW}请编辑配置文件后重新运行${NC}"
        echo ""
        echo "  vim configs/config.yaml"
        echo ""
        exit 1
    else
        echo -e "${RED}错误: 配置示例文件不存在${NC}"
        exit 1
    fi
fi

# 检查二进制文件
if [ ! -f "bin/diet-assistant" ]; then
    echo -e "${YELLOW}二进制文件不存在，开始编译...${NC}"
    if command -v make &> /dev/null; then
        make build
    else
        go build -o bin/diet-assistant cmd/server/main.go
    fi
    echo -e "${GREEN}✓ 编译完成${NC}"
    echo ""
fi

# 创建必要的目录
echo "创建必要的目录..."
mkdir -p logs
mkdir -p uploads
mkdir -p web/static
mkdir -p web/templates
echo -e "${GREEN}✓ 目录创建完成${NC}"
echo ""

# 检查端口是否被占用
PORT=${SERVER_PORT:-9090}
if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1 ; then
    echo -e "${RED}错误: 端口 $PORT 已被占用${NC}"
    echo "请停止占用该端口的进程或修改配置文件中的端口号"
    echo ""
    echo "查看占用进程: lsof -i :$PORT"
    exit 1
fi

# 显示启动信息
echo -e "${GREEN}==================================="
echo "启动服务..."
echo "===================================${NC}"
echo ""
echo -e "服务地址: ${BLUE}http://localhost:$PORT${NC}"
echo -e "健康检查: ${BLUE}http://localhost:$PORT/health${NC}"
echo -e "API 文档: ${BLUE}http://localhost:$PORT/api/docs${NC}"
echo ""
echo -e "${YELLOW}按 Ctrl+C 停止服务${NC}"
echo ""

# 启动服务
./bin/diet-assistant
