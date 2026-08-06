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
file and starting the API.

The Compose files live under `infrastructure/compose`; `make up` is the
equivalent command when running from the repository root.

Endpoints:

- Nginx entry point: `http://localhost:8088/en`
- Frontend: `http://localhost:3000/en`
- Frontend health: `http://localhost:3000/api/health`
- API liveness: `http://localhost:8080/health/live`
- API readiness: `http://localhost:8080/health/ready`
- MinIO console: `http://localhost:9001`

Routine shutdown uses `make down`. Do not add `-v`; PostgreSQL and
MinIO data must be preserved.

## Database

Milestone 1 deliberately creates no business tables. Future milestones add
versioned migrations before generating SQLBoiler models. Database migrations
must remain a dedicated deployment step and must not run during API startup.
# aisha2
