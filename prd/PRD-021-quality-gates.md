# PRD-21 · CI, Contracts, Accessibility & Quality Gates

| | |
| --- | --- |
| **Epic** | 21 — CI, Contracts, Accessibility & Quality Gates |
| **Stories** | Required end-to-end, security, accessibility, concurrency, and deployment tests |
| **Priority** | High |
| **Status** | 🟡 Partial |
| **Surfaces** | Jenkins, API/frontend test suites, OpenAPI, migration checks |
| **API** | OpenAPI contract and health/smoke checks |

## 1. Summary

Make workflow correctness, API contracts, security, accessibility, responsive UI, and performance release-blocking quality checks.

## 2. Problem

A green unit-test build can still ship broken ownership rules, stale contracts, failed migrations, inaccessible RTL pages, or concurrency bugs.

## 3. Goals

- Test each completed feature at domain, application, repository, API, integration, concurrency, and E2E levels as appropriate.
- Validate OpenAPI, migrations, generated mocks, builds, and accessibility.
- Make evidence visible before changing Jira status.

### Non-goals

- Claiming production certification from local tests alone.

## 4. Users

- Developer receiving fast feedback.
- Reviewer checking traceable acceptance evidence.
- Release manager deciding readiness.
- Customer relying on accessible and stable workflows.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Run Go format, vet, unit, integration, concurrency, and build checks. | Quality requirements |
| R2 | Run frontend lint, type-check, component, integration, E2E, visual, accessibility, and build checks. | Frontend testing rules |
| R3 | Validate OpenAPI coverage for every endpoint and migration up/down behavior. | API and database rules |
| R4 | Test authentication, ownership, file security, webhooks, inventory races, and idempotency. | Security and concurrency stories |
| R5 | Keep evidence linked before marking a Jira item implemented or done. | Jira rules |

## 6. Flow

```text
Change → static checks → unit → integration/test stack → API contract → browser/accessibility → build → evidence → tracker update
```

## 7. Technical notes

- Application tests use generated mocks; repository tests use real test services.
- CI must not use production credentials.
- Visual and E2E checks cover Arabic RTL and mobile critical paths.

## 8. Success metrics

- Required pipeline checks pass on every merge.
- Critical workflow regression rate.
- Percentage of completed PRDs with linked acceptance evidence.

## 9. Risks & open questions

- Test-container runtime and CI duration need tuning.
- Coverage thresholds and browser matrix need explicit release policy.

## Source traceability

`docs/architecture.md` §§27, 56–60, `docs/deployment.md`, `docs/security.md`.
