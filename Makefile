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


# Generic plugin lifecycle test for all enabled plugins in registry
test-plugins:
> @set -e; \
> echo "[test-plugins] Starting dev infra (Keycloak + CP database)"; \
> $(MAKE) dev-infra-up; \
> \
> set -a; \
> . deployments/.env.local; \
> set +a; \
> \
> export ENV=$${CP_ENV:-development}; \
> export DB_HOST=localhost; \
> export DB_PORT=5433; \
> export DB_NAME=$${POSTGRES_CP_DB}; \
> export DB_USERNAME=$${POSTGRES_CP_USER}; \
> export DB_PASSWORD=$${POSTGRES_CP_PASSWORD}; \
> export DB_SSL_MODE=disable; \
> export AUTH_ISSUER=$${AUTH_ISSUER:-http://localhost:8280/realms/provisionhub}; \
> export AUTH_CLIENT_ID=$${AUTH_CLIENT_ID:-$${CP_CLIENT_ID}}; \
> export AUTH_CLIENT_SECRET=$${AUTH_CLIENT_SECRET:-$${CP_CLIENT_SECRET}}; \
> export AUTH_REDIRECT_URL=$${AUTH_REDIRECT_URL:-$${CP_REDIRECT_URL}}; \
> export AUTH_SCOPES=$${AUTH_SCOPES:-$${CP_SCOPES}}; \
> export AUTH_LOGOUT_REDIRECT_URL=$${AUTH_LOGOUT_REDIRECT_URL:-$${LOGOUT_REDIRECT}}; \
> \
> if [ -z "$$DB_NAME" ] || [ -z "$$DB_USERNAME" ] || [ -z "$$DB_PASSWORD" ]; then \
>   echo "[test-plugins] Missing DB vars in deployments/.env.local"; \
>   exit 1; \
> fi; \
> \
> if [ -z "$$AUTH_CLIENT_ID" ] || [ -z "$$AUTH_CLIENT_SECRET" ] || [ -z "$$AUTH_REDIRECT_URL" ] || [ -z "$$AUTH_LOGOUT_REDIRECT_URL" ]; then \
>   echo "[test-plugins] Missing auth vars in deployments/.env.local"; \
>   exit 1; \
> fi; \
> \
> cd apps/control-plane; \
> go run ./cmd/plugin-lifecycle -registry $${PLUGINS_REGISTRY_PATH:-plugins.json}