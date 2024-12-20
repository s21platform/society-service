-- +goose Up
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

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS societies_subscribers;
-- +goose StatementEnd
