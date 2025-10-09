-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    url VARCHAR(500),
    start_date DATE NOT NULL,
    end_date DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger to update updated_at timestamp
CREATE TRIGGER IF NOT EXISTS update_projects_updated_at
    AFTER UPDATE ON projects
    FOR EACH ROW
BEGIN
    UPDATE projects SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_projects_updated_at;
DROP TABLE IF EXISTS projects;
-- +goose StatementEnd
