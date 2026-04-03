# UI Rebuild Guide (Internal)

## Objective

This guide documents how to rebuild or replace the standard ProvisionHub UI internally without changing Control Plane or Worker runtime components.

The UI is a replaceable reference component. Core system behavior must remain stable and accessible through API/CLI contracts.

---

## Non-Negotiable Rules

1. API-first boundary

- UI must consume Control Plane APIs only.
- UI must not embed orchestration logic that belongs to backend plugins.

2. No runtime coupling

- UI build/release cycle is independent from Control Plane and Workers.
- Backend deployments must not require frontend redeploys.

3. Contract compatibility

- Rebuilt UI must follow versioned API contracts.
- Breaking API changes require ADR and migration plan.

---

## Required Inputs Before Rebuild

- Current architecture references in docs/architecture/architecture.md
- Active ADRs in docs/adr/
- API contract inventory (REST and/or gRPC gateway endpoints)
- UX requirements and access-control model

---

## Rebuild Process

1. Map user journeys to backend capabilities

- List each screen/workflow and corresponding API calls.
- Identify required read/write permissions for each action.

2. Generate API compatibility matrix

- For each endpoint: method, payload, response schema, errors, auth scope.
- Mark required vs optional fields.

3. Implement frontend shell

- Authentication/session handling
- Navigation and route guards
- Global error handling and retries for idempotent reads

4. Implement core modules

- Provision request form
- Workflow/job status views
- Audit and history views (if exposed)

5. Add observability hooks

- Structured client logs
- Request correlation IDs passed to backend
- Client-side metrics for latency and error rate

6. Validate against contracts

- Run contract tests against staging APIs
- Verify authorization boundaries and forbidden operations
- Verify behavior under partial failure/timeouts

7. Rollout

- Release behind feature flag if replacing existing UI
- Monitor backend error rates and user flows
- Prepare rollback plan to previous UI build

---

## Minimal Validation Checklist

- UI can create provisioning request end-to-end
- UI can read workflow/job status end-to-end
- UI handles 4xx/5xx errors with actionable messages
- UI does not require backend code changes
- UI can be deployed independently

---

## Ownership Model

- Frontend team owns UI implementation and release cadence.
- Platform/backend team owns contracts, orchestration behavior, and execution guarantees.
- Changes that cross boundary (contract shape/semantics) require ADR review.

---

## When to Create an ADR

Create an ADR if rebuild requires:

- New API semantics (not just visual change)
- Breaking contract changes
- Changes in authentication or authorization model
- New dependency between UI lifecycle and backend runtime

---

## References

- docs/overview.md
- docs/architecture/architecture.md
- docs/adr/0001-go-monorepo.md
- docs/adr/0002-grpc-plugin-contract.md
