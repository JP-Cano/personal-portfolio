package repository

import (
	"errors"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/models"
	"gorm.io/gorm"
)

// ExperienceRepository defines the interface for experience data operations
type ExperienceRepository interface {
	FindAll() ([]models.Experience, error)
	FindByID(id uint) (*models.Experience, error)
	Create(experience *models.Experience) error
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
}

// experienceRepository implements ExperienceRepository interface
type experienceRepository struct {
	db *gorm.DB
}

// NewExperienceRepository creates a new instance of ExperienceRepository
func NewExperienceRepository(db *gorm.DB) ExperienceRepository {
	return &experienceRepository{db: db}
}

// FindAll retrieves all experiences from the database
func (r *experienceRepository) FindAll() ([]models.Experience, error) {
	var experiences []models.Experience

	result := r.db.Order("start_date DESC").Find(&experiences)

	if result.Error != nil {
		return nil, result.Error
	}

	return experiences, nil
}

// FindByID retrieves a single experience by ID
func (r *experienceRepository) FindByID(id uint) (*models.Experience, error) {
	var experience models.Experience

	result := r.db.First(&experience, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}

	return &experience, nil
}

// Create inserts a new experience into the database
func (r *experienceRepository) Create(experience *models.Experience) error {
	result := r.db.Create(experience)
	return result.Error
}

// Update modifies an existing experience in the database
func (r *experienceRepository) Update(id uint, updates map[string]interface{}) error {
	result := r.db.Model(&models.Experience{}).Where("id = ?", id).Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete removes an experience from the database
func (r *experienceRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Experience{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
