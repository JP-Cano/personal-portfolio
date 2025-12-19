package repository

import (
	"errors"
	"time"

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

// NewAuthRepository creates a new instance of AuthRepository
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// FindUserByEmail retrieves a user by their email address
func (r *authRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User

	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

// CreateSession inserts a new session into the database
func (r *authRepository) CreateSession(session *models.Session) error {
	result := r.db.Create(session)
	return result.Error
}

// FindSessionByID retrieves a session by its ID with the associated user preloaded
func (r *authRepository) FindSessionByID(sessionID string) (*models.Session, error) {
	var session models.Session

	result := r.db.Preload("User").Where("id = ?", sessionID).First(&session)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}

	return &session, nil
}

// DeleteSession removes a session from the database
func (r *authRepository) DeleteSession(sessionID string) error {
	result := r.db.Where("id = ?", sessionID).Delete(&models.Session{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// DeleteExpiredSessions removes all expired sessions from the database
func (r *authRepository) DeleteExpiredSessions() error {
	result := r.db.Where("expires_at < ?", time.Now()).Delete(&models.Session{})
	return result.Error
}
