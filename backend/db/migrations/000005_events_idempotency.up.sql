CREATE TABLE audit_events (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type text NOT NULL,
    actor_user_id uuid,
    target_type text NOT NULL,
    target_id uuid,
    correlation_id text,
    reason text,
    previous_state jsonb,
    new_state jsonb,
    source_metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
    occurred_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT audit_events_actor_fk
        FOREIGN KEY (actor_user_id) REFERENCES users (id) ON DELETE RESTRICT,
    CONSTRAINT audit_events_type_not_blank CHECK (btrim(event_type) <> ''),
    CONSTRAINT audit_events_target_type_not_blank CHECK (btrim(target_type) <> '')
);
CREATE INDEX audit_events_actor_time_idx ON audit_events (actor_user_id, occurred_at DESC);
CREATE INDEX audit_events_target_idx ON audit_events (target_type, target_id, occurred_at DESC);
CREATE INDEX audit_events_type_time_idx ON audit_events (event_type, occurred_at DESC);
CREATE INDEX audit_events_correlation_idx ON audit_events (correlation_id) WHERE correlation_id IS NOT NULL;

CREATE FUNCTION reject_audit_event_mutation()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    RAISE EXCEPTION 'audit_events are append-only';
END;
$$;

CREATE TRIGGER audit_events_reject_update
    BEFORE UPDATE ON audit_events
    FOR EACH ROW EXECUTE FUNCTION reject_audit_event_mutation();

CREATE TRIGGER audit_events_reject_delete
    BEFORE DELETE ON audit_events
    FOR EACH ROW EXECUTE FUNCTION reject_audit_event_mutation();

CREATE TABLE outbox_events (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type text NOT NULL,
    aggregate_type text NOT NULL,
    aggregate_id uuid NOT NULL,
    payload jsonb NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    available_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    attempt_count integer NOT NULL DEFAULT 0,
    locked_at timestamptz,
    processed_at timestamptz,
    last_error text,
    CONSTRAINT outbox_events_type_not_blank CHECK (btrim(event_type) <> ''),
    CONSTRAINT outbox_events_aggregate_type_not_blank CHECK (btrim(aggregate_type) <> ''),
    CONSTRAINT outbox_events_attempt_count_non_negative CHECK (attempt_count >= 0)
);
CREATE INDEX outbox_events_pending_idx
    ON outbox_events (available_at, created_at)
    WHERE processed_at IS NULL;
CREATE INDEX outbox_events_aggregate_idx ON outbox_events (aggregate_type, aggregate_id);

CREATE TABLE idempotency_keys (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_user_id uuid,
    scope text NOT NULL,
    idempotency_key text NOT NULL,
    request_hash char(64) NOT NULL,
    state text NOT NULL DEFAULT 'PROCESSING',
    response_status integer,
    response_body jsonb,
    resource_type text,
    resource_id uuid,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at timestamptz NOT NULL,
    completed_at timestamptz,
    CONSTRAINT idempotency_keys_actor_fk
        FOREIGN KEY (actor_user_id) REFERENCES users (id) ON DELETE RESTRICT,
    CONSTRAINT idempotency_keys_scope_not_blank CHECK (btrim(scope) <> ''),
    CONSTRAINT idempotency_keys_key_not_blank CHECK (btrim(idempotency_key) <> ''),
    CONSTRAINT idempotency_keys_state_valid CHECK (state IN ('PROCESSING', 'COMPLETED', 'FAILED')),
    CONSTRAINT idempotency_keys_response_status_valid
        CHECK (response_status IS NULL OR response_status BETWEEN 100 AND 599),
    CONSTRAINT idempotency_keys_expiry_after_creation CHECK (expires_at > created_at),
    CONSTRAINT idempotency_keys_actor_scope_key_unique UNIQUE NULLS NOT DISTINCT
        (actor_user_id, scope, idempotency_key)
);
CREATE INDEX idempotency_keys_expiry_idx ON idempotency_keys (expires_at);
CREATE INDEX idempotency_keys_resource_idx
    ON idempotency_keys (resource_type, resource_id)
    WHERE resource_id IS NOT NULL;
