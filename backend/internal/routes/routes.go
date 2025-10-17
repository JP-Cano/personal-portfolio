package routes

import (
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/handlers"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/handlers/dto"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, experienceHandler *handlers.ExperienceHandler, projectHandler *handlers.ProjectHandler) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Experience routes
		experiences := v1.Group("/experiences")
		{
			experiences.GET("", experienceHandler.GetAllExperiences)
			experiences.GET("/:id", experienceHandler.GetExperienceByID)
			experiences.POST("", middleware.ValidateRequest[dto.ExperienceRequest](), experienceHandler.CreateExperience)
			experiences.PATCH("/:id", middleware.ValidateRequest[dto.UpdateExperienceRequest](), experienceHandler.UpdateExperience)
			experiences.DELETE("/:id", experienceHandler.DeleteExperience)
		}

		// Project routes
		projects := v1.Group("/projects")
		{
			projects.GET("", projectHandler.GetAllProjects)
			projects.GET("/:id", projectHandler.GetProjectById)
			projects.POST("", middleware.ValidateRequest[dto.ProjectRequest](), projectHandler.CreateProject)
			projects.PATCH("/:id", middleware.ValidateRequest[dto.UpdateProjectRequest](), projectHandler.UpdateProject)
			projects.DELETE("/:id", projectHandler.DeleteProject)
		}
	}
}
