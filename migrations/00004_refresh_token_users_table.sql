-- +goose Up
ALTER TABLE users ADD COLUMN refresh_token TEXT;
ALTER TABLE users ADD COLUMN refresh_token_expiry TIMESTAMP;

-- +goose Down
ALTER TABLE users DROP COLUMN refresh_token;
ALTER TABLE users DROP COLUMN refresh_token_expiry;