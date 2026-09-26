package inventory

import (
	"context"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/inventory/data/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ReservationExpiryWorker releases expired holds independently of checkout.
// The database row locks make repeated ticks and checkout entry safe.
type ReservationExpiryWorker struct {
	repository *repositories.PostgresRepository
	interval   time.Duration
}

func NewReservationExpiryWorker(pool *pgxpool.Pool, interval time.Duration) *ReservationExpiryWorker {
	if interval <= 0 {
		interval = time.Minute
	}
	return &ReservationExpiryWorker{repository: repositories.NewPostgresRepository(pool), interval: interval}
}

func (w *ReservationExpiryWorker) Run(ctx context.Context) {
	if err := w.repository.ReleaseExpired(ctx); err != nil && ctx.Err() == nil { /* next tick retries */
	}
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = w.repository.ReleaseExpired(ctx)
		}
	}
}
