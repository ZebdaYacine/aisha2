# AISHA MVP Deployment and CI/CD Guide

## 1. Environments

Minimum environments:

- Local development.
- CI test.
- Staging.
- Production.

Never use production credentials in development or CI.

## 2. Docker Services

```text
frontend
api
worker
postgres
redis
minio
minio-init
nginx
```

Optional:

- mail testing service in development.
- migration one-shot service.
- automated backup service.

## 3. Docker Compose Requirements

### PostgreSQL

- Persistent named volume.
- Health check using `pg_isready`.
- Database not publicly exposed in production.
- Runtime user with least privilege.
- Migration user handled separately if possible.

### Redis

- Persistent configuration only if required.
- Password in non-local environments.
- Private network.
- Health check using `redis-cli ping`.

### MinIO

- Persistent data volume.
- Private internal endpoint.
- Separate access key and secret.
- Initial buckets:
  - `product-public`
  - `product-private`
  - `artisan-private`
- Public bucket policy only for approved published media, or use reverse-proxied controlled access.

### API

- Multi-stage Go build.
- Non-root runtime user.
- Health endpoints.
- No source code or compiler in final image.
- Environment variables injected at runtime.

### Frontend

- Multi-stage Next.js standalone build.
- Non-root runtime user.
- Health endpoint or HTTP check.
- Runtime/public environment values clearly separated.

### Worker

- Uses the same application code and domain rules.
- Processes outbox and reservation expiry.
- Independent health/readiness where practical.

## 4. Health Endpoints

### API

```text
GET /health/live
GET /health/ready
```

`live` confirms process is alive.

`ready` verifies required dependencies:

- PostgreSQL.
- Redis if required for normal traffic.
- MinIO if uploads are critical.

Do not include secrets in health response.

### Frontend

A lightweight HTTP route must return success after application startup.

## 5. Environment Variables

Example categories:

```text
APP_ENV
APP_NAME
APP_URL
API_PORT
LOG_LEVEL

DATABASE_URL
DATABASE_MAX_OPEN_CONNS
DATABASE_MAX_IDLE_CONNS

REDIS_URL

MINIO_ENDPOINT
MINIO_PUBLIC_ENDPOINT
MINIO_ACCESS_KEY
MINIO_SECRET_KEY
MINIO_USE_SSL
MINIO_PUBLIC_BUCKET
MINIO_PRIVATE_BUCKET

SESSION_SECRET
ACCESS_TOKEN_TTL
REFRESH_TOKEN_TTL

CASBIN_MODEL_PATH
CASBIN_POLICY_SOURCE

PAYMENT_PROVIDER
PAYMENT_API_KEY
PAYMENT_WEBHOOK_SECRET

SHIPPING_PROVIDER
SHIPPING_API_KEY

SMTP_HOST
SMTP_PORT
SMTP_USER
SMTP_PASSWORD
SMTP_FROM

FRONTEND_URL
CORS_ALLOWED_ORIGINS
```

For local development the example allows `localhost:3033` and
`127.0.0.1:3033`. The production Compose override defaults to these local
origins plus the VPS frontend origin `http://167.86.79.16`; set
`CORS_ALLOWED_ORIGINS` explicitly in the Jenkins production environment
credential when the allowed-origin policy should be narrower. Credentials are
enabled for the configured origins only.

Provide `.env.example` with no real values.

## 6. Local Development Commands

Suggested Makefile targets:

```bash
make bootstrap
make up
make down
make logs
make migrate-up
make migrate-down
make sqlboiler
make test
make lint
make build
make smoke
```

Bootstrap:

1. Copy `.env.example` to `.env`.
2. Start PostgreSQL, Redis, and MinIO.
3. Run migrations.
4. Create buckets.
5. Seed the deterministic development roles, demo accounts, categories, and
   application data (local development only; the seed resets development data).
6. Start API, worker, and frontend.

## 7. Jenkins Pipeline

Recommended stages:

```text
Checkout
-> Validate repository
-> Backend dependencies
-> Backend format/vet/test/build
-> Frontend install/lint/type-check/test/build
-> OpenAPI validation
-> Migration validation
-> Docker build
-> Container scan
-> Compose config validation
-> Integration tests
-> Publish images
-> Deploy staging
-> Health checks
-> Smoke tests
-> Manual production approval
-> Production deployment
-> Health checks
-> Rollback on failure
```

### VPS deployment pipeline

The repository `Jenkinsfile` keeps deployment opt-in. A normal job runs the
validation stages only. Start a build with `DEPLOY_VPS=true` to deploy the
validated commit to the configured VPS.

Configure these Jenkins credentials before enabling the deploy parameter:

- `aisha-vps-ssh`: SSH private-key credential for the VPS user. The Jenkins
  agent must already contain the VPS host key in its `known_hosts`; the
  pipeline intentionally uses strict host-key checking.
- `aisha-vps-production-env`: Secret-file credential containing the production
  `.env`. It is copied over SSH with mode `0600` and is never committed or
  printed in the build log.

The deploy stage requires `rsync`, `ssh`, `scp`, Docker Compose, and the
Jenkins SSH Agent plugin on the build agent. It synchronizes source without
deleting the VPS application directory, backs up PostgreSQL, runs migrations
without development seed data, builds images tagged with the commit SHA, and
runs API, frontend, reverse-proxy, and readiness smoke checks. PostgreSQL,
Redis, and MinIO volumes are preserved. The remote deploy script records the
last successful image tag and restores it if the new containers fail health
checks. Database migrations are not automatically rolled back; they must be
backward-compatible before production approval.

The VPS must have Docker Compose, `curl`, and a production `.env` policy with
private PostgreSQL, Redis, and MinIO ports. Do not place the VPS password in
the Jenkinsfile; use an SSH key credential and rotate any temporary password
after key-based access is configured.

## 8. Example Jenkins Stage Requirements

### Repository validation

- Required files exist.
- No committed `.env`.
- No obvious secrets.
- Docker Compose parses.

### Backend

```bash
go mod download
test -z "$(gofmt -l .)"
go vet ./...
go test -race ./...
go build ./...
```

### Frontend

```bash
npm ci
npm run lint
npm run type-check
npm test -- --runInBand
npm run build
```

### Database

- Start temporary PostgreSQL.
- Run all migrations up.
- Run integration tests.
- Run down migrations where safe.
- Run up again.

### Docker

```bash
make compose-config
make build
```

## 9. Deployment Strategy

For MVP, use a safe recreate strategy with persistent volumes:

1. Pull new images.
2. Back up PostgreSQL.
3. Run migration pre-check.
4. Run migrations.
5. Start new API and worker.
6. Start frontend.
7. Verify readiness.
8. Run smoke tests.
9. Keep previous image tags for rollback.

Do not execute:

```bash
docker compose -f infrastructure/compose/docker-compose.yml down -v
```

in production.

## 10. Image Tagging

Use immutable tags:

```text
aisha-api:<git-sha>
aisha-worker:<git-sha>
aisha-frontend:<git-sha>
```

Also maintain controlled environment aliases such as `staging` only if deployment tooling requires them.

Do not deploy only `latest` without retaining the exact commit tag.

## 11. Database Backup

Before production migration:

- Create timestamped PostgreSQL backup.
- Encrypt or protect backup storage.
- Verify backup command succeeds.
- Define retention.
- Perform restore drills.

Example operational workflow:

```bash
pg_dump --format=custom --file=backup.dump "$DATABASE_URL"
```

Credentials must not appear in command logs.

## 12. MinIO Backup

- Preserve MinIO volume.
- Replicate or copy objects to backup storage.
- Back up bucket policies and lifecycle configuration.
- Test restoration of private and public media metadata.

Database backup alone is not sufficient because object files are external.

## 13. Rollback

Rollback requires:

- Previous image tags.
- Migration compatibility analysis.
- Database backup.
- Deployment script.
- Smoke test.

Preferred migrations are backward-compatible so application rollback does not require immediate database rollback.

If a migration is destructive, deployment must stop until an explicit rollback and backup plan exists.

## 14. Smoke Tests

After deployment:

1. Frontend home returns `200`.
2. API live and ready return success.
3. Public categories load.
4. Public products load.
5. Login works with staging test user.
6. Protected endpoint rejects unauthenticated access.
7. MinIO upload/read test succeeds in staging.
8. Redis rate limit or cache test succeeds.
9. Create and cancel development checkout.
10. Worker processes a test outbox event.

Production smoke tests must avoid real payment and destructive warehouse actions.

## 15. Monitoring

Minimum observability:

- Structured application logs.
- Reverse proxy access logs.
- API latency and error count.
- PostgreSQL health.
- Redis health.
- MinIO health and storage usage.
- Worker pending/failed outbox count.
- Active expired reservations count.
- Payment webhook failures.
- Disk space.
- Container restart count.

## 16. Production Security

- HTTPS only.
- HSTS after TLS is verified.
- Database, Redis, and MinIO are private.
- Firewall permits only required ports.
- Real secrets come from Jenkins credentials or secret manager.
- Containers run non-root where practical.
- No development debug mode.
- CORS allows only production frontend.
- Default admin password is changed before release.

## 17. Deployment Success Criteria

A deployment is successful only when:

- Migrations succeeded.
- Containers are healthy.
- API readiness passes.
- Frontend responds.
- Worker runs.
- Smoke tests pass.
- No severe errors appear in logs.
- Previous release remains available for rollback.
