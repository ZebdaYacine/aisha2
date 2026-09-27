# PRD-11 · Inventory Ledger & Reservations

| | |
| --- | --- |
| **Epic** | 11 — Inventory Ledger & Reservations |
| **Stories** | US-WH-003; US-WH-004; US-CHECK-002 |
| **Priority** | Critical |
| **Status** | 🟡 Partial |
| **Surfaces** | `apps/web/app/[locale]/admin/inventory`, product availability, checkout |
| **API** | `GET /api/v1/warehouse/inventory`, `POST /api/v1/warehouse/inventory/:productId/adjust`, checkout reservation endpoints |

## 1. Summary

Provide an authoritative, concurrency-safe stock ledger that supports availability, checkout reservations, expiry, fulfilment, and permitted workshop views.

## 2. Problem

Direct quantity editing causes overselling, untraceable corrections, and inconsistent customer availability. Reservations also need safe handling when payments fail or requests race for the last unit.

## 3. Goals

- Use append-only movements as the source of truth.
- Keep available stock non-negative and every correction attributable.
- Reserve accepted stock transactionally with expiry, release, and commit states.
- Prove last-unit correctness under concurrency.

### Non-goals

- Redis as the authoritative inventory store.
- Forecasting or automated replenishment in the MVP.

## 4. Users

- Warehouse agent or inventory manager adjusting stock.
- Artisan viewing permitted workshop inventory.
- Customer seeing authoritative availability.
- Checkout service reserving stock.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Record append-only reception, inspection, reservation, release, commit, shipment, return, and adjustment movements. | US-WH-003 |
| R2 | Track on-hand, available, reserved, quarantined, damaged, rejected, and shipped quantities. | REQ-INV-001 |
| R3 | Require actor and reason for corrections; available stock cannot become negative. | US-WH-003 |
| R4 | Reserve accepted stock transactionally with expiry, release, and commit behavior. | US-CHECK-002 |
| R5 | Protect concurrent purchase of the final unit so only one reservation succeeds. | US-CHECK-002 |
| R6 | Show only permitted inventory to artisans and warehouse users. | US-WH-004 |

## 6. Flow

```text
Accepted stock → available balance → checkout lock/reservation → payment success/expiry → commit/release → shipment
```

## 7. Technical notes

- PostgreSQL movement records are authoritative; balance tables are only read optimizations.
- Use transactions and product row locks for reservations; checkout rechecks availability after the lock is acquired.
- Expired holds are released transactionally at checkout entry and by the standalone reservation-expiry worker. Both paths use row locks and idempotent release movements.
- Returned order items are recorded as `RETURN` movements in the `QUARANTINED` bucket, so a return never becomes sellable without a later inspection decision. Migration `000019_commerce_hardening` also makes return recording one-per-order and concurrency-safe.
- Artisan inventory is filtered by owned workshop; warehouse agents and administrators have operational scope.
- `POST /warehouse/inventory/:productId/adjust` appends an `ADJUSTMENT` movement and writes audit/outbox records in the same transaction.

## 8. Success metrics

- Zero negative available balances.
- Successful last-unit concurrency test with one winner.
- Reconciliation of every displayed balance to movement history.

## 9. Risks & open questions

- Reservation duration is currently the documented 30-minute pending-checkout hold; changing it remains a business decision.
- Partial fulfilment and split-shipment behavior must be defined before final ledger transitions.
- Payment confirmation and fulfilment still need to call the `COMMITTED`/`SHIPPED` transitions once PRD-15/16 are implemented; those epics remain blocked on provider and shipping decisions.

## Source traceability

`REQ-INV-001`, `US-WH-003`, `US-WH-004`, `US-CHECK-002`, `docs/security.md`.
