package middleware

import (
	"net/http"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/services"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/constants"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

const (
	UserContextKey    = "user"
	SessionContextKey = "session"
)

// AuthMiddleware creates a middleware that validates session cookies
// and sets the user in the Gin context
func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie(constants.GetSessionCookieName())
		if err != nil {
			utils.RespondWithError(c, http.StatusUnauthorized, "Authentication required", nil)
			c.Abort()
			return
		}

		session, err := authService.ValidateSession(sessionID)
		if err != nil {
			utils.RespondWithError(c, http.StatusUnauthorized, "Invalid or expired session", err)
			c.Abort()
			return
		}

		c.Set(SessionContextKey, session)
		c.Set(UserContextKey, &session.User)

		c.Next()
	}
}
