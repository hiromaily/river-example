# Load .env file if it exists
-include .env
export

# Environment variables
export DATABASE_URL ?= postgres://postgres:postgres@localhost:15432/river?sslmode=disable

# Colors for output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m


.PHONY: install-river
install-river: ## Install River CLI tool
	@echo "$(GREEN)Installing River CLI...$(RESET)"
	go install github.com/riverqueue/river/cmd/river@latest

##@ Database
.PHONY: db-migrate
db-migrate: ## Run database migrations
	@echo "$(GREEN)Running database migrations...$(RESET)"
	river migrate-up --database-url "$(DATABASE_URL)"

.PHONY: db-migrate-down
db-migrate-down: ## Rollback database migrations
	@echo "$(YELLOW)Rolling back database migrations...$(RESET)"
	river migrate-down --database-url "$(DATABASE_URL)"

.PHONY: run-producer
run-producer: ## Run the job producer
	@echo "$(GREEN)Starting producer...$(RESET)"
	go run cmd/producer/main.go

.PHONY: run-worker
run-worker: ## Run the job worker
	@echo "$(GREEN)Starting worker...$(RESET)"
	go run cmd/worker/main.go
