DROP TABLE IF EXISTS category_translations;

ALTER TABLE artisan_profile_translations
    DROP CONSTRAINT IF EXISTS artisan_profile_translations_display_name_not_blank,
    DROP COLUMN IF EXISTS display_name;
