# PRD-23 · MVP Release Readiness

| | |
| --- | --- |
| **Epic** | 23 — MVP Release Readiness |
| **Stories** | MVP completion criteria and cross-feature acceptance |
| **Priority** | Critical |
| **Status** | ⬜ Not Started |
| **Surfaces** | Full product, CI/CD, staging, runbooks, Jira tracker |
| **API** | All MVP API contracts and health/smoke checks |

## 1. Summary

Provide one release gate proving that the documented AISHA MVP workflows, security controls, UI contracts, operations, and open decisions are ready for an approved launch.

## 2. Problem

Individual screens or epics can appear complete while the managed-stock journey, buyer/seller authorization, payment, shipping, backup, accessibility, or rollback path remains unverified.

## 3. Goals

- Confirm every included PRD has implementation and acceptance evidence.
- Run the end-to-end customer, artisan, moderation, warehouse, payment, fulfilment, and administration journeys.
- Resolve or explicitly defer every open business decision.
- Approve staging before production deployment.

### Non-goals

- Marking incomplete work complete to meet a date.
- Launching with unconfirmed payment, shipping, tax, currency, payout, or legal rules.

## 4. Users

- Product owner approving scope and decisions.
- Engineering and QA validating workflows.
- Security and operations reviewers.
- Customer, artisan, warehouse, and administrator relying on the release.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | All MVP requirements and linked user stories have passing acceptance evidence. | MVP completion criteria |
| R2 | Customer purchase and artisan membership preserve ownership, idempotency, and state invariants end to end. | US-AUTH-001..003, US-ART-001..004, US-CHECK-001..003 |
| R3 | Product approval, warehouse inspection, accepted stock, reservation, payment, fulfilment, shipment, and tracking are connected. | US-PROD-005..007, US-WH-001..004, US-PAY-001..005, US-SHIP-001..002 |
| R4 | Security, privacy, accessibility, Arabic RTL, performance, backup, monitoring, and rollback checks pass. | Cross-cutting requirements |
| R5 | External decisions are documented and approved before dependent epics are released. | Open questions in source PRDs |

## 6. Flow

```text
Scope review → feature evidence → integrated test stack → security/accessibility/performance review → staging smoke → approval → production release
```

## 7. Technical notes

- Jira status changes must link tests, review evidence, or deployment evidence.
- Release artifacts include image tags, migration version, OpenAPI spec, runbook, backup, and rollback plan.
- Production smoke tests must avoid real payments and destructive warehouse actions.

## 8. Success metrics

- Percentage of included PRDs with accepted evidence.
- Critical end-to-end and security test pass rate.
- Successful staging release and restore/rollback drill.

## 9. Risks & open questions

- Launch countries, currencies, payment provider, carrier, taxes, refunds, payouts, legal documents, and certification rules remain release blockers until approved.
- Any missing evidence or unresolved critical vulnerability blocks release readiness.

## Source traceability

MVP completion criteria in `docs/userstory.md`, `docs/requirements.md`, `docs/deployment.md`, `docs/security.md`, and `docs/jira_project.md`.
