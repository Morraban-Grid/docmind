.PHONY: help up down logs clean test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

up: ## Start all services
	docker compose -f deployments/docker/docker-compose.yml up -d

down: ## Stop all services
	docker compose -f deployments/docker/docker-compose.yml down

logs: ## View logs from all services
	docker compose -f deployments/docker/docker-compose.yml logs -f

ps: ## Show status of all services
	docker compose -f deployments/docker/docker-compose.yml ps

clean: ## Stop and remove all containers, volumes, and images
	docker compose -f deployments/docker/docker-compose.yml down -v --rmi all

restart: ## Restart all services
	docker compose -f deployments/docker/docker-compose.yml restart

build: ## Rebuild all services
	docker compose -f deployments/docker/docker-compose.yml build

test-go: ## Run Go tests
	cd services/go-user-service && go test ./...

test-python: ## Run Python tests
	cd services/python-rag-service && pytest

test: test-go test-python ## Run all tests

init: ## Initialize project (copy .env.example to .env)
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env file created. Please edit it with your configuration."; \
	else \
		echo ".env file already exists."; \
	fi
