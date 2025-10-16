package services

import (
	"errors"
	"fmt"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/models"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrExperienceNotFound = errors.New("experience not found")
	ErrInvalidInput       = errors.New("invalid input data")
)

// ExperienceService defines the interface for experience business logic
type ExperienceService interface {
	GetAllExperiences() ([]models.Experience, error)
	GetExperienceByID(id uint) (*models.Experience, error)
	CreateExperience(experience *models.Experience) error
	UpdateExperience(id uint, experience *models.Experience) error
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
	experiences, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch experiences: %w", err)
	}

	return experiences, nil
}

// GetExperienceByID retrieves a single experience by ID
func (s *experienceService) GetExperienceByID(id uint) (*models.Experience, error) {
	experience, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrExperienceNotFound
		}
		return nil, fmt.Errorf("failed to fetch experience: %w", err)
	}

	return experience, nil
}

// CreateExperience creates a new experience
func (s *experienceService) CreateExperience(experience *models.Experience) error {
	if err := s.validateExperience(experience); err != nil {
		return err
	}

	if err := s.repo.Create(experience); err != nil {
		return fmt.Errorf("failed to create experience: %w", err)
	}

	return nil
}

// UpdateExperience updates an existing experience
func (s *experienceService) UpdateExperience(id uint, experience *models.Experience) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrExperienceNotFound
		}
		return fmt.Errorf("failed to fetch experience: %w", err)
	}

	if err := s.validateExperience(experience); err != nil {
		return err
	}

	experience.ID = existing.ID

	if err := s.repo.Update(experience); err != nil {
		return fmt.Errorf("failed to update experience: %w", err)
	}

	return nil
}

// DeleteExperience deletes an experience by ID
func (s *experienceService) DeleteExperience(id uint) error {
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrExperienceNotFound
		}
		return fmt.Errorf("failed to delete experience: %w", err)
	}

	return nil
}

// validateExperience validates the experience data
func (s *experienceService) validateExperience(experience *models.Experience) error {
	if experience.Title == "" {
		return fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if experience.Company == "" {
		return fmt.Errorf("%w: company is required", ErrInvalidInput)
	}

	if experience.Type == "" {
		return fmt.Errorf("%w: type is required", ErrInvalidInput)
	}

	validTypes := map[models.WorkType]bool{
		models.Remote: true,
		models.OnSite: true,
		models.Hybrid: true,
	}

	if !validTypes[experience.Type] {
		return fmt.Errorf("%w: invalid work type", ErrInvalidInput)
	}

	if experience.StartDate.IsZero() {
		return fmt.Errorf("%w: start date is required", ErrInvalidInput)
	}

	if experience.EndDate != nil && experience.EndDate.Before(experience.StartDate) {
		return fmt.Errorf("%w: end date must be after start date", ErrInvalidInput)
	}

	return nil
}
