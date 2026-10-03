#!/usr/bin/env bash

set -Eeuo pipefail

APP_DIR="${1:?application directory is required}"
IMAGE_TAG="${2:?immutable image tag is required}"

COMPOSE_FILE="$APP_DIR/infrastructure/compose/docker-compose.yml"
COMPOSE_PROD_FILE="$APP_DIR/infrastructure/compose/docker-compose.prod.yml"
ENV_FILE="$APP_DIR/.env"
BACKUP_DIR="$APP_DIR/backups/postgres"
STATE_FILE="$APP_DIR/.deploy-last-successful-tag"

if [[ ! -f "$ENV_FILE" || ! -f "$COMPOSE_FILE" || ! -f "$COMPOSE_PROD_FILE" ]]; then
  echo "Deployment files are missing from $APP_DIR" >&2
  exit 1
fi

require_env_value() {
  local key="$1"
  if ! awk -F= -v key="$key" '
    $1 == key {
      value = $0
      sub(/^[^=]*=/, "", value)
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", value)
      if (value ~ /^".*"$/) {
        sub(/^"/, "", value)
        sub(/"$/, "", value)
      }
      if (value != "") found = 1
    }
    END { exit(found ? 0 : 1) }
  ' "$ENV_FILE"; then
    return 0
  fi
  return 1
}

if ! require_env_value "APP_ENV"; then
  echo "Production environment is missing a non-empty APP_ENV" >&2
  exit 1
fi
if ! grep -Eq "^[[:space:]]*APP_ENV[[:space:]]*=[[:space:]]*(production|prod)[[:space:]]*$" "$ENV_FILE"; then
  echo "Production deployment requires APP_ENV=production" >&2
  exit 1
fi
if grep -Eiq '^[[:space:]]*WEB_BASE_URL[[:space:]]*=.*(localhost|127\.0\.0\.1)' "$ENV_FILE"; then
  echo "Production WEB_BASE_URL cannot point to localhost or 127.0.0.1" >&2
  exit 1
fi
require_env_value "SMTP_HOST"
require_env_value "SMTP_USER"
if ! (require_env_value "SMTP_PASSWORD" || require_env_value "SMTP_PASS"); then
  echo "Production environment needs SMTP_PASSWORD or SMTP_PASS" >&2
  exit 1
fi
if ! (require_env_value "SMTP_FROM" || require_env_value "MAIL_FROM"); then
  echo "Production environment needs SMTP_FROM or MAIL_FROM" >&2
  exit 1
fi
if awk -F= '$1 == "SMTP_USER" || $1 == "SMTP_FROM" || $1 == "MAIL_FROM" { if ($0 ~ /\\@/) found = 1 } END { exit(found ? 0 : 1) }' "$ENV_FILE"; then
  echo "SMTP email values must use @ directly; remove a literal backslash before @" >&2
  exit 1
fi

mkdir -p "$BACKUP_DIR"

compose() {
  IMAGE_TAG="$IMAGE_TAG" docker compose \
    --project-directory "$APP_DIR/infrastructure/compose" \
    --env-file "$ENV_FILE" \
    -f "$COMPOSE_FILE" \
    -f "$COMPOSE_PROD_FILE" \
    "$@"
}

previous_tag=""
if [[ -s "$STATE_FILE" ]]; then
  previous_tag="$(head -n 1 "$STATE_FILE")"
fi

backup_file="$BACKUP_DIR/aisha-$(date -u +%Y%m%dT%H%M%SZ).dump"
rollback() {
  local exit_code=$?
  if ((exit_code != 0)) && [[ -n "$previous_tag" && "$previous_tag" != "$IMAGE_TAG" ]]; then
    echo "Deployment failed; restoring previously healthy image tag $previous_tag" >&2
    IMAGE_TAG="$previous_tag" docker compose \
      --project-directory "$APP_DIR/infrastructure/compose" \
      --env-file "$ENV_FILE" \
      -f "$COMPOSE_FILE" \
      -f "$COMPOSE_PROD_FILE" \
      up -d --force-recreate --wait api frontend nginx || true
  fi
  exit "$exit_code"
}
trap rollback EXIT

cd "$APP_DIR"

echo "Validating production Compose configuration"
compose config --quiet

echo "Starting stateful dependencies without removing volumes"
compose up -d --wait postgres redis minio
compose run --rm --no-deps minio-init

echo "Creating PostgreSQL backup before migration"
compose exec -T postgres pg_dump -U aisha -d aisha --format=custom > "$backup_file"
test -s "$backup_file"

echo "Building immutable images: $IMAGE_TAG"
compose build api frontend

echo "Applying pending migrations without development seed data"
compose --profile tools run --rm --no-deps migration-tools sh -ec \
  'cd /workspace/infrastructure/scripts && go run ./cmd/migrate -path /workspace/db/migrations -database "$DATABASE_URL" up'

echo "Starting application and reverse proxy"
compose up -d --force-recreate --wait api frontend nginx

echo "Running production smoke checks"
curl --fail --silent --show-error --retry 20 --retry-delay 2 http://127.0.0.1/health >/dev/null
compose exec -T api /bin/busybox wget -qO- http://127.0.0.1:8088/health/ready >/dev/null
compose exec -T frontend node -e "fetch('http://127.0.0.1:3033/api/health').then(response => { if (!response.ok) process.exit(1) }).catch(() => process.exit(1))"

printf '%s\n' "$IMAGE_TAG" > "$STATE_FILE"
chmod 600 "$STATE_FILE"
echo "Deployment completed successfully. Backup: $backup_file"

trap - EXIT
