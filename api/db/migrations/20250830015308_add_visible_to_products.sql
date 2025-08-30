-- +goose Up
-- +goose StatementBegin
ALTER TABLE products ADD COLUMN visible BOOLEAN NOT NULL DEFAULT TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products DROP COLUMN visible;
-- +goose StatementEnd
