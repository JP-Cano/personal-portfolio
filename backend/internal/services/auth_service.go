package services

import (
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/models"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/repository"
)

type AuthService interface {
	Login(email, password string) (*models.AuthResponse, error)
	Logout(sessionId string) error
	CleanUpExpiredSessions() error
	ValidateSession(sessionId string) error
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (a authService) CleanUpExpiredSessions() error {
	//TODO implement me
	panic("implement me")
}

func (a authService) ValidateSession(sessionId string) error {
	//TODO implement me
	panic("implement me")
}

func (a authService) Login(email, password string) (*models.AuthResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a authService) Logout(sessionId string) error {
	//TODO implement me
	panic("implement me")
}
