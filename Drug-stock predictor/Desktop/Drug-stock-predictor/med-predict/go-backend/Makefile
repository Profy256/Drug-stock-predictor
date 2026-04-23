.PHONY: help dev build run test docker-up docker-down clean

help:
	@echo "Available commands:"
	@echo "  make dev          - Run in development mode with hot reload"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the built application"
	@echo "  make test         - Run tests"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make clean        - Clean build artifacts"

dev:
	docker-compose up -d db
	go run ./cmd/api/main.go

build:
	go build -o bin/med-predict ./cmd/api

run: build
	./bin/med-predict

test:
	go test ./...

docker-build:
	docker build -t med-predict:latest .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f backend

clean:
	rm -rf bin/
	go clean

deps:
	go mod download
	go mod tidy

fmt:
	go fmt ./...

lint:
	golangci-lint run

init-db:
	docker-compose exec db psql -U postgres -d med_predict -f /docker-entrypoint-initdb.d/001_init_schema.sql
