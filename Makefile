# Shawty-UR Makefile

.PHONY: help docker-up docker-down docker-logs migrate-up migrate-down migrate-create migrate-status build run clean

# Variables
APP_NAME=shawty-ur
DOCKER_COMPOSE=docker-compose
GOOSE=goose
DB_DSN=postgres://howl:turnip_man1234@localhost:5432/social?sslmode=disable
MIGRATIONS_DIR=./migrations

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## Docker Commands
docker-up: ## Start PostgreSQL and Redis containers
	$(DOCKER_COMPOSE) up -d
	@echo "Waiting for databases to be ready..."
	@sleep 5
	@echo "PostgreSQL and Redis are running!"

docker-down: ## Stop and remove containers
	$(DOCKER_COMPOSE) down

docker-logs: ## View container logs
	$(DOCKER_COMPOSE) logs -f

docker-clean: ## Remove containers and volumes
	$(DOCKER_COMPOSE) down -v

## Database Migration Commands
migrate-install: ## Install goose migration tool
	@echo "Installing goose..."
	go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "Goose installed successfully!"

migrate-up: ## Run all pending migrations
	@echo "Running migrations..."
	cd $(MIGRATIONS_DIR) && $(GOOSE) postgres "$(DB_DSN)" up

migrate-down: ## Rollback the last migration
	@echo "Rolling back migration..."
	cd $(MIGRATIONS_DIR) && $(GOOSE) postgres "$(DB_DSN)" down

migrate-reset: ## Rollback all migrations
	@echo "Resetting database..."
	cd $(MIGRATIONS_DIR) && $(GOOSE) postgres "$(DB_DSN)" reset

migrate-status: ## Show migration status
	cd $(MIGRATIONS_DIR) && $(GOOSE) postgres "$(DB_DSN)" status

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=create_users_table)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=create_users_table"; \
		exit 1; \
	fi
	cd $(MIGRATIONS_DIR) && $(GOOSE) create $(NAME) sql

## Application Commands
deps: ## Download Go dependencies
	go mod download
	go mod tidy

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) ./api

run: ## Run the application
	@echo "Starting $(APP_NAME)..."
	./bin/$(APP_NAME)

dev: ## Run application in development mode (without building)
	go run ./api

clean: ## Clean build artifacts
	rm -rf bin/

## Combined Commands
setup: docker-up migrate-install migrate-up ## Complete setup: start Docker, install goose, run migrations
	@echo "Setup complete! Run 'make build && make run' to start the app"

start: docker-up build run ## Start everything (Docker + App)

restart: docker-down docker-up ## Restart Docker containers

test: ## Run tests
	go test -v ./...

lint: ## Run linter
	golangci-lint run
