# Database schema

The source of truth is `backend/migrations`. SQLBoiler models are generated into
`backend/internal/platform/database/models` after migrations are applied. Generated
models are infrastructure-only and must never be returned directly by HTTP handlers.
