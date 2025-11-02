.PHONY: build clean test run help

# Binary name
BINARY_NAME=blockchain
BUILD_DIR=./build
CMD_DIR=./cmd/blockchain

# Build the application
build:
	@echo "Building..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts and data
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -rf ./tmp
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

# Run the application (requires build first)
run: build
	@$(BUILD_DIR)/$(BINARY_NAME)

# Development build with race detector
dev:
	@echo "Building with race detector..."
	@go build -race -o $(BUILD_DIR)/$(BINARY_NAME)-dev $(CMD_DIR)/main.go

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@golint ./...

# Run vet
vet:
	@echo "Running vet..."
	@go vet ./...

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker-compose build

docker-up:
	@echo "Starting blockchain network..."
	@docker-compose up -d
	@echo "Network started! View logs with: docker-compose logs -f"

docker-down:
	@echo "Stopping blockchain network..."
	@docker-compose down

docker-clean:
	@echo "Cleaning Docker containers and volumes..."
	@docker-compose down -v
	@docker system prune -f

docker-logs:
	@docker-compose logs -f

docker-test:
	@./scripts/docker-test.sh

# Network demo
network-demo:
	@./scripts/network-demo.sh

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  clean          - Remove build artifacts and data"
	@echo "  test           - Run tests"
	@echo "  deps           - Install dependencies"
	@echo "  run            - Build and run the application"
	@echo "  dev            - Build with race detector"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  vet            - Run go vet"
	@echo ""
	@echo "Docker commands:"
	@echo "  docker-build   - Build Docker images"
	@echo "  docker-up      - Start blockchain network"
	@echo "  docker-down    - Stop blockchain network"
	@echo "  docker-clean   - Clean containers and volumes"
	@echo "  docker-logs    - View network logs"
	@echo "  docker-test    - Run full Docker test"
	@echo ""
	@echo "Network:"
	@echo "  network-demo   - Setup local network demo"
	@echo ""
	@echo "  help           - Show this help message"
