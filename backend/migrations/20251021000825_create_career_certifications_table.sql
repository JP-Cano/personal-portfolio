-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS career_certifications
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    title          VARCHAR(255) NOT NULL,
    issuer         VARCHAR(255) NOT NULL,
    issue_date     DATETIME     NOT NULL,
    expiry_date    DATETIME,
    credential_id  VARCHAR(255),
    credential_url VARCHAR(500),
    file_url       VARCHAR(500) NOT NULL,
    file_name      VARCHAR(255) NOT NULL,
    original_name  VARCHAR(255) NOT NULL,
    file_size      INTEGER      NOT NULL,
    mime_type      VARCHAR(100) NOT NULL,
    description    TEXT,
    created_at     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at     DATETIME
);

CREATE INDEX idx_career_certifications_deleted_at ON career_certifications (deleted_at);
CREATE INDEX idx_career_certifications_issue_date ON career_certifications (issue_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_career_certifications_issue_date;
DROP INDEX IF EXISTS idx_career_certifications_deleted_at;
DROP TABLE IF EXISTS career_certifications;
-- +goose StatementEnd
