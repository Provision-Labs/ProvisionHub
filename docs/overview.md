# ProvisionHub — Overview

## What is ProvisionHub

ProvisionHub is a lightweight, extensible platform for provisioning and managing developer systems through a plugin-driven architecture. It separates orchestration from execution, so platform teams can define workflows in the Control Plane while distributed Workers execute heavy operations.

The platform is designed to be:

- Plugin-first
- Language-agnostic (via gRPC contracts)
- Distributed by default
- Extensible without modifying the core
- UI-decoupled (reference UI can be replaced independently)
- Suitable for internal developer platforms
- Lightweight compared to monolithic portals

ProvisionHub standardizes how repositories, services, infrastructure, and integrations are created and managed.

---

## Problem Statement

Organizations typically struggle with:

- Manual service provisioning
- Inconsistent repository setup
- Multiple CI/CD patterns
- Repeated infrastructure boilerplate
- Lack of standardized templates
- Hardcoded platform logic
- Poor extensibility for new providers

Existing solutions often:

- Are too heavy
- Require frontend-first plugin systems
- Are difficult to extend
- Mix orchestration with execution
- Do not support distributed workers well

ProvisionHub addresses these issues with a minimal core plus a robust plugin system.

---

## Goals

ProvisionHub aims to:

- Provide a plugin-based provisioning platform
- Separate orchestration from execution
- Support multiple providers (GitHub, GitLab, cloud vendors, internal systems)
- Allow distributed execution via workers
- Enable dynamic plugin installation
- Support versioned plugins
- Provide CLI-driven plugin management
- Keep the core minimal and stable
- Allow multi-language plugin development

---

## Non-Goals

ProvisionHub is not:

- A full CI/CD system
- A Kubernetes replacement
- A GitOps controller
- A full developer portal UI
- A workflow engine replacement
- A secrets manager
- A monolithic platform

These capabilities may exist as plugins, not in the core.

ProvisionHub may include a standard reference UI in the monorepo, but this UI is optional and fully decoupled from core runtime behavior.

---

## Core Building Blocks

### Control Plane

Responsible for orchestration, plugin lifecycle management, job scheduling, API exposure, and system state. The Control Plane does not execute heavy provisioning operations.

### Workers

Stateless, horizontally scalable executors that consume jobs from the queue and run provider-specific plugins.

### Plugins

Two plugin classes:

- Control Plane plugins: orchestration logic, auth providers, templates, high-level workflow decisions
- Worker plugins: provider integrations for Git, deploy, build, and infrastructure execution

### Jobs

A Job is the unit of work sent from Control Plane to Workers. Jobs include type, provider, payload, metadata, and optional secrets.

### Registry

ProvisionHub supports dynamic plugin installation through public and private registries.

---

## High-Level Architecture

See [docs/architecture/architecture.md](docs/architecture/architecture.md) for system boundaries, macro decisions, and end-to-end flow.

For internal teams that need to rebuild or replace the UI, see [docs/architecture/ui-rebuild.md](docs/architecture/ui-rebuild.md).

---

## Example End-to-End Flow

1. User requests a new service.
2. Control Plane plugin validates intent and composes a provisioning plan.
3. Control Plane emits jobs (for example, create repo, setup CI, deploy baseline infra).
4. Workers consume jobs and execute provider actions via Worker plugins.
5. Job results are reported back to Control Plane and exposed through APIs.

---

## Target Use Cases

- Service scaffolding
- Repository creation
- CI/CD bootstrapping
- Infrastructure provisioning
- Template-based services
- Internal developer platforms
- Multi-provider deployments

---

## Architecture Decision Records

Major decisions and trade-offs are documented in [docs/adr/0001-go-monorepo.md](docs/adr/0001-go-monorepo.md) and [docs/adr/0002-grpc-plugin-contract.md](docs/adr/0002-grpc-plugin-contract.md).

---

## Summary

ProvisionHub is a plugin-driven provisioning platform that keeps orchestration and execution separate, scales through workers, and evolves safely through explicit architectural decisions.
