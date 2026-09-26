.PHONY: bootstrap install fmt lint type-check test test-migrations build check up down logs compose-config smoke wire tools migrate-up migrate-down migrate-reset seed sqlboiler

COMPOSE_FILE := infrastructure/compose/docker-compose.yml
COMPOSE_PROD_FILE := infrastructure/compose/docker-compose.prod.yml
COMPOSE := docker compose -f $(COMPOSE_FILE)
BACKEND_BIN := $(CURDIR)/apps/api/bin
MIGRATE := $(BACKEND_BIN)/migrate
SQLBOILER := $(BACKEND_BIN)/sqlboiler

bootstrap: install
	@test -f .env || cp .env.example .env

install:
	cd apps/web && npm ci
	cd apps/api && go mod download

fmt:
	cd apps/api && gofmt -w .

lint:
	cd apps/web && npm run lint
	cd apps/api && test -z "$$(gofmt -l .)"
	cd apps/api && go vet ./...

type-check:
	cd apps/web && npm run type-check

test:
	cd apps/web && npm test
	cd apps/api && go test -race ./...

test-migrations:
	$(COMPOSE) up -d postgres
	$(COMPOSE) exec -T postgres sh -ec 'dropdb --if-exists -U aisha aisha_migration_test && createdb -U aisha aisha_migration_test'
	$(COMPOSE) --profile tools run --rm migration-tools go test ./tests/integration -run TestMigrationsRoundTrip -count=1

build:
	cd apps/web && npm run build
	cd apps/api && go build ./...

compose-config:
	$(COMPOSE) config --quiet
	docker compose -f $(COMPOSE_FILE) -f $(COMPOSE_PROD_FILE) config --quiet

check: lint type-check test build compose-config

up:
	$(COMPOSE) up --build

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

smoke:
	curl --fail --silent http://localhost:8088/health/live
	curl --fail --silent http://localhost:8088/health/ready
	curl --fail --silent http://localhost:3033/api/health

wire:
	cd apps/api && go run github.com/google/wire/cmd/wire ./cmd

tools:
	mkdir -p "$(BACKEND_BIN)"
	cd infrastructure/scripts && go build -o "$(MIGRATE)" ./cmd/migrate
	cd infrastructure/scripts && go build -o "$(SQLBOILER)" github.com/volatiletech/sqlboiler/v4
	cd infrastructure/scripts && go build -o "$(BACKEND_BIN)/sqlboiler-psql" github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql

migrate-up:
	$(COMPOSE) --profile tools run --rm migration-tools sh -ec 'cd /workspace/infrastructure/scripts && go run ./cmd/migrate -path /workspace/db/migrations -database "$$DATABASE_URL" up'

migrate-down:
	$(COMPOSE) --profile tools run --rm migration-tools sh -ec 'cd /workspace/infrastructure/scripts && go run ./cmd/migrate -path /workspace/db/migrations -database "$$DATABASE_URL" down'

migrate-reset: migrate-down migrate-up

seed: migrate-up
	$(COMPOSE) exec -T postgres psql -v ON_ERROR_STOP=1 -U aisha -d aisha < apps/api/db/seeds/development.sql

sqlboiler:
	$(COMPOSE) --profile tools run --rm migration-tools sh -ec 'mkdir -p /workspace/bin && cd /workspace/infrastructure/scripts && go build -o /workspace/bin/sqlboiler github.com/volatiletech/sqlboiler/v4 && go build -o /workspace/bin/sqlboiler-psql github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql && cd /workspace && PATH="/workspace/bin:$$PATH" bin/sqlboiler psql --config db/sqlboiler.toml'
