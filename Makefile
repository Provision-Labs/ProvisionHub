COMPOSE_BASE = -f deployments/docker-compose.yaml
COMPOSE_QAS = $(COMPOSE_BASE) -f deployments/docker-compose.dev.yaml
COMPOSE_PROD = $(COMPOSE_BASE) -f deployments/docker-compose.prod.yaml
ENV_FILE = --env-file deployments/.env.local

.RECIPEPREFIX := >

# Dev

dev-infra-up:
> docker compose $(ENV_FILE) $(COMPOSE_QAS) up -d keycloak-database keycloak cp-database

dev-infra-down:
> docker compose $(ENV_FILE) $(COMPOSE_QAS) down

dev-cp:
> cd apps/control-plane && \
> air --build.cmd "go build -o ./tmp/server ./cmd/server" --build.bin "./tmp/server"

dev-web:
> cd apps/web && pnpm i && pnpm dev

dev-local:
> @$(MAKE) dev-infra-up
> @set -e; \
> ( $(MAKE) dev-cp ) & pid_cp=$$!; \
> ( $(MAKE) dev-web ) & pid_web=$$!; \
> trap 'kill $$pid_cp $$pid_web 2>/dev/null || true' INT TERM EXIT; \
> wait $$pid_cp; \
> wait $$pid_web

# Qas
qas:
> docker compose $(ENV_FILE) $(COMPOSE_QAS) up -d

qas-build:
> docker compose $(ENV_FILE) $(COMPOSE_QAS) up -d --build

qas-down:
> docker compose $(ENV_FILE) $(COMPOSE_QAS) down

qas-logs:
> docker compose $(ENV_FILE) $(COMPOSE_QAS) logs -f

# Prod
prod:
> docker compose $(ENV_FILE) $(COMPOSE_PROD) up -d

prod-down:
> docker compose $(ENV_FILE) $(COMPOSE_PROD) down

# Utils
qas-ps:
> docker compose $(ENV_FILE) $(COMPOSE_QAS) ps

qas-clean:
> docker compose $(ENV_FILE) $(COMPOSE_QAS) down -v  # removes volumes too, careful