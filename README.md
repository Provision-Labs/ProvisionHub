<p align="center">
  <img src="./docs/img/provision-labs-logo-wo-bg.png" width="160"/>
</p>

<h1 align="center">ProvisionHub</h1>

<p align="center">
Self-service Platform Provisioning • Git-native • Async • Extensible
</p>

---

## 🚀 Overview

**ProvisionHub** is an open-source platform provisioning system designed to scaffold and orchestrate application and infrastructure components using Git-native workflows and asynchronous execution.

It enables developers and platform teams to:

- Create systems and services quickly
- Generate infrastructure-ready repositories
- Automate provisioning workflows
- Integrate with GitOps pipelines
- Build internal developer platforms

ProvisionHub is designed to be **adaptable**, **extensible**, and **environment-agnostic**.

---

## ✨ Core Capabilities

- System scaffolding (WIP)
- Component generation (backend, frontend, database, async, etc.) (WIP)
- Git-native repository creation (WIP)
- Asynchronous provisioning (queue + workers) (WIP)
- Execution tracking and logs (WIP)
- Policy-aware workflows (approval / automation) (WIP)
- Optional GitOps integration (ArgoCD) (WIP)
- Template-driven architecture (WIP)

---

## 🧱 Architecture

ProvisionHub follows a **Control Plane + Worker Plane** architecture.

### Control Plane (Go)
Responsible for:
- Authentication (OIDC / Keycloak)
- System & Component management
- Git provider integration (GitLab)
- Publishing provisioning jobs
- Tracking provisioning runs

### Worker Plane (Go)
Responsible for:
- Executing provisioning jobs
- Rendering templates
- Creating repositories
- Committing & pushing changes
- Updating run status

### Core Infrastructure
- PostgreSQL → state & audit
- RabbitMQ → async job processing
- Git provider → source of truth
- Optional GitOps → deployment automation

---

## 🧠 Core Concepts

### System
Logical project container that groups multiple components.

### Component
A deployable unit generated from templates (backend, frontend, database, queue, etc.).

### Blueprint
Configuration that defines how a system or component should be generated.

### Provisioning Run
Tracks execution of a provisioning request and its steps.

---

## 📦 Repository Structure

```
apps/
  control-plane/        # Go API
  worker/               # Go async worker
  web/                  # Next.js frontend (optional)

catalog/
  modules/              # Component definitions

templates/
  helm/
  kustomize/

docs/
deployments/
```

---

## ⚙️ Getting Started (Local)

### Requirements
- Go 1.22+
- Docker & Docker Compose
- Next.js

### Start infrastructure

```bash
docker compose -f deployments/docker-compose.yaml up -d
```

Services:
- PostgreSQL
- RabbitMQ
- Keycloak (for auth)

### Run control plane

```bash
cd apps/control-plane
go run ./cmd/server
```

### Run worker

```bash
cd apps/worker
go run ./cmd/worker
```

---

## 🔐 Authentication

ProvisionHub uses **OIDC (Keycloak)** for authentication.

Flow:
1. Login via browser
2. Obtain JWT access token
3. Call API using `Authorization: Bearer <token>`

---

## 🔄 Provisioning Flow

1. User creates a System or Component
2. Control plane validates blueprint
3. Job published to RabbitMQ
4. Worker executes provisioning steps
5. Repository generated
6. Status & logs updated
7. Optional: GitOps deployment

---

## 🧩 Extensibility

ProvisionHub is designed to be modular.

You can extend:
- New component types
- New templates
- New Git providers
- New provisioning steps
- GitOps integrations
- Cloud provisioning plugins

---

## 🧪 Roadmap

### v0.1
- Auth (Keycloak)
- System creation
- Component scaffolding
- Git repo generation
- Async provisioning
- Run tracking

### v0.2
- Approval workflows
- GitOps compatibility
- Retry + DLQ
- Policy engine

### v1.0
- Module versioning
- Multi-environment support
- CLI
- Kubernetes Operator
- Plugin system

---

## 🤝 Contributing

We welcome contributions.

Steps:
1. Fork repository
2. Create branch
3. Submit PR

See `CONTRIBUTING.md` for details.

---

## 📜 License

ProvisionHub is licensed under the Apache License 2.0.

---

## 🧭 Vision

ProvisionHub aims to enable **adaptive platform engineering** — where infrastructure, automation, and developer experience converge into programmable, self-service systems.

---

<p align="center">
Built with adaptability in mind 🦎
</p>
