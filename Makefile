.PHONY: help build run test clean docker-build docker-compose-up docker-compose-down

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building application..."
	go build -o main ./cmd/api

run: ## Run the application locally
	@echo "Running application..."
	go run ./cmd/api

dev: ## Run the application with hot reload (requires air)
	@command -v air >/dev/null 2>&1 || go install github.com/cosmtrek/air@latest
	air

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f main
	rm -f coverage.out coverage.html
	go clean

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -f build/package/Dockerfile -t golang-docker:latest .

docker-compose-up: ## Start Docker Compose services
	@echo "Starting Docker Compose services..."
	docker-compose -f deployments/docker-compose.yml up -d

docker-compose-down: ## Stop Docker Compose services
	@echo "Stopping Docker Compose services..."
	docker-compose -f deployments/docker-compose.yml down

docker-compose-logs: ## View Docker Compose logs
	docker-compose -f deployments/docker-compose.yml logs -f

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

lint: ## Run golangci-lint (requires golangci-lint)
	@command -v golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

all: clean build test ## Run clean, build and test