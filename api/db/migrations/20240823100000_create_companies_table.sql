-- +goose Up
CREATE TABLE companies (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    address_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_companies_address
        FOREIGN KEY(address_id)
        REFERENCES addresses(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_companies_address_id ON companies(address_id);

-- +goose Down
DROP TABLE IF EXISTS companies;