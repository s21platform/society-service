-- +goose Up
CREATE TABLE IF NOT EXISTS post_permission (
    id              SERIAL PRIMARY KEY,
    post_permission TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO post_permission (post_permission) VALUES
    ('Owner/Admin/Moderator, comment OFF'),
    ('All, comment OFF'),
    ('Owner/Admin/Moderator, comment ON'),
    ('All, comment ON');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS post_permission CASCADE;
-- +goose StatementEnd
