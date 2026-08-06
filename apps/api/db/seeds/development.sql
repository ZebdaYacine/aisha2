-- Deterministic development-only data for every application table.
-- Authentication and workflow records are intentionally inert.
BEGIN;

INSERT INTO roles (id, code) VALUES
    ('10000000-0000-0000-0000-000000000002', 'customer'),
    ('10000000-0000-0000-0000-000000000003', 'artisan'),
    ('10000000-0000-0000-0000-000000000006', 'administrator')
ON CONFLICT DO NOTHING;

INSERT INTO categories (id, slug, display_name, sort_order) VALUES
    ('20000000-0000-0000-0000-000000000001', 'decoration-and-art', 'Decoration and art', 10),
    ('20000000-0000-0000-0000-000000000003', 'jewellery-and-beauty-accessories', 'Jewellery and beauty accessories', 30),
    ('20000000-0000-0000-0000-000000000005', 'carpets-and-textiles', 'Carpets and textiles', 50)
ON CONFLICT DO NOTHING;

INSERT INTO users (id, email, status, display_name, email_verified_at, phone, created_at, updated_at) VALUES
    ('90000000-0000-0000-0000-000000000001', 'admin.seed@example.test', 'ACTIVE', 'Seed Administrator', '2025-01-01T00:00:00Z', NULL, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000002', 'customer.seed@example.test', 'ACTIVE', 'Seed Customer', '2025-01-01T00:00:00Z', '+213000000000', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000003', 'artisan.seed@example.test', 'ACTIVE', 'Atelier Tala', '2025-01-01T00:00:00Z', NULL, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id, assigned_by_user_id, assigned_at) VALUES
    ('90000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000006', NULL, '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000003', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO sessions (id, user_id, token_hash, family_id, expires_at, revoked_at, created_at) VALUES
    ('90000000-0000-0000-0000-000000000010', '90000000-0000-0000-0000-000000000002', 'seed-revoked-session-token-hash', '90000000-0000-0000-0000-000000000011', '2025-01-02T00:00:00Z', '2025-01-01T12:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO artisan_profiles (id, user_id, public_display_name, internal_name, workshop_name, wilaya, location_text, contact_email, contact_visibility, status, profile_image_object_key, created_at, updated_at, approved_at, submitted_at) VALUES
    ('90000000-0000-0000-0000-000000000020', '90000000-0000-0000-0000-000000000003', 'Atelier Tala', 'Seed Artisan', 'Atelier Tala', 'Tizi Ouzou', 'Kabylia', 'artisan.seed@example.test', 'PRIVATE', 'APPROVED', '/images/aisha/O5.jpg', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-01T12:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO artisan_profile_translations (artisan_profile_id, locale, biography, created_at, updated_at) VALUES
    ('90000000-0000-0000-0000-000000000020', 'en', 'A workshop preserving hand-shaped Algerian craft traditions.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000020', 'fr', 'Un atelier qui préserve les traditions artisanales algériennes façonnées à la main.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000020', 'ar', 'ورشة تحافظ على تقاليد الحرف الجزائرية اليدوية.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000020', 'es', 'Un taller que preserva las tradiciones artesanales argelinas hechas a mano.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO artisan_documents (id, artisan_profile_id, document_type, object_key, original_filename, media_type, size_bytes, checksum_sha256, created_at) VALUES
    ('90000000-0000-0000-0000-000000000030', '90000000-0000-0000-0000-000000000020', 'SEED_APPLICATION', 'seed/private/artisan-application.pdf', 'artisan-application.pdf', 'application/pdf', 1, repeat('0', 64), '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO artisan_media (id, artisan_profile_id, media_kind, object_key, original_filename, media_type, size_bytes, checksum_sha256, sort_order, visibility, created_at) VALUES
    ('90000000-0000-0000-0000-000000000040', '90000000-0000-0000-0000-000000000020', 'IMAGE', '/images/aisha/O5.jpg', 'O5.jpg', 'image/jpeg', 1, repeat('1', 64), 0, 'PUBLIC', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO artisan_profile_categories (artisan_profile_id, category_id, created_at) VALUES
    ('90000000-0000-0000-0000-000000000020', '20000000-0000-0000-0000-000000000003', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO products (id, artisan_profile_id, category_id, product_type, status, price_minor, currency, materials, production_method, intended_use, dimensions, weight_grams, region_of_origin, created_at, updated_at, published_at) VALUES
    ('90000000-0000-0000-0000-000000000050', '90000000-0000-0000-0000-000000000020', '20000000-0000-0000-0000-000000000003', 'ARTISAN_SPECIFIC', 'ACTIVE', 12500, 'EUR', 'Silver and enamel', 'Hand engraving and enamelling', 'Adornment', 'Seed sample', 120, 'Kabylia', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO product_translations (product_id, locale, name, description, story, cultural_context, created_at, updated_at) VALUES
    ('90000000-0000-0000-0000-000000000050', 'en', 'Kabyle silver brooch', 'A hand-enamelled silver brooch.', 'Engraved and enamelled by hand.', 'A development catalogue sample.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000050', 'fr', 'Broche kabyle en argent', 'Une broche en argent émaillée à la main.', 'Gravée et émaillée à la main.', 'Un exemple de catalogue de développement.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000050', 'ar', 'بروش فضي قبائلي', 'بروش فضي مطلي بالمينا يدويا.', 'منقوش ومطلي بالمينا يدويا.', 'عينة لكتالوج التطوير.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000050', 'es', 'Broche cabila de plata', 'Un broche de plata esmaltado a mano.', 'Grabado y esmaltado a mano.', 'Una muestra del catálogo de desarrollo.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO product_media (id, product_id, media_kind, object_key, original_filename, media_type, size_bytes, width_pixels, height_pixels, checksum_sha256, alt_text, sort_order, visibility, created_at) VALUES
    ('90000000-0000-0000-0000-000000000060', '90000000-0000-0000-0000-000000000050', 'IMAGE', '/images/aisha/O5.jpg', 'O5.jpg', 'image/jpeg', 1, 1080, 1080, repeat('2', 64), 'Kabyle silver brooch', 0, 'PUBLIC', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO addresses (id, user_id, full_name, phone, line1, city, postal_code, country, is_default, created_at, updated_at) VALUES
    ('90000000-0000-0000-0000-000000000070', '90000000-0000-0000-0000-000000000002', 'Seed Customer', '+213000000000', 'Development address', 'Tizi Ouzou', '00000', 'Algeria', true, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO password_reset_tokens (id, user_id, token_hash, expires_at, consumed_at, created_at) VALUES
    ('90000000-0000-0000-0000-000000000080', '90000000-0000-0000-0000-000000000002', 'seed-consumed-password-reset-token-hash', '2025-01-02T00:00:00Z', '2025-01-01T12:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO audit_events (id, event_type, actor_user_id, target_type, target_id, correlation_id, reason, previous_state, new_state, source_metadata, occurred_at) VALUES
    ('90000000-0000-0000-0000-000000000090', 'DEVELOPMENT_SEED_APPLIED', '90000000-0000-0000-0000-000000000001', 'database_seed', NULL, 'development-seed-v1', 'Deterministic development fixture', NULL, '{"seeded":true}'::jsonb, '{"source":"apps/api/db/seeds/development.sql"}'::jsonb, '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO outbox_events (id, event_type, aggregate_type, aggregate_id, payload, created_at, available_at, attempt_count, processed_at) VALUES
    ('90000000-0000-0000-0000-000000000100', 'DEVELOPMENT_SEED_COMPLETED', 'database_seed', '90000000-0000-0000-0000-000000000090', '{"seeded":true}'::jsonb, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z', 0, '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO idempotency_keys (id, actor_user_id, scope, idempotency_key, request_hash, state, response_status, response_body, resource_type, resource_id, created_at, expires_at, completed_at) VALUES
    ('90000000-0000-0000-0000-000000000110', '90000000-0000-0000-0000-000000000001', 'development-seed', 'development-seed-v1', repeat('3', 64), 'COMPLETED', 200, '{"seeded":true}'::jsonb, 'database_seed', '90000000-0000-0000-0000-000000000090', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

COMMIT;
