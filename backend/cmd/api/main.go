package main

import (
	"log"
	"os"

	_ "github.com/JuanPabloCano/personal-portfolio/backend/docs"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/handlers"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/repository"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/routes"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/services"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/database"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Personal Portfolio API
// @version 1.0
// @description API for managing personal portfolio experiences and projects
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

func main() {
	// Database configuration
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

	db := database.GetDB()

	experienceRepo := repository.NewExperienceRepository(db)
	experienceService := services.NewExperienceService(experienceRepo)
	experienceHandler := handlers.NewExperienceHandler(experienceService)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.SetupRoutes(r, experienceHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	log.Printf("Swagger documentation available at: http://localhost:%s/swagger/index.html", port)
	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
