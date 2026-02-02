-- +goose Up
-- +goose StatementBegin
CREATE TABLE vehicles (
    id SERIAL PRIMARY KEY,
    type VARCHAR(100),
    vendor VARCHAR(100),
    model VARCHAR(100)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE vehicles;
-- +goose StatementEnd
