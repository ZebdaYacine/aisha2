# PRD-04 · Authentication & Authorization

| | |
| --- | --- |
| **Epic** | 4 — Authentication & Authorization |
| **Stories** | US-AUTH-001 … US-AUTH-003; REQ-RBAC-001 |
| **Priority** | Critical |
| **Status** | ✅ Implemented |
| **Surfaces** | `apps/web/app/[locale]/login`, `register`, `account`; `apps/api/internal/features/auth` |
| **API** | `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `POST /api/v1/auth/refresh`, `POST /api/v1/auth/logout`, `GET /api/v1/me` |

## 1. Summary

Allow visitors to create and use secure accounts while ensuring protected actions are authenticated, authorized, and ownership-checked.

## 2. Problem

Weak identity handling can expose customer data, permit unauthorized seller actions, and make sessions difficult to revoke or audit.

## 3. Goals

- Provide secure registration, sign-in, refresh, logout, and password reset.
- Enforce backend authentication, Casbin authorization, and resource ownership.
- Return stable public errors without exposing credentials or internal details.

### Non-goals

- Social login or an unapproved identity provider in the MVP.

## 4. Users

- Visitor registering for customer capabilities.
- Customer signing in and managing private data.
- Artisan using the same account for seller capabilities.
- Administrator managing permitted roles and sessions.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Validate registration, hash passwords, enforce uniqueness, and rate-limit repeated attempts. | US-AUTH-001 |
| R2 | Authenticate active users with generic invalid-credential responses. | US-AUTH-002 |
| R3 | Rotate short-lived access/refresh credentials and revoke sessions on logout. | US-AUTH-002..003 |
| R4 | Enforce Casbin permission and ownership on every protected route. | REQ-RBAC-001 |
| R5 | Prevent suspended users from signing in while preserving permitted historical access. | US-AUTH-002 |

## 6. Flow

```text
Register/sign in → create session → load backend capabilities → access protected route → revoke on logout
```

## 7. Technical notes

- Passwords, tokens, and reset secrets are never logged.
- Secure HTTP-only cookies, refresh rotation, rate limits, request IDs, and stable error codes are required.
- Frontend guards improve UX but never authorize a request.

## 8. Success metrics

- Successful registration and sign-in rate without security-test regressions.
- Protected routes consistently return `401` or `403`.
- Session revocation and refresh-reuse tests pass.

## 9. Risks & open questions

- Cookie/CSRF policy must remain aligned between frontend and API deployments.
- Password verification and email-delivery operational limits need monitoring.

## Source traceability

`REQ-AUTH-001`, `REQ-AUTH-002`, `REQ-RBAC-001`, `US-AUTH-001..003`, `docs/security.md`.
