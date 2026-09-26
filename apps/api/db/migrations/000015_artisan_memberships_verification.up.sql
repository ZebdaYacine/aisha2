CREATE TABLE artisan_memberships (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL UNIQUE REFERENCES users(id) ON DELETE RESTRICT,
    artisan_profile_id uuid NOT NULL UNIQUE REFERENCES artisan_profiles(id) ON DELETE RESTRICT,
    status text NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE','SUSPENDED','CLOSED')),
    reason text,
    activated_at timestamptz,
    suspended_at timestamptz,
    closed_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT artisan_memberships_reason_not_blank CHECK (reason IS NULL OR btrim(reason) <> '')
);
CREATE INDEX artisan_memberships_status_idx ON artisan_memberships(status);

INSERT INTO artisan_memberships(user_id, artisan_profile_id, status, activated_at)
SELECT user_id, id, 'ACTIVE', COALESCE(approved_at, CURRENT_TIMESTAMP)
FROM artisan_profiles
WHERE status = 'APPROVED'
ON CONFLICT DO NOTHING;

CREATE TABLE artisan_verifications (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    artisan_membership_id uuid NOT NULL UNIQUE REFERENCES artisan_memberships(id) ON DELETE RESTRICT,
    status text NOT NULL DEFAULT 'NOT_SUBMITTED' CHECK (status IN ('NOT_SUBMITTED','PENDING','VERIFIED','CHANGES_REQUESTED','REJECTED')),
    reason text,
    decided_by_user_id uuid REFERENCES users(id) ON DELETE SET NULL,
    decided_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT artisan_verifications_reason_not_blank CHECK (reason IS NULL OR btrim(reason) <> '')
);
CREATE INDEX artisan_verifications_status_idx ON artisan_verifications(status, updated_at DESC);

INSERT INTO artisan_verifications(artisan_membership_id)
SELECT id FROM artisan_memberships
ON CONFLICT DO NOTHING;
