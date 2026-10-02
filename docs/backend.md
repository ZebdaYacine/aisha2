# AISHA MVP Backend Development Guide

## 1. Goal

Build a Go Fiber API that implements the business workflow without placing business logic in handlers or SQLBoiler models.

## 2. Required Packages

Suggested core packages:

```text
github.com/gofiber/fiber/v3
github.com/casbin/casbin/v2
github.com/volatiletech/sqlboiler/v4
github.com/google/wire
github.com/jackc/pgx/v5
github.com/redis/go-redis/v9
github.com/minio/minio-go/v7
github.com/golang-migrate/migrate/v4
github.com/google/uuid
github.com/go-playground/validator/v10
```

Choose exact versions compatible with the repository and current Go version. Do not mix multiple HTTP routers or ORMs.

## 3. Coding Rules

- Run `gofmt`.
- Handle every error.
- Wrap errors with context.
- Use `context.Context` in application and repository methods.
- Keep interfaces small and consumer-owned.
- Use constructor injection.
- Avoid global mutable state.
- Use UTC timestamps.
- Use structured logs with correlation IDs.
- Use integer minor units for money.
- Never expose SQLBoiler models as API responses.

## 4. HTTP Handler Pattern

Example responsibility:

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

    result, err := h.createProduct.Execute(
        c.Context(),
        principal,
        req.ToCommand(),
    )
    if err != nil {
        return h.errors.Write(c, err)
    }

    return c.Status(fiber.StatusCreated).JSON(ProductResponseFrom(result))
}
```

Handlers do not:

- Open SQL queries.
- Decide product state transitions.
- Calculate trusted prices.
- Implement Casbin policy manually.
- Call MinIO directly.

## 5. Application Use Case Pattern

```go
type CreateProductCommand struct {
    ArtisanID string
    CategoryID string
    PriceMinor int64
    Currency string
    Translations []ProductTranslationInput
}

type CreateProduct interface {
    Execute(ctx context.Context, actor Principal, cmd CreateProductCommand) (ProductDTO, error)
}
```

Use case responsibilities:

1. Authorise action.
2. Validate command and domain invariants.
3. Load required data.
4. Start transaction where needed.
5. Call repositories.
6. Insert audit/outbox event.
7. Commit.
8. Return application DTO.

## 6. Repository Pattern

Domain interface:

```go
type ProductRepository interface {
    GetByID(ctx context.Context, id string) (*Product, error)
    Create(ctx context.Context, product *Product) error
    Update(ctx context.Context, product *Product) error
}
```

Transaction-aware ports may receive a `Tx` abstraction or repositories created from a transaction scope.

Do not leak `*sql.Tx` throughout domain code.

## 7. Domain Rules

### Product

A product can be submitted only when:

- Artisan is approved and active.
- Required translation fields exist.
- Category exists.
- Price is positive.
- Currency is supported by configuration.
- Required media exists.
- Current status permits submission.

A product can be active only when:

- Moderation is approved.
- Product is not suspended.
- Accepted inventory is available.

### Inventory

Use append-only movements:

```text
RECEPTION
INSPECTION_ACCEPT
INSPECTION_REJECT
INSPECTION_QUARANTINE
INSPECTION_DAMAGE
RESERVATION
RESERVATION_RELEASE
ORDER_COMMIT
SHIPMENT
ADJUSTMENT_IN
ADJUSTMENT_OUT
RETURN
```

Maintain an optimised `inventory_balances` table if required, but every balance change must be explainable by movements.

### Reservation

Required fields:

- ID.
- Product or stock item.
- Order or checkout reference.
- Quantity.
- Status.
- Expires at.
- Created at.
- Consumed/released at.

A background worker expires old active reservations; checkout also performs a locked expiry sweep as a safety net.

### Order

Order lines are snapshots. Do not join current product price to show historical order totals.

### Payment

Store:

- Internal payment ID.
- Order ID.
- Provider.
- Provider payment reference.
- Amount minor.
- Currency.
- Status.
- Idempotency key.
- Created and updated timestamps.

Store provider events separately with a unique external event ID.

## 8. Database Migrations

Migration rules:

- Every schema change has `up` and `down`.
- Prefer additive changes.
- Add explicit foreign keys.
- Add unique constraints.
- Add indexes for frequent lookups.
- Define deletion behaviour.
- Production migrations never run invisibly inside API startup.
- Regenerate SQLBoiler models after schema changes.

Suggested migration sequence:

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
000018_cart_wishlist
000023_notifications
000024_commerce_fulfilment
000025_email_verification_tokens
000026_notification_email_delivery
```

## 9. API Endpoints

### Authentication

```text
POST   /api/v1/auth/register
GET    /api/v1/auth/activate?token=<single-use-token>
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password
GET    /api/v1/me
```

Registration stores only a hash of the single-use email-verification token and
publishes an outbox email containing the configured `WEB_BASE_URL` activation
link. The API defaults that base URL to `http://localhost:3033` in development
and `https://aishasouk.com` in production when it is not explicitly configured.
Activation links expire after 120 seconds. The activation endpoint consumes the
token transactionally and marks the email as verified; it returns an activation
result without creating a session. The web page then shows the activated state
and a Sign in button.

### Notifications

```text
GET    /api/v1/notifications
POST   /api/v1/notifications/{id}/read
POST   /api/v1/notifications/read-all
GET    /api/v1/notifications/ws-ticket
GET    /api/v1/notifications/ws?ticket=<one-time-ticket>
```

The API worker claims pending outbox records with a lease, creates recipient-scoped
notifications idempotently, retries persistence failures with backoff, and publishes
committed notifications to connected users over WebSocket. The browser receives a one-minute,
one-time ticket through the authenticated BFF so HttpOnly access cookies stay private.
Workflow events also fan out to the relevant operational peers: moderators and
administrators for review queues, warehouse agents for stock transitions, and the
customer/artisan/fulfilment participants of an order.

When SMTP is configured, the same committed workflow notifications are sent to each
recipient's account email through the transactional mail adapter. Registration,
artisan application and moderation, payment, warehouse, inventory, order, account,
and future workflow events use the same outbox path. Hostinger's implicit TLS mode
is used for port 465. Each notification has a durable email-delivery state and is
claimed once, so SMTP failures are recorded without retrying the entire outbox event
or resending messages already delivered. Email bodies are responsive HTML messages with inline-safe CSS,
event summaries, and an allowlisted set of business details. Internal identifiers
such as user, payment, product, workshop, actor, and aggregate IDs are excluded
from email content; they remain available only to protected internal workflows.

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
POST   /api/v1/artisan-applications/draft
POST   /api/v1/artisan-applications/me/submit
GET    /api/v1/artisan-applications/me
PATCH  /api/v1/artisan/profile
POST   /api/v1/artisan/products
PATCH  /api/v1/artisan/products/:id
POST   /api/v1/artisan/products/:id/media
GET    /api/v1/artisan/profile/media
POST   /api/v1/artisan/profile/media
PATCH  /api/v1/artisan/profile/media/:id
DELETE /api/v1/artisan/profile/media/:id
GET    /api/v1/artisan-applications/me/documents
POST   /api/v1/artisan-applications/me/documents
POST   /api/v1/artisan/products/:id/submit
GET    /api/v1/artisan/products
GET    /api/v1/artisan/inventory
```

An existing customer starts artisan onboarding by saving a draft. The draft
contains the proposed workshop details and profile information. The customer
then uploads at least one private application document and one private profile
media item. The final submit endpoint checks those server-side requirements
before moving the application to `SUBMITTED`; a browser cannot bypass that
check by calling the endpoint directly.

### Moderation

```text
GET    /api/v1/admin/product-submissions
POST   /api/v1/admin/product-submissions/:id/approve
POST   /api/v1/admin/product-submissions/:id/request-changes
POST   /api/v1/admin/product-submissions/:id/reject
POST   /api/v1/admin/products/:id/suspend
GET    /api/v1/admin/media/users
GET    /api/v1/admin/media/products
GET    /api/v1/admin/artisan-applications/:id/documents
GET    /api/v1/admin/artisan-applications/:id/media
```

### Warehouse and Inventory

```text
POST   /api/v1/warehouse/receptions
GET    /api/v1/warehouse/receptions
POST   /api/v1/warehouse/receptions/:id/inspect
POST   /api/v1/warehouse/receptions/:id/evidence
GET    /api/v1/warehouse/receptions/:id/evidence
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
POST   /api/v1/cart/merge
PATCH  /api/v1/cart/items/:productId
DELETE /api/v1/cart/items/:productId
GET    /api/v1/wishlist
POST   /api/v1/wishlist/items/:productId
DELETE /api/v1/wishlist/items/:productId
POST   /api/v1/checkout
GET    /api/v1/orders
GET    /api/v1/orders/:id
POST   /api/v1/orders/:id/cancel
GET    /api/v1/payments/:paymentId
POST   /api/v1/payments/:paymentId/confirm
```

### Payment and Shipping Callbacks

```text
POST   /api/v1/payments/:paymentId/retry
POST   /api/v1/webhooks/payments/:provider
POST   /api/v1/webhooks/shipments/:provider

The current MVP uses an internal manual development adapter. Payment
confirmation never accepts card data and only succeeds when the locked
server-side payment amount and order total match. Confirmation commits held
reservations with their committed timestamp, marks the order PAID, records a
manual provider reference, and emits a PAYMENT_CONFIRMED outbox event for the
customer plus warehouse/admin recipients. The notification worker delivers
the in-app event and configured SMTP email; a real provider/webhook adapter
may be added later without changing the order contract.
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

## 10. Validation

Validate at two levels:

### Request DTO

- Required fields.
- String length.
- Numeric range.
- UUID format.
- Supported file metadata.
- ISO currency format.
- Pagination bounds.

### Domain/Application

- Actor permission.
- Resource ownership.
- Valid status transition.
- Stock availability.
- Product readiness.
- Amount consistency.
- Duplicate callback.
- Idempotency.

Do not rely on frontend validation.

## 11. Error Catalogue

Suggested stable codes:

```text
VALIDATION_ERROR
INVALID_CREDENTIALS
AUTHENTICATION_REQUIRED
FORBIDDEN
RESOURCE_NOT_FOUND
CONFLICT
ARTISAN_NOT_APPROVED
PRODUCT_NOT_EDITABLE
PRODUCT_NOT_SELLABLE
INVALID_STATE_TRANSITION
OUT_OF_STOCK
RESERVATION_EXPIRED
IDEMPOTENCY_CONFLICT
PAYMENT_AMOUNT_MISMATCH
INVALID_WEBHOOK_SIGNATURE
DUPLICATE_WEBHOOK
UNSUPPORTED_FILE_TYPE
FILE_TOO_LARGE
RATE_LIMITED
INTERNAL_ERROR
```

## 12. Idempotency

For checkout:

1. Client sends `Idempotency-Key`.
2. Store key, actor, route, request hash, status, and response.
3. Same key plus same request returns saved response.
4. Same key plus different request returns `IDEMPOTENCY_CONFLICT`.

Apply similar protection to:

- Payment creation.
- Refund.
- Shipment creation.
- External webhook event processing.

## 13. File Upload

Upload workflow:

1. Validate request size before full processing.
2. Validate MIME using content signature.
3. Generate object key.
4. Upload to private bucket.
5. Store metadata.
6. Scan if scanner is available.
7. Create signed URL only for authorised read.
8. Publish selected product media only after moderation.

Never trust the original filename as the storage path.

The API uses `MINIO_ENDPOINT` for storage operations and `MINIO_PUBLIC_ENDPOINT`
for signed URLs returned to browsers. The presigner is configured with MinIO's
default `us-east-1` region so generating a URL does not require a bucket-location
request through the browser-facing endpoint.

## 14. OpenAPI

Every route must document:

- Authentication.
- Casbin permission.
- Request schema.
- Response schema.
- Error codes.
- Pagination.
- Idempotency header.
- Example payload.

OpenAPI must be validated in CI.

## 15. Testing

### Unit

- Domain state transitions.
- Money calculations.
- Inventory invariants.
- Order totals.
- Permission policies.
- DTO mapping.

### Integration

- PostgreSQL repositories.
- Transaction rollback.
- Casbin enforcement.
- MinIO file storage.
- Redis rate limiting.
- Migration up/down.

### Concurrency

- Last-unit reservation.
- Duplicate checkout.
- Duplicate webhook.
- Reservation expiry versus payment confirmation.

### E2E

- Product to warehouse to sale.
- Customer purchase.
- Payment failure.
- Shipment flow.

## 16. Required Commands

```bash
go mod download
gofmt -w .
go vet ./...
go test ./...
go build ./...
migrate -path migrations -database "$DATABASE_URL" up
```
