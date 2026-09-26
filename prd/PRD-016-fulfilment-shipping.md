# PRD-16 · Fulfilment & Shipping

| | |
| --- | --- |
| **Epic** | 16 — Fulfilment & Shipping |
| **Stories** | US-FUL-001; US-SHIP-001 … US-SHIP-002 |
| **Priority** | High |
| **Status** | ⛔ Blocked |
| **Surfaces** | Warehouse orders, shipment creation, customer tracking, artisan purchases |
| **API** | `GET /api/v1/warehouse/orders`, `POST /api/v1/warehouse/orders/:id/prepare`, `POST /api/v1/warehouse/orders/:id/ship`, shipment callback endpoints |

## 1. Summary

Move paid order allocations through pick, pack, shipment handover, carrier events, and customer-visible tracking.

## 2. Problem

A paid order is not complete until stock is physically prepared and handed to a carrier. Missing shipment rules or unscoped tracking can expose private data and create untraceable fulfilment.

## 3. Goals

- Pick and pack only paid, eligible allocations.
- Record carrier/manual shipment data and tracking events.
- Keep customer tracking and artisan purchase tracking scoped and reliable.

### Non-goals

- Hard-coding a carrier, tariff, destination, or split-shipment policy before approval.

## 4. Users

- Warehouse agent picking and packing.
- Fulfilment agent creating shipment handover.
- Customer tracking their order.
- Artisan tracking purchases made through the same account.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Show only paid and fulfilment-eligible orders to warehouse actors. | US-FUL-001 |
| R2 | Prevent picked quantity from exceeding allocation. | US-FUL-001 |
| R3 | Record carrier/manual shipment, tracking reference, status, and handover date. | US-SHIP-001 |
| R4 | Expose permitted shipment events to customers. | US-SHIP-002 |
| R5 | Reject or explicitly handle duplicate and backwards events. | US-SHIP-002 |
| R6 | Make shipment creation idempotent and auditable. | Reliability rules |

## 6. Flow

```text
Paid order → pick → pack → ready to ship → shipment handover → in transit → delivered/failed/returned
```

## 7. Technical notes

- Shipment status transitions belong in the application/domain layer.
- External callbacks require signature verification and idempotency where supported.
- Customer and artisan views expose only permitted tracking data.

## 8. Success metrics

- Paid-order-to-shipment turnaround.
- Shipment events visible without duplicate notifications.
- Zero over-picking or unauthorized tracking reads.

## 9. Risks & open questions

- Carrier/manual workflow, destinations, tariffs, split shipments, and delivery estimates are open decisions.
- Return and failed-delivery transitions must align with inventory and refund policies.

## Source traceability

`REQ-SHIP-001`, `US-FUL-001`, `US-SHIP-001..002`, `docs/deployment.md`.
