package database

import (
	"log"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(dbPath string) error {
	var err error

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	log.Println("Database connection established")

	err = AutoMigrate()
	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.Project{},
		&models.Experience{},
	)
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
