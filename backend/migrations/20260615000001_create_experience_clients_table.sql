-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS experience_clients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    experience_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(500),
    start_date DATE NOT NULL,
    end_date DATETIME,
    description TEXT,
    achievements TEXT DEFAULT '[]',
    responsibilities TEXT DEFAULT '[]',
    technologies TEXT DEFAULT '[]',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (experience_id) REFERENCES experiences(id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS update_experience_clients_updated_at
    AFTER UPDATE ON experience_clients
    FOR EACH ROW
BEGIN
    UPDATE experience_clients SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

CREATE INDEX IF NOT EXISTS idx_experience_clients_experience_id ON experience_clients(experience_id);
CREATE INDEX IF NOT EXISTS idx_experience_clients_deleted_at ON experience_clients(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_experience_clients_updated_at;
DROP INDEX IF EXISTS idx_experience_clients_deleted_at;
DROP INDEX IF EXISTS idx_experience_clients_experience_id;
DROP TABLE IF EXISTS experience_clients;
-- +goose StatementEnd
