-- PRD-08 requires every product to belong to exactly one active/owned-workshop lineage.
-- Existing rows are assigned to each artisan's default workshop before the constraint is tightened.
INSERT INTO workshops (artisan_profile_id, name, wilaya, location_text, status, is_default, is_public)
SELECT a.id, COALESCE(NULLIF(btrim(a.workshop_name), ''), a.public_display_name), a.wilaya, a.location_text,
       CASE WHEN a.status = 'APPROVED' THEN 'ACTIVE' ELSE 'INACTIVE' END, true, a.status = 'APPROVED'
FROM artisan_profiles a
WHERE NOT EXISTS (
    SELECT 1 FROM workshops w WHERE w.artisan_profile_id = a.id AND w.is_default = true
);

UPDATE products p
SET workshop_id = w.id
FROM workshops w
WHERE p.workshop_id IS NULL
  AND w.artisan_profile_id = p.artisan_profile_id
  AND w.is_default = true;

ALTER TABLE products
    ALTER COLUMN workshop_id SET NOT NULL;
