-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_permissions
(
    id BIGINT PRIMARY KEY,
    name TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_permissions;
-- +goose StatementEnd
