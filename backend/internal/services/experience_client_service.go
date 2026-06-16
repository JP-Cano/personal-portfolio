package services

import (
	"errors"
	"fmt"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/models"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/repository"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/constants"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/logger"
	"gorm.io/gorm"
)

type ExperienceClientService interface {
	GetClientsByExperienceID(experienceID uint) ([]models.ExperienceClient, error)
	GetClientByID(id uint) (*models.ExperienceClient, error)
	CreateClient(experienceID uint, client *models.ExperienceClient) error
	UpdateClient(id uint, updates map[string]interface{}) error
	DeleteClient(id uint) error
}

type experienceClientService struct {
	repo repository.ExperienceClientRepository
}

func NewExperienceClientService(repo repository.ExperienceClientRepository) ExperienceClientService {
	return &experienceClientService{repo: repo}
}

func (s *experienceClientService) GetClientsByExperienceID(experienceID uint) ([]models.ExperienceClient, error) {
	logger.Debug("Fetching clients for experience ID: %d", experienceID)
	clients, err := s.repo.FindByExperienceID(experienceID)
	if err != nil {
		logger.Error("Failed to get clients for experience %d: %v", experienceID, err)
		return nil, fmt.Errorf("getting clients by experience ID: %w", err)
	}
	logger.Info("Successfully fetched %d clients for experience %d", len(clients), experienceID)
	return clients, nil
}

func (s *experienceClientService) GetClientByID(id uint) (*models.ExperienceClient, error) {
	logger.Debug("Fetching client with ID: %d", id)
	client, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Experience client not found: %d", id)
			return nil, constants.ErrExperienceClientNotFound
		}
		logger.Error("Failed to get client %d: %v", id, err)
		return nil, fmt.Errorf("getting client by ID: %w", err)
	}
	logger.Info("Successfully fetched client: %d", id)
	return client, nil
}

func (s *experienceClientService) CreateClient(experienceID uint, client *models.ExperienceClient) error {
	logger.Info("Creating new client for experience ID: %d", experienceID)
	client.ExperienceID = experienceID
	if err := s.repo.Create(client); err != nil {
		logger.Error("Failed to create client for experience %d: %v", experienceID, err)
		return fmt.Errorf("creating client: %w", err)
	}
	logger.Info("Successfully created client with ID: %d for experience %d", client.ID, experienceID)
	return nil
}

func (s *experienceClientService) UpdateClient(id uint, updates map[string]interface{}) error {
	logger.Info("Updating client with ID: %d", id)
	if err := s.repo.Update(id, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Experience client not found for update: %d", id)
			return constants.ErrExperienceClientNotFound
		}
		logger.Error("Failed to update client %d: %v", id, err)
		return fmt.Errorf("updating client: %w", err)
	}
	logger.Info("Successfully updated client: %d", id)
	return nil
}

func (s *experienceClientService) DeleteClient(id uint) error {
	logger.Info("Deleting client with ID: %d", id)
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Experience client not found for deletion: %d", id)
			return constants.ErrExperienceClientNotFound
		}
		logger.Error("Failed to delete client %d: %v", id, err)
		return fmt.Errorf("deleting client: %w", err)
	}
	logger.Info("Successfully deleted client: %d", id)
	return nil
}
