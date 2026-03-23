-- +goose Up
ALTER TABLE users
ADD COLUMN password TEXT DEFAULT '';

UPDATE users
SET password= ''
WHERE password is NULL;
-- +goose Down
ALTER TABLE users
DROP COLUMN password;
