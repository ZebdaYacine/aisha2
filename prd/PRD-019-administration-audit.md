# PRD-19 · Administration, Suspension & Audit

| | |
| --- | --- |
| **Epic** | 19 — Administration, Suspension & Audit |
| **Stories** | US-ADMIN-001 … US-ADMIN-005 |
| **Priority** | Critical |
| **Status** | 🟡 Partial |
| **Surfaces** | `apps/web/app/[locale]/admin`, moderation, audit, verification queues |
| **API** | Admin role, suspension, audit, and verification-review endpoints |

## 1. Summary

Give authorized administrators, moderators, warehouse, finance, support, and auditors the least-privilege tools needed to control the marketplace and trace sensitive actions.

## 2. Problem

Administrative actions can stop sales, expose private records, or alter access. Without scoped permissions, reasons, and append-only audit events, the platform cannot be safely operated.

## 3. Goals

- Separate customer capability, artisan membership, and privileged roles.
- Support reasoned suspension and verification decisions without deleting history.
- Provide filtered, protected audit access.

### Non-goals

- Unbounded administrator access or silent data correction.

## 4. Users

- Administrator managing roles and account/entity status.
- Moderator reviewing products.
- Verification reviewer handling optional artisan evidence.
- Auditor searching protected event history.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Manage roles and permissions through backend authorization. | US-ADMIN-001 |
| R2 | Suspend users, artisans, workshops, and products with reasons and preserved history. | US-ADMIN-002..003 |
| R3 | Search audit events by actor, target, event type, date, user, workshop, product, order, and correlation ID. | US-ADMIN-004 |
| R4 | Review optional artisan verification privately with actor, reason, timestamp, and decision audit. | US-ADMIN-005 |
| R5 | Deny lower-privilege access to administrative data and actions. | US-ADMIN-001..005 |
| R6 | Present users, artisan applications, product moderation, and event audit as nested dashboard sections with consistent paginated rows and Details modals. | US-ADMIN-001..005 |
| R7 | Let administrators update their name, phone, and personal account information in a dashboard modal, including a secure password change. | US-ADMIN-001 |
| R8 | Keep administrator accounts out of artisan onboarding, artisan routing, and artisan self-service capabilities. | US-ADMIN-001 |
| R9 | Present approved artisans as a shop/status/action table and keep internal IDs and non-actionable fields out of the Details modal. | US-ADMIN-001..003 |
| R10 | Show authorized artisan user media in artisan Details and allow administrators to delete user or product media from the Media section. | US-ADMIN-004..005 |

## 6. Flow

```text
Authorized admin → scoped queue/search → validate action and reason → transaction → audit/outbox → visible result
```

## 7. Technical notes

- Casbin covers role/action/resource; application services enforce ownership and state rules. Administrator routes now support reasoned user status changes (`ACTIVE`, `SUSPENDED`, `DISABLED`) and workshop lifecycle controls (`ACTIVE`, `INACTIVE`, `ARCHIVED`), each with transactional session/access effects, audit events, and outbox records.
- The admin UI uses nested overview, users, artisan applications, product moderation, and audit sections. Operational rows use a shared Details-modal pattern; the profile and password controls stay in a dashboard modal.
- Audit records are append-only and include correlation IDs. The API supports actor, target, correlation, and time-range filters. Administrator account summaries explicitly return `NOT_STARTED` artisan state, and artisan self-service services reject administrator principals.
- Private documents are never exposed through public admin shortcuts; administrator document reads use authorized signed URLs.
- Artisan application Details includes authorized uploaded-media previews when available. Approved-artisan rows use shop name, status badges, and right-aligned actions; Details shows only public shop/status information and management controls.
- The admin Media section provides a confirmation-protected delete action for user and product media; deletion removes the database record, records an audit event, and removes the private MinIO object.

## 8. Success metrics

- Protected administrative action success/error rate.
- Audit search completeness and event correlation rate.
- Unauthorized admin attempts rejected and logged.

## 9. Risks & open questions

- Exact role matrix and support scope need approval.
- Audit retention, export, and privacy rules need an operational owner.

## Source traceability

`REQ-AUDIT-001`, `US-ADMIN-001..005`, `docs/security.md`.
