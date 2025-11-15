.PHONY: help build run test clean migrate-up migrate-down dev

help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## 编译后端
	@echo "编译后端..."
	@mkdir -p bin
	go build -o bin/diet-assistant cmd/server/main.go
	@echo "编译完成: bin/diet-assistant"

build-frontend: ## 构建前端
	@echo "构建前端..."
	cd web/frontend && npm run build
	@echo "前端构建完成"

build-all: ## 构建前后端（集成模式）
	@echo "构建所有组件..."
	./scripts/build-all.sh

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
