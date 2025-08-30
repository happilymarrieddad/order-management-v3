-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_attribute_values (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    commodity_attribute_id BIGINT NOT NULL,
    value VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_pav_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT fk_pav_commodity_attribute FOREIGN KEY (commodity_attribute_id) REFERENCES commodity_attributes(id) ON DELETE CASCADE,
    CONSTRAINT uq_pav_product_attribute UNIQUE (product_id, commodity_attribute_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_attribute_values;
-- +goose StatementEnd
