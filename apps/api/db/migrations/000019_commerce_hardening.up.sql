-- A return is one lifecycle event per order.  The unique index makes the
-- idempotency guarantee safe when two fulfilment workers race.
CREATE UNIQUE INDEX order_returns_one_per_order_idx
    ON order_returns (order_id);
