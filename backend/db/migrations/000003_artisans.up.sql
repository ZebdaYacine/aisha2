CREATE TABLE artisan_profiles (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL UNIQUE,
    public_display_name text NOT NULL,
    internal_name text,
    workshop_name text,
    wilaya text,
    location_text text,
    contact_email text,
    contact_phone text,
    contact_visibility text NOT NULL DEFAULT 'PRIVATE',
    status text NOT NULL DEFAULT 'DRAFT',
    profile_image_object_key text,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approved_at timestamptz,
    suspended_at timestamptz,
    CONSTRAINT artisan_profiles_user_fk
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,
    CONSTRAINT artisan_profiles_display_name_not_blank CHECK (btrim(public_display_name) <> ''),
    CONSTRAINT artisan_profiles_contact_visibility_valid
        CHECK (contact_visibility IN ('PRIVATE', 'PUBLIC')),
    CONSTRAINT artisan_profiles_status_valid
        CHECK (status IN ('DRAFT', 'SUBMITTED', 'UNDER_REVIEW', 'CHANGES_REQUESTED', 'APPROVED', 'REJECTED', 'SUSPENDED'))
);
CREATE INDEX artisan_profiles_status_idx ON artisan_profiles (status);
CREATE INDEX artisan_profiles_public_name_idx ON artisan_profiles (public_display_name);

CREATE TABLE artisan_profile_translations (
    artisan_profile_id uuid NOT NULL,
    locale varchar(5) NOT NULL,
    biography text,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (artisan_profile_id, locale),
    CONSTRAINT artisan_profile_translations_profile_fk
        FOREIGN KEY (artisan_profile_id) REFERENCES artisan_profiles (id) ON DELETE CASCADE,
    CONSTRAINT artisan_profile_translations_locale_valid
        CHECK (locale IN ('ar', 'fr', 'en', 'es'))
);

CREATE TABLE artisan_documents (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    artisan_profile_id uuid NOT NULL,
    document_type text NOT NULL,
    object_key text NOT NULL UNIQUE,
    original_filename text,
    media_type text NOT NULL,
    size_bytes bigint NOT NULL,
    checksum_sha256 char(64),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT artisan_documents_profile_fk
        FOREIGN KEY (artisan_profile_id) REFERENCES artisan_profiles (id) ON DELETE RESTRICT,
    CONSTRAINT artisan_documents_type_not_blank CHECK (btrim(document_type) <> ''),
    CONSTRAINT artisan_documents_object_key_not_blank CHECK (btrim(object_key) <> ''),
    CONSTRAINT artisan_documents_media_type_not_blank CHECK (btrim(media_type) <> ''),
    CONSTRAINT artisan_documents_size_positive CHECK (size_bytes > 0)
);
CREATE INDEX artisan_documents_profile_id_idx ON artisan_documents (artisan_profile_id);

CREATE TABLE artisan_media (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    artisan_profile_id uuid NOT NULL,
    media_kind text NOT NULL,
    object_key text NOT NULL UNIQUE,
    original_filename text,
    media_type text NOT NULL,
    size_bytes bigint NOT NULL,
    checksum_sha256 char(64),
    sort_order integer NOT NULL DEFAULT 0,
    visibility text NOT NULL DEFAULT 'PRIVATE',
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT artisan_media_profile_fk
        FOREIGN KEY (artisan_profile_id) REFERENCES artisan_profiles (id) ON DELETE RESTRICT,
    CONSTRAINT artisan_media_kind_valid CHECK (media_kind IN ('IMAGE', 'VIDEO')),
    CONSTRAINT artisan_media_object_key_not_blank CHECK (btrim(object_key) <> ''),
    CONSTRAINT artisan_media_media_type_not_blank CHECK (btrim(media_type) <> ''),
    CONSTRAINT artisan_media_size_positive CHECK (size_bytes > 0),
    CONSTRAINT artisan_media_visibility_valid CHECK (visibility IN ('PRIVATE', 'PUBLIC'))
);
CREATE INDEX artisan_media_profile_sort_idx
    ON artisan_media (artisan_profile_id, sort_order, created_at);
