# Database schema

The source of truth is `backend/db/migrations`. SQLBoiler models are generated into
`backend/db/sqlboiler/models` after migrations are applied. Generated
models are infrastructure-only and must never be returned directly by HTTP handlers.
