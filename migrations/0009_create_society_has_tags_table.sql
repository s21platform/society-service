-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS access_level CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS society_has_tags (
    id              SERIAL PRIMARY KEY,
    society_id      INT NOT NULL,
    tag_id          INT,
    is_active       BOOL DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_society FOREIGN KEY (society_id) REFERENCES society (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS society_has_tags CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS access_level
(
    id          SERIAL PRIMARY KEY,
    level_name  VARCHAR(100),
    description TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

