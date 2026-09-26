DROP TABLE IF EXISTS warehouse_evidence;
DROP TABLE IF EXISTS warehouse_inspections;
DROP TABLE IF EXISTS warehouse_receptions;
DROP INDEX IF EXISTS inventory_movements_product_bucket_created_idx;
ALTER TABLE inventory_movements
    DROP CONSTRAINT IF EXISTS inventory_movements_type_valid,
    ADD CONSTRAINT inventory_movements_type_valid
        CHECK (movement_type IN ('ACCEPTED', 'ADJUSTMENT', 'RESERVED', 'RELEASED', 'COMMITTED', 'SHIPPED')),
    DROP CONSTRAINT IF EXISTS inventory_movements_stock_bucket_valid,
    DROP COLUMN IF EXISTS stock_bucket;
