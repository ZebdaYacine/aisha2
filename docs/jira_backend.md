# AISHA Backend Milestone Tracker

This tracker follows `docs/backend.md`. Complete milestones in order and do not start a later milestone until the current milestone meets its definition of done.

## Status Legend

| Status | Meaning |
|---|---|
| ✅ Implemented | Code and relevant tests are present in the repository. |
| 🟡 Partial | Some required work exists, but the milestone is not complete. |
| ⬜ Not started | No complete implementation is present. |

## Milestones

| Order | Jira Epic | Milestone | Status | Implemented now | Remaining scope / definition of done | Depends on |
|---:|---|---|---|---|---|---|
| 1 | `BE-EPIC-001` | Platform foundation | ✅ Implemented | Go Fiber API, configuration, Wire composition, PostgreSQL pool, Redis and MinIO clients, health endpoints, structured request logging, Docker foundation, SQLBoiler configuration | Keep health and dependency checks passing; do not add business logic to handlers or generated models | — |
| 2 | `BE-EPIC-002` | Standard API errors and DTOs | ✅ Implemented | Stable API error envelope, error codes, field-error support, health DTO mapping, and unit tests | Reuse the same envelope and DTO boundary in every later endpoint | Platform foundation |
| 3 | `BE-EPIC-003` | Database foundation | ✅ Implemented | Seven paired migrations; users, roles, sessions, password resets, artisan profiles, catalogue, media, audit, outbox, and idempotency tables; reference seeds; generated SQLBoiler models; migration structure tests | Run the PostgreSQL up/down integration test with `TEST_DATABASE_URL` in CI; update generated models after every schema change | Platform foundation |
| 4 | `BE-EPIC-004` | Authentication and Casbin authorization | ✅ Implemented | Register, login, refresh, logout, forgot/reset password, `/me`, hashed passwords and tokens, session rotation, field-level DTO validation, default Casbin policies and protected-route enforcement, Redis-backed rate limits, PostgreSQL/Redis integration tests, route and policy tests, and validated OpenAPI documentation | Keep later protected routes behind authentication and explicit Casbin permissions; run conditional integration tests with `TEST_DATABASE_URL` and `TEST_REDIS_URL` in CI | API errors, database foundation |
| 5 | `BE-EPIC-005` | Public catalogue | ✅ Implemented | Category, product, product-detail, artisan, and artisan-detail endpoints expose localized explicit DTOs with active/published/approved/public-media rules, bounded pagination, filters, OpenAPI, service/handler tests, and PostgreSQL publication-rule coverage | Inventory availability remains in the inventory milestone | Authentication completion, database foundation |
| 6 | `BE-EPIC-006` | Artisan applications and profiles | 🟡 Partial | Application submission/status/resubmission, approved profile updates, administrator review decisions, category/translations, role assignment, audit/outbox events, private document metadata reads, Casbin/ownership, DTOs, service tests, migration, and OpenAPI are implemented | Add secure MinIO upload/signed document reads and PostgreSQL/HTTP integration coverage | Authentication, public catalogue |
| 7 | `BE-EPIC-007` | Artisan products and media | ⬜ Not started | Product, translation, and media schema exists; MinIO client exists | Implement create/edit/list/submit product workflows and secure media upload; validate MIME signature, size, object key, translations, price, currency, category, ownership, and product state; use transactions, audit/outbox events, signed URLs, tests, and OpenAPI | Artisan profiles |
| 8 | `BE-EPIC-008` | Product moderation | ⬜ Not started | Product status column supports moderation-related states | Add moderation migration changes if required; implement submission queue, approve, request changes, reject, and suspend; enforce Casbin; keep transitions in domain/application code; write audit events and transition/authorization tests | Artisan products |
| 9 | `BE-EPIC-009` | Warehouse and inventory | ⬜ Not started | — | Add warehouse, reception, inspection, append-only inventory movement, balance, and reservation migrations; implement reception, inspection, adjustment, inventory query, reservation expiry, and warehouse order endpoints; use explicit transactions and locking; add invariant and concurrency tests | Moderation |
| 10 | `BE-EPIC-010` | Customer profile and addresses | ✅ Implemented | Paired migration and generated model, authenticated profile/address CRUD, explicit DTO validation, principal-derived ownership, Casbin enforcement, transactional single-default changes, service/HTTP tests, and PostgreSQL ownership/default-invariant coverage are implemented | Addresses are consumed by checkout in its later milestone | Authentication |
| 11 | `BE-EPIC-011` | Cart and checkout | ⬜ Not started | Idempotency table exists | Add cart and checkout schema; implement cart item operations and checkout; trust server-side prices; reserve stock transactionally; persist idempotent responses and detect conflicting request hashes; add rollback and duplicate-checkout tests | Inventory, customer addresses |
| 12 | `BE-EPIC-012` | Orders | ⬜ Not started | — | Add order and immutable order-line snapshot migrations; implement customer order list/detail/cancellation and warehouse prepare/ship flows; protect state changes with transactions and Casbin; test totals, transitions, ownership, and rollback | Checkout |
| 13 | `BE-EPIC-013` | Payments | ⬜ Not started | — | After a provider is explicitly selected, add payments and provider-event schema; implement status/retry and signed webhook processing; verify amounts, unique external event IDs, and idempotency; add duplicate webhook, failure, and concurrency tests | Orders |
| 14 | `BE-EPIC-014` | Shipments | ⬜ Not started | — | After shipping integration details are explicitly defined, add shipment schema, transitions, callbacks, authorization, idempotency, audit/outbox events, tests, and OpenAPI | Orders; payment rules where applicable |
| 15 | `BE-EPIC-015` | Custom orders | ⬜ Not started | — | Add custom-order, message, quote, and quote-acceptance schema and workflows; enforce ownership, Casbin, valid transitions, immutable accepted quote amounts, transactions, tests, and OpenAPI | Authentication, artisan profiles, orders/payments as required |
| 16 | `BE-EPIC-016` | Notifications and outbox processing | 🟡 Partial | Outbox table and authentication password-reset notifier port exist | Implement reliable outbox claiming, retry/backoff, delivery status, worker lifecycle, notification adapters selected by configuration, idempotent consumers, and integration tests | Database foundation; consuming workflows |
| 17 | `BE-EPIC-017` | Backend release hardening | 🟡 Partial | Unit tests, migration checks, health checks, and basic OpenAPI exist | Validate OpenAPI in CI; run PostgreSQL, Redis, MinIO, Casbin, migration, transaction, concurrency, and E2E suites; verify `gofmt`, `go vet ./...`, `go test ./...`, and `go build ./...`; document deployment migrations and rollback | All MVP milestones |

## Current Sprint Recommendation

| Priority | Jira Issue | Work item | Status | Acceptance criteria |
|---:|---|---|---|---|
| 1 | `BE-401` | Complete request DTO validation | ✅ Implemented | All authentication request fields use consistent validator rules and return the standard field-error envelope; unit and handler tests pass |
| 2 | `BE-402` | Integrate Casbin with protected routes | ✅ Implemented | Authenticated routes invoke Casbin through application/authorization code; allowed and forbidden integration tests pass |
| 3 | `BE-403` | Add authentication rate limiting | ✅ Implemented | Login, registration, and password-reset endpoints use Redis-backed limits and return `RATE_LIMITED`; integration tests pass |
| 4 | `BE-404` | Complete authentication persistence tests | ✅ Implemented | Registration, login, refresh rotation, logout, password reset, and session revocation are covered by a PostgreSQL integration test |
| 5 | `BE-405` | Complete authentication OpenAPI | ✅ Implemented | Authentication, permissions, schemas, stable errors, and examples are documented and validated |

After `BE-EPIC-004` is complete, start `BE-EPIC-005` Public Catalogue.

## Milestone Completion Checklist

Apply this checklist to every milestone:

- [ ] Fiber handlers contain no business rules or SQL queries.
- [ ] Application/domain code owns authorization, invariants, and state transitions.
- [ ] Casbin protects every non-public action.
- [ ] SQLBoiler models remain inside infrastructure repositories.
- [ ] API responses use explicit DTOs and the standard error envelope.
- [ ] Money uses integer minor units.
- [ ] Critical multi-record/state operations use explicit transactions.
- [ ] Every migration has paired up/down files, constraints, indexes, and deletion behavior.
- [ ] SQLBoiler models are regenerated after schema changes.
- [ ] Unit and relevant integration/concurrency tests are present.
- [ ] OpenAPI documents routes, authentication, permission, pagination, idempotency, errors, and examples.
- [ ] `gofmt`, `go vet ./...`, `go test ./...`, and `go build ./...` pass before marking implemented.
