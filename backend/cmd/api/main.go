package main

import (
	"os"

	_ "github.com/JuanPabloCano/personal-portfolio/backend/docs"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/handlers"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/middleware"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/repository"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/routes"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/services"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/constants"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/database"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"gorm.io/gorm"
)

// @title Personal Portfolio API
// @version 1.0
// @description API for managing personal portfolio experiences and projects

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
func main() {
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using environment variables")
	}
	logLevel := logger.INFO
	if os.Getenv("DEBUG") == "true" {
		logLevel = logger.DEBUG
	}

	logger.Init(logger.Config{
		Level:      logLevel,
		Output:     os.Stdout,
		UseColor:   true,
		IncludePos: true,
	})

	logger.Info("Starting Personal Portfolio API...")

	dbConfig := configDatabaseDriver()

	err := database.InitDB(dbConfig)
	if err != nil {
		logger.Fatal("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := database.CloseDB(); err != nil {
			logger.Error("Error closing database: %v", err)
		}
	}()

	logger.Info("Database initialized successfully")

	err = database.RunMigrations(dbConfig)
	if err != nil {
		logger.Fatal("Failed to run migrations: %v", err)
	}

	logger.Info("Database migrations completed")

	db := database.GetDB()

	projectHandler, experienceHandler, careerCertificationHandler := registerDependencies(db)

	// Set Gin to release mode if not in debug
	if os.Getenv("DEBUG") != "true" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.MaxMultipartMemory = 10 << 20 // 10 MB

	// Add custom middleware
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.LoggerMiddleware())

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	r.Use(cors.New(config))

	r.Static("/certifications", constants.CareerCertificationsDir)
	logger.Info("Serving static files from: %s", constants.CareerCertificationsDir)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.SetupRoutes(
		r,
		experienceHandler,
		projectHandler,
		careerCertificationHandler,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Starting server on port: %s", port)
	logger.Info("Swagger documentation available at: http://localhost:%s/swagger/index.html", port)

	err = r.Run(":" + port)
	if err != nil {
		logger.Fatal("Failed to start server: %v", err)
	}
}

func configDatabaseDriver() database.Config {
	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "sqlite"
	}

	dbConfig := database.Config{
		Driver:        dbDriver,
		MigrationsDir: "migrations",
	}

	switch dbDriver {
	case "turso":
		tursoURL := os.Getenv("TURSO_DATABASE_URL")
		tursoToken := os.Getenv("TURSO_AUTH_TOKEN")

		if tursoURL == "" || tursoToken == "" {
			logger.Fatal("TURSO_DATABASE_URL and TURSO_AUTH_TOKEN must be set when using turso driver")
		}

		dbConfig.TursoURL = tursoURL
		dbConfig.TursoToken = tursoToken
		logger.Info("Using Turso database (production)")

	case "sqlite":
		dbPath := os.Getenv("DATABASE_PATH")
		if dbPath == "" {
			dbPath = "portfolio.db"
		}

		dbConfig.SQLitePath = dbPath
		logger.Info("Using SQLite database (development): %s", dbPath)

	default:
		logger.Fatal("Invalid DB_DRIVER: %s (use 'sqlite' or 'turso')", dbDriver)
	}

	return dbConfig
}

// registerDependencies initializes and registers all necessary dependencies for project and experience handlers.
// It returns the initialized ProjectHandler and ExperienceHandler instances.
func registerDependencies(db *gorm.DB) (*handlers.ProjectHandler, *handlers.ExperienceHandler, *handlers.CareerCertificationHandler) {
	experienceRepo := repository.NewExperienceRepository(db)
	experienceService := services.NewExperienceService(experienceRepo)
	experienceHandler := handlers.NewExperienceHandler(experienceService)

	projectRepo := repository.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepo)
	projectHandler := handlers.NewProjectHandler(projectService)

	careerCertificationRepo := repository.NewCareerCertificationRepository(db)
	careerCertificationService := services.NewCareerCertificationService(careerCertificationRepo)
	careerCertificationHandler := handlers.NewCareerCertificationHandler(careerCertificationService)

	return projectHandler, experienceHandler, careerCertificationHandler

}
