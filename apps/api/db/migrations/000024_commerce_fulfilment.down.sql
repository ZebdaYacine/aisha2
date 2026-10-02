DROP INDEX IF EXISTS shipment_events_idempotency_uq;
DROP TABLE IF EXISTS payment_refunds;
DROP INDEX IF EXISTS payment_attempts_provider_reference_uq;
ALTER TABLE payment_attempts DROP COLUMN IF EXISTS provider_reference, DROP COLUMN IF EXISTS failure_reason;
UPDATE orders SET status = 'PAID' WHERE status IN ('PREPARING', 'READY_TO_SHIP', 'SHIPPED', 'DELIVERED', 'PARTIALLY_REFUNDED', 'REFUNDED');
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_status_check;
ALTER TABLE orders
    ADD CONSTRAINT orders_status_check CHECK (status IN (
        'PENDING_PAYMENT', 'PAID', 'CANCELLED', 'PAYMENT_FAILED', 'RETURNED'
    ));
