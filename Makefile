.PHONY: help deps run build test clean docker-build docker-run docker-down lint fmt vet

# Variables
APP_NAME=inventory-api
DOCKER_IMAGE=inventory-system
MAIN_PATH=./cmd/api
BIN_DIR=./bin

# Default target
help: ## Show this help message
	@echo "Go Inventory System - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod verify
	@echo "Dependencies installed successfully"

run: ## Run the application locally
	@echo "Starting application..."
	go run $(MAIN_PATH)/main.go

build: ## Build the application binary
	@echo "Building application..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BIN_DIR)/$(APP_NAME)"

test: ## Run tests with coverage
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo "Coverage report:"
	go tool cover -func=coverage.out

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	go test -v -tags=integration ./tests/integration/...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html
	go clean
	@echo "Clean complete"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -f deployments/docker/Dockerfile -t $(DOCKER_IMAGE):latest .
	@echo "Docker image built: $(DOCKER_IMAGE):latest"

docker-build-optimized: ## Build optimized multi-stage Docker image
	@echo "Building optimized Docker image..."
	docker build -f deployments/docker/Dockerfile.multistage -t $(DOCKER_IMAGE):optimized .
	@echo "Optimized Docker image built: $(DOCKER_IMAGE):optimized"

docker-run: ## Run application with Docker Compose
	@echo "Starting services with Docker Compose..."
	docker compose -f deployments/docker-compose.yml up -d
	@echo "Services started. API available at http://localhost:8080"

docker-down: ## Stop Docker Compose services
	@echo "Stopping services..."
	docker compose -f deployments/docker-compose.yml down
	@echo "Services stopped"

docker-logs: ## View Docker Compose logs
	docker compose -f deployments/docker-compose.yml logs -f

lint: ## Run golangci-lint
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install it from https://golangci-lint.run/usage/install/"; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...
	@echo "Vet complete"

setup: ## Run setup script
	@echo "Running setup script..."
	./scripts/setup.sh

migrate-up: ## Run database migrations (via app startup)
	@echo "Migrations are run automatically on application startup"
	@echo "Start the application with 'make run' or 'make docker-run'"

seed: ## Seed the database with sample data
	@echo "Seeding database..."
	@if [ -z "$$DB_HOST" ]; then export DB_HOST=localhost; fi
	@if [ -z "$$DB_USER" ]; then export DB_USER=postgres; fi
	@if [ -z "$$DB_NAME" ]; then export DB_NAME=inventory_db; fi
	PGPASSWORD=postgres psql -h $$DB_HOST -U $$DB_USER -d $$DB_NAME -f scripts/seed.sql
	@echo "Database seeded"

all: deps fmt vet test build ## Run all checks and build
