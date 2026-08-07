# AISHA AI Development Guide

## Project Documentation

Before making any code changes, always read these documents in order:

1. docs/requirements.md
2. docs/userstory.md
3. docs/architecture.md
4. docs/backend.md
5. docs/frontend.md
6. docs/security.md
7. docs/deployment.md
8. docs/ui_design.md
9. docs/jira_project.md

These documents are the single source of truth.

---

# Development Workflow

Always follow this workflow:

1. Read the related requirements.
2. Read the corresponding user stories.
3. Review the architecture.
4. Check backend/frontend rules.
5. Read jira_project.md.
6. Implement only the selected task.
7. Add tests.
8. Update OpenAPI.
9. Update documentation.
10. Update jira_project.md.

Never skip steps.

---

# Jira Rules

docs/jira_project.md contains the implementation tracker.

Implement only tasks marked:

- ⬜ Not Started
- 🟡 Partial

Never modify completed features unless fixing a bug or adding a dependency.

After finishing:

- change task status
- update epic status
- update remaining work

---

# Architecture Rules

Always follow Clean Architecture.

Handlers

- HTTP only
- validation
- response mapping

Services

- business logic only

Repositories

- persistence only

Models

- SQLBoiler models never leave the repository layer.

DTOs

- Never expose database models.

---

# Backend Rules

Use:

- Go
- Fiber
- SQLBoiler
- PostgreSQL
- Redis
- MinIO
- Casbin

Always:

- transactions where required
- dependency injection
- validation
- structured logging
- request ids
- context propagation

---

# Frontend Rules

Use:

- Next.js App Router
- React Server Components
- shadcn/ui
- Tailwind
- TypeScript
- React Hook Form
- Zod

No business logic in UI.

Use server actions or API routes only.

---

# Database Rules

Every schema change must include:

- migration up
- migration down
- indexes
- constraints
- foreign keys

Never edit production migrations.

---

# API Rules

Every endpoint must include:

- request DTO
- response DTO
- validation
- OpenAPI
- authorization
- error responses

---

# Testing Rules

Every completed feature requires:

- unit tests
- integration tests
- repository tests
- API tests

Critical workflows also require concurrency tests.

---

# Security Rules

Always enforce:

- authentication
- Casbin authorization
- ownership validation
- input validation
- rate limiting
- secure cookies
- signed URLs for private files

Never trust frontend data.

---

# File Upload Rules

Use MinIO.

Always:

- validate mime type
- validate size
- generate object keys
- never trust filenames

Private files must use signed URLs.

---

# Coding Style

- Small functions.
- Thin handlers.
- Clear interfaces.
- Dependency injection.
- No duplicated logic.
- Prefer composition.
- Keep packages cohesive.

---

# Before Completing Any Task

Verify:

- Requirements satisfied
- User story implemented
- API working
- UI connected
- Tests passing
- OpenAPI updated
- Documentation updated
- Jira tracker updated

Only then consider the task complete.

---

# Never

Never:

- invent business rules
- skip tests
- skip migrations
- ignore documentation
- expose SQLBoiler models
- hardcode secrets
- modify unrelated modules
- mark incomplete work as complete

If information is missing, stop and ask for clarification instead of guessing.