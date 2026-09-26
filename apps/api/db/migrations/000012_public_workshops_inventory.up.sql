CREATE TABLE workshops (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    artisan_profile_id uuid NOT NULL,
    name text NOT NULL,
    description text,
    wilaya text,
    location_text text,
    status text NOT NULL DEFAULT 'ACTIVE',
    is_default boolean NOT NULL DEFAULT false,
    is_public boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT workshops_artisan_fk FOREIGN KEY (artisan_profile_id) REFERENCES artisan_profiles (id) ON DELETE RESTRICT,
    CONSTRAINT workshops_status_valid CHECK (status IN ('ACTIVE', 'INACTIVE', 'ARCHIVED')),
    CONSTRAINT workshops_name_not_blank CHECK (btrim(name) <> '')
);

CREATE UNIQUE INDEX workshops_one_default_per_artisan_idx
    ON workshops (artisan_profile_id) WHERE is_default = true;
CREATE INDEX workshops_public_status_idx
    ON workshops (status, is_public, artisan_profile_id);

CREATE TABLE workshop_translations (
    workshop_id uuid NOT NULL,
    locale varchar(5) NOT NULL,
    name text,
    description text,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (workshop_id, locale),
    CONSTRAINT workshop_translations_workshop_fk FOREIGN KEY (workshop_id) REFERENCES workshops (id) ON DELETE CASCADE,
    CONSTRAINT workshop_translations_locale_valid CHECK (locale IN ('ar', 'fr', 'en', 'es'))
);
CREATE INDEX workshop_translations_locale_idx ON workshop_translations (locale, workshop_id);

INSERT INTO workshops (artisan_profile_id, name, description, wilaya, location_text, status, is_default, is_public)
SELECT id, COALESCE(NULLIF(btrim(workshop_name), ''), public_display_name), NULL, wilaya, location_text,
       CASE WHEN status = 'APPROVED' THEN 'ACTIVE' ELSE 'INACTIVE' END, true, status = 'APPROVED'
FROM artisan_profiles;

ALTER TABLE products ADD COLUMN workshop_id uuid;
ALTER TABLE products
    ADD CONSTRAINT products_workshop_fk FOREIGN KEY (workshop_id) REFERENCES workshops (id) ON DELETE RESTRICT;
CREATE INDEX products_workshop_status_idx ON products (workshop_id, status);

UPDATE products p
SET workshop_id = w.id
FROM workshops w
WHERE w.artisan_profile_id = p.artisan_profile_id AND w.is_default = true;

CREATE TABLE inventory_movements (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id uuid NOT NULL,
    movement_type text NOT NULL,
    quantity_delta bigint NOT NULL,
    reference_key text NOT NULL UNIQUE,
    reason text NOT NULL,
    actor_user_id uuid,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT inventory_movements_product_fk FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE RESTRICT,
    CONSTRAINT inventory_movements_actor_fk FOREIGN KEY (actor_user_id) REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT inventory_movements_type_valid CHECK (movement_type IN ('ACCEPTED', 'ADJUSTMENT', 'RESERVED', 'RELEASED', 'COMMITTED', 'SHIPPED')),
    CONSTRAINT inventory_movements_delta_nonzero CHECK (quantity_delta <> 0),
    CONSTRAINT inventory_movements_reason_not_blank CHECK (btrim(reason) <> '')
);
CREATE INDEX inventory_movements_product_created_idx
    ON inventory_movements (product_id, created_at, id);
