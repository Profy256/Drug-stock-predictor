.PHONY: help build run test migrate docker-up docker-down clean lint fmt

help:
	@echo "Med Predict Backend - Available Commands:"
	@echo ""
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make dev           - Run with hot reload (requires air)"
	@echo "  make test          - Run tests"
	@echo "  make migrate       - Run database migrations"
	@echo "  make docker-up     - Start Docker services"
	@echo "  make docker-down   - Stop Docker services"
	@echo "  make docker-logs   - View Docker logs"
	@echo "  make lint          - Run linter (requires golangci-lint)"
	@echo "  make fmt           - Format code"
	@echo "  make clean         - Remove build artifacts"

build:
	@echo "Building application..."
	@go build -o med-predict-backend cmd/server/main.go
	@echo "Build complete: ./med-predict-backend"

run:
	@echo "Running application..."
	@go run cmd/server/main.go

dev:
	@echo "Running with hot reload..."
	@air

test:
	@echo "Running tests..."
	@go test -v ./...

migrate:
	@echo "Running migrations..."
	@psql -U postgres -d medpredict -f migrations/001_init_schema.sql

docker-build:
	@echo "Building Docker image..."
	@docker build -t med-predict-backend:latest .

docker-up:
	@echo "Starting Docker services..."
	@docker-compose up -d

docker-down:
	@echo "Stopping Docker services..."
	@docker-compose down

docker-down-v:
	@echo "Stopping Docker services and removing volumes..."
	@docker-compose down -v

docker-logs:
	@docker-compose logs -f

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	@go fmt ./...

clean:
	@echo "Cleaning..."
	@rm -f med-predict-backend
	@go clean

install-deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

psql:
	@psql -U postgres -d medpredict

db-reset:
	@echo "Resetting database..."
	@psql -U postgres -d medpredict -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@$(MAKE) migrate
