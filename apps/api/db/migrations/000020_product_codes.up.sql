ALTER TABLE products ADD COLUMN product_code text;

UPDATE products
SET product_code = 'AISHA-' || upper(substr(replace(id::text, '-', ''), 1, 10))
WHERE product_code IS NULL;

ALTER TABLE products
  ADD CONSTRAINT products_product_code_format
  CHECK (product_code IS NULL OR product_code ~ '^AISHA-[A-Z0-9]{10}$');

CREATE UNIQUE INDEX products_product_code_uq
  ON products(product_code)
  WHERE product_code IS NOT NULL;
