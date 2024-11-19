-- +goose Up
-- +goose StatementBegin
ALTER TABLE societies DROP COLUMN owner_id;
ALTER TABLE societies ADD COLUMN owner_uuid UUID;



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE societies ADD COLUMN owner_uuid UUID;
ALTER TABLE societies DROP COLUMN owner_id;
-- +goose StatementEnd
