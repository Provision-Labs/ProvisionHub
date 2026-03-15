COMPOSE_BASE = -f deployments/docker-compose.yaml
COMPOSE_QAS = $(COMPOSE_BASE) -f deployments/docker-compose.dev.yaml
COMPOSE_PROD = $(COMPOSE_BASE) -f deployments/docker-compose.prod.yaml

# Dev

# Qas
qas:
	docker compose $(COMPOSE_QAS) up -d

qas-build:
	docker compose $(COMPOSE_QAS) up -d --build

qas-down:
	docker compose $(COMPOSE_QAS) down

qas-logs:
	docker compose $(COMPOSE_QAS) logs -f

# Prod
prod:
	docker compose $(COMPOSE_PROD) up -d

prod-down:
	docker compose $(COMPOSE_PROD) down

# Utils
qas-ps:
	docker compose $(COMPOSE_QAS) ps

qas-clean:
	docker compose $(COMPOSE_QAS) down -v  # removes volumes too, careful