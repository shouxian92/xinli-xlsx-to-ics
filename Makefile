# Makefile for XLSX to ICS Telegram Bot
# Linux/macOS development helper

.PHONY: help build run test clean deps docker-build docker-run deploy

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the bot binary"
	@echo "  run          - Run the bot locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  deploy       - Build for deployment"

# Build the bot
build:
	@echo "Building bot..."
	go build -o xlsx-to-ics .

# Run the bot
run:
	@echo "Running bot..."
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Please copy config.env.example to .env and add your bot token."; \
		exit 1; \
	fi
	@source .env && go run .

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	go clean
	rm -f xlsx-to-ics xlsx-to-ics-* *.exe

# Download and tidy dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	go mod verify

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t xlsx-to-ics-bot .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Please copy config.env.example to .env and add your bot token."; \
		exit 1; \
	fi
	docker run -p 8080:8080 --env-file .env xlsx-to-ics-bot

# Build for deployment
deploy:
	@echo "Building for deployment..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o xlsx-to-ics-linux .

# Install dependencies (Ubuntu/Debian)
install-deps-ubuntu:
	@echo "Installing dependencies for Ubuntu/Debian..."
	sudo apt update
	sudo apt install -y golang-go git make

# Install dependencies (CentOS/RHEL/Fedora)
install-deps-centos:
	@echo "Installing dependencies for CentOS/RHEL/Fedora..."
	sudo dnf install -y golang git make

# Install dependencies (macOS)
install-deps-macos:
	@echo "Installing dependencies for macOS..."
	brew install go git make

# Setup development environment
setup: install-deps-ubuntu deps
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then \
		cp config.env.example .env; \
		echo "Created .env file. Please edit it with your bot token."; \
	fi
	chmod +x run_bot.sh
	@echo "Setup complete! Edit .env with your bot token and run 'make run'"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Check code quality
check: fmt lint test
	@echo "Code quality checks completed!"

# Show bot status
status:
	@echo "Bot Status:"
	@if pgrep -f "xlsx-to-ics" > /dev/null; then \
		echo "  Running: Yes"; \
		ps aux | grep xlsx-to-ics | grep -v grep; \
	else \
		echo "  Running: No"; \
	fi
	@echo "  Port 8080: $(shell netstat -tlnp 2>/dev/null | grep :8080 || echo 'Not listening')"
	@echo "  Health check: $(shell curl -s http://localhost:8080/health 2>/dev/null || echo 'Not responding')"
