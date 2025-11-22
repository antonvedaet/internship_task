.PHONY: build run stop clean test-db health migrate lint local-build local-run

build:
	docker-compose build

run:
	docker-compose up

run-detached:
	docker-compose up -d

stop:
	docker-compose down

clean:
	docker-compose down -v

test-db:
	docker-compose exec db psql -U postgres -d pr_reviewer -c "SELECT version();"

health:
	curl http://localhost:8080/health

migrate:
	docker-compose exec db psql -U postgres -d pr_reviewer -f /docker-entrypoint-initdb.d/001_init.sql

lint:
	golangci-lint run --config .golangci.yml

local-build:
	go build -o bin/server ./cmd/server

local-run:
	go run ./cmd/server

help:
	@echo "Available commands:"
	@echo "  build         - Build docker images"
	@echo "  run           - Start the application"
	@echo "  stop          - Stop the application"
	@echo "  clean         - Stop and remove volumes"
	@echo "  test-db       - Test database connection"
	@echo "  health        - Health check"
	@echo "  migrate       - Run database migrations"
	@echo "  lint          - Run linter"
	@echo "  local-build   - Build locally (postgres required) "
	@echo "  local-run     - Run locally (postgres required)"
	@echo "  help          - Show this help"