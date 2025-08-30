-- +goose Up
ALTER TABLE companies ADD COLUMN order_prefix TEXT NOT NULL DEFAULT '';
ALTER TABLE companies ADD COLUMN order_postfix TEXT NOT NULL DEFAULT '';
ALTER TABLE companies ADD COLUMN default_order_number INTEGER NOT NULL DEFAULT 100000;

-- +goose Down
ALTER TABLE companies DROP COLUMN order_prefix;
ALTER TABLE companies DROP COLUMN order_postfix;
ALTER TABLE companies DROP COLUMN default_order_number;