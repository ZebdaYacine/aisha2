# PRD-22 · Staging, Production Operations & Rollback

| | |
| --- | --- |
| **Epic** | 22 — Staging, Production Operations & Rollback |
| **Stories** | Deployment, backups, monitoring, health, and rollback workflows |
| **Priority** | Critical |
| **Status** | 🟡 Partial |
| **Surfaces** | Docker Compose, Jenkins, Nginx, staging, production, worker |
| **API** | `GET /health/live`, `GET /health/ready`, deployment smoke endpoints |

## 1. Summary

Operate AISHA safely across local, CI, staging, and production with private infrastructure, backups, health checks, monitoring, immutable images, and a tested rollback path.

## 2. Problem

A deployment can corrupt data, expose services, lose object files, or leave workers unhealthy if migration, backup, secret, readiness, and rollback procedures are informal.

## 3. Goals

- Provide repeatable staging and production deployment steps.
- Back up PostgreSQL and MinIO before risky changes.
- Verify readiness and smoke journeys after deployment.
- Retain previous images and a recoverable rollback plan.

### Non-goals

- Zero-downtime multi-region failover or fully automated disaster recovery in the MVP.

## 4. Users

- Operator deploying a release.
- Release manager approving production.
- On-call engineer diagnosing health or failed workers.
- Customer relying on a stable service.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Separate local, CI, staging, and production credentials and configuration. | Deployment requirements |
| R2 | Keep database, Redis, and MinIO private behind the reverse proxy/network boundary. | Security requirements |
| R3 | Run backup, migration pre-check, deployment, readiness, smoke, and rollback steps. | Deployment requirements |
| R4 | Preserve PostgreSQL and MinIO data through releases and test restoration. | Deployment requirements |
| R5 | Monitor API errors/latency, database/Redis/MinIO health, worker failures, reservations, webhooks, disk, and restarts. | Monitoring requirements |

## 6. Flow

```text
Build immutable images → backup → migration pre-check → migrate → start API/worker/frontend → readiness → smoke → approve or rollback
```

## 7. Technical notes

- Never use `docker compose down -v` in production.
- Use non-root images, HTTPS, HSTS after TLS verification, and secret injection.
- Migrations should remain backward-compatible where possible.

## 8. Success metrics

- Deployment success and rollback time.
- Backup/restore drill success.
- Post-deploy smoke and readiness pass rate.

## 9. Risks & open questions

- Hosting, TLS, secret manager, backup retention, and alert ownership need confirmation.
- Destructive migration rollback requires an explicit approved plan.

## Source traceability

`docs/deployment.md`, `docs/security.md`, `docs/architecture.md`.
