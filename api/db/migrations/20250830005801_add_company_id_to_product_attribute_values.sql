-- +goose Up
-- +goose StatementBegin
ALTER TABLE product_attribute_values ADD COLUMN company_id BIGINT NOT NULL;
ALTER TABLE product_attribute_values ADD CONSTRAINT fk_pav_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE product_attribute_values DROP CONSTRAINT fk_pav_company;
ALTER TABLE product_attribute_values DROP COLUMN company_id;
-- +goose StatementEnd
