#!/bin/bash

# AI Diet Assistant - Build All (Frontend + Backend)
# 此脚本构建前后端并集成到一起

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # 无颜色

echo -e "${BLUE}==================================="
echo "AI Diet Assistant - 构建所有组件"
echo "===================================${NC}"
echo ""

# 获取脚本目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

# 检查 Node.js
if ! command -v node &> /dev/null; then
    echo -e "${RED}错误：未安装 Node.js！${NC}"
    echo "请先安装 Node.js：https://nodejs.org/"
    exit 1
fi

# 检查 Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误：未安装 Go！${NC}"
    echo "请先安装 Go：https://golang.org/"
    exit 1
fi

echo -e "${GREEN}✓ Node.js 版本: $(node --version)${NC}"
echo -e "${GREEN}✓ Go 版本: $(go version | awk '{print $3}')${NC}"
echo ""

# 构建前端
echo -e "${BLUE}[1/3] 构建前端...${NC}"
cd "$PROJECT_ROOT/web/frontend"

# 检查依赖
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}安装前端依赖...${NC}"
    npm install
fi

# 构建前端
echo "正在构建 Next.js 应用..."
npm run build

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 前端构建成功${NC}"
else
    echo -e "${RED}✗ 前端构建失败${NC}"
    exit 1
fi
echo ""

# 构建后端
echo -e "${BLUE}[2/3] 构建后端...${NC}"
cd "$PROJECT_ROOT"

echo "正在编译 Go 应用..."
make build

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 后端构建成功${NC}"
else
    echo -e "${RED}✗ 后端构建失败${NC}"
    exit 1
fi
echo ""

# 验证构建
echo -e "${BLUE}[3/3] 验证构建...${NC}"

# 检查后端二进制文件
if [ -f "$PROJECT_ROOT/bin/diet-assistant" ]; then
    echo -e "${GREEN}✓ 后端二进制文件: bin/diet-assistant${NC}"
else
    echo -e "${RED}✗ 后端二进制文件未找到${NC}"
    exit 1
fi

# 检查前端构建文件
if [ -d "$PROJECT_ROOT/web/frontend/.next" ]; then
    echo -e "${GREEN}✓ 前端构建文件: web/frontend/.next${NC}"
else
    echo -e "${RED}✗ 前端构建文件未找到${NC}"
    exit 1
fi

if [ -d "$PROJECT_ROOT/web/frontend/.next/standalone" ]; then
    echo -e "${GREEN}✓ 前端独立模式: web/frontend/.next/standalone${NC}"
else
    echo -e "${YELLOW}⚠ 前端独立模式未启用${NC}"
fi

echo ""
echo -e "${GREEN}==================================="
echo "构建完成！"
echo "===================================${NC}"
echo ""
echo "下一步："
echo "1. 确保 .env 中 SERVER_MODE=release"
echo "2. 启动应用：./bin/diet-assistant"
echo "3. 访问：http://localhost:9090"
echo ""
echo -e "${YELLOW}注意：前后端已集成，只需启动后端即可访问完整应用${NC}"
echo ""
