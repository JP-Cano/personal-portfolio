-- +goose Up
-- +goose StatementBegin
ALTER TABLE experiences
    ADD COLUMN url VARCHAR(500);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE experiences DROP COLUMN url;
-- +goose StatementEnd
