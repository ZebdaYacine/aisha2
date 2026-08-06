# AISHA Clean Feature-Based Architecture

## 1. Purpose

This document defines the implementation architecture of the AISHA monorepo.

AISHA uses:

```text
Monorepo
+ Feature-Based Architecture
+ Clean Architecture
+ Modular Monolith Backend
+ Next.js App Router Frontend
+ Ports and Adapters
+ Dependency Injection
+ Optional MVVM-style ViewModel layer on the frontend
```

Business logic must not be placed in Next.js route files, React presentational components, frontend API clients, Fiber handlers, route registration files, SQLBoiler generated models, migrations, or generic utility folders.

This document has two parts and an appendix:

- **Part I — Frontend** (§4–§11): Next.js App Router, feature-based, ViewModel-optional.
- **Part II — Backend** (§12–§40): Go Fiber v3, feature-based Clean Architecture,
  fully bound to concrete packages, code shapes, domain rules, and CI commands —
  this is the MVP-ready version of the backend, not just the conceptual layering.

## 2. Dependency Rules

```text
Outer layers may depend on inner layers.
Inner layers must never depend on outer layers.
```

Backend:

```text
Server / Data / Infrastructure
             ↓
        Application
             ↓
           Domain
```

Frontend:

```text
App Routes
    ↓
Views and ViewModels
    ↓
Application / Domain
    ↓
Repository Contracts
    ↑
Data Repository Implementations
    ↓
Shared API Infrastructure
```

The frontend communicates with the backend only through stable API contracts.

## 3. Monorepo Tree

```text
aisha/
├── apps/
│   ├── web/
│   └── api/
├── packages/
│   ├── contracts/
│   ├── eslint-config/
│   ├── typescript-config/
│   └── testing/
├── infrastructure/
│   ├── docker/
│   ├── compose/
│   ├── scripts/
│   └── monitoring/
├── docs/
│   ├── architecture.md
│   ├── backend.md
│   ├── frontend.md
│   ├── workflows.md
│   ├── requirements.md
│   ├── jira_project.md
│   ├── deployment.md
│   ├── security.md
│   ├── ui_design.md
│   ├── userstory.md
├── .github/workflows/
├── .env.example
├── Makefile
├── README.md
```

# Part I — Frontend

## 4. Frontend Base Structure

```text
apps/web/
├── app/
├── core/
├── features/
├── messages/
├── public/
├── tests/
├── next.config.ts
├── middleware.ts
├── package.json
├── tsconfig.json
└── .env.example
```

### Responsibilities

```text
app/       Routes, layouts, metadata, loading, errors, page composition
core/      Shared technical infrastructure and reusable UI
features/  Business feature modules
messages/  Translation files
public/    Static assets and PWA files
tests/     Unit, integration, and E2E tests
```

## 5. Frontend Routes

Route files must stay thin. They may read parameters, define metadata, and compose feature views. They must not implement repositories, perform scattered raw fetch calls, duplicate backend authorization, or contain trusted calculations.

```text
apps/web/app/
├── layout.tsx
├── page.tsx
├── globals.css
├── login/
├── register/
├── init-email/
├── confirmOTP/
├── init-password/
├── products/
│   ├── page.tsx
│   └── [slug]/
├── categories/[slug]/
├── collections/[slug]/
├── regions/[slug]/
├── artisans/
│   ├── page.tsx
│   └── [slug]/
├── search/
├── wishlist/
├── cart/
├── checkout/
│   ├── page.tsx
│   └── success/
├── payment/[paymentId]/
├── account/
│   ├── layout.tsx
│   ├── page.tsx
│   ├── profile/
│   ├── addresses/
│   ├── orders/
│   │   ├── page.tsx
│   │   └── [id]/
│   ├── wishlist/
│   └── reviews/
├── custom-orders/
│   ├── page.tsx
│   ├── new/
│   └── [id]/
├── artisan/
│   ├── layout.tsx
│   ├── page.tsx
│   ├── profile/
│   ├── products/
│   │   ├── page.tsx
│   │   ├── new/
│   │   └── [id]/
│   ├── orders/
│   └── custom-orders/
└── admin/
    ├── layout.tsx
    ├── page.tsx
    ├── users/
    ├── artisans/
    ├── catalogue/
    ├── moderation/
    ├── warehouse/
    ├── inventory/
    ├── orders/
    ├── payments/
    ├── shipments/
    ├── custom-orders/
    └── audit/
```

## 6. Frontend Core

```text
apps/web/core/
├── api/
│   ├── client.ts
│   ├── server-client.ts
│   ├── endpoints.ts
│   ├── errors.ts
│   └── token-manager.ts
├── context/
│   ├── AuthContext.tsx
│   ├── CartContext.tsx
│   └── AppProviders.tsx
├── components/
│   ├── ui/
│   ├── layout/
│   ├── commerce/
│   ├── editorial/
│   ├── feedback/
│   ├── guards/
│   └── forms/
├── hooks/
├── lib/
├── providers/
├── zod/
├── translation.ts
└── types.ts
```

`core` contains generic reusable infrastructure only. Feature-specific rules stay in their feature.

## 7. Frontend Feature Structure

```text
apps/web/features/<feature>/
├── index.ts
├── domain/
│   ├── entities/
│   ├── value-objects/
│   ├── rules/
│   ├── errors/
│   ├── repositories/
│   └── services/
├── application/
│   ├── commands/
│   ├── queries/
│   ├── dto/
│   ├── mappers/
│   └── use-cases/
├── data/
│   ├── api/
│   ├── repositories/
│   ├── mappers/
│   ├── cache/
│   └── storage/
├── viewmodel/
│   ├── hooks/
│   ├── state/
│   ├── actions/
│   └── types/
├── view/
├── components/
│   ├── presentational/
│   ├── containers/
│   ├── forms/
│   └── feedback/
├── schemas/
├── hooks/
├── types/
├── utils/
├── constants/
└── tests/
    ├── domain/
    ├── application/
    ├── data/
    ├── viewmodel/
    └── view/
```

Not every simple feature needs every folder.

## 8. Frontend Layer Responsibilities

### Domain

Contains entities, value objects, framework-independent rules, repository contracts, and domain errors.

The frontend domain must not import React, Next.js, TanStack Query, Axios, `fetch`, browser APIs, or UI components.

### Application

Optional layer for commands, queries, use cases, DTOs, and reusable orchestration shared across several ViewModels.

### Data

Contains API calls, request/response types, repository implementations, API-to-domain mappers, cache adapters, and storage adapters.

### ViewModel

Coordinates loading, empty, error, filter, pagination, form, selection, and action state.

Allowed forms:

```text
Custom React hook
Plain TypeScript class
Reducer with actions
TanStack Query wrapper
State-machine adapter
```

The ViewModel must not render JSX, use database concepts, call random endpoints when a repository exists, or duplicate trusted backend rules.

### View

Consumes ViewModels and composes the feature screen. Views do not implement repositories, call raw endpoints, or calculate trusted business values.

### Components

```text
components/
├── presentational/
├── containers/
├── forms/
└── feedback/
```

Presentational components only receive props and emit callbacks.

## 9. Frontend Features

```text
apps/web/features/
├── auth/
├── account/
├── address/
├── catalogue/
├── category/
├── collection/
├── region/
├── product/
├── artisan/
├── search/
├── cart/
├── checkout/
├── order/
├── payment/
├── shipment/
├── wishlist/
├── review/
├── custom-order/
├── moderation/
├── warehouse/
├── inventory/
└── admin/
```

## 10. Example Frontend Product Feature

```text
features/product/
├── index.ts
├── domain/
│   ├── entities/product.ts
│   ├── value-objects/product-price.ts
│   ├── rules/product-rules.ts
│   ├── errors/product-errors.ts
│   └── repositories/product-repository.ts
├── application/
│   ├── queries/get-products.ts
│   ├── queries/get-product-by-slug.ts
│   ├── commands/
│   ├── dto/
│   └── mappers/
├── data/
│   ├── api/product-api.ts
│   ├── api/product-api.types.ts
│   ├── api/endpoints.ts
│   ├── repositories/api-product-repository.ts
│   └── mappers/product-mapper.ts
├── viewmodel/
│   ├── product-list-viewmodel.ts
│   ├── product-details-viewmodel.ts
│   ├── product-state.ts
│   ├── product-actions.ts
│   └── hooks/
│       ├── use-product-list-viewmodel.ts
│       └── use-product-details-viewmodel.ts
├── view/
│   ├── ProductListingView.tsx
│   ├── ProductDetailsView.tsx
│   └── index.ts
├── components/
│   ├── presentational/
│   │   ├── ProductCard.tsx
│   │   ├── ProductGrid.tsx
│   │   ├── ProductGallery.tsx
│   │   └── ProductPrice.tsx
│   ├── containers/
│   ├── forms/
│   └── feedback/
├── schemas/
├── hooks/
├── types/
├── utils/
└── tests/
```

## 11. Internationalization and RTL

Supported locales:

```text
ar
fr
en
es
```

```text
apps/web/messages/
├── ar.json
├── fr.json
├── en.json
└── es.json
```

All visible text uses translation keys. Arabic uses `dir="rtl"`. Use logical CSS properties and test complete pages, forms, tables, filters, drawers, breadcrumbs, and carousels in RTL.

# Part II — Backend

Part II is written to be directly buildable: every layer below names the exact
package it is implemented with, the exact code shape handlers/use
cases/repositories take, and the exact commands CI runs. It supersedes any
more abstract "framework-independent" phrasing from earlier drafts wherever
the two disagree — concrete wins.

## 12. Backend Goal and Required Packages

Build a Go Fiber API that implements the business workflow without placing
business logic in handlers or SQLBoiler models.

```text
github.com/gofiber/fiber/v3            → server/, feature server/ handlers
github.com/casbin/casbin/v2            → internal/pkg/authorization
github.com/volatiletech/sqlboiler/v4   → features/<feature>/data/models, data/repositories
github.com/google/wire                 → internal/container (compile-time DI)
github.com/jackc/pgx/v5                → internal/pkg/database (Postgres driver under SQLBoiler)
github.com/redis/go-redis/v9           → rate limiting, OTP state, short-lived cache
github.com/minio/minio-go/v7           → internal/pkg/storage
github.com/golang-migrate/migrate/v4   → db/migrations tooling (Makefile target, never in-process)
github.com/google/uuid                 → domain IDs
github.com/go-playground/validator/v10 → server/request struct tags
```

Choose exact versions compatible with the repository and current Go version.
Do not mix multiple HTTP routers or ORMs, and do not substitute an equivalent
(no Gin, no GORM, no manual JWT parsing) — these choices are fixed.

Fiber v3 uses value-receiver contexts (`fiber.Ctx`, not `*fiber.Ctx`) —
handler signatures across every feature follow this consistently (§20).

### Coding rules

- Run `gofmt`.
- Handle every error; wrap errors with context.
- Use `context.Context` in application and repository methods.
- Keep interfaces small and consumer-owned.
- Use constructor injection; avoid global mutable state.
- Use UTC timestamps.
- Use structured logs with correlation IDs.
- Use integer minor units for money.
- Never expose SQLBoiler models as API responses.

## 13. Backend Base Structure

```text
apps/api/
├── cmd/
│   └── main.go
├── internal/
│   ├── bootstrap/
│   ├── config/
│   ├── container/
│   ├── server/
│   ├── pkg/
│   └── features/
├── db/
│   ├── schema/
│   ├── migrations/
│   ├── seeds/
│   ├── queries/
│   └── sqlboiler.toml
├── tests/
│   ├── unit/
│   ├── integration/
│   ├── concurrency/
│   └── e2e/
├── docs/
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
└── .env.example
```

`cmd/main.go` contains startup only — a single call into `bootstrap`:

```go
func main() {
    app := bootstrap.NewApp()
    defer app.Close()
    app.Server.Run()
}
```

### `internal/bootstrap/` (owns startup sequencing)

```text
apps/api/internal/bootstrap/
├── app.go        # Assembles config + container + server, returns a runnable App
├── database.go   # Opens the Postgres pool, runs health check, wires it into container
└── env.go        # Loads .env / process env into internal/config structs, fails fast on missing keys
```

Distinct from `internal/container/` (DI wiring of feature modules) and
`internal/config/` (typed config values). Rule: `bootstrap` may import
`config`, `container`, and `server`. Nothing in `container`, `server`, `pkg`,
or `features` may import `bootstrap` — the dependency arrow points one way.
`bootstrap` is also the only package allowed to call `os.Exit` or panic on
startup misconfiguration.

## 14. Backend Shared Structure

```text
apps/api/internal/
├── config/
│   ├── config.go
│   ├── database.go
│   ├── redis.go
│   ├── minio.go
│   ├── jwt.go
│   ├── mail.go
│   ├── payment.go
│   ├── shipping.go
│   └── validation.go
├── container/
│   ├── wire.go            # //go:build wireinject — provider sets, e.g. ProductProviderSet
│   ├── wire_gen.go         # generated by `wire`, committed to VCS, never hand-edited
│   ├── infrastructure.go
│   ├── features.go
│   └── lifecycle.go
├── server/
│   ├── server.go
│   ├── router.go
│   ├── middleware/
│   ├── routes/
│   └── error_handler/
├── pkg/
│   ├── apperror/
│   ├── auth/
│   ├── authorization/
│   ├── clock/
│   ├── database/
│   ├── email/
│   ├── events/
│   ├── idempotency/
│   ├── money/
│   ├── pagination/
│   ├── response/
│   ├── storage/
│   ├── transaction/
│   ├── validator/
│   └── utils/
└── features/
```

## 15. Backend Feature Structure

```text
apps/api/internal/features/<feature>/
├── module.go
├── domain/
│   ├── entity.go
│   ├── aggregate.go
│   ├── value_objects.go
│   ├── repository.go       # interfaces declared here — see §16, §20
│   ├── services.go
│   ├── policy.go
│   ├── rules.go
│   ├── events.go
│   └── errors.go
├── application/
│   ├── commands/
│   │   ├── command.go
│   │   └── handler.go
│   ├── queries/
│   │   ├── query.go
│   │   └── handler.go
│   ├── usecases/
│   ├── ports/                # external-provider interfaces, e.g. PaymentGateway
│   ├── dto.go
│   ├── mapper.go
│   └── service.go
├── data/
│   ├── repositories/
│   │   └── sql_repository.go
│   ├── mappers/
│   │   └── persistence_mapper.go
│   ├── models/
│   ├── queries/
│   ├── cache/
│   └── adapters/
├── server/
│   ├── router.go
│   ├── routes.go
│   ├── handler.go
│   ├── middleware/
│   ├── request/
│   ├── response/
│   └── presenter/
└── mocks/                    # generated — see §37, never hand-edited
    ├── mock_<feature>_repository.go
    └── mock_<port>.go
```

## 16. Backend Layer Responsibilities

### Domain

Owns entities, aggregates, value objects, invariants, state transitions,
repository contracts, policies, services, events, and errors.

Every repository interface a feature's data layer implements is declared
directly in that feature's `domain/repository.go` — not inferred from naming
convention:

```go
// features/product/domain/repository.go
package domain

type ProductRepository interface {
    GetByID(ctx context.Context, id string) (*Product, error)
    GetBySlug(ctx context.Context, slug string) (*Product, error)
    Create(ctx context.Context, product *Product) error
    Update(ctx context.Context, product *Product) error
    List(ctx context.Context, f ProductFilter) ([]*Product, int, error)
}
```

Method naming is fixed: `GetByID` / `Create` / `Update` / `List`, consistently
across every feature — not `FindByID`/`Save` in one feature and `GetByID`/
`Create` in another.

Domain must not import Fiber, SQLBoiler, PostgreSQL drivers, Redis, MinIO,
Casbin implementation, provider SDKs, or HTTP structs.

### Application

Owns commands, queries, handlers, use cases, DTOs, transaction boundaries,
authorization orchestration, ownership checks, idempotency, audit, and
outbox coordination. Execution order inside a use case: authorise → validate
command and domain invariants → load required data → open transaction where
needed → call repositories → write audit/outbox event → commit → return DTO.

### Data

Owns SQLBoiler repositories, persistence mappers, optimized queries, Redis
adapters, MinIO adapters, cache, and external provider adapters. SQLBoiler
models remain inside data and are never exposed as API responses.

```go
// features/product/data/repositories/sql_product_repository.go
package repositories

type sqlProductRepository struct{ db *sql.DB }

func NewSQLProductRepository(db *sql.DB) domain.ProductRepository {
    return &sqlProductRepository{db: db}
}
```

`container/features.go` (generated via `wire`, §14) wires the concrete type
behind the interface at startup — nothing above `data/` ever imports
`sqlProductRepository` directly, only `domain.ProductRepository`. Transaction-
aware ports may receive a `Tx` abstraction, or a repository constructed from
a transaction scope — `*sql.Tx`/pgx equivalents must never leak above `data/`.

### Server

Owns router, route registration, handlers, feature middleware, request
structs, response structs, validation, mapping, HTTP status handling, and
presenters. See §20 for the concrete handler shape.

Handlers do not: open SQL queries, decide product state transitions,
calculate trusted prices, implement Casbin policy manually, or call MinIO
directly.

## 17. Global Server Structure

```text
apps/api/internal/server/
├── server.go
├── router.go
├── middleware/
│   ├── request_id.go
│   ├── recovery.go
│   ├── logging.go
│   ├── cors.go
│   ├── authentication.go
│   ├── rate_limit.go
│   └── security_headers.go
├── routes/
│   ├── health.go
│   ├── public.go
│   ├── authenticated.go
│   └── admin.go
└── error_handler/
    └── error_handler.go
```

Global middleware belongs here. JWT verification happens exactly once, in
`middleware/authentication.go` — feature middleware only adds
authorization/ownership checks on top of an already-verified principal.

## 18. Feature Server Structure

```text
features/<feature>/server/
├── router.go
├── routes.go
├── handler.go
├── middleware/
│   ├── authentication.go
│   ├── authorization.go
│   └── validation.go
├── request/
│   ├── create_request.go
│   ├── update_request.go
│   ├── filter_request.go
│   └── mapper.go
├── response/
│   ├── response.go
│   ├── list_response.go
│   ├── detail_response.go
│   └── mapper.go
└── presenter/
    └── presenter.go
```

Feature middleware contains only feature-specific concerns. Do not duplicate
global authentication middleware.

## 19. Router, Request, and Response Rules

### Router

Groups feature routes and delegates to handlers.

### Request

Owns JSON names, validation tags, path/query/header parsing, and mapping to
application commands or queries. It contains no business logic.

### Response

Defines the public JSON contract and maps from application DTOs. It never
exposes SQLBoiler models, secrets, or provider internals.

### Handler

Extracts context and principal, parses input, validates transport data, maps
to a command/query, calls the application layer, maps the result, and
returns the correct status.

## 20. Handler / Use Case / Repository Pattern (concrete shapes)

### Handler (`features/<feature>/server/handler.go`)

```go
func (h *ProductHandler) Create(c fiber.Ctx) error {
    principal, err := h.auth.RequirePrincipal(c)
    if err != nil {
        return h.errors.Write(c, err)
    }

    var req CreateProductRequest
    if err := c.Bind().Body(&req); err != nil {
        return h.errors.Write(c, ErrInvalidJSON)
    }
    if err := h.validator.Struct(req); err != nil {
        return h.errors.Write(c, NewValidationError(err))
    }

    result, err := h.createProduct.Execute(c.Context(), principal, req.ToCommand())
    if err != nil {
        return h.errors.Write(c, err)
    }
    return c.Status(fiber.StatusCreated).JSON(ProductResponseFrom(result))
}
```

### Use case (`features/<feature>/application/usecases/create_product.go`)

```go
type CreateProductCommand struct {
    ArtisanID    string
    CategoryID   string
    PriceMinor   int64
    Currency     string
    Translations []ProductTranslationInput
}

type CreateProduct interface {
    Execute(ctx context.Context, actor Principal, cmd CreateProductCommand) (ProductDTO, error)
}
```

### Repository

Method names and shape are fixed by §16 (`GetByID` / `Create` / `Update` /
`List`). See the `sqlProductRepository` example in §16.

### Request flow — public endpoint (no auth)

```text
Client
  → router.go (feature)                     [route registration only]
  → middleware/ (global: request_id, recovery, logging, cors, rate_limit)
  → handler.go                              [parse + validate request struct]
  → application/queries/handler.go          [orchestration, no I/O logic itself]
  → data/repositories/sql_repository.go     [SQLBoiler query, mapped via persistence_mapper]
  → domain/entity.go                        [returned entity, business invariants already enforced]
  ← application/dto.go + mapper.go          [entity → DTO]
  ← server/response/mapper.go               [DTO → public JSON contract]
  ← Client
```

### Request flow — authenticated endpoint (JWT access token)

```text
Client (Authorization: Bearer <accessToken>)
  → router.go
  → middleware/authentication.go (global)   [verify signature + expiry, extract principal]
  → middleware/authorization.go (feature)   [role/policy check via internal/pkg/authorization]
  → handler.go                              [principal now in context]
  → application/commands/handler.go         [ownership check, transaction boundary, idempotency]
  → domain/policy.go                        [authoritative business rule — frontend guard is UX only]
  → data/repositories/sql_repository.go
  ← ... same response path as above
```

Refresh flow: access tokens are short-lived, refresh tokens are stored as
session records (Postgres) and rotated on use; `internal/pkg/auth` owns
signing/verification, `features/auth` owns the issue/refresh/revoke use
cases. Redis is used only for OTP state and rate limiting, never as the
token source of truth.

## 21. Backend Features

```text
apps/api/internal/features/
├── auth/
├── user/
├── address/
├── artisan/
├── category/
├── collection/
├── region/
├── product/
├── product_media/
├── moderation/
├── warehouse/
├── inventory/
├── reservation/
├── cart/
├── checkout/
├── order/
├── payment/
├── shipment/
├── custom_order/
├── wishlist/
├── review/
├── notification/
├── audit/
└── admin/
```

## 22. Example Backend Product Feature

```text
features/product/
├── module.go
├── domain/
│   ├── product.go
│   ├── product_id.go
│   ├── product_status.go
│   ├── product_repository.go   # GetByID / GetBySlug / Create / Update / List
│   ├── product_policy.go
│   ├── product_rules.go        # submission / activation invariants, see §23
│   ├── product_events.go
│   └── product_errors.go
├── application/
│   ├── commands/
│   │   ├── create_product.go
│   │   ├── update_product.go
│   │   ├── submit_product.go
│   │   └── suspend_product.go
│   ├── queries/
│   │   ├── get_product.go
│   │   └── list_products.go
│   ├── ports/
│   ├── product_dto.go
│   └── product_mapper.go
├── data/
│   ├── repositories/sql_product_repository.go
│   ├── mappers/product_persistence_mapper.go
│   ├── queries/product_queries.go
│   └── models/
├── server/
│   ├── router.go
│   ├── handler.go
│   ├── middleware/product_access.go
│   ├── request/
│   │   ├── create_product_request.go
│   │   ├── update_product_request.go
│   │   ├── product_filter_request.go
│   │   └── request_mapper.go
│   ├── response/
│   │   ├── product_response.go
│   │   ├── product_list_response.go
│   │   └── response_mapper.go
│   └── presenter/product_presenter.go
└── mocks/
    ├── mock_product_repository.go
    └── mock_payment_gateway.go   # if product feature consumes a port
```

## 23. Domain Rules by Feature

Each feature's `domain/rules.go` and `domain/policy.go` carry these
invariants. This is the pattern every feature follows, illustrated on the
features currently in MVP scope.

**`features/product/domain/`** — submission is allowed only when: artisan is
approved and active, required translation fields exist, category exists,
price is positive, currency is supported by config, required media exists,
and current status permits submission. Activation is allowed only when:
moderation is approved, product is not suspended, and accepted inventory is
available.

**`features/inventory/domain/`** — append-only movement log, never a
mutable balance edited in place:

```text
RECEPTION · INSPECTION_ACCEPT · INSPECTION_REJECT · INSPECTION_QUARANTINE
INSPECTION_DAMAGE · RESERVATION · RESERVATION_RELEASE · ORDER_COMMIT
SHIPMENT · ADJUSTMENT_IN · ADJUSTMENT_OUT · RETURN
```

An `inventory_balances` table may exist in `data/` purely as a read
optimisation; every value in it must be reconstructable from the movement
log, and the movement log — not the balance table — is the source of truth.

**`features/reservation/domain/`** — required fields: ID, product/stock item
reference, order or checkout reference, quantity, status, `expires_at`,
`created_at`, `consumed_at`/`released_at`. A background job expires stale
active reservations; this expiry job lives in `internal/pkg` or a
feature-owned worker, not inside a request handler.

**`features/order/domain/`** — order lines are immutable snapshots taken at
checkout time. Never join current product price to render a historical
order total — that is a correctness bug, not a style preference.

**`features/payment/domain/`** — persisted fields: internal payment ID,
order ID, provider, provider payment reference, amount minor, currency,
status, idempotency key, created/updated timestamps. Provider webhook events
are stored separately from the payment record itself, keyed by a unique
external event ID, so a provider retry cannot double-apply.

## 24. Database Migrations

Migration rules:

- Every schema change has `up` and `down`.
- Prefer additive changes.
- Add explicit foreign keys.
- Add unique constraints.
- Add indexes for frequent lookups.
- Define deletion behaviour explicitly.
- Production migrations never run invisibly inside API startup — `bootstrap/`
  (§13) opens connections and health-checks, it does not migrate.
- Regenerate SQLBoiler models after schema changes — a stale generated model
  is treated as a build break, same policy as a stale mock (§37).

Migration sequence:

```text
000001_extensions
000002_users_sessions
000003_artisans
000004_catalogue
000005_moderation
000006_warehouse
000007_inventory
000008_cart_addresses
000009_orders
000010_payments
000011_shipments
000012_custom_orders
000013_notifications_outbox
000014_audit_idempotency
```

```bash
migrate -path db/migrations -database "$DATABASE_URL" up
```

## 25. Cross-Feature Communication

A feature must never import another feature's concrete data implementation.

Allowed mechanisms:

```text
Application port
Exported application query
Domain service interface
Stable shared contract
Domain event
Integration event
Outbox event
```

## 26. Feature Mapping

| Capability | Frontend | Backend |
|---|---|---|
| Authentication | `auth` | `auth` |
| Profile | `account` | `user` |
| Addresses | `address` | `address` |
| Artisan | `artisan` | `artisan` |
| Catalogue | `catalogue` | `category`, `collection`, `region` |
| Products | `product` | `product`, `product_media` |
| Search | `search` | catalogue queries |
| Moderation | `moderation` | `moderation` |
| Warehouse | `warehouse` | `warehouse` |
| Inventory | `inventory` | `inventory`, `reservation` |
| Cart | `cart` | `cart` |
| Checkout | `checkout` | `checkout`, `reservation` |
| Orders | `order` | `order` |
| Payments | `payment` | `payment` |
| Shipments | `shipment` | `shipment` |
| Wishlist | `wishlist` | `wishlist` |
| Reviews | `review` | `review` |
| Custom orders | `custom-order` | `custom_order` |
| Administration | `admin` | `admin`, `audit` |

## 27. API Contract

Every endpoint defines method, path, authentication, permission, request
schema, response schema, stable error codes, pagination, idempotency, and
examples.

```json
{
  "data": {},
  "meta": {},
  "requestId": "uuid"
}
```

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "The request is invalid.",
    "fields": {
      "email": "INVALID_EMAIL"
    }
  },
  "requestId": "uuid"
}
```

## 28. API Endpoints and Routing

Every endpoint group below registers in its feature's
`features/<feature>/server/routes.go`, then is mounted from the global
`internal/server/routes/{public,authenticated,admin}.go` (§17).

### Authentication

```text
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password
GET    /api/v1/me
```

### Public Catalogue

```text
GET    /api/v1/categories
GET    /api/v1/products
GET    /api/v1/products/:id
GET    /api/v1/artisans
GET    /api/v1/artisans/:id
```

### Artisan

```text
POST   /api/v1/artisan-applications
GET    /api/v1/artisan-applications/me
PATCH  /api/v1/artisan/profile
POST   /api/v1/artisan/products
PATCH  /api/v1/artisan/products/:id
POST   /api/v1/artisan/products/:id/media
POST   /api/v1/artisan/products/:id/submit
GET    /api/v1/artisan/products
GET    /api/v1/artisan/inventory
```

### Moderation

```text
GET    /api/v1/admin/product-submissions
POST   /api/v1/admin/product-submissions/:id/approve
POST   /api/v1/admin/product-submissions/:id/request-changes
POST   /api/v1/admin/product-submissions/:id/reject
POST   /api/v1/admin/products/:id/suspend
```

### Warehouse and Inventory

```text
POST   /api/v1/warehouse/receptions
GET    /api/v1/warehouse/receptions
POST   /api/v1/warehouse/receptions/:id/inspect
GET    /api/v1/warehouse/inventory
POST   /api/v1/warehouse/inventory/:productId/adjust
GET    /api/v1/warehouse/orders
POST   /api/v1/warehouse/orders/:id/prepare
POST   /api/v1/warehouse/orders/:id/ship
```

### Cart and Checkout

```text
GET    /api/v1/cart
POST   /api/v1/cart/items
PATCH  /api/v1/cart/items/:id
DELETE /api/v1/cart/items/:id
POST   /api/v1/checkout
GET    /api/v1/orders
GET    /api/v1/orders/:id
POST   /api/v1/orders/:id/cancel
```

### Payment and Shipping Callbacks

```text
POST   /api/v1/payments/:paymentId/retry
GET    /api/v1/payments/:paymentId
POST   /api/v1/webhooks/payments/:provider
POST   /api/v1/webhooks/shipments/:provider
```

### Custom Orders

```text
POST   /api/v1/custom-orders
GET    /api/v1/custom-orders
GET    /api/v1/custom-orders/:id
POST   /api/v1/custom-orders/:id/messages
POST   /api/v1/custom-orders/:id/quotes
POST   /api/v1/custom-orders/:id/quotes/:quoteId/accept
```

### Routing table

| Route group | Feature | Global route file |
|---|---|---|
| `/auth/*`, `/me` | `auth` | `public.go` (register/login/refresh) + `authenticated.go` (`/me`) |
| `/categories`, `/products*`, `/artisans*` (GET) | `catalogue`, `product`, `artisan` | `public.go` |
| `/artisan-applications*`, `/artisan/*` | `artisan` | `authenticated.go` |
| `/admin/product-submissions*`, `/admin/products/:id/suspend` | `moderation` | `admin.go` |
| `/warehouse/*` | `warehouse`, `inventory` | `admin.go` (warehouse_agent role) |
| `/cart*`, `/checkout`, `/orders*` | `cart`, `checkout`, `order` | `authenticated.go` |
| `/payments/*` (non-webhook) | `payment` | `authenticated.go` |
| `/webhooks/payments/:provider`, `/webhooks/shipments/:provider` | `payment`, `shipment` | `public.go`, with signature-verification middleware instead of JWT (§32) |
| `/custom-orders*` | `custom_order` | `authenticated.go` |

Handlers for admin-only groups additionally pass through
`middleware/authorization.go` with the relevant Casbin role.

## 29. Authentication and Authorization

Authentication includes access tokens, refresh tokens, session records,
rotation, revocation, OTP, and rate limiting.

Recommended roles:

```text
customer
artisan
warehouse_agent
moderator
admin
```

Frontend guards improve UX. Backend application policies are authoritative.
See §20 for the concrete public/authenticated request-flow diagrams.

## 30. Validation

Validate at two levels — neither substitutes for the other:

### Request DTO (`server/request/*.go`, `validator/v10` struct tags)

- Required fields.
- String length.
- Numeric range.
- UUID format.
- Supported file metadata.
- ISO currency format.
- Pagination bounds.

This layer never sees a `Principal` or touches the database.

### Domain/Application (`domain/policy.go`, `application/commands/handler.go`)

- Actor permission.
- Resource ownership.
- Valid status transition.
- Stock availability.
- Product readiness.
- Amount consistency.
- Duplicate callback.
- Idempotency.

Do not rely on frontend validation (§8's ViewModel layer) — it is UX only
and is never trusted here.

## 31. Error Catalogue

```text
VALIDATION_ERROR · INVALID_CREDENTIALS · AUTHENTICATION_REQUIRED · FORBIDDEN
RESOURCE_NOT_FOUND · CONFLICT · ARTISAN_NOT_APPROVED · PRODUCT_NOT_EDITABLE
PRODUCT_NOT_SELLABLE · INVALID_STATE_TRANSITION · OUT_OF_STOCK
RESERVATION_EXPIRED · IDEMPOTENCY_CONFLICT · PAYMENT_AMOUNT_MISMATCH
INVALID_WEBHOOK_SIGNATURE · DUPLICATE_WEBHOOK · UNSUPPORTED_FILE_TYPE
FILE_TOO_LARGE · RATE_LIMITED · INTERNAL_ERROR
```

`internal/pkg/apperror/codes.go` defines these as typed constants;
`internal/server/error_handler/error_handler.go` is the single place that
maps a code to an HTTP status and the public error envelope from §27.
Feature code raises a code, never an HTTP status directly.

## 32. Idempotency

Checkout, payment creation, refund, shipment creation, and external webhook
processing all go through `internal/pkg/idempotency`:

1. Client sends an `Idempotency-Key` header (checkout) or the provider sends
   a unique external event ID (webhooks).
2. `pkg/idempotency` stores key, actor, route, request hash, status, and
   response — this table is part of the `000014_audit_idempotency` migration.
3. Same key + same request hash → return the saved response, no re-execution.
4. Same key + different request hash → `IDEMPOTENCY_CONFLICT`.

Webhook signature verification (`INVALID_WEBHOOK_SIGNATURE`) happens in
feature middleware *before* idempotency is even checked — an unverified
webhook is rejected outright, not deduplicated.

## 33. File Upload

Order of operations, enforced in that order:

1. Validate request size before full processing.
2. Validate MIME by content signature, not by trusting the extension.
3. Generate the object key server-side — the original filename is never
   used as or folded into the storage path.
4. Upload to a private MinIO bucket (§35: buckets private by default).
5. Store metadata in Postgres.
6. Scan if a scanner is configured.
7. Issue a signed URL only for an authorised read.
8. For product media specifically: publish to the public catalogue only
   after moderation approval, not at upload time.

Applies to `product_media`, artisan verification documents, custom-order
attachments, and moderation files alike.

## 34. OpenAPI

Every route documents: authentication requirement, Casbin permission,
request schema, response schema, possible error codes (§31), pagination,
`Idempotency-Key` usage where applicable, and one example payload. The spec
is validated in CI — a route merged without a matching OpenAPI entry fails
the pipeline, not just a code review comment.

## 35. Data Responsibilities

PostgreSQL is the source of truth for users, sessions, artisans, catalogue,
products, moderation, inventory, reservations, carts, orders, payments,
shipments, custom orders, reviews, notifications, audit, idempotency, and
outbox.

Redis is used for rate limiting, short-lived OTP state, justified caching,
session assistance, and short-lived coordination.

MinIO stores product media, artisan documents, custom-order attachments, and
moderation files. Buckets are private by default.

## 36. Testing

```text
apps/api/tests/
├── unit/          # domain state transitions, money calc, inventory invariants,
│                  # order totals, permission policies, DTO mapping — no mocks needed
├── integration/   # real Postgres repositories, transaction rollback, Casbin
│                  # enforcement, MinIO storage, Redis rate limiting, migration up/down
├── concurrency/   # last-unit reservation race, duplicate checkout, duplicate
│                  # webhook, reservation-expiry-vs-payment-confirmation race
└── e2e/           # product → warehouse → sale; customer purchase; payment
                   # failure; full shipment flow
```

`unit/` mirrors `application/` tests using generated mocks (§37); `integration/`
and `concurrency/` require the real test-container stack and never use mocks.

Frontend tests cover schemas, mappers, ViewModels, formatting, authentication,
cart, checkout, filters, admin tables, E2E flows, and Arabic RTL navigation.

## 37. Mock Generation Workflow

Every interface in `domain/repository.go` and every port in
`application/ports/` must have a generated mock, so application and domain
tests never touch a real database, Redis instance, or provider SDK.

```bash
# Regenerate mocks for a single feature's domain + application ports
mockery --dir=internal/features/product/domain --output=internal/features/product/mocks --outpkg=mocks --all
mockery --dir=internal/features/product/application/ports --output=internal/features/product/mocks --outpkg=mocks --all
```

Add to `apps/api/Makefile`:

```makefile
mocks:
	@for f in $$(ls internal/features); do \
		mockery --dir=internal/features/$$f/domain --output=internal/features/$$f/mocks --outpkg=mocks --all 2>/dev/null; \
		mockery --dir=internal/features/$$f/application/ports --output=internal/features/$$f/mocks --outpkg=mocks --all 2>/dev/null; \
	done

wire:
	wire ./internal/container/...
```

Test layering rule:

- **Domain tests**: no mocks needed — pure functions/invariants, no I/O.
- **Application tests**: use generated mocks for repositories and ports;
  assert via `testify/assert` and `testify/require`; verify call expectations
  via `testify/mock`.
- **Data tests**: hit a real (test-container) Postgres/Redis instance —
  mocks are not used here.
- Regenerate mocks whenever a repository or port interface signature changes;
  a stale mock is a build break.

## 38. Required Commands

```bash
go mod download
gofmt -w .
go vet ./...
go test ./...
go build ./...
migrate -path db/migrations -database "$DATABASE_URL" up
```

Combined with the `mocks` and `wire` Makefile targets (§37), this is the full
local + CI command set for the backend.

## 39. Architecture Checklist

### Backend

- Domain imports no framework or infrastructure package.
- Application depends on interfaces.
- Server handlers contain no business rules.
- SQLBoiler models remain inside data.
- Transactions are controlled by application use cases.
- Request and response structs stay in server.
- Feature middleware is separate from global middleware.
- External callbacks are verified and idempotent.
- Money uses integer minor units.
- UTC is used internally.
- Every domain repository interface and application port has a
  corresponding entry in `mocks/`, regenerated on signature change.
- `bootstrap/` contains no business logic and is the only package allowed
  to call `os.Exit` or panic on startup misconfiguration.
- Application-layer tests use mocks; data-layer tests use real
  test-container instances — never the reverse.
- JWT verification happens once, in global `middleware/authentication.go`;
  feature middleware only adds authorization/ownership checks on top.
- Repository interfaces use `GetByID` / `Create` / `Update` naming
  consistently across every feature.
- Every inventory balance is explainable from the movement log, never
  edited independently of it.
- Every route has a matching OpenAPI entry validated in CI; every error
  path uses a code from §31, never a raw HTTP status string.
- Idempotency keys are checked before checkout/payment/refund/shipment
  logic runs, and webhook signatures are verified before idempotency is
  checked.
- File uploads never trust the client-provided filename as a storage path;
  product media is published only post-moderation.
- `go vet`, `gofmt -w .`, `go test ./...`, and `migrate ... up` all pass
  before merge.

### Frontend

- Route files remain thin.
- Domain imports no React or HTTP library.
- Repository contracts live inward.
- Data implements repository contracts.
- Views consume ViewModels.
- Presentational components do not call APIs.
- ViewModels do not render JSX.
- Backend rules are not duplicated as trusted frontend rules.
- All visible text is translated.
- Arabic RTL is tested.
- Loading, empty, error, forbidden, not-found, and offline states exist.

## 40. Definition of Done

The architecture is correctly applied when frontend features support domain,
repository contracts, data implementations, optional application use cases,
ViewModels, views, and components where needed; backend features support
domain, application, data, server, and generated mocks; the server includes
router, middleware, request, response, handler, and presenter
responsibilities; DI wiring is generated via `wire` rather than hand-written;
framework dependencies stay outside domain; persistence models do not leak;
features communicate through stable public interfaces; every route has a
matching OpenAPI entry and stable error code; idempotency and webhook
signature verification are in place for every external callback; and the
system can evolve feature by feature without coupling the whole codebase.
