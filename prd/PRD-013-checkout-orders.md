# PRD-13 · Checkout, Orders, Cancellation & Returns

| | |
| --- | --- |
| **Epic** | 13 — Checkout, Orders, Cancellation & Returns |
| **Stories** | US-CHECK-001 … US-CHECK-003; US-ORD-001 … US-ORD-003; US-CAN-001 … US-CAN-002; US-RET-001 |
| **Priority** | Critical |
| **Status** | 🟡 Partial |
| **Surfaces** | `apps/web/app/[locale]/checkout`, `account/orders`, order detail |
| **API** | `POST /api/v1/checkout`, `GET /api/v1/orders`, `GET /api/v1/orders/:id`, `POST /api/v1/orders/:id/cancel` |

## 1. Summary

Turn a validated cart into a pending order with trusted totals, reservations, immutable line snapshots, payment intent, ownership-safe history, cancellation, and return handling.

## 2. Problem

Checkout is a high-risk boundary: client prices can be manipulated, concurrent requests can oversell stock, and buyer and seller views can expose unrelated private order data.

## 3. Goals

- Recalculate product, workshop, artisan, address, availability, and totals on the backend.
- Create reservations, pending orders, and payment attempts transactionally and idempotently.
- Preserve immutable order history and separate buyer from seller visibility.
- Support eligible cancellation and return records without silently changing history.

### Non-goals

- Final payment-provider behavior, tax, shipping tariffs, and split-shipment rules before they are approved.

## 4. Users

- Customer checking out and viewing their order.
- Artisan purchasing as a customer.
- Artisan viewing only related seller-side order items.
- Warehouse, support, and finance actors with scoped access.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Require authentication, reload trusted data, validate address, calculate totals, reserve stock, create pending order, and create payment attempt. | US-CHECK-001 |
| R2 | Protect final-unit checkout concurrency so only one reservation succeeds. | US-CHECK-002 |
| R3 | Preserve product, artisan, workshop, price, currency, quantity, address, and totals as immutable snapshots. | REQ-ORDER-001, US-ORD-001 |
| R4 | Let customers read only their own orders and artisans read only permitted order-item data. | US-ORD-001..003 |
| R5 | Make checkout, cancellation, and return operations idempotent and state-validated. | US-CAN-001..002, US-RET-001 |
| R6 | Release unpaid reservations on eligible cancellation or expiry. | US-CAN-001 |
| R7 | Prevent shipped orders from being silently cancelled; returned stock requires inspection. | US-CAN-002, US-RET-001 |

## 6. Flow

```text
Cart → auth/address → trusted totals → reservation → pending order → payment
Paid → preparing → ready to ship → shipped → delivered
Eligible cancel/return → validated transition → release/refund/inventory decision
```

## 7. Technical notes

- Use `Idempotency-Key` for checkout and other duplicate-sensitive mutations.
- Migration `000014_moderation_orders` adds orders, immutable line/address snapshots, stock reservations, payment attempts, returns, and shipment events.
- `POST /checkout`, buyer order reads, seller item reads, cancellation, and authorized return recording are connected through the authenticated API and Next.js BFF routes.
- Buyer order detail reads persisted shipment events, cancellation is available for eligible pending/paid orders, and repeated cancellation/return operations are safe.
- Application services own transaction boundaries and state transitions.
- Seller visibility is item-scoped; seller status never grants full buyer-order access.

## 8. Success metrics

- Checkout completion and duplicate-request prevention rate.
- Zero negative stock or duplicate orders under concurrency tests.
- Correct buyer/seller authorization outcomes.

## 9. Risks & open questions

- Payment, tax, shipping, return-window, and multi-artisan order policies are open.
- Partial fulfilment and split-shipment behavior must be decided before final order states.
- Reservation expiry cleanup is handled by the inventory worker. Real payment confirmation/refunds, shipping tariffs, and tax remain intentionally deferred to their owning epics.

## Source traceability

`REQ-CHECKOUT-001`, `REQ-ORDER-001..002`, `US-CHECK-001..003`, `US-ORD-001..003`, `US-CAN-001..002`, `US-RET-001`.
