-- +goose Up
-- +goose StatementBegin
CREATE TABLE commodity_batches (
    id BIGSERIAL PRIMARY KEY,
    commodity_id BIGINT NOT NULL,
    company_id BIGINT NOT NULL,
    quantity DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_commodity_batches_commodity FOREIGN KEY (commodity_id) REFERENCES commodities(id) ON DELETE CASCADE,
    CONSTRAINT fk_commodity_batches_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS commodity_batches;
-- +goose StatementEnd
