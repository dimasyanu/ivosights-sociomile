-- +goose Up
-- +goose StatementBegin
INSERT INTO users
(id, `name`, email, password_hash, roles, created_at, created_by, updated_at, updated_by, deleted_at)
VALUES (
    0x6AB04FE7D763470F82D562F53C157723,
    'Administrator',
    'admin@mail.com',
    '$2a$10$pvybJ1KmPcZkaM8TPD/51O9gIVt7Jazqrro1BsPHeZKn5qwnNY1A.',
    'admin',
    NOW(),
    'system',
    NOW(),
    'system',
    NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE email = 'admin@mail.com';
-- +goose StatementEnd
