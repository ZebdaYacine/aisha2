CREATE TABLE product_submissions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id uuid NOT NULL,
    version integer NOT NULL,
    submitted_by_user_id uuid NOT NULL,
    snapshot jsonb NOT NULL,
    submitted_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT product_submissions_product_fk
        FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE,
    CONSTRAINT product_submissions_user_fk
        FOREIGN KEY (submitted_by_user_id) REFERENCES users (id) ON DELETE RESTRICT,
    CONSTRAINT product_submissions_version_positive CHECK (version > 0),
    CONSTRAINT product_submissions_product_version_unique UNIQUE (product_id, version)
);

CREATE INDEX product_submissions_product_time_idx
    ON product_submissions (product_id, submitted_at DESC);

ALTER TABLE artisan_documents
    ADD COLUMN visibility text NOT NULL DEFAULT 'PRIVATE',
    ADD CONSTRAINT artisan_documents_visibility_valid
        CHECK (visibility IN ('PRIVATE'));
