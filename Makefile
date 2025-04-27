SHELL := /bin/sh

###############################################################################

# ----------------- ENVS -----------------

include .env
export

# ----------------- GENERAL -----------------

install:
	go mod download

tests:
	go test -v ./test/...

lint:
	golangci-lint run

coverage:
	go test -coverprofile=cover.out ./test/...
	go tool cover -html=cover.out

# ----------------- DOCKER -----------------

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-down-v:
	docker compose down -v

docker-logs:
	docker compose logs -f

docker-ps:
	docker compose ps

docker-test:
	docker compose run --rm test

# ----------------- DEVELOPMENT -----------------

run-dev:
	go run ./cmd/service/main.go

# ----------------- MIGRATIONS -----------------

migrate-up:
	@echo "Running migrations..."
	docker compose run --rm migrate -path /migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/${POSTGRES_DB}?sslmode=disable" up

migrate-down:
	@echo "Rolling back migrations..."
	docker compose run --rm migrate -path /migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/${POSTGRES_DB}?sslmode=disable" down 1

# ----------------- CONFIGS -----------------

setup-env:
	@if [ ! -f .env ]; then \
		echo "Creating .env file from .env.example"; \
		cp .env.example .env; \
	else \
		echo ".env file already exists"; \
	fi

# ----------------- HELPERS -----------------

GREEN=\033[0;32m
RESET=\033[0m

help:
	@echo -e "$(GREEN)Available commands:$(RESET)"
	@echo "  install           - Download Go dependencies"
	@echo "  tests             - Run tests"
	@echo "  lint              - Run linter"
	@echo "  coverage          - Generate test coverage report"
	@echo ""
	@echo "  docker-build      - Build Docker images"
	@echo "  docker-up         - Start Docker containers"
	@echo "  docker-down       - Stop Docker containers"
	@echo "  docker-down-v     - Stop Docker containers and remove volumes"
	@echo "  docker-logs       - Show Docker container logs"
	@echo "  docker-ps         - Show running Docker containers"
	@echo "  docker-test       - Run tests in Docker container"
	@echo ""
	@echo "  run-dev           - Run application locally"
	@echo ""
	@echo "  migrate-up        - Apply database migrations"
	@echo "  migrate-down      - Rollback last database migration"
	@echo ""
	@echo "  setup-env         - Create .env file from .env.example if it doesn't exist"