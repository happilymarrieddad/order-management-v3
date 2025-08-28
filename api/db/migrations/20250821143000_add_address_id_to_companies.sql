-- +goose Up
-- This assumes an 'addresses' table exists with an 'id' primary key.
-- ON DELETE RESTRICT prevents deleting an address if it is still in use by a company.
ALTER TABLE companies ADD CONSTRAINT fk_companies_address_id FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE RESTRICT;

-- +goose Down
ALTER TABLE companies DROP CONSTRAINT fk_companies_address_id;