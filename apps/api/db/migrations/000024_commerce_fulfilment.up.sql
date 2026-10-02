-- Complete the documented manual payment and warehouse fulfilment lifecycle.
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_status_check;
ALTER TABLE orders
    ADD CONSTRAINT orders_status_check CHECK (status IN (
        'PENDING_PAYMENT', 'PAID', 'PREPARING', 'READY_TO_SHIP', 'SHIPPED',
        'DELIVERED', 'CANCELLED', 'PAYMENT_FAILED', 'PARTIALLY_REFUNDED',
        'REFUNDED', 'RETURNED'
    ));

ALTER TABLE payment_attempts
    ADD COLUMN provider_reference text,
    ADD COLUMN failure_reason text;

CREATE UNIQUE INDEX payment_attempts_provider_reference_uq
    ON payment_attempts(provider_reference)
    WHERE provider_reference IS NOT NULL;

CREATE TABLE payment_refunds (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_attempt_id uuid NOT NULL REFERENCES payment_attempts(id) ON DELETE RESTRICT,
    order_id uuid NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    actor_user_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    idempotency_key text NOT NULL,
    amount_minor bigint NOT NULL CHECK (amount_minor > 0),
    currency char(3) NOT NULL CHECK (currency ~ '^[A-Z]{3}$'),
    reason text NOT NULL CHECK (btrim(reason) <> ''),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(actor_user_id, idempotency_key)
);
CREATE INDEX payment_refunds_order_idx ON payment_refunds(order_id, created_at DESC);

CREATE UNIQUE INDEX shipment_events_idempotency_uq
    ON shipment_events(order_id, status, COALESCE(tracking_reference, ''));
