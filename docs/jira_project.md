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
- One-account customer authentication, role-aware login destinations, Casbin authorization, and authenticated/admin artisan application routes.
- Public localized category, product, and artisan catalogue APIs with pagination and publication filtering.
- Next.js public storefront, search, filters, product/category/artisan pages, authentication pages, account pages, artisan/admin application screens, artisan media controls, and product authoring screens.
- English, French, Arabic, and Spanish catalogue translations; Arabic RTL support.
- Client-side cart and development deployment/rebuild script.

Not yet connected to the API:

- Artisan membership activation with a mandatory first workshop, multi-workshop ownership, workshop lifecycle management, optional verification states, combined customer/artisan dashboard capabilities, and workshop-linked product ownership.
- Product moderation, warehouse, inventory, reservations, server cart, checkout, orders, payments, shipments, custom orders, reviews, and wishlist persistence.
- The corresponding backend feature directories exist as composition scaffolds in many cases; they do not represent completed functionality.

## Epic Tracker

| Order | Jira Epic | Epic | Requirements and stories | Status | Present in the system | Remaining definition of done | Depends on |
|---:|---|---|---|---|---|---|---|
| 1 | `PRJ-EPIC-001` | Product definition and architecture | Product vision, roles, MVP scope, architecture, security, UI, and deployment documentation | ✅ Implemented | The required planning and architecture documents exist and define the business workflow, layers, state machines, and non-functional rules | Keep this tracker and contracts synchronized as implementation changes | — |
| 2 | `PRJ-EPIC-002` | Platform and development foundation | Deployment baseline; local development; CI foundation | ✅ Implemented | Go API, Next.js frontend, Compose services, health checks, Makefile, Jenkinsfile, migrations, MinIO buckets, and `script.sh` are present | Add the worker runtime and complete staging/production operational controls under `PRJ-EPIC-022` | `PRJ-EPIC-001` |
| 3 | `PRJ-EPIC-003` | Design system, multilingual UI, RTL, and responsive foundation | `REQ-I18N-001`; `US-VIS-001..003`; accessibility and frontend guide | ✅ Implemented | Shared UI primitives, responsive storefront, English/French/Arabic/Spanish copy, locale routing, RTL handling, loading/error/empty states, and frontend tests exist | Apply the same quality bar to future operational screens and complete the full accessibility audit under `PRJ-EPIC-021` | `PRJ-EPIC-001` |
| 4 | `PRJ-EPIC-004` | Authentication and authorization | `REQ-AUTH-001`, `REQ-AUTH-002`, `REQ-RBAC-001`; `US-AUTH-001..003` | ✅ Implemented | Registration, login, refresh rotation, logout, password reset, HTTP-only session handling, rate limiting, Casbin checks, stable errors, and tests exist | Execute PostgreSQL/Redis integration tests in CI and extend authorization tests to every later protected workflow | `PRJ-EPIC-002` |
| 5 | `PRJ-EPIC-005` | Customer profile, addresses, and combined account | `REQ-AUTH-001`, `REQ-AUTH-002`; `US-ACC-001`, `US-AUTH-001..003` | 🟡 Partial | Customer profile/address APIs, authentication, and role-aware login exist; the combined customer/artisan dashboard and complete customer-mode guarantees remain | Connect the combined dashboard and verify that artisan membership never removes customer capabilities | `PRJ-EPIC-004` |
| 6 | `PRJ-EPIC-006` | Public catalogue and discovery | `US-VIS-001..003`; `REQ-PROD-002`, `REQ-I18N-001` | 🟡 Partial | Categories, products, artisans, localized names/content, public media, pagination, category filtering, search, real API adapter, and storefront integration exist | Add public workshop profiles and authoritative stock/availability from `PRJ-EPIC-011`; do not invent availability in the frontend | `PRJ-EPIC-003`, `PRJ-EPIC-004` |
| 7 | `PRJ-EPIC-007` | Artisan membership, onboarding, workshops, and administration | `REQ-ART-001`, `REQ-RBAC-001`; `US-ART-001..008`, `US-WS-001..006`, `US-ADMIN-001..005` | 🟡 Partial | Existing application/profile workflow, private documents/profile media, admin review, role assignment, audit query, and role-aware UI exist | Implement same-account membership activation, mandatory first workshop, multi-workshop CRUD/lifecycle/ownership, optional verification workflow, membership suspension, and complete admin state transitions | `PRJ-EPIC-004`, `PRJ-EPIC-020` |
| 8 | `PRJ-EPIC-008` | Workshop-linked product authoring and media | `REQ-PROD-001`, `REQ-MEDIA-001`; `US-PROD-001..004`, `US-PROD-006`, `US-PROD-008` | 🟡 Partial | Artisan-owned draft CRUD, four-locale validation, price/category rules, private MinIO media, immutable submission snapshots, editor UI, ownership tests, and OpenAPI exist | Add workshop ownership to products, product move rules, workshop status checks, requested-change/archival flows, and historical workshop snapshots | `PRJ-EPIC-007`, `PRJ-EPIC-020` |
| 9 | `PRJ-EPIC-009` | Product moderation and activation | `REQ-MOD-001`; `US-PROD-005`, `US-PROD-007` | ⬜ Not started | Product statuses and submission snapshots are represented in the schema | Implement moderator queue, approve/reject/request-changes/suspend/archive transitions, mandatory reasons, workshop/membership activation checks, stock gating, Casbin policies, audit/outbox events, notifications, UI, and transition tests | `PRJ-EPIC-008` |
| 10 | `PRJ-EPIC-010` | Warehouse reception and inspection | `REQ-WH-001`, `REQ-WH-002`; `US-WH-001`, `US-WH-002` | ⬜ Not started | Warehouse and inspection modules are scaffolds; no operational schema or routes are connected | Implement reception batches, inspection evidence, accepted/rejected/quarantined/damaged quantities, quantity invariant, transactions, roles, UI, tests, and OpenAPI | `PRJ-EPIC-009`, `PRJ-EPIC-020` |
| 11 | `PRJ-EPIC-011` | Inventory ledger and reservations | `REQ-INV-001`; `US-WH-003`, `US-WH-004`, `US-CHECK-002` | ⬜ Not started | Inventory and reservation module scaffolds exist only | Implement append-only movements, balances, adjustments, workshop views, reservation expiry, accepted-stock activation, locking, background expiry, warehouse/customer availability, and last-unit concurrency tests | `PRJ-EPIC-010` |
| 12 | `PRJ-EPIC-012` | Cart and wishlist persistence | `REQ-CART-001`; `US-CART-001`; architecture wishlist feature | 🟡 Partial | Client cart context, cart page/drawer, quantity controls, persistence, wishlist page, and frontend tests exist | Add server/anonymous cart persistence, merge-on-login, active-product revalidation, backend totals, real availability errors, wishlist API, ownership checks, and E2E tests | `PRJ-EPIC-004`, `PRJ-EPIC-006`, `PRJ-EPIC-011` |
| 13 | `PRJ-EPIC-013` | Checkout, orders, cancellation, and returns | `REQ-CHECKOUT-001`, `REQ-ORDER-001`, `REQ-ORDER-002`; `US-CHECK-001..003`, `US-ORD-001..003`, `US-CAN-001..002`, `US-RET-001` | 🟡 Partial | Checkout and customer order/tracking screens exist as frontend surfaces; backend modules are scaffolds and no checkout/order API is registered | Implement server price reload, address selection, transactional reservations, pending orders, immutable line snapshots, idempotency, state transitions, customer/artisan visibility contexts, cancellation, returns, real account integration, and concurrency tests | `PRJ-EPIC-005`, `PRJ-EPIC-011`, `PRJ-EPIC-012` |
| 14 | `PRJ-EPIC-014` | Reviews and customer engagement | Architecture review/wishlist features; post-purchase extension | ⬜ Not started | Review backend package is a scaffold; no review API or persistent review UI is connected | Confirm MVP inclusion, then implement verified-purchase rules, moderation, ownership, rating aggregation, localized UI, and tests | `PRJ-EPIC-013` |
| 15 | `PRJ-EPIC-015` | Payment abstraction, confirmation, and refunds | `REQ-PAY-001`; `US-PAY-001..005` | ⛔ Blocked | Payment requirements and module scaffold exist; no provider adapter or webhook endpoint is implemented | Confirm provider, currencies, payment states, retry/refund rules, and webhook contract; then implement trusted totals, signed/idempotent webhooks, payment events, failure/expiry release, refunds, UI, and tests | `PRJ-EPIC-011`, `PRJ-EPIC-013`; provider decision |
| 16 | `PRJ-EPIC-016` | Fulfilment and shipments | `REQ-SHIP-001`; `US-FUL-001`, `US-SHIP-001`, `US-SHIP-002` | ⛔ Blocked | Shipment module scaffold and development tracking screens exist; warehouse pick/pack and shipment APIs do not | Confirm carrier/manual workflow, destinations, tariffs, and split-shipment rules; then implement prepare, handover, tracking events, callbacks, audit/outbox events, UI, and tests | `PRJ-EPIC-010`, `PRJ-EPIC-013`; shipping decision |
| 17 | `PRJ-EPIC-017` | Made-to-order requests | `REQ-CUSTOM-001`; `US-CUSTOM-001..003` | ⬜ Not started | Custom-order module scaffold and documented states exist | Confirm MVP depth; implement request, messages, quote, acceptance, ownership/Casbin, attachments, artisan/customer UI, notifications, transactions, tests, and OpenAPI | `PRJ-EPIC-007`, `PRJ-EPIC-020` |
| 18 | `PRJ-EPIC-018` | Notifications and background workers | Reliability requirements; `US-NOTIF-001..003` | 🟡 Partial | Outbox schema and password-reset event producer exist; notification module has no worker or delivery adapter | Add membership, commerce, and artisan notification records, worker claiming/retry/backoff, configured adapters, reservation expiry, health/metrics, and idempotent consumer tests | `PRJ-EPIC-002`, producing epics |
| 19 | `PRJ-EPIC-019` | Administration, suspension, verification, and audit | `REQ-AUDIT-001`; `US-ADMIN-001..005` | 🟡 Partial | Artisan review administration, user roles, audit/outbox tables, protected mutation triggers, and Casbin foundation exist | Add user/artisan/workshop/product suspension, optional verification review, category/operational administration, audit privacy controls, and complete filtered query tests | `PRJ-EPIC-004`, protected workflows |
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
| `PRJ-CUST-001` | Customer profile and combined account | `US-ACC-001`, user capabilities §3.1 | 🟡 | Profile/address APIs and account pages exist; one dashboard exposing customer and active artisan modes remains |
| `PRJ-CUST-002` | Customer addresses | `US-CHECK-001` prerequisite | ✅ | CRUD, ownership, single-default transaction, and tests exist |
| `PRJ-CAT-001` | Categories and localized catalogue | `REQ-PROD-002`, `REQ-I18N-001` | ✅ | Public API, category translations, seeded data, and storefront adapter exist |
| `PRJ-CAT-002` | Products, artisans, and workshops public pages | `US-VIS-001..003` | 🟡 | Real public API, publication filtering, media paths, pages, search, and tests exist; public workshop profiles remain |
| `PRJ-CAT-003` | Search, filters, collections, and regions | Visitor discovery requirements | 🟡 | Search/category filters exist; collection/region backend APIs are not connected |
| `PRJ-ART-001` | Start artisan onboarding | `US-ART-001` | 🟡 | Existing application entry exists; same-account onboarding draft and resume flow remain |
| `PRJ-ART-002` | Complete artisan profile | `REQ-ART-001`, `US-ART-002` | 🟡 | Profile fields, translations, categories, and validation exist; consent/declarations and complete membership model remain |
| `PRJ-ART-003` | First mandatory workshop and membership activation | `US-ART-003..004` | ⬜ | No workshop entity or atomic first-workshop/membership activation workflow exists |
| `PRJ-ART-004` | Optional verification documents | `US-ART-005..006` | 🟡 | Private upload and admin reads exist; verification states, reviewer decisions, and user-visible reasons remain |
| `PRJ-ART-005` | Edit and suspend artisan membership | `US-ART-007..008`, `US-ADMIN-003` | 🟡 | Approved profile edits exist; additive membership suspension and explicit restrictions remain |
| `PRJ-WS-001` | Workshop ownership and CRUD | `US-WS-001..006` | ⬜ | Implement workshop list/create/update/status/delete rules, ownership checks, audit, and tests |
| `PRJ-PROD-001` | Workshop-linked product draft CRUD | `REQ-PROD-001`, `US-PROD-001..003` | 🟡 | Owned draft CRUD exists; workshop ownership, move rules, active destination checks, and history rules remain |
| `PRJ-PROD-002` | Product media upload | `REQ-MEDIA-001`, `US-PROD-001` | ✅ | Private MinIO upload/delete APIs, generated keys, signature/size checks, signed reads, UI, and tests exist |
| `PRJ-PROD-003` | Product submission and requested changes | `US-PROD-004`, `US-PROD-006` | 🟡 | Four-locale validation, `PENDING_REVIEW`, immutable snapshots, audit/outbox event, OpenAPI, and tests exist; requested-change editing/resubmission remains |
| `PRJ-PROD-004` | Product archive and workshop history | `US-PROD-007..008` | ⬜ | Add sellable activation gates, archive rules, and historical workshop/order preservation |
| `PRJ-MOD-001` | Product moderation and activation | `REQ-MOD-001`, `US-PROD-005`, `US-PROD-007` | ⬜ | Add queue, decisions, reasons, policies, stock/workshop/membership gates, audit, notifications, and UI |
| `PRJ-WH-001` | Reception | `REQ-WH-001`, `US-WH-001` | ⬜ | Add reception/batch records and warehouse routes |
| `PRJ-WH-002` | Inspection | `REQ-WH-002`, `US-WH-002` | ⬜ | Add quantity equation, evidence, transaction, and inspection UI |
| `PRJ-INV-001` | Inventory ledger and workshop view | `REQ-INV-001`, `US-WH-003..004` | ⬜ | Add append-only movements, balances, adjustments, workshop views, and tests |
| `PRJ-INV-002` | Reservations and expiry | `US-CHECK-002` | ⬜ | Add locking, expiry worker, release, commit, and last-unit concurrency tests |
| `PRJ-CART-001` | Client cart | `REQ-CART-001`, `US-CART-001` | ✅ | Local cart context and responsive cart UI exist |
| `PRJ-CART-002` | Server cart and merge-on-login | `REQ-CART-001`, `US-CART-002` | ⬜ | Add persistence, anonymous cart strategy, merge rules, and API |
| `PRJ-CHECK-001` | Checkout and idempotency | `REQ-CHECKOUT-001`, `US-CHECK-001..003` | 🟡 | Checkout UI exists; server totals, reservations, order creation, idempotency, and artisan-as-customer tests are missing |
| `PRJ-ORDER-001` | Orders, seller visibility, and immutable snapshots | `REQ-ORDER-001..002`, `US-ORD-001..003` | 🟡 | Account screens exist; order schema/API/state transitions and buyer/seller contexts are missing |
| `PRJ-ORDER-002` | Cancellation and returns | `US-CAN-001..002`, `US-RET-001` | ⬜ | Add eligibility, reservation release, refund/return records, stock movements, ownership, and audit rules |
| `PRJ-PAY-001` | Payment adapter, webhooks, and refunds | `REQ-PAY-001`, `US-PAY-001..005` | ⛔ | Waiting for provider and business rules |
| `PRJ-SHIP-001` | Pick, pack, shipment, tracking | `REQ-SHIP-001`, `US-FUL-001`, `US-SHIP-001..002` | ⛔ | Waiting for carrier/manual workflow and destination rules |
| `PRJ-CUSTOM-001` | Made-to-order request and quote | `REQ-CUSTOM-001`, `US-CUSTOM-001..003` | ⬜ | Implement request, messages, quote, acceptance, and notifications |
| `PRJ-ENGAGE-001` | Wishlist persistence | Architecture wishlist feature | 🟡 | Frontend page exists; no authenticated persistence API |
| `PRJ-ENGAGE-002` | Reviews | Architecture review feature | ⬜ | Backend module is a scaffold; confirm MVP inclusion and implement rules |
| `PRJ-NOTIFY-001` | Outbox and notifications | Reliability requirements; `US-NOTIF-001..003` | 🟡 | Outbox/password-reset producer exists; membership, commerce, artisan notifications, worker, and delivery are missing |
| `PRJ-ADMIN-001` | Permissions and account suspension | `US-ADMIN-001..003` | 🟡 | User/role administration exists; user, artisan, workshop, and product suspension workflows remain |
| `PRJ-ADMIN-002` | Audit and verification review | `REQ-AUDIT-001`, `US-ADMIN-004..005` | 🟡 | Protected filtered/paginated audit API and UI exist; optional verification review and complete privacy controls remain |
| `PRJ-SEC-001` | Upload, private data, and ownership access | `REQ-RBAC-001`, `REQ-MEDIA-001`, user-story §17 | 🟡 | Signature checks, generated keys, signed reads, and cleanup exist; workshop ownership, verification access, and broader integration tests remain |
| `PRJ-QUALITY-001` | Backend/frontend automated checks | Requirements §6, user-story §19–§20, Deployment §7–§8 | 🟡 | Jenkins runs format, vet, tests, builds, migrations, OpenAPI, and Docker checks; current membership/workshop/commerce E2E and quality gates remain |
| `PRJ-OPS-001` | Deployment, backups, monitoring, rollback | Deployment §9–§17 | 🟡 | Compose/rebuild/health foundation exists; production operations remain |

## Recommended Delivery Order

1. Complete `PRJ-EPIC-007`: same-account artisan membership, mandatory first workshop, workshop ownership/lifecycle, optional verification, and suspension.
2. Complete `PRJ-EPIC-008`: link products to owned workshops and finish product change, archive, and history rules.
3. Implement `PRJ-EPIC-009` moderation and sellable-product activation.
4. Implement `PRJ-EPIC-010` reception and inspection, then `PRJ-EPIC-011` inventory and reservations.
5. Replace the client-only cart and checkout surfaces with `PRJ-EPIC-012` and `PRJ-EPIC-013` APIs.
6. Resolve payment and shipping decisions before starting `PRJ-EPIC-015` and `PRJ-EPIC-016`.
7. Add notifications, broader administration, security tests, and release operations continuously with each workflow.

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
- One user account preserves customer capabilities while active artisan membership adds seller capabilities.
- Artisan membership activation requires a valid first workshop; products belong to exactly one owned workshop.
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
