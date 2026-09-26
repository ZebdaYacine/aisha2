# PRD-02 · Platform & Development Foundation

| | |
| --- | --- |
| **Epic** | 2 — Platform & Development Foundation |
| **Stories** | Local development, CI, deployment baseline |
| **Priority** | High |
| **Status** | ✅ Implemented |
| **Surfaces** | `infrastructure/`, `apps/api`, `apps/web`, `Makefile`, `Jenkinsfile` |
| **API** | `GET /health/live`, `GET /health/ready` |

## 1. Summary

Provide a repeatable runtime for the Go API, Next.js storefront, PostgreSQL, Redis, MinIO, Nginx, migrations, workers, and quality checks.

## 2. Problem

Without a reproducible environment, local work, CI validation, and deployment can drift. Service failures can remain hidden and production data can be endangered by implicit migrations or disposable volumes.

## 3. Goals

- Bootstrap the application consistently across local, CI, staging, and production.
- Expose liveness and readiness without leaking secrets.
- Keep migrations explicit, reversible, and operationally visible.

### Non-goals

- Multi-region deployment, autoscaling, or a fully managed cloud platform in the MVP.

## 4. Users

- Developer running the stack locally.
- CI pipeline validating a change.
- Operator deploying or rolling back a release.
- Support engineer checking service health.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Run frontend, API, worker, PostgreSQL, Redis, MinIO, and reverse proxy with health checks. | US-OPS-001 |
| R2 | Persist required data and object-storage volumes. | US-OPS-002 |
| R3 | Provide explicit migration up/down, seed, test, build, and smoke commands. | US-OPS-003 |
| R4 | Build non-root multi-stage images and inject secrets at runtime. | US-OPS-004 |
| R5 | Propagate request IDs and structured logs through the API. | US-OPS-005 |

## 6. Flow

```text
Checkout → bootstrap services → migrate → seed → test → build → health check → smoke test
```

## 7. Technical notes

- Production startup must not run migrations invisibly.
- PostgreSQL, Redis, and MinIO stay on a private network.
- Immutable image tags and persistent volumes are required for rollback.

## 8. Success metrics

- Clean checkout reaches a healthy local stack with documented commands.
- CI executes backend, frontend, OpenAPI, migration, and container checks.
- Health checks detect dependency failures before traffic is accepted.

## 9. Risks & open questions

- Worker deployment and staging/production provisioning still require operational evidence.
- Backup retention, restore cadence, and alert thresholds need an explicit owner.

## Source traceability

`docs/deployment.md`, `docs/architecture.md`, `docs/backend.md`.
