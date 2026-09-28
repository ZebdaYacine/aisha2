CREATE TABLE notifications (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    recipient_user_id uuid NOT NULL,
    event_type text NOT NULL,
    title_key text NOT NULL,
    body_key text NOT NULL,
    payload jsonb NOT NULL DEFAULT '{}'::jsonb,
    dedupe_key text NOT NULL,
    read_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT notifications_recipient_fk
        FOREIGN KEY (recipient_user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT notifications_event_type_not_blank CHECK (btrim(event_type) <> ''),
    CONSTRAINT notifications_title_key_not_blank CHECK (btrim(title_key) <> ''),
    CONSTRAINT notifications_body_key_not_blank CHECK (btrim(body_key) <> ''),
    CONSTRAINT notifications_dedupe_key_not_blank CHECK (btrim(dedupe_key) <> '')
);

CREATE UNIQUE INDEX notifications_recipient_dedupe_uq
    ON notifications (recipient_user_id, dedupe_key);
CREATE INDEX notifications_recipient_created_idx
    ON notifications (recipient_user_id, created_at DESC, id DESC);
CREATE INDEX notifications_recipient_unread_idx
    ON notifications (recipient_user_id, created_at DESC)
    WHERE read_at IS NULL;
