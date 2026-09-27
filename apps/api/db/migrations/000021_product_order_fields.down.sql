ALTER TABLE products
    DROP CONSTRAINT IF EXISTS products_order_total_positive,
    DROP CONSTRAINT IF EXISTS products_planned_quantity_positive,
    DROP COLUMN IF EXISTS order_total_minor,
    DROP COLUMN IF EXISTS planned_quantity;
