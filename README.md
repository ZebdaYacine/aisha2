# AISHA

Milestone 1 establishes the platform foundation for AISHA without implementing
authentication, catalogue, payments, orders, inventory, or other business
features.

## Included

- Next.js 16, React 19, strict TypeScript, Tailwind CSS, and four-locale routing.
- Go 1.25 and Fiber v3 API foundation.
- Google Wire composition root.
- PostgreSQL, Redis, and MinIO clients with API readiness checks.
- SQLBoiler configuration with no generated business models.
- Docker Compose with persistent PostgreSQL and MinIO volumes.
- MinIO bucket initialization and an Nginx development entry point.
- Frontend and API health endpoints.
- Makefile and Jenkins validation foundations.

## Local development

Requirements: Node 22+, Go 1.25+, Docker with Compose.

```bash
cp .env.example .env
make install
make check
make up
```

To rebuild and redeploy the latest API and frontend code while preserving the
PostgreSQL and MinIO containers and volumes:

```bash
./script.sh
```

Use `./script.sh --prod` to apply the production Compose overrides. The script
does not remove volumes and applies pending migrations through a dedicated
one-shot migration container before loading the idempotent development seed
file and starting the API. In local development, loading the seed intentionally
resets application data first so the demo environment stays deterministic; it
does not remove Docker volumes or affect production.

The Compose files live under `infrastructure/compose`; `make up` is the
equivalent command when running from the repository root.

Endpoints:

- Nginx entry point: `http://localhost:8089/en`
- Frontend: `http://localhost:3033/en`
- Frontend health: `http://localhost:3033/api/health`
- API liveness: `http://localhost:8088/health/live`
- API readiness: `http://localhost:8088/health/ready`
- MinIO console: `http://localhost:9001`

Routine shutdown uses `make down`. Do not add `-v`; PostgreSQL and
MinIO data must be preserved.

## Local demo accounts

The development seed creates six Algerian demo accounts. They all use the
development-only password `Yassine1996@Got`:

| Account | Email | Role |
| --- | --- | --- |
| Saad | `saad.admin@example.test` | Administrator |
| Kader | `kader.agent@example.test` | Warehouse agent |
| Yassine | `yassine.warehouse@example.test` | Warehouse agent |
| Lyna | `lyna.customer@example.test` | Customer |
| Oussama | `oussama.artisan@example.test` | Artisan |
| Youcef | `youcef.moderator@example.test` | Moderator |

The application currently models agent/warehouse access with the single
`warehouse_agent` role. These credentials are for local development only and
must never be reused in production.

## Database

Database migrations remain a dedicated deployment step and must not run during
API startup. The development seed is safe to replay because it resets the
development tables and recreates the deterministic fixture set.
# aisha2
