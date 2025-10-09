-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS experiences (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    company VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    type VARCHAR(50) NOT NULL CHECK(type IN ('Remote', 'On Site', 'Hybrid')),
    start_date DATE NOT NULL,
    end_date DATETIME,
    description TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger to update updated_at timestamp
CREATE TRIGGER IF NOT EXISTS update_experiences_updated_at
    AFTER UPDATE ON experiences
    FOR EACH ROW
BEGIN
    UPDATE experiences SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_experiences_updated_at;
DROP TABLE IF EXISTS experiences;
-- +goose StatementEnd
