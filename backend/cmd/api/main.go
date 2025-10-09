package main

import (
	"log"
	"net/http"
	"os"

	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "portfolio.db"
	}

	migrationsDir := "migrations"

	err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := database.CloseDB(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	err = database.RunMigrations(dbPath, migrationsDir)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	log.Println("Starting server on :8080")
	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
