# PRD-15 · Payment Confirmation & Refunds

| | |
| --- | --- |
| **Epic** | 15 — Payment Confirmation & Refunds |
| **Stories** | US-PAY-001 … US-PAY-005 |
| **Priority** | Critical |
| **Status** | ⛔ Blocked |
| **Surfaces** | Payment step, payment return page, order detail, finance tools |
| **API** | `GET /api/v1/payments/:paymentId`, `POST /api/v1/payments/:paymentId/retry`, `POST /api/v1/webhooks/payments/:provider` |

## 1. Summary

Abstract payment providers behind a trusted adapter, confirm payment through verified provider events, and handle failure, expiry, and refunds without duplicating side effects.

## 2. Problem

A browser redirect or frontend amount cannot prove payment. Duplicate webhooks, mismatched amounts, and incomplete refund state can corrupt orders and inventory.

## 3. Goals

- Use backend-trusted amount and currency for every payment attempt.
- Verify, deduplicate, and audit provider webhooks.
- Model payment failure, expiry, partial refund, and full refund safely.

### Non-goals

- Selecting a provider, currencies, tax rules, or cash-on-delivery policy before explicit approval.

## 4. Users

- Customer starting or returning from payment.
- Finance/support actor handling refunds.
- Payment provider sending signed events.
- Order and inventory services consuming trusted payment state.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Send trusted backend amount and currency to the provider adapter. | US-PAY-001 |
| R2 | Treat browser return as informational; only trusted verification/webhook can confirm payment. | US-PAY-002 |
| R3 | Verify signature/timestamp, match order/amount/currency, and reject unknown or duplicate events. | US-PAY-003 |
| R4 | Release reservations on payment failure or expiry. | US-PAY-004 |
| R5 | Validate refund amount and update payment/order state transactionally. | US-PAY-005 |

## 6. Flow

```text
Pending order → create provider attempt → browser return may be pending
Provider event → verify → idempotency → match amount/currency → capture → commit reservation
Failure/expiry → release reservation; refund → update payment/order/audit
```

## 7. Technical notes

- Store provider events separately with a unique external event ID.
- Never store raw card data or log provider secrets.
- Use stable errors such as `INVALID_WEBHOOK_SIGNATURE`, `DUPLICATE_WEBHOOK`, and `PAYMENT_AMOUNT_MISMATCH`.

## 8. Success metrics

- Verified payment success rate.
- Duplicate webhook side effects equal zero.
- Payment/order/reservation states reconcile after failure and refund tests.

## 9. Risks & open questions

- Provider, credentials, supported currencies, refund window, tax, and cash-on-delivery are unconfirmed; this blocks implementation.
- Provider-specific retry and dispute behavior needs a signed contract.

## Source traceability

`REQ-PAY-001`, `US-PAY-001..005`, `docs/security.md`, `docs/backend.md`.
