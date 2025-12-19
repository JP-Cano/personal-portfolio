-- +goose Up
-- +goose StatementBegin
-- Drop the old sessions table and recreate with correct foreign key
DROP TABLE IF EXISTS sessions;

CREATE TABLE IF NOT EXISTS sessions
(
    id         TEXT PRIMARY KEY,
    user_id    INTEGER  NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions (expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- This is a fix migration, no down needed
-- +goose StatementEnd
