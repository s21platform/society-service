-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS access_level
(
    id          SERIAL PRIMARY KEY,
    level_name  VARCHAR(100),
    description TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS access_level;
-- +goose StatementEnd
