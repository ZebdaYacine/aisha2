package auth

import (
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth/data/adapters"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxResetNotifier = adapters.OutboxResetNotifier

func NewOutboxResetNotifier(pool *pgxpool.Pool) *OutboxResetNotifier {
	return adapters.NewOutboxResetNotifier(pool)
}
