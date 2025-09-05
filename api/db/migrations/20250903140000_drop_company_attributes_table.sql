-- +goose Up
DROP TABLE IF EXISTS company_attributes;

-- +goose Down
CREATE TABLE company_attributes (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL,
    key VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_ca_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT uq_ca_company_key UNIQUE (company_id, key)
);