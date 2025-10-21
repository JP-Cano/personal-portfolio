package repository

import (
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/models"
	"gorm.io/gorm"
)

type CareerCertificationRepository interface {
	Create(certification *models.CareerCertification) error
	FindAll() ([]models.CareerCertification, error)
	FindByID(id uint) (*models.CareerCertification, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
}

// careerCertificationRepository provides methods to interact with the career certifications data in the database.
// It is a concrete implementation of the CareerCertificationRepository interface.
// Uses gorm.DB for database operations.
type careerCertificationRepository struct {
	db *gorm.DB
}

// NewCareerCertificationRepository initializes and returns an implementation of CareerCertificationRepository.
func NewCareerCertificationRepository(db *gorm.DB) CareerCertificationRepository {
	return &careerCertificationRepository{db: db}
}

// Create inserts a new CareerCertification record into the database and returns an error if the operation fails.
func (r *careerCertificationRepository) Create(certification *models.CareerCertification) error {
	return r.db.Create(certification).Error
}

// FindAll retrieves all career certifications from the database sorted by issue date in descending order.
func (r *careerCertificationRepository) FindAll() ([]models.CareerCertification, error) {
	var certifications []models.CareerCertification
	err := r.db.Order("issue_date DESC").Find(&certifications).Error
	return certifications, err
}

// FindByID retrieves a CareerCertification by its ID from the database. Returns the record or an error if not found.
func (r *careerCertificationRepository) FindByID(id uint) (*models.CareerCertification, error) {
	var certification models.CareerCertification
	err := r.db.First(&certification, id).Error
	return &certification, err
}

// Update updates fields of a CareerCertification in the database identified by the given ID using the provided map of updates.
func (r *careerCertificationRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.CareerCertification{}).Where("id = ?", id).Updates(updates).Error
}

// Delete removes a CareerCertification record from the database by its ID and returns an error if the operation fails.
func (r *careerCertificationRepository) Delete(id uint) error {
	return r.db.Delete(&models.CareerCertification{}, id).Error
}
