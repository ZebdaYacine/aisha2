# AISHA Codex Instructions

Before changing code, read in order:

1. `docs/requirements.md`
2. `docs/userstory.md`
3. `docs/architecture.md`
4. The relevant technical guide.
5. Existing source and tests.

Rules:

- Do not invent payment providers, shipping tariffs, taxes, commissions, countries, refund windows, or legal requirements.
- Keep Fiber handlers thin.
- Keep business rules in application/domain code.
- Enforce Casbin authorisation in the backend.
- Never expose SQLBoiler models directly.
- Use integer minor units for money.
- Use explicit transactions for inventory, checkout, payment, and order state changes.
- Preserve PostgreSQL and MinIO volumes.
- Add tests for every critical workflow or bug fix.
- Do not claim success for tests or builds that were not run.
- Stop after the requested MVP scope is complete.

Final report:

```md
## Summary
## Requirements
## Files Changed
## Validation
## Security and Authorization
## Assumptions
## Unresolved Issues
```
