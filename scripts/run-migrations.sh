#!/bin/bash

# 简单的数据库迁移脚本

echo "=== 运行数据库迁移 ==="
echo ""

# 从配置文件读取数据库信息
DB_NAME="ai_diet_assistant"

# 检查 MySQL 连接（无密码）
echo "1. 检查数据库连接..."
if ! mysql -u root -e "SELECT 1;" > /dev/null 2>&1; then
    echo "❌ 无法连接到 MySQL（尝试无密码连接）"
    echo ""
    echo "如果你的 MySQL 有密码，请手动运行："
    echo "  mysql -u root -p $DB_NAME < migrations/001_init_up.sql"
    echo "  mysql -u root -p $DB_NAME < migrations/002_add_password_version_up.sql"
    exit 1
fi
echo "✅ 数据库连接成功"

echo ""
echo "2. 创建数据库（如果不存在）..."
mysql -u root -e "CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
echo "✅ 数据库已准备就绪"

echo ""
echo "3. 运行迁移..."
for migration in migrations/*_up.sql; do
    if [ -f "$migration" ]; then
        echo "   执行: $(basename $migration)"
        mysql -u root $DB_NAME < "$migration" 2>&1 | grep -v "Warning: Using a password on the command line"
    fi
done
echo "✅ 迁移完成"

echo ""
echo "=== 完成 ==="
echo "现在可以重启应用了"
