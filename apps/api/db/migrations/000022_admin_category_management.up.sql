ALTER TABLE categories
    ADD COLUMN benefit_rate_basis_points integer NOT NULL DEFAULT 0,
    ADD CONSTRAINT categories_benefit_rate_valid
        CHECK (benefit_rate_basis_points >= 0 AND benefit_rate_basis_points <= 10000);
