# PRD-09 · Product Moderation & Sellable Activation

| | |
| --- | --- |
| **Epic** | 9 — Product Moderation & Sellable Activation |
| **Stories** | US-PROD-005; US-PROD-007 |
| **Priority** | High |
| **Status** | ✅ Shipped |
| **Surfaces** | `apps/web/app/[locale]/admin/moderation`, artisan submission feedback |
| **API** | `GET /api/v1/admin/product-submissions`, approve/request-changes/reject/suspend endpoints |

## 1. Summary

Ensure that only approved, active, workshop-owned products with accepted stock become publicly purchasable.

## 2. Problem

Approval alone does not prove that a product is in stock, sellable, or safe to expose. Uncontrolled transitions can create unavailable products and erase moderation history.

## 3. Goals

- Provide an auditable moderator queue and decision workflow.
- Require reasons for rejection, requested changes, and suspension.
- Gate activation on moderation, membership, workshop, media, price, and accepted-stock rules.

### Non-goals

- Automated moderation or public product editing by moderators.

## 4. Users

- Moderator reviewing submitted product versions.
- Artisan receiving decisions and resubmitting changes.
- Customer seeing only valid active products.
- Administrator auditing protected decisions.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Support `DRAFT`, `PENDING_REVIEW`, `CHANGES_REQUESTED`, `APPROVED`, `ACTIVE`, `SUSPENDED`, and `ARCHIVED`. | US-PROD-005..007 |
| R2 | Require a reason for rejection, requested changes, and suspension. | US-PROD-005 |
| R3 | Prevent publication before approval. | US-PROD-005 |
| R4 | Activate only when workshop, artisan membership, price, media, and accepted-stock gates pass. | US-PROD-007 |
| R5 | Emit audit, outbox, and notification events for protected decisions. | US-PROD-005 |
| R6 | Show product media in the moderation Details modal through authorized signed URLs without exposing internal product identifiers. | US-PROD-005 |
| R7 | Generate a stable `AISHA-XXXXXXXXXX` warehouse product code when a moderator approves a product. | US-PROD-005 |
| R8 | Keep an approved product unavailable until a warehouse inspection records accepted stock; the warehouse flow may then activate it atomically. | US-PROD-007 |

## 6. Flow

```text
Submission → moderator queue → approve/request changes/reject → approval code → warehouse receive/inspect → activation gates → public active product
```

## 7. Technical notes

- State transitions belong in domain/application services, not handlers.
- Casbin and ownership checks protect queue actions.
- Product media becomes public only after approved activation.
- Migration `000014_moderation_orders` stores moderator decisions; the API exposes a status-filtered, paginated queue and decision endpoint. Approval leaves a product `APPROVED` when stock or another gate is missing; `ACTIVATE` requires every gate and publishes atomically.
- Decisions append `audit_events` and `outbox_events` in the same transaction. The outbox event is the notification-delivery boundary; worker delivery remains a cross-cutting follow-up.
- Repository integration coverage verifies queue filtering, approval, activation-gate rejection, media publication, and audit/outbox writes.
- The moderator queue returns product names and media metadata; the Details modal renders image/video/file previews from short-lived signed URLs.
- Warehouse inspection reevaluates the same activation gates in its transaction. If accepted stock is available for an approved product, the product moves to `ACTIVE`, its media becomes public, and audit/outbox events record the automatic activation; otherwise it remains `APPROVED`.
- Approval assigns a stable, unique warehouse-facing product code derived from the product identity. Warehouse workflows use this code for human identification and do not require operators to enter an internal UUID.

## 8. Success metrics

- Median moderation turnaround time.
- Percentage of active products satisfying all activation gates.
- Invalid transitions and unauthorized decisions rejected.

## 9. Risks & open questions

- A product with no accepted stock must remain unavailable and must not be activated.
- Reviewer assignment and notification channels require operational decisions.
- The current activation implementation records public media visibility against the existing object key; a later storage publication step may be needed if private/public buckets are separated in production.

## Source traceability

`REQ-MOD-001`, `US-PROD-005`, `US-PROD-007`, product state machine in `docs/userstory.md`.
