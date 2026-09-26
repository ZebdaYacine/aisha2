DROP TABLE IF EXISTS inventory_movements;

ALTER TABLE products
    DROP CONSTRAINT IF EXISTS products_workshop_fk,
    DROP COLUMN IF EXISTS workshop_id;

DROP TABLE IF EXISTS workshop_translations;
DROP TABLE IF EXISTS workshops;
