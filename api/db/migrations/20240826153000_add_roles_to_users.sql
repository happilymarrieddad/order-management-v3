-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN roles text[] NOT NULL DEFAULT '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN roles;
-- +goose StatementEnd