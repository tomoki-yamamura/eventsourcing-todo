-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    aggregate_id CHAR(36) NOT NULL,
    event_id CHAR(36) NOT NULL UNIQUE,
    event_type VARCHAR(255) NOT NULL,
    event_data JSON NOT NULL,
    version INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE INDEX unique_aggregate_version (aggregate_id, version)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd