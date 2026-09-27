# PRD-07 · Artisan Membership, Onboarding & Workshops

| | |
| --- | --- |
| **Epic** | 7 — Artisan Membership, Onboarding & Workshops |
| **Stories** | US-ART-001 … US-ART-008; US-WS-001 … US-WS-006 |
| **Priority** | High |
| **Status** | ✅ Shipped |
| **Surfaces** | `apps/web/app/[locale]/artisan`, `apps/web/app/[locale]/admin/artisans` |
| **API** | `POST /api/v1/artisan-applications`, `GET /api/v1/artisan-applications/me`, `PATCH /api/v1/artisan/profile`, `POST /api/v1/artisan/membership/activate`, `GET/POST /api/v1/artisan-applications/me/documents`, `GET/POST/PATCH/DELETE /api/v1/artisan/profile/media`, workshop CRUD/status endpoints, verification and admin membership endpoints |

## 1. Summary

Let an existing customer become an artisan, create the mandatory first workshop, manage additional workshops, and participate in controlled approval and verification workflows.

## 2. Problem

Seller onboarding is unsafe when it creates duplicate accounts, allows products without an owned workshop, or exposes private verification data.

## 3. Goals

- Add artisan membership to the existing user account.
- Require and atomically create a valid first workshop.
- Support workshop ownership, lifecycle management, optional verification, and audited administration.

### Non-goals

- Automated artisan payouts or mandatory verification documents unless a later policy requires them.

## 4. Users

- Registered customer applying to become an artisan.
- Active artisan managing profiles and workshops.
- Reviewer or administrator approving, suspending, or auditing access.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Reuse the existing user account and resume an onboarding draft. | US-ART-001 |
| R2 | Collect validated public, operational, and private artisan information. | US-ART-002 |
| R3 | Require a valid first workshop before membership becomes `ACTIVE`; approval activates the submitted workshop as the default workshop. | US-ART-003..004 |
| R4 | Allow an artisan to create, update, activate, deactivate, and delete only eligible owned workshops. | US-WS-001..005 |
| R5 | Prevent cross-artisan reads and writes. | US-WS-006 |
| R6 | Support optional verification states and private documents. | US-ART-005..006 |
| R7 | Preserve customer capabilities when artisan membership is suspended. | US-ART-008 |
| R8 | Present owned workshops in a responsive management table with actions in the first column, one Add workshop action, read-only Details, and protected Update/Delete dialogs. | US-WS-001..005 |
| R9 | Group the approved artisan workspace into Workshops, Private files, and Product authoring tabs, hiding seller-only tabs when membership capabilities are unavailable. | US-ACC-001, US-ART-004, US-ART-008 |
| R10 | Present private media in rows with name, upload date, View, Update, and Delete actions; Add and Update use modal forms with a media-kind combobox and local file picker. | US-ART-005..006 |

## 6. Flow

```text
Customer → Become an Artisan → profile + first workshop → administrator approval → ACTIVE membership + default workshop → artisan dashboard
                                                     └─ optional verification/review
```

## 7. Technical notes

- Migration `000015_artisan_memberships_verification` adds one membership and one verification state per artisan profile, backfilling existing approved profiles.
- Approval creates or reactivates the ACTIVE membership, promotes the submitted workshop to the active default workshop, and records verification, role, audit, and outbox state in one transaction. `POST /artisan/membership/activate` remains idempotent for legacy approved profiles that do not yet have a membership.
- Workshop CRUD/status operations re-check membership status, ownership, default-workshop protection, and product/order history in the repository.
- Public catalogue, product authoring, moderation activation, and checkout require an ACTIVE membership; suspension leaves customer capabilities intact.
- Casbin permission is combined with ownership checks.
- Application documents are sent as multipart form data to `/artisan-applications/me/documents`; profile media uses `/artisan/profile/media`. The API validates file signatures and size, generates an object key, writes to the configured private MinIO bucket (`MINIO_ARTISAN_BUCKET`, default `artisan-private`), persists metadata only after the object write succeeds, cleans up the object on database failure, and returns a 15-minute authorized signed URL. Successful document uploads move verification from NOT_SUBMITTED or CHANGES_REQUESTED to PENDING and enqueue a reviewer notification event.
- Profile media replacement and deletion are owner-scoped. Replacement writes a new generated MinIO object before atomically updating metadata, then removes the old object; deletion removes the metadata and private object. Media responses include `createdAt` for the row upload date and never expose object keys.
- The responsive artisan workspace covers activation, tabular workshop CRUD, first-column lifecycle actions, Details/Update/Delete dialogs, verification state, and suspended read-only behaviour; the administrator screen covers membership suspension/reactivation and verification decisions.
- Unit, OpenAPI, migration, repository, ownership, and dependency-backed concurrency tests cover the workflow.

## 8. Success metrics

- Onboarding completion rate and duplicate-submission rate.
- Percentage of active artisans with a valid owned workshop.
- Cross-owner access attempts rejected and audited.

## 9. Risks & open questions

- Verification is intentionally optional: reviewer decisions and user-facing reasons are available through the API, while document policy and reviewer UI remain follow-up work.
- Admin membership and verification screens still need a complete responsive UI and deeper transition/concurrency coverage.
- Suspension rules distinguish artisan selling access from customer access; public listings and seller mutations are blocked while customer profile, address, and purchase capabilities remain available.

## Source traceability

`REQ-ART-001`, `US-ART-001..008`, `US-WS-001..006`, `US-ADMIN-001..005`.
