package services

import (
	"errors"
	"time"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/models"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// SessionDuration Session duration - 2 hours
const SessionDuration = 2 * time.Hour

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionExpired     = errors.New("session has expired")
)

type AuthService interface {
	Login(email, password string) (*models.Session, *models.AuthResponse, error)
	Logout(sessionId string) error
	CleanUpExpiredSessions() error
	ValidateSession(sessionId string) (*models.Session, error)
	GetUserBySessionID(sessionId string) (*models.UserResponse, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

// Login authenticates a user and creates a new session
func (a *authService) Login(email, password string) (*models.Session, *models.AuthResponse, error) {
	user, err := a.repo.FindUserByEmail(email)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	session := &models.Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(SessionDuration),
		CreatedAt: time.Now(),
	}

	if err := a.repo.CreateSession(session); err != nil {
		return nil, nil, err
	}

	response := &models.AuthResponse{
		User: models.UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
		Message: "Login successful",
	}

	return session, response, nil
}

// Logout removes a session
func (a *authService) Logout(sessionId string) error {
	return a.repo.DeleteSession(sessionId)
}

// CleanUpExpiredSessions removes all expired sessions
func (a *authService) CleanUpExpiredSessions() error {
	return a.repo.DeleteExpiredSessions()
}

// ValidateSession checks if a session exists and is not expired
func (a *authService) ValidateSession(sessionId string) (*models.Session, error) {
	session, err := a.repo.FindSessionByID(sessionId)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		_ = a.repo.DeleteSession(sessionId)
		return nil, ErrSessionExpired
	}

	return session, nil
}

// GetUserBySessionID returns user data for a valid session
func (a *authService) GetUserBySessionID(sessionId string) (*models.UserResponse, error) {
	session, err := a.ValidateSession(sessionId)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:    session.User.ID,
		Email: session.User.Email,
	}, nil
}
