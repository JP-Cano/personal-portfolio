package main

import (
	"log"
	"net/http"
	"os"

	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	databasePath := os.Getenv("DATABASE_PATH")
	err := database.InitDB(databasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		err := database.CloseDB()
		if err != nil {
			log.Fatalf("Failed to close database: %v", err)
		}
	}()

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
