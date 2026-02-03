.PHONY: help run build migrate-up migrate-down migrate-status migrate-create test clean

# Load .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

run: ## Run the application
	go run web/cmd/api/main.go

run-worker: ## Run the worker service
	go run worker/cmd/api/main.go

run-ui:
	cd web/ui && npm run dev

run-landing:
	cd landing && npm run dev

build: ## Build the application
	go build -o bin/api web/cmd/api/main.go

build-worker: ## Build the worker service
	go build -o bin/worker worker/cmd/api/main.go

build-all: ## Build both services
	go build -o bin/api web/cmd/api/main.go
	go build -o bin/worker worker/cmd/api/main.go

migrate-up: ## Run database migrations
	dbmate --migrations-dir ./migrations up

migrate-down: ## Rollback last migration
	dbmate --migrations-dir ./migrations down

migrate-status: ## Show migration status
	dbmate --migrations-dir ./migrations status

migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	dbmate --migrations-dir ./migrations new $(NAME)

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

deps: ## Download dependencies
	go mod download
	go mod tidy

dev: ## Run in development mode with auto-reload (requires air)
	air

.DEFAULT_GOAL := help
