-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN deleted_at DATETIME;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
