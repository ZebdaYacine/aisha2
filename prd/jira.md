# AISHA PRD Jira Tracker

This is the implementation view for all AISHA product domains. Every PRD uses the same nine-section product format and maps requirements to user stories.

## Status legend

| Status | Meaning |
| --- | --- |
| ✅ Implemented | The scoped workflow or foundation exists with relevant evidence. |
| 🟡 Partial | Some implementation exists, but the complete workflow is not connected. |
| ⬜ Not Started | Requirements are documented; implementation is not connected. |
| ⛔ Blocked | An unresolved external decision prevents implementation or acceptance. |

## Board

| PRD | Epic | Stories | Status | Depends on |
| --- | --- | --- | --- | --- |
| [PRD-01](./PRD-001-product-definition.md) | Product definition and architecture | Vision, roles, MVP rules | ✅ Implemented | — |
| [PRD-02](./PRD-002-platform-foundation.md) | Platform and development foundation | Local, CI, deployment baseline | ✅ Implemented | PRD-01 |
| [PRD-03](./PRD-003-design-i18n.md) | Design, localization, and RTL | US-VIS-001..003, REQ-I18N-001 | ✅ Implemented | PRD-01 |
| [PRD-04](./PRD-004-authentication-authorization.md) | Authentication and authorization | US-AUTH-001..003, REQ-RBAC-001 | ✅ Implemented | PRD-02 |
| [PRD-05](./PRD-005-customer-account.md) | Customer account | US-ACC-001, US-CHECK-001, US-CHECK-003 | ✅ Implemented | PRD-04, PRD-07, PRD-13 |
| [PRD-06](./PRD-006-catalogue-discovery.md) | Public catalogue and discovery | US-VIS-001..003 | ✅ Implemented | PRD-03, PRD-04 |
| [PRD-07](./PRD-007-artisan-membership.md) | Artisan membership and workshops | US-ART-001..008, US-WS-001..006 | ✅ Shipped | Approval now activates the submitted default workshop and membership atomically; the product authoring workshop combobox refreshes when its editor opens, with legacy approved profiles repaired through the admin Activate action |
| [PRD-08](./PRD-008-product-authoring.md) | Product authoring and media | US-PROD-001..004, US-PROD-006, US-PROD-008 | ✅ Shipped | Workshop-owned authoring, quantity/total-price fields, single-selected-language submission, four-media cap, private media, immutable submissions, owner archive, and protected history are implemented and tested |
| [PRD-09](./PRD-009-moderation-activation.md) | Product moderation and activation | US-PROD-005, US-PROD-007 | ✅ Implemented | Moderator approval generates a stable warehouse product code exposed to the artisan; accepted warehouse stock can auto-activate the approved product; PRD-08 |
| [PRD-10](./PRD-010-warehouse-inspection.md) | Warehouse reception and inspection | US-WH-001..002 | ✅ Implemented | Warehouse product search by artisan phone/workshop/code and automatic activation after accepted inspection are connected; PRD-09, PRD-20 |
| [PRD-11](./PRD-011-inventory-reservations.md) | Inventory and reservations | US-WH-003..004, US-CHECK-002 | 🟡 Partial — inventory ledger and workshop view are implemented; checkout reservation concurrency and fulfilment transitions remain | PRD-10 |
| [PRD-12](./PRD-012-cart-wishlist.md) | Cart and wishlist | US-CART-001..002 | ✅ Shipped | Owner-scoped persistence, anonymous merge, wishlist, availability revalidation, and active-product guards are implemented and tested | PRD-04, PRD-06, PRD-11 |
| [PRD-13](./PRD-013-checkout-orders.md) | Checkout, orders, cancellation, returns | US-CHECK-001..003, US-ORD-001..003, US-CAN-001..002, US-RET-001 | 🟡 Partial | PRD-05, PRD-11, PRD-12 |
| [PRD-14](./PRD-014-reviews-engagement.md) | VReviews and engagement | Review and wishlist capabilities | ⬜ Not Started | PRD-13 |
| [PRD-15](./PRD-015-payments.md) | Payment confirmation and refunds | US-PAY-001..005 | ⛔ Blocked | PRD-11, PRD-13; provider decision |
| [PRD-16](./PRD-016-fulfilment-shipping.md) | Fulfilment and shipping | US-FUL-001, US-SHIP-001..002 | ⛔ Blocked | PRD-10, PRD-13; shipping decision |
| [PRD-17](./PRD-017-made-to-order.md) | Made-to-order requests | US-CUSTOM-001..003 | ⬜ Not Started | PRD-07, PRD-20 |
| [PRD-18](./PRD-018-notifications-workers.md) | Notifications and workers | US-NOTIF-001..003 | 🟡 Partial | Notification store/API, outbox consumer with retry/dedupe, reservation worker, and authenticated WebSocket delivery implemented; email/push adapters and operational dead-letter UI remain |
| [PRD-19](./PRD-019-administration-audit.md) | Administration and audit | US-ADMIN-001..005 | 🟡 Partial | Nested tab-panel administration now uses Users, Workers, and Warehouses alongside artisan/product moderation, inventory, media, audit, four-language category CRUD with benefit rates, and paginated order review; backup/restore, explicit multi-warehouse assignment, and deeper retention/concurrency coverage remain |
| [PRD-20](./PRD-020-media-security.md) | File and media security | Media upload and private-file rules | 🟡 Partial | PRD-04, PRD-07, PRD-08 |
| [PRD-21](./PRD-021-quality-gates.md) | CI, contracts, accessibility, quality | Cross-cutting acceptance tests | 🟡 Partial | All implemented workflows |
| [PRD-22](./PRD-022-operations-rollback.md) | Operations and rollback | Deployment, backups, monitoring | 🟡 Partial | PRD-21 and workflow PRDs |
| [PRD-23](./PRD-023-release-readiness.md) | MVP release readiness | MVP completion criteria | ⬜ Not Started | PRD-01..PRD-22 |

## Release sequence

1. Foundation: PRD-01 → PRD-02 → PRD-03 → PRD-04.
2. Identity and catalogue: PRD-05 → PRD-06 → PRD-07 → PRD-08.
3. Managed stock: PRD-09 → PRD-10 → PRD-11.
4. Commerce: PRD-12 → PRD-13 → PRD-15 → PRD-16.
5. Extensions and controls: PRD-14, PRD-17, PRD-18, PRD-19, PRD-20.
6. Release gates: PRD-21 → PRD-22 → PRD-23.

## Update rules

- Update the matching PRD and this board in the same change.
- Do not mark a workflow complete because a screen or scaffold exists alone.
- Keep user-story IDs, API paths, dependencies, and status evidence synchronized.
- Do not invent unresolved payment, shipping, tax, payout, currency, or upload-policy decisions.
