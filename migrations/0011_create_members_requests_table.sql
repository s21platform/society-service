-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS members_requests (
    id      SERIAL PRIMARY KEY,
    user_uuid  UUID NOT NULL,
    society_id INT NOT NULL,
    status_id INT NOT NULL,
    create_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_society FOREIGN KEY (society_id) REFERENCES society (id),
    CONSTRAINT fk_status_requests FOREIGN KEY (status_id) REFERENCES status_requests (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS members_requests CASCADE;
-- +goose StatementEnd
