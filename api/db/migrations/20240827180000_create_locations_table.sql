-- +goose Up
-- +goose StatementBegin
CREATE TABLE locations (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL,
    address_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_company
        FOREIGN KEY(company_id)
        REFERENCES companies(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_address
        FOREIGN KEY(address_id)
        REFERENCES addresses(id)
        ON DELETE RESTRICT,

    UNIQUE (company_id, name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS locations;
-- +goose StatementEnd