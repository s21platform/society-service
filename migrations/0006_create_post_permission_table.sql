-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS user_permissions CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS post_permission (
    id              SERIAL PRIMARY KEY,
    post_permission TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS post_permission CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_permissions
(
    id SERIAL PRIMARY KEY,
    name TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

