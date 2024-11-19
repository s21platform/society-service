-- +goose Up
-- +goose StatementBegin
ALTER TABLE societies DROP COLUMN owner_id;
ALTER TABLE societies ADD COLUMN owner_uuid UUID;



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE societies DROP COLUMN ownet_uuid;
ALTER TABLE societies ADD COLUMN owner_id INT;
-- +goose StatementEnd
