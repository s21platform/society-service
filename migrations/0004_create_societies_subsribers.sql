-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS societies_subscribers
(
    id BIGINT PRIMARY KEY,
    society_id BIGINT NOT NULL,
    user_uuid UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_society FOREIGN KEY (society_id) REFERENCES societies (id)
<<<<<<< HEAD
    );
=======
);
>>>>>>> main
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS societies_subscribers;
<<<<<<< HEAD
-- +goose StatementEnd
=======
-- +goose StatementEnd
>>>>>>> main
