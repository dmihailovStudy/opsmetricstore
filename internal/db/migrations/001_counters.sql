-- +goose Up
CREATE TABLE counters (
    timestamp TIMESTAMP,
    name VARCHAR(255),
    value BIGINT
);

-- +goose Down
DROP TABLE counters;