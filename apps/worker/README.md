# ProvisionHub Worker

The **Worker** is the asynchronous job executor of ProvisionHub, built in Go. It processes provisioning requests, renders templates, and creates Git repositories.

---

## 🎯 Purpose

The Worker is responsible for the actual execution of provisioning tasks:

- **Job Consumption** - Listens to RabbitMQ for provisioning jobs
- **Template Rendering** - Processes templates with user-defined values
- **Repository Creation** - Creates Git repositories via GitLab API
- **Code Generation** - Scaffolds project structure
- **Git Operations** - Commits and pushes changes
- **Status Reporting** - Updates provisioning run status
- **Error Handling** - Manages failures and retries

---

## 🏗️ Architecture

### Components

- **Message Consumer** - RabbitMQ subscriber
- **Template Engine** - Renders templates (Go templates / Jinja2)
- **Git Client** - GitLab API integration
- **File System Handler** - Local workspace management
- **Status Reporter** - Callback to Control Plane API

### Technology Stack

- **Language**: Go 1.22+
- **Message Queue**: RabbitMQ
- **Template Engine**: Go `text/template` or external
- **Git Provider**: GitLab
- **File System**: Local temporary workspaces

---

## 🔄 Execution Flow

1. **Consume Job** - Worker picks up message from RabbitMQ queue
2. **Validate Blueprint** - Checks provisioning configuration
3. **Prepare Workspace** - Creates temporary directory
4. **Render Templates** - Processes all template files with provided values
5. **Create Repository** - Uses GitLab API to create new repo
6. **Commit & Push** - Commits generated files and pushes to GitLab
7. **Update Status** - Reports success/failure to Control Plane
8. **Cleanup** - Removes temporary workspace

---

## 🚀 Getting Started

### Prerequisites

- Go 1.22 or higher
- RabbitMQ running (via Docker Compose)
- GitLab access (token and group)
- Control Plane API running

### Environment Variables

Create a `.env` file or set the following variables:

```bash
# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_QUEUE=provisioning-jobs
RABBITMQ_PREFETCH_COUNT=1
RABBITMQ_RETRY_DELAY=5s

# Control Plane
CONTROL_PLANE_URL=http://localhost:8080
CONTROL_PLANE_API_KEY=worker-api-key

# GitLab
GITLAB_URL=https://gitlab.com
GITLAB_TOKEN=your-gitlab-token
GITLAB_GROUP_ID=your-group-id

# Worker
WORKER_ID=worker-001
WORKER_CONCURRENCY=5
WORKSPACE_DIR=/tmp/provisionhub-workspaces
TEMPLATE_CACHE_DIR=/tmp/provisionhub-templates

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### Installation

```bash
cd apps/worker

# Install dependencies
go mod download

# Run the worker
go run ./cmd/worker
```

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
apps/worker/
├── cmd/
│   └── worker/          # Main worker entrypoint
├── internal/
│   ├── consumer/       # RabbitMQ consumer logic
│   ├── executor/       # Provisioning execution
│   ├── template/       # Template rendering
│   ├── git/            # Git operations
│   ├── reporter/       # Status updates to Control Plane
│   └── workspace/      # Workspace management
├── pkg/                # Shared utilities
├── go.mod
├── go.sum
└── README.md
```

---

## ⚙️ Job Message Format

Workers consume JSON messages from RabbitMQ:

```json
{
  "run_id": "uuid-of-provisioning-run",
  "type": "system" | "component",
  "blueprint": {
    "template": "go-backend",
    "name": "my-service",
    "namespace": "my-team",
    "values": {
      "port": 8080,
      "database": "postgres"
    }
  },
  "git_config": {
    "provider": "gitlab",
    "group_id": "123",
    "visibility": "private"
  }
}
```

---

## 🎨 Template Engine

Workers use Go's `text/template` for rendering:

### Template Variables

Templates have access to:
- `.Name` - Component/system name
- `.Namespace` - Namespace or group
- `.Values` - User-provided values
- `.Meta` - Metadata (timestamp, version, etc.)

### Example Template

```go
// main.go.tmpl
package main

import "fmt"

func main() {
    fmt.Println("{{ .Name }} starting on port {{ .Values.port }}")
}
```

---

## 🔄 Concurrent Processing

Workers can process multiple jobs concurrently:

```go
// Configure via environment
WORKER_CONCURRENCY=5  // Process 5 jobs at once
```

Each job runs in an isolated workspace to prevent conflicts.

---

## 🚨 Error Handling

### Retry Strategy

- **Transient Errors**: Requeue with exponential backoff
- **Template Errors**: Mark as failed, report to Control Plane
- **Git Errors**: Retry up to 3 times
- **Dead Letter Queue**: Failed jobs after max retries

### Logging

All execution steps are logged and sent to Control Plane:

```
[INFO] Starting provisioning for run: abc-123
[INFO] Rendering template: go-backend
[INFO] Creating repository: my-team/my-service
[INFO] Committing files (42 files)
[INFO] Pushing to GitLab
[SUCCESS] Provisioning completed in 12.5s
```

---

## 🧪 Testing

```bash
# Unit tests
go test ./internal/...

# Integration tests (requires RabbitMQ)
go test -tags=integration ./...

# Test template rendering
go test ./internal/template/...

# Mock tests
go test -v ./...
```

---

## 🐳 Docker

Build and run the Worker in a container:

```bash
# Build image
docker build -t provisionhub-worker .

# Run container
docker run \
  --env-file .env \
  provisionhub-worker
```

---

## 📊 Monitoring

### Health Checks

Workers expose health endpoints:
- `GET /health` - Worker health
- `GET /metrics` - Prometheus metrics

### Metrics

Key metrics exposed:
- `jobs_processed_total` - Total jobs processed
- `jobs_failed_total` - Total failed jobs
- `job_duration_seconds` - Processing time histogram
- `queue_depth` - Current RabbitMQ queue depth

---

## ⚡ Scaling

Scale workers horizontally for increased throughput:

```bash
# Docker Compose
docker-compose up --scale worker=5

# Kubernetes
kubectl scale deployment provisionhub-worker --replicas=5
```

Workers coordinate via RabbitMQ's built-in load balancing.

---

## 🔧 Configuration

### Worker Modes

- **Standard Mode**: Process all job types
- **Specialized Mode**: Process specific types only
- **Batch Mode**: Process multiple jobs per execution

### Resource Limits

Configure resource constraints:
```bash
MAX_TEMPLATE_SIZE=10MB
MAX_WORKSPACE_SIZE=500MB
MAX_JOB_DURATION=10m
```

---

## 🤝 Contributing

See the main [CONTRIBUTING.md](../../CONTRIBUTING.md) for contribution guidelines.

---

## 📜 License

Apache License 2.0 - See [LICENSE](../../LICENSE)

---

## 🔗 Related

- [Main Project Documentation](../../README.md)
- [Control Plane Documentation](../control-plane/README.md)
- [Web UI Documentation](../web/README.md)
