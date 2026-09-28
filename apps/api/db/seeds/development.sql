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

-- Development-only password for every seeded account: Yassine1996@Got.
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

COMMIT;
