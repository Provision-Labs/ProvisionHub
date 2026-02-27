# ProvisionHub Control Plane

The **Control Plane** is the core API service of ProvisionHub, built in Go. It handles authentication, system management, and orchestrates provisioning workflows.

---

## 🎯 Purpose

The Control Plane serves as the central coordinator for the ProvisionHub platform, responsible for:

- **Authentication & Authorization** - OIDC integration with Keycloak
- **System & Component Management** - CRUD operations for provisioning entities
- **Git Provider Integration** - GitLab repository management
- **Job Publishing** - Sends provisioning jobs to RabbitMQ
- **Execution Tracking** - Monitors and logs provisioning runs
- **API Gateway** - RESTful API for all platform operations

---

## 🏗️ Architecture

### Components

- **HTTP Server** - RESTful API endpoints
- **Database Layer** - PostgreSQL for state persistence
- **Message Publisher** - RabbitMQ integration for async jobs
- **Git Client** - GitLab API integration
- **Auth Middleware** - JWT validation and OIDC flow

### Technology Stack

- **Language**: Go 1.22+
- **HTTP Framework**: TBD (net/http)
- **Database**: PostgreSQL
- **Message Queue**: RabbitMQ
- **Authentication**: OIDC (Keycloak)
- **Git Provider**: GitLab

---

## 🔌 API Endpoints

### Authentication

- `POST /auth/login` - Initiate OIDC login
- `POST /auth/callback` - OIDC callback handler
- `POST /auth/logout` - End user session

### Systems

- `GET /api/v1/systems` - List systems
- `POST /api/v1/systems` - Create system
- `GET /api/v1/systems/:id` - Get system details
- `PUT /api/v1/systems/:id` - Update system
- `DELETE /api/v1/systems/:id` - Delete system

### Components

- `GET /api/v1/systems/:id/components` - List components
- `POST /api/v1/systems/:id/components` - Create component
- `GET /api/v1/components/:id` - Get component details
- `PUT /api/v1/components/:id` - Update component
- `DELETE /api/v1/components/:id` - Delete component

### Provisioning Runs

- `GET /api/v1/runs` - List provisioning runs
- `GET /api/v1/runs/:id` - Get run details
- `GET /api/v1/runs/:id/logs` - Stream run logs
- `POST /api/v1/runs/:id/retry` - Retry failed run

---

## 🚀 Getting Started

### Prerequisites

- Go 1.22 or higher
- PostgreSQL running (via Docker Compose)
- RabbitMQ running (via Docker Compose)
- Keycloak running (via Docker Compose)

### Environment Variables

Create a `.env` file or set the following variables:

```bash
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=provisionhub
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSL_MODE=disable

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# Keycloak
OIDC_ISSUER_URL=http://localhost:8081/realms/provisionhub
OIDC_CLIENT_ID=provisionhub-api
OIDC_CLIENT_SECRET=your-secret-here

# GitLab
GITLAB_URL=https://gitlab.com
GITLAB_TOKEN=your-gitlab-token
GITLAB_GROUP_ID=your-group-id
```

### Installation

```bash
cd apps/control-plane

# Install dependencies
go mod download

# Run database migrations
go run ./cmd/migrate up

# Start the server
go run ./cmd/server
```

The API will be available at `http://localhost:8080`

### Development

```bash
# Run tests
go test ./...

# Run with hot reload (requires air)
air

# Format code
go fmt ./...

# Lint
golangci-lint run
```

---

## 📁 Project Structure

```
apps/control-plane/
├── cmd/
│   ├── server/          # Main API server entrypoint
│   └── migrate/         # Database migration tool
├── internal/
│   ├── api/            # HTTP handlers and routes
│   ├── auth/           # Authentication logic
│   ├── domain/         # Business logic and entities
│   ├── git/            # Git provider integration
│   ├── queue/          # RabbitMQ publisher
│   └── store/          # Database repositories
├── pkg/                # Shared utilities
├── migrations/         # SQL migrations
├── go.mod
├── go.sum
└── README.md
```

---

## 🔐 Authentication Flow

1. User initiates login via `/auth/login`
2. Redirect to Keycloak login page
3. User authenticates with Keycloak
4. Keycloak redirects to `/auth/callback` with authorization code
5. Control Plane exchanges code for access token
6. JWT token returned to client
7. Client includes token in `Authorization: Bearer <token>` header

---

## 🔄 Provisioning Workflow

1. User submits provisioning request (System or Component)
2. Control Plane validates blueprint schema
3. Creates a Provisioning Run record in database
4. Publishes job message to RabbitMQ queue
5. Returns Run ID to user
6. Worker picks up job and executes
7. Worker updates run status via callback API
8. User can query run status and logs

---

## 🧪 Testing

```bash
# Unit tests
go test ./internal/...

# Integration tests (requires infrastructure)
go test -tags=integration ./...

# Test coverage
go test -cover ./...
```

---

## 🐳 Docker

Build and run the Control Plane in a container:

```bash
# Build image
docker build -t provisionhub-control-plane .

# Run container
docker run -p 8080:8080 \
  --env-file .env \
  provisionhub-control-plane
```

---

## 📊 Monitoring & Observability

- **Health Check**: `GET /health`
- **Readiness**: `GET /ready`
- **Metrics**: `GET /metrics` (Prometheus format)

---

## 🤝 Contributing

See the main [CONTRIBUTING.md](../../CONTRIBUTING.md) for contribution guidelines.

---

## 📜 License

Apache License 2.0 - See [LICENSE](../../LICENSE)

---

## 🔗 Related

- [Main Project Documentation](../../README.md)
- [Worker Documentation](../worker/README.md)
- [Web UI Documentation](../web/README.md)
