# PRD-01 · Product Definition & Architecture

| | |
| --- | --- |
| **Epic** | 1 — Product Definition & Architecture |
| **Stories** | Product vision, roles, MVP boundaries, platform rules |
| **Priority** | High |
| **Status** | ✅ Implemented |
| **Surfaces** | `docs/requirements.md`, `docs/architecture.md` |
| **API** | Cross-cutting contract; no standalone endpoint |

## 1. Summary

Define AISHA as a multilingual marketplace for Algerian handmade and artistic products. The product keeps artisan identity, provenance, warehouse inspection, managed stock, fulfilment, and shipment visible across the shopping journey.

## 2. Problem

A generic marketplace cannot provide enough trust for culturally significant products or physically managed stock. Customers need provenance and quality context; operators need a controlled path from product draft to delivery.

## 3. Goals

- Keep artisan, workshop, product, and provenance information connected.
- Preserve one account model where customer capabilities remain available to artisans.
- Separate moderation approval from warehouse stock acceptance.
- Define a traceable flow from approved product to shipment.

### Non-goals

- Advanced promotions, marketplace syndication, multi-warehouse routing, automated payouts, or a native mobile app in the MVP.

## 4. Users

- Visitor discovering products and artisans.
- Customer browsing, purchasing, and tracking orders.
- Artisan managing workshops and product drafts.
- Moderator, warehouse agent, and administrator operating protected workflows.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Support visitor, customer, artisan, moderator, warehouse, and administrator capabilities. | US-VIS-001, US-AUTH-001, US-ART-001, US-ADMIN-001 |
| R2 | Preserve customer capabilities when artisan membership is active. | US-ART-004, US-CHECK-003 |
| R3 | Keep product approval and warehouse acceptance as separate controls. | US-PROD-005, US-WH-002 |
| R4 | Enforce clean boundaries between handlers, application services, repositories, DTOs, and persistence models. | Architecture rules |
| R5 | Preserve the managed-stock flow: approved → received → inspected → accepted → reserved → paid → fulfilled → shipped. | US-WH-001, US-CHECK-001, US-FUL-001 |

## 6. Flow

```text
Visitor → catalogue → customer account → cart → checkout
Artisan → workshop → product draft → moderation → warehouse inspection
Accepted stock → reservation → trusted payment → fulfilment → shipment
```

## 7. Technical notes

- PostgreSQL is authoritative for permanent business data.
- Backend authorization, ownership, validation, audit, and idempotency are authoritative.
- SQLBoiler models remain inside repositories; public responses use DTOs.

## 8. Success metrics

- Percentage of active products with complete provenance data.
- Percentage of stock that follows the inspection path before sale.
- Checkout-to-delivery traceability across order records.

## 9. Risks & open questions

- Payment provider, shipping rules, tax, commission, and payout decisions remain unconfirmed.
- Multiple artisans per order and split-shipment rules must be decided before fulfilment implementation.

## Source traceability

`docs/requirements.md`, `docs/userstory.md`, `docs/architecture.md`, `docs/security.md`.
