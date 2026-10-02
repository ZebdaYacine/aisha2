-- Deterministic development-only data for every application table.
-- Authentication and workflow records are intentionally inert.
BEGIN;

-- Development only: make every seed run a clean, deterministic reset.
-- The migration metadata table is intentionally not included.
TRUNCATE TABLE
    audit_events, outbox_events, notifications, idempotency_keys,
    addresses, password_reset_tokens,
    inventory_movements,
    product_media, product_translations, product_submissions, products,
    workshop_translations, workshops,
    artisan_profile_categories, artisan_media, artisan_documents,
    artisan_profile_translations, artisan_verifications, artisan_memberships,
    artisan_profiles,
    category_translations, categories,
    order_returns, shipment_events, payment_attempts, stock_reservations,
    order_items, orders, product_moderation_decisions,
    warehouse_evidence, warehouse_inspections, warehouse_receptions,
    cart_items, carts, wishlist_items,
    sessions, user_roles, users, roles
RESTART IDENTITY CASCADE;

INSERT INTO roles (id, code) VALUES
    ('10000000-0000-0000-0000-000000000001', 'visitor'),
    ('10000000-0000-0000-0000-000000000002', 'customer'),
    ('10000000-0000-0000-0000-000000000003', 'artisan'),
    ('10000000-0000-0000-0000-000000000004', 'moderator'),
    ('10000000-0000-0000-0000-000000000005', 'warehouse_agent'),
    ('10000000-0000-0000-0000-000000000006', 'administrator')
ON CONFLICT DO NOTHING;

INSERT INTO categories (id, slug, display_name, sort_order) VALUES
    ('20000000-0000-0000-0000-000000000001', 'decoration-and-art', 'Decoration and art', 10),
    ('20000000-0000-0000-0000-000000000002', 'domestic-use', 'Domestic use', 20),
    ('20000000-0000-0000-0000-000000000003', 'jewellery-and-beauty-accessories', 'Jewellery and beauty accessories', 30),
    ('20000000-0000-0000-0000-000000000004', 'traditional-copper-products', 'Traditional copper products', 40),
    ('20000000-0000-0000-0000-000000000005', 'carpets-and-textiles', 'Carpets and textiles', 50),
    ('20000000-0000-0000-0000-000000000006', 'glassware', 'Glassware', 60),
    ('20000000-0000-0000-0000-000000000007', 'bamboo-and-halfa-products', 'Bamboo and halfa products', 70),
    ('20000000-0000-0000-0000-000000000008', 'traditional-musical-instruments', 'Traditional musical instruments', 80),
    ('20000000-0000-0000-0000-000000000009', 'souvenirs', 'Souvenirs', 90),
    ('20000000-0000-0000-0000-000000000010', 'leather-goods', 'Leather goods', 100),
    ('20000000-0000-0000-0000-000000000011', 'pottery-and-porcelain', 'Pottery and porcelain', 110),
    ('20000000-0000-0000-0000-000000000012', 'embroidery', 'Embroidery', 120)
ON CONFLICT DO NOTHING;

INSERT INTO category_translations (category_id, locale, name, created_at, updated_at) VALUES
    ('20000000-0000-0000-0000-000000000001', 'en', 'Decoration and art', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000001', 'fr', 'Décoration et art', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000001', 'ar', 'الديكور والفنون', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000001', 'es', 'Decoración y arte', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000002', 'en', 'Domestic use', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000002', 'fr', 'Usage domestique', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000002', 'ar', 'الاستعمال المنزلي', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000002', 'es', 'Uso doméstico', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000003', 'en', 'Jewellery and beauty accessories', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000003', 'fr', 'Bijoux et accessoires de beauté', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000003', 'ar', 'المجوهرات وإكسسوارات التجميل', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000003', 'es', 'Joyería y accesorios de belleza', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000004', 'en', 'Traditional copper products', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000004', 'fr', 'Produits traditionnels en cuivre', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000004', 'ar', 'المنتجات النحاسية التقليدية', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000004', 'es', 'Productos tradicionales de cobre', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000005', 'en', 'Carpets and textiles', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000005', 'fr', 'Tapis et textiles', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000005', 'ar', 'السجاد والمنسوجات', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000005', 'es', 'Alfombras y textiles', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000006', 'en', 'Glassware', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000006', 'fr', 'Verrerie', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000006', 'ar', 'الزجاجيات', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000006', 'es', 'Vidriería', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000007', 'en', 'Bamboo and halfa products', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000007', 'fr', 'Produits en bambou et en alfa', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000007', 'ar', 'منتجات الخيزران والحلفاء', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000007', 'es', 'Productos de bambú y esparto', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000008', 'en', 'Traditional musical instruments', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000008', 'fr', 'Instruments de musique traditionnels', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000008', 'ar', 'الآلات الموسيقية التقليدية', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000008', 'es', 'Instrumentos musicales tradicionales', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000009', 'en', 'Souvenirs', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000009', 'fr', 'Souvenirs', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000009', 'ar', 'التذكارات', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000009', 'es', 'Recuerdos', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000010', 'en', 'Leather goods', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000010', 'fr', 'Articles en cuir', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000010', 'ar', 'المنتجات الجلدية', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000010', 'es', 'Artículos de cuero', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000011', 'en', 'Pottery and porcelain', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000011', 'fr', 'Poterie et porcelaine', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000011', 'ar', 'الفخار والخزف', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000011', 'es', 'Cerámica y porcelana', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000012', 'en', 'Embroidery', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000012', 'fr', 'Broderie', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000012', 'ar', 'التطريز', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('20000000-0000-0000-0000-000000000012', 'es', 'Bordado', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT (category_id, locale) DO UPDATE SET name=EXCLUDED.name, updated_at=EXCLUDED.updated_at;

-- Development-only accounts use the bcrypt hash below. Keep the seed password
-- out of source control and provide it through the local development runbook.
INSERT INTO users (id, email, password_hash, status, display_name, email_verified_at, phone, created_at, updated_at) VALUES
    ('90000000-0000-0000-0000-000000000001', 'saad.admin@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Saad', '2025-01-01T00:00:00Z', '+213550000001', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000002', 'kader.agent@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Kader', '2025-01-01T00:00:00Z', '+213550000002', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000003', 'yassine.warehouse@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Yassine', '2025-01-01T00:00:00Z', '+213550000003', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT (id) DO UPDATE SET password_hash=EXCLUDED.password_hash;

INSERT INTO user_roles (user_id, role_id, assigned_by_user_id, assigned_at) VALUES
    ('90000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000006', NULL, '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000005', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000005', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO sessions (id, user_id, token_hash, family_id, expires_at, revoked_at, created_at) VALUES
    ('90000000-0000-0000-0000-000000000010', '90000000-0000-0000-0000-000000000002', 'seed-revoked-session-token-hash', '90000000-0000-0000-0000-000000000011', '2025-01-02T00:00:00Z', '2025-01-01T12:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO artisan_profiles (id, user_id, public_display_name, internal_name, workshop_name, wilaya, location_text, contact_email, contact_visibility, status, profile_image_object_key, created_at, updated_at, approved_at, submitted_at) VALUES
    ('90000000-0000-0000-0000-000000000020', '90000000-0000-0000-0000-000000000003', 'Oussama Atelier', 'Oussama Artisan', 'Oussama Atelier', 'Tizi Ouzou', 'Kabylia', 'oussama.artisan@example.test', 'PRIVATE', 'APPROVED', '/images/aisha/O5.jpg', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-01T12:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO artisan_profile_translations (artisan_profile_id, locale, display_name, biography, created_at, updated_at) VALUES
    ('90000000-0000-0000-0000-000000000020', 'en', 'Oussama Atelier', 'A workshop preserving hand-shaped Algerian craft traditions.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000020', 'fr', 'Atelier Oussama', 'Un atelier qui préserve les traditions artisanales algériennes façonnées à la main.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000020', 'ar', 'ورشة أسامة', 'ورشة تحافظ على تقاليد الحرف الجزائرية اليدوية.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000020', 'es', 'Taller Oussama', 'Un taller que preserva las tradiciones artesanales argelinas hechas a mano.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT (artisan_profile_id, locale) DO UPDATE SET display_name=EXCLUDED.display_name, biography=EXCLUDED.biography, updated_at=EXCLUDED.updated_at;

INSERT INTO artisan_documents (id, artisan_profile_id, document_type, object_key, original_filename, media_type, size_bytes, checksum_sha256, created_at) VALUES
    ('90000000-0000-0000-0000-000000000030', '90000000-0000-0000-0000-000000000020', 'SEED_APPLICATION', 'seed/private/artisan-application.pdf', 'artisan-application.pdf', 'application/pdf', 1, repeat('0', 64), '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO artisan_media (id, artisan_profile_id, media_kind, object_key, original_filename, media_type, size_bytes, checksum_sha256, sort_order, visibility, created_at) VALUES
    ('90000000-0000-0000-0000-000000000040', '90000000-0000-0000-0000-000000000020', 'IMAGE', '/images/aisha/O5.jpg', 'O5.jpg', 'image/jpeg', 1, repeat('1', 64), 0, 'PUBLIC', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO artisan_profile_categories (artisan_profile_id, category_id, created_at) VALUES
    ('90000000-0000-0000-0000-000000000020', '20000000-0000-0000-0000-000000000003', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO workshops (artisan_profile_id, name, description, wilaya, location_text, status, is_default, is_public) VALUES
    ('90000000-0000-0000-0000-000000000020', 'Oussama Atelier', 'A public development workshop preserving hand-shaped Algerian craft traditions.', 'Tizi Ouzou', 'Kabylia', 'ACTIVE', true, true)
ON CONFLICT (artisan_profile_id) WHERE is_default DO UPDATE
SET name = EXCLUDED.name, description = EXCLUDED.description, wilaya = EXCLUDED.wilaya,
    location_text = EXCLUDED.location_text, status = EXCLUDED.status, is_public = EXCLUDED.is_public,
    updated_at = CURRENT_TIMESTAMP;

INSERT INTO products (workshop_id, id, artisan_profile_id, category_id, product_type, status, price_minor, currency, materials, production_method, intended_use, dimensions, weight_grams, region_of_origin, created_at, updated_at, published_at) VALUES
    ((SELECT id FROM workshops WHERE artisan_profile_id = '90000000-0000-0000-0000-000000000020' AND is_default = true), '90000000-0000-0000-0000-000000000050', '90000000-0000-0000-0000-000000000020', '20000000-0000-0000-0000-000000000003', 'ARTISAN_SPECIFIC', 'ACTIVE', 12500, 'EUR', 'Silver and enamel', 'Hand engraving and enamelling', 'Adornment', 'Seed sample', 120, 'Kabylia', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z')
ON CONFLICT DO NOTHING;

UPDATE products p
SET workshop_id = w.id
FROM workshops w
WHERE w.artisan_profile_id = p.artisan_profile_id AND w.is_default = true AND p.workshop_id IS NULL;

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
    ('90000000-0000-0000-0000-000000000070', '90000000-0000-0000-0000-000000000002', 'Kader', '+213550000002', 'Development address', 'Tizi Ouzou', '15000', 'Algeria', true, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')
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

-- Additional development catalogue fixtures backed by images in
-- apps/web/public/images/aisha.
INSERT INTO users (id, email, password_hash, status, display_name, email_verified_at, phone, created_at, updated_at) VALUES
    ('90000000-0000-0000-0000-000000000004', 'lyna.customer@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Lyna', '2025-01-01T00:00:00Z', '+213550000004', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000005', 'oussama.artisan@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Oussama', '2025-01-01T00:00:00Z', '+213550000005', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000006', 'youcef.moderator@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Youcef', '2025-01-01T00:00:00Z', '+213550000006', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000007', 'legacy.removed@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Legacy fixture', '2025-01-01T00:00:00Z', NULL, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT (id) DO UPDATE SET password_hash=EXCLUDED.password_hash;

INSERT INTO user_roles (user_id, role_id, assigned_by_user_id, assigned_at) VALUES
    ('90000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000003', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000006', '10000000-0000-0000-0000-000000000004', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO artisan_profiles (id, user_id, public_display_name, internal_name, workshop_name, wilaya, location_text, contact_email, contact_visibility, status, profile_image_object_key, created_at, updated_at, approved_at, submitted_at) VALUES
    ('90000000-0000-0000-0000-000000000120', '90000000-0000-0000-0000-000000000004', 'Atelier Noura', 'Seed Jewellery Artisan', 'Atelier Noura', 'Sétif', 'High Plateaus', 'noura.seed@example.test', 'PRIVATE', 'APPROVED', '/images/aisha/654382945Y.jpg', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-01T12:00:00Z'),
    ('90000000-0000-0000-0000-000000000121', '90000000-0000-0000-0000-000000000005', 'Maison Tifawt', 'Seed Embroidery Artisan', 'Maison Tifawt', 'Tlemcen', 'Old Town', 'tifawt.seed@example.test', 'PRIVATE', 'APPROVED', '/images/aisha/C2.png', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-01T12:00:00Z'),
    ('90000000-0000-0000-0000-000000000122', '90000000-0000-0000-0000-000000000006', 'Atelier Zina', 'Seed Woodwork Artisan', 'Atelier Zina', 'Blida', 'Mitidja', 'zina.seed@example.test', 'PRIVATE', 'APPROVED', '/images/aisha/plateu en bois 1.jpg', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-01T12:00:00Z'),
    ('90000000-0000-0000-0000-000000000123', '90000000-0000-0000-0000-000000000007', 'Dar Tissili', 'Seed Textile Artisan', 'Dar Tissili', 'Ghardaïa', 'M''zab', 'tissili.seed@example.test', 'PRIVATE', 'APPROVED', '/images/aisha/Comme l’âme trouve sa paix dans un bel espace, le tapis révèle toute sa beauté là où il se sent chez lui.#homedecor #interiordesign #rugs #carpet#dubai.jpg', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-01T12:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO workshops (artisan_profile_id, name, description, wilaya, location_text, status, is_default, is_public) VALUES
    ('90000000-0000-0000-0000-000000000020', 'Oussama Atelier', 'A public development workshop preserving hand-shaped Algerian craft traditions.', 'Tizi Ouzou', 'Kabylia', 'ACTIVE', true, true),
    ('90000000-0000-0000-0000-000000000120', 'Atelier Noura', 'A public development workshop for jewellery and ceremonial pieces.', 'Sétif', 'High Plateaus', 'ACTIVE', true, true),
    ('90000000-0000-0000-0000-000000000121', 'Maison Tifawt', 'A public development atelier for embroidery and festive textiles.', 'Tlemcen', 'Old Town', 'ACTIVE', true, true),
    ('90000000-0000-0000-0000-000000000122', 'Atelier Zina', 'A public development workshop for painted woodwork.', 'Blida', 'Mitidja', 'ACTIVE', true, true),
    ('90000000-0000-0000-0000-000000000123', 'Dar Tissili', 'A public development textile house for woven decoration.', 'Ghardaïa', 'M''zab', 'ACTIVE', true, true)
ON CONFLICT (artisan_profile_id) WHERE is_default DO UPDATE
SET name = EXCLUDED.name, description = EXCLUDED.description, wilaya = EXCLUDED.wilaya,
    location_text = EXCLUDED.location_text, status = EXCLUDED.status, is_public = EXCLUDED.is_public,
    updated_at = CURRENT_TIMESTAMP;

INSERT INTO workshop_translations (workshop_id, locale, name, description)
SELECT w.id, 'en', w.name, w.description
FROM workshops w
WHERE w.is_default = true
ON CONFLICT (workshop_id, locale) DO UPDATE
SET name = EXCLUDED.name, description = EXCLUDED.description, updated_at = CURRENT_TIMESTAMP;

UPDATE products p
SET workshop_id = w.id
FROM workshops w
WHERE w.artisan_profile_id = p.artisan_profile_id AND w.is_default = true AND p.workshop_id IS NULL;

INSERT INTO artisan_profile_translations (artisan_profile_id, locale, display_name, biography, created_at, updated_at) VALUES
    ('90000000-0000-0000-0000-000000000120', 'en', 'Atelier Noura', 'A development workshop creating pearl and ceremonial jewellery.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000120', 'fr', 'Atelier Noura', 'Un atelier de développement qui crée des bijoux en perles et de cérémonie.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000120', 'ar', 'ورشة نورة', 'ورشة تطوير تصنع مجوهرات من اللؤلؤ وقطعاً احتفالية.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000120', 'es', 'Taller Noura', 'Un taller de desarrollo que crea joyas de perlas y piezas ceremoniales.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000121', 'en', 'Tifawt House', 'A development atelier translating festive dress and embroidery into contemporary pieces.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000121', 'fr', 'Maison Tifawt', 'Un atelier de développement qui revisite la tenue de fête et la broderie.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000121', 'ar', 'دار تيفاوت', 'ورشة تطوير تعيد تقديم أزياء المناسبات والتطريز بأسلوب معاصر.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000121', 'es', 'Casa Tifawt', 'Un taller de desarrollo que reinterpreta la ropa festiva y el bordado.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000122', 'en', 'Atelier Zina', 'A development workshop shaping wooden serving pieces with colourful geometric inlay.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000122', 'fr', 'Atelier Zina', 'Un atelier de développement qui façonne des pièces en bois décorées de motifs géométriques.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000122', 'ar', 'ورشة زينة', 'ورشة تطوير تشكل قطعاً خشبية مزينة بزخارف هندسية ملونة.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000122', 'es', 'Taller Zina', 'Un taller de desarrollo que crea piezas de madera con incrustaciones geométricas.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000123', 'en', 'Tissili House', 'A development textile house creating expressive woven pieces for modern interiors.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000123', 'fr', 'Maison Tissili', 'Une maison textile de développement qui crée des pièces tissées pour les intérieurs contemporains.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000123', 'ar', 'دار تيسيلي', 'دار نسيج للتطوير تصنع قطعاً منسوجة نابضة بالحياة للمساحات المعاصرة.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000123', 'es', 'Casa Tissili', 'Una casa textil de desarrollo que crea piezas tejidas expresivas para interiores modernos.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT (artisan_profile_id, locale) DO UPDATE SET display_name=EXCLUDED.display_name, biography=EXCLUDED.biography, updated_at=EXCLUDED.updated_at;

INSERT INTO artisan_media (id, artisan_profile_id, media_kind, object_key, original_filename, media_type, size_bytes, checksum_sha256, sort_order, visibility, created_at) VALUES
    ('90000000-0000-0000-0000-000000000170', '90000000-0000-0000-0000-000000000120', 'IMAGE', '/images/aisha/654382945Y.jpg', '654382945Y.jpg', 'image/jpeg', 101095, repeat('4', 64), 0, 'PUBLIC', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000171', '90000000-0000-0000-0000-000000000121', 'IMAGE', '/images/aisha/C2.png', 'C2.png', 'image/png', 2405470, repeat('5', 64), 0, 'PUBLIC', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000172', '90000000-0000-0000-0000-000000000122', 'IMAGE', '/images/aisha/plateu en bois 1.jpg', 'plateu en bois 1.jpg', 'image/jpeg', 266984, repeat('6', 64), 0, 'PUBLIC', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000173', '90000000-0000-0000-0000-000000000123', 'IMAGE', '/images/aisha/Comme l’âme trouve sa paix dans un bel espace, le tapis révèle toute sa beauté là où il se sent chez lui.#homedecor #interiordesign #rugs #carpet#dubai.jpg', 'carpet-development-sample.jpg', 'image/jpeg', 230685, repeat('7', 64), 0, 'PUBLIC', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO artisan_profile_categories (artisan_profile_id, category_id, created_at) VALUES
    ('90000000-0000-0000-0000-000000000120', '20000000-0000-0000-0000-000000000003', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000121', '20000000-0000-0000-0000-000000000012', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000122', '20000000-0000-0000-0000-000000000002', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000123', '20000000-0000-0000-0000-000000000005', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO products (workshop_id, id, artisan_profile_id, category_id, product_type, status, price_minor, currency, materials, production_method, intended_use, dimensions, weight_grams, region_of_origin, created_at, updated_at, published_at) VALUES
    ((SELECT id FROM workshops WHERE artisan_profile_id = '90000000-0000-0000-0000-000000000120' AND is_default = true), '90000000-0000-0000-0000-000000000150', '90000000-0000-0000-0000-000000000120', '20000000-0000-0000-0000-000000000003', 'ARTISAN_SPECIFIC', 'ACTIVE', 9800, 'EUR', 'Pearls and gold-tone metal', 'Hand assembly and finishing', 'Ceremonial wear', 'Development sample', 85, 'High Plateaus', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z'),
    ((SELECT id FROM workshops WHERE artisan_profile_id = '90000000-0000-0000-0000-000000000121' AND is_default = true), '90000000-0000-0000-0000-000000000151', '90000000-0000-0000-0000-000000000121', '20000000-0000-0000-0000-000000000012', 'ARTISAN_SPECIFIC', 'ACTIVE', 24500, 'EUR', 'Silk, textile and embroidery thread', 'Hand embroidery and tailoring', 'Festive wear', 'Development sample', 480, 'Tlemcen', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z'),
    ((SELECT id FROM workshops WHERE artisan_profile_id = '90000000-0000-0000-0000-000000000122' AND is_default = true), '90000000-0000-0000-0000-000000000152', '90000000-0000-0000-0000-000000000122', '20000000-0000-0000-0000-000000000002', 'ARTISAN_SPECIFIC', 'ACTIVE', 7600, 'EUR', 'Wood and painted inlay', 'Hand cutting, sanding and painting', 'Serving and display', 'Development sample', 950, 'Mitidja', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z'),
    ((SELECT id FROM workshops WHERE artisan_profile_id = '90000000-0000-0000-0000-000000000123' AND is_default = true), '90000000-0000-0000-0000-000000000153', '90000000-0000-0000-0000-000000000123', '20000000-0000-0000-0000-000000000005', 'ARTISAN_SPECIFIC', 'ACTIVE', 18900, 'EUR', 'Wool and cotton', 'Hand weaving and finishing', 'Interior decoration', 'Development sample', 1800, 'Mzab', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z')
ON CONFLICT DO NOTHING;

UPDATE products p
SET workshop_id = w.id
FROM workshops w
WHERE w.artisan_profile_id = p.artisan_profile_id AND w.is_default = true AND p.workshop_id IS NULL;

INSERT INTO product_translations (product_id, locale, name, description, story, cultural_context, created_at, updated_at) VALUES
    ('90000000-0000-0000-0000-000000000150', 'en', 'Pearl ceremonial necklace', 'A development sample of layered pearl jewellery.', 'Assembled by hand in a small workshop.', 'A contemporary fixture inspired by Algerian celebration jewellery.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000150', 'fr', 'Collier de cérémonie en perles', 'Un échantillon de développement de bijou en perles.', 'Assemblé à la main dans un petit atelier.', 'Une pièce contemporaine inspirée des bijoux de fête algériens.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000150', 'ar', 'قلادة احتفالية من اللؤلؤ', 'عينة تطوير لمجوهرات من اللؤلؤ المتعدد الطبقات.', 'جمعت يدوياً في ورشة صغيرة.', 'قطعة معاصرة مستوحاة من مجوهرات المناسبات الجزائرية.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000150', 'es', 'Collar ceremonial de perlas', 'Una muestra de desarrollo de joyería de perlas.', 'Montado a mano en un pequeño taller.', 'Una pieza contemporánea inspirada en la joyería festiva argelina.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000151', 'en', 'Violet embroidered caftan', 'A development sample with a vivid embroidered finish.', 'A study in festive colour, textile and handwork.', 'A contemporary interpretation of Algerian occasion dress.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000151', 'fr', 'Caftan violet brodé', 'Un échantillon de développement aux finitions brodées éclatantes.', 'Une étude de couleur, de textile et de travail manuel.', 'Une interprétation contemporaine de la tenue algérienne de cérémonie.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000151', 'ar', 'قفطان بنفسجي مطرز', 'عينة تطوير بلمسات تطريز زاهية.', 'دراسة في اللون والنسيج والعمل اليدوي.', 'تفسير معاصر لزي المناسبات الجزائري.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000151', 'es', 'Caftán violeta bordado', 'Una muestra de desarrollo con un acabado bordado intenso.', 'Un estudio de color, textil y trabajo manual.', 'Una interpretación contemporánea de la vestimenta festiva argelina.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000152', 'en', 'Painted wood serving tray', 'A hand-finished serving tray with colourful geometric decoration.', 'Woodwork and painted pattern meet in a practical household object.', 'A development piece for everyday hosting and display.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000152', 'fr', 'Plateau de service en bois peint', 'Un plateau de service fini à la main avec une décoration géométrique colorée.', 'Le travail du bois et le motif peint se rencontrent dans un objet du quotidien.', 'Une pièce de développement pour recevoir et décorer.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000152', 'ar', 'صينية تقديم خشبية مطلية', 'صينية تقديم مشغولة يدوياً بزخارف هندسية ملونة.', 'يجتمع العمل الخشبي والزخرفة المطلية في قطعة منزلية عملية.', 'قطعة تطوير للاستقبال وتزيين المنزل.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000152', 'es', 'Bandeja de madera pintada', 'Una bandeja de servicio acabada a mano con decoración geométrica colorida.', 'La madera y el patrón pintado se unen en un objeto práctico.', 'Una pieza de desarrollo para recibir y decorar.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000153', 'en', 'Woven colour field rug', 'A development sample of a bold woven rug for contemporary interiors.', 'Layered colour and texture give the room a warm focal point.', 'A modern textile study rooted in handwoven decoration.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000153', 'fr', 'Tapis tissé aux couleurs vives', 'Un échantillon de développement d’un tapis tissé pour les intérieurs contemporains.', 'Les couleurs et les textures superposées réchauffent l’espace.', 'Une étude textile moderne inspirée de la décoration tissée à la main.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000153', 'ar', 'سجادة منسوجة بألوان نابضة', 'عينة تطوير لسجادة منسوجة للمساحات المعاصرة.', 'تمنح طبقات اللون والملمس المكان نقطة دفء مميزة.', 'دراسة نسيجية عصرية مستوحاة من الزخرفة المنسوجة يدوياً.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000153', 'es', 'Alfombra tejida de colores', 'Una muestra de desarrollo de una alfombra tejida para interiores contemporáneos.', 'El color y la textura aportan calidez al espacio.', 'Un estudio textil moderno inspirado en la decoración tejida a mano.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO product_media (id, product_id, media_kind, object_key, original_filename, media_type, size_bytes, width_pixels, height_pixels, checksum_sha256, alt_text, sort_order, visibility, created_at) VALUES
    ('90000000-0000-0000-0000-000000000160', '90000000-0000-0000-0000-000000000150', 'IMAGE', '/images/aisha/654382945Y.jpg', '654382945Y.jpg', 'image/jpeg', 101095, 1080, 1146, repeat('8', 64), 'Pearl ceremonial necklace', 0, 'PUBLIC', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000161', '90000000-0000-0000-0000-000000000151', 'IMAGE', '/images/aisha/C2.png', 'C2.png', 'image/png', 2405470, 1024, 1536, repeat('9', 64), 'Violet embroidered caftan', 0, 'PUBLIC', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000162', '90000000-0000-0000-0000-000000000152', 'IMAGE', '/images/aisha/plateu en bois 1.jpg', 'plateu en bois 1.jpg', 'image/jpeg', 266984, 1080, 1080, repeat('a', 64), 'Painted wood serving tray', 0, 'PUBLIC', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000163', '90000000-0000-0000-0000-000000000153', 'IMAGE', '/images/aisha/Comme l’âme trouve sa paix dans un bel espace, le tapis révèle toute sa beauté là où il se sent chez lui.#homedecor #interiordesign #rugs #carpet#dubai.jpg', 'carpet-development-sample.jpg', 'image/jpeg', 230685, 1080, 1350, repeat('b', 64), 'Woven colour field rug', 0, 'PUBLIC', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

-- Keep the reset dataset focused on the six named demo accounts above. The
-- catalogue block remains useful as a source of representative fixtures, but
-- these older identities must not leak into a fresh local environment.
DELETE FROM product_media
WHERE product_id IN (
    '90000000-0000-0000-0000-000000000150',
    '90000000-0000-0000-0000-000000000151',
    '90000000-0000-0000-0000-000000000152',
    '90000000-0000-0000-0000-000000000153'
);
DELETE FROM product_translations
WHERE product_id IN (
    '90000000-0000-0000-0000-000000000150',
    '90000000-0000-0000-0000-000000000151',
    '90000000-0000-0000-0000-000000000152',
    '90000000-0000-0000-0000-000000000153'
);
DELETE FROM product_submissions
WHERE product_id IN (
    '90000000-0000-0000-0000-000000000150',
    '90000000-0000-0000-0000-000000000151',
    '90000000-0000-0000-0000-000000000152',
    '90000000-0000-0000-0000-000000000153'
);
DELETE FROM product_moderation_decisions
WHERE product_id IN (
    '90000000-0000-0000-0000-000000000150',
    '90000000-0000-0000-0000-000000000151',
    '90000000-0000-0000-0000-000000000152',
    '90000000-0000-0000-0000-000000000153'
);
DELETE FROM products
WHERE id IN (
    '90000000-0000-0000-0000-000000000150',
    '90000000-0000-0000-0000-000000000151',
    '90000000-0000-0000-0000-000000000152',
    '90000000-0000-0000-0000-000000000153'
);
DELETE FROM workshop_translations
WHERE workshop_id IN (
    SELECT id FROM workshops
    WHERE artisan_profile_id IN (
        '90000000-0000-0000-0000-000000000120',
        '90000000-0000-0000-0000-000000000121',
        '90000000-0000-0000-0000-000000000122',
        '90000000-0000-0000-0000-000000000123'
    )
);
DELETE FROM workshops
WHERE artisan_profile_id IN (
    '90000000-0000-0000-0000-000000000120',
    '90000000-0000-0000-0000-000000000121',
    '90000000-0000-0000-0000-000000000122',
    '90000000-0000-0000-0000-000000000123'
);
DELETE FROM artisan_profile_categories
WHERE artisan_profile_id IN (
    '90000000-0000-0000-0000-000000000120',
    '90000000-0000-0000-0000-000000000121',
    '90000000-0000-0000-0000-000000000122',
    '90000000-0000-0000-0000-000000000123'
);
DELETE FROM artisan_media
WHERE artisan_profile_id IN (
    '90000000-0000-0000-0000-000000000120',
    '90000000-0000-0000-0000-000000000121',
    '90000000-0000-0000-0000-000000000122',
    '90000000-0000-0000-0000-000000000123'
);
DELETE FROM artisan_documents
WHERE artisan_profile_id IN (
    '90000000-0000-0000-0000-000000000120',
    '90000000-0000-0000-0000-000000000121',
    '90000000-0000-0000-0000-000000000122',
    '90000000-0000-0000-0000-000000000123'
);
DELETE FROM artisan_profile_translations
WHERE artisan_profile_id IN (
    '90000000-0000-0000-0000-000000000120',
    '90000000-0000-0000-0000-000000000121',
    '90000000-0000-0000-0000-000000000122',
    '90000000-0000-0000-0000-000000000123'
);
DELETE FROM artisan_profiles
WHERE id IN (
    '90000000-0000-0000-0000-000000000120',
    '90000000-0000-0000-0000-000000000121',
    '90000000-0000-0000-0000-000000000122',
    '90000000-0000-0000-0000-000000000123'
);
DELETE FROM user_roles
WHERE user_id = '90000000-0000-0000-0000-000000000007';
DELETE FROM users
WHERE id = '90000000-0000-0000-0000-000000000007';

-- Expanded development directory: customers, approved artisans, moderators,
-- and warehouse agents. Accounts use the shared development bcrypt hash.
INSERT INTO users (id, email, password_hash, status, display_name, email_verified_at, phone, created_at, updated_at) VALUES
    ('90000000-0000-0000-0000-000000000200', 'nour.customer@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Nour', '2025-01-01T00:00:00Z', '+213560000200', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000201', 'saad.customer@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Saad Kader', '2025-01-01T00:00:00Z', '+213560000201', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000202', 'lyna.customer2@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Lyna B.', '2025-01-01T00:00:00Z', '+213560000202', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000203', 'amina.moderator@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Amina', '2025-01-01T00:00:00Z', '+213560000203', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000204', 'kader.moderator@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Kader Moderator', '2025-01-01T00:00:00Z', '+213560000204', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000205', 'yassine.agent@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Yassine Agent', '2025-01-01T00:00:00Z', '+213560000205', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000206', 'youcef.agent@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Youcef Agent', '2025-01-01T00:00:00Z', '+213560000206', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000207', 'karim.artisan@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Karim', '2025-01-01T00:00:00Z', '+213560000207', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000208', 'sarah.artisan@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Sarah', '2025-01-01T00:00:00Z', '+213560000208', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000209', 'oussama.artisan2@example.test', '$2a$12$NQsAT76BPVFvuz/4nUT2QO.WezqW96ahPazYygd9.uhWqrKJ05Dke', 'ACTIVE', 'Oussama K.', '2025-01-01T00:00:00Z', '+213560000209', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT (id) DO UPDATE SET
    email = EXCLUDED.email, display_name = EXCLUDED.display_name, phone = EXCLUDED.phone,
    status = EXCLUDED.status, password_hash = EXCLUDED.password_hash;

INSERT INTO user_roles (user_id, role_id, assigned_by_user_id, assigned_at) VALUES
    ('90000000-0000-0000-0000-000000000200', '10000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000201', '10000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000202', '10000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000203', '10000000-0000-0000-0000-000000000004', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000204', '10000000-0000-0000-0000-000000000004', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000205', '10000000-0000-0000-0000-000000000005', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000206', '10000000-0000-0000-0000-000000000005', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000207', '10000000-0000-0000-0000-000000000003', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000208', '10000000-0000-0000-0000-000000000003', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000209', '10000000-0000-0000-0000-000000000003', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

-- Three additional approved artisan profiles, each with a public workshop and
-- two images from the checked-in development image collection.
INSERT INTO artisan_profiles (id, user_id, public_display_name, internal_name, workshop_name, wilaya, location_text, contact_email, contact_visibility, status, profile_image_object_key, created_at, updated_at, approved_at, submitted_at) VALUES
    ('90000000-0000-0000-0000-000000000220', '90000000-0000-0000-0000-000000000207', 'Atelier Karim', 'Karim Woodcraft', 'Atelier Karim', 'Algiers', 'Bab El Oued', 'karim.artisan@example.test', 'PRIVATE', 'APPROVED', '/images/aisha/plateu en bois 1.jpg', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-01T12:00:00Z'),
    ('90000000-0000-0000-0000-000000000221', '90000000-0000-0000-0000-000000000208', 'Atelier Sarah', 'Sarah Textile Artisan', 'Atelier Sarah', 'Oran', 'Sidi El Houari', 'sarah.artisan@example.test', 'PRIVATE', 'APPROVED', '/images/aisha/C1.jpg', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-01T12:00:00Z'),
    ('90000000-0000-0000-0000-000000000222', '90000000-0000-0000-0000-000000000209', 'Atelier Oussama Kabyle', 'Oussama Kabyle Artisan', 'Atelier Oussama Kabyle', 'Béjaïa', 'Akbou', 'oussama.artisan2@example.test', 'PRIVATE', 'APPROVED', '/images/aisha/133T695O4.jpg', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-01T12:00:00Z')
ON CONFLICT (id) DO UPDATE SET
    public_display_name = EXCLUDED.public_display_name, workshop_name = EXCLUDED.workshop_name,
    wilaya = EXCLUDED.wilaya, location_text = EXCLUDED.location_text,
    profile_image_object_key = EXCLUDED.profile_image_object_key, status = EXCLUDED.status,
    approved_at = EXCLUDED.approved_at;

INSERT INTO artisan_profile_translations (artisan_profile_id, locale, display_name, biography, created_at, updated_at) VALUES
    ('90000000-0000-0000-0000-000000000220', 'en', 'Atelier Karim', 'A woodcraft studio creating practical pieces with Algerian geometric motifs.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000220', 'fr', 'Atelier Karim', 'Un atelier de bois qui crée des pièces pratiques aux motifs géométriques algériens.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000220', 'ar', 'ورشة كريم', 'ورشة خشب تصنع قطعاً عملية بزخارف هندسية جزائرية.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000220', 'es', 'Taller Karim', 'Un taller de madera que crea piezas prácticas con motivos geométricos argelinos.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000221', 'en', 'Atelier Sarah', 'A textile studio combining festive colour with careful hand embroidery.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000221', 'fr', 'Atelier Sarah', 'Un atelier textile qui associe couleurs festives et broderie soignée.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000221', 'ar', 'ورشة سارة', 'ورشة نسيج تجمع بين الألوان الاحتفالية والتطريز المتقن.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000221', 'es', 'Taller Sarah', 'Un taller textil que combina color festivo con bordado cuidadoso.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000222', 'en', 'Atelier Oussama Kabyle', 'A Kabyle studio celebrating silver, colour and hand-finished ornament.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000222', 'fr', 'Atelier Oussama Kabyle', 'Un atelier kabyle qui célèbre l’argent, la couleur et l’ornement façonné à la main.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000222', 'ar', 'ورشة أسامة القبائلية', 'ورشة قبائلية تحتفي بالفضة والألوان والزخارف المصنوعة يدوياً.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000222', 'es', 'Taller Oussama Cabilia', 'Un taller cabilio que celebra la plata, el color y el adorno hecho a mano.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT (artisan_profile_id, locale) DO UPDATE SET display_name = EXCLUDED.display_name, biography = EXCLUDED.biography, updated_at = EXCLUDED.updated_at;

INSERT INTO workshops (id, artisan_profile_id, name, description, wilaya, location_text, status, is_default, is_public) VALUES
    ('90000000-0000-0000-0000-000000000230', '90000000-0000-0000-0000-000000000220', 'Atelier Karim', 'Woodcraft and painted serving pieces from Algiers.', 'Algiers', 'Bab El Oued', 'ACTIVE', true, true),
    ('90000000-0000-0000-0000-000000000231', '90000000-0000-0000-0000-000000000221', 'Atelier Sarah', 'Embroidered festive textiles from Oran.', 'Oran', 'Sidi El Houari', 'ACTIVE', true, true),
    ('90000000-0000-0000-0000-000000000232', '90000000-0000-0000-0000-000000000222', 'Atelier Oussama Kabyle', 'Silver and colourful Kabyle ornament from Béjaïa.', 'Béjaïa', 'Akbou', 'ACTIVE', true, true)
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, status = EXCLUDED.status, is_public = EXCLUDED.is_public;

INSERT INTO workshop_translations (workshop_id, locale, name, description) VALUES
    ('90000000-0000-0000-0000-000000000230', 'en', 'Atelier Karim', 'Woodcraft and painted serving pieces from Algiers.'),
    ('90000000-0000-0000-0000-000000000230', 'fr', 'Atelier Karim', 'Objets en bois et plateaux peints fabriqués à Alger.'),
    ('90000000-0000-0000-0000-000000000230', 'ar', 'ورشة كريم', 'قطع خشبية وصوانٍ مطلية من الجزائر.'),
    ('90000000-0000-0000-0000-000000000230', 'es', 'Taller Karim', 'Piezas de madera y bandejas pintadas de Argel.'),
    ('90000000-0000-0000-0000-000000000231', 'en', 'Atelier Sarah', 'Embroidered festive textiles from Oran.'),
    ('90000000-0000-0000-0000-000000000231', 'fr', 'Atelier Sarah', 'Textiles festifs brodés d’Oran.'),
    ('90000000-0000-0000-0000-000000000231', 'ar', 'ورشة سارة', 'منسوجات احتفالية مطرزة من وهران.'),
    ('90000000-0000-0000-0000-000000000231', 'es', 'Taller Sarah', 'Textiles festivos bordados de Orán.'),
    ('90000000-0000-0000-0000-000000000232', 'en', 'Atelier Oussama Kabyle', 'Silver and colourful Kabyle ornament from Béjaïa.'),
    ('90000000-0000-0000-0000-000000000232', 'fr', 'Atelier Oussama Kabyle', 'Ornements kabyles en argent et en couleur de Béjaïa.'),
    ('90000000-0000-0000-0000-000000000232', 'ar', 'ورشة أسامة القبائلية', 'زخارف قبائلية فضية وملونة من بجاية.'),
    ('90000000-0000-0000-0000-000000000232', 'es', 'Taller Oussama Cabilia', 'Adornos cabilios de plata y color de Béjaïa.')
ON CONFLICT (workshop_id, locale) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, updated_at = CURRENT_TIMESTAMP;

INSERT INTO artisan_profile_categories (artisan_profile_id, category_id, created_at) VALUES
    ('90000000-0000-0000-0000-000000000220', '20000000-0000-0000-0000-000000000002', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000221', '20000000-0000-0000-0000-000000000012', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000222', '20000000-0000-0000-0000-000000000003', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO artisan_media (id, artisan_profile_id, media_kind, object_key, original_filename, media_type, size_bytes, checksum_sha256, sort_order, visibility, created_at) VALUES
    ('90000000-0000-0000-0000-000000000240', '90000000-0000-0000-0000-000000000220', 'IMAGE', '/images/aisha/plateu en bois 1.jpg', 'karim-woodwork.jpg', 'image/jpeg', 266984, repeat('c', 64), 0, 'PUBLIC', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000241', '90000000-0000-0000-0000-000000000220', 'IMAGE', '/images/aisha/Plateau en Bois avec Poignées ☕️#plateaux #plateau #bois #cuisinealgérienne #decocuisine (1).jpg', 'karim-tray.jpg', 'image/jpeg', 311240, repeat('d', 64), 1, 'PUBLIC', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000242', '90000000-0000-0000-0000-000000000221', 'IMAGE', '/images/aisha/C1.jpg', 'sarah-textile.jpg', 'image/jpeg', 182400, repeat('e', 64), 0, 'PUBLIC', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000243', '90000000-0000-0000-0000-000000000221', 'IMAGE', '/images/aisha/C2.png', 'sarah-embroidered-piece.png', 'image/png', 2405470, repeat('f', 64), 1, 'PUBLIC', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000244', '90000000-0000-0000-0000-000000000222', 'IMAGE', '/images/aisha/133T695O4.jpg', 'oussama-silver.jpg', 'image/jpeg', 215300, repeat('a', 64), 0, 'PUBLIC', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000245', '90000000-0000-0000-0000-000000000222', 'IMAGE', '/images/aisha/654382945Y.jpg', 'oussama-jewellery.jpg', 'image/jpeg', 101095, repeat('b', 64), 1, 'PUBLIC', '2025-01-01T00:00:00Z')
ON CONFLICT (id) DO NOTHING;

INSERT INTO artisan_memberships (user_id, artisan_profile_id, status, activated_at)
VALUES
    ('90000000-0000-0000-0000-000000000207', '90000000-0000-0000-0000-000000000220', 'ACTIVE', '2025-01-02T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000208', '90000000-0000-0000-0000-000000000221', 'ACTIVE', '2025-01-02T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000209', '90000000-0000-0000-0000-000000000222', 'ACTIVE', '2025-01-02T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO artisan_verifications (artisan_membership_id, status, decided_by_user_id, decided_at)
SELECT id, 'VERIFIED', '90000000-0000-0000-0000-000000000001', '2025-01-02T00:00:00Z'
FROM artisan_memberships
WHERE artisan_profile_id IN ('90000000-0000-0000-0000-000000000220', '90000000-0000-0000-0000-000000000221', '90000000-0000-0000-0000-000000000222')
ON CONFLICT DO NOTHING;

-- Oussama is the sole artisan demo account. The migration creates a membership
-- automatically for approved profiles, so the explicit upsert below also makes
-- this fixture safe when the seed is loaded against a migrated database.
UPDATE artisan_profiles
SET user_id = '90000000-0000-0000-0000-000000000005',
    public_display_name = 'Oussama Atelier',
    internal_name = 'Oussama Artisan',
    workshop_name = 'Oussama Atelier',
    contact_email = 'oussama.artisan@example.test'
WHERE id = '90000000-0000-0000-0000-000000000020';

INSERT INTO artisan_memberships (user_id, artisan_profile_id, status, activated_at)
VALUES ('90000000-0000-0000-0000-000000000005', '90000000-0000-0000-0000-000000000020', 'ACTIVE', '2025-01-02T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO artisan_verifications (artisan_membership_id, status, decided_by_user_id, decided_at)
SELECT id, 'VERIFIED', '90000000-0000-0000-0000-000000000001', '2025-01-02T00:00:00Z'
FROM artisan_memberships
WHERE artisan_profile_id = '90000000-0000-0000-0000-000000000020'
ON CONFLICT DO NOTHING;

INSERT INTO inventory_movements (product_id, movement_type, quantity_delta, reference_key, reason)
SELECT id, 'ACCEPTED', 5, 'development-seed:' || id::text, 'Accepted development seed stock'
FROM products
WHERE status = 'ACTIVE'
ON CONFLICT (reference_key) DO NOTHING;

-- Expanded development fixture requested for workflow testing.
-- The bcrypt hash below is for the development-only password supplied out of
-- band. Do not reuse it outside local development.
-- Counts represented by this block: 7 approved artisans, 3 moderators,
-- 2 warehouse agents, 10 customers, 20 orders, and 5 submitted applications.

-- Keep exactly two users assigned to the warehouse-agent role in this fixture.
DELETE FROM user_roles
WHERE role_id = '10000000-0000-0000-0000-000000000005'
  AND user_id IN ('90000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000003');
DELETE FROM password_reset_tokens WHERE user_id IN ('90000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000003');
DELETE FROM addresses WHERE user_id IN ('90000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000003');
DELETE FROM sessions WHERE user_id IN ('90000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000003');
DELETE FROM users WHERE id IN ('90000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000003');

INSERT INTO users (id, email, password_hash, status, display_name, email_verified_at, phone, created_at, updated_at) VALUES
    ('90000000-0000-0000-0000-000000000210', 'yacine.belkacem@example.test', '$2a$12$NeiBLeftB6nO684/gyPqlutzCffPZ9hSKhXKHPwd0n3JlkYdqqp6e', 'ACTIVE', 'Yacine Belkacem', '2025-01-01T00:00:00Z', '+213560000210', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000211', 'meriem.saidi@example.test', '$2a$12$NeiBLeftB6nO684/gyPqlutzCffPZ9hSKhXKHPwd0n3JlkYdqqp6e', 'ACTIVE', 'Meriem Saidi', '2025-01-01T00:00:00Z', '+213560000211', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000212', 'walid.amrani@example.test', '$2a$12$NeiBLeftB6nO684/gyPqlutzCffPZ9hSKhXKHPwd0n3JlkYdqqp6e', 'ACTIVE', 'Walid Amrani', '2025-01-01T00:00:00Z', '+213560000212', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000213', 'amel.benali@example.test', '$2a$12$NeiBLeftB6nO684/gyPqlutzCffPZ9hSKhXKHPwd0n3JlkYdqqp6e', 'ACTIVE', 'Amel Benali', '2025-01-01T00:00:00Z', '+213560000213', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000214', 'sofiane.haddad@example.test', '$2a$12$NeiBLeftB6nO684/gyPqlutzCffPZ9hSKhXKHPwd0n3JlkYdqqp6e', 'ACTIVE', 'Sofiane Haddad', '2025-01-01T00:00:00Z', '+213560000214', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000215', 'ines.touati@example.test', '$2a$12$NeiBLeftB6nO684/gyPqlutzCffPZ9hSKhXKHPwd0n3JlkYdqqp6e', 'ACTIVE', 'Ines Touati', '2025-01-01T00:00:00Z', '+213560000215', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000216', 'rachid.meziane@example.test', '$2a$12$NeiBLeftB6nO684/gyPqlutzCffPZ9hSKhXKHPwd0n3JlkYdqqp6e', 'ACTIVE', 'Rachid Meziane', '2025-01-01T00:00:00Z', '+213560000216', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000217', 'lina.boudiaf@example.test', '$2a$12$NeiBLeftB6nO684/gyPqlutzCffPZ9hSKhXKHPwd0n3JlkYdqqp6e', 'ACTIVE', 'Lina Boudiaf', '2025-01-01T00:00:00Z', '+213560000217', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000218', 'farid.cherif@example.test', '$2a$12$NeiBLeftB6nO684/gyPqlutzCffPZ9hSKhXKHPwd0n3JlkYdqqp6e', 'ACTIVE', 'Farid Cherif', '2025-01-01T00:00:00Z', '+213560000218', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT (id) DO UPDATE SET
    email = EXCLUDED.email, password_hash = EXCLUDED.password_hash,
    display_name = EXCLUDED.display_name, phone = EXCLUDED.phone,
    status = EXCLUDED.status, updated_at = EXCLUDED.updated_at;

INSERT INTO user_roles (user_id, role_id, assigned_by_user_id, assigned_at) VALUES
    ('90000000-0000-0000-0000-000000000210', '10000000-0000-0000-0000-000000000003', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000211', '10000000-0000-0000-0000-000000000003', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000212', '10000000-0000-0000-0000-000000000003', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000213', '10000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000214', '10000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000215', '10000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000216', '10000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000217', '10000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000218', '10000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000001', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

-- Three more approved artisans bring the approved artisan directory to seven.
INSERT INTO artisan_profiles (id, user_id, public_display_name, internal_name, workshop_name, wilaya, location_text, contact_email, contact_visibility, status, profile_image_object_key, created_at, updated_at, approved_at, submitted_at) VALUES
    ('90000000-0000-0000-0000-000000000223', '90000000-0000-0000-0000-000000000210', 'Atelier Yacine', 'Yacine Ceramics', 'Atelier Yacine', 'Constantine', 'Old Medina', 'yacine.belkacem@example.test', 'PRIVATE', 'APPROVED', '/images/aisha/A1.jpg', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-01T12:00:00Z'),
    ('90000000-0000-0000-0000-000000000224', '90000000-0000-0000-0000-000000000211', 'Atelier Meriem', 'Meriem Weaving', 'Atelier Meriem', 'Bouira', 'Lakhdaria', 'meriem.saidi@example.test', 'PRIVATE', 'APPROVED', '/images/aisha/C1.jpg', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-01T12:00:00Z'),
    ('90000000-0000-0000-0000-000000000225', '90000000-0000-0000-0000-000000000212', 'Atelier Walid', 'Walid Copperwork', 'Atelier Walid', 'Tlemcen', 'Mansourah', 'walid.amrani@example.test', 'PRIVATE', 'APPROVED', '/images/aisha/plateu en bois 1.jpg', '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-01T12:00:00Z')
ON CONFLICT (id) DO UPDATE SET
    public_display_name = EXCLUDED.public_display_name, internal_name = EXCLUDED.internal_name,
    workshop_name = EXCLUDED.workshop_name, wilaya = EXCLUDED.wilaya,
    location_text = EXCLUDED.location_text, contact_email = EXCLUDED.contact_email,
    profile_image_object_key = EXCLUDED.profile_image_object_key, status = EXCLUDED.status,
    approved_at = EXCLUDED.approved_at;

INSERT INTO artisan_profile_translations (artisan_profile_id, locale, display_name, biography, created_at, updated_at) VALUES
    ('90000000-0000-0000-0000-000000000223', 'en', 'Atelier Yacine', 'A ceramics studio shaped by the colours and forms of eastern Algeria.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000223', 'fr', 'Atelier Yacine', 'Un atelier de céramique inspiré par les couleurs de l’est algérien.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000223', 'ar', 'ورشة ياسين', 'ورشة خزف مستوحاة من ألوان وأشكال شرق الجزائر.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000223', 'es', 'Taller Yacine', 'Un estudio de cerámica inspirado en los colores del este de Argelia.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000224', 'en', 'Atelier Meriem', 'A weaving studio working with wool, texture and patient hand-finishing.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000224', 'fr', 'Atelier Meriem', 'Un atelier de tissage qui travaille la laine et les finitions patientes.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000224', 'ar', 'ورشة مريم', 'ورشة نسج تعمل بالصوف والملمس والتشطيبات اليدوية المتأنية.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000224', 'es', 'Taller Meriem', 'Un estudio de tejido que trabaja la lana y los acabados hechos a mano.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000225', 'en', 'Atelier Walid', 'A copper workshop combining traditional forms with a clean contemporary line.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000225', 'fr', 'Atelier Walid', 'Un atelier de cuivre qui associe formes traditionnelles et ligne contemporaine.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000225', 'ar', 'ورشة وليد', 'ورشة نحاس تجمع بين الأشكال التقليدية والخط المعاصر.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000225', 'es', 'Taller Walid', 'Un taller de cobre que combina formas tradicionales con una línea contemporánea.', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT (artisan_profile_id, locale) DO UPDATE SET display_name = EXCLUDED.display_name, biography = EXCLUDED.biography, updated_at = EXCLUDED.updated_at;

INSERT INTO workshops (id, artisan_profile_id, name, description, wilaya, location_text, status, is_default, is_public) VALUES
    ('90000000-0000-0000-0000-000000000233', '90000000-0000-0000-0000-000000000223', 'Atelier Yacine', 'Ceramics and hand-painted vessels from Constantine.', 'Constantine', '旧 Medina', 'ACTIVE', true, true),
    ('90000000-0000-0000-0000-000000000234', '90000000-0000-0000-0000-000000000224', 'Atelier Meriem', 'Woven wool and textile pieces from Bouira.', 'Bouira', 'Lakhdaria', 'ACTIVE', true, true),
    ('90000000-0000-0000-0000-000000000235', '90000000-0000-0000-0000-000000000225', 'Atelier Walid', 'Hand-finished copper pieces from Tlemcen.', 'Tlemcen', 'Mansourah', 'ACTIVE', true, true)
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, status = EXCLUDED.status, is_public = EXCLUDED.is_public;

INSERT INTO workshop_translations (workshop_id, locale, name, description) VALUES
    ('90000000-0000-0000-0000-000000000233', 'en', 'Atelier Yacine', 'Ceramics and hand-painted vessels from Constantine.'),
    ('90000000-0000-0000-0000-000000000233', 'fr', 'Atelier Yacine', 'Céramiques et vases peints à la main de Constantine.'),
    ('90000000-0000-0000-0000-000000000233', 'ar', 'ورشة ياسين', 'خزف وأوانٍ مطلية يدوياً من قسنطينة.'),
    ('90000000-0000-0000-0000-000000000233', 'es', 'Taller Yacine', 'Cerámica y vasijas pintadas a mano de Constantina.'),
    ('90000000-0000-0000-0000-000000000234', 'en', 'Atelier Meriem', 'Woven wool and textile pieces from Bouira.'),
    ('90000000-0000-0000-0000-000000000234', 'fr', 'Atelier Meriem', 'Pièces tissées en laine et textile de Bouira.'),
    ('90000000-0000-0000-0000-000000000234', 'ar', 'ورشة مريم', 'قطع منسوجة من الصوف والنسيج من البويرة.'),
    ('90000000-0000-0000-0000-000000000234', 'es', 'Taller Meriem', 'Piezas tejidas de lana y textil de Bouira.'),
    ('90000000-0000-0000-0000-000000000235', 'en', 'Atelier Walid', 'Hand-finished copper pieces from Tlemcen.'),
    ('90000000-0000-0000-0000-000000000235', 'fr', 'Atelier Walid', 'Pièces en cuivre finies à la main de Tlemcen.'),
    ('90000000-0000-0000-0000-000000000235', 'ar', 'ورشة وليد', 'قطع نحاسية مشغولة يدوياً من تلمسان.'),
    ('90000000-0000-0000-0000-000000000235', 'es', 'Taller Walid', 'Piezas de cobre acabadas a mano de Tremecén.')
ON CONFLICT (workshop_id, locale) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, updated_at = CURRENT_TIMESTAMP;

INSERT INTO artisan_profile_categories (artisan_profile_id, category_id, created_at) VALUES
    ('90000000-0000-0000-0000-000000000223', '20000000-0000-0000-0000-000000000011', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000224', '20000000-0000-0000-0000-000000000005', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000225', '20000000-0000-0000-0000-000000000004', '2025-01-01T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO artisan_media (id, artisan_profile_id, media_kind, object_key, original_filename, media_type, size_bytes, checksum_sha256, sort_order, visibility, created_at) VALUES
    ('90000000-0000-0000-0000-000000000246', '90000000-0000-0000-0000-000000000223', 'IMAGE', '/images/aisha/A1.jpg', 'yacine-ceramics.jpg', 'image/jpeg', 160000, repeat('c', 64), 0, 'PUBLIC', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000247', '90000000-0000-0000-0000-000000000224', 'IMAGE', '/images/aisha/Algeria art.jpg', 'meriem-weaving.jpg', 'image/jpeg', 182400, repeat('d', 64), 0, 'PUBLIC', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000248', '90000000-0000-0000-0000-000000000225', 'IMAGE', '/images/aisha/Snapinsta.app_471934208_18033457835416001_4267644784453957928_n_1080.jpg', 'walid-copperwork.jpg', 'image/jpeg', 266984, repeat('e', 64), 0, 'PUBLIC', '2025-01-01T00:00:00Z')
ON CONFLICT (id) DO NOTHING;

INSERT INTO artisan_memberships (user_id, artisan_profile_id, status, activated_at) VALUES
    ('90000000-0000-0000-0000-000000000210', '90000000-0000-0000-0000-000000000223', 'ACTIVE', '2025-01-02T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000211', '90000000-0000-0000-0000-000000000224', 'ACTIVE', '2025-01-02T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000212', '90000000-0000-0000-0000-000000000225', 'ACTIVE', '2025-01-02T00:00:00Z')
ON CONFLICT DO NOTHING;

INSERT INTO artisan_verifications (artisan_membership_id, status, decided_by_user_id, decided_at)
SELECT id, 'VERIFIED', '90000000-0000-0000-0000-000000000001', '2025-01-02T00:00:00Z'
FROM artisan_memberships
WHERE artisan_profile_id IN ('90000000-0000-0000-0000-000000000223', '90000000-0000-0000-0000-000000000224', '90000000-0000-0000-0000-000000000225')
ON CONFLICT DO NOTHING;

-- Development catalogue: exactly ten public workshops and thirty active
-- products across all twelve seeded categories. The additional workshops are
-- non-default locations owned by approved artisans; the four legacy product
-- ids are restored because the order fixtures below reference them.
INSERT INTO workshops (id, artisan_profile_id, name, description, wilaya, location_text, status, is_default, is_public) VALUES
    ('90000000-0000-0000-0000-000000000236', '90000000-0000-0000-0000-000000000223', 'Yacine Ceramic Annex', 'A second studio for kiln work and hand-painted vessels.', 'Constantine', 'El Khroub', 'ACTIVE', false, true),
    ('90000000-0000-0000-0000-000000000237', '90000000-0000-0000-0000-000000000224', 'Meriem Textile Room', 'A small weaving room for wool and linen collections.', 'Bouira', 'Sour El Ghozlane', 'ACTIVE', false, true),
    ('90000000-0000-0000-0000-000000000238', '90000000-0000-0000-0000-000000000225', 'Walid Copper Studio', 'A finishing studio for engraved and polished copper pieces.', 'Tlemcen', 'Chetouane', 'ACTIVE', false, true)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name, description = EXCLUDED.description, wilaya = EXCLUDED.wilaya,
    location_text = EXCLUDED.location_text, status = EXCLUDED.status, is_default = EXCLUDED.is_default,
    is_public = EXCLUDED.is_public, updated_at = CURRENT_TIMESTAMP;

INSERT INTO workshop_translations (workshop_id, locale, name, description) VALUES
    ('90000000-0000-0000-0000-000000000236', 'en', 'Yacine Ceramic Annex', 'A second studio for kiln work and hand-painted vessels.'),
    ('90000000-0000-0000-0000-000000000236', 'fr', 'Annexe céramique Yacine', 'Un second atelier pour la cuisson et les pièces peintes à la main.'),
    ('90000000-0000-0000-0000-000000000236', 'ar', 'ملحقة ياسين للخزف', 'ورشة ثانية للفخار والقطع المطلية يدوياً.'),
    ('90000000-0000-0000-0000-000000000236', 'es', 'Anexo cerámico Yacine', 'Un segundo estudio para piezas de cerámica pintadas a mano.'),
    ('90000000-0000-0000-0000-000000000237', 'en', 'Meriem Textile Room', 'A small weaving room for wool and linen collections.'),
    ('90000000-0000-0000-0000-000000000237', 'fr', 'Salle textile Meriem', 'Une petite salle de tissage pour les collections en laine et en lin.'),
    ('90000000-0000-0000-0000-000000000237', 'ar', 'قاعة مريم للنسيج', 'قاعة صغيرة لنسج مجموعات الصوف والكتان.'),
    ('90000000-0000-0000-0000-000000000237', 'es', 'Sala textil Meriem', 'Una pequeña sala de tejido para colecciones de lana y lino.'),
    ('90000000-0000-0000-0000-000000000238', 'en', 'Walid Copper Studio', 'A finishing studio for engraved and polished copper pieces.'),
    ('90000000-0000-0000-0000-000000000238', 'fr', 'Atelier cuivre Walid', 'Un atelier de finition pour des pièces en cuivre gravées et polies.'),
    ('90000000-0000-0000-0000-000000000238', 'ar', 'ورشة وليد للنحاس', 'ورشة لتشطيب القطع النحاسية المنقوشة والمصقولة.'),
    ('90000000-0000-0000-0000-000000000238', 'es', 'Estudio de cobre Walid', 'Un estudio de acabado para piezas de cobre grabadas y pulidas.')
ON CONFLICT (workshop_id, locale) DO UPDATE SET
    name = EXCLUDED.name, description = EXCLUDED.description, updated_at = CURRENT_TIMESTAMP;

UPDATE products
SET product_code = 'AISHA-' || right(replace(id::text, '-', ''), 10),
    planned_quantity = 12,
    order_total_minor = price_minor * 12
WHERE id = '90000000-0000-0000-0000-000000000050';

WITH seed_catalog(product_id, workshop_id, artisan_profile_id, category_id, price_minor, planned_quantity, name_en, name_fr, name_ar, name_es, media_key) AS (VALUES
    ('90000000-0000-0000-0000-000000000150'::uuid, '90000000-0000-0000-0000-000000000233'::uuid, '90000000-0000-0000-0000-000000000223'::uuid, '20000000-0000-0000-0000-000000000011'::uuid, 8900, 18, 'Constantine painted bowl', 'Bol peint de Constantine', 'وعاء مطلي من قسنطينة', 'Cuenco pintado de Constantina', '/images/aisha/654382945Y.jpg'),
    ('90000000-0000-0000-0000-000000000151'::uuid, '90000000-0000-0000-0000-000000000234'::uuid, '90000000-0000-0000-0000-000000000224'::uuid, '20000000-0000-0000-0000-000000000005'::uuid, 24500, 10, 'Bouira woven runner', 'Chemin de table tissé de Bouira', 'مفرش منسوج من البويرة', 'Camino de mesa tejido de Bouira', '/images/aisha/C2.png'),
    ('90000000-0000-0000-0000-000000000152'::uuid, '90000000-0000-0000-0000-000000000235'::uuid, '90000000-0000-0000-0000-000000000225'::uuid, '20000000-0000-0000-0000-000000000004'::uuid, 7600, 24, 'Tlemcen copper tray', 'Plateau en cuivre de Tlemcen', 'صينية نحاسية من تلمسان', 'Bandeja de cobre de Tremecén', '/images/aisha/plateu en bois 1.jpg'),
    ('90000000-0000-0000-0000-000000000153'::uuid, '90000000-0000-0000-0000-000000000230'::uuid, '90000000-0000-0000-0000-000000000220'::uuid, '20000000-0000-0000-0000-000000000002'::uuid, 18200, 14, 'Algiers geometric serving board', 'Planche de service géométrique d Alger', 'لوح تقديم هندسي من الجزائر', 'Tabla geométrica de servir de Argel', '/images/aisha/133T695O4.jpg'),
    ('90000000-0000-0000-0000-000000000400'::uuid, '90000000-0000-0000-0000-000000000231'::uuid, '90000000-0000-0000-0000-000000000221'::uuid, '20000000-0000-0000-0000-000000000012'::uuid, 12800, 16, 'Oran embroidered cushion', 'Coussin brodé d Oran', 'وسادة مطرزة من وهران', 'Cojín bordado de Orán', '/images/aisha/53567.png'),
    ('90000000-0000-0000-0000-000000000401'::uuid, '90000000-0000-0000-0000-000000000232'::uuid, '90000000-0000-0000-0000-000000000222'::uuid, '20000000-0000-0000-0000-000000000003'::uuid, 15700, 15, 'Kabyle silver pendant', 'Pendentif kabyle en argent', 'قلادة فضية قبائلية', 'Colgante cabilio de plata', '/images/aisha/6.jpg'),
    ('90000000-0000-0000-0000-000000000402'::uuid, '90000000-0000-0000-0000-000000000236'::uuid, '90000000-0000-0000-0000-000000000223'::uuid, '20000000-0000-0000-0000-000000000001'::uuid, 11200, 20, 'Blue ceramic wall tile', 'Carreau mural en céramique bleue', 'بلاطة جدارية خزفية زرقاء', 'Azulejo mural de cerámica azul', '/images/aisha/A1.jpg'),
    ('90000000-0000-0000-0000-000000000403'::uuid, '90000000-0000-0000-0000-000000000237'::uuid, '90000000-0000-0000-0000-000000000224'::uuid, '20000000-0000-0000-0000-000000000007'::uuid, 6900, 30, 'Halfa market basket', 'Panier de marché en alfa', 'سلة سوق من الحلفاء', 'Cesta de mercado de esparto', '/images/aisha/M1.jpg'),
    ('90000000-0000-0000-0000-000000000404'::uuid, '90000000-0000-0000-0000-000000000238'::uuid, '90000000-0000-0000-0000-000000000225'::uuid, '20000000-0000-0000-0000-000000000004'::uuid, 21400, 9, 'Engraved copper coffee set', 'Service à café en cuivre gravé', 'طقم قهوة نحاسي منقوش', 'Juego de café de cobre grabado', '/images/aisha/B1.jpg'),
    ('90000000-0000-0000-0000-000000000405'::uuid, '90000000-0000-0000-0000-000000000230'::uuid, '90000000-0000-0000-0000-000000000220'::uuid, '20000000-0000-0000-0000-000000000010'::uuid, 9800, 22, 'Leather card wallet', 'Porte-cartes en cuir', 'محفظة بطاقات جلدية', 'Cartera de tarjetas de cuero', '/images/aisha/jewee.jpg'),
    ('90000000-0000-0000-0000-000000000406'::uuid, '90000000-0000-0000-0000-000000000231'::uuid, '90000000-0000-0000-0000-000000000221'::uuid, '20000000-0000-0000-0000-000000000005'::uuid, 33500, 8, 'Oran wool wall hanging', 'Tissage mural en laine d Oran', 'نسيج جداري من صوف وهران', 'Tapiz mural de lana de Orán', '/images/aisha/C3.png'),
    ('90000000-0000-0000-0000-000000000407'::uuid, '90000000-0000-0000-0000-000000000232'::uuid, '90000000-0000-0000-0000-000000000222'::uuid, '20000000-0000-0000-0000-000000000009'::uuid, 5200, 35, 'Béjaïa keepsake magnet', 'Aimant souvenir de Béjaïa', 'مغناطيس تذكاري من بجاية', 'Imán recuerdo de Béjaïa', '/images/aisha/C1.jpg'),
    ('90000000-0000-0000-0000-000000000408'::uuid, '90000000-0000-0000-0000-000000000233'::uuid, '90000000-0000-0000-0000-000000000223'::uuid, '20000000-0000-0000-0000-000000000011'::uuid, 14600, 11, 'Hand-thrown serving dish', 'Plat de service tourné à la main', 'طبق تقديم مشغول يدوياً', 'Fuente de servir torneada a mano', '/images/aisha/C3.png'),
    ('90000000-0000-0000-0000-000000000409'::uuid, '90000000-0000-0000-0000-000000000234'::uuid, '90000000-0000-0000-0000-000000000224'::uuid, '20000000-0000-0000-0000-000000000012'::uuid, 17500, 13, 'Festive embroidered pouch', 'Pochette brodée de fête', 'حقيبة صغيرة مطرزة للمناسبات', 'Bolsa bordada festiva', '/images/aisha/CP1.jpg'),
    ('90000000-0000-0000-0000-000000000410'::uuid, '90000000-0000-0000-0000-000000000235'::uuid, '90000000-0000-0000-0000-000000000225'::uuid, '20000000-0000-0000-0000-000000000006'::uuid, 8300, 26, 'Tlemcen blue glass', 'Verre bleu de Tlemcen', 'كأس زجاجي أزرق من تلمسان', 'Vaso de vidrio azul de Tremecén', '/images/aisha/D1.jpg'),
    ('90000000-0000-0000-0000-000000000411'::uuid, '90000000-0000-0000-0000-000000000236'::uuid, '90000000-0000-0000-0000-000000000223'::uuid, '20000000-0000-0000-0000-000000000008'::uuid, 18900, 7, 'Ceramic frame drum', 'Tambour sur cadre en céramique', 'طبل إطار خزفي', 'Tambor de marco de cerámica', '/images/aisha/B2.jpg'),
    ('90000000-0000-0000-0000-000000000412'::uuid, '90000000-0000-0000-0000-000000000237'::uuid, '90000000-0000-0000-0000-000000000224'::uuid, '20000000-0000-0000-0000-000000000007'::uuid, 7400, 28, 'Palm fibre bread basket', 'Corbeille à pain en fibres de palmier', 'سلة خبز من ألياف النخيل', 'Cesta de pan de fibra de palma', '/images/aisha/output-onlinepngtools.png'),
    ('90000000-0000-0000-0000-000000000413'::uuid, '90000000-0000-0000-0000-000000000238'::uuid, '90000000-0000-0000-0000-000000000225'::uuid, '20000000-0000-0000-0000-000000000003'::uuid, 26800, 6, 'Polished silver earrings', 'Boucles d oreilles en argent poli', 'أقراط فضية مصقولة', 'Pendientes de plata pulida', '/images/aisha/765432.jpg'),
    ('90000000-0000-0000-0000-000000000414'::uuid, '90000000-0000-0000-0000-000000000230'::uuid, '90000000-0000-0000-0000-000000000220'::uuid, '20000000-0000-0000-0000-000000000002'::uuid, 10400, 18, 'Painted wooden tray', 'Plateau en bois peint', 'صينية خشبية مطلية', 'Bandeja de madera pintada', '/images/aisha/J1.jpg'),
    ('90000000-0000-0000-0000-000000000415'::uuid, '90000000-0000-0000-0000-000000000231'::uuid, '90000000-0000-0000-0000-000000000221'::uuid, '20000000-0000-0000-0000-000000000005'::uuid, 28900, 9, 'Handwoven wool cushion', 'Coussin en laine tissé à la main', 'وسادة صوفية منسوجة يدوياً', 'Cojín de lana tejido a mano', '/images/aisha/Algeria art.jpg'),
    ('90000000-0000-0000-0000-000000000416'::uuid, '90000000-0000-0000-0000-000000000232'::uuid, '90000000-0000-0000-0000-000000000222'::uuid, '20000000-0000-0000-0000-000000000009'::uuid, 6100, 32, 'Kabyle embroidered bookmark', 'Marque-page kabyle brodé', 'فاصل كتاب قبائلي مطرز', 'Marcapáginas cabilio bordado', '/images/aisha/Snapinsta.app_447895980_18152353162314221_3689087632318507403_n_1080.jpg'),
    ('90000000-0000-0000-0000-000000000417'::uuid, '90000000-0000-0000-0000-000000000233'::uuid, '90000000-0000-0000-0000-000000000223'::uuid, '20000000-0000-0000-0000-000000000001'::uuid, 13200, 15, 'Constantine ceramic vase', 'Vase en céramique de Constantine', 'مزهرية خزفية من قسنطينة', 'Jarrón de cerámica de Constantina', '/images/aisha/Snapinsta.app_410159085_18133321372314221_4586148325445799764_n_1080.jpg'),
    ('90000000-0000-0000-0000-000000000418'::uuid, '90000000-0000-0000-0000-000000000234'::uuid, '90000000-0000-0000-0000-000000000224'::uuid, '20000000-0000-0000-0000-000000000012'::uuid, 19600, 12, 'Linen embroidered tablecloth', 'Nappe en lin brodée', 'مفرش طاولة من الكتان المطرز', 'Mantel de lino bordado', '/images/aisha/Snapinsta.app_470795950_18209401987292380_4883742063971950135_n_1080.jpg'),
    ('90000000-0000-0000-0000-000000000419'::uuid, '90000000-0000-0000-0000-000000000235'::uuid, '90000000-0000-0000-0000-000000000225'::uuid, '20000000-0000-0000-0000-000000000004'::uuid, 22500, 8, 'Copper incense holder', 'Brûle-parfum en cuivre', 'مبخرة نحاسية', 'Quemador de incienso de cobre', '/images/aisha/Snapinsta.app_455131776_903354394960467_3623373114008527603_n_1080.jpg'),
    ('90000000-0000-0000-0000-000000000420'::uuid, '90000000-0000-0000-0000-000000000236'::uuid, '90000000-0000-0000-0000-000000000223'::uuid, '20000000-0000-0000-0000-000000000011'::uuid, 9900, 21, 'Hand-painted ceramic cup', 'Tasse en céramique peinte à la main', 'كوب خزفي مطلي يدوياً', 'Taza de cerámica pintada a mano', '/images/aisha/Snapinsta.app_470984486_18032514275416001_5509346141223573205_n_1080.jpg'),
    ('90000000-0000-0000-0000-000000000421'::uuid, '90000000-0000-0000-0000-000000000237'::uuid, '90000000-0000-0000-0000-000000000224'::uuid, '20000000-0000-0000-0000-000000000007'::uuid, 5800, 40, 'Woven storage basket', 'Panier de rangement tissé', 'سلة تخزين منسوجة', 'Cesta de almacenamiento tejida', '/images/aisha/Snapinsta.app_409974168_18133323640314221_5974977855731823254_n_1080.jpg'),
    ('90000000-0000-0000-0000-000000000422'::uuid, '90000000-0000-0000-0000-000000000238'::uuid, '90000000-0000-0000-0000-000000000225'::uuid, '20000000-0000-0000-0000-000000000010'::uuid, 14300, 17, 'Leather belt with brass buckle', 'Ceinture en cuir avec boucle en laiton', 'حزام جلدي بإبزيم نحاسي', 'Cinturón de cuero con hebilla de latón', '/images/aisha/Snapinsta.app_450074780_18155221924314221_934381476914220603_n_1080.jpg'),
    ('90000000-0000-0000-0000-000000000423'::uuid, '90000000-0000-0000-0000-000000000230'::uuid, '90000000-0000-0000-0000-000000000220'::uuid, '20000000-0000-0000-0000-000000000002'::uuid, 11800, 19, 'Wooden spice box', 'Boîte à épices en bois', 'علبة توابل خشبية', 'Caja de especias de madera', '/images/aisha/Copy of 366366855_681414600673448_8672563481987395173_n.jpg'),
    ('90000000-0000-0000-0000-000000000424'::uuid, '90000000-0000-0000-0000-000000000231'::uuid, '90000000-0000-0000-0000-000000000221'::uuid, '20000000-0000-0000-0000-000000000006'::uuid, 8700, 23, 'Hand-blown glass tumbler', 'Verre soufflé à la main', 'كأس زجاجي منفخ يدوياً', 'Vaso de vidrio soplado a mano', '/images/aisha/Plateaux Rond 🥏#plateaux #plateau #cuisine #decohome (3).jpg')
),
inserted_products AS (
    INSERT INTO products (id, workshop_id, artisan_profile_id, category_id, product_code, product_type, status, price_minor, currency, materials, production_method, intended_use, dimensions, weight_grams, region_of_origin, planned_quantity, order_total_minor, created_at, updated_at, published_at)
    SELECT s.product_id, s.workshop_id, s.artisan_profile_id, s.category_id,
           'AISHA-' || right(replace(s.product_id::text, '-', ''), 10),
           'ARTISAN_SPECIFIC', 'ACTIVE', s.price_minor, 'EUR', 'Locally sourced craft materials', 'Handmade in the workshop', 'Home and personal use', 'Development collection item', 220, 'Algeria', s.planned_quantity,
           s.price_minor * s.planned_quantity, '2025-01-01T00:00:00Z', '2025-01-02T00:00:00Z', '2025-01-02T00:00:00Z'
    FROM seed_catalog s
    ON CONFLICT (id) DO UPDATE SET
        workshop_id = EXCLUDED.workshop_id, artisan_profile_id = EXCLUDED.artisan_profile_id,
        category_id = EXCLUDED.category_id, product_code = EXCLUDED.product_code,
        status = EXCLUDED.status, price_minor = EXCLUDED.price_minor,
        planned_quantity = EXCLUDED.planned_quantity, order_total_minor = EXCLUDED.order_total_minor,
        updated_at = EXCLUDED.updated_at, published_at = EXCLUDED.published_at
    RETURNING id
),
inserted_translations AS (
    INSERT INTO product_translations (product_id, locale, name, description, story, cultural_context, created_at, updated_at)
    SELECT s.product_id, l.locale,
           CASE l.locale WHEN 'en' THEN s.name_en WHEN 'fr' THEN s.name_fr WHEN 'ar' THEN s.name_ar ELSE s.name_es END,
           CASE l.locale WHEN 'en' THEN 'A handmade Algerian craft piece from a verified development workshop.' WHEN 'fr' THEN 'Une pièce artisanale algérienne fabriquée dans un atelier de développement vérifié.' WHEN 'ar' THEN 'قطعة حرفية جزائرية مصنوعة في ورشة تطوير موثوقة.' ELSE 'Una pieza artesanal argelina hecha en un taller de desarrollo verificado.' END,
           CASE l.locale WHEN 'en' THEN 'Made slowly by hand with local materials.' WHEN 'fr' THEN 'Fabriquée lentement à la main avec des matériaux locaux.' WHEN 'ar' THEN 'مصنوعة يدوياً بعناية باستعمال مواد محلية.' ELSE 'Hecha a mano con materiales locales.' END,
           CASE l.locale WHEN 'en' THEN 'A development catalogue sample.' WHEN 'fr' THEN 'Un échantillon du catalogue de développement.' WHEN 'ar' THEN 'عينة من كتالوج التطوير.' ELSE 'Una muestra del catálogo de desarrollo.' END,
           '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'
    FROM seed_catalog s
    CROSS JOIN (VALUES ('en'::varchar(5)), ('fr'::varchar(5)), ('ar'::varchar(5)), ('es'::varchar(5))) AS l(locale)
    JOIN inserted_products p ON p.id = s.product_id
    ON CONFLICT (product_id, locale) DO UPDATE SET
        name = EXCLUDED.name, description = EXCLUDED.description, story = EXCLUDED.story,
        cultural_context = EXCLUDED.cultural_context, updated_at = EXCLUDED.updated_at
    RETURNING product_id
)
INSERT INTO product_media (id, product_id, media_kind, object_key, original_filename, media_type, size_bytes, width_pixels, height_pixels, checksum_sha256, alt_text, sort_order, visibility, created_at)
SELECT md5(s.product_id::text || ':development-media')::uuid, s.product_id, 'IMAGE', s.media_key,
       'development-' || right(replace(s.product_id::text, '-', ''), 6) || '.jpg',
       CASE WHEN right(lower(s.media_key), 4) = '.png' THEN 'image/png' ELSE 'image/jpeg' END,
       120000, 1080, 1080, repeat('9', 64), s.name_en, 0, 'PUBLIC', '2025-01-01T00:00:00Z'
FROM seed_catalog s
JOIN inserted_products p ON p.id = s.product_id
ON CONFLICT (object_key) DO NOTHING;

-- The product rows above are active catalogue fixtures, so give each one
-- accepted development stock for customer availability and warehouse search.
INSERT INTO inventory_movements (product_id, movement_type, quantity_delta, reference_key, reason)
SELECT id, 'ACCEPTED', planned_quantity, 'development-catalogue:' || id::text, 'Accepted development catalogue stock'
FROM products
WHERE id IN (
    '90000000-0000-0000-0000-000000000150', '90000000-0000-0000-0000-000000000151',
    '90000000-0000-0000-0000-000000000152', '90000000-0000-0000-0000-000000000153',
    '90000000-0000-0000-0000-000000000400', '90000000-0000-0000-0000-000000000401',
    '90000000-0000-0000-0000-000000000402', '90000000-0000-0000-0000-000000000403',
    '90000000-0000-0000-0000-000000000404', '90000000-0000-0000-0000-000000000405',
    '90000000-0000-0000-0000-000000000406', '90000000-0000-0000-0000-000000000407',
    '90000000-0000-0000-0000-000000000408', '90000000-0000-0000-0000-000000000409',
    '90000000-0000-0000-0000-000000000410', '90000000-0000-0000-0000-000000000411',
    '90000000-0000-0000-0000-000000000412', '90000000-0000-0000-0000-000000000413',
    '90000000-0000-0000-0000-000000000414', '90000000-0000-0000-0000-000000000415',
    '90000000-0000-0000-0000-000000000416', '90000000-0000-0000-0000-000000000417',
    '90000000-0000-0000-0000-000000000418', '90000000-0000-0000-0000-000000000419',
    '90000000-0000-0000-0000-000000000420', '90000000-0000-0000-0000-000000000421',
    '90000000-0000-0000-0000-000000000422', '90000000-0000-0000-0000-000000000423',
    '90000000-0000-0000-0000-000000000424'
)
ON CONFLICT (reference_key) DO NOTHING;

-- Five submitted artisan applications are kept separate from the seven
-- approved artisan accounts so the moderator queue has realistic work.
INSERT INTO artisan_profiles (id, user_id, public_display_name, internal_name, workshop_name, wilaya, location_text, contact_email, contact_visibility, status, created_at, updated_at, submitted_at) VALUES
    ('90000000-0000-0000-0000-000000000260', '90000000-0000-0000-0000-000000000213', 'Amel Benali Studio', 'Amel Benali Application', 'Atelier Amel', 'Algiers', 'Kouba', 'amel.benali@example.test', 'PRIVATE', 'SUBMITTED', '2025-01-03T00:00:00Z', '2025-01-03T00:00:00Z', '2025-01-03T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000261', '90000000-0000-0000-0000-000000000214', 'Sofiane Leather', 'Sofiane Haddad Application', 'Maison Sofiane', 'Batna', 'Lambèse', 'sofiane.haddad@example.test', 'PRIVATE', 'SUBMITTED', '2025-01-03T00:00:00Z', '2025-01-03T00:00:00Z', '2025-01-03T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000262', '90000000-0000-0000-0000-000000000215', 'Ines Embroidery', 'Ines Touati Application', 'Atelier Ines', 'Médéa', 'Berrouaghia', 'ines.touati@example.test', 'PRIVATE', 'SUBMITTED', '2025-01-03T00:00:00Z', '2025-01-03T00:00:00Z', '2025-01-03T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000263', '90000000-0000-0000-0000-000000000216', 'Rachid Wood Studio', 'Rachid Meziane Application', 'Atelier Rachid', 'Béjaïa', 'Sidi Aïch', 'rachid.meziane@example.test', 'PRIVATE', 'SUBMITTED', '2025-01-03T00:00:00Z', '2025-01-03T00:00:00Z', '2025-01-03T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000264', '90000000-0000-0000-0000-000000000217', 'Lina Basketry', 'Lina Boudiaf Application', 'Atelier Lina', 'Djelfa', 'Aïn Oussera', 'lina.boudiaf@example.test', 'PRIVATE', 'SUBMITTED', '2025-01-03T00:00:00Z', '2025-01-03T00:00:00Z', '2025-01-03T00:00:00Z')
ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, submitted_at = EXCLUDED.submitted_at, updated_at = EXCLUDED.updated_at;

INSERT INTO artisan_profile_translations (artisan_profile_id, locale, display_name, biography, created_at, updated_at) VALUES
    ('90000000-0000-0000-0000-000000000260', 'en', 'Amel Benali Studio', 'A new studio applying hand decoration to everyday objects.', '2025-01-03T00:00:00Z', '2025-01-03T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000261', 'en', 'Sofiane Leather', 'A leather workshop focused on durable small goods.', '2025-01-03T00:00:00Z', '2025-01-03T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000262', 'en', 'Ines Embroidery', 'An embroidery application rooted in family techniques.', '2025-01-03T00:00:00Z', '2025-01-03T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000263', 'en', 'Rachid Wood Studio', 'A proposed wood studio creating useful household pieces.', '2025-01-03T00:00:00Z', '2025-01-03T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000264', 'en', 'Lina Basketry', 'A basketry application using locally sourced fibres.', '2025-01-03T00:00:00Z', '2025-01-03T00:00:00Z')
ON CONFLICT (artisan_profile_id, locale) DO UPDATE SET display_name = EXCLUDED.display_name, biography = EXCLUDED.biography, updated_at = EXCLUDED.updated_at;

INSERT INTO artisan_documents (id, artisan_profile_id, document_type, object_key, original_filename, media_type, size_bytes, checksum_sha256, created_at) VALUES
    ('90000000-0000-0000-0000-000000000270', '90000000-0000-0000-0000-000000000260', 'IDENTITY_DOCUMENT', 'seed/private/applications/amel-benali-id.pdf', 'amel-benali-id.pdf', 'application/pdf', 12000, repeat('a', 64), '2025-01-03T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000271', '90000000-0000-0000-0000-000000000261', 'IDENTITY_DOCUMENT', 'seed/private/applications/sofiane-haddad-id.pdf', 'sofiane-haddad-id.pdf', 'application/pdf', 12000, repeat('b', 64), '2025-01-03T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000272', '90000000-0000-0000-0000-000000000262', 'IDENTITY_DOCUMENT', 'seed/private/applications/ines-touati-id.pdf', 'ines-touati-id.pdf', 'application/pdf', 12000, repeat('c', 64), '2025-01-03T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000273', '90000000-0000-0000-0000-000000000263', 'IDENTITY_DOCUMENT', 'seed/private/applications/rachid-meziane-id.pdf', 'rachid-meziane-id.pdf', 'application/pdf', 12000, repeat('d', 64), '2025-01-03T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000274', '90000000-0000-0000-0000-000000000264', 'IDENTITY_DOCUMENT', 'seed/private/applications/lina-boudiaf-id.pdf', 'lina-boudiaf-id.pdf', 'application/pdf', 12000, repeat('e', 64), '2025-01-03T00:00:00Z')
ON CONFLICT (id) DO NOTHING;

INSERT INTO addresses (id, user_id, full_name, phone, line1, city, postal_code, country, is_default, created_at, updated_at) VALUES
    ('90000000-0000-0000-0000-000000000700', '90000000-0000-0000-0000-000000000004', 'Lyna', '+213550000004', '12 Rue Didouche Mourad', 'Algiers', '16000', 'Algeria', true, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000701', '90000000-0000-0000-0000-000000000200', 'Nour', '+213560000200', '8 Rue Larbi Ben M''hidi', 'Oran', '31000', 'Algeria', true, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000702', '90000000-0000-0000-0000-000000000201', 'Saad Kader', '+213560000201', '5 Rue de la Paix', 'Tizi Ouzou', '15000', 'Algeria', true, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000703', '90000000-0000-0000-0000-000000000202', 'Lyna B.', '+213560000202', '21 Rue Emir Abdelkader', 'Sétif', '19000', 'Algeria', true, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000704', '90000000-0000-0000-0000-000000000213', 'Amel Benali', '+213560000213', '14 Rue des Palmiers', 'Algiers', '16050', 'Algeria', true, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000705', '90000000-0000-0000-0000-000000000214', 'Sofiane Haddad', '+213560000214', '3 Rue des Aurès', 'Batna', '05000', 'Algeria', true, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000706', '90000000-0000-0000-0000-000000000215', 'Ines Touati', '+213560000215', '7 Rue de la Poste', 'Médéa', '26000', 'Algeria', true, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000707', '90000000-0000-0000-0000-000000000216', 'Rachid Meziane', '+213560000216', '10 Rue de la Casbah', 'Béjaïa', '06000', 'Algeria', true, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000708', '90000000-0000-0000-0000-000000000217', 'Lina Boudiaf', '+213560000217', '19 Rue des Jardins', 'Djelfa', '17000', 'Algeria', true, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z'),
    ('90000000-0000-0000-0000-000000000709', '90000000-0000-0000-0000-000000000218', 'Farid Cherif', '+213560000218', '2 Rue Ibn Khaldoun', 'Blida', '09000', 'Algeria', true, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')
ON CONFLICT (id) DO UPDATE SET full_name = EXCLUDED.full_name, phone = EXCLUDED.phone, line1 = EXCLUDED.line1, city = EXCLUDED.city, postal_code = EXCLUDED.postal_code, country = EXCLUDED.country, is_default = EXCLUDED.is_default, updated_at = EXCLUDED.updated_at;

-- Twenty realistic orders distributed across the ten customer accounts.
WITH seed_orders(order_id, order_number, user_id, address_id, status, product_id, quantity, created_at) AS (VALUES
    ('90000000-0000-0000-0000-000000000300', 'AISHA-DEV-0300', '90000000-0000-0000-0000-000000000004', '90000000-0000-0000-0000-000000000700', 'PAID', '90000000-0000-0000-0000-000000000050', 1, '2025-02-01T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000301', 'AISHA-DEV-0301', '90000000-0000-0000-0000-000000000200', '90000000-0000-0000-0000-000000000701', 'PAID', '90000000-0000-0000-0000-000000000150', 1, '2025-02-02T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000302', 'AISHA-DEV-0302', '90000000-0000-0000-0000-000000000201', '90000000-0000-0000-0000-000000000702', 'PREPARING', '90000000-0000-0000-0000-000000000151', 1, '2025-02-03T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000303', 'AISHA-DEV-0303', '90000000-0000-0000-0000-000000000202', '90000000-0000-0000-0000-000000000703', 'SHIPPED', '90000000-0000-0000-0000-000000000152', 2, '2025-02-04T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000304', 'AISHA-DEV-0304', '90000000-0000-0000-0000-000000000213', '90000000-0000-0000-0000-000000000704', 'DELIVERED', '90000000-0000-0000-0000-000000000153', 1, '2025-02-05T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000305', 'AISHA-DEV-0305', '90000000-0000-0000-0000-000000000214', '90000000-0000-0000-0000-000000000705', 'PAID', '90000000-0000-0000-0000-000000000050', 1, '2025-02-06T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000306', 'AISHA-DEV-0306', '90000000-0000-0000-0000-000000000215', '90000000-0000-0000-0000-000000000706', 'PAID', '90000000-0000-0000-0000-000000000150', 2, '2025-02-07T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000307', 'AISHA-DEV-0307', '90000000-0000-0000-0000-000000000216', '90000000-0000-0000-0000-000000000707', 'READY_TO_SHIP', '90000000-0000-0000-0000-000000000151', 1, '2025-02-08T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000308', 'AISHA-DEV-0308', '90000000-0000-0000-0000-000000000217', '90000000-0000-0000-0000-000000000708', 'CANCELLED', '90000000-0000-0000-0000-000000000152', 1, '2025-02-09T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000309', 'AISHA-DEV-0309', '90000000-0000-0000-0000-000000000218', '90000000-0000-0000-0000-000000000709', 'PAYMENT_FAILED', '90000000-0000-0000-0000-000000000153', 1, '2025-02-10T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000310', 'AISHA-DEV-0310', '90000000-0000-0000-0000-000000000004', '90000000-0000-0000-0000-000000000700', 'PENDING_PAYMENT', '90000000-0000-0000-0000-000000000150', 1, '2025-02-11T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000311', 'AISHA-DEV-0311', '90000000-0000-0000-0000-000000000200', '90000000-0000-0000-0000-000000000701', 'PAID', '90000000-0000-0000-0000-000000000151', 2, '2025-02-12T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000312', 'AISHA-DEV-0312', '90000000-0000-0000-0000-000000000201', '90000000-0000-0000-0000-000000000702', 'DELIVERED', '90000000-0000-0000-0000-000000000152', 1, '2025-02-13T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000313', 'AISHA-DEV-0313', '90000000-0000-0000-0000-000000000202', '90000000-0000-0000-0000-000000000703', 'SHIPPED', '90000000-0000-0000-0000-000000000153', 1, '2025-02-14T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000314', 'AISHA-DEV-0314', '90000000-0000-0000-0000-000000000213', '90000000-0000-0000-0000-000000000704', 'PREPARING', '90000000-0000-0000-0000-000000000050', 2, '2025-02-15T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000315', 'AISHA-DEV-0315', '90000000-0000-0000-0000-000000000214', '90000000-0000-0000-0000-000000000705', 'PAID', '90000000-0000-0000-0000-000000000150', 1, '2025-02-16T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000316', 'AISHA-DEV-0316', '90000000-0000-0000-0000-000000000215', '90000000-0000-0000-0000-000000000706', 'READY_TO_SHIP', '90000000-0000-0000-0000-000000000151', 1, '2025-02-17T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000317', 'AISHA-DEV-0317', '90000000-0000-0000-0000-000000000216', '90000000-0000-0000-0000-000000000707', 'DELIVERED', '90000000-0000-0000-0000-000000000152', 2, '2025-02-18T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000318', 'AISHA-DEV-0318', '90000000-0000-0000-0000-000000000217', '90000000-0000-0000-0000-000000000708', 'PAID', '90000000-0000-0000-0000-000000000153', 1, '2025-02-19T10:00:00Z'::timestamptz),
    ('90000000-0000-0000-0000-000000000319', 'AISHA-DEV-0319', '90000000-0000-0000-0000-000000000218', '90000000-0000-0000-0000-000000000709', 'PENDING_PAYMENT', '90000000-0000-0000-0000-000000000050', 1, '2025-02-20T10:00:00Z'::timestamptz)
)
INSERT INTO orders (id, order_number, user_id, status, currency, subtotal_minor, shipping_minor, total_minor, address_snapshot, created_at, updated_at)
SELECT s.order_id::uuid, s.order_number, s.user_id::uuid, s.status, p.currency, p.price_minor * s.quantity, 0, p.price_minor * s.quantity,
       jsonb_build_object('fullName', a.full_name, 'phone', a.phone, 'line1', a.line1, 'city', a.city, 'postalCode', a.postal_code, 'country', a.country), s.created_at, s.created_at
FROM seed_orders s JOIN products p ON p.id = s.product_id::uuid JOIN addresses a ON a.id = s.address_id::uuid
ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, total_minor = EXCLUDED.total_minor, address_snapshot = EXCLUDED.address_snapshot, updated_at = EXCLUDED.updated_at;

WITH seed_lines(order_id, product_id, quantity) AS (VALUES
    ('90000000-0000-0000-0000-000000000300', '90000000-0000-0000-0000-000000000050', 1), ('90000000-0000-0000-0000-000000000301', '90000000-0000-0000-0000-000000000150', 1),
    ('90000000-0000-0000-0000-000000000302', '90000000-0000-0000-0000-000000000151', 1), ('90000000-0000-0000-0000-000000000303', '90000000-0000-0000-0000-000000000152', 2),
    ('90000000-0000-0000-0000-000000000304', '90000000-0000-0000-0000-000000000153', 1), ('90000000-0000-0000-0000-000000000305', '90000000-0000-0000-0000-000000000050', 1),
    ('90000000-0000-0000-0000-000000000306', '90000000-0000-0000-0000-000000000150', 2), ('90000000-0000-0000-0000-000000000307', '90000000-0000-0000-0000-000000000151', 1),
    ('90000000-0000-0000-0000-000000000308', '90000000-0000-0000-0000-000000000152', 1), ('90000000-0000-0000-0000-000000000309', '90000000-0000-0000-0000-000000000153', 1),
    ('90000000-0000-0000-0000-000000000310', '90000000-0000-0000-0000-000000000150', 1), ('90000000-0000-0000-0000-000000000311', '90000000-0000-0000-0000-000000000151', 2),
    ('90000000-0000-0000-0000-000000000312', '90000000-0000-0000-0000-000000000152', 1), ('90000000-0000-0000-0000-000000000313', '90000000-0000-0000-0000-000000000153', 1),
    ('90000000-0000-0000-0000-000000000314', '90000000-0000-0000-0000-000000000050', 2), ('90000000-0000-0000-0000-000000000315', '90000000-0000-0000-0000-000000000150', 1),
    ('90000000-0000-0000-0000-000000000316', '90000000-0000-0000-0000-000000000151', 1), ('90000000-0000-0000-0000-000000000317', '90000000-0000-0000-0000-000000000152', 2),
    ('90000000-0000-0000-0000-000000000318', '90000000-0000-0000-0000-000000000153', 1), ('90000000-0000-0000-0000-000000000319', '90000000-0000-0000-0000-000000000050', 1)
)
INSERT INTO order_items (id, order_id, product_id, artisan_profile_id, workshop_id, product_name, artisan_name, workshop_name, unit_price_minor, currency, quantity, subtotal_minor)
SELECT gen_random_uuid(), s.order_id::uuid, p.id, p.artisan_profile_id, p.workshop_id,
       COALESCE(pt.name, p.product_type), a.public_display_name, w.name, p.price_minor, p.currency, s.quantity, p.price_minor * s.quantity
FROM seed_lines s JOIN products p ON p.id = s.product_id::uuid JOIN artisan_profiles a ON a.id = p.artisan_profile_id JOIN workshops w ON w.id = p.workshop_id
LEFT JOIN product_translations pt ON pt.product_id = p.id AND pt.locale = 'en'
WHERE NOT EXISTS (SELECT 1 FROM order_items oi WHERE oi.order_id = s.order_id::uuid);

INSERT INTO payment_attempts (order_id, provider, status, amount_minor, currency, failure_reason)
SELECT o.id, CASE WHEN o.status IN ('PENDING_PAYMENT', 'PAYMENT_FAILED', 'CANCELLED') THEN 'PENDING_PROVIDER' ELSE 'manual' END,
       CASE WHEN o.status = 'PENDING_PAYMENT' THEN 'PENDING' WHEN o.status = 'PAYMENT_FAILED' THEN 'FAILED' WHEN o.status = 'CANCELLED' THEN 'CANCELLED' ELSE 'CONFIRMED' END,
       o.total_minor, o.currency, CASE WHEN o.status = 'PAYMENT_FAILED' THEN 'Development payment failure fixture' ELSE NULL END
FROM orders o
WHERE o.order_number LIKE 'AISHA-DEV-%'
  AND NOT EXISTS (SELECT 1 FROM payment_attempts p WHERE p.order_id = o.id);

INSERT INTO shipment_events (order_id, status, occurred_at)
SELECT o.id, CASE WHEN o.status = 'PENDING_PAYMENT' THEN 'PENDING' ELSE o.status END, o.created_at
FROM orders o
WHERE o.order_number LIKE 'AISHA-DEV-%'
  AND NOT EXISTS (SELECT 1 FROM shipment_events e WHERE e.order_id = o.id);

INSERT INTO stock_reservations (order_id, product_id, quantity, status, expires_at)
SELECT oi.order_id, oi.product_id, oi.quantity, 'HELD', CURRENT_TIMESTAMP + INTERVAL '30 minutes'
FROM order_items oi JOIN orders o ON o.id = oi.order_id
WHERE o.status = 'PENDING_PAYMENT'
  AND NOT EXISTS (SELECT 1 FROM stock_reservations sr WHERE sr.order_id = oi.order_id AND sr.product_id = oi.product_id);

COMMIT;
