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

// ExperienceService defines the interface for experience business logic
type ExperienceService interface {
	GetAllExperiences() ([]models.Experience, error)
	GetExperienceByID(id uint) (*models.Experience, error)
	CreateExperience(experience *models.Experience) error
	UpdateExperience(id uint, updates map[string]interface{}) error
	DeleteExperience(id uint) error
}

// experienceService implements ExperienceService interface
type experienceService struct {
	repo repository.ExperienceRepository
}

// NewExperienceService creates a new instance of ExperienceService
func NewExperienceService(repo repository.ExperienceRepository) ExperienceService {
	return &experienceService{repo: repo}
}

// GetAllExperiences retrieves all experiences
func (s *experienceService) GetAllExperiences() ([]models.Experience, error) {
	logger.Debug("Fetching all experiences")
	experiences, err := s.repo.FindAll()
	if err != nil {
		logger.Error("Failed to fetch experiences: %v", err)
		return nil, fmt.Errorf("failed to fetch experiences: %w", err)
	}

	logger.Info("Successfully fetched %d experiences", len(experiences))
	return experiences, nil
}

// GetExperienceByID retrieves a single experience by ID
func (s *experienceService) GetExperienceByID(id uint) (*models.Experience, error) {
	logger.Debug("Fetching experience with ID: %d", id)
	experience, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Experience not found: %d", id)
			return nil, constants.ErrExperienceNotFound
		}
		logger.Error("Failed to fetch experience %d: %v", id, err)
		return nil, fmt.Errorf("failed to fetch experience: %w", err)
	}

	logger.Info("Successfully fetched experience: %d", id)
	return experience, nil
}

// CreateExperience creates a new experience
func (s *experienceService) CreateExperience(experience *models.Experience) error {
	logger.Info("Creating new experience: %s at %s", experience.Title, experience.Company)
	if err := s.repo.Create(experience); err != nil {
		logger.Error("Failed to create experience: %v", err)
		return fmt.Errorf("failed to create experience: %w", err)
	}

	logger.Info("Successfully created experience with ID: %d", experience.ID)
	return nil
}

// UpdateExperience updates an existing experience
func (s *experienceService) UpdateExperience(id uint, updates map[string]interface{}) error {
	logger.Info("Updating experience with ID: %d", id)
	if err := s.repo.Update(id, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Experience not found for update: %d", id)
			return constants.ErrExperienceNotFound
		}
		logger.Error("Failed to update experience %d: %v", id, err)
		return fmt.Errorf("failed to update experience: %w", err)
	}

	logger.Info("Successfully updated experience: %d", id)
	return nil
}

// DeleteExperience deletes an experience by ID
func (s *experienceService) DeleteExperience(id uint) error {
	logger.Info("Deleting experience with ID: %d", id)
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Experience not found for deletion: %d", id)
			return constants.ErrExperienceNotFound
		}
		logger.Error("Failed to delete experience %d: %v", id, err)
		return fmt.Errorf("failed to delete experience: %w", err)
	}

	logger.Info("Successfully deleted experience: %d", id)
	return nil
}
