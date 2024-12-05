-- +goose Up
CREATE TABLE IF NOT EXISTS users_info (
    id UUID PRIMARY KEY,
    ip VARCHAR(50) NOT NULL,
    hashed_refresh_token BYTEA NOT NULL,
);

-- +goose Down
DROP TABLE IF EXISTS users_info;
