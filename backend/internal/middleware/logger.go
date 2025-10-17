package middleware

import (
	"time"

	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs HTTP requests and responses
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get status code
		statusCode := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get request method and path
		method := c.Request.Method
		path := c.Request.URL.Path

		// Create logger with context fields
		log := logger.WithFields(map[string]interface{}{
			"status":     statusCode,
			"method":     method,
			"path":       path,
			"ip":         clientIP,
			"latency":    latency.String(),
			"user_agent": c.Request.UserAgent(),
		})

		// Log based on status code
		switch {
		case statusCode >= 500:
			log.Error("Server error")
		case statusCode >= 400:
			log.Warn("Client error")
		case statusCode >= 300:
			log.Info("Redirection")
		case statusCode >= 200:
			log.Info("Success")
		default:
			log.Info("Request processed")
		}
	}
}

// RecoveryMiddleware recovers from panics and logs them
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.WithFields(map[string]interface{}{
					"error":  err,
					"path":   c.Request.URL.Path,
					"method": c.Request.Method,
				}).Error("Panic recovered: %v", err)

				c.JSON(500, gin.H{
					"success": false,
					"error":   "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
