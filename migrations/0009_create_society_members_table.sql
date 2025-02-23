-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS society_members (
    id                      SERIAL PRIMARY KEY,
    society_id              UUID NOT NULL,
    user_uuid               UUID NOT NULL,
    role                    INT NOT NULL,
    create_at               TIMESTAMP DEFAULT NOW(),
    expires_at              TIMESTAMP DEFAULT NULL,
    payment_status          INT NOT NULL DEFAULT 1,
    CONSTRAINT fk_society FOREIGN KEY (society_id) REFERENCES society (id),
    CONSTRAINT fk_role FOREIGN KEY (role) REFERENCES role_members (id),
    CONSTRAINT fk_payment_members FOREIGN KEY (payment_status) REFERENCES payment_members (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS society_members CASCADE;
-- +goose StatementEnd
