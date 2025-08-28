-- +goose Up
-- This migration changes the 'id' column of the 'users' table from UUID to BIGINT.
--
-- WARNING: This is a DESTRUCTIVE operation. It is intended for early development
-- stages where resetting the database is acceptable. If you have existing user data
-- and foreign key relationships, this will break them. A more complex data
-- migration strategy would be required for a production environment.

-- The CASCADE will drop any foreign key constraints that reference users.id
ALTER TABLE users DROP COLUMN id CASCADE;

-- Add the new ID column as a BIGSERIAL, which is an auto-incrementing 8-byte integer.
ALTER TABLE users ADD COLUMN id BIGSERIAL PRIMARY KEY;


-- +goose Down
-- Reverting this migration is also DESTRUCTIVE and will result in data loss.

-- You may need to ensure the uuid-ossp extension is available in your database.
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
ALTER TABLE users DROP COLUMN id CASCADE;