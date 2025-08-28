-- +goose Up
-- If you have existing users, you must provide a default company_id
-- or update them before this migration can succeed as the column is NOT NULL.
-- For example, you might run these commands instead of the simple ADD COLUMN:
-- ALTER TABLE users ADD COLUMN company_id BIGINT;
-- UPDATE users SET company_id = <your_default_company_id> WHERE company_id IS NULL;
-- ALTER TABLE users ALTER COLUMN company_id SET NOT NULL;
ALTER TABLE users ADD COLUMN company_id BIGINT NOT NULL;

-- This assumes a 'companies' table exists with an 'id' primary key.
ALTER TABLE users ADD CONSTRAINT fk_users_company_id FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE users DROP CONSTRAINT fk_users_company_id;
ALTER TABLE users DROP COLUMN company_id;