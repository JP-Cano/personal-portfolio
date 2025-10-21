package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Config holds database configuration
type Config struct {
	Driver        string // "sqlite" or "turso"
	SQLitePath    string // Path for SQLite database
	TursoURL      string // Turso database URL
	TursoToken    string // Turso auth token
	MigrationsDir string
}

// InitDB initializes the database connection based on the driver type
func InitDB(config Config) error {
	var err error

	switch config.Driver {
	case "turso":
		log.Printf("Connecting to Turso database...")
		// For Turso, construct the libsql:// DSN
		dsn := fmt.Sprintf("%s?authToken=%s", config.TursoURL, config.TursoToken)

		// Open with libsql driver
		sqlDB, err := sql.Open("libsql", dsn)
		if err != nil {
			return fmt.Errorf("failed to open turso database: %w", err)
		}

		// Test the connection
		if err := sqlDB.Ping(); err != nil {
			return fmt.Errorf("failed to ping turso database: %w", err)
		}

		// Use GORM with the sql.DB instance
		DB, err = gorm.Open(sqlite.Dialector{
			Conn: sqlDB,
		}, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return fmt.Errorf("failed to connect to turso database with gorm: %w", err)
		}

	case "sqlite":
		log.Printf("Connecting to SQLite database: %s", config.SQLitePath)
		DB, err = gorm.Open(sqlite.Open(config.SQLitePath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return fmt.Errorf("failed to connect to sqlite database: %w", err)
		}

	default:
		return fmt.Errorf("unsupported database driver: %s (use 'sqlite' or 'turso')", config.Driver)
	}

	log.Println("Database connection established")
	return nil
}

// RunMigrations runs database migrations using goose
func RunMigrations(config Config) error {
	var db *sql.DB
	var err error

	switch config.Driver {
	case "turso":
		log.Println("Running migrations on Turso database...")
		// For Turso, construct the libsql:// DSN
		dsn := fmt.Sprintf("%s?authToken=%s", config.TursoURL, config.TursoToken)
		db, err = sql.Open("libsql", dsn)
		if err != nil {
			return fmt.Errorf("failed to open turso database for migrations: %w", err)
		}

	case "sqlite":
		log.Println("Running migrations on SQLite database...")
		db, err = sql.Open("sqlite3", config.SQLitePath)
		if err != nil {
			return fmt.Errorf("failed to open sqlite database for migrations: %w", err)
		}

	default:
		return fmt.Errorf("unsupported database driver: %s", config.Driver)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("failed to close database connection: %v", err)
		}
	}(db)

	// Both SQLite and Turso use sqlite3 dialect
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, config.MigrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
