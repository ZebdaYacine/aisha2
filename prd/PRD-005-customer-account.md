# PRD-05 · Customer Profile, Addresses & Combined Account

| | |
| --- | --- |
| **Epic** | 5 — Customer Profile, Addresses & Combined Account |
| **Stories** | US-ACC-001; US-AUTH-001 … US-AUTH-003 |
| **Priority** | High |
| **Status** | ✅ Implemented |
| **Surfaces** | `apps/web/app/[locale]/account`, `apps/web/app/[locale]/account/addresses` |
| **API** | `GET /api/v1/me`, profile and address CRUD endpoints |

## 1. Summary

Give every authenticated user one trusted account workspace for profile, addresses, purchases, and tracking, without removing customer capabilities when that user becomes an artisan.

## 2. Problem

Customer data, seller capabilities, orders, and delivery addresses become unsafe and confusing when users need separate accounts or when ownership is enforced only in the UI.

## 3. Goals

- Keep customer and artisan capabilities on one account.
- Let users manage their own profile and delivery addresses.
- Provide clear account navigation for purchases, tracking, and capability-gated artisan work.
- Preserve ownership and default-address integrity.

### Non-goals

- A second identity for artisan membership.
- Administrative access to customer data without explicit permission.

## 4. Users

**Customer** — manages profile, addresses, purchases, and tracking.
**Artisan** — uses all customer features while also accessing permitted seller areas.
**Administrator** — performs explicitly authorized support or account actions.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Account pages show backend-derived user status, capabilities, and artisan status. | US-ACC-001 |
| R2 | Users can view and update only their own profile. | US-ACC-001 |
| R3 | Users can create, update, delete, and select one default delivery address. | US-CHECK-001 |
| R4 | Account navigation exposes profile, addresses, orders, tracking, and permitted artisan work. | US-ACC-001 |
| R5 | An artisan retains browsing, cart, checkout, payment, cancellation, and tracking capabilities. | US-ART-004, US-CHECK-003 |
| R6 | Loading, empty, pending, suspended, forbidden, and error states are explicit. | US-ACC-001 |

## 6. Flow

```text
Sign in → account overview → profile/addresses/orders/tracking
                 └─ if artisan membership is active → artisan workspace
```

## 7. Technical notes

- Account status and capabilities come from the backend, not duplicated frontend rules.
- Default-address changes require an ownership check and transaction.
- Orders and tracking remain separate customer-owned resources. Profile management keeps optional phone values valid, handles expired sessions by clearing cookies and returning to login, and logout always clears local state and redirects.

## 8. Success metrics

- Address CRUD completion and default-address correctness.
- Percentage of artisan users retaining successful customer journeys.
- Unauthorized account and address access attempts rejected.

## 9. Risks & open questions

- Canonical membership activation and real buyer order/tracking APIs remain downstream dependencies tracked by PRD-07 and PRD-13.
- Account capability labels must stay synchronized with Casbin and membership states as those epics evolve.

## Source traceability

`US-ACC-001`, `US-AUTH-001..003`, `US-CHECK-001`, `US-CHECK-003`, `docs/jira_project.md`.
