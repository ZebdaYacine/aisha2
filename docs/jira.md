# AISHA Jira Delivery Tracker

Last reviewed: 2026-08-07

This scoped update completes `PRJ-EPIC-007` and everything listed under
`PRJ-EPIC-008`. The full project tracker remains in
[`jira_project.md`](./jira_project.md).

## Completed scope

| Jira item | Status | Completed work |
|---|---|---|
| `PRJ-EPIC-007` | ✅ Complete | Artisan applications, multilingual profile data, category ownership, approved profile edits, private document/profile-media uploads, authorized signed reads, admin review decisions, user-role administration, audit filtering, audit/outbox events, frontend screens, OpenAPI, migrations, and tests. |
| `PRJ-ART-001` | ✅ Complete | Application submit/resubmit and ownership-protected profile workflow. |
| `PRJ-ART-002` | ✅ Complete | Admin review, decision reasons, role assignment, and protected routes. |
| `PRJ-ART-003` | ✅ Complete | Generated MinIO keys, file signature/size validation, private storage, signed reads, cleanup, UI, and tests. |
| `PRJ-ADMIN-001` | ✅ Complete for this epic | Paginated users, role assignment, Casbin authorization, audit/outbox event, and UI. Broader catalogue administration is later scope. |
| `PRJ-ADMIN-002` | ✅ Complete for this epic | Filtered/paginated audit API and admin audit UI. |
| `PRJ-EPIC-008` | ✅ Complete | Artisan-owned product draft CRUD, four-locale validation, active-category/type/price rules, private product media, media deletion, submission snapshots, `PENDING_REVIEW`, audit/outbox event, UI, OpenAPI, migrations, and ownership/storage tests. |
| `PRJ-PROD-001` | ✅ Complete | Draft create, edit, list, and get APIs plus the localized editor. |
| `PRJ-PROD-002` | ✅ Complete | Private image/video upload and deletion with generated object keys and signed reads. |
| `PRJ-PROD-003` | ✅ Complete | Mandatory locale/content/media checks and immutable submission version storage. |

## Explicitly not started

`PRJ-EPIC-009` Product moderation and every later epic remain not started. No
moderator queue, approve/reject/request-changes workflow, publication
activation, warehouse, inventory, cart, checkout, order, payment, shipment,
review, or notification worker work was started in this update.

Products submitted by artisans remain `PENDING_REVIEW` and their media remains
private until the separately scoped moderation and publication work is approved.

## Validation

- Backend: `GOCACHE=/tmp/aisha-go-cache go test ./...`
- Frontend: `npm run type-check`
- Frontend lint: `npm run lint`
- Frontend tests: `npm test` — 13 suites, 23 tests
- Frontend production build: `npm run build`
- OpenAPI and migration checks are included in the backend test suite.
