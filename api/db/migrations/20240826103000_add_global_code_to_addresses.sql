-- +goose Up
-- This section is executed when the migration is applied.
ALTER TABLE addresses ADD COLUMN global_code VARCHAR(255);

-- +goose Down
-- This section is executed when the migration is rolled back.
ALTER TABLE addresses DROP COLUMN global_code;