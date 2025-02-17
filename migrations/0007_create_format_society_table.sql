-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS societies_subscribers CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS format_society (
    id              SERIAL PRIMARY KEY,
    format_name     TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS format_society CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS societies_subscribers
(
    id SERIAL PRIMARY KEY,
    society_id BIGINT NOT NULL,
    user_uuid UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (society_id, user_uuid),
    CONSTRAINT fk_society FOREIGN KEY (society_id) REFERENCES societies (id)
);
-- +goose StatementEnd

