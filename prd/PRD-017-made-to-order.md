# PRD-17 · Made-to-Order Requests

| | |
| --- | --- |
| **Epic** | 17 — Made-to-Order Requests |
| **Stories** | US-CUSTOM-001 … US-CUSTOM-003 |
| **Priority** | Medium |
| **Status** | ⬜ Not Started |
| **Surfaces** | `apps/web/app/[locale]/custom-orders`, artisan workspace |
| **API** | `POST /api/v1/custom-orders`, `GET /api/v1/custom-orders`, `GET /api/v1/custom-orders/:id`, quote/message endpoints |

## 1. Summary

Let customers request custom work from an artisan, exchange clarification, receive an immutable quote, and accept it into a trusted checkout path.

## 2. Problem

Custom work needs private context, clear ownership, versioned quotes, and explicit acceptance. A simple product cart cannot safely represent scope, dates, files, or artisan responses.

## 3. Goals

- Capture a structured, private custom request.
- Restrict visibility to authorized customer, artisan, and staff participants.
- Support clarification, decline, quote, expiry, and acceptance states.

### Non-goals

- Automated production planning or payment before quote acceptance.

## 4. Users

- Customer requesting custom work.
- Artisan responding for an owned profile, workshop, or eligible product.
- Support or administrator reviewing permitted request history.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Capture artisan/workshop/product, description, quantity, desired date, references, and contact preference. | US-CUSTOM-001 |
| R2 | Restrict request and attachments to authorized participants. | US-CUSTOM-001 |
| R3 | Let the owning artisan request clarification, decline, or send a quote. | US-CUSTOM-002 |
| R4 | Keep quote versions immutable and include price, scope, expiry, production time, and workshop. | US-CUSTOM-002 |
| R5 | Prevent expired quotes from being accepted and convert accepted quotes into trusted checkout data. | US-CUSTOM-003 |

## 6. Flow

```text
Customer request → artisan clarification/decline/quote → customer accepts valid quote → trusted checkout
```

## 7. Technical notes

- Use private MinIO objects and signed URLs for reference files.
- Quote acceptance requires ownership, expiry, and state validation.
- Messages and quote changes should emit reliable notifications.

## 8. Success metrics

- Request-to-response time.
- Quote acceptance rate before expiry.
- Unauthorized request/file access attempts rejected.

## 9. Risks & open questions

- Confirm MVP depth, messaging retention, attachment limits, and production/payment handoff.
- Multi-artisan or workshop routing requires explicit ownership rules.

## Source traceability

`REQ-CUSTOM-001`, `US-CUSTOM-001..003`, `docs/security.md`.
