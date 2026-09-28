# AISHA MVP Jira Project Tracker

Last reviewed: 2026-09-26

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
- Public localized category, product, artisan, and workshop catalogue APIs with pagination and publication filtering.
- Public product availability derived from the append-only inventory movement ledger; workshop read pages and workshop-filtered product discovery are connected.
- Next.js public storefront, search, filters, product/category/artisan pages, authentication pages, account pages, artisan/admin application screens, artisan media controls, and product authoring screens.
- English, French, Arabic, and Spanish catalogue translations; Arabic RTL support.
- Client and authenticated server cart, anonymous-cart merge, wishlist persistence, and development deployment/rebuild script.

Not yet connected to the API:

- Payment/shipping provider workflows, custom orders, and reviews remain to be connected. PRD-07 artisan membership/workshop activation, lifecycle, verification, suspension, and admin review are connected. Warehouse reception/inspection, inventory adjustments/reservations, server cart/merge, wishlist persistence, checkout, buyer/seller orders, cancellation, and return recording are connected. Product moderation has a connected status-filtered queue, protected decisions, activation gates, and audit/outbox writes. The catalogue consumes read-only workshop and inventory data.
- The corresponding backend feature directories exist as composition scaffolds in many cases; they do not represent completed functionality.

## Epic Tracker

| Order | Jira Epic | Epic | Requirements and stories | Status | Present in the system | Remaining definition of done | Depends on |
|---:|---|---|---|---|---|---|---|
| 1 | `PRJ-EPIC-001` | Product definition and architecture | Product vision, roles, MVP scope, architecture, security, UI, and deployment documentation | ✅ Implemented | The required planning and architecture documents exist and define the business workflow, layers, state machines, and non-functional rules | Keep this tracker and contracts synchronized as implementation changes | — |
| 2 | `PRJ-EPIC-002` | Platform and development foundation | Deployment baseline; local development; CI foundation | ✅ Implemented | Go API, Next.js frontend, Compose services, health checks, Makefile, Jenkinsfile, migrations, MinIO buckets, and `script.sh` are present | Add the worker runtime and complete staging/production operational controls under `PRJ-EPIC-022` | `PRJ-EPIC-001` |
| 3 | `PRJ-EPIC-003` | Design system, multilingual UI, RTL, and responsive foundation | `REQ-I18N-001`; `US-VIS-001..003`; accessibility and frontend guide | ✅ Implemented | Shared UI primitives, responsive storefront, English/French/Arabic/Spanish copy, locale routing, RTL handling, loading/error/empty states, and frontend tests exist | Apply the same quality bar to future operational screens and complete the full accessibility audit under `PRJ-EPIC-021` | `PRJ-EPIC-001` |
| 4 | `PRJ-EPIC-004` | Authentication and authorization | `REQ-AUTH-001`, `REQ-AUTH-002`, `REQ-RBAC-001`; `US-AUTH-001..003` | ✅ Implemented | Registration, login, refresh rotation, logout, password reset, HTTP-only session handling, rate limiting, Casbin checks, stable errors, and tests exist | Execute PostgreSQL/Redis integration tests in CI and extend authorization tests to every later protected workflow | `PRJ-EPIC-002` |
| 5 | `PRJ-EPIC-005` | Customer profile, addresses, and combined account | `REQ-AUTH-001`, `REQ-AUTH-002`; `US-ACC-001`, `US-AUTH-001..003` | ✅ Implemented | Backend-derived `/me` account summary, owner-scoped profile/address CRUD, transactional single-default address handling, combined customer/artisan navigation, explicit account/profile/address loading-empty-error states, localized responsive UI, OpenAPI response schemas, and tests | Canonical membership activation and real buyer orders/tracking remain owned by `PRJ-EPIC-007` and `PRJ-EPIC-013`; they are dependencies, not gaps in the PRD-05 account surface | `PRJ-EPIC-004`, `PRJ-EPIC-007`, `PRJ-EPIC-013` |
| 6 | `PRJ-EPIC-006` | Public catalogue and discovery | `US-VIS-001..003`; `REQ-PROD-002`, `REQ-I18N-001` | ✅ Implemented | Public APIs filter active/published products, approved artisans, and public workshops; localized craft categories, server-side product/workshop search, inventory-ledger availability, category/availability filters, URL-preserving pagination, responsive product/workshop pages, empty/error states, signed public product-media URLs, visible available quantities, OpenAPI, and tests are connected | Inventory mutation/reservation workflows remain in `PRJ-EPIC-011`; catalogue owns the public read surface and PRD-07 owns the completed private workshop lifecycle | `PRJ-EPIC-003`, `PRJ-EPIC-004` |
| 7 | `PRJ-EPIC-007` | Artisan membership, onboarding, workshops, and administration | `REQ-ART-001`, `REQ-RBAC-001`; `US-ART-001..008`, `US-WS-001..006`, `US-ADMIN-001..005` | ✅ Implemented | Migration 000015, atomic/idempotent same-account activation, mandatory default workshop, owned workshop CRUD/status/delete protection, verification state/decision APIs, suspension gates, audit/outbox events, responsive artisan workspace, administrator review UI, and dependency-backed ownership/concurrency tests are connected | Notification delivery and document-policy operations remain cross-cutting follow-up work in `PRJ-EPIC-018` and `PRJ-EPIC-020`; they do not block membership activation | `PRJ-EPIC-004`, `PRJ-EPIC-018`, `PRJ-EPIC-020` |
| 8 | `PRJ-EPIC-008` | Workshop-linked product authoring and media | `REQ-PROD-001`, `REQ-MEDIA-001`; `US-PROD-001..004`, `US-PROD-006`, `US-PROD-008` | ✅ Implemented | Workshop-owned draft CRUD, ACTIVE-membership/workshop gates, single-selected-language validation, private MinIO media, immutable workshop-aware submission snapshots, requested-change resubmission, owner archive controls, editor UI, migration, tests, and OpenAPI are connected | Continue moderation notifications and operational media policy under PRD-09/PRD-20; these do not block PRD-08 authoring | `PRJ-EPIC-007`, `PRJ-EPIC-020` |
| 9 | `PRJ-EPIC-009` | Product moderation and activation | `REQ-MOD-001`; `US-PROD-005`, `US-PROD-007` | ✅ Implemented | Migration 000014, status-filtered paginated moderator queue/decision API, approval-generated unique warehouse product codes, activation gates, Casbin policies, audit/outbox events, moderation UI, OpenAPI, and repository integration coverage are connected | Add notification delivery, reviewer assignment, and deeper transition/concurrency coverage; storage-bucket publication remains an operational follow-up if separate public storage is required | `PRJ-EPIC-008` |
| 10 | `PRJ-EPIC-010` | Warehouse reception and inspection | `REQ-WH-001`, `REQ-WH-002`; `US-WH-001`, `US-WH-002` | ✅ Implemented | Migration 000016 plus 000020 product codes, idempotent reception API, protected artisan-phone/workshop/product-code search, private evidence upload/signed reads, transactional inspection, quantity invariant, accepted/rejected/quarantined/damaged stock buckets, Casbin warehouse roles, admin UI, OpenAPI, unit tests, and PostgreSQL integration coverage are connected | Add operational reinspection policy for returned stock and evidence retention rules | `PRJ-EPIC-009`, `PRJ-EPIC-020` |
| 11 | `PRJ-EPIC-011` | Inventory ledger and reservations | `REQ-INV-001`; `US-WH-003`, `US-WH-004`, `US-CHECK-002` | 🟡 Partial | Migration 000017, append-only adjustments, audit/outbox writes, workshop-scoped balance API, product-row-locked checkout reservations, expiry worker, quarantined return movements, and migration 000019 return idempotency are connected and tested | Add payment/fulfilment commit and shipment movements plus dependency-backed last-unit integration coverage; blocked policy decisions remain explicit | `PRJ-EPIC-010` |
| 12 | `PRJ-EPIC-012` | Cart and wishlist persistence | `REQ-CART-001`; `US-CART-001`; architecture wishlist feature | ✅ Implemented | Migration 000018, owner-scoped server cart CRUD, local anonymous cart, merge-on-login, current product/availability DTOs, active-product repository guards, API/BFF routes, wishlist API, product controls, account UI, and frontend regression tests are connected | Add back-in-stock notifications only if the capability is included in launch scope | `PRJ-EPIC-004`, `PRJ-EPIC-006`, `PRJ-EPIC-011` |
| 13 | `PRJ-EPIC-013` | Checkout, orders, cancellation, and returns | `REQ-CHECKOUT-001`, `REQ-ORDER-001`, `REQ-ORDER-002`; `US-CHECK-001..003`, `US-ORD-001..003`, `US-CAN-001..002`, `US-RET-001` | 🟡 Partial | Migration 000014 plus 000019, trusted transactional checkout, UUID/idempotency validation, reservations, immutable snapshots, buyer/seller visibility, persisted shipment-event display, cancellation, quarantined and concurrency-safe returns, BFF routes, and live order screens are connected and tested | Payment confirmation/refunds, shipment integration, split-shipment policy, and dependency-backed concurrency remain blocked by the external provider/shipping decisions | `PRJ-EPIC-005`, `PRJ-EPIC-011`, `PRJ-EPIC-012` |
| 14 | `PRJ-EPIC-014` | Reviews and customer engagement | Architecture review/wishlist features; post-purchase extension | ⬜ Not started | Review backend package is a scaffold; no review API or persistent review UI is connected | Confirm MVP inclusion, then implement verified-purchase rules, moderation, ownership, rating aggregation, localized UI, and tests | `PRJ-EPIC-013` |
| 15 | `PRJ-EPIC-015` | Payment abstraction, confirmation, and refunds | `REQ-PAY-001`; `US-PAY-001..005` | ⛔ Blocked | Payment requirements and module scaffold exist; no provider adapter or webhook endpoint is implemented | Confirm provider, currencies, payment states, retry/refund rules, and webhook contract; then implement trusted totals, signed/idempotent webhooks, payment events, failure/expiry release, refunds, UI, and tests | `PRJ-EPIC-011`, `PRJ-EPIC-013`; provider decision |
| 16 | `PRJ-EPIC-016` | Fulfilment and shipments | `REQ-SHIP-001`; `US-FUL-001`, `US-SHIP-001`, `US-SHIP-002` | ⛔ Blocked | Shipment module scaffold and development tracking screens exist; warehouse pick/pack and shipment APIs do not | Confirm carrier/manual workflow, destinations, tariffs, and split-shipment rules; then implement prepare, handover, tracking events, callbacks, audit/outbox events, UI, and tests | `PRJ-EPIC-010`, `PRJ-EPIC-013`; shipping decision |
| 17 | `PRJ-EPIC-017` | Made-to-order requests | `REQ-CUSTOM-001`; `US-CUSTOM-001..003` | ⬜ Not started | Custom-order module scaffold and documented states exist | Confirm MVP depth; implement request, messages, quote, acceptance, ownership/Casbin, attachments, artisan/customer UI, notifications, transactions, tests, and OpenAPI | `PRJ-EPIC-007`, `PRJ-EPIC-020` |
| 18 | `PRJ-EPIC-018` | Notifications and background workers | Reliability requirements; `US-NOTIF-001..003` | 🟡 Partial | Notification persistence/API, recipient-scoped idempotent outbox consumer, retry/backoff, reservation-expiry worker, and authenticated WebSocket delivery are implemented | Add configured email/push adapters, dead-letter/admin visibility, health metrics, and producers for any future payment/shipment/custom-order events | `PRJ-EPIC-002`, producing epics |
| 19 | `PRJ-EPIC-019` | Administration, suspension, verification, and audit | `REQ-AUDIT-001`; `US-ADMIN-001..005` | 🟡 Partial | Admin dashboard has nested Users/Workers/artisan-application/moderation/audit sections, consistent responsive table panels with horizontal scrolling, top-left contextual actions, paginated rows with Details modals, in-dashboard admin profile/password management, searchable role/status controls, semantic status badges, locale-aware full dates, artisan review administration, administrator workshop lifecycle controls, product moderation suspension, verification review, audit/outbox tables, protected mutation triggers, and Casbin authorization; administrator artisan capabilities are explicitly denied | Add broader suspension-policy coverage, audit privacy/retention controls, and complete filtered/concurrency integration tests | `PRJ-EPIC-004`, protected workflows |
| 20 | `PRJ-EPIC-020` | File and media security | `REQ-MEDIA-001`; security upload requirements | 🟡 Partial | MinIO client, persistent buckets, initialization, public/private bucket separation, and media metadata tables exist | Add content-signature validation, configured size limits, generated keys, private authorized reads, signed URLs, moderated publication, cleanup/rollback, and MinIO integration tests | `PRJ-EPIC-004`, `PRJ-EPIC-007`, `PRJ-EPIC-008` |
| 21 | `PRJ-EPIC-021` | CI, contracts, accessibility, and quality gates | Non-functional security, reliability, accessibility, performance; required E2E tests | 🟡 Partial | Jenkins validates backend/frontend checks, migrations, OpenAPI, and Docker builds; unit/integration/browser tests exist | Run dependency-backed integrations in CI, add business E2E/concurrency/security/accessibility suites, container scanning, coverage reporting, and release-blocking quality gates | All implemented workflows |
| 22 | `PRJ-EPIC-022` | Staging, production operations, and rollback | Deployment and production security guide | 🟡 Partial | Compose production overrides, multi-stage non-root images, health checks, Nginx, persistent volumes, rebuild script, and an opt-in Jenkins VPS deployment with SSH credentials, PostgreSQL backup, commit-tagged images, migration prechecks, smoke tests, and image rollback exist | Add staging/production provisioning, secret management, TLS/HSTS, worker deployment, MinIO backup/restore, monitoring, approvals, and full rollback drills | `PRJ-EPIC-021`, all MVP workflows |
| 23 | `PRJ-EPIC-023` | MVP release readiness | MVP scope and required end-to-end tests | ⬜ Not started | — | Resolve launch decisions, complete all included epics, run required workflow/security/accessibility/performance tests, verify backup/rollback, approve staging, and complete production smoke tests | All included epics |

## Feature-Level Task Status

| Jira task | Feature / workflow | Source traceability | Status | Evidence or next action |
|---|---|---|---|---|
| `PRJ-AUTH-001` | Customer registration | `REQ-AUTH-001`, `US-AUTH-001` | ✅ | API validation, password hashing, rate limiting, frontend form, and tests exist |
| `PRJ-AUTH-002` | Login, refresh, logout, password reset | `REQ-AUTH-002`, `US-AUTH-002..003` | ✅ | Session rotation/revocation, reset flow, HTTP-only cookies, BFF routes, configurable HTTP/HTTPS cookie security for the current VPS, and tests exist |
| `PRJ-AUTH-003` | Casbin authorization and ownership | `REQ-RBAC-001` | ✅ | Protected API routes enforce authentication and Casbin; extend policies as new routes are added |
| `PRJ-CUST-001` | Customer profile and combined account | `US-ACC-001`, user capabilities §3.1 | ✅ | Backend-derived `/me` summary, same-user customer/artisan capability model, responsive/RTL navigation, explicit loading/error/unavailable states, profile and address ownership protections, localized UI, OpenAPI, and tests; canonical membership and buyer order/tracking remain in their owning epics |
| `PRJ-CUST-002` | Customer addresses | `US-CHECK-001` prerequisite | ✅ | CRUD, ownership, single-default transaction, and tests exist |
| `PRJ-CAT-001` | Categories and localized catalogue | `REQ-PROD-002`, `REQ-I18N-001` | ✅ | Public API, category translations, seeded data, and storefront adapter exist |
| `PRJ-DATA-001` | Development fixture directory | Local development workflow | ✅ | The development seed now includes sixteen named customer, artisan, moderator, and warehouse-agent accounts, three additional approved artisan workshops, four-language workshop/profile copy, and six public artisan media records backed by checked-in images |
| `PRJ-CAT-002` | Products, artisans, and workshops public pages | `US-VIS-001..003` | ✅ | Public product, artisan, and workshop APIs/pages, publication filtering, localized craft categories, privacy-safe verification display, media paths, search, and tests are connected |
| `PRJ-CAT-003` | Search, filters, collections, and regions | Visitor discovery requirements | 🟡 | Localized server-side product/workshop search, category and availability filters, artisan search/region/craft filters, and URL-preserving pagination exist; collection/region backend APIs are not connected |
| `PRJ-ART-001` | Start artisan onboarding | `US-ART-001` | ✅ | Existing customers save a draft first, upload at least one supporting document and one profile media item, then use a server-validated final submit; active-membership dashboard routing and customer capability preservation remain connected |
| `PRJ-ART-002` | Complete artisan profile | `REQ-ART-001`, `US-ART-002` | ✅ | Owner-scoped profile fields, translations, categories, workshop details, private MinIO documents/media, draft/resume persistence, required-file validation, and review state are connected |
| `PRJ-ART-003` | First mandatory workshop and membership activation | `US-ART-003..004` | ✅ | Approval now atomically creates or reactivates the ACTIVE membership, submitted default workshop, verification row, role grant, and audit/outbox records; the existing idempotent `POST /artisan/membership/activate` remains available for legacy approved profiles and is connected to the responsive activation form |
| `PRJ-ART-004` | Optional verification documents | `US-ART-005..006` | ✅ | Private upload advances verification to PENDING and queues a reviewer outbox event; user state and responsive admin list/decision UI expose all verification states with reasons and audited decisions; artisan private files use dated rows on desktop and labelled stacked cards on mobile, signed View links, modal Add/Update, owner-scoped Delete, and MinIO cleanup | Notification delivery and legal-document policy remain in `PRJ-EPIC-018`/`PRJ-EPIC-020` |
| `PRJ-ART-005` | Edit and suspend artisan membership | `US-ART-007..008`, `US-ADMIN-003` | ✅ | Profile edits remain owner-scoped; responsive admin membership controls record required reasons/audit state, preserve customer capabilities, and gate seller/public catalogue/checkout actions | Broader account/product suspension policy remains in `PRJ-EPIC-019` |
| `PRJ-WS-001` | Workshop ownership and CRUD | `US-WS-001..006` | ✅ | Workshop list/create/update/status/delete APIs and BFF/UI are connected with owner joins, default-workshop protection, product/order history deletion protection, idempotency, workshop audit/outbox events, cross-owner tests, suspended-member gating, responsive first-column-actions tables with Details/Update/Delete dialogs, and the approved-workspace Workshops/Private files/Product authoring tabs |
| `PRJ-PROD-001` | Workshop-linked product draft CRUD | `REQ-PROD-001`, `US-PROD-001..003` | ✅ | Products require an owned active workshop; owner-only create/read/update is enforced in repository queries; moves are allowed only between owned active workshops and blocked after inventory movements; the artisan workspace now provides the responsive product table, Details/Update actions, history-preserving Archive action, and modal New draft/Update editor |
| `PRJ-PROD-002` | Product media upload | `REQ-MEDIA-001`, `US-PROD-001` | ✅ | Private MinIO upload/delete APIs, generated keys, signature/size checks, signed reads, UI, and tests exist |
| `PRJ-PROD-003` | Product submission and requested changes | `US-PROD-004`, `US-PROD-006` | ✅ | Four-locale/media validation, active-workshop checks, `PENDING_REVIEW`, immutable workshop-aware snapshots, CHANGES_REQUESTED editing/resubmission, audit/outbox events, OpenAPI, and tests exist |
| `PRJ-PROD-004` | Product archive and workshop history | `US-PROD-007..008` | ✅ | Owner archive is authenticated and audited, admin archive/activation decisions are protected, order lines preserve product/artisan/workshop snapshots, archived products clear publication state, and workshop deletion is blocked once product/inventory/order history exists; repository integration coverage verifies archive idempotency and history protection | Continue deeper cross-epic order/return coverage under `PRJ-EPIC-013` |
| `PRJ-MOD-001` | Product moderation and activation | `REQ-MOD-001`, `US-PROD-005`, `US-PROD-007` | ✅ | Status-filtered paginated queue, protected decisions, approval-generated stable warehouse product code, mandatory reasons, activation gates, Casbin policies, atomic media publication state, audit/outbox writes, OpenAPI, UI, signed product-media previews in Details, unit tests, and repository integration coverage are implemented | Add notification delivery, reviewer assignment, and deeper transition/concurrency tests |
| `PRJ-WH-001` | Reception | `REQ-WH-001`, `US-WH-001` | ✅ | Idempotent reception records, protected artisan-phone/workshop/product-code search, UUID-safe workshop/product filtering, local/international Algerian phone matching, selected product context, approved product code and available quantity in the workshop combobox, private evidence upload, warehouse authorization, API, UI, tests, and automatic hand-off to the inspection details step are connected |
| `PRJ-WH-002` | Inspection | `REQ-WH-002`, `US-WH-002` | ✅ | One-time transactional inspection, mandatory reason/evidence, quantity equation, append-only bucketed movements, audit/outbox events, API, UI, client-side completion guards, and tests are connected |
| `PRJ-INV-001` | Inventory ledger and workshop view | `REQ-INV-001`, `US-WH-003..004` | ✅ | Append-only movement migration, balance aggregation, reasoned adjustment API, audit/outbox events, workshop ownership filtering, UUID-safe inventory loading, product-code display without internal IDs, admin workshop combobox/table/modal, OpenAPI, and unit tests are connected |
| `PRJ-INV-002` | Reservations and expiry | `US-CHECK-002` | 🟡 | Checkout and the standalone inventory worker release expired holds with row locks and idempotent movements; commit/shipment transitions and dependency-backed last-unit concurrency coverage remain blocked by payment/fulfilment decisions |
| `PRJ-CART-001` | Client cart | `REQ-CART-001`, `US-CART-001` | ✅ | Local cart context and responsive cart UI exist |
| `PRJ-CART-002` | Server cart and merge-on-login | `REQ-CART-001`, `US-CART-002` | ✅ | Authenticated cart persistence, owner-scoped CRUD, anonymous local-cart merge-on-login, product/availability revalidation DTOs, public cart/wishlist images and saved dates, compact image thumbnails with cursor-position lightbox zoom, API/BFF routes, frontend synchronization, UUID-free cart display metadata, and localized unauthenticated wishlist feedback are connected |
| `PRJ-CHECK-001` | Checkout and idempotency | `REQ-CHECKOUT-001`, `US-CHECK-001..003` | 🟡 | Server trusted totals, address ownership, reservations, pending orders, idempotency, BFF checkout, expiry worker, cart-metadata subtotal rendering, and focused service coverage are connected | Payment confirmation and dependency-backed concurrency tests remain under the blocked payment epic |
| `PRJ-ORDER-001` | Orders, seller visibility, and immutable snapshots | `REQ-ORDER-001..002`, `US-ORD-001..003` | 🟡 | Buyer order API/live screens now render persisted shipment events; seller item-scoped API, status transitions, and immutable snapshots remain connected | Add dependency-backed integration/concurrency coverage and fulfilment transitions |
| `PRJ-ORDER-002` | Cancellation and returns | `US-CAN-001..002`, `US-RET-001` | 🟡 | Buyer cancellation releases held reservations and records audit/outbox; authorized fulfilment return recording is idempotent, audited, and remains non-sellable pending inspection | Add refund provider integration, explicit return-inspection workflow, and policy-backed expiry handling |
| `PRJ-PAY-001` | Payment adapter, webhooks, and refunds | `REQ-PAY-001`, `US-PAY-001..005` | ⛔ | Waiting for provider and business rules |
| `PRJ-SHIP-001` | Pick, pack, shipment, tracking | `REQ-SHIP-001`, `US-FUL-001`, `US-SHIP-001..002` | ⛔ | Waiting for carrier/manual workflow and destination rules |
| `PRJ-CUSTOM-001` | Made-to-order request and quote | `REQ-CUSTOM-001`, `US-CUSTOM-001..003` | ⬜ | Implement request, messages, quote, acceptance, and notifications |
| `PRJ-ENGAGE-001` | Wishlist persistence | Architecture wishlist feature | ✅ | Authenticated owner-scoped wishlist API, idempotent save/remove, BFF routes, product-card controls, and account wishlist UI are connected |
| `PRJ-ENGAGE-002` | Reviews | Architecture review feature | ⬜ | Backend module is a scaffold; confirm MVP inclusion and implement rules |
| `PRJ-NOTIFY-001` | Outbox and notifications | Reliability requirements; `US-NOTIF-001..003` | 🟡 | Outbox/password-reset producer exists; membership, commerce, artisan notifications, worker, and delivery are missing |
| `PRJ-ADMIN-001` | Permissions and account suspension | `US-ADMIN-001..003` | 🟡 | Role-scoped back-office access is connected: the Overview tab was removed in favour of Users and Workers tabs; moderators review products, warehouse agents manage reception/inventory, and administrators can create/update accounts, assign customer/artisan/moderator/warehouse-agent/administrator roles, activate/suspend/disable accounts, and use sanitized paginated Details modals | Add broader suspension-policy coverage and concurrency integration tests |
| `PRJ-ADMIN-002` | Audit and verification review | `REQ-AUDIT-001`, `US-ADMIN-004..005` | 🟡 | Protected filtered/paginated audit API/UI with row Details modals, artisan application and approved-artisan user-media previews, product moderation media Details, approved-artisan shop/status/action table with sanitized Details modal, and confirmation-protected user/product media deletion with audit and MinIO cleanup exist; admin media list queries and browser-facing MinIO presigning are verified | Add audit privacy/retention controls and deeper filtered integration tests |
| `PRJ-ADMIN-003` | Administrative catalogue and commerce controls | `US-ADMIN-001..005` | 🟡 | Added the nested Categories tab with CRUD, four-language labels, benefit-rate basis points, soft deactivation, audit events, API/BFF routes, OpenAPI, responsive overflow, and tests; added a paginated Orders tab with sanitized customer/status/total/date details and the existing warehouse, inventory, moderation, media, application, and audit tabs | Add explicit warehouse entity/agent assignment, order fulfilment mutations, and backup/restore operations after the deployment security contract is defined |
| `PRJ-SEC-001` | Upload, private data, and ownership access | `REQ-RBAC-001`, `REQ-MEDIA-001`, user-story §17 | 🟡 | Signature checks, generated keys, browser-resolvable public-endpoint signed reads, legacy public asset compatibility, media DTO redaction, cleanup, and admin application media access exist; workshop ownership, verification access, and broader integration tests remain |
| `PRJ-QUALITY-001` | Backend/frontend automated checks | Requirements §6, user-story §19–§20, Deployment §7–§8 | 🟡 | Jenkins runs format, vet, tests, builds, migrations, OpenAPI, and Docker checks; shared AISHA theme tokens, Light/Dark/System switching with keyboard/outside-click handling, route-preserving flag/abbreviation language switching, searchable comboboxes, semantic badges, full localized date formatting, responsive/RTL controls, mobile-safe header/menu/account layouts, overflow-safe modals and media/table surfaces, persistent customer/artisan mode switching without losing either capability, approved-artisan routing to the existing workshop, first-visit-only branded splash, and frontend validation are connected; current membership/workshop/commerce E2E and quality gates remain |
| `PRJ-I18N-002` | Locale-complete operational UI | `REQ-I18N-001`, `US-VIS-001..003` | 🟡 | Added a shared locale-aware copy/status layer and connected language switching to preserve the active route, query string, and hash. Public, artisan, notification, catalogue, order, inventory, and administration surfaces now use translated table labels, statuses, controls, pagination, loading, empty, and modal copy; remaining legacy hardcoded copy is tracked for the next accessibility/localization pass |
| `PRJ-UI-003` | Light-mode tab panels and mobile tab navigation | Frontend responsive/UI guide | ✅ | Artisan, account, and administration tab panels now use readable light-mode brand tokens, visible active states, card contrast, focus-safe controls, and horizontal overflow that keeps mobile tab labels and content inside the viewport |
| `PRJ-OPS-001` | Deployment, backups, monitoring, rollback | Deployment §9–§17 | 🟡 | Compose/rebuild/health foundation exists; VPS HTTP session-cookie configuration and production operations remain |

## Product-to-Warehouse Handoff Update (2026-09-27)

- `PRJ-EPIC-008` / `PRJ-PROD-001`: product drafts now persist planned quantity and total order price, require one selected supported locale payload at submission, and enforce a transactional maximum of four media files.
- `PRJ-EPIC-009` / `PRJ-MOD-001`: moderator approval exposes the stable warehouse product code to the artisan workflow.
- `PRJ-EPIC-010` / `PRJ-WH-002`: accepted warehouse inspection stock automatically activates an approved product and publishes its media in the same transaction when the existing gates pass; otherwise the product remains unavailable.
- Evidence: migration `000021_product_order_fields`, product repository/service/handler, warehouse inspection repository, OpenAPI schemas, artisan product workspace, and Go/frontend test suites.

## Recommended Delivery Order

1. Complete notification delivery and operational order/return coverage under `PRJ-EPIC-018` and `PRJ-EPIC-013`.
2. Complete `PRJ-EPIC-011` reservation workers, payment/fulfilment transitions, and concurrency coverage on top of the accepted stock ledger.
3. Replace the client-only cart and checkout surfaces with `PRJ-EPIC-012` and `PRJ-EPIC-013` APIs.
4. Resolve payment and shipping decisions before starting `PRJ-EPIC-015` and `PRJ-EPIC-016`.
5. Add broader administration, security tests, and release operations continuously with each workflow.

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
