#!/bin/bash

# AI Diet Assistant 项目初始化脚本

set -e

echo "=== AI Diet Assistant 项目初始化 ==="

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "错误: Go 未安装，请先安装 Go"
    exit 1
fi

echo "✓ Go 版本: $(go version)"

# 下载依赖
echo "下载 Go 依赖..."
go mod download
go mod tidy

# 创建必要的目录
echo "创建必要的目录..."
mkdir -p logs
mkdir -p uploads

# 复制配置文件模板
if [ ! -f configs/config.yaml ]; then
    echo "创建配置文件..."
    cp configs/config.yaml.example configs/config.yaml 2>/dev/null || echo "请手动创建 configs/config.yaml"
fi

echo ""
echo "=== 初始化完成 ==="
echo ""
echo "下一步:"
echo "1. 配置 configs/config.yaml 文件"
echo "2. 创建数据库: CREATE DATABASE ai_diet_assistant;"
echo "3. 运行数据库迁移: make migrate-up"
echo "4. 启动服务: make run"
echo ""
