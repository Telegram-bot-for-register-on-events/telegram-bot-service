-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    chat_id     BIGINT NOT NULL PRIMARY KEY,
    username    VARCHAR NOT NULL,
    created_at  TIMESTAMP NOT NULL
    );

-- +goose Down
DROP TABLE IF EXISTS users;

