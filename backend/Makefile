.PHONY: help run build test clean migrate-up migrate-down migrate-status migrate-create migrate-reset db-setup

# Variables
APP_NAME=portfolio-api
DB_PATH=portfolio.db
MIGRATIONS_DIR=migrations
CMD_DIR=cmd/api

# Colors for output
BLUE=\033[0;34m
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

run: ## Run the application
	@echo "$(BLUE)Starting application...$(NC)"
	@go run $(CMD_DIR)/main.go

build: ## Build the application
	@echo "$(BLUE)Building application...$(NC)"
	@go build -o $(APP_NAME) $(CMD_DIR)/main.go
	@echo "$(GREEN)Build complete: $(APP_NAME)$(NC)"

test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...

clean: ## Clean build artifacts and database
	@echo "$(YELLOW)Cleaning up...$(NC)"
	@rm -f $(APP_NAME)
	@rm -f $(DB_PATH)
	@echo "$(GREEN)Cleanup complete$(NC)"

migrate-up: ## Run all pending migrations
	@echo "$(BLUE)Running migrations...$(NC)"
	@goose -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) up
	@echo "$(GREEN)Migrations complete$(NC)"

migrate-down: ## Rollback the last migration
	@echo "$(YELLOW)Rolling back last migration...$(NC)"
	@goose -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) down
	@echo "$(GREEN)Rollback complete$(NC)"

migrate-status: ## Show migration status
	@echo "$(BLUE)Migration status:$(NC)"
	@goose -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) status

migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)Error: NAME is required. Usage: make migrate-create NAME=migration_name$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)Creating migration: $(NAME)$(NC)"
	@goose -dir $(MIGRATIONS_DIR) create $(NAME) sql
	@echo "$(GREEN)Migration created$(NC)"

migrate-reset: ## Reset database (rollback all migrations and re-run them)
	@echo "$(YELLOW)Resetting database...$(NC)"
	@goose -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) reset
	@goose -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) up
	@echo "$(GREEN)Database reset complete$(NC)"

db-setup: ## Setup database (create and run migrations)
	@echo "$(BLUE)Setting up database...$(NC)"
	@if [ -f $(DB_PATH) ]; then \
		echo "$(YELLOW)Database already exists at $(DB_PATH)$(NC)"; \
		read -p "Do you want to reset it? [y/N] " answer; \
		if [ "$$answer" = "y" ] || [ "$$answer" = "Y" ]; then \
			$(MAKE) clean; \
			$(MAKE) migrate-up; \
		else \
			echo "$(GREEN)Keeping existing database$(NC)"; \
		fi \
	else \
		$(MAKE) migrate-up; \
	fi
	@echo "$(GREEN)Database setup complete$(NC)"

install-tools: ## Install required tools (goose)
	@echo "$(BLUE)Installing goose...$(NC)"
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "$(GREEN)Tools installed$(NC)"

deps: ## Download dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Dependencies downloaded$(NC)"

dev: ## Run in development mode with auto-reload (requires air)
	@if ! command -v air > /dev/null; then \
		echo "$(YELLOW)Installing air for hot reload...$(NC)"; \
		go install github.com/air-verse/air@latest; \
	fi
	@echo "$(BLUE)Starting development server with hot reload...$(NC)"
	@air

swagger: ## Generate swagger documentation
	@echo "$(BLUE)Generating swagger documentation...$(NC)"
	@swag init -g cmd/api/main.go -o docs
	@echo "$(GREEN)Swagger documentation generated at docs/$(NC)"
	@echo "$(YELLOW)View at: http://localhost:8080/swagger/index.html$(NC)"

swagger-install: ## Install swag CLI tool
	@echo "$(BLUE)Installing swag...$(NC)"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "$(GREEN)Swag installed$(NC)"