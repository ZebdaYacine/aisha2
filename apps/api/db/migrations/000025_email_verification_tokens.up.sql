CREATE TABLE email_verification_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    token_hash text NOT NULL UNIQUE,
    expires_at timestamptz NOT NULL,
    consumed_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT email_verification_tokens_user_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT email_verification_tokens_hash_not_blank CHECK (btrim(token_hash) <> ''),
    CONSTRAINT email_verification_tokens_expiry_after_creation CHECK (expires_at > created_at)
);

CREATE INDEX email_verification_tokens_user_idx ON email_verification_tokens (user_id, created_at DESC);
CREATE INDEX email_verification_tokens_active_idx ON email_verification_tokens (expires_at) WHERE consumed_at IS NULL;
