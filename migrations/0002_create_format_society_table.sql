-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS format_society (
    id              SERIAL PRIMARY KEY,
    format_name     TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO format_society (format_name) VALUES
                                                  ('open'),
                                                  ('close'),
                                                  ('paid');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS format_society CASCADE;
-- +goose StatementEnd
