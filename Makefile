# Digital Wallet Demo - Root Makefile
# Unified commands for managing the entire microservices system

.PHONY: help
help: ## Show this help message
	@echo "Digital Wallet Demo - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# =============================================================================
# DOCKER COMPOSE COMMANDS
# =============================================================================

.PHONY: up
up: ## Start all services with Docker Compose
	docker-compose up -d
	@echo "✅ All services are starting up..."
	@echo "📊 Wallet Service: http://localhost:8081"
	@echo "💳 Transaction Service: http://localhost:8082"
	@echo "🌐 Kong API Gateway: http://localhost:8000"
	@echo "📝 Redis Cache: localhost:6379"

.PHONY: down
down: ## Stop all services
	docker-compose down
	@echo "🛑 All services stopped"

.PHONY: restart
restart: ## Restart all services
	docker-compose restart
	@echo "🔄 All services restarted"

.PHONY: logs
logs: ## Show logs for all services
	docker-compose logs -f

.PHONY: status
status: ## Show status of all services
	docker-compose ps

.PHONY: clean
clean: ## Stop and remove all containers, networks, and volumes
	docker-compose down -v --remove-orphans
	docker system prune -f
	@echo "🧹 System cleaned up"

# =============================================================================
# DATABASE COMMANDS
# =============================================================================

.PHONY: migrate
migrate: ## Run database migrations for both services
	@echo "🗄️  Running wallet service migrations..."
	@cd services/wallets && make migrate
	@echo "🗄️  Running transaction service migrations..."
	@cd services/transactions && make migrate
	@echo "✅ All migrations completed"

.PHONY: reset-db
reset-db: ## Reset databases for both services
	@echo "🔄 Resetting wallet database..."
	@cd services/wallets && make reset-db
	@echo "🔄 Resetting transaction database..."
	@cd services/transactions && make reset-db
	@echo "✅ All databases reset"

# =============================================================================
# TESTING COMMANDS
# =============================================================================

.PHONY: test
test: ## Run tests for both services
	@echo "🧪 Running wallet service tests..."
	@cd services/wallets && make test
	@echo "🧪 Running transaction service tests..."
	@cd services/transactions && make test
	@echo "✅ All tests completed"

.PHONY: test-wallet
test-wallet: ## Run tests for wallet service only
	@cd services/wallets && make test

.PHONY: test-transaction
test-transaction: ## Run tests for transaction service only
	@cd services/transactions && make test

.PHONY: test-ci
test-ci: ## Run CI tests for both services with coverage
	@echo "🧪 Running wallet service CI tests..."
	@cd services/wallets && make test-ci
	@echo "🧪 Running transaction service CI tests..."
	@cd services/transactions && make test-ci
	@echo "✅ All CI tests completed"

# =============================================================================
# DEVELOPMENT COMMANDS
# =============================================================================

.PHONY: serve-wallet
serve-wallet: ## Start wallet service locally
	@cd services/wallets && make serve

.PHONY: serve-transaction
serve-transaction: ## Start transaction service locally
	@cd services/transactions && make serve

.PHONY: build
build: ## Build Docker images for all services
	docker-compose build
	@echo "🔨 All services built"

.PHONY: swagger-wallet
swagger-wallet: ## Generate Swagger docs for wallet service
	@cd services/wallets && make swagger
	@echo "📚 Wallet service Swagger docs generated"

.PHONY: swagger-transaction
swagger-transaction: ## Generate Swagger docs for transaction service
	@cd services/transactions && make swagger
	@echo "📚 Transaction service Swagger docs generated"

.PHONY: swagger
swagger: swagger-wallet swagger-transaction ## Generate Swagger docs for all services

# =============================================================================
# HEALTH CHECK COMMANDS
# =============================================================================

.PHONY: health
health: ## Check health of all services
	@echo "🏥 Checking service health..."
	@echo "Wallet Service:"
	@curl -s http://localhost:8081/health || echo "❌ Wallet service not responding"
	@echo "\nTransaction Service:"
	@curl -s http://localhost:8082/health || echo "❌ Transaction service not responding"
	@echo "\nKong Gateway:"
	@curl -s http://localhost:8000 || echo "❌ Kong gateway not responding"

.PHONY: redis-check
redis-check: ## Check Redis cache status
	@echo "📝 Checking Redis cache..."
	@docker exec redis_cache redis-cli ping || echo "❌ Redis not responding"
	@docker exec redis_cache redis-cli KEYS "*" | head -10

# =============================================================================
# QUICK START COMMANDS
# =============================================================================

.PHONY: dev
dev: clean up migrate ## Quick development setup (clean + up + migrate)
	@echo "🚀 Development environment ready!"
	@echo "📖 Run 'make health' to verify all services are running"

.PHONY: demo
demo: dev ## Setup demo environment with sample data
	@echo "🎭 Setting up demo environment..."
	@sleep 10  # Wait for services to be ready
	@echo "💰 Creating demo wallet..."
	@curl -X POST http://localhost:8081/api/v1/wallets \
		-H "Content-Type: application/json" \
		-d '{"user_id":"demo-user","acnt_type":"user"}' || true
	@echo "💵 Making sample deposit..."
	@curl -X POST http://localhost:8081/api/v1/wallets/deposit \
		-H "Content-Type: application/json" \
		-d '{"user_id":"demo-user","amount":1000}' || true
	@echo "\n✅ Demo environment ready with sample data!"
	@echo "🔗 Try: curl http://localhost:8081/api/v1/wallets/demo-user"

# Default target
.DEFAULT_GOAL := help