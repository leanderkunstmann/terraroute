.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m - %s\n", $$1, $$2}'

.PHONY: backend
backend: ## Run the backend server
	@echo "Starting backend server..."
	@if [ -f "backend/.env" ]; then \
		echo "Loading environment variables from backend/.env"; \
		source backend/.env; \
	fi && \
	cd backend && go run . --config config.json

.PHONY: frontend
frontend: $(NODE_MODULES) ## Run the frontend server
	@echo "Starting frontend server..."
	@cd frontend && npm run dev

NODE_MODULES := frontend/node_modules
$(NODE_MODULES): ## Install frontend dependencies
	@echo "Installing frontend dependencies..."
	@cd frontend && npm install
