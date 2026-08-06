# Database schema

The source of truth is `apps/api/db/migrations`. SQLBoiler models are generated into
`apps/api/db/sqlboiler/models` after migrations are applied. Generated
models are infrastructure-only and must never be returned directly by HTTP handlers.
