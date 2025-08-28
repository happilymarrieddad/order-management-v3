-- +goose Up
-- This migration adds a nullable foreign key for a primary address to the users table.

-- Add the address_id column to the users table.
-- It's nullable because a user might not have a primary address.
ALTER TABLE users ADD COLUMN address_id BIGINT;

-- Add a foreign key constraint to ensure data integrity.
-- This assumes an 'addresses' table with a primary key 'id'.
-- ON DELETE SET NULL means if the referenced address is deleted, this user's address_id will be set to NULL.
ALTER TABLE users ADD CONSTRAINT fk_users_address FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE SET NULL;

-- Add an index on the new column for faster lookups, especially for joins.
CREATE INDEX idx_users_address_id ON users(address_id);

-- +goose Down
-- This migration removes the address_id column and its associated constraints from the users table.
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_address;
DROP INDEX IF EXISTS idx_users_address_id;
ALTER TABLE users DROP COLUMN IF EXISTS address_id;