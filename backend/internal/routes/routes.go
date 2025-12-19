package routes

import (
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/handlers"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/handlers/dto"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/middleware"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	router *gin.Engine,
	experienceHandler *handlers.ExperienceHandler,
	projectHandler *handlers.ProjectHandler,
	uploadCertificatesHandler *handlers.CareerCertificationHandler,
	authHandler *handlers.AuthHandler,
	authService services.AuthService,
) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Health check for API
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "API is running",
				"version": "1.0",
			})
		})

		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", authHandler.GetCurrentUser)
			auth.POST("/logout", authHandler.Logout)
		}

		// Experience routes
		experiences := v1.Group("/experiences")
		{
			// Public routes
			experiences.GET("", experienceHandler.GetAllExperiences)
			experiences.GET("/:id", experienceHandler.GetExperienceByID)

			// Protected routes
			experiences.POST("",
				middleware.AuthMiddleware(authService),
				middleware.ValidateRequest[dto.ExperienceRequest](),
				experienceHandler.CreateExperience,
			)
			experiences.PATCH("/:id",
				middleware.AuthMiddleware(authService),
				middleware.ValidateRequest[dto.UpdateExperienceRequest](),
				experienceHandler.UpdateExperience,
			)
			experiences.DELETE("/:id",
				middleware.AuthMiddleware(authService),
				experienceHandler.DeleteExperience,
			)
		}

		// Project routes
		projects := v1.Group("/projects")
		{
			// Public routes
			projects.GET("", projectHandler.GetAllProjects)
			projects.GET("/:id", projectHandler.GetProjectById)

			// Protected routes
			projects.POST("",
				middleware.AuthMiddleware(authService),
				middleware.ValidateRequest[dto.ProjectRequest](),
				projectHandler.CreateProject,
			)
			projects.PATCH("/:id",
				middleware.AuthMiddleware(authService),
				middleware.ValidateRequest[dto.UpdateProjectRequest](),
				projectHandler.UpdateProject,
			)
			projects.DELETE("/:id",
				middleware.AuthMiddleware(authService),
				projectHandler.DeleteProject,
			)
		}

		// Upload Certificates
		uploadCertificates := v1.Group("/upload-certificates")
		{
			// Public routes
			uploadCertificates.GET("", uploadCertificatesHandler.GetAllCertifications)
			uploadCertificates.GET("/:id", uploadCertificatesHandler.GetCertificationByID)

			// Protected routes
			uploadCertificates.POST("",
				middleware.AuthMiddleware(authService),
				middleware.ValidateQuery[dto.UploadCertificatesRequest](),
				uploadCertificatesHandler.UploadAcademicCertificates,
			)
			uploadCertificates.DELETE("/:id",
				middleware.AuthMiddleware(authService),
				uploadCertificatesHandler.DeleteCertification,
			)
		}
	}
}
