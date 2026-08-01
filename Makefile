.PHONY: bootstrap install fmt lint type-check test test-migrations build check up down logs compose-config smoke wire tools migrate-up migrate-down migrate-reset sqlboiler

BACKEND_BIN := $(CURDIR)/backend/bin
MIGRATE := $(BACKEND_BIN)/migrate
SQLBOILER := $(BACKEND_BIN)/sqlboiler

bootstrap: install
	@test -f .env || cp .env.example .env

install:
	cd frontend && npm ci
	cd backend && go mod download

fmt:
	cd backend && gofmt -w .

lint:
	cd frontend && npm run lint
	cd backend && test -z "$$(gofmt -l .)"
	cd backend && go vet ./...

type-check:
	cd frontend && npm run type-check

test:
	cd frontend && npm test
	cd backend && go test -race ./...

test-migrations:
	docker compose up -d postgres
	docker compose exec -T postgres sh -ec 'dropdb --if-exists -U aisha aisha_migration_test && createdb -U aisha aisha_migration_test'
	docker compose --profile tools run --rm migration-tools go test ./tests -run TestMigrationsRoundTrip -count=1

build:
	cd frontend && npm run build
	cd backend && go build ./...

compose-config:
	docker compose config --quiet
	docker compose -f docker-compose.yml -f docker-compose.prod.yml config --quiet

check: lint type-check test build compose-config

up:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f

smoke:
	curl --fail --silent http://localhost:8080/health/live
	curl --fail --silent http://localhost:8080/health/ready
	curl --fail --silent http://localhost:3000/api/health

wire:
	cd backend && go run github.com/google/wire/cmd/wire ./cmd/api

tools:
	mkdir -p "$(BACKEND_BIN)"
	cd backend/tools && go build -o "$(MIGRATE)" ./cmd/migrate
	cd backend/tools && go build -o "$(SQLBOILER)" github.com/volatiletech/sqlboiler/v4
	cd backend/tools && go build -o "$(BACKEND_BIN)/sqlboiler-psql" github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql

migrate-up:
	docker compose --profile tools run --rm migration-tools sh -ec 'cd tools && go run ./cmd/migrate -path ../migrations -database "$$DATABASE_URL" up'

migrate-down:
	docker compose --profile tools run --rm migration-tools sh -ec 'cd tools && go run ./cmd/migrate -path ../migrations -database "$$DATABASE_URL" down'

migrate-reset: migrate-down migrate-up

sqlboiler:
	docker compose --profile tools run --rm migration-tools sh -ec 'mkdir -p bin && cd tools && go build -o ../bin/sqlboiler github.com/volatiletech/sqlboiler/v4 && go build -o ../bin/sqlboiler-psql github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql && cd /workspace && PATH="/workspace/bin:$$PATH" bin/sqlboiler psql --config db/sqlboiler.toml'
