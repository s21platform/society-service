-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payment_members (
    id      SERIAL PRIMARY KEY,
    status  TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment_members CASCADE;
-- +goose StatementEnd
