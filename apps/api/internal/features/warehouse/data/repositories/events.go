package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func warehouseEvent(ctx context.Context, tx pgx.Tx, actor, eventType, targetID, previousStatus, nextStatus string, payload map[string]any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal warehouse event: %w", err)
	}
	if _, err = tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,previous_state,new_state) VALUES($1,$2,'warehouse_reception',$3,jsonb_build_object('status',$4::text),jsonb_build_object('status',$5::text))`, eventType, actor, targetID, previousStatus, nextStatus); err != nil {
		return fmt.Errorf("write warehouse audit event: %w", err)
	}
	if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES($1,'warehouse_reception',$2,$3)`, eventType, targetID, data); err != nil {
		return fmt.Errorf("write warehouse outbox event: %w", err)
	}
	return nil
}
