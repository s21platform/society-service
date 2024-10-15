-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS societies
(
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(255)            NOT NULL,
    description  TEXT                    NOT NULL,
    is_private   BOOLEAN   DEFAULT FALSE NOT NULL,
    direction_id INT                     NOT NULL,
    owner_id     INT                     NOT NULL,
    photo_url    TEXT,
    access_id    INT                     NOT NULL REFERENCES access_level (id),
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS societies;
-- +goose StatementEnd
