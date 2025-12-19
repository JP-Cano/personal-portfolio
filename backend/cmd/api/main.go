package main

import (
	"os"
	"strings"
	"time"

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
	if err := godotenv.Load("../.env"); err != nil {
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

	projectHandler, experienceHandler, careerCertificationHandler, authHandler, authService := registerDependencies(db)

	cleanExpiredSessions(authService)

	// Set Gin to release mode if not in debug
	if os.Getenv("DEBUG") != "true" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.MaxMultipartMemory = 10 << 20 // 10 MB

	// Add custom middleware
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.Throttler(middleware.MaxRequestsPerSecond))

	// CORS configuration
	origins := os.Getenv("ALLOWED_ORIGINS")
	config := cors.DefaultConfig()
	config.AllowOrigins = strings.Split(origins, ",")
	config.AllowCredentials = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization", "Cookie")

	for _, origin := range config.AllowOrigins {
		logger.Debug("Allowed origin: %s", origin)
	}

	r.Use(cors.New(config))

	r.Static("/certifications", constants.CareerCertificationsDir)
	logger.Info("Serving static files from: %s", constants.CareerCertificationsDir)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.SetupRoutes(
		r,
		experienceHandler,
		projectHandler,
		careerCertificationHandler,
		authHandler,
		authService,
	)

	port := os.Getenv("SERVER_PORT")
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

// cleanExpiredSessions periodically removes expired sessions using the provided AuthService at a 30-day interval.
func cleanExpiredSessions(authService services.AuthService) {
	go func() {
		ticker := time.NewTicker(30 * 24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			if err := authService.CleanUpExpiredSessions(); err != nil {
				logger.Error("Failed to cleanup expired sessions: %s", err.Error())
			} else {
				logger.Info("Cleaned up expired sessions")
			}
		}
	}()
}

// configDatabaseDriver initializes and returns a database.Config based on environment variables or defaults.
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

// registerDependencies initializes and registers all necessary dependencies for handlers.
// It returns the initialized handlers and auth service.
func registerDependencies(db *gorm.DB) (
	*handlers.ProjectHandler,
	*handlers.ExperienceHandler,
	*handlers.CareerCertificationHandler,
	*handlers.AuthHandler,
	services.AuthService,
) {
	// Experience dependencies
	experienceRepo := repository.NewExperienceRepository(db)
	experienceService := services.NewExperienceService(experienceRepo)
	experienceHandler := handlers.NewExperienceHandler(experienceService)

	// Project dependencies
	projectRepo := repository.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepo)
	projectHandler := handlers.NewProjectHandler(projectService)

	// Career certification dependencies
	careerCertificationRepo := repository.NewCareerCertificationRepository(db)
	careerCertificationService := services.NewCareerCertificationService(careerCertificationRepo)
	careerCertificationHandler := handlers.NewCareerCertificationHandler(careerCertificationService)

	// Auth dependencies
	authRepo := repository.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	return projectHandler, experienceHandler, careerCertificationHandler, authHandler, authService
}
