# PRD-10 · Warehouse Reception & Inspection

| | |
| --- | --- |
| **Epic** | 10 — Warehouse Reception & Inspection |
| **Stories** | US-WH-001; US-WH-002 |
| **Priority** | High |
| **Status** | ✅ Shipped |
| **Surfaces** | `apps/web/app/[locale]/admin/warehouse` |
| **API** | `GET /api/v1/warehouse/products`, `POST /api/v1/warehouse/receptions`, `GET /api/v1/warehouse/receptions`, `POST /api/v1/warehouse/receptions/:id/inspect` |

## 1. Summary

Record physical stock arriving at AISHA and make quality decisions before any quantity becomes available for sale.

## 2. Problem

Received goods are not automatically acceptable goods. Without a controlled inspection record, customers may buy stock that has not passed quality checks and operators cannot explain inventory decisions.

## 3. Goals

- Record reception, batch, quantity, actor, date, and evidence.
- Keep received stock unavailable until inspection succeeds.
- Classify every inspected unit and preserve the quantity equation.

### Non-goals

- Automated visual quality grading or supplier payment automation.

## 4. Users

- Warehouse agent receiving a parcel or batch.
- Inventory manager reviewing inspection outcomes.
- Artisan or administrator viewing permitted stock context.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Record product, artisan/supplier, quantity, reception reference, date, agent, batch, notes, and evidence. | US-WH-001 |
| R2 | Keep received quantity in `RECEIVED_PENDING_INSPECTION`. | US-WH-001 |
| R3 | Record accepted, rejected, quarantined, and damaged quantities. | US-WH-002 |
| R4 | Enforce `accepted + rejected + quarantined + damaged = inspected quantity`. | US-WH-002 |
| R5 | Use authorization, transactions, audit events, and safe duplicate handling. | US-WH-001..002 |
| R6 | Let an authorized warehouse agent search by artisan phone, choose a workshop, and select a moderator-validated product code before receiving stock. | US-WH-001 |
| R7 | After accepted inspection quantity is committed, automatically activate an approved product when artisan, workshop, price, media, and stock gates pass; otherwise keep it unavailable. | US-WH-002 |

## 6. Flow

```text
Receive by approval code → pending inspection → enter outcomes/evidence → validate totals → commit inventory movements → activate approved product when ready
```

## 7. Technical notes

- Only accepted quantity can feed availability; rejected, quarantined, and damaged units use separate stock buckets.
- Inspection writes are transactional and replay-safe. Reception references are unique idempotency keys and repeated inspections return the committed outcome without duplicating movements.
- Evidence uses generated private MinIO object keys and short-lived signed reads. Evidence is required before inspection can be committed.
- Migration `000016_warehouse_reception_inspection` adds reception, inspection, evidence, and stock-bucket persistence. The API exposes list/create reception, inspect, and evidence upload/read operations.
- Migration `000020_product_codes` adds a unique warehouse-facing product code. The validated-product search joins product, artisan, phone, and workshop context and is protected by the warehouse reception read permission. The UI never asks an agent to enter or copy an internal product UUID.
- The inspection transaction updates accepted stock, product status, media visibility, audit, and outbox records together. A product is activated only from `APPROVED`; an already active product is unchanged.

## 8. Success metrics

- Percentage of received stock inspected before availability.
- Invalid quantity submissions rejected.
- Inspection records with complete actor, reason, and evidence data.

## 9. Risks & open questions

- Evidence retention and warehouse exception policy need explicit operating rules.
- Reinspection of returned stock must be coordinated with returns and inventory epics.

## Source traceability

`REQ-WH-001`, `REQ-WH-002`, `US-WH-001`, `US-WH-002`.
