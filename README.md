# ProvisionHub

> Lightweight, plugin-driven provisioning platform for internal developer platforms

It separates orchestration from execution so platform teams can define workflows in the Control Plane while distributed Workers perform the heavy work.

## What this repository contains

This repository is the documentation and architecture home for ProvisionHub.

- The product vision and problem statement
- The macro architecture and system boundaries
- The architectural decisions captured as ADRs
- Supporting docs for jobs, registry, and plugin development

## Why ProvisionHub exists

ProvisionHub is designed to solve common platform problems:

- Manual service provisioning
- Inconsistent repository setup
- Repeated CI/CD boilerplate
- Hardcoded platform logic
- Poor extensibility for new providers

It aims to stay small at the core and move most implementation detail into plugins.

## Core principles

- Plugin-first
- Distributed by default
- Provider-agnostic
- Language-agnostic through gRPC contracts
- Minimal core, extensible edge

## How it works

```text
UI / CLI → Control Plane → Queue → Workers → Provider Plugins
```

The Control Plane orchestrates workflows, the queue decouples execution, and Workers perform provider-specific actions through plugins.

## Start here

1. Read the product overview: [docs/overview.md](docs/overview.md)
2. Review the macro architecture: [docs/architecture/architecture.md](docs/architecture/architecture.md)
3. Explore the decision log: [docs/adr](docs/adr)

## Documentation map

- [Overview](docs/overview.md) — product vision, problem statement, goals, and core concepts
- [Architecture](docs/architecture/architecture.md) — macro view, system boundaries, and end-to-end flow
- [ADR-0001 — Go monorepo](docs/adr/0001-go-monorepo.md) — read this if you are contributing to the core runtime or shared packages
- [ADR-0002 — gRPC plugin contract](docs/adr/0002-grpc-plugin-contract.md) — read this if you are changing plugin boundaries or contracts
- [UI rebuild guide](docs/architecture/ui-rebuild.md) — read this if you need to replace or reimplement the reference UI

## Status / where we are

Documentation-first.
System boundaries are being defined before implementation starts.
No production code yet.

## Contributing

This repository is documentation-first.

- Add new architecture docs under [docs/architecture](docs/architecture)
- Add new ADRs under [docs/adr](docs/adr) using the next available number
- Use [docs/adr/template.md](docs/adr/template.md) as a starting point for new decisions
- Keep docs small, specific, and linked from this README when they are entry points
