# Makefile for WhoKnows Variations

# ===============================
# Variables
# ===============================
DEV_COMPOSE = docker compose -f docker-compose.dev.yml
AIR_FILE = .air.toml
SWAG ?= swag

# ===============================
# Helper: check if Docker daemon is running
# ===============================
define check_docker
	@docker info > /dev/null 2>&1 || { \
		echo "[FAIL] Docker does not appear to be running."; \
		echo "[INFO] Please start Docker Desktop (or the Docker daemon) and try again."; \
		exit 1; \
	}
endef

# ===============================
# Help
# ===============================
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make dev          Run development environment with auto-reload"
	@echo "  make dev-d        Run development environment with auto-reload in detached mode"
	@echo "  make stop-dev     Stop development environment"
	@echo "  make reset-dev    Stop and remove dev volumes"
	@echo "  make dev-test     Run development environment with auto-reload and tests"
	@echo "  make swagger      Generate swagger.json/swagger.yaml from annotated handlers"

# ===============================
# Development mode (auto-reload with Air)
# ===============================
.PHONY: dev
dev:
	$(call check_docker)
	@if [ ! -f $(AIR_FILE) ]; then \
		echo "[FAIL] .air.toml not found! Please create it before running dev mode."; \
		exit 1; \
	fi
	@echo "[STARTING] Starting development environment with Air..."
	$(DEV_COMPOSE) up --build whoknows_variations_dev nginx_dev postgres grafana

.PHONY: dev-d
dev-d:
	$(call check_docker)
	@if [ ! -f $(AIR_FILE) ]; then \
		echo "[FAIL] .air.toml not found! Please create it before running dev mode."; \
		exit 1; \
	fi
	@echo "[STARTING] Starting development environment with Air in detached mode..."
	$(DEV_COMPOSE) up -d --build whoknows_variations_dev nginx_dev postgres grafana

.PHONY: stop-dev
stop-dev:
	@echo "[STOPPING] Stopping development environment..."
	$(DEV_COMPOSE) down
	@echo "[DONE] Development environment stopped."

.PHONY: reset-dev
reset-dev:
	@echo "[STOPPING] Removing dev containers and volumes..."
	$(DEV_COMPOSE) down -v
	@echo "[DONE] Dev containers and volumes removed."

.PHONY: dev-test
dev-test:
	$(call check_docker)
	@if [ ! -f $(AIR_FILE) ]; then \
		echo "[FAIL] .air.toml not found! Please create it before running dev mode."; \
		exit 1; \
	fi
	@echo "[STARTING] Starting development environment with Air and tests..."
	$(DEV_COMPOSE) run --rm test

.PHONY: swagger
swagger:
	@echo "[GEN] Generating swagger docs..."
	$(SWAG) init -g cmd/main.go -o docs/api --parseDependency --parseInternal --outputTypes json,yaml
