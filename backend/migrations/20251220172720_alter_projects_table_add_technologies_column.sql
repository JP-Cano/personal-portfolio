-- +goose Up
-- +goose StatementBegin
ALTER TABLE projects ADD COLUMN technologies TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE projects DROP COLUMN technologies;
-- +goose StatementEnd
