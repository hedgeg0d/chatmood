.PHONY: help build run dev test clean docker-build docker-run deploy setup deps

# Default target
help:
	@echo "ChatMood - Telegram Mini App"
	@echo ""
	@echo "Available commands:"
	@echo "  setup          - Initial project setup"
	@echo "  deps           - Download Go dependencies"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application locally"
	@echo "  dev            - Run in development mode with auto-reload"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker"
	@echo "  docker-dev     - Run development environment with Docker"
	@echo "  deploy         - Deploy to production"
	@echo "  lint           - Run linters"
	@echo "  format         - Format code"
	@echo "  security       - Run security checks"

# Initial setup
setup:
	@echo "Setting up ChatMood project..."
	@cp .env.example .env
	@echo "Created .env file - please update with your values"
	@go mod download
	@echo "Dependencies downloaded"
	@echo "Setup complete! Update .env with your Telegram bot token"

# Download dependencies
deps:
	@echo "Downloading Go dependencies..."
	@go mod download
	@go mod tidy

# Build the application
build:
	@echo "Building ChatMood..."
	@go build -o bin/chatmood cmd/server/main.go
	@echo "Build complete: bin/chatmood"

# Run the application
run: build
	@echo "Starting ChatMood server..."
	@./bin/chatmood

# Development mode with file watching
dev:
	@echo "Starting development server..."
	@if command -v air > /dev/null 2>&1; then \
		air; \
	else \
		echo "Installing air for live reload..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...
	@go test -race ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f chatmood
	@rm -f main
	@rm -f coverage.out coverage.html
	@go clean

# Format code
format:
	@echo "Formatting code..."
	@go fmt ./...
	@if command -v goimports > /dev/null 2>&1; then \
		goimports -w .; \
	fi

# Lint code
lint:
	@echo "Running linters..."
	@if command -v golangci-lint > /dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

# Security checks
security:
	@echo "Running security checks..."
	@if command -v gosec > /dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "Installing gosec..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t chatmood:latest .

docker-run: docker-build
	@echo "Running with Docker..."
	@docker run -p 8080:8080 --env-file .env chatmood:latest

docker-dev:
	@echo "Starting development environment with Docker..."
	@docker-compose up --build

docker-prod:
	@echo "Starting production environment with Docker..."
	@docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Stop Docker containers
docker-stop:
	@echo "Stopping Docker containers..."
	@docker-compose down

# View Docker logs
docker-logs:
	@docker-compose logs -f chatmood

# Database commands (for future use)
db-migrate:
	@echo "Running database migrations..."
	@echo "Database migrations not implemented yet"

db-seed:
	@echo "Seeding database..."
	@echo "Database seeding not implemented yet"

# Deployment
deploy:
	@echo "Deploying to production..."
	@echo "Deployment script not implemented yet"
	@echo "Manual deployment steps:"
	@echo "1. Build: make build"
	@echo "2. Copy binary and web/ to server"
	@echo "3. Set environment variables"
	@echo "4. Run: ./chatmood"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "Development tools installed"

# Generate project statistics
stats:
	@echo "Project Statistics:"
	@echo "==================="
	@echo "Go files:" $$(find . -name "*.go" | wc -l)
	@echo "Lines of Go code:" $$(find . -name "*.go" -exec wc -l {} + | tail -n1 | awk '{print $$1}')
	@echo "JavaScript files:" $$(find . -name "*.js" | wc -l)
	@echo "Lines of JS code:" $$(find . -name "*.js" -exec wc -l {} + | tail -n1 | awk '{print $$1}')
	@echo "Total files:" $$(find . -type f | grep -v ".git" | wc -l)

# Create release
release:
	@echo "Creating release build..."
	@mkdir -p releases
	@GOOS=linux GOARCH=amd64 go build -o releases/chatmood-linux-amd64 cmd/server/main.go
	@GOOS=darwin GOARCH=amd64 go build -o releases/chatmood-darwin-amd64 cmd/server/main.go
	@GOOS=windows GOARCH=amd64 go build -o releases/chatmood-windows-amd64.exe cmd/server/main.go
	@cp -r web releases/
	@echo "Release builds created in releases/ directory"

# Backup
backup:
	@echo "Creating backup..."
	@tar -czf chatmood-backup-$$(date +%Y%m%d-%H%M%S).tar.gz \
		--exclude='.git' \
		--exclude='bin' \
		--exclude='releases' \
		--exclude='*.log' \
		--exclude='.env' \
		.
	@echo "Backup created"

# Show project info
info:
	@echo "ChatMood Project Information"
	@echo "============================"
	@echo "Go version:" $$(go version)
	@echo "Project path:" $$(pwd)
	@echo "Git branch:" $$(git branch --show-current 2>/dev/null || echo "Not a git repository")
	@echo "Git commit:" $$(git rev-parse --short HEAD 2>/dev/null || echo "Not a git repository")
	@echo "Port: $$(grep PORT .env 2>/dev/null | cut -d'=' -f2 || echo '8080')"
