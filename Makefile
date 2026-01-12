.PHONY: help setup build run stop logs dev test clean

# Default compose command
COMPOSE := $(shell command -v docker-compose 2> /dev/null || echo "docker compose")

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

setup: ## Run the setup wizard
	@bash scripts/setup_docker.sh

build: ## Build Docker image
	@$(COMPOSE) build

run: ## Start the bot
	@$(COMPOSE) up -d
	@echo "✓ Bot running at http://localhost:8080"

stop: ## Stop the bot
	@$(COMPOSE) down

logs: ## Show logs
	@$(COMPOSE) logs -f

dev: ## Start with ngrok tunnel for development
	@$(COMPOSE) --profile dev up -d --build
	@sleep 3
	@echo "✓ Bot running!"
	@echo "  Local:  http://localhost:8080/webhook"
	@curl -s http://localhost:4040/api/tunnels 2>/dev/null | python3 -c "import sys,json; t=json.load(sys.stdin)['tunnels']; print('  Public:', t[0]['public_url']+'/webhook' if t else 'ngrok not ready')" || echo "  Ngrok: http://localhost:4040"

test: ## Send a test message (usage: make test PHONE=14155551234)
	@bash scripts/send_test.sh $(PHONE)

restart: ## Restart the bot
	@$(COMPOSE) restart

rebuild: ## Rebuild and restart
	@$(COMPOSE) up -d --build

shell: ## Open shell in container
	@$(COMPOSE) exec whatsapp-bot sh

clean: ## Remove containers and images
	@$(COMPOSE) down --rmi local -v
	@echo "✓ Cleaned up"

status: ## Show container status
	@$(COMPOSE) ps
