CREATE TABLE roles (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code text NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT roles_code_not_blank CHECK (btrim(code) <> '')
);

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email text,
    password_hash text,
    status text NOT NULL DEFAULT 'ACTIVE',
    display_name text,
    email_verified_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_email_not_blank CHECK (email IS NULL OR btrim(email) <> ''),
    CONSTRAINT users_status_valid CHECK (status IN ('ACTIVE', 'SUSPENDED', 'DISABLED')),
    CONSTRAINT users_password_hash_not_blank CHECK (password_hash IS NULL OR btrim(password_hash) <> '')
);

CREATE UNIQUE INDEX users_email_unique_ci
    ON users (lower(email))
    WHERE email IS NOT NULL;
CREATE INDEX users_status_idx ON users (status);

CREATE TABLE user_roles (
    user_id uuid NOT NULL,
    role_id uuid NOT NULL,
    assigned_by_user_id uuid,
    assigned_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id),
    CONSTRAINT user_roles_user_fk
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT user_roles_role_fk
        FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE RESTRICT,
    CONSTRAINT user_roles_assigned_by_fk
        FOREIGN KEY (assigned_by_user_id) REFERENCES users (id) ON DELETE SET NULL
);
CREATE INDEX user_roles_role_id_idx ON user_roles (role_id);

CREATE TABLE sessions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    token_hash text NOT NULL UNIQUE,
    family_id uuid NOT NULL,
    expires_at timestamptz NOT NULL,
    last_used_at timestamptz,
    revoked_at timestamptz,
    replaced_by_session_id uuid,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT sessions_user_fk
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT sessions_replacement_fk
        FOREIGN KEY (replaced_by_session_id) REFERENCES sessions (id) ON DELETE SET NULL,
    CONSTRAINT sessions_token_hash_not_blank CHECK (btrim(token_hash) <> ''),
    CONSTRAINT sessions_expiry_after_creation CHECK (expires_at > created_at)
);
CREATE INDEX sessions_user_id_idx ON sessions (user_id);
CREATE INDEX sessions_family_id_idx ON sessions (family_id);
CREATE INDEX sessions_active_expiry_idx ON sessions (expires_at) WHERE revoked_at IS NULL;
