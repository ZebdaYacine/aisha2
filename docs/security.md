# AISHA MVP Security Requirements

## 1. Security Objectives

Protect:

- Customer identities and addresses.
- Artisan documents.
- Product draft media.
- Inventory integrity.
- Order and payment state.
- Provider credentials.
- Administrative actions.
- Audit records.

The highest-risk flows are authentication, authorisation, uploads, inventory, checkout, payment webhooks, refunds, and administration.

## 2. Authentication

### Passwords

- Use Argon2id or bcrypt with an approved cost.
- Store only the hash.
- Enforce minimum length and block obviously weak passwords.
- Never log password or password reset token.
- Password reset tokens are random, short-lived, single-use, and stored hashed where practical.

### Sessions

Use one consistent design.

Recommended web design:

- Secure, HttpOnly session or refresh cookie.
- `SameSite=Lax` or stricter unless business flow requires otherwise.
- `Secure` in production.
- Short-lived access credentials.
- Refresh rotation.
- Revocation on logout, password reset, and administrator action.

Detect refresh-token reuse and revoke the related session family.

## 3. Authorisation

Casbin policies must cover:

- Role.
- Resource.
- Action.
- Ownership or tenant-like relationship where applicable.

Examples:

```text
customer, order, read, owner
artisan, product, update, owner
moderator, product_submission, review, any
warehouse, reception, create, any
admin, user_role, update, permitted_scope
```

Do not rely only on route prefix. Check the actual resource.

## 4. Input Validation

Validate:

- JSON size.
- Required fields.
- String lengths.
- Numeric bounds.
- UUIDs.
- Enum values.
- Currency.
- Pagination.
- File metadata.
- URL fields.
- Rich text.

Use allowlists for status and category-like values.

Escape output and avoid rendering untrusted HTML. If rich text is introduced, sanitise it on input or output with a proven library.

## 5. SQL and Data Access

- Use SQLBoiler parameterisation.
- Do not build SQL from untrusted text.
- Use least-privilege database accounts.
- Separate migration credentials from runtime credentials where practical.
- Do not return SQL errors to clients.
- Add database constraints as a second line of defence.

## 6. Payment Security

- Do not store raw card data.
- Do not log provider secrets.
- Do not trust frontend amount, price, currency, or paid status.
- Do not trust browser redirect.
- Verify webhook signature.
- Validate timestamp to reduce replay.
- Store and uniquely constrain provider event ID.
- Compare provider amount and currency to the order.
- Make event handling idempotent.
- Reject invalid state transitions.
- Redact provider payloads in logs.

## 7. Inventory Security and Integrity

Threats:

- Overselling.
- Duplicate reservation.
- Unauthorised adjustment.
- Manipulated quantity.
- Race conditions.
- Replayed warehouse request.

Controls:

- Database transaction.
- Row lock or safe conditional update.
- Idempotency key.
- Non-negative constraints.
- Append-only movement.
- Authorised warehouse roles.
- Mandatory actor and reason.
- Audit event.

## 8. File Upload Security

Required checks:

1. Maximum request size.
2. Actual file size.
3. File signature or magic bytes.
4. Allowed MIME type.
5. Safe generated object key.
6. Private bucket by default.
7. Malware scan when available.
8. Image re-encoding or metadata stripping where practical.
9. Authorised signed read URL.
10. No execution from upload storage.

Reject:

- Executables.
- Scriptable documents unless explicitly required.
- Double-extension tricks.
- Path traversal filenames.
- Oversized images designed to exhaust memory.
- Unsupported archive formats.

## 9. CSRF

When cookie-based authentication is used:

- Use SameSite cookies.
- Protect unsafe methods with CSRF token or verified origin strategy.
- Validate `Origin` or `Referer` for browser state-changing requests where suitable.
- Webhooks use signature authentication, not browser CSRF protection.

## 10. CORS

Production CORS must allow only known frontend origins.

Do not use:

```text
Access-Control-Allow-Origin: *
```

with credentials.

Allow only required methods and headers.

## 11. Security Headers

Configure at Nginx and/or application level:

- Strict-Transport-Security.
- Content-Security-Policy.
- X-Content-Type-Options: nosniff.
- Referrer-Policy.
- Permissions-Policy.
- Frame protection through CSP `frame-ancestors`.
- Secure cache rules for private pages.

CSP must be compatible with Next.js without becoming `unsafe-*` everywhere.

## 12. Rate Limiting

Apply Redis-backed rate limits to:

- Login.
- Registration.
- Password reset.
- Verification.
- File upload.
- Checkout.
- Payment retry.
- Custom-order messages.
- Search abuse if necessary.
- Webhook endpoints by provider/IP only when it does not block legitimate callbacks.

Return `429` with a stable code and optional retry information.

## 13. Secrets

Secrets include:

- Database credentials.
- Redis password.
- MinIO keys.
- Session and JWT secrets.
- Payment keys.
- Webhook secrets.
- Email credentials.
- Jenkins tokens.

Rules:

- `.env` with real secrets is not committed.
- Commit `.env.example` with placeholders.
- Production secrets live in Jenkins credentials, Docker secrets, or approved secret manager.
- Never print environment variables in CI logs.
- Rotate compromised credentials.

## 14. Logging and Privacy

Structured logs may include:

- Timestamp.
- Level.
- Request ID.
- Actor ID.
- Route.
- Method.
- Status.
- Duration.
- Safe error code.

Do not log:

- Password.
- Session token.
- Refresh token.
- Password reset token.
- Payment secret.
- Full webhook secret-bearing payload.
- Private artisan documents.
- Full customer address unless explicitly needed and protected.
- Database connection string with password.

## 15. Audit Logs

Audit events are append-only.

Required data:

- Event type.
- Actor.
- Target.
- UTC timestamp.
- Correlation ID.
- Previous/new state where safe.
- Reason.
- Source metadata when policy permits.

Access to audits is restricted and audited.

## 16. Error Handling

Public API responses never contain:

- Stack traces.
- SQL text.
- Internal paths.
- Provider secret.
- Go panic details.

Unexpected errors return:

```json
{
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "An unexpected error occurred.",
    "details": {},
    "requestId": "..."
  }
}
```

## 17. Container Security

- Use minimal images.
- Run as non-root where practical.
- Pin major image versions.
- Do not bake secrets into images.
- Use read-only filesystem where compatible.
- Drop unnecessary Linux capabilities.
- Add health checks.
- Scan dependencies and images in CI.
- Keep PostgreSQL, Redis, and MinIO off the public internet.

## 18. Production Network

Recommended exposure:

```text
Internet -> 443 reverse proxy
Reverse proxy -> frontend and API internal network
API -> PostgreSQL/Redis/MinIO private network
```

Only reverse proxy ports are publicly exposed.

## 19. Security Tests

Required:

- Customer cannot read another customer's order.
- Artisan cannot edit another artisan's product.
- Unapproved artisan cannot submit product.
- Warehouse route denies customer.
- Admin role management denies lower roles.
- Invalid webhook signature rejected.
- Duplicate webhook handled once.
- Amount mismatch rejected.
- File with spoofed extension rejected.
- Oversized file rejected.
- SQL injection strings do not alter query.
- Rate limit activates.
- Refresh-token reuse protection.
- Inventory race test.
- Sensitive logs are redacted.

## 20. Incident-Ready Behaviour

- Correlation IDs permit tracing.
- Security events have stable event types.
- Sessions can be revoked.
- Users and products can be suspended.
- Provider callbacks can be replayed safely in an authorised reconciliation tool.
- Failed outbox events are visible.
- Backups can be restored.
