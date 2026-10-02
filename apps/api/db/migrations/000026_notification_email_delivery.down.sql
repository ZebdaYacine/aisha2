DROP INDEX IF EXISTS notifications_email_delivery_pending_idx;

ALTER TABLE notifications
    DROP CONSTRAINT IF EXISTS notifications_email_delivery_status_check,
    DROP COLUMN IF EXISTS email_delivery_error,
    DROP COLUMN IF EXISTS email_delivery_attempted_at,
    DROP COLUMN IF EXISTS email_delivery_status;
