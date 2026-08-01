ALTER TABLE users ADD COLUMN phone text;
ALTER TABLE users ADD CONSTRAINT users_phone_not_blank CHECK (phone IS NULL OR btrim(phone) <> '');

CREATE TABLE addresses (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    full_name text NOT NULL,
    phone text,
    line1 text NOT NULL,
    line2 text,
    city text NOT NULL,
    postal_code text NOT NULL,
    country text NOT NULL,
    is_default boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT addresses_user_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT addresses_full_name_not_blank CHECK (btrim(full_name) <> ''),
    CONSTRAINT addresses_line1_not_blank CHECK (btrim(line1) <> ''),
    CONSTRAINT addresses_city_not_blank CHECK (btrim(city) <> ''),
    CONSTRAINT addresses_postal_code_not_blank CHECK (btrim(postal_code) <> ''),
    CONSTRAINT addresses_country_not_blank CHECK (btrim(country) <> '')
);
CREATE INDEX addresses_user_created_idx ON addresses (user_id, created_at DESC);
CREATE UNIQUE INDEX addresses_one_default_per_user_idx ON addresses (user_id) WHERE is_default=true;
