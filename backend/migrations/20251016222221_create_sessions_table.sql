-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sessions
(
    id         TEXT PRIMARY KEY,
    user_id    INTEGER  NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user (id) ON DELETE CASCADE
);

-- Create index on deleted_at for better query performance with soft deletes
CREATE INDEX idx_sessions_expires_at ON sessions (expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd
