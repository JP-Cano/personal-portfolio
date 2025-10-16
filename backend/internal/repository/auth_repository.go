package repository

import (
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindUserByEmail(email string) (*models.User, error)
	CreateSession(session *models.Session) error
	FindSessionByID(sessionID string) (*models.Session, error)
	DeleteSession(sessionID string) error
	DeleteExpiredSessions() error
}

type authRepository struct {
	db *gorm.DB
}
