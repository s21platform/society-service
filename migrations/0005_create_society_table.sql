-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS societies;

-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS society
(
    id                      SERIAL PRIMARY KEY,
    name                    VARCHAR (255)       NOT NULL,
    description             TEXT,
    owner_uuid              UUID                NOT NULL,
    photo_url               TEXT,
    create_at               TIMESTAMP DEFAULT NOW(),
    update_at               TIMESTAMP DEFAULT NOW(),
    format_id               INT                 NOT NULL,
    post_permission_id      INT                 NOT NULL DEFAULT 1,
    is_search               BOOL DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS society;
-- +goose StatementEnd