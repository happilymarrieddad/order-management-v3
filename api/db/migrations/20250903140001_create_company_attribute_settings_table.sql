-- +goose Up
CREATE TABLE company_attribute_settings (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL,
    commodity_attribute_id BIGINT NOT NULL,
    display_order INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_cas_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT fk_cas_commodity_attribute FOREIGN KEY (commodity_attribute_id) REFERENCES commodity_attributes(id) ON DELETE CASCADE,
    CONSTRAINT uq_cas_company_attribute UNIQUE (company_id, commodity_attribute_id),
    CONSTRAINT uq_cas_company_display_order UNIQUE (company_id, display_order)
);

-- +goose Down
DROP TABLE IF EXISTS company_attribute_settings;