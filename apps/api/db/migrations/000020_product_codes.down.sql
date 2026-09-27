DROP INDEX IF EXISTS products_product_code_uq;

ALTER TABLE products
  DROP CONSTRAINT IF EXISTS products_product_code_format;

ALTER TABLE products DROP COLUMN IF EXISTS product_code;
