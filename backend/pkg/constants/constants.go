package constants

import (
	"errors"
	"os"
)

var (
	ErrExperienceNotFound = errors.New("experience not found")
	ErrProjectNotFound    = errors.New("project not found")
)

const CareerCertificationsDir = "pkg/assets/career-certifications"

// DefaultSessionCookieName is the default name for the session cookie
const DefaultSessionCookieName = "portfolio_session"

// GetSessionCookieName returns the session cookie name from env or default
func GetSessionCookieName() string {
	name := os.Getenv("SESSION_COOKIE_NAME")
	if name == "" {
		return DefaultSessionCookieName
	}
	return name
}
