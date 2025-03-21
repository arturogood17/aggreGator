-- +goose Up
ALTER TABLE feeds
ADD COLUMN last_fetched_at TIMESTAMP;
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
ALTER TABLE feeds
DROP COLUMN last_fetched_at;
-- +goose StatementBegin
-- +goose StatementEnd
