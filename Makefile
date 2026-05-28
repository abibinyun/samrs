.PHONY: help up down build logs restart clean seed test

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

up: ## Start all services
	docker-compose up -d

up-dev: ## Start all services in development mode (with hot reload)
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

up-tools: ## Start all services with tools (pgAdmin)
	docker-compose --profile tools up -d

up-monitoring: ## Start all services with monitoring
	docker-compose --profile monitoring up -d

up-all: ## Start all services with tools and monitoring
	docker-compose --profile tools --profile monitoring up -d

down: ## Stop all services
	docker-compose down

build: ## Build all services
	docker-compose build

rebuild: ## Rebuild and restart all services
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

logs: ## Show logs (all services)
	docker-compose logs -f

logs-backend: ## Show backend logs
	docker-compose logs -f backend

logs-frontend: ## Show frontend logs
	docker-compose logs -f frontend

logs-bff: ## Show BFF logs
	docker-compose logs -f bff

restart: ## Restart all services
	docker-compose restart

restart-backend: ## Restart backend only
	docker-compose restart backend

restart-frontend: ## Restart frontend only
	docker-compose restart frontend

ps: ## Show running services
	docker-compose ps

seed: ## Run database seed
	docker compose exec backend ./seed

test: ## Run backend tests
	docker-compose exec backend go test ./...

test-coverage: ## Run backend tests with coverage
	docker-compose exec backend go test -cover ./...

shell-backend: ## Open backend shell
	docker-compose exec backend sh

shell-db: ## Open database shell
	docker-compose exec postgres psql -U samrs_user -d samrs_db

clean: ## Stop and remove all containers, volumes, and images
	docker-compose down -v --rmi all --remove-orphans

backup-db: ## Backup database
	docker-compose exec postgres pg_dump -U samrs_user samrs_db > backup_$$(date +%Y%m%d_%H%M%S).sql
	@echo "Database backed up to backup_$$(date +%Y%m%d_%H%M%S).sql"

restore-db: ## Restore database (usage: make restore-db FILE=backup.sql)
	@if [ -z "$(FILE)" ]; then echo "Usage: make restore-db FILE=backup.sql"; exit 1; fi
	docker-compose exec -T postgres psql -U samrs_user -d samrs_db < $(FILE)
	@echo "Database restored from $(FILE)"

health: ## Check services health
	@echo "Checking services..."
	@curl -s http://localhost:8090/ping > /dev/null && echo "Backend: UP" || echo "Backend: DOWN"
	@curl -s http://localhost:5173 > /dev/null && echo "Frontend: UP" || echo "Frontend: DOWN"
	@curl -s http://localhost:3000 > /dev/null && echo "BFF: UP" || echo "BFF: DOWN"
