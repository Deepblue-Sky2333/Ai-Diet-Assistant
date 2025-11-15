#!/bin/bash

# AI Diet Assistant 状态检查脚本

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "==================================="
echo "AI Diet Assistant 状态检查"
echo "==================================="
echo ""

# 检查进程
PID=$(pgrep -f "diet-assistant" | head -n 1)

if [ -z "$PID" ]; then
    echo -e "服务状态: ${RED}未运行${NC}"
    exit 1
else
    echo -e "服务状态: ${GREEN}运行中${NC}"
    echo "进程 PID: $PID"
    
    # 显示进程信息
    ps -p $PID -o pid,ppid,%cpu,%mem,etime,cmd
    echo ""
fi

# 检查端口
PORT=${SERVER_PORT:-9090}
if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1 ; then
    echo -e "端口 $PORT: ${GREEN}监听中${NC}"
else
    echo -e "端口 $PORT: ${RED}未监听${NC}"
fi
echo ""

# 检查健康状态
if command -v curl &> /dev/null; then
    echo "健康检查..."
    HEALTH=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:$PORT/health 2>/dev/null)
    
    if [ "$HEALTH" = "200" ]; then
        echo -e "健康状态: ${GREEN}正常 (HTTP $HEALTH)${NC}"
    else
        echo -e "健康状态: ${YELLOW}异常 (HTTP $HEALTH)${NC}"
    fi
fi
