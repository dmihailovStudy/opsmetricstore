-- +goose Up
CREATE TABLE gauges (
    timestamp TIMESTAMP,
    name VARCHAR(255),
    value DOUBLE PRECISION
);

-- +goose Down
DROP TABLE gauges;