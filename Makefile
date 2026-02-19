.PHONY: help build run test clean docker-up docker-down lint format migrate

APP_NAME=cadastral-service
API_BINARY=main
MOCK_BINARY=mock-server
DOCKER_COMPOSE=docker-compose

RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

help: 
	@echo "$(YELLOW)Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

build: ## build the application
	@echo "$(BLUE)Building application...$(NC)"
	go build -o $(API_BINARY) ./cmd/api
	go build -o $(MOCK_BINARY) ./cmd/mock-server
	@echo "$(GREEN)Build successful!$(NC)"

run: ## run the API server
	@echo "$(BLUE)Starting API server...$(NC)"
	DATABASE_URL=postgres://cadastral:cadastral123@localhost:5432/cadastral_db?sslmode=disable \
	PORT=8080 \
	go run ./cmd/api/main.go

run-mock: 
	@echo "$(BLUE)Starting mock server...$(NC)"
	PORT=8081 \
	DELAY_MAX=60 \
	go run ./cmd/mock-server/main.go

dev: ## run both API and mock servers in development mode (requires tmux)
	@echo "$(BLUE)Starting development environment...$(NC)"
	tmux new-session -d -s cadastral-service \
		"make run" \; \
	split-window -h \
		"make run-mock" \; \
	attach-session -t cadastral-service

test:
	@echo "$(BLUE)Running tests...$(NC)"
	go test -v ./test/...
	@echo "$(GREEN)Tests completed!$(NC)"

test-coverage:
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

lint: ## run linter
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)golangci-lint not found, installing...$(NC)"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2; \
		golangci-lint run ./...; \
	fi

format: ## format Go code
	@echo "$(BLUE)Formatting code...$(NC)"
	go fmt ./...
	@echo "$(GREEN)Formatting completed!$(NC)"

clean: ## clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	rm -f $(API_BINARY) $(MOCK_BINARY)
	rm -f coverage.out coverage.html
	find . -name "*.test" -type f -delete
	@echo "$(GREEN)Clean completed!$(NC)"

# database commands
migrate: 
	@echo "$(BLUE)Running migrations...$(NC)"
	@docker exec -it $$(docker ps -qf "name=postgres") psql -U cadastral -d cadastral_db -c "\
	CREATE TABLE IF NOT EXISTS users (\
		id VARCHAR(255) PRIMARY KEY,\
		username VARCHAR(255) UNIQUE NOT NULL,\
		password_hash TEXT NOT NULL,\
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP\
	);\
	CREATE TABLE IF NOT EXISTS queries (\
		id VARCHAR(255) PRIMARY KEY,\
		cadastral_number VARCHAR(255) NOT NULL,\
		latitude DOUBLE PRECISION NOT NULL,\
		longitude DOUBLE PRECISION NOT NULL,\
		status VARCHAR(50) NOT NULL DEFAULT 'pending',\
		result BOOLEAN,\
		user_id VARCHAR(255),\
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,\
		completed_at TIMESTAMP WITH TIME ZONE\
	);\
	CREATE INDEX IF NOT EXISTS idx_queries_cadastral ON queries(cadastral_number);\
	CREATE INDEX IF NOT EXISTS idx_queries_user_id ON queries(user_id);\
	CREATE INDEX IF NOT EXISTS idx_queries_created_at ON queries(created_at DESC);\
	CREATE INDEX IF NOT EXISTS idx_queries_status ON queries(status);"
	@echo "$(GREEN)Migrations completed!$(NC)"

migrate-init: 
	@echo "$(BLUE)Initializing database with sample data...$(NC)"
	@docker exec -it $$(docker ps -qf "name=postgres") psql -U cadastral -d cadastral_db -c "\
	INSERT INTO users (id, username, password_hash, created_at) VALUES \
	('admin_001', 'admin', '\$2a\$10\$N9qo8uLOickgx2ZMRZoMyeS7.2Y5Z1e8Z5c6W5q5k5n5v5c5n5v5c5n', CURRENT_TIMESTAMP) \
	ON CONFLICT (username) DO NOTHING;"
	@echo "$(GREEN)Database initialized!$(NC)"

docker-up: ## start Docker containers
	@echo "$(BLUE)Starting Docker containers...$(NC)"
	$(DOCKER_COMPOSE) up -d --build
	@echo "$(GREEN)Docker containers started!$(NC)"
	@echo "$(YELLOW)API Server: http://localhost:8080$(NC)"
	@echo "$(YELLOW)Mock Server: http://localhost:8081$(NC)"
	@echo "$(YELLOW)PostgreSQL: localhost:5432$(NC)"

docker-down: ## stop Docker containers
	@echo "$(BLUE)Stopping Docker containers...$(NC)"
	$(DOCKER_COMPOSE) down
	@echo "$(GREEN)Docker containers stopped!$(NC)"

docker-logs: ## show Docker logs
	@echo "$(BLUE)Showing Docker logs...$(NC)"
	$(DOCKER_COMPOSE) logs -f

docker-logs-api:
	@echo "$(BLUE)Showing API server logs...$(NC)"
	$(DOCKER_COMPOSE) logs -f api

docker-logs-mock:
	@echo "$(BLUE)Showing mock server logs...$(NC)"
	$(DOCKER_COMPOSE) logs -f mock-server

docker-restart: ## restart Docker containers
	@echo "$(BLUE)Restarting Docker containers...$(NC)"
	$(DOCKER_COMPOSE) restart
	@echo "$(GREEN)Docker containers restarted!$(NC)"

docker-clean: ## clean Docker containers and volumes
	@echo "$(BLUE)Cleaning Docker containers and volumes...$(NC)"
	$(DOCKER_COMPOSE) down -v --rmi all
	@echo "$(GREEN)Docker cleaned!$(NC)"

# API testing commands
test-api: 
	@echo "$(BLUE)Testing API endpoints...$(NC)"
	@echo "$(YELLOW)Testing /ping...$(NC)"
	@curl -s http://localhost:8080/ping | jq . || echo "Server not running"
	@echo ""
	@echo "$(YELLOW)Testing POST /query...$(NC)"
	@curl -s -X POST http://localhost:8080/query \
		-H "Content-Type: application/json" \
		-d '{"cadastral_number": "77:01:0001010:1234", "latitude": 55.7558, "longitude": 37.6176}' | jq . || echo "Request failed"
	@echo ""
	@echo "$(YELLOW)Testing GET /history...$(NC)"
	@curl -s http://localhost:8080/history | jq . || echo "Request failed"

test-mock: 
	@echo "$(BLUE)Testing mock server...$(NC)"
	@echo "$(YELLOW)Testing POST /api/result...$(NC)"
	@curl -s -X POST http://localhost:8081/api/result \
		-H "Content-Type: application/json" \
		-d '{"cadastral_number": "77:01:0001010:1234", "latitude": 55.7558, "longitude": 37.6176}' | jq . || echo "Request failed"

# health checks
health: ## check service health
	@echo "$(BLUE)Checking service health...$(NC)"
	@echo "$(YELLOW)API Server:$(NC)"
	@curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" http://localhost:8080/ping || echo "API Server not running"
	@echo ""
	@echo "$(YELLOW)Mock Server:$(NC)"
	@curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" http://localhost:8081/ping || echo "Mock Server not running"
	@echo ""
	@echo "$(YELLOW)PostgreSQL:$(NC)"
	@docker exec -it $$(docker ps -qf "name=postgres") pg_isready -U cadastral && echo "PostgreSQL is ready" || echo "PostgreSQL not ready"

setup: ## setup development environment
	@echo "$(BLUE)Setting up development environment...$(NC)"
	@if [ ! -f go.mod ]; then \
		echo "$(YELLOW)Initializing Go module...$(NC)"; \
		go mod init cadastral-service; \
	fi
	@echo "$(YELLOW)Installing dependencies...$(NC)"
	go get github.com/gin-gonic/gin
	go get github.com/lib/pq
	go get github.com/stretchr/testify
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	go mod tidy
	@echo "$(YELLOW)Creating necessary directories...$(NC)"
	mkdir -p cmd/{api,mock-server} internal/{api,config,models,repository,service} pkg/{database,logger} test
	@echo "$(GREEN)Setup completed!$(NC)"

# documentation
docs: 
	@echo "$(BLUE)Generating documentation...$(NC)"
	@if command -v swag >/dev/null 2>&1; then \
		echo "$(YELLOW)Generating Swagger docs...$(NC)"; \
		swag init -g cmd/api/main.go -o docs; \
	else \
		echo "$(YELLOW)Installing swag...$(NC)"; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		swag init -g cmd/api/main.go -o docs; \
	fi
	@echo "$(GREEN)Documentation generated in docs/ folder$(NC)"

# benchmarks
bench: 
	@echo "$(BLUE)Running benchmarks...$(NC)"
	go test -bench=. -benchmem ./...

# dependency management
deps-update: 
	@echo "$(BLUE)Updating dependencies...$(NC)"
	go get -u ./...
	go mod tidy
	@echo "$(GREEN)Dependencies updated!$(NC)"

deps-check: 
	@echo "$(BLUE)Checking for outdated dependencies...$(NC)"
	@if command -v go-mod-outdated >/dev/null 2>&1; then \
		go-mod-outdated; \
	else \
		echo "$(YELLOW)go-mod-outdated not found, installing...$(NC)"; \
		go install github.com/psampaz/go-mod-outdated@latest; \
		go-mod-outdated; \
	fi

# build for production
release:
	@echo "$(BLUE)Building production release...$(NC)"
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/$(APP_NAME)-linux-amd64 ./cmd/api
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/$(APP_NAME)-darwin-amd64 ./cmd/api
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/$(APP_NAME)-windows-amd64.exe ./cmd/api
	@echo "$(GREEN)Production builds created in dist/ folder$(NC)"

# code quality
audit: ## run code audit (lint + test + format)
	@echo "$(BLUE)Running code audit...$(NC)"
	@make lint
	@make test
	@make format
	@echo "$(GREEN)Code audit completed!$(NC)"

quickstart: 
	@make setup
	@make docker-up
	@sleep 10
	@make migrate
	@make migrate-init
	@echo "$(GREEN)Quick start completed!$(NC)"
	@echo "$(YELLOW)API Server: http://localhost:8080$(NC)"
	@echo "$(YELLOW)Mock Server: http://localhost:8081$(NC)"
	@echo "$(YELLOW)Try: make test-api$(NC)"