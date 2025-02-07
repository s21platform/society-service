-- +goose Up
-- +goose StatementBegin
ALTER TABLE societies_subscribers ADD COLUMN id_user_permission INT NOT NULL REFERENCES user_permissions(id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS societies_subscribers;

-- +goose StatementEnd