ALTER TABLE inventory_movements
    ADD COLUMN stock_bucket text NOT NULL DEFAULT 'AVAILABLE',
    ADD CONSTRAINT inventory_movements_stock_bucket_valid
        CHECK (stock_bucket IN ('AVAILABLE', 'REJECTED', 'QUARANTINED', 'DAMAGED'));

ALTER TABLE inventory_movements
    DROP CONSTRAINT inventory_movements_type_valid,
    ADD CONSTRAINT inventory_movements_type_valid
        CHECK (movement_type IN ('RECEPTION', 'INSPECTION_ACCEPT', 'INSPECTION_REJECT', 'INSPECTION_QUARANTINE', 'INSPECTION_DAMAGE', 'ACCEPTED', 'ADJUSTMENT', 'RESERVED', 'RELEASED', 'COMMITTED', 'SHIPPED', 'RETURN'));

CREATE INDEX inventory_movements_product_bucket_created_idx
    ON inventory_movements (product_id, stock_bucket, created_at, id);

CREATE TABLE warehouse_receptions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id uuid NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    artisan_profile_id uuid NOT NULL REFERENCES artisan_profiles(id) ON DELETE RESTRICT,
    workshop_id uuid NOT NULL REFERENCES workshops(id) ON DELETE RESTRICT,
    supplier_name text,
    received_quantity bigint NOT NULL CHECK (received_quantity > 0),
    reference_key text NOT NULL UNIQUE,
    parcel_reference text,
    notes text,
    status text NOT NULL DEFAULT 'RECEIVED_PENDING_INSPECTION'
        CHECK (status IN ('RECEIVED_PENDING_INSPECTION', 'INSPECTED')),
    received_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    received_by_user_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    inspected_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT warehouse_receptions_reference_not_blank CHECK (btrim(reference_key) <> ''),
    CONSTRAINT warehouse_receptions_supplier_not_blank CHECK (supplier_name IS NULL OR btrim(supplier_name) <> ''),
    CONSTRAINT warehouse_receptions_parcel_not_blank CHECK (parcel_reference IS NULL OR btrim(parcel_reference) <> '')
);
CREATE INDEX warehouse_receptions_status_received_idx
    ON warehouse_receptions (status, received_at DESC);
CREATE INDEX warehouse_receptions_product_received_idx
    ON warehouse_receptions (product_id, received_at DESC);

CREATE TABLE warehouse_inspections (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    reception_id uuid NOT NULL UNIQUE REFERENCES warehouse_receptions(id) ON DELETE RESTRICT,
    accepted_quantity bigint NOT NULL CHECK (accepted_quantity >= 0),
    rejected_quantity bigint NOT NULL CHECK (rejected_quantity >= 0),
    quarantined_quantity bigint NOT NULL CHECK (quarantined_quantity >= 0),
    damaged_quantity bigint NOT NULL CHECK (damaged_quantity >= 0),
    reason text NOT NULL,
    inspected_by_user_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    inspected_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT warehouse_inspections_reason_not_blank CHECK (btrim(reason) <> '')
);
CREATE INDEX warehouse_inspections_inspected_at_idx
    ON warehouse_inspections (inspected_at DESC);

CREATE TABLE warehouse_evidence (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    reception_id uuid NOT NULL REFERENCES warehouse_receptions(id) ON DELETE RESTRICT,
    object_key text NOT NULL UNIQUE,
    original_filename text,
    media_type text NOT NULL,
    size_bytes bigint NOT NULL CHECK (size_bytes > 0),
    checksum_sha256 char(64),
    uploaded_by_user_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT warehouse_evidence_object_key_not_blank CHECK (btrim(object_key) <> ''),
    CONSTRAINT warehouse_evidence_media_type_not_blank CHECK (btrim(media_type) <> '')
);
CREATE INDEX warehouse_evidence_reception_created_idx
    ON warehouse_evidence (reception_id, created_at);
