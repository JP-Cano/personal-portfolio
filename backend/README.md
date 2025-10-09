# 🚀 Personal Portfolio API

A RESTful API built with Go, Gin, and SQLite for managing portfolio projects and work experiences.

## 📋 Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Database Migrations](#database-migrations)
- [Makefile Commands](#makefile-commands)
- [Database Schema](#database-schema)
- [Environment Variables](#environment-variables)
- [Development](#development)
- [Production](#production)
- [Troubleshooting](#troubleshooting)

## ✨ Features

- ✅ RESTful API with Gin framework
- ✅ SQLite database with GORM ORM
- ✅ Database migrations with Goose
- ✅ Automatic timestamps (created_at, updated_at)
- ✅ Work type enum validation (Remote, On Site, Hybrid)
- ✅ Comprehensive Makefile for development
- ✅ Clean architecture with separation of concerns

## 🛠️ Tech Stack

- **Language**: Go 1.25.1
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: SQLite
- **Migrations**: [Goose](https://github.com/pressly/goose)

## 📁 Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── models/                     # Domain models
│   │   ├── project.go
│   │   └── experience.go
│   ├── handlers/                   # HTTP handlers (controllers)
│   ├── services/                   # Business logic
│   ├── repository/                 # Database access layer
│   ├── routes/                     # Route definitions
│   └── middleware/                 # Custom middleware
├── pkg/
│   ├── database/                   # Database configuration
│   │   └── database.go
│   └── utils/                      # Utility functions
├── migrations/                     # Database migration files
│   ├── 20251009172324_create_projects_table.sql
│   └── 20251009172325_create_experiences_table.sql
├── Makefile                        # Development commands
├── README.md                       # This file
└── go.mod                          # Go dependencies

Generated files (ignored by git):
├── portfolio.db                    # SQLite database file
└── portfolio-api                   # Compiled binary
```

## 🚀 Getting Started

### Prerequisites

- Go 1.25.1 or higher
- Make (for using Makefile commands)
- SQLite3 (for database inspection)

### Installation

1. **Clone the repository**
   ```bash
   cd backend
   ```

2. **Install dependencies**
   ```bash
   make deps
   ```

3. **Install Goose (migration tool)**
   ```bash
   make install-tools
   ```

4. **Run migrations**
   ```bash
   make migrate-up
   ```
   This creates the `portfolio.db` file and sets up the database schema.

5. **Run the application**
   ```bash
   make run
   ```
   The API will start on `http://localhost:8080`

6. **Test it**
   ```bash
   curl http://localhost:8080/ping
   # Response: {"message":"pong"}
   ```

## 🗄️ Database Migrations

Migrations are managed with [Goose](https://github.com/pressly/goose) and stored in the `migrations/` directory.

### How Migrations Work

Migrations are version-controlled database schema changes. Each migration has:
- A **timestamp** (ensures correct execution order)
- An **Up** section (applies changes)
- A **Down** section (reverts changes)

### Creating a New Migration

```bash
make migrate-create NAME=add_skills_table
```

This generates a file like: `migrations/20251009180000_add_skills_table.sql`

### Edit the Migration File

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE skills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    level VARCHAR(50),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS skills;
-- +goose StatementEnd
```

### Deploying Migrations

**Option 1: Automatic (on app start)**
```bash
make run
```
The app automatically runs pending migrations on startup.

**Option 2: Manual**
```bash
make migrate-up
```

### Check Migration Status

```bash
make migrate-status
```

Output:
```
Applied At                  Migration
=======================================
Thu Oct  9 17:26:53 2025 -- 20251009172324_create_projects_table.sql
Thu Oct  9 17:26:53 2025 -- 20251009172325_create_experiences_table.sql
```

### Rollback Last Migration

```bash
make migrate-down
```

### Reset Database (Fresh Start)

```bash
make migrate-reset
```

## 📝 Makefile Commands

Run `make help` to see all available commands:

### Application Commands

| Command | Description |
|---------|-------------|
| `make run` | Run the application (auto-migrates database) |
| `make build` | Build the application binary |
| `make test` | Run tests |
| `make clean` | Clean build artifacts and database |
| `make dev` | Run with hot reload (requires air) |

### Migration Commands

| Command | Description |
|---------|-------------|
| `make migrate-up` | Run all pending migrations |
| `make migrate-down` | Rollback the last migration |
| `make migrate-status` | Show migration status |
| `make migrate-create NAME=xxx` | Create a new migration file |
| `make migrate-reset` | Rollback all migrations and re-run |
| `make db-setup` | Interactive database setup |

### Setup Commands

| Command | Description |
|---------|-------------|
| `make install-tools` | Install required tools (goose) |
| `make deps` | Download and tidy dependencies |
| `make help` | Show all available commands |

## 🗃️ Database Schema

### Projects Table

Stores portfolio projects.

```sql
CREATE TABLE projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    url VARCHAR(500),
    start_date DATE NOT NULL,
    end_date DATETIME,                                    -- NULL for ongoing projects
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Go Model:**
```go
type Project struct {
    ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
    Name        string     `json:"name" gorm:"type:varchar(255);not null"`
    Description string     `json:"description" gorm:"type:text"`
    URL         string     `json:"url" gorm:"type:varchar(500)"`
    StartDate   time.Time  `json:"start_date" gorm:"type:date;not null"`
    EndDate     *time.Time `json:"end_date,omitempty"`
    CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
```

### Experiences Table

Stores work experiences.

```sql
CREATE TABLE experiences (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    company VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    type VARCHAR(50) NOT NULL CHECK(type IN ('Remote', 'On Site', 'Hybrid')),
    start_date DATE NOT NULL,
    end_date DATETIME,                                    -- NULL for current position
    description TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Go Model:**
```go
type WorkType string

const (
    Remote WorkType = "Remote"
    OnSite WorkType = "On Site"
    Hybrid WorkType = "Hybrid"
)

type Experience struct {
    ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
    Title       string     `json:"title" gorm:"type:varchar(255);not null"`
    Company     string     `json:"company" gorm:"type:varchar(255);not null"`
    Location    string     `json:"location" gorm:"type:varchar(255)"`
    Type        WorkType   `json:"type" gorm:"type:varchar(50);not null"`
    StartDate   time.Time  `json:"start_date" gorm:"type:date;not null"`
    EndDate     *time.Time `json:"end_date,omitempty"`
    Description string     `json:"description" gorm:"type:text"`
    CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
```

### Automatic Timestamps

Both tables have triggers that automatically update `updated_at` on any UPDATE operation:

```sql
CREATE TRIGGER update_projects_updated_at
    AFTER UPDATE ON projects
    FOR EACH ROW
BEGIN
    UPDATE projects SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
```

## 🔧 Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_PATH` | `portfolio.db` | Path to SQLite database file |

**Example:**
```bash
DATABASE_PATH=/path/to/custom.db make run
```

## 💻 Development

### Running Locally

```bash
# Start the application
make run

# Or with hot reload (installs air if needed)
make dev
```

### Accessing the Database

**Option 1: SQLite CLI**
```bash
sqlite3 portfolio.db

.tables                    # List all tables
.schema experiences        # Show table structure
.schema projects
SELECT * FROM experiences; # Query data
.quit                      # Exit
```

**Option 2: DB Browser for SQLite**
1. Download from https://sqlitebrowser.org/
2. Open the app
3. Click "Open Database"
4. Select `portfolio.db` in the backend directory

**Option 3: DBeaver**
1. Install DBeaver
2. New Database Connection → SQLite
3. Browse to `portfolio.db`

### Creating a New Migration

```bash
# 1. Create migration file
make migrate-create NAME=add_categories_table

# 2. Edit the generated file in migrations/
vim migrations/TIMESTAMP_add_categories_table.sql

# 3. Apply the migration
make migrate-up

# 4. Verify
make migrate-status
```

### Testing Migrations

```bash
# Test up migration
make migrate-up

# Test down migration (rollback)
make migrate-down

# Re-apply
make migrate-up
```

## 🚢 Production

### Building for Production

```bash
# Build binary
make build

# This creates: portfolio-api
```

### Running in Production

```bash
# Option 1: Run migrations separately (recommended)
make migrate-up
./portfolio-api

# Option 2: Let app auto-migrate
./portfolio-api
```

### Production Checklist

- [ ] Set `GIN_MODE=release` environment variable
- [ ] Use a persistent volume for `portfolio.db`
- [ ] Back up database regularly
- [ ] Run migrations in a separate step (not on app startup)
- [ ] Use proper error logging
- [ ] Set up monitoring

## 🐛 Troubleshooting

### Migration Fails

```bash
# Check what went wrong
make migrate-status

# Rollback if needed
make migrate-down

# Fix the migration file, then:
make migrate-up
```

### Database Locked

This happens when another process is using the database.

```bash
# Check if app is running
lsof -ti:8080

# Kill the process if needed
lsof -ti:8080 | xargs kill -9
```

### Start Fresh

```bash
# Delete everything and start over
make clean
make migrate-up
make run
```

### Port Already in Use

```bash
# Kill process using port 8080
lsof -ti:8080 | xargs kill -9

# Then run again
make run
```

### Goose Not Found

```bash
# Install goose
make install-tools

# Or manually:
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## 📚 Migration Best Practices

### ✅ DO:

1. **Always write Down migrations**
   ```sql
   -- +goose Down
   DROP TABLE IF EXISTS skills;
   ```

2. **One logical change per migration**
   - ✅ Good: `add_skills_table.sql`
   - ❌ Bad: `add_skills_and_categories_and_fix_projects.sql`

3. **Test both Up and Down**
   ```bash
   make migrate-up
   make migrate-down
   make migrate-up
   ```

4. **Use descriptive names**
   - ✅ `add_featured_flag_to_projects`
   - ❌ `update_projects`

### ❌ DON'T:

1. **Never edit applied migrations**
   - Create a new migration to fix issues
   
2. **Don't forget Down migrations**
   - Always make migrations reversible

3. **Don't use database-specific syntax**
   - Keep it SQLite-compatible

## 📖 Additional Resources

- [Gin Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [Goose Documentation](https://github.com/pressly/goose)
- [SQLite Documentation](https://www.sqlite.org/docs.html)

## 🤝 Contributing

1. Create a new migration for schema changes
2. Test both up and down migrations
3. Ensure code follows Go best practices
4. Write tests for new features

## 📄 License

This project is part of a personal portfolio.

---

## 🎯 Quick Reference

```bash
# First time setup
make deps
make install-tools
make migrate-up

# Daily development
make run                                    # Start app
make migrate-create NAME=add_feature       # New migration
make migrate-up                             # Apply migrations
make migrate-status                         # Check status

# Troubleshooting
make clean                                  # Start fresh
make migrate-reset                          # Reset database
make help                                   # See all commands
```

**Happy coding! 🚀**
