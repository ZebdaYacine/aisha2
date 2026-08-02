ALTER TABLE artisan_profiles
    ADD COLUMN review_reason text,
    ADD COLUMN submitted_at timestamptz;

ALTER TABLE artisan_profiles
    ADD CONSTRAINT artisan_profiles_review_reason_not_blank
    CHECK (review_reason IS NULL OR btrim(review_reason) <> '');
