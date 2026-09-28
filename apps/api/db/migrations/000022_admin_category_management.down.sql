ALTER TABLE categories
    DROP CONSTRAINT IF EXISTS categories_benefit_rate_valid,
    DROP COLUMN IF EXISTS benefit_rate_basis_points;
