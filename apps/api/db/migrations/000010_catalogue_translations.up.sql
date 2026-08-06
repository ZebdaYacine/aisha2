ALTER TABLE artisan_profile_translations
    ADD COLUMN display_name text;

UPDATE artisan_profile_translations t
SET display_name = a.public_display_name
FROM artisan_profiles a
WHERE a.id = t.artisan_profile_id
  AND t.display_name IS NULL;

ALTER TABLE artisan_profile_translations
    ADD CONSTRAINT artisan_profile_translations_display_name_not_blank
    CHECK (display_name IS NULL OR btrim(display_name) <> '');

CREATE TABLE category_translations (
    category_id uuid NOT NULL,
    locale varchar(5) NOT NULL,
    name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (category_id, locale),
    CONSTRAINT category_translations_category_fk
        FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE,
    CONSTRAINT category_translations_locale_valid
        CHECK (locale IN ('ar', 'fr', 'en', 'es')),
    CONSTRAINT category_translations_name_not_blank
        CHECK (btrim(name) <> '')
);

CREATE INDEX category_translations_locale_idx
    ON category_translations (locale, category_id);
