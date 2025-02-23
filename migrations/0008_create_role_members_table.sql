-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS role_members (
    id      SERIAL PRIMARY KEY,
    status  TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO role_members (status) VALUES
                                      ('owner'),
                                      ('admin'),
                                      ('moderator'),
                                      ('member');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS role_members CASCADE;
-- +goose StatementEnd
