#!/usr/bin/env bash

set -Eeuo pipefail

ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"

COMPOSE_FILE="$ROOT_DIR/infrastructure/compose/docker-compose.yml"
COMPOSE_PROD_FILE="$ROOT_DIR/infrastructure/compose/docker-compose.prod.yml"
COMPOSE_PROJECT_DIR="$ROOT_DIR/infrastructure/compose"
ENV_FILE="$ROOT_DIR/.env"
SEED_FILE="$ROOT_DIR/apps/api/db/seeds/development.sql"
COMPOSE=()

PRODUCTION=false

log() {
  printf '\n\033[1;34m==>\033[0m %s\n' "$1"
}

error() {
  printf '\n\033[1;31mError:\033[0m %s\n' "$1" >&2
}

on_error() {
  local exit_code=$?

  error "Deployment failed at line ${BASH_LINENO[0]}."

  if command -v docker >/dev/null 2>&1 && ((${#COMPOSE[@]} > 0)); then
    echo
    echo "Current container status:"
    "${COMPOSE[@]}" ps 2>/dev/null || true
    echo
    echo "Recent API logs:"
    "${COMPOSE[@]}" logs --tail=80 api 2>/dev/null || true
  fi

  exit "$exit_code"
}

trap on_error ERR

usage() {
  cat <<'EOF'
Usage:
  ./script.sh
  ./script.sh --prod
  ./script.sh --help

Description:
  Rebuild and redeploy the AISHA Docker Compose project.

  The script recreates the containers while preserving named volumes:

    - postgres_data
    - redis_data
    - minio_data
    - go_modules

Options:
  --prod    Apply docker-compose.prod.yml overrides.
  --help    Display this help message.
EOF
}

case "${1:-}" in
  "")
    ;;
  --prod)
    PRODUCTION=true
    ;;
  --help|-h)
    usage
    exit 0
    ;;
  *)
    error "Unknown option: ${1}"
    usage >&2
    exit 2
    ;;
esac

if [[ ! -f "$COMPOSE_FILE" ]]; then
  error "Compose file not found: $COMPOSE_FILE"
  exit 1
fi

if [[ "$PRODUCTION" == true && ! -f "$COMPOSE_PROD_FILE" ]]; then
  error "Production Compose file not found: $COMPOSE_PROD_FILE"
  exit 1
fi

configured_environment="$(awk -F= '
  $1 == "APP_ENV" {
    value = $0
    sub(/^[^=]*=/, "", value)
    gsub(/^[[:space:]]+|[[:space:]]+$/, "", value)
    if (value ~ /^".*"$/) {
      sub(/^"/, "", value)
      sub(/"$/, "", value)
    }
    print tolower(value)
    exit
  }
' "$ENV_FILE" 2>/dev/null || true)"
if [[ "$PRODUCTION" == false && ("$configured_environment" == "production" || "$configured_environment" == "prod") ]]; then
  PRODUCTION=true
  log "APP_ENV=production detected; applying the production Compose override"
fi

read_root_env_value() {
  local key="$1"
  awk -F= -v key="$key" '
    $1 == key {
      value = $0
      sub(/^[^=]*=/, "", value)
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", value)
      if (value ~ /^".*"$/) {
        sub(/^"/, "", value)
        sub(/"$/, "", value)
      }
      print value
      exit
    }
  ' "$ENV_FILE" 2>/dev/null || true
}

if [[ "$PRODUCTION" == true ]]; then
  if [[ ! -f "$ENV_FILE" ]]; then
    error "Production deployment requires the repository-root .env file."
    exit 2
  fi

  production_database_password="$(read_root_env_value POSTGRES_PASSWORD)"
  production_minio_password="$(read_root_env_value MINIO_ROOT_PASSWORD)"
  production_signing_key="$(read_root_env_value AUTH_SIGNING_KEY)"

  if [[ -z "$production_database_password" || "$production_database_password" == "aisha_dev" || "$production_database_password" == "-aisha_dev" ]]; then
    error "Set POSTGRES_PASSWORD in the root .env to the real production database password."
    exit 2
  fi
  if [[ -z "$production_minio_password" || "$production_minio_password" == "aisha_minio_dev" || "$production_minio_password" == "replace-me-too" ]]; then
    error "Set MINIO_ROOT_PASSWORD in the root .env to the real production MinIO password."
    exit 2
  fi
  if [[ -z "$production_signing_key" || "$production_signing_key" == "aisha-development-signing-key-change-me" || ${#production_signing_key} -lt 32 ]]; then
    error "Set AUTH_SIGNING_KEY in the root .env to a production value of at least 32 characters."
    exit 2
  fi
fi

if [[ "$PRODUCTION" == false && ! -f "$SEED_FILE" ]]; then
  error "Development seed file not found: $SEED_FILE"
  exit 1
fi

if ! command -v docker >/dev/null 2>&1; then
  error "Docker is not installed or is not available in PATH."
  exit 1
fi

if ! docker info >/dev/null 2>&1; then
  error "Docker daemon is not running or the current user cannot access it."
  exit 1
fi

if ! docker compose version >/dev/null 2>&1; then
  error "Docker Compose v2 is not available."
  exit 1
fi

COMPOSE=(
  docker compose
  --project-directory "$COMPOSE_PROJECT_DIR"
  -f "$COMPOSE_FILE"
)

if [[ -f "$ENV_FILE" ]]; then
  COMPOSE+=(--env-file "$ENV_FILE")
fi

if [[ "$PRODUCTION" == true ]]; then
  COMPOSE+=(-f "$COMPOSE_PROD_FILE")
fi

cd "$ROOT_DIR"

log "Validating Docker Compose configuration"
"${COMPOSE[@]}" config --quiet

log "Starting PostgreSQL, Redis, and MinIO without removing containers or volumes"
"${COMPOSE[@]}" up \
  -d \
  --wait \
  postgres \
  redis \
  minio

log "Creating the required MinIO buckets"
"${COMPOSE[@]}" run \
  --rm \
  --no-deps \
  minio-init

log "Applying pending database migrations"
"${COMPOSE[@]}" --profile tools run \
  --rm \
  --no-deps \
  migration-tools \
  sh -ec 'cd /workspace/infrastructure/scripts && go run ./cmd/migrate -path /workspace/db/migrations -database "$DATABASE_URL" up'

if [[ "$PRODUCTION" == false ]]; then
  log "Loading development seed data"
  "${COMPOSE[@]}" exec \
    -T \
    postgres \
    psql -v ON_ERROR_STOP=1 -U aisha -d aisha < "$SEED_FILE"
else
  log "Skipping development seed data in production mode"
fi

log "Building API and frontend images"
"${COMPOSE[@]}" build \
  api \
  frontend

log "Recreating API and frontend containers"
"${COMPOSE[@]}" up \
  -d \
  --no-deps \
  --force-recreate \
  --wait \
  api \
  frontend

log "Recreating the Nginx container"
"${COMPOSE[@]}" up \
  -d \
  --no-deps \
  --force-recreate \
  --wait \
  nginx

echo
echo "Deployment completed successfully."
echo
echo "Current service status:"
"${COMPOSE[@]}" ps

echo
echo "Services:"
echo "  Frontend:      http://localhost:3033"
echo "  API:           http://localhost:8088"
echo "  Nginx:         http://localhost:8089"
echo "  MinIO API:     http://localhost:9000"
echo "  MinIO Console: http://localhost:9001"
echo "  PostgreSQL:    127.0.0.1:5433"
echo "  Redis:         127.0.0.1:6379"
