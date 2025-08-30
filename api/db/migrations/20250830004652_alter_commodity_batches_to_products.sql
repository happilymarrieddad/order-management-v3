-- +goose Up
-- +goose StatementBegin
ALTER TABLE commodity_batches RENAME TO products;

ALTER TABLE products DROP CONSTRAINT fk_commodity_batches_commodity;
ALTER TABLE products DROP COLUMN commodity_id;

ALTER TABLE products ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE products ADD COLUMN commodity_id BIGINT NOT NULL;
ALTER TABLE products ADD CONSTRAINT fk_products_commodity FOREIGN KEY (commodity_id) REFERENCES commodities(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products DROP CONSTRAINT fk_products_commodity;
ALTER TABLE products DROP COLUMN commodity_id;
ALTER TABLE products DROP COLUMN name;

ALTER TABLE products ADD COLUMN commodity_id BIGINT NOT NULL;
ALTER TABLE products ADD CONSTRAINT fk_commodity_batches_commodity FOREIGN KEY (commodity_id) REFERENCES commodities(id) ON DELETE CASCADE;

ALTER TABLE products RENAME TO commodity_batches;
-- +goose StatementEnd
