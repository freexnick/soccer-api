include ./configs/.env

GIT_COMMIT_SHA ?= $(shell git rev-parse --short HEAD)
APP_VERSION ?= $(shell git describe --tags --always --dirty --match='v*' 2>/dev/null || echo "dev")
MIGRATION_PATH ?= internal/infrastructure/database/migrations
MIGRATE_IMAGE ?= migrate/migrate:v4.18.3
PROJECT_NAME ?= soccer-api
COMPOSE_FILE := ./deployments/docker-compose.yaml
STATIC_CHECKS := all,-ST1000

.PHONY: lint create-migration apply-migrations rollback-migration run build stop logs

help:
	@echo "Usage:"
	@echo "  make create-migration NAME=<migration_name>  Create a new migration file"
	@echo "  make apply-migrations                        Apply all pending migrations"
	@echo "  make rollback-migration                      Rollback the last migration"
	@echo "  make lint                                    Lint project"
	@echo "  make run                                     Start containers"
	@echo "  make logs                                    Container logs"
	@echo "  make stop                                    Stop containers"

check-db-url:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Error: DB_URL is not set. Please define DB_URL in your environment or ./configs/.env"; \
		echo "Example: export DB_URL=\"postgresql://user:pass@db_host:port/dbname?sslmode=disable&search_path=public\""; \
		exit 1; \
	fi

create-migration:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: Please provide a migration name with NAME=<migration_name>"; \
		exit 1; \
	fi
	@echo "Creating new migration: $(NAME) in $(MIGRATION_PATH)"
	@mkdir -p $(MIGRATION_PATH)
	# Run migration tool using official Docker image for golang-migrate
	@docker run --rm \
		-v $(PWD)/$(MIGRATION_PATH):/app/migrations \
		$(MIGRATE_IMAGE) create -ext sql -dir /app/migrations -seq $(NAME)
	@echo "Migration file created. Please edit the .up.sql and .down.sql files."

apply-migrations: check-db-url
	@echo "Applying migrations to $(DB_URL) using path $(MIGRATION_PATH)..."
	@docker run --rm \
		--network $(PROJECT_NAME) \
		-v $(PWD)/$(MIGRATION_PATH):/app/migrations \
		$(MIGRATE_IMAGE) \
		-path /app/migrations -database "$(DB_URL)" up
	@echo "Migrations applied."

rollback-migration: check-db-url
	@echo "Rolling back the last migration on $(DB_URL) using path $(MIGRATION_PATH)..."
	@docker run --rm \
		--network $(PROJECT_NAME) \
		-v "$(PWD)/$(MIGRATION_PATH):/app/migrations" \
		$(MIGRATE_IMAGE) \
		-path=/app/migrations -database "$(DB_URL)" down 1
	@echo "Last migration rolled back."

lint:
	@echo "Running go vet..."
	@go vet ./...

	@echo "Running staticcheck..."
	@go run honnef.co/go/tools/cmd/staticcheck@latest -checks=${STATIC_CHECKS} ./...

run:
	@echo "starting services..."
	@echo "Using APP_VERSION=$(APP_VERSION) and GIT_COMMIT_SHA=$(GIT_COMMIT_SHA) for build."
	APP_VERSION=$(APP_VERSION) GIT_COMMIT_SHA=$(GIT_COMMIT_SHA) \
		docker-compose -f $(COMPOSE_FILE) -p $(PROJECT_NAME) up --build -d

build:
	@echo "Building Docker images with APP_VERSION=$(APP_VERSION) and GIT_COMMIT_SHA=$(GIT_COMMIT_SHA)..."
	APP_VERSION=$(APP_VERSION) GIT_COMMIT_SHA=$(GIT_COMMIT_SHA) \
		docker-compose -f $(COMPOSE_FILE) -p $(PROJECT_NAME) build --no-cache server swagger-ui

stop:
	@echo "Stopping and removing services, networks, and volumes..."
	docker-compose -f $(COMPOSE_FILE) -p $(PROJECT_NAME) down -v --remove-orphans

logs:
	docker-compose -f $(COMPOSE_FILE) -p $(PROJECT_NAME) logs -f --tail=50

.DEFAULT_GOAL := help