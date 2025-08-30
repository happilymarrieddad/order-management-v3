-- +goose Up
-- +goose StatementBegin
ALTER TABLE commodities ALTER COLUMN commodity_type TYPE INT USING
    CASE commodity_type
        WHEN 'Unknown' THEN 0
        WHEN 'Produce' THEN 1
        ELSE 0 -- Default to Unknown if unexpected value
    END;
ALTER TABLE commodity_attributes ALTER COLUMN commodity_type_name TYPE INT USING
    CASE commodity_type_name
        WHEN 'Unknown' THEN 0
        WHEN 'Produce' THEN 1
        ELSE 0 -- Default to Unknown if unexpected value
    END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE commodities ALTER COLUMN commodity_type TYPE VARCHAR(255) USING
    CASE commodity_type
        WHEN 0 THEN 'Unknown'
        WHEN 1 THEN 'Produce'
        ELSE 'Unknown' -- Default to Unknown if unexpected value
    END;
ALTER TABLE commodity_attributes ALTER COLUMN commodity_type_name TYPE VARCHAR(255) USING
    CASE commodity_type_name
        WHEN 0 THEN 'Unknown'
        WHEN 1 THEN 'Produce'
        ELSE 'Unknown' -- Default to Unknown if unexpected value
    END;
-- +goose StatementEnd
