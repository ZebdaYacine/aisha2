CREATE TABLE categories (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    slug text NOT NULL UNIQUE,
    display_name text NOT NULL,
    parent_id uuid,
    sort_order integer NOT NULL DEFAULT 0,
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT categories_parent_fk
        FOREIGN KEY (parent_id) REFERENCES categories (id) ON DELETE RESTRICT,
    CONSTRAINT categories_slug_not_blank CHECK (btrim(slug) <> ''),
    CONSTRAINT categories_display_name_not_blank CHECK (btrim(display_name) <> ''),
    CONSTRAINT categories_not_own_parent CHECK (parent_id IS NULL OR parent_id <> id)
);
CREATE INDEX categories_active_sort_idx ON categories (is_active, sort_order, slug);
CREATE INDEX categories_parent_id_idx ON categories (parent_id);

CREATE TABLE artisan_profile_categories (
    artisan_profile_id uuid NOT NULL,
    category_id uuid NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (artisan_profile_id, category_id),
    CONSTRAINT artisan_profile_categories_profile_fk
        FOREIGN KEY (artisan_profile_id) REFERENCES artisan_profiles (id) ON DELETE CASCADE,
    CONSTRAINT artisan_profile_categories_category_fk
        FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE RESTRICT
);
CREATE INDEX artisan_profile_categories_category_idx
    ON artisan_profile_categories (category_id);

CREATE TABLE products (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    artisan_profile_id uuid NOT NULL,
    category_id uuid NOT NULL,
    product_type text NOT NULL,
    status text NOT NULL DEFAULT 'DRAFT',
    price_minor bigint NOT NULL,
    currency char(3) NOT NULL,
    materials text,
    production_method text,
    intended_use text,
    dimensions text,
    weight_grams integer,
    country_of_origin text,
    region_of_origin text,
    eco_friendly_verified boolean NOT NULL DEFAULT false,
    fair_trade_verified boolean NOT NULL DEFAULT false,
    made_to_order_eligible boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    published_at timestamptz,
    CONSTRAINT products_artisan_fk
        FOREIGN KEY (artisan_profile_id) REFERENCES artisan_profiles (id) ON DELETE RESTRICT,
    CONSTRAINT products_category_fk
        FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE RESTRICT,
    CONSTRAINT products_type_valid
        CHECK (product_type IN ('ARTISAN_SPECIFIC', 'STANDARD_TRADITIONAL')),
    CONSTRAINT products_status_valid
        CHECK (status IN ('DRAFT', 'PENDING_REVIEW', 'CHANGES_REQUESTED', 'APPROVED', 'ACTIVE', 'SUSPENDED', 'ARCHIVED')),
    CONSTRAINT products_price_positive CHECK (price_minor > 0),
    CONSTRAINT products_currency_format CHECK (currency ~ '^[A-Z]{3}$'),
    CONSTRAINT products_weight_positive CHECK (weight_grams IS NULL OR weight_grams > 0)
);
CREATE INDEX products_artisan_id_idx ON products (artisan_profile_id);
CREATE INDEX products_category_status_idx ON products (category_id, status);
CREATE INDEX products_status_created_idx ON products (status, created_at DESC);

CREATE TABLE product_translations (
    product_id uuid NOT NULL,
    locale varchar(5) NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    story text,
    cultural_context text,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (product_id, locale),
    CONSTRAINT product_translations_product_fk
        FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE,
    CONSTRAINT product_translations_locale_valid CHECK (locale IN ('ar', 'fr', 'en', 'es')),
    CONSTRAINT product_translations_name_not_blank CHECK (btrim(name) <> ''),
    CONSTRAINT product_translations_description_not_blank CHECK (btrim(description) <> '')
);

CREATE TABLE product_media (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id uuid NOT NULL,
    media_kind text NOT NULL,
    object_key text NOT NULL UNIQUE,
    original_filename text,
    media_type text NOT NULL,
    size_bytes bigint NOT NULL,
    width_pixels integer,
    height_pixels integer,
    duration_seconds integer,
    checksum_sha256 char(64),
    alt_text text,
    sort_order integer NOT NULL DEFAULT 0,
    visibility text NOT NULL DEFAULT 'PRIVATE',
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT product_media_product_fk
        FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE RESTRICT,
    CONSTRAINT product_media_kind_valid CHECK (media_kind IN ('IMAGE', 'VIDEO')),
    CONSTRAINT product_media_object_key_not_blank CHECK (btrim(object_key) <> ''),
    CONSTRAINT product_media_media_type_not_blank CHECK (btrim(media_type) <> ''),
    CONSTRAINT product_media_size_positive CHECK (size_bytes > 0),
    CONSTRAINT product_media_width_positive CHECK (width_pixels IS NULL OR width_pixels > 0),
    CONSTRAINT product_media_height_positive CHECK (height_pixels IS NULL OR height_pixels > 0),
    CONSTRAINT product_media_duration_positive CHECK (duration_seconds IS NULL OR duration_seconds > 0),
    CONSTRAINT product_media_visibility_valid CHECK (visibility IN ('PRIVATE', 'PUBLIC'))
);
CREATE INDEX product_media_product_sort_idx ON product_media (product_id, sort_order, created_at);
