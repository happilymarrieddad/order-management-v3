-- +goose Up
-- +goose StatementBegin
CREATE TABLE company_attributes (
    id BIGSERIAL PRIMARY KEY,
    position INT NOT NULL,
    commodity_attribute_id BIGINT NOT NULL,
    company_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_ca_commodity_attribute FOREIGN KEY (commodity_attribute_id) REFERENCES commodity_attributes(id) ON DELETE CASCADE,
    CONSTRAINT fk_ca_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT uq_ca_commodity_attribute_company UNIQUE (commodity_attribute_id, company_id),
    CONSTRAINT uq_ca_company_position UNIQUE (company_id, position)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_company_attribute_position()
RETURNS TRIGGER AS $$
BEGIN
    -- Lock the table to prevent race conditions when calculating the next position.
    -- This ensures that concurrent inserts for the same company are serialized.
    LOCK TABLE company_attributes IN EXCLUSIVE MODE;
    -- Find the current max position for the given company_id and add 1.
    -- COALESCE handles the case where it's the first attribute for the company (MAX would be NULL).
    SELECT COALESCE(MAX(position), 0) + 1 INTO NEW.position
    FROM company_attributes
    WHERE company_id = NEW.company_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER set_position_before_insert
BEFORE INSERT ON company_attributes
FOR EACH ROW
EXECUTE FUNCTION set_company_attribute_position();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_position_before_insert ON company_attributes;
-- +goose StatementEnd

-- +goose StatementBegin
DROP FUNCTION IF EXISTS set_company_attribute_position();
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS company_attributes;
-- +goose StatementEnd