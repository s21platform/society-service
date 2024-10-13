-- +goose Up
-- +goose StatementBegin
    CREATE TABLE IF NOT EXISTS societies (
        id PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        description TEXT,
        is_private BOOLEAN DEFAULT FALSE,
        direction_id INT,
        owner_id INT,
        photo_url TEXT,
        access_level VARCHAR(50) NOT NULL references cities(name),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS societies;
-- +goose StatementEnd
