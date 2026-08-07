# AISHA MVP Jira Project Tracker

Last reviewed: 2026-08-07

This is the implementation tracker for the AISHA MVP. It is based on
`docs/requirements.md`, `docs/userstory.md`, `docs/architecture.md`,
`docs/backend.md`, `docs/frontend.md`, `docs/security.md`, and
`docs/deployment.md`.

The status reflects the repository as it exists today. A package directory,
database table, mock screen, or empty module scaffold is not counted as a
completed feature until the required workflow is connected end to end.

## Status Legend

| Status | Meaning |
|---|---|
| ✅ Implemented | The scoped API/UI workflow exists and has relevant automated validation. |
| 🟡 Partial | Some schema, API, UI, or infrastructure exists, but the required workflow is incomplete. |
| ⬜ Not started | Only documentation, a migration, a placeholder module, or no implementation exists. |
| ⛔ Blocked | Completion requires an unconfirmed provider, business rule, credential, or operational decision. |

## Current System Snapshot

Implemented and connected today:

- Go Fiber API foundation with health, error handling, validation, logging, and request IDs.
- PostgreSQL migrations, reference categories, development seed data, Redis, MinIO, Nginx, and Docker Compose.
- Customer registration, login, refresh, logout, password reset, `/me`, profile, and address APIs.
- Casbin authorization and authenticated/admin artisan application routes.
- Public localized category, product, and artisan catalogue APIs with pagination and publication filtering.
- Next.js public storefront, search, filters, product/category/artisan pages, authentication pages, account pages, and artisan/admin application screens.
- English, French, Arabic, and Spanish catalogue translations; Arabic RTL support.
- Client-side cart and development deployment/rebuild script.

Not yet connected to the API:

- Product authoring, product media upload, product moderation, warehouse, inventory, reservations, server cart, checkout, orders, payments, shipments, custom orders, reviews, and wishlist persistence.
- The corresponding backend feature directories exist as composition scaffolds in many cases; they do not represent completed functionality.

## Epic Tracker

| Order | Jira Epic | Epic | Requirements and stories | Status | Present in the system | Remaining definition of done | Depends on |
|---:|---|---|---|---|---|---|---|
| 1 | `PRJ-EPIC-001` | Product definition and architecture | Product vision, roles, MVP scope, architecture, security, UI, and deployment documentation | ✅ Implemented | The required planning and architecture documents exist and define the business workflow, layers, state machines, and non-functional rules | Keep this tracker and contracts synchronized as implementation changes | — |
| 2 | `PRJ-EPIC-002` | Platform and development foundation | Deployment baseline; local development; CI foundation | ✅ Implemented | Go API, Next.js frontend, Compose services, health checks, Makefile, Jenkinsfile, migrations, MinIO buckets, and `script.sh` are present | Add the worker runtime and complete staging/production operational controls under `PRJ-EPIC-022` | `PRJ-EPIC-001` |
| 3 | `PRJ-EPIC-003` | Design system, multilingual UI, RTL, and responsive foundation | `REQ-I18N-001`; `US-VIS-001`; accessibility and frontend guide | ✅ Implemented | Shared UI primitives, responsive storefront, English/French/Arabic/Spanish copy, locale routing, RTL handling, loading/error/empty states, and frontend tests exist | Apply the same quality bar to future operational screens and complete the full accessibility audit under `PRJ-EPIC-021` | `PRJ-EPIC-001` |
| 4 | `PRJ-EPIC-004` | Authentication and authorization | `REQ-AUTH-001`, `REQ-AUTH-002`, `REQ-RBAC-001`; `US-AUTH-001..003` | ✅ Implemented | Registration, login, refresh rotation, logout, password reset, HTTP-only session handling, rate limiting, Casbin checks, stable errors, and tests exist | Execute PostgreSQL/Redis integration tests in CI and extend authorization tests to every later protected workflow | `PRJ-EPIC-002` |
| 5 | `PRJ-EPIC-005` | Customer profile and addresses | Customer role stories; address ownership rules | ✅ Implemented | Profile and address API/BFF routes, principal-derived ownership, single-default address transaction, forms, and tests exist | Reuse the real address API from checkout and order creation | `PRJ-EPIC-004` |
| 6 | `PRJ-EPIC-006` | Public catalogue and discovery | `US-VIS-001`, `US-VIS-002`; `REQ-PROD-002`; public catalogue requirements | ✅ Implemented | Categories, products, artisans, localized names/content, public media, pagination, category filtering, search, real API adapter, and storefront integration exist | Add authoritative stock/availability from `PRJ-EPIC-011` without inventing availability in the frontend | `PRJ-EPIC-003`, `PRJ-EPIC-004` |
| 7 | `PRJ-EPIC-007` | Artisan onboarding and administration | `REQ-ART-001`; `US-ART-001`, `US-ART-002`, `US-ADMIN-001..003` | 🟡 Partial | Artisan application submission/resubmission, multilingual biography, category selection, approved profile edit, admin review/decision routes, reasons, role assignment, audit/outbox events, and UI exist | Add secure document/profile-media upload and authorized reads; complete user/role administration and audit query UI | `PRJ-EPIC-004`, `PRJ-EPIC-020` |
| 8 | `PRJ-EPIC-008` | Artisan product authoring and media | `REQ-PROD-001`; `US-PROD-001`, `US-PROD-002` | ⬜ Not started | Product, translation, and media tables plus public product presentation exist; the product backend module is still a scaffold | Implement artisan-owned draft create/edit/list/submit APIs, four-locale validation, price/category rules, media uploads, editor UI, ownership tests, and OpenAPI | `PRJ-EPIC-007`, `PRJ-EPIC-020` |
| 9 | `PRJ-EPIC-009` | Product moderation | `REQ-MOD-001`; `US-PROD-003`, `US-PROD-004` | ⬜ Not started | Product statuses are represented in the schema | Implement submission snapshots, moderator queue, approve/request-changes/reject/suspend transitions, mandatory reasons, Casbin policies, audit/outbox events, notifications, UI, and transition tests | `PRJ-EPIC-008` |
| 10 | `PRJ-EPIC-010` | Warehouse reception and inspection | `REQ-WH-001`, `REQ-WH-002`; `US-WH-001`, `US-WH-002` | ⬜ Not started | Warehouse and inspection modules are scaffolds; no operational schema or routes are connected | Implement reception batches, inspection evidence, accepted/rejected/quarantined/damaged quantities, quantity invariant, transactions, roles, UI, tests, and OpenAPI | `PRJ-EPIC-009`, `PRJ-EPIC-020` |
| 11 | `PRJ-EPIC-011` | Inventory ledger and reservations | `REQ-INV-001`; `US-WH-003`, `US-CHECK-002` | ⬜ Not started | Inventory and reservation module scaffolds exist only | Implement append-only movements, balances, adjustments, reservation expiry, accepted-stock activation, locking, background expiry, warehouse/customer availability, and last-unit concurrency tests | `PRJ-EPIC-010` |
| 12 | `PRJ-EPIC-012` | Cart and wishlist persistence | `REQ-CART-001`; `US-CART-001`; architecture wishlist feature | 🟡 Partial | Client cart context, cart page/drawer, quantity controls, persistence, wishlist page, and frontend tests exist | Add server/anonymous cart persistence, merge-on-login, active-product revalidation, backend totals, real availability errors, wishlist API, ownership checks, and E2E tests | `PRJ-EPIC-004`, `PRJ-EPIC-006`, `PRJ-EPIC-011` |
| 13 | `PRJ-EPIC-013` | Checkout, orders, and cancellation | `REQ-CHECKOUT-001`, `REQ-ORDER-001`, `REQ-ORDER-002`; `US-CHECK-001`, `US-CHECK-002`, `US-ORD-001..003` | 🟡 Partial | Checkout and customer order/tracking screens exist as frontend surfaces; backend modules are scaffolds and no checkout/order API is registered | Implement server price reload, address selection, transactional reservations, pending orders, immutable line snapshots, idempotency, state transitions, ownership, cancellation, real account integration, and concurrency tests | `PRJ-EPIC-005`, `PRJ-EPIC-011`, `PRJ-EPIC-012` |
| 14 | `PRJ-EPIC-014` | Reviews and customer engagement | Architecture review/wishlist features; post-purchase extension | ⬜ Not started | Review backend package is a scaffold; no review API or persistent review UI is connected | Confirm MVP inclusion, then implement verified-purchase rules, moderation, ownership, rating aggregation, localized UI, and tests | `PRJ-EPIC-013` |
| 15 | `PRJ-EPIC-015` | Payment abstraction and payment confirmation | `REQ-PAY-001`; `US-PAY-001..004` | ⛔ Blocked | Payment requirements and module scaffold exist; no provider adapter or webhook endpoint is implemented | Confirm provider, currencies, payment states, retry/refund rules, and webhook contract; then implement trusted totals, signed/idempotent webhooks, payment events, failure/expiry release, UI, and tests | `PRJ-EPIC-011`, `PRJ-EPIC-013`; provider decision |
| 16 | `PRJ-EPIC-016` | Fulfilment and shipments | `REQ-SHIP-001`; `US-FUL-001`, `US-SHIP-001`, `US-SHIP-002` | ⛔ Blocked | Shipment module scaffold and development tracking screens exist; warehouse pick/pack and shipment APIs do not | Confirm carrier/manual workflow, destinations, tariffs, and split-shipment rules; then implement prepare, handover, tracking events, callbacks, audit/outbox events, UI, and tests | `PRJ-EPIC-010`, `PRJ-EPIC-013`; shipping decision |
| 17 | `PRJ-EPIC-017` | Made-to-order requests | `REQ-CUSTOM-001`; `US-CUSTOM-001..003` | ⬜ Not started | Custom-order module scaffold and documented states exist | Confirm MVP depth; implement request, messages, quote, acceptance, ownership/Casbin, attachments, artisan/customer UI, notifications, transactions, tests, and OpenAPI | `PRJ-EPIC-007`, `PRJ-EPIC-020` |
| 18 | `PRJ-EPIC-018` | Notifications and background workers | Outbox/reliability requirements; notification stories | 🟡 Partial | Outbox schema and password-reset event producer exist; notification module has no worker or delivery adapter | Add worker service, claiming/retry/backoff, notification records, configured adapters, reservation expiry, health/metrics, and idempotent consumer tests | `PRJ-EPIC-002`, producing epics |
| 19 | `PRJ-EPIC-019` | Administration and audit | `REQ-AUDIT-001`; `US-ADMIN-001..003` | 🟡 Partial | Artisan review administration, audit/outbox tables, protected mutation triggers, and Casbin foundation exist | Add role/user/category/product/warehouse administration, audit application service, audit query API/UI, filters/pagination, privacy controls, and tests | `PRJ-EPIC-004`, protected workflows |
| 20 | `PRJ-EPIC-020` | File and media security | `REQ-MEDIA-001`; security upload requirements | 🟡 Partial | MinIO client, persistent buckets, initialization, public/private bucket separation, and media metadata tables exist | Add content-signature validation, configured size limits, generated keys, private authorized reads, signed URLs, moderated publication, cleanup/rollback, and MinIO integration tests | `PRJ-EPIC-004`, `PRJ-EPIC-007`, `PRJ-EPIC-008` |
| 21 | `PRJ-EPIC-021` | CI, contracts, accessibility, and quality gates | Non-functional security, reliability, accessibility, performance; required E2E tests | 🟡 Partial | Jenkins validates backend/frontend checks, migrations, OpenAPI, and Docker builds; unit/integration/browser tests exist | Run dependency-backed integrations in CI, add business E2E/concurrency/security/accessibility suites, container scanning, coverage reporting, and release-blocking quality gates | All implemented workflows |
| 22 | `PRJ-EPIC-022` | Staging, production operations, and rollback | Deployment and production security guide | 🟡 Partial | Compose production overrides, multi-stage non-root images, health checks, Nginx, persistent volumes, and rebuild script exist | Add staging/production provisioning, secret management, TLS/HSTS, worker deployment, immutable image publishing, backups/restores, migration prechecks, monitoring, smoke tests, approvals, and rollback automation | `PRJ-EPIC-021`, all MVP workflows |
| 23 | `PRJ-EPIC-023` | MVP release readiness | MVP scope and required end-to-end tests | ⬜ Not started | — | Resolve launch decisions, complete all included epics, run required workflow/security/accessibility/performance tests, verify backup/rollback, approve staging, and complete production smoke tests | All included epics |

## Feature-Level Task Status

| Jira task | Feature / workflow | Source traceability | Status | Evidence or next action |
|---|---|---|---|---|
| `PRJ-AUTH-001` | Customer registration | `REQ-AUTH-001`, `US-AUTH-001` | ✅ | API validation, password hashing, rate limiting, frontend form, and tests exist |
| `PRJ-AUTH-002` | Login, refresh, logout, password reset | `REQ-AUTH-002`, `US-AUTH-002..003` | ✅ | Session rotation/revocation, reset flow, HTTP-only cookies, BFF routes, and tests exist |
| `PRJ-AUTH-003` | Casbin authorization and ownership | `REQ-RBAC-001` | ✅ | Protected API routes enforce authentication and Casbin; extend policies as new routes are added |
| `PRJ-CUST-001` | Customer profile | Customer role requirements | ✅ | Real profile API and localized account form exist |
| `PRJ-CUST-002` | Customer addresses | `US-CHECK-001` prerequisite | ✅ | CRUD, ownership, single-default transaction, and tests exist |
| `PRJ-CAT-001` | Categories and localized catalogue | `REQ-PROD-002`, `REQ-I18N-001` | ✅ | Public API, category translations, seeded data, and storefront adapter exist |
| `PRJ-CAT-002` | Products and artisans public pages | `US-VIS-001..002` | ✅ | Real public API, publication filtering, media paths, pages, search, and tests exist |
| `PRJ-CAT-003` | Search, filters, collections, and regions | Visitor discovery requirements | 🟡 | Search/category filters exist; collection/region backend APIs are not connected |
| `PRJ-ART-001` | Artisan application submission | `US-ART-001` | ✅ | Authenticated submit/resubmit, translations, categories, and status exist |
| `PRJ-ART-002` | Artisan application review | `US-ART-002` | ✅ | Admin list/decision/reason, role assignment, and authorization exist |
| `PRJ-ART-003` | Artisan documents and profile media | `REQ-ART-001` | 🟡 | Metadata exists; upload validation and private authorized reads remain |
| `PRJ-PROD-001` | Artisan product draft CRUD | `REQ-PROD-001`, `US-PROD-001` | ⬜ | Product package is scaffold-only; implement owned draft APIs and editor |
| `PRJ-PROD-002` | Product media upload | `REQ-MEDIA-001` | ⬜ | Media schema exists; implement MinIO upload and validation |
| `PRJ-PROD-003` | Product submission | `US-PROD-002` | ⬜ | Add validation, `PENDING_REVIEW`, immutable submitted content, and tests |
| `PRJ-MOD-001` | Product moderation | `REQ-MOD-001`, `US-PROD-003` | ⬜ | Add queue, decisions, reasons, policies, audit, and UI |
| `PRJ-WH-001` | Reception | `REQ-WH-001`, `US-WH-001` | ⬜ | Add reception/batch records and warehouse routes |
| `PRJ-WH-002` | Inspection | `REQ-WH-002`, `US-WH-002` | ⬜ | Add quantity equation, evidence, transaction, and inspection UI |
| `PRJ-INV-001` | Inventory ledger | `REQ-INV-001`, `US-WH-003` | ⬜ | Add append-only movements, balances, adjustments, and tests |
| `PRJ-INV-002` | Reservations and expiry | `US-CHECK-002` | ⬜ | Add locking, expiry worker, release, commit, and last-unit concurrency tests |
| `PRJ-CART-001` | Client cart | `REQ-CART-001`, `US-CART-001` | ✅ | Local cart context and responsive cart UI exist |
| `PRJ-CART-002` | Server cart and merge-on-login | `REQ-CART-001` | ⬜ | Add persistence, anonymous cart strategy, merge rules, and API |
| `PRJ-CHECK-001` | Checkout and idempotency | `REQ-CHECKOUT-001`, `US-CHECK-001` | 🟡 | Checkout UI exists; server totals, reservations, order creation, and idempotency are missing |
| `PRJ-ORDER-001` | Orders and immutable snapshots | `REQ-ORDER-001..002`, `US-ORD-001..003` | 🟡 | Account screens exist; order schema/API/state transitions are missing |
| `PRJ-PAY-001` | Payment adapter and webhooks | `REQ-PAY-001`, `US-PAY-001..004` | ⛔ | Waiting for provider and business rules |
| `PRJ-SHIP-001` | Pick, pack, shipment, tracking | `REQ-SHIP-001`, `US-FUL-001`, `US-SHIP-001..002` | ⛔ | Waiting for carrier/manual workflow and destination rules |
| `PRJ-CUSTOM-001` | Made-to-order request and quote | `REQ-CUSTOM-001`, `US-CUSTOM-001..003` | ⬜ | Implement request, messages, quote, acceptance, and notifications |
| `PRJ-ENGAGE-001` | Wishlist persistence | Architecture wishlist feature | 🟡 | Frontend page exists; no authenticated persistence API |
| `PRJ-ENGAGE-002` | Reviews | Architecture review feature | ⬜ | Backend module is a scaffold; confirm MVP inclusion and implement rules |
| `PRJ-NOTIFY-001` | Outbox and notifications | Reliability requirements | 🟡 | Outbox/password-reset producer exists; worker and delivery are missing |
| `PRJ-ADMIN-001` | Users, roles, catalogue administration | Administrator role requirements | 🟡 | Artisan review administration exists; broader admin CRUD is missing |
| `PRJ-ADMIN-002` | Audit query and review | `REQ-AUDIT-001`, `US-ADMIN-003` | 🟡 | Audit storage exists; application audit service and query UI are missing |
| `PRJ-SEC-001` | Upload, private data, and media access | Security §8, `REQ-MEDIA-001` | 🟡 | MinIO foundation exists; signature checks, signed reads, and scanning hooks remain |
| `PRJ-QUALITY-001` | Backend/frontend automated checks | Deployment §7–§8 | 🟡 | Jenkins runs format, vet, tests, builds, migrations, OpenAPI, and Docker checks; full dependency-backed and business E2E gates remain |
| `PRJ-OPS-001` | Deployment, backups, monitoring, rollback | Deployment §9–§17 | 🟡 | Compose/rebuild/health foundation exists; production operations remain |

## Recommended Delivery Order

1. Complete `PRJ-EPIC-008` product drafts and media security.
2. Complete `PRJ-EPIC-009` moderation and publication transitions.
3. Implement `PRJ-EPIC-010` reception and inspection.
4. Implement `PRJ-EPIC-011` inventory and reservations.
5. Replace the client-only cart and checkout surfaces with `PRJ-EPIC-012` and `PRJ-EPIC-013` APIs.
6. Resolve payment and shipping decisions before starting `PRJ-EPIC-015` and `PRJ-EPIC-016`.
7. Add notifications, audit queries, security tests, and release operations continuously with each workflow.

Do not mark a later epic complete because an earlier screen is present. The
managed-stock workflow must remain:

```text
approved product
-> received stock
-> inspected stock
-> accepted stock
-> reservation
-> trusted payment
-> fulfilment
-> shipment
```

## Epic Definition of Done

- Requirements and linked user stories are satisfied without inventing open business rules.
- API, application/domain rules, persistence, UI, and tests are connected for the scoped workflow.
- Protected actions enforce authentication, Casbin permission, and ownership.
- Handlers stay thin; SQLBoiler models never become API responses.
- Product and artisan content supports the supported locales; Arabic RTL is verified.
- Money uses integer minor units and backend-trusted totals.
- Critical state changes use explicit transactions, idempotency where required, audit events, and outbox events.
- Migrations have paired `up`/`down` files, constraints, indexes, and explicit deletion behaviour.
- Private files are validated, stored under generated keys, and read only through authorized access.
- Relevant unit, integration, concurrency, E2E, accessibility, and security tests pass.
- OpenAPI, frontend contracts, environment examples, deployment documentation, and this tracker are synchronized.

## External Decisions Register

| Decision | Status | Jira impact |
|---|---|---|
| Payment provider | Unconfirmed | Blocks `PRJ-EPIC-015`; do not invent a provider or webhook contract |
| Shipping provider/manual process and tariffs | Unconfirmed | Blocks `PRJ-EPIC-016`; do not hard-code carrier, tariff, or destination rules |
| Launch countries and currencies | Unconfirmed | Affects checkout, payment, and shipment validation |
| Tax rules | Unconfirmed | Do not add tax calculations or displays |
| Commission and artisan payout | Unconfirmed | Do not implement split payment or payout logic |
| Refund window and cash on delivery | Unconfirmed | Do not expose policy or payment options until approved |
| Multiple artisans per order and split shipments | Unconfirmed | Affects order, reservation, and fulfilment design |
| Maximum upload sizes and legal documents | Unconfirmed | Configure only after an operational decision |
| Fair-trade, origin, and eco-friendly certification rules | Unconfirmed | Do not expose verified claims without approval |
