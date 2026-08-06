package adapters

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxResetNotifier struct{ pool *pgxpool.Pool }

func NewOutboxResetNotifier(pool *pgxpool.Pool) *OutboxResetNotifier {
	return &OutboxResetNotifier{pool: pool}
}

func (n *OutboxResetNotifier) SendPasswordReset(ctx context.Context, userID, email, token string) error {
	payload, err := json.Marshal(map[string]string{"email": email, "token": token})
	if err != nil {
		return fmt.Errorf("marshal password reset event: %w", err)
	}
	_, err = n.pool.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES('PASSWORD_RESET_REQUESTED','user',$1,$2)`, userID, payload)
	if err != nil {
		return fmt.Errorf("insert password reset outbox event: %w", err)
	}
	return nil
}
