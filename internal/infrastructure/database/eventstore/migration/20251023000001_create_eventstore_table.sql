-- +goose Up
-- +goose StatementBegin
CREATE TABLE eventstore (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    aggregate_id CHAR(36) NOT NULL,
    event_id CHAR(36) NOT NULL UNIQUE,
    event_type VARCHAR(100) NOT NULL,
    event_data JSON NOT NULL,
    version INT NOT NULL,
    occurred_at TIMESTAMP(6) NOT NULL,
    created_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6),
    INDEX idx_aggregate_id_version (aggregate_id, version),
    INDEX idx_event_id (event_id),
    INDEX idx_occurred_at (occurred_at)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE eventstore;
-- +goose StatementEnd