# AISHA Full Project Epic Tracker

This tracker follows the approved AISHA requirements, architecture, frontend, backend, design, security, and deployment guides. Complete epics in dependency order. A frontend screen backed only by mock data is not considered a completed end-to-end epic.

For backend implementation details, use [`docs/jira_backend.md`](./jira_backend.md).

## Status Legend

| Status | Meaning |
|---|---|
| ✅ Implemented | The scoped code and relevant automated tests are present in the repository. |
| 🟡 Partial | Some required layers or workflows exist, but the epic is not complete end to end. |
| ⬜ Not started | No complete implementation is present. |
| ⛔ Blocked | A documented external decision or credential is required before implementation can finish. |

## Project Epics

| Order | Jira Epic | Epic | Status | Implemented now | Remaining scope / definition of done | Depends on |
|---:|---|---|---|---|---|---|
| 1 | `PRJ-EPIC-001` | Product requirements and architecture | ✅ Implemented | MVP requirements, user stories, modular-monolith architecture, backend/frontend guides, design system specification, security rules, and deployment guide are documented | Keep documentation synchronized with implemented contracts and explicitly resolve open business decisions before dependent work | — |
| 2 | `PRJ-EPIC-002` | Development platform foundation | ✅ Implemented | Next.js/TypeScript frontend, Go Fiber API, PostgreSQL, Redis, MinIO, Nginx, Docker Compose, health checks, Make targets, and Jenkins validation foundation | Preserve persistent volumes and keep local/CI commands reproducible as later services are added | Product architecture |
| 3 | `PRJ-EPIC-003` | Shared frontend design system | ✅ Implemented | Design tokens, explicit light color scheme, global styles, buttons, badges, skeletons, accessible modal/sheet foundation, containers, breadcrumbs, headings, empty/error/form feedback states, commerce/editorial components, responsive layouts, focus management, reduced-motion behavior, and component tests exist | Extend the same primitives rather than adding page-specific alternatives when later operational workspaces are implemented | Platform foundation |
| 4 | `PRJ-EPIC-004` | Multilingual and RTL experience | ✅ Implemented | English, French, Arabic, and Spanish routing/messages exist; public interaction copy is translated without English fallback; locale-aware content/number/currency utilities, Arabic RTL layout, logical positioning, mirrored directional icons, RTL overlays, and desktop/mobile E2E tests exist | Maintain complete translation keys and RTL behavior as later customer, artisan, and admin workflows are implemented | Frontend design system |
| 5 | `PRJ-EPIC-005` | Backend database foundation | ✅ Implemented | Versioned migrations, users/sessions, artisan/catalogue schemas, audit/outbox/idempotency tables, reference seeds, SQLBoiler models, and migration tests exist | Run migration round-trip tests against PostgreSQL in CI and regenerate models after every schema change | Platform foundation |
| 6 | `PRJ-EPIC-006` | API errors, DTOs, authentication, and authorization | ✅ Implemented | Stable API errors, DTO validation, registration/login/session rotation/logout/password reset, `/me`, Casbin enforcement, Redis rate limiting, OpenAPI documentation, and automated tests exist | Run conditional PostgreSQL and Redis integration tests in CI; require authentication plus explicit Casbin permission on every later protected route | Database foundation |
| 7 | `PRJ-EPIC-007` | Editorial public storefront | ✅ Implemented | Responsive editorial home with categories, selected products, artisan story, regional discovery, made-to-order campaign, trust and newsletter sections; accessible header/footer/search/filter overlays; category/product/artisan/search pages; gallery zoom; category-matched imagery; localized metadata; URL-backed filters; route loading/error/empty states; and desktop/mobile/RTL E2E coverage exist | Replace the typed storefront data adapter with real queries under `PRJ-EPIC-008` without changing presentational component contracts | Design system, multilingual foundation; real data depends on public catalogue API |
| 8 | `PRJ-EPIC-008` | Public catalogue API | ✅ Implemented | Five public endpoints, localized fallback queries, publication/approval/public-media filtering, bounded pagination and category filtering, explicit DTOs and OpenAPI, PostgreSQL/handler/service tests, and real API integration across home, product, category, search, artisan, navigation, and filter surfaces are implemented | Inventory availability remains intentionally owned by `PRJ-EPIC-015`; the catalogue UI does not invent stock | Authentication/authorization, database foundation |
| 9 | `PRJ-EPIC-009` | Customer authentication experience | ✅ Implemented | Registration, login, refresh rotation, logout, forgot/reset password, protected navigation, HTTP-only access/refresh cookies, same-origin BFF routes, localized validation and rate-limit feedback, reset completion, unit tests, and desktop/mobile browser success/failure coverage are implemented | Maintain session and authorization coverage as later protected workspaces are added | Authentication API, multilingual frontend |
| 10 | `PRJ-EPIC-010` | Customer profile and addresses | ✅ Implemented | Paired migration and generated model, profile/address backend use cases and DTOs, principal-derived ownership, Casbin checks, transactional single-default writes, localized real forms, service/handler/PostgreSQL tests, and desktop/mobile address CRUD browser coverage are implemented | Reuse these owned addresses during checkout under `PRJ-EPIC-017` | Customer authentication |
| 11 | `PRJ-EPIC-011` | Artisan onboarding and profile | ✅ Implemented | Authenticated application submission/resubmission and status, category and multilingual biography data, approved-profile editing, administrator queue and decisions, mandatory decision reasons, atomic artisan-role assignment, audit/outbox events, private reviewer-only document metadata, Casbin/ownership enforcement, responsive artisan/admin workspaces, service and PostgreSQL transaction tests, migration, generated models, and OpenAPI are implemented | Secure document/profile-media upload, content-signature validation, and authorised signed reads remain explicitly tracked under `PRJ-EPIC-024` | Authentication, file security, notifications |
| 12 | `PRJ-EPIC-012` | Artisan products and media | ⬜ Not started | Product/translation/media schema and public product UI components exist | Implement protected draft creation/edit/list/submit flows, four-locale content, ownership checks, server-validated pricing/category rules, secure MinIO upload, frontend artisan editor, tests, and OpenAPI | Artisan onboarding, public catalogue |
| 13 | `PRJ-EPIC-013` | Product moderation | ⬜ Not started | Product statuses can represent the main moderation lifecycle | Implement immutable/reviewable submissions, moderation queue, approve/request-changes/reject/suspend transitions, mandatory reasons, Casbin policies, audit/outbox events, admin UI, notifications, and transition/authorization tests | Artisan products |
| 14 | `PRJ-EPIC-014` | Warehouse reception and inspection | ⬜ Not started | — | Add warehouse/reception/inspection schema and workflows; implement quantity invariants, evidence, Casbin roles, explicit transactions, audit/outbox events, warehouse UI, tests, and OpenAPI | Product moderation, file security |
| 15 | `PRJ-EPIC-015` | Inventory and reservations | ⬜ Not started | — | Implement append-only movements, balances, adjustments, expiring reservations, accepted-stock activation, transaction/locking rules, background expiry worker, warehouse/customer availability UI, invariant and last-unit concurrency tests | Warehouse inspection |
| 16 | `PRJ-EPIC-016` | Shopping cart | 🟡 Partial | Client-side cart context, cart drawer/page, quantity controls, line items, totals, persistence behavior, and unit tests exist | Add server/anonymous cart strategy, API persistence where required, merge-on-login rules, active-product revalidation, backend-calculated display totals, real availability errors, and E2E tests | Public catalogue, customer authentication |
| 17 | `PRJ-EPIC-017` | Checkout and idempotency | 🟡 Partial | Responsive checkout form and order-summary interface exist | Implement address selection, server-side price reload/calculation, transactional reservations, pending order/payment attempt, `Idempotency-Key` handling, frontend API integration, rollback/duplicate/concurrency tests, and OpenAPI | Inventory, cart, customer addresses |
| 18 | `PRJ-EPIC-018` | Orders and customer history | 🟡 Partial | Account order-list/detail/tracking mock screens exist | Add order and immutable line-snapshot schema, totals and transition rules, customer ownership, cancel behavior, warehouse prepare flow, real account API integration, tests, and OpenAPI | Checkout |
| 19 | `PRJ-EPIC-019` | Payments | ⛔ Blocked | Payment abstraction requirements are documented | Obtain explicit provider and business-rule approval; then implement provider adapter, internal payment/events schema, trusted amount verification, signed/idempotent webhooks, retry/status UI, failure/release flows, and tests without storing card data | Orders, inventory; provider decision |
| 20 | `PRJ-EPIC-020` | Shipments and tracking | ⛔ Blocked | Manual or adapter-based shipment requirements and mock tracking UI are documented | Confirm operational carrier/rules or choose the documented manual MVP flow; then implement shipment/events schema, warehouse handover, customer tracking, callbacks where applicable, audit/outbox events, tests, and OpenAPI | Paid orders; shipping decision |
| 21 | `PRJ-EPIC-021` | Custom orders | ⬜ Not started | Requirements and recommended states are documented | Confirm MVP inclusion depth; implement request/message/quote schema, secure attachments, ownership/Casbin, quote acceptance, frontend customer/artisan experiences, notifications, transactions, tests, and OpenAPI | Artisan profiles, authentication, files; payments if quotes are payable |
| 22 | `PRJ-EPIC-022` | Notifications and background workers | 🟡 Partial | Outbox schema and password-reset outbox producer exist | Add worker runtime, reliable claiming, retry/backoff, delivery status, notification records, configured adapters, reservation-expiry jobs, health/metrics, idempotent consumers, and integration tests | Database foundation and producing workflows |
| 23 | `PRJ-EPIC-023` | Audit and administration | 🟡 Partial | Append-only audit schema and database mutation protection exist | Add application audit service, records for every required sensitive action, administrator query API/UI, Casbin policies, filters/pagination, privacy controls, and integration tests | Protected business workflows |
| 24 | `PRJ-EPIC-024` | File and media security | 🟡 Partial | MinIO client, private/public/artisan buckets, persistent volume, initialization, and media metadata schema exist | Implement content-signature validation, configured size limits, generated object keys, private authorized reads, signed URLs, moderated publication, optional scanning hook, cleanup/rollback, and MinIO integration tests | Authentication, product/artisan workflows |
| 25 | `PRJ-EPIC-025` | Accessibility and responsive quality | 🟡 Partial | Semantic responsive storefront components, keyboard-oriented controls, visible form errors, mobile layouts, and RTL support exist | Complete WCAG audit for keyboard/focus/labels/contrast/live regions, reduced-motion checks, screen-reader verification, responsive tests across all operational dashboards, and automated accessibility checks | All frontend workflows |
| 26 | `PRJ-EPIC-026` | Security hardening | 🟡 Partial | Password/token hashing, session revocation, Casbin foundation, rate limits, safe error envelopes, security headers, CORS configuration, private data network, and production credential checks exist | Complete threat-model checklist, least-privilege DB/migration users, Redis authentication outside local, secret scanning, upload security, log redaction tests, dependency/container scanning, TLS/HSTS deployment, and authorization bypass E2E tests | All protected workflows, deployment |
| 27 | `PRJ-EPIC-027` | CI, integration, and E2E quality gates | 🟡 Partial | Jenkins stages, backend/frontend lint/type/test/build, OpenAPI validation, migrations, Docker builds, frontend Jest tests, Playwright storefront tests, and Make targets exist | Run PostgreSQL/Redis/MinIO integrations in CI, add complete business E2E/concurrency suites, container scanning, coverage/quality reporting, staging smoke tests, and block release on failures | Implemented workflows |
| 28 | `PRJ-EPIC-028` | Staging and production operations | 🟡 Partial | Development/production Compose files, multi-stage images, non-root targets, health checks, Nginx routing, persistent PostgreSQL/Redis/MinIO volumes, and deployment guidance exist | Provision staging/production, configure secret management and TLS, add worker, immutable image publishing, backup/restore drills, migration prechecks, monitoring/alerts, smoke tests, manual approval, and rollback automation | CI quality gates, all MVP workflows |
| 29 | `PRJ-EPIC-029` | MVP release readiness | ⬜ Not started | — | Resolve every launch-blocking open question, complete included MVP workflows, run required E2E/security/accessibility/performance tests, verify backups and rollback, approve staging, publish immutable release, and complete production smoke tests | All MVP epics |

## Current Sprint Recommendation

| Priority | Jira Issue | Work item | Status | Acceptance criteria |
|---:|---|---|---|---|
| 1 | `PRJ-701` | Implement public catalogue backend queries | ⬜ Not started | Five public catalogue endpoints return only publishable localized DTOs with bounded pagination and repository/handler tests |
| 2 | `PRJ-702` | Integrate storefront with catalogue API | ⬜ Not started | Home, product, category, search, and artisan surfaces load real API data with accessible loading, empty, and error states |
| 3 | `PRJ-703` | Integrate frontend authentication | ⬜ Not started | Register/login/refresh/logout/forgot/reset and `/me` work through the real API with secure session handling and field/toast errors |
| 4 | `PRJ-704` | Run authentication infrastructure tests in CI | ⬜ Not started | PostgreSQL and Redis integration tests execute rather than skip, and failures block the pipeline |
| 5 | `PRJ-705` | Complete public/auth E2E coverage | ⬜ Not started | Four-locale catalogue smoke tests and authentication success/failure/authorization scenarios pass in Playwright |

After the current sprint, proceed to customer profile/addresses and artisan onboarding. Do not begin payments or shipping integrations until their external decisions are confirmed.

## Epic Completion Checklist

Apply this checklist before marking any project epic implemented:

- [ ] The implemented scope matches the approved requirements and excludes unapproved business rules.
- [ ] Frontend screens use real APIs unless the epic is explicitly UI-only.
- [ ] All visible text is translated in English, French, Arabic, and Spanish.
- [ ] Arabic layout works correctly in RTL.
- [ ] Responsive layouts work on mobile, tablet, and desktop.
- [ ] Forms provide accessible client feedback and preserve backend field errors.
- [ ] Backend handlers remain thin and never query PostgreSQL directly.
- [ ] Application/domain code owns authorization, ownership, invariants, and state transitions.
- [ ] Every protected action enforces authentication and an explicit Casbin permission.
- [ ] SQLBoiler models stay inside infrastructure repositories and never become API responses.
- [ ] Money uses integer minor units and trusted totals are calculated by the backend.
- [ ] Critical state changes use explicit transactions, audit records, and outbox events where required.
- [ ] Schema changes have paired migrations, constraints, indexes, deletion behavior, and regenerated models.
- [ ] File uploads are private by default and validated by content signature and configured size.
- [ ] Unit, integration, concurrency, E2E, accessibility, and security tests appropriate to the epic pass.
- [ ] OpenAPI and frontend API contracts are synchronized and validated.
- [ ] `gofmt`, `go vet ./...`, `go test ./...`, `go build ./...`, frontend lint/type-check/test/build, and relevant Playwright suites pass.
- [ ] Documentation, environment examples, deployment steps, and Jira statuses are updated.

## External Decisions Register

| Decision | Status | Rule |
|---|---|---|
| Payment provider | Unconfirmed | Do not implement a real provider until explicitly approved. |
| Shipping provider and tariffs | Unconfirmed | Do not hard-code a carrier, tariff, or destination rule. |
| Supported launch countries | Unconfirmed | Do not restrict or promise destinations without approval. |
| Supported currencies | Unconfirmed | Keep currency validation configuration-driven. |
| Tax rules | Unconfirmed | Do not calculate or display invented tax behavior. |
| Commission and artisan payout | Unconfirmed | Do not implement split payments or payout calculations. |
| Refund window | Unconfirmed | Do not enforce an invented time window. |
| Cash on delivery | Unconfirmed | Do not expose it as an option without approval. |
| Multi-artisan orders/shipments | Unconfirmed | Do not assume splitting or aggregation behavior. |
| Maximum upload sizes | Unconfirmed | Use configuration only after an approved operational limit is supplied. |
| Artisan legal documents | Unconfirmed | Do not invent mandatory legal-document types. |
