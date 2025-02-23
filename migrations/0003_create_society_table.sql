-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS society (
    id                      UUID PRIMARY KEY,
    name                    VARCHAR (255) NOT NULL,
    description             TEXT,
    owner_uuid              UUID NOT NULL,
    photo_url               TEXT DEFAULT 'https://storage.yandexcloud.net/space21/avatars/default/logo-discord.jpeg',
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
