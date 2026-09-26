DROP INDEX IF EXISTS stock_reservations_held_expiry_idx;

ALTER TABLE stock_reservations
    DROP CONSTRAINT IF EXISTS stock_reservations_state_timestamps_valid,
    DROP COLUMN IF EXISTS committed_at;

ALTER TABLE inventory_movements
    DROP CONSTRAINT IF EXISTS inventory_movements_stock_bucket_valid,
    ADD CONSTRAINT inventory_movements_stock_bucket_valid
        CHECK (stock_bucket IN ('AVAILABLE', 'REJECTED', 'QUARANTINED', 'DAMAGED'));
