.PHONY: fmt vet test build check wire mocks migrate-up migrate-down

fmt:
	gofmt -w cmd internal tests

vet:
	go vet ./...

test:
	go test ./...

build:
	go build ./...

check: fmt vet test build

wire:
	wire ./internal/container/...

mocks:
	@echo "Generate feature mocks with mockery as feature interfaces are added."

migrate-up:
	migrate -path db/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DATABASE_URL)" down
