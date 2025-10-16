-- +goose Up
-- +goose StatementBegin
ALTER TABLE experiences ADD COLUMN deleted_at DATETIME;

-- Create index on deleted_at for better query performance with soft deletes
CREATE INDEX idx_experiences_deleted_at ON experiences(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_experiences_deleted_at;

-- SQLite doesn't support DROP COLUMN directly in older versions
-- If you need to rollback, you would need to recreate the table
-- For now, we'll just note this limitation
-- ALTER TABLE experiences DROP COLUMN deleted_at;
-- +goose StatementEnd
