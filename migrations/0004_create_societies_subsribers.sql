-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS societies_subscribers
(
    id BIGINT PRIMARY KEY,
    society_id BIGINT NOT NULL FOREIGN KEY(societies.id),
    user_uuid BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS societies_subscribers;
-- +goose StatementEnd
