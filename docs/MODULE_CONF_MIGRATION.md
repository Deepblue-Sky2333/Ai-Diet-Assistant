# Module Configuration File Migration

## 变更说明

`module.conf` 文件已从 `.kiro/module.conf` 移动到 `configs/module.conf`。

### 原因
- 将配置文件统一放在 `configs/` 目录中，更符合项目结构规范
- `.kiro/` 目录主要用于 Kiro IDE 的规范和任务管理
- `configs/` 目录是项目配置文件的标准位置

### 影响的文件
以下脚本已更新以使用新路径：
- `scripts/install.sh`
- `scripts/update-module-path.sh`
- `scripts/update-module-path-auto.sh`

### 迁移步骤
如果您有旧的 `.kiro/module.conf` 文件：
```bash
# 文件已自动移动，无需手动操作
# 如果需要手动迁移：
mv .kiro/module.conf configs/module.conf
```

### 使用方法
配置文件的使用方法保持不变：
```bash
# 更新模块路径
./scripts/update-module-path.sh
```

### 备份文件
备份文件现在会创建在 `configs/` 目录中：
- 格式: `configs/module.conf.backup.YYYYMMDD_HHMMSS`
- 这些备份文件已添加到 `.gitignore` 中

## 文件内容
`configs/module.conf` 包含项目的 Go 模块路径配置：
```properties
MODULE_PATH=github.com/Deepblue-Sky2333/Ai-Diet-Assistant
```

更新此值后，运行 `./scripts/update-module-path.sh` 以应用更改。
