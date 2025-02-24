-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS society_has_tags (
    id              SERIAL PRIMARY KEY,
    society_id      UUID NOT NULL,
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
