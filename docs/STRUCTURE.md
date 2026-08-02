# AISHA Backend Project Structure

## 1. Purpose

This document defines the backend folder structure for the AISHA MVP.

The project uses a **feature-based architecture**. Each business feature owns its domain rules, data access, repository implementations, HTTP handlers, request DTOs, response DTOs, routes, and tests.

The architecture must keep:

- Business logic out of Fiber handlers.
- SQL and SQLBoiler details out of the domain layer.
- Infrastructure configuration out of features.
- Shared technical code inside `core`.
- Database schema, migrations, and seed data inside `db`.
- Unit and integration tests inside a dedicated `test` folder.

---

## 2. Root Project Structure

```text
aisha-backend/
├── cmd/
│   └── main.go
│
├── features/
│   ├── auth/
│   ├── users/
│   ├── artisans/
│   ├── categories/
│   ├── products/
│   ├── moderation/
│   ├── warehouse/
│   ├── inventory/
│   ├── cart/
│   ├── checkout/
│   ├── orders/
│   ├── payments/
│   ├── shipments/
│   ├── customorders/
│   ├── notifications/
│   └── admin/
│
├── server/
│   ├── server.go
│   ├── router.go
│   ├── routes.go
│   ├── handlers/
│   ├── middleware/
│   ├── security/
│   ├── request/
│   ├── response/
│   └── errors/
│
├── core/
│   ├── container/
│   ├── config/
│   ├── database/
│   ├── cache/
│   ├── storage/
│   ├── mail/
│   ├── logger/
│   ├── security/
│   ├── transaction/
│   ├── validator/
│   ├── clock/
│   ├── pagination/
│   ├── idempotency/
│   ├── audit/
│   ├── outbox/
│   └── utils/
│
├── db/
│   ├── schema/
│   ├── migrations/
│   ├── seeds/
│   ├── runner/
│   ├── sqlboiler/
│   └── queries/
│
├── test/
│   ├── unit/
│   ├── integration/
│   ├── concurrency/
│   ├── e2e/
│   ├── fixtures/
│   ├── mocks/
│   └── testutil/
│
├── jobs/
│   ├── reservationexpiry/
│   ├── outboxprocessor/
│   ├── notifications/
│   └── webhookretry/
│
├── openapi/
│   ├── openapi.yaml
│   └── examples/
│
├── policies/
│   └── casbin/
│
├── scripts/
│   ├── generate-models.sh
│   ├── migrate.sh
│   ├── seed.sh
│   └── test.sh
│
├── .env
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── sqlboiler.toml
├── wire.go
└── wire_gen.go
```

---

## 3. `cmd` Folder

```text
cmd/
└── main.go
```

`cmd/main.go` is the single application entry point.

It must remain small. It starts the application but does not contain business logic, route definitions, SQL queries, or infrastructure implementation details.

### Main startup sequence

```text
1. Load environment variables.
2. Validate application configuration.
3. Initialise structured logger.
4. Connect to PostgreSQL.
5. Run database migrations.
6. Run configured database seeds.
7. Connect to Redis.
8. Connect to MinIO.
9. Initialise mail sender.
10. Initialise Casbin and authentication services.
11. Build the dependency container.
12. Create feature services and handlers.
13. Register middleware and routes.
14. Start background jobs.
15. Start the Fiber HTTP server.
16. Wait for termination signal.
17. Stop jobs and close all connections gracefully.
```

Example:

```go
package main

import (
    "context"
    "log"
    "os/signal"
    "syscall"

    "aisha/core/container"
)

func main() {
    ctx, stop := signal.NotifyContext(
        context.Background(),
        syscall.SIGINT,
        syscall.SIGTERM,
    )
    defer stop()

    app, cleanup, err := container.Build(ctx)
    if err != nil {
        log.Fatalf("build application: %v", err)
    }
    defer cleanup()

    if err := app.Run(ctx); err != nil {
        log.Fatalf("run application: %v", err)
    }
}
```

### Migration and seed configuration

Migrations and seeds are started from the application bootstrap, but their behaviour must be controlled through environment variables.

```env
DB_MIGRATE_ON_START=true
DB_SEED_ON_START=true
DB_SEED_SET=development
```

Recommended behaviour:

- Development: run migrations and development seeds automatically.
- Test: run migrations and deterministic test seeds automatically.
- Production: migrations may run on startup only when explicitly enabled.
- Production reference data seeds must be idempotent.
- Demo or fake user seeds must never run in production.

---

## 4. Feature-Based Architecture

Every business capability is placed inside `features`.

```text
features/
├── auth/
├── products/
├── admin/
├── orders/
└── ...
```

Each feature is independent and should expose only the minimum required public API.

### Standard feature structure

```text
features/products/
├── domain/
│   ├── product.go
│   ├── status.go
│   ├── errors.go
│   ├── events.go
│   ├── rules.go
│   └── repository.go
│
├── application/
│   ├── commands/
│   │   ├── create_product.go
│   │   ├── update_product.go
│   │   └── submit_product.go
│   ├── queries/
│   │   ├── get_product.go
│   │   └── list_products.go
│   ├── dto/
│   │   └── product_dto.go
│   └── service.go
│
├── data/
│   ├── models/
│   ├── mapper/
│   ├── queries/
│   └── filters/
│
├── repository/
│   ├── postgres_repository.go
│   ├── cache_repository.go
│   └── transaction_repository.go
│
├── server/
│   ├── handler.go
│   ├── routes.go
│   ├── request/
│   │   ├── create_product.go
│   │   └── update_product.go
│   └── response/
│       ├── product.go
│       └── product_list.go
│
├── module.go
└── providers.go
```

A small feature may omit folders it does not need, but it must preserve the dependency direction.

```text
server -> application -> domain
repository/data -> domain interfaces
core -> shared technical capabilities
```

The domain must not import Fiber, PostgreSQL, Redis, MinIO, SQLBoiler, SMTP, or HTTP packages.

---

## 5. Feature Responsibilities

## 5.1 `domain`

The `domain` folder contains pure business rules.

```text
features/products/domain/
├── product.go
├── status.go
├── errors.go
├── events.go
├── rules.go
└── repository.go
```

It contains:

- Entities.
- Value objects.
- Domain services.
- Domain errors.
- State transitions.
- Business invariants.
- Repository interfaces.
- Domain events.

Example:

```go
type ProductRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
    Create(ctx context.Context, product *Product) error
    Update(ctx context.Context, product *Product) error
}
```

The interface belongs to the domain or application consumer. PostgreSQL implements it in the repository folder.

The domain must not use SQLBoiler models directly.

---

## 5.2 `application`

The `application` folder orchestrates use cases.

```text
features/products/application/
├── commands/
├── queries/
├── dto/
└── service.go
```

It is responsible for:

- Authorisation.
- Command validation.
- Loading domain entities.
- Calling domain methods.
- Starting transactions.
- Calling repositories.
- Writing audit records.
- Writing outbox events.
- Returning application DTOs.

Example:

```go
type CreateProductCommand struct {
    ArtisanID   uuid.UUID
    CategoryID  uuid.UUID
    PriceMinor  int64
    Currency    string
    Translations []TranslationInput
}

type CreateProduct interface {
    Execute(
        ctx context.Context,
        actor auth.Principal,
        cmd CreateProductCommand,
    ) (ProductDTO, error)
}
```

Application DTOs must not expose SQLBoiler models.

---

## 5.3 `data`

The `data` folder contains persistence-oriented structures needed by the feature.

```text
features/products/data/
├── models/
├── mapper/
├── queries/
└── filters/
```

It may contain:

- Database filter structures.
- Query parameters.
- SQLBoiler-to-domain mappers.
- Domain-to-SQLBoiler mappers.
- Database row helper structures.
- Feature-specific query builders.

Example:

```go
func ToDomain(model *models.Product) (*domain.Product, error) {
    if model == nil {
        return nil, errors.New("product model is nil")
    }

    return domain.RestoreProduct(
        model.ID,
        model.ArtisanID,
        model.CategoryID,
        model.PriceMinor,
        model.Currency,
        domain.ProductStatus(model.Status),
        model.CreatedAt.UTC(),
        model.UpdatedAt.UTC(),
    )
}
```

Generated SQLBoiler models remain in `db/sqlboiler/models` and are mapped before reaching the application layer.

---

## 5.4 `repository`

The `repository` folder implements domain repository interfaces.

```text
features/products/repository/
├── postgres_repository.go
├── cache_repository.go
└── transaction_repository.go
```

It is responsible for:

- PostgreSQL queries.
- SQLBoiler operations.
- Redis-backed caching when required.
- Converting database models to domain entities.
- Mapping database errors to stable application errors.
- Respecting the transaction supplied by the application layer.

Example:

```go
type PostgresRepository struct {
    db database.Executor
}

func NewPostgresRepository(db database.Executor) *PostgresRepository {
    return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetByID(
    ctx context.Context,
    id uuid.UUID,
) (*domain.Product, error) {
    model, err := models.FindProduct(ctx, r.db, id.String())
    if err != nil {
        return nil, mapDatabaseError(err)
    }

    product, err := data.ToDomain(model)
    if err != nil {
        return nil, fmt.Errorf("map product %s: %w", id, err)
    }

    return product, nil
}
```

Repositories must not return Fiber responses or request DTOs.

---

## 5.5 Feature `server`

Each feature contains its own HTTP delivery code.

```text
features/products/server/
├── handler.go
├── routes.go
├── request/
└── response/
```

### `handler.go`

Handlers perform only HTTP-related tasks:

1. Read the authenticated principal.
2. Parse path, query, and body data.
3. Validate request DTOs.
4. Convert requests into application commands.
5. Execute use cases.
6. Convert application DTOs into response DTOs.
7. Return the HTTP response through the shared error writer.

Handlers must not:

- Execute SQL.
- Open database transactions.
- Calculate trusted prices.
- Apply state transitions directly.
- Call MinIO directly.
- Send mail directly.
- Create Casbin policies manually.

### `request`

Contains feature-specific input DTOs.

```go
type CreateProductRequest struct {
    ArtisanID  string `json:"artisanId" validate:"required,uuid"`
    CategoryID string `json:"categoryId" validate:"required,uuid"`
    PriceMinor int64  `json:"priceMinor" validate:"required,gt=0"`
    Currency   string `json:"currency" validate:"required,len=3"`
}
```

### `response`

Contains stable API response DTOs.

```go
type ProductResponse struct {
    ID         string `json:"id"`
    ArtisanID  string `json:"artisanId"`
    CategoryID string `json:"categoryId"`
    PriceMinor int64  `json:"priceMinor"`
    Currency   string `json:"currency"`
    Status     string `json:"status"`
}
```

### `routes.go`

Registers feature routes only.

```go
func RegisterRoutes(router fiber.Router, handler *Handler) {
    group := router.Group("/products")

    group.Get("/", handler.List)
    group.Get("/:id", handler.GetByID)
    group.Post("/", handler.Create)
    group.Patch("/:id", handler.Update)
    group.Post("/:id/submit", handler.Submit)
}
```

---

## 6. Suggested Features

```text
features/
├── auth/
├── users/
├── artisans/
├── categories/
├── products/
├── moderation/
├── warehouse/
├── inventory/
├── cart/
├── checkout/
├── orders/
├── payments/
├── shipments/
├── customorders/
├── notifications/
└── admin/
```

### `auth`

Contains:

- Registration.
- Login.
- Refresh tokens.
- Logout.
- Password reset.
- Session management.
- Password hashing.
- Token creation and validation.

### `products`

Contains:

- Product creation and update.
- Product translations.
- Product media metadata.
- Product submission.
- Product readiness rules.
- Public product catalogue queries.

### `admin`

Contains administrative orchestration that does not naturally belong to another feature.

Examples:

- Admin dashboard summaries.
- User role administration.
- System settings.
- Audit log access.
- Global operational reports.

Product approval should remain in `moderation`, warehouse operations in `warehouse`, and payments in `payments`. Do not move every privileged operation into `admin` merely because an administrator executes it.

---

## 7. Global `server` Folder

```text
server/
├── server.go
├── router.go
├── routes.go
├── handlers/
├── middleware/
├── security/
├── request/
├── response/
└── errors/
```

This folder contains shared HTTP server infrastructure. Business-specific handlers stay inside their feature folders.

## 7.1 `server.go`

Creates and starts the Fiber application.

Responsibilities:

- Configure Fiber.
- Set body limits.
- Configure timeouts.
- Register middleware.
- Register API routes.
- Register health endpoints.
- Start and stop the HTTP server.

## 7.2 `router.go`

Creates route groups.

```text
/api/v1
/api/v1/auth
/api/v1/admin
/api/v1/warehouse
/api/v1/webhooks
/health
/ready
```

## 7.3 `routes.go`

Collects each feature route registrar.

```go
func RegisterRoutes(app *fiber.App, modules Modules) {
    api := app.Group("/api/v1")

    authserver.RegisterRoutes(api, modules.AuthHandler)
    productserver.RegisterRoutes(api, modules.ProductHandler)
    orderserver.RegisterRoutes(api, modules.OrderHandler)
    adminserver.RegisterRoutes(api, modules.AdminHandler)
}
```

## 7.4 `middleware`

```text
server/middleware/
├── correlation_id.go
├── request_logger.go
├── recovery.go
├── authentication.go
├── authorization.go
├── rate_limit.go
├── idempotency.go
├── body_limit.go
├── timeout.go
├── cors.go
└── security_headers.go
```

Recommended middleware order:

```text
1. Recovery.
2. Correlation ID.
3. Structured request logging.
4. Security headers.
5. CORS.
6. Body-size limit.
7. Request timeout.
8. Authentication.
9. Rate limiting.
10. Idempotency where required.
11. Route handler.
```

## 7.5 `security`

```text
server/security/
├── principal.go
├── auth_context.go
├── permission.go
├── casbin.go
├── jwt.go
├── password.go
├── webhook_signature.go
└── csrf.go
```

Responsibilities:

- Extract the current principal.
- Verify access and refresh tokens.
- Verify Casbin permissions.
- Verify webhook signatures.
- Hash and compare passwords.
- Attach authentication data to Fiber context.

Middleware may reject unauthenticated requests, but use cases must still enforce ownership and business-level authorisation.

## 7.6 Shared handlers

```text
server/handlers/
├── health.go
├── readiness.go
├── not_found.go
└── method_not_allowed.go
```

Only global technical handlers belong here.

## 7.7 Shared request and response types

```text
server/request/
├── pagination.go
├── idempotency.go
└── upload.go

server/response/
├── envelope.go
├── pagination.go
├── error.go
└── health.go
```

Feature-specific requests and responses stay inside each feature.

---

## 8. `core` Folder

```text
core/
├── container/
├── config/
├── database/
├── cache/
├── storage/
├── mail/
├── logger/
├── security/
├── transaction/
├── validator/
├── clock/
├── pagination/
├── idempotency/
├── audit/
├── outbox/
└── utils/
```

The `core` folder contains shared technical components. It must not become a folder for unrelated business rules.

---

## 8.1 `core/container`

```text
core/container/
├── container.go
├── providers.go
├── modules.go
├── cleanup.go
├── wire.go
└── wire_gen.go
```

Responsibilities:

- Build the dependency graph.
- Initialise shared infrastructure.
- Initialise feature modules.
- Expose the runnable application.
- Register cleanup functions.

Example container:

```go
type Container struct {
    Config   *config.Config
    DB       *sql.DB
    Redis    *redis.Client
    MinIO    *minio.Client
    Mailer   mail.Sender
    Logger   *slog.Logger
    Server   *server.Server
    Jobs     []jobs.Job
}
```

Use Google Wire or explicit constructor injection. Do not use a global service locator.

---

## 8.2 `core/config`

```text
core/config/
├── config.go
├── app.go
├── http.go
├── database.go
├── redis.go
├── minio.go
├── mail.go
├── auth.go
├── payment.go
├── shipping.go
└── loader.go
```

All configuration is loaded from environment variables.

Example:

```go
type Config struct {
    App      AppConfig
    HTTP     HTTPConfig
    Database DatabaseConfig
    Redis    RedisConfig
    MinIO    MinIOConfig
    Mail     MailConfig
    Auth     AuthConfig
    Payment  PaymentConfig
    Shipping ShippingConfig
}
```

Rules:

- Validate configuration before connecting to services.
- Never read environment variables inside handlers or repositories.
- Never log secrets.
- Use typed durations and typed limits.
- Provide safe defaults only for non-sensitive development settings.

---

## 8.3 `core/database`

```text
core/database/
├── postgres.go
├── health.go
├── executor.go
├── transaction.go
└── errors.go
```

Responsibilities:

- Open PostgreSQL connection pool.
- Configure pool limits.
- Ping the database.
- Provide executor abstractions.
- Manage transaction scopes.
- Close the pool during shutdown.

Avoid passing `*sql.Tx` through domain code.

Example abstraction:

```go
type Executor interface {
    ExecContext(context.Context, string, ...any) (sql.Result, error)
    QueryContext(context.Context, string, ...any) (*sql.Rows, error)
    QueryRowContext(context.Context, string, ...any) *sql.Row
}
```

---

## 8.4 `core/cache`

```text
core/cache/
├── redis.go
├── keys.go
├── health.go
├── lock.go
└── rate_limit.go
```

Used for:

- Rate limiting.
- Short-lived caches.
- Distributed locks where justified.
- Session or token metadata.
- Job coordination.

Redis must not be the source of truth for orders, payments, inventory, or audit records.

---

## 8.5 `core/storage`

```text
core/storage/
├── minio.go
├── object_store.go
├── upload.go
├── signed_url.go
├── mime.go
└── keys.go
```

Responsibilities:

- Connect to MinIO.
- Ensure required buckets exist.
- Upload private objects.
- Generate secure object keys.
- Generate signed URLs.
- Validate MIME signatures.
- Delete or quarantine objects when necessary.

Features depend on an interface, not directly on `*minio.Client`.

---

## 8.6 `core/mail`

```text
core/mail/
├── sender.go
├── smtp.go
├── templates.go
├── message.go
└── mock.go
```

Example interface:

```go
type Sender interface {
    Send(ctx context.Context, message Message) error
}
```

Used for:

- Email verification.
- Password reset.
- Order confirmation.
- Payment updates.
- Shipment notifications.

Use cases should usually write an outbox event. A notification worker then sends the email after transaction commit.

---

## 8.7 `core/utils`

```text
core/utils/
├── strings.go
├── slices.go
├── pointer.go
├── hash.go
└── retry.go
```

Only small, generic, stateless helpers belong here.

Do not place business logic such as product price calculation, stock validation, order state transitions, or payment rules in `utils`.

A helper used by only one feature should remain inside that feature.

---

## 9. Database Folder

```text
db/
├── schema/
├── migrations/
├── seeds/
├── runner/
├── sqlboiler/
└── queries/
```

---

## 9.1 `db/schema`

Contains the latest intended database schema for documentation and local database creation.

```text
db/schema/
├── extensions.sql
├── users.sql
├── artisans.sql
├── catalogue.sql
├── inventory.sql
├── orders.sql
├── payments.sql
├── shipments.sql
├── notifications.sql
└── audit.sql
```

The migration history remains the source used to upgrade existing databases.

---

## 9.2 `db/migrations`

```text
db/migrations/
├── 000001_extensions.up.sql
├── 000001_extensions.down.sql
├── 000002_users_sessions.up.sql
├── 000002_users_sessions.down.sql
├── 000003_artisans.up.sql
├── 000003_artisans.down.sql
├── 000004_catalogue.up.sql
├── 000004_catalogue.down.sql
├── 000005_moderation.up.sql
├── 000005_moderation.down.sql
├── 000006_warehouse.up.sql
├── 000006_warehouse.down.sql
├── 000007_inventory.up.sql
├── 000007_inventory.down.sql
├── 000008_cart_addresses.up.sql
├── 000008_cart_addresses.down.sql
├── 000009_orders.up.sql
├── 000009_orders.down.sql
├── 000010_payments.up.sql
├── 000010_payments.down.sql
├── 000011_shipments.up.sql
├── 000011_shipments.down.sql
├── 000012_custom_orders.up.sql
├── 000012_custom_orders.down.sql
├── 000013_notifications_outbox.up.sql
├── 000013_notifications_outbox.down.sql
├── 000014_audit_idempotency.up.sql
└── 000014_audit_idempotency.down.sql
```

Migration rules:

- Every migration has `up` and `down` files.
- Use explicit foreign keys.
- Define deletion behaviour.
- Add unique constraints.
- Add indexes for common lookups.
- Prefer additive and backward-compatible changes.
- Never edit an already deployed migration.
- Create a new migration for every schema change.
- Regenerate SQLBoiler models after migrations change the schema.

---

## 9.3 `db/seeds`

Seeds are organised by environment and table or feature.

```text
db/seeds/
├── reference/
│   ├── 001_roles.sql
│   ├── 002_permissions.sql
│   ├── 003_currencies.sql
│   ├── 004_countries.sql
│   └── 005_categories.sql
├── development/
│   ├── 101_users.sql
│   ├── 102_artisans.sql
│   ├── 103_products.sql
│   └── 104_inventory.sql
├── test/
│   ├── 201_users.sql
│   ├── 202_products.sql
│   ├── 203_orders.sql
│   └── 204_payments.sql
└── seed_manifest.go
```

Seed rules:

- Reference seeds must be idempotent.
- Use stable identifiers for fixed records.
- Development seeds may create demo accounts.
- Test seeds must be deterministic.
- Never insert plain-text passwords.
- Never run development seeds in production.
- Keep one seed file focused on one table or closely related aggregate.

Example idempotent seed:

```sql
INSERT INTO roles (id, code, name)
VALUES
    ('00000000-0000-0000-0000-000000000001', 'CUSTOMER', 'Customer'),
    ('00000000-0000-0000-0000-000000000002', 'ARTISAN', 'Artisan'),
    ('00000000-0000-0000-0000-000000000003', 'ADMIN', 'Administrator')
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name;
```

---

## 9.4 `db/runner`

```text
db/runner/
├── migrate.go
├── seed.go
├── manifest.go
└── lock.go
```

This package is called by the application container during startup.

Responsibilities:

- Run pending migrations.
- Prevent concurrent migration execution.
- Select the correct seed set.
- Track applied seed versions.
- Return contextual errors.

Example:

```go
func Prepare(ctx context.Context, cfg config.DatabaseConfig) error {
    if cfg.MigrateOnStart {
        if err := RunMigrations(ctx, cfg.URL); err != nil {
            return fmt.Errorf("run database migrations: %w", err)
        }
    }

    if cfg.SeedOnStart {
        if err := RunSeeds(ctx, cfg.SeedSet); err != nil {
            return fmt.Errorf("run database seeds: %w", err)
        }
    }

    return nil
}
```

Run migrations before SQLBoiler repositories are used and before the HTTP server accepts traffic.

---

## 9.5 `db/sqlboiler`

```text
db/sqlboiler/
├── models/
├── boil_queries.go
├── boil_table_names.go
└── boil_types.go
```

This folder contains generated code.

Rules:

- Do not add business methods to generated files.
- Do not edit generated files manually.
- Do not return generated models from API handlers.
- Map generated models to domain entities.
- Regenerate models after schema changes.

Suggested command:

```bash
sqlboiler psql --output db/sqlboiler/models --wipe
```

---

## 10. Test Folder

```text
test/
├── unit/
├── integration/
├── concurrency/
├── e2e/
├── fixtures/
├── mocks/
└── testutil/
```

The test folder mirrors the application features.

```text
test/unit/
├── auth/
├── products/
├── inventory/
├── orders/
└── payments/

test/integration/
├── auth/
├── products/
├── inventory/
├── orders/
├── payments/
├── postgres/
├── redis/
└── minio/
```

Go package-local tests may also be placed beside source files when access to unexported functions is useful. The central `test` folder is used for cross-package tests, integration tests, fixtures, and complete workflows.

---

## 10.1 Unit tests

Unit tests must not require PostgreSQL, Redis, MinIO, SMTP, or external network access.

Test:

- Domain state transitions.
- Product readiness rules.
- Order totals.
- Money calculations.
- Inventory movement rules.
- Reservation expiry decisions.
- Payment amount validation.
- DTO mapping.
- Permission rules.

Example:

```text
test/unit/products/product_submit_test.go
test/unit/inventory/reservation_test.go
test/unit/orders/order_total_test.go
test/unit/payments/payment_state_test.go
```

---

## 10.2 Integration tests

Integration tests use real technical dependencies, normally through Docker Compose or Testcontainers.

Test:

- PostgreSQL repositories.
- Migration up and down.
- Seed execution.
- Transaction rollback.
- Redis rate limiting.
- MinIO upload and signed URLs.
- Casbin policies.
- SMTP adapter or mail capture service.
- HTTP routes and middleware.

---

## 10.3 Concurrency tests

```text
test/concurrency/
├── last_unit_reservation_test.go
├── duplicate_checkout_test.go
├── duplicate_webhook_test.go
└── payment_vs_expiry_test.go
```

Test:

- Two customers reserving the last item.
- Duplicate checkout requests.
- Duplicate payment callbacks.
- Reservation expiration during payment confirmation.
- Concurrent warehouse inventory adjustments.

---

## 10.4 E2E tests

```text
test/e2e/
├── customer_purchase_test.go
├── product_to_sale_test.go
├── payment_failure_test.go
├── shipment_flow_test.go
└── custom_order_test.go
```

E2E tests should call the API and verify the complete business result.

---

## 10.5 Test utilities

```text
test/testutil/
├── app.go
├── database.go
├── redis.go
├── minio.go
├── auth.go
├── factories.go
└── cleanup.go
```

Utilities may:

- Start a test application.
- Reset the database.
- Run migrations.
- Load test seeds.
- Generate authenticated principals.
- Create entities through factories.
- Clean storage buckets.

Do not copy production business logic into test helpers.

---

## 11. Jobs Folder

```text
jobs/
├── reservationexpiry/
├── outboxprocessor/
├── notifications/
└── webhookretry/
```

Each job should implement a small interface.

```go
type Job interface {
    Name() string
    Run(ctx context.Context) error
}
```

Examples:

- Expire active inventory reservations.
- Release expired stock.
- Process outbox events.
- Send queued emails.
- Retry temporary webhook failures.

Jobs use application services or dedicated use cases. They must not duplicate business rules.

---

## 12. Feature Module Registration

Each feature exposes a module that groups its dependencies.

```go
type ProductModule struct {
    Handler *productserver.Handler
    Service *productapplication.Service
}

func NewProductModule(
    db database.Executor,
    tx transaction.Manager,
    auth security.Authorizer,
    validator *validator.Validate,
) (*ProductModule, error) {
    repository := productrepository.NewPostgresRepository(db)

    service := productapplication.NewService(
        repository,
        tx,
        auth,
    )

    handler := productserver.NewHandler(service, validator)

    return &ProductModule{
        Handler: handler,
        Service: service,
    }, nil
}
```

For a larger project, Google Wire provider sets may replace manual construction.

---

## 13. Request Execution Flow

Example: create a product.

```text
HTTP request
    ↓
Global Fiber middleware
    ↓
Authentication principal
    ↓
Product feature handler
    ↓
CreateProductRequest validation
    ↓
CreateProductCommand
    ↓
Product application use case
    ↓
Casbin and ownership authorisation
    ↓
Product domain rules
    ↓
Transaction manager
    ↓
Product repository
    ↓
SQLBoiler model and PostgreSQL
    ↓
Audit and outbox records
    ↓
Commit transaction
    ↓
Product application DTO
    ↓
ProductResponse
    ↓
HTTP 201 response
```

No layer may be skipped to place business logic directly in a handler or generated database model.

---

## 14. Dependency Rules

Allowed dependencies:

```text
cmd
  -> core/container

core/container
  -> core infrastructure
  -> server
  -> feature modules
  -> jobs

server
  -> feature server packages
  -> shared core services

feature/server
  -> feature/application
  -> feature request/response DTOs

feature/application
  -> feature/domain
  -> repository interfaces
  -> shared ports

feature/repository
  -> feature/domain
  -> feature/data
  -> db/sqlboiler
  -> core/database

feature/domain
  -> Go standard library only where possible
```

Forbidden dependencies:

```text
feature/domain      -> Fiber
feature/domain      -> SQLBoiler
feature/domain      -> PostgreSQL
feature/domain      -> Redis
feature/domain      -> MinIO
feature/application -> Fiber context
feature/server      -> SQL queries
feature/server      -> SQLBoiler models
core/utils          -> business feature rules
```

---

## 15. Error Handling

Feature domain and application layers return typed errors.

```text
features/products/domain/errors.go
features/orders/domain/errors.go
features/payments/domain/errors.go
```

The global server error writer converts them to stable HTTP responses.

```json
{
  "error": {
    "code": "PRODUCT_NOT_EDITABLE",
    "message": "The product cannot be edited in its current state.",
    "correlationId": "d42dc15a-96c8-45bd-86a4-b7e83ba21a78",
    "details": null
  }
}
```

Rules:

- Wrap technical errors with context.
- Never expose SQL errors or stack traces to clients.
- Log the internal cause with correlation ID.
- Return stable public error codes.
- Map validation errors to field-level details.

---

## 16. Environment File

Example `.env.example`:

```env
APP_ENV=development
APP_NAME=aisha-api
APP_VERSION=0.1.0

HTTP_HOST=0.0.0.0
HTTP_PORT=8080
HTTP_BODY_LIMIT_MB=10
HTTP_READ_TIMEOUT=15s
HTTP_WRITE_TIMEOUT=30s
HTTP_SHUTDOWN_TIMEOUT=15s

DATABASE_URL=postgres://aisha:aisha@postgres:5432/aisha?sslmode=disable
DATABASE_MAX_OPEN_CONNECTIONS=25
DATABASE_MAX_IDLE_CONNECTIONS=10
DATABASE_CONNECTION_MAX_LIFETIME=30m
DB_MIGRATE_ON_START=true
DB_SEED_ON_START=true
DB_SEED_SET=development

REDIS_URL=redis://redis:6379/0

MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=aisha
MINIO_SECRET_KEY=change-me
MINIO_USE_SSL=false
MINIO_PRIVATE_BUCKET=aisha-private
MINIO_PUBLIC_BUCKET=aisha-public

MAIL_DRIVER=smtp
MAIL_HOST=mailpit
MAIL_PORT=1025
MAIL_USERNAME=
MAIL_PASSWORD=
MAIL_FROM_ADDRESS=no-reply@aisha.local
MAIL_FROM_NAME=AISHA

JWT_ACCESS_SECRET=change-me
JWT_REFRESH_SECRET=change-me-too
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=720h

CASBIN_MODEL_PATH=policies/casbin/model.conf
CASBIN_POLICY_PATH=policies/casbin/policy.csv
```

Never commit the real `.env` file.

---

## 17. Required Development Commands

```bash
go mod download
gofmt -w .
go vet ./...
go test ./...
go build ./...
```

Database commands:

```bash
migrate -path db/migrations -database "$DATABASE_URL" up
migrate -path db/migrations -database "$DATABASE_URL" down 1
sqlboiler psql --output db/sqlboiler/models --wipe
```

Suggested Makefile commands:

```bash
make run
make build
make fmt
make vet
make test
make test-unit
make test-integration
make migrate-up
make migrate-down
make seed
make generate-models
```

---

## 18. Initial Implementation Order

Build the backend in this order:

```text
1. Create the root folder structure.
2. Create configuration loading and validation.
3. Create PostgreSQL, Redis, and MinIO clients.
4. Create the migration and seed runner.
5. Create the dependency container.
6. Create the Fiber server and global middleware.
7. Implement shared authentication and security.
8. Implement the auth feature.
9. Implement users and artisans.
10. Implement categories and products.
11. Implement moderation.
12. Implement warehouse and inventory.
13. Implement cart and checkout.
14. Implement orders and payments.
15. Implement shipments.
16. Implement custom orders and notifications.
17. Implement admin operations.
18. Add background jobs.
19. Add OpenAPI documentation.
20. Complete unit, integration, concurrency, and E2E tests.
```

---

## 19. Final Architecture Summary

```text
cmd/main.go
    starts the complete application

features/
    owns all business capabilities

feature/domain/
    owns entities and business rules

feature/application/
    owns commands, queries, and use cases

feature/data/
    maps persistence structures

feature/repository/
    implements database and cache access

feature/server/
    owns handlers, routes, requests, and responses

server/
    owns global Fiber setup, middleware, security, and HTTP errors

core/container/
    assembles all dependencies

core/config/
    loads and validates environment configuration

core/database, cache, storage, mail/
    provide PostgreSQL, Redis, MinIO, and mail adapters

core/utils/
    contains only generic helpers

db/schema/
    documents the current schema

db/migrations/
    contains versioned up/down migrations

db/seeds/
    contains reference, development, and test seed data

db/runner/
    runs migrations and seeds during application startup

test/
    contains unit, integration, concurrency, and E2E tests
```

This structure keeps AISHA modular, testable, and maintainable while allowing each feature to evolve without mixing business rules, HTTP code, and infrastructure code.
