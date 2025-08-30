-- +goose Up
-- +goose StatementBegin
ALTER TABLE products DROP COLUMN quantity;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products ADD COLUMN quantity DOUBLE PRECISION;
-- +goose StatementEnd
