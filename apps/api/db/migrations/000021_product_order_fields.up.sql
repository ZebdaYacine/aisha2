ALTER TABLE products
    ADD COLUMN planned_quantity integer,
    ADD COLUMN order_total_minor bigint;

UPDATE products
SET planned_quantity = 1,
    order_total_minor = price_minor
WHERE planned_quantity IS NULL OR order_total_minor IS NULL;

ALTER TABLE products
    ALTER COLUMN planned_quantity SET DEFAULT 1,
    ALTER COLUMN planned_quantity SET NOT NULL,
    ALTER COLUMN order_total_minor SET DEFAULT 1,
    ALTER COLUMN order_total_minor SET NOT NULL;

ALTER TABLE products
    ADD CONSTRAINT products_planned_quantity_positive CHECK (planned_quantity > 0),
    ADD CONSTRAINT products_order_total_positive CHECK (order_total_minor > 0);
