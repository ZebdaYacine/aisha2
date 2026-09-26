# PRD-18 · Notifications & Background Workers

| | |
| --- | --- |
| **Epic** | 18 — Notifications & Background Workers |
| **Stories** | US-NOTIF-001 … US-NOTIF-003 |
| **Priority** | High |
| **Status** | 🟡 Partial |
| **Surfaces** | Account notifications, artisan workspace, worker runtime |
| **API** | Notification read/list endpoints; worker and outbox consumers |

## 1. Summary

Reliably deliver account, membership, artisan, commerce, shipment, and custom-order notifications through an outbox-backed worker.

## 2. Problem

Direct notification sends can be lost when a transaction fails or duplicated when providers retry. Users also need private, role-appropriate messages rather than raw event payloads.

## 3. Goals

- Publish notifications only after successful domain commits.
- Claim, retry, back off, and surface failed outbox work.
- Prevent duplicate delivery and recipient privacy leaks.

### Non-goals

- Real-time chat or an unconfirmed external messaging provider.

## 4. Users

- Customer receiving order and payment updates.
- Artisan receiving product, request, and sales updates.
- Administrator monitoring failed events.
- Worker process consuming outbox records.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Deliver membership, verification, suspension, order, payment, shipment, product, and custom-order events as applicable. | US-NOTIF-001..003 |
| R2 | Use an outbox or equivalent reliable mechanism for required domain events. | Reliability requirements |
| R3 | Make event consumption and notification creation idempotent. | US-NOTIF-002..003 |
| R4 | Restrict content to authorized recipients and redact private data. | Security rules |
| R5 | Provide worker health, retry, backoff, and failed-event visibility. | Deployment requirements |

## 6. Flow

```text
Domain transaction → outbox event → worker claim → render localized notification → deliver/store → retry or dead-letter
```

## 7. Technical notes

- Worker uses the same application/domain rules as the API.
- Reservation expiry belongs in the worker or feature-owned job, not a handler.
- Notification templates use translation keys and stable event types.

## 8. Success metrics

- Outbox-to-notification delivery latency.
- Duplicate notification rate.
- Failed events retried and resolved within the operational target.

## 9. Risks & open questions

- Email, in-app, push, and SMS channels need a launch decision.
- Worker scaling and dead-letter retention need operational limits.

## Source traceability

`US-NOTIF-001..003`, reliability requirements in `docs/requirements.md`, `docs/deployment.md`.
