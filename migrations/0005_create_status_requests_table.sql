-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS status_requests (
    id      SERIAL PRIMARY KEY,
    status  TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO status_requests (status) VALUES
                                         ('pending'),
                                         ('approved'),
                                         ('rejected');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS status_requests CASCADE;
-- +goose StatementEnd
