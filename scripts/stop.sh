#!/bin/bash

# AI Diet Assistant 停止脚本

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "==================================="
echo "AI Diet Assistant 停止脚本"
echo "==================================="
echo ""

# 查找进程
PID=$(pgrep -f "diet-assistant" | head -n 1)

if [ -z "$PID" ]; then
    echo -e "${YELLOW}未找到运行中的服务${NC}"
    exit 0
fi

echo "找到进程 PID: $PID"
echo "正在停止服务..."

# 尝试优雅停止
kill -TERM $PID 2>/dev/null

# 等待进程结束
for i in {1..10}; do
    if ! ps -p $PID > /dev/null 2>&1; then
        echo -e "${GREEN}✓ 服务已停止${NC}"
        exit 0
    fi
    sleep 1
done

# 如果还在运行，强制停止
if ps -p $PID > /dev/null 2>&1; then
    echo -e "${YELLOW}优雅停止失败，强制停止...${NC}"
    kill -9 $PID 2>/dev/null
    sleep 1
    
    if ! ps -p $PID > /dev/null 2>&1; then
        echo -e "${GREEN}✓ 服务已强制停止${NC}"
    else
        echo -e "${RED}错误: 无法停止服务${NC}"
        exit 1
    fi
fi
