ALTER TABLE notifications
    ADD COLUMN email_delivery_status text NOT NULL DEFAULT 'PENDING',
    ADD COLUMN email_delivery_attempted_at timestamptz,
    ADD COLUMN email_delivery_error text,
    ADD CONSTRAINT notifications_email_delivery_status_check
        CHECK (email_delivery_status IN ('PENDING', 'SENDING', 'SENT', 'FAILED'));

CREATE INDEX notifications_email_delivery_pending_idx
    ON notifications (created_at)
    WHERE email_delivery_status IN ('PENDING', 'SENDING');
