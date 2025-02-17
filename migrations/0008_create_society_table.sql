-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS societies CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS society (
    id                      SERIAL PRIMARY KEY,
    name                    VARCHAR (255) NOT NULL,
    description             TEXT,
    owner_uuid              UUID NOT NULL,
    photo_url               TEXT,
    create_at               TIMESTAMP DEFAULT NOW(),
    update_at               TIMESTAMP DEFAULT NOW(),
    format_id               INT NOT NULL,
    post_permission_id      INT NOT NULL DEFAULT 1,
    is_search               BOOL DEFAULT FALSE,
    CONSTRAINT fk_format_society FOREIGN KEY (format_id) REFERENCES format_society (id),
    CONSTRAINT fk_post_permission FOREIGN KEY (post_permission_id) REFERENCES post_permission (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS society CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS societies
(
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(255)            NOT NULL,
    description  TEXT                    NOT NULL,
    is_private   BOOLEAN   DEFAULT FALSE NOT NULL,
    direction_id INT                     NOT NULL,
    owner_uuid   UUID                    NOT NULL,
    photo_url    TEXT,
    access_id    INT                     NOT NULL REFERENCES access_level (id),
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
-- +goose StatementEnd

