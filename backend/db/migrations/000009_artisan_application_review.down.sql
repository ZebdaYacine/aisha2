ALTER TABLE artisan_profiles DROP CONSTRAINT IF EXISTS artisan_profiles_review_reason_not_blank;
ALTER TABLE artisan_profiles DROP COLUMN IF EXISTS submitted_at, DROP COLUMN IF EXISTS review_reason;
