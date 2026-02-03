-- +goose Up
-- +goose StatementBegin
INSERT INTO vehicles (type, vendor, model) VALUES
('Cargo', 'Volkswagen', 'Caddy'),
('Passenger', 'Toyota', 'Camry');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM vehicles WHERE model in ('Caddy', 'Camry');
-- +goose StatementEnd
