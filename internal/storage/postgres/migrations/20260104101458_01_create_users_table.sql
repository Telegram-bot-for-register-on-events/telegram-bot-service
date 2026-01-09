-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    chat_id     BIGINT PRIMARY KEY,
    username    VARCHAR,
    created_at  TIMESTAMP
    );

-- +goose Down
DROP TABLE IF EXISTS users;

