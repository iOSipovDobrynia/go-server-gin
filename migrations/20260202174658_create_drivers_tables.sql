-- +goose Up
-- +goose StatementBegin
CREATE TABLE drivers (
    id SERIAL PRIMARY KEY ,
    name VARCHAR(100),
    score INTEGER,
    vehicle_id INTEGER REFERENCES vehicles(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE drivers;
-- +goose StatementEnd
