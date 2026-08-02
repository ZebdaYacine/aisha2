# AISHA MVP Architecture

## 1. Architecture Style

Use a modular monolith with Clean Architecture boundaries.

Do not start with microservices. Separate modules in code and data ownership so they may be extracted later if actual scale or organisational need appears.

## 2. Runtime Architecture

```text
Browser
  |
  | HTTPS
  v
Nginx / Traefik
  |
  +------------------+
  |                  |
  v                  v
Next.js Web       Go Fiber API
                     |
        +------------+-------------+
        |            |             |
        v            v             v
   PostgreSQL      Redis         MinIO
        |
        v
   Outbox records
        |
        v
 Background worker
        |
        +--> Email / Push adapter
        +--> Payment adapter
        +--> Shipping adapter
```

## 3. Application Modules

Minimum backend modules:

```text
identity
users
artisans
catalogue
moderation
warehouse
inventory
cart
checkout
payments
orders
shipments
customorders
notifications
audit
files
```

Each module should contain its own application services, interfaces, DTOs, policies, and data implementations.

## 4. Clean Architecture Layers

### Domain

Contains enterprise rules and domain concepts.

Examples:

- Product.
- Artisan.
- Inventory movement.
- Reservation.
- Order.
- Payment.
- Shipment.
- Status transition rules.
- Money value object.

Must not depend on:

- Fiber.
- PostgreSQL driver.
- Redis.
- MinIO SDK.
- SQLBoiler models.
- External providers.

### Application

Contains use cases and orchestration.

Examples:

- `SubmitProductForReview`.
- `InspectReceptionBatch`.
- `ReserveInventory`.
- `CreateCheckout`.
- `ProcessPaymentWebhook`.
- `PrepareOrder`.
- `CreateShipment`.

Responsibilities:

- Validate workflow transitions.
- Open explicit transactions.
- Invoke repository interfaces.
- Enforce domain rules.
- Invoke policy/authorisation service.
- Create outbox events.

### Interface / Delivery

Contains:

- Fiber handlers.
- Request DTOs.
- Response DTOs.
- Middleware.
- OpenAPI mapping.

Handlers must:

1. Parse request.
2. Validate request DTO.
3. Read authenticated principal.
4. Call one application use case.
5. Map domain/application errors to API errors.
6. Return response DTO.

Handlers must not query PostgreSQL directly.

### Infrastructure

Contains:

- SQLBoiler repository implementations.
- PostgreSQL transactions.
- Redis adapters.
- MinIO adapter.
- Casbin adapter.
- Payment provider.
- Shipping provider.
- Email/push provider.
- Logger.
- Clock and ID generators.

## 5. Dependency Direction

```text
delivery -> application -> domain
infrastructure -> application/domain interfaces
```

The domain never imports infrastructure.

Application interfaces are implemented by infrastructure.

Wire assembles all dependencies in the composition root.

## 6. Backend Folder Layout

`STRUCTURE.md` is the canonical detailed backend layout. Implemented packages are
organised by responsibility without creating empty placeholders for future MVP work.

```text
backend/
├── cmd/
│   └── main.go
├── features/
│   ├── auth/
│   ├── artisans/
│   ├── users/
│   └── catalogue/
├── server/
├── core/
│   ├── cache/
│   ├── config/
│   ├── database/
│   ├── health/
│   ├── security/
│   └── storage/
├── db/
│   ├── migrations/
│   ├── seeds/
│   ├── schema/
│   └── sqlboiler/models/
├── openapi/
├── test/integration/
├── tools/
├── sqlboiler.toml
├── go.mod
└── Dockerfile
```

## 7. Frontend Architecture

```text
frontend/
├── app/
│   ├── [locale]/
│   │   ├── (public)/
│   │   ├── (customer)/
│   │   ├── artisan/
│   │   ├── admin/
│   │   └── layout.tsx
│   └── api/
├── components/
│   ├── ui/
│   └── shared/
├── features/
│   ├── auth/
│   ├── catalogue/
│   ├── cart/
│   ├── checkout/
│   ├── artisan/
│   ├── warehouse/
│   └── admin/
├── lib/
│   ├── api/
│   ├── auth/
│   ├── i18n/
│   ├── validation/
│   └── utils/
├── messages/
│   ├── ar.json
│   ├── fr.json
│   ├── en.json
│   └── es.json
└── tests/
```

Server Components are the default. Client Components are used only for interactivity, browser APIs, local state, or client-only libraries.

## 8. Data Architecture

Main tables:

```text
users
sessions
roles / casbin_rules
artisan_profiles
artisan_documents
categories
products
product_translations
product_media
product_submissions
warehouse_receptions
warehouse_reception_items
warehouse_inspections
inventory_movements
inventory_balances
inventory_reservations
carts
cart_items
addresses
orders
order_items
payments
payment_events
shipments
shipment_events
custom_order_requests
custom_order_messages
custom_order_quotes
notifications
outbox_events
audit_events
idempotency_keys
```

## 9. Transaction Boundaries

Transactions are mandatory for:

- Artisan approval plus role assignment.
- Product moderation state update plus audit/outbox event.
- Inspection plus inventory movements.
- Reservation creation.
- Checkout order plus reservations plus payment attempt.
- Payment confirmation plus order update plus reservation consumption.
- Payment failure plus reservation release.
- Cancellation plus inventory and payment consequences.
- Shipment state update plus order state update.

External provider calls should generally occur outside the core database transaction. Store an intent/outbox record, call the provider safely, and reconcile with idempotency.

## 10. Outbox Pattern

For critical side effects:

1. Update business state.
2. Insert `outbox_events` in the same transaction.
3. Commit.
4. Worker claims event.
5. Worker sends notification or provider request.
6. Mark event processed.
7. Retry failures with backoff.

Required fields:

- ID.
- Event type.
- Aggregate type and ID.
- Payload.
- Created at.
- Available at.
- Attempt count.
- Locked at.
- Processed at.
- Last error.

## 11. Caching

Redis may cache:

- Public product details.
- Category lists.
- Artisan public profiles.
- Rate-limit counters.
- Session metadata if the selected session design requires it.
- Short-lived idempotency or lock data.

Redis must not be the sole source of:

- Orders.
- Payments.
- Inventory.
- Artisan approval.
- Product moderation.
- Audit events.

## 12. API Style

- Base path: `/api/v1`.
- JSON request and response DTOs.
- Stable machine-readable errors.
- Cursor or page-based pagination used consistently.
- Correlation ID on every request.
- OpenAPI generated or maintained with implementation.
- Idempotency key header on checkout, payment, refund, and shipment creation.

Standard error:

```json
{
  "error": {
    "code": "OUT_OF_STOCK",
    "message": "The requested quantity is no longer available.",
    "details": {
      "productId": "..."
    },
    "requestId": "..."
  }
}
```

## 13. Scalability Path

Keep these interfaces extractable:

- Payment provider.
- Shipping provider.
- Notification provider.
- Search.
- File storage.
- Inventory repository.

Do not introduce a queue, search cluster, or microservice only because it may be useful later. The MVP uses the simplest reliable infrastructure.
