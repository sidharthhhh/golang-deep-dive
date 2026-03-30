.PHONY: help build run test clean docker-build docker-run docker-stop k8s-deploy k8s-delete migrate-up migrate-down deps lint

# Variables
APP_NAME=auth-service
DOCKER_IMAGE=auth-service:latest
DOCKER_COMPOSE=docker-compose
KUBECTL=kubectl

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build the application
	go build -o bin/$(APP_NAME) cmd/server/main.go

run: ## Run the application locally
	go run cmd/server/main.go

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-coverage: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test-unit: ## Run unit tests only
	go test -v ./internal/...

test-integration: ## Run integration tests
	go test -v ./tests/integration_test.go ./tests/setup_test.go

test-auth: ## Run authentication tests
	go test -v ./tests/auth_test.go ./tests/setup_test.go

test-health: ## Run health check tests
	go test -v ./tests/health_test.go ./tests/setup_test.go

test-user: ## Run user management tests
	go test -v ./tests/user_test.go ./tests/setup_test.go

test-password: ## Run password management tests
	go test -v ./tests/password_test.go ./tests/setup_test.go

test-admin: ## Run admin tests
	go test -v ./tests/admin_test.go ./tests/setup_test.go

test-all: ## Run all tests with verbose output
	bash tests/run_tests.sh

test-quick: ## Run quick tests (auth + health)
	bash tests/run_tests.sh quick

lint: ## Run linter
	golangci-lint run ./...

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

# Docker commands
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run with Docker Compose
	$(DOCKER_COMPOSE) up -d

docker-stop: ## Stop Docker Compose
	$(DOCKER_COMPOSE) down

docker-logs: ## View Docker logs
	$(DOCKER_COMPOSE) logs -f auth-service

docker-clean: ## Clean Docker resources
	$(DOCKER_COMPOSE) down -v
	docker rmi $(DOCKER_IMAGE) || true

# Database migrations
migrate-up: ## Run database migrations up
	@echo "Running migrations..."
	mysql -h localhost -u root -p -D auth_service < migrations/000_init_database.sql
	mysql -h localhost -u root -p -D auth_service < migrations/001_create_users_table.sql
	mysql -h localhost -u root -p -D auth_service < migrations/003_create_token_blacklist.sql
	mysql -h localhost -u root -p -D auth_service < migrations/004_create_password_reset_tokens.sql
	mysql -h localhost -u root -p -D auth_service < migrations/005_add_indexes.sql
	mysql -h localhost -u root -p -D auth_service < migrations/006_add_password_changed_at.sql

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	mysql -h localhost -u root -p -D auth_service -e "DROP TABLE IF EXISTS password_reset_tokens;"
	mysql -h localhost -u root -p -D auth_service -e "DROP TABLE IF EXISTS token_blacklist;"
	mysql -h localhost -u root -p -D auth_service -e "DROP TABLE IF EXISTS users;"

# Kubernetes commands
k8s-deploy: ## Deploy to Kubernetes
	$(KUBECTL) apply -f k8s/configmap.yaml
	$(KUBECTL) apply -f k8s/secrets.yaml
	$(KUBECTL) apply -f k8s/mysql.yaml
	$(KUBECTL) apply -f k8s/deployment.yaml
	$(KUBECTL) apply -f k8s/service.yaml
	$(KUBECTL) apply -f k8s/pdb.yaml
	$(KUBECTL) apply -f k8s/ingress.yaml
	$(KUBECTL) apply -f k8s/hpa.yaml

k8s-delete: ## Delete from Kubernetes
	$(KUBECTL) delete -f k8s/hpa.yaml || true
	$(KUBECTL) delete -f k8s/ingress.yaml || true
	$(KUBECTL) delete -f k8s/pdb.yaml || true
	$(KUBECTL) delete -f k8s/service.yaml || true
	$(KUBECTL) delete -f k8s/deployment.yaml || true
	$(KUBECTL) delete -f k8s/mysql.yaml || true
	$(KUBECTL) delete -f k8s/secrets.yaml || true
	$(KUBECTL) delete -f k8s/configmap.yaml || true

k8s-logs: ## View Kubernetes logs
	$(KUBECTL) logs -f -l app=$(APP_NAME)

k8s-status: ## Check Kubernetes status
	$(KUBECTL) get pods -l app=$(APP_NAME)
	$(KUBECTL) get svc $(APP_NAME)
	$(KUBECTL) get ingress

# Development helpers
dev: ## Run in development mode with hot reload
	air

format: ## Format code
	go fmt ./...
	goimports -w .

vet: ## Run go vet
	go vet ./...

security: ## Run security checks
	gosec ./...

# API testing
api-test: ## Test API endpoints
	@echo "Testing health endpoint..."
	curl -s http://localhost:8080/health | jq
	@echo "\nTesting readiness endpoint..."
	curl -s http://localhost:8080/health/ready | jq

# Complete setup
setup: deps migrate-up ## Complete setup (deps + migrations)
	@echo "Setup complete!"

# Production build
prod-build: ## Build for production
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags='-w -s -extldflags "-static"' \
		-o bin/$(APP_NAME) \
		cmd/server/main.go
