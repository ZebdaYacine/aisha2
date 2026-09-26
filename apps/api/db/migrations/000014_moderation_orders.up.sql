CREATE TABLE product_moderation_decisions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id uuid NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    submission_id uuid REFERENCES product_submissions(id) ON DELETE SET NULL,
    actor_user_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    action text NOT NULL CHECK (action IN ('APPROVE','REQUEST_CHANGES','REJECT','SUSPEND','ACTIVATE','ARCHIVE')),
    previous_status text NOT NULL,
    new_status text NOT NULL,
    reason text,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT product_moderation_reason_required CHECK (action NOT IN ('REQUEST_CHANGES','REJECT','SUSPEND') OR NULLIF(btrim(reason),'') IS NOT NULL)
);
CREATE INDEX product_moderation_product_created_idx ON product_moderation_decisions(product_id, created_at DESC);
CREATE INDEX product_moderation_action_created_idx ON product_moderation_decisions(action, created_at DESC);

CREATE TABLE orders (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number text NOT NULL UNIQUE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status text NOT NULL DEFAULT 'PENDING_PAYMENT' CHECK (status IN ('PENDING_PAYMENT','PAID','PREPARING','READY_TO_SHIP','SHIPPED','DELIVERED','CANCELLED','PAYMENT_FAILED','RETURNED')),
    currency char(3) NOT NULL CHECK (currency ~ '^[A-Z]{3}$'),
    subtotal_minor bigint NOT NULL CHECK (subtotal_minor >= 0),
    shipping_minor bigint NOT NULL DEFAULT 0 CHECK (shipping_minor >= 0),
    total_minor bigint NOT NULL CHECK (total_minor >= 0),
    address_snapshot jsonb NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    cancelled_at timestamptz
);
CREATE INDEX orders_user_created_idx ON orders(user_id, created_at DESC);
CREATE INDEX orders_status_created_idx ON orders(status, created_at);

CREATE TABLE order_items (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id uuid NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id uuid NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    artisan_profile_id uuid NOT NULL REFERENCES artisan_profiles(id) ON DELETE RESTRICT,
    workshop_id uuid NOT NULL REFERENCES workshops(id) ON DELETE RESTRICT,
    product_name text NOT NULL,
    artisan_name text NOT NULL,
    workshop_name text NOT NULL,
    unit_price_minor bigint NOT NULL CHECK (unit_price_minor > 0),
    currency char(3) NOT NULL CHECK (currency ~ '^[A-Z]{3}$'),
    quantity integer NOT NULL CHECK (quantity > 0),
    subtotal_minor bigint NOT NULL CHECK (subtotal_minor > 0),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX order_items_order_idx ON order_items(order_id);
CREATE INDEX order_items_artisan_idx ON order_items(artisan_profile_id, created_at DESC);

CREATE TABLE stock_reservations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id uuid NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id uuid NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity integer NOT NULL CHECK (quantity > 0),
    status text NOT NULL DEFAULT 'HELD' CHECK (status IN ('HELD','RELEASED','COMMITTED')),
    expires_at timestamptz NOT NULL,
    released_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(order_id, product_id)
);
CREATE INDEX stock_reservations_product_status_idx ON stock_reservations(product_id, status, expires_at);

CREATE TABLE payment_attempts (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id uuid NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    provider text NOT NULL DEFAULT 'PENDING_PROVIDER',
    status text NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING','CONFIRMED','FAILED','CANCELLED','REFUND_PENDING','REFUNDED')),
    amount_minor bigint NOT NULL CHECK (amount_minor >= 0),
    currency char(3) NOT NULL CHECK (currency ~ '^[A-Z]{3}$'),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX payment_attempts_order_idx ON payment_attempts(order_id, created_at DESC);

CREATE TABLE order_returns (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id uuid NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    actor_user_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    reason text NOT NULL,
    status text NOT NULL DEFAULT 'RECEIVED' CHECK (status IN ('REQUESTED','RECEIVED','INSPECTED','REFUNDED')),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT order_returns_reason_not_blank CHECK (btrim(reason) <> '')
);
CREATE INDEX order_returns_order_idx ON order_returns(order_id, created_at DESC);

CREATE TABLE shipment_events (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id uuid NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    status text NOT NULL,
    tracking_reference text,
    actor_user_id uuid REFERENCES users(id) ON DELETE SET NULL,
    occurred_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX shipment_events_order_time_idx ON shipment_events(order_id, occurred_at DESC);
