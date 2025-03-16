-- +goose Up
CREATE TABLE users(id INTEGER PRIMARY KEY,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
name TEXT NOT NULL UNIQUE)
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE users;
-- +goose StatementBegin
-- +goose StatementEnd
