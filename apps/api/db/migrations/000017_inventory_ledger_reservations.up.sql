ALTER TABLE inventory_movements
    DROP CONSTRAINT inventory_movements_stock_bucket_valid,
    ADD CONSTRAINT inventory_movements_stock_bucket_valid
        CHECK (stock_bucket IN ('AVAILABLE', 'REJECTED', 'QUARANTINED', 'DAMAGED', 'SHIPPED'));

ALTER TABLE stock_reservations
    ADD COLUMN committed_at timestamptz,
    ADD CONSTRAINT stock_reservations_state_timestamps_valid CHECK (
        (status <> 'RELEASED' OR released_at IS NOT NULL)
        AND (status <> 'COMMITTED' OR committed_at IS NOT NULL)
    );

CREATE INDEX stock_reservations_held_expiry_idx
    ON stock_reservations (expires_at, id)
    WHERE status = 'HELD';
