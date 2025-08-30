-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status_enum AS ENUM (
    'pending_acceptance',
    'pending_booking',
    'hold',
    'booked',
    'shipped_in_transit',
    'delivered',
    'ready_to_invoice',
    'invoiced',
    'rejected',
    'cancelled',
    'hold_for_pod',
    'order_template',
    'paid_in_full'
);
-- +goose StatementEnd

-- +goose Down
DROP TYPE IF EXISTS order_status_enum;