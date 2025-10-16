package routes

import (
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, experienceHandler *handlers.ExperienceHandler) {
	// Health check endpoint (outside of API versioning)
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
			experiences.POST("", experienceHandler.CreateExperience)
			experiences.PUT("/:id", experienceHandler.UpdateExperience)
			experiences.DELETE("/:id", experienceHandler.DeleteExperience)
		}
	}
}
