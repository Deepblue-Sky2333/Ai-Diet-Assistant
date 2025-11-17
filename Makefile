.PHONY: help build create-user run run-bin dev test clean migrate-up migrate-down deps fmt lint docker-up docker-down

help: ## 显示帮助信息
	@echo "AI Diet Assistant - 可用命令:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "快速开始:"
	@echo "  make build      # 编译服务"
	@echo "  make run        # 运行服务"
	@echo "  make test       # 运行测试"

build: ## 编译服务
	@echo "编译服务..."
	@mkdir -p bin
	go build -o bin/diet-assistant cmd/server/main.go
	@echo "编译完成: bin/diet-assistant"

create-user: ## 编译用户创建工具
	@echo "编译用户创建工具..."
	@mkdir -p bin
	go build -o bin/create-user cmd/create-user/main.go
	@echo "编译完成: bin/create-user"
	@echo ""
	@echo "使用方法:"
	@echo "  ./bin/create-user -username <用户名> -password <密码> [-email <邮箱>] [-role admin|user]"
	@echo ""
	@echo "示例:"
	@echo "  # 创建第一个用户（自动成为管理员）"
	@echo "  ./bin/create-user -username admin -password adminpass123 -email admin@example.com"
	@echo ""
	@echo "  # 创建普通用户"
	@echo "  ./bin/create-user -username testuser -password userpass123"
	@echo ""
	@echo "  # 显式指定角色"
	@echo "  ./bin/create-user -username admin2 -password admin2pass -role admin"

run: ## 运行服务
	@echo "启动服务..."
	go run cmd/server/main.go

run-bin: ## 运行编译后的二进制文件
	@echo "运行服务..."
	./bin/diet-assistant

dev: ## 开发模式运行（带热重载需要安装 air）
	@echo "开发模式启动..."
	air

test: ## 运行测试
	@echo "运行测试..."
	go test -v ./...

clean: ## 清理编译文件
	@echo "清理..."
	rm -rf bin/

migrate-up: ## 执行数据库迁移（升级）
	@echo "执行数据库迁移..."
	@for file in migrations/*_up.sql; do \
		echo "执行: $$file"; \
		mysql -h localhost -u root -p < $$file; \
	done

migrate-down: ## 回滚数据库迁移
	@echo "回滚数据库迁移..."
	@for file in migrations/*_down.sql; do \
		echo "执行: $$file"; \
		mysql -h localhost -u root -p < $$file; \
	done

deps: ## 下载依赖
	@echo "下载依赖..."
	go mod download
	go mod tidy

fmt: ## 格式化代码
	@echo "格式化代码..."
	go fmt ./...

lint: ## 代码检查
	@echo "代码检查..."
	golangci-lint run

docker-up: ## 启动 Docker 容器
	docker-compose up -d

docker-down: ## 停止 Docker 容器
	docker-compose down
