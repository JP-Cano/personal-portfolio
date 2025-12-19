package handlers

import (
	"errors"
	"net/http"
	"os"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/handlers/dto"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/services"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/constants"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/logger"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// CookieMaxAge Cookie max age in seconds (2 hours)
const CookieMaxAge = 2 * 60 * 60

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	service services.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Login godoc
// @Summary Login user
// @Description Authenticates a user with email and password, sets session cookie
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Login credentials"
// @Success 200 {object} utils.SuccessResponse{data=models.AuthResponse} "Login successful"
// @Failure 400 {object} utils.ErrorResponse "Invalid request body"
// @Failure 401 {object} utils.ErrorResponse "Invalid credentials"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	session, authResponse, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			utils.RespondWithError(c, http.StatusUnauthorized, "Invalid email or password", err)
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, "Login failed", err)
		return
	}

	h.setSessionCookie(c, session.ID)

	utils.RespondWithSuccess(c, http.StatusOK, authResponse, "Login successful")
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Returns the currently authenticated user based on session cookie
// @Tags auth
// @Produce json
// @Success 200 {object} utils.SuccessResponse{data=models.UserResponse} "Current user data"
// @Failure 401 {object} utils.ErrorResponse "Not authenticated"
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	sessionID, err := c.Cookie(constants.GetSessionCookieName())
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "Not authenticated", err)
		return
	}

	user, err := h.service.GetUserBySessionID(sessionID)
	if err != nil {
		if errors.Is(err, services.ErrSessionNotFound) || errors.Is(err, services.ErrSessionExpired) {
			h.clearSessionCookie(c)
			utils.RespondWithError(c, http.StatusUnauthorized, "Session expired or invalid", err)
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get user", err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, user, "")
}

// Logout godoc
// @Summary Logout user
// @Description Logs out the current user by deleting their session
// @Tags auth
// @Produce json
// @Success 200 {object} utils.SuccessResponse "Logout successful"
// @Failure 401 {object} utils.ErrorResponse "Not authenticated"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID, err := c.Cookie(constants.GetSessionCookieName())
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "Not authenticated", err)
		return
	}

	if err := h.service.Logout(sessionID); err != nil {
		logger.Warn("Failed to delete session during logout: %s", err.Error())
	}

	h.clearSessionCookie(c)

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Logout successful")
}

// setSessionCookie sets the session cookie with appropriate security settings
func (h *AuthHandler) setSessionCookie(c *gin.Context, sessionID string) {
	secure := os.Getenv("COOKIE_SECURE") == "true"
	domain := os.Getenv("COOKIE_DOMAIN")
	sameSite := http.SameSiteLaxMode

	if os.Getenv("COOKIE_SAME_SITE") == "none" {
		sameSite = http.SameSiteNoneMode
	}

	c.SetSameSite(sameSite)
	c.SetCookie(
		constants.GetSessionCookieName(),
		sessionID,
		CookieMaxAge,
		"/",
		domain,
		secure,
		true,
	)
}

// clearSessionCookie clears the session cookie
func (h *AuthHandler) clearSessionCookie(c *gin.Context) {
	secure := os.Getenv("COOKIE_SECURE") == "true"
	domain := os.Getenv("COOKIE_DOMAIN")
	sameSite := http.SameSiteLaxMode

	if os.Getenv("COOKIE_SAME_SITE") == "none" {
		sameSite = http.SameSiteNoneMode
	}

	c.SetSameSite(sameSite)
	c.SetCookie(
		constants.GetSessionCookieName(),
		"",
		-1,
		"/",
		domain,
		secure,
		true,
	)
}
