#!/bin/bash

# 数据库设置脚本
# 用于创建数据库和运行迁移

echo "=== AI Diet Assistant 数据库设置 ==="
echo ""

# 数据库配置（从 config.yaml 读取或使用默认值）
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-3306}
DB_USER=${DB_USER:-root}
DB_NAME=${DB_NAME:-ai_diet_assistant}

# 提示输入密码
echo "请输入 MySQL $DB_USER 用户的密码（如果没有密码，直接按回车）:"
read -s DB_PASSWORD

# 构建 MySQL 命令
if [ -z "$DB_PASSWORD" ]; then
    MYSQL_CMD="mysql -h $DB_HOST -P $DB_PORT -u $DB_USER"
else
    MYSQL_CMD="mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD"
fi

echo ""
echo "1. 检查数据库连接..."
if ! $MYSQL_CMD -e "SELECT 1;" > /dev/null 2>&1; then
    echo "❌ 无法连接到 MySQL 数据库"
    echo "请检查："
    echo "  - MySQL 服务是否运行"
    echo "  - 用户名和密码是否正确"
    exit 1
fi
echo "✅ 数据库连接成功"

echo ""
echo "2. 创建数据库（如果不存在）..."
$MYSQL_CMD -e "CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
echo "✅ 数据库 $DB_NAME 已准备就绪"

echo ""
echo "3. 运行数据库迁移..."
for migration in migrations/*_up.sql; do
    if [ -f "$migration" ]; then
        echo "   执行: $(basename $migration)"
        $MYSQL_CMD $DB_NAME < "$migration"
    fi
done
echo "✅ 数据库迁移完成"

echo ""
echo "=== 设置完成 ==="
echo ""
echo "现在你可以运行应用："
echo "  ./bin/diet-assistant"
echo ""
echo "默认测试账号："
echo "  用户名: test"
echo "  密码: 1145141919810"
echo ""
