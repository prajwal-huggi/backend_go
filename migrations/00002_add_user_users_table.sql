-- +goose Up
INSERT INTO users (name, email, role)
VALUES
('Alice Johnson', 'alice@example.com', 'admin'),
('Bob Smith', 'bob@example.com', 'manager'),
('Charlie Brown', 'charlie@example.com', 'user'),
('Diana Prince', 'diana@example.com', 'user'),
('Ethan Clark', 'ethan@example.com', 'manager');

-- +goose Down
DELETE FROM users
WHERE id IN (1,2,3,4,5);
