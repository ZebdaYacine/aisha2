ALTER TABLE artisan_documents
    DROP CONSTRAINT IF EXISTS artisan_documents_visibility_valid,
    DROP COLUMN IF EXISTS visibility;

DROP TABLE IF EXISTS product_submissions;
