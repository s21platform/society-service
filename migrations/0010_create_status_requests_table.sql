-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS status_requests (
    id      SERIAL PRIMARY KEY,
    status  TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS status_requests CASCADE;
-- +goose StatementEnd
