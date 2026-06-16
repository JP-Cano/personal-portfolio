package repository

import (
	"gorm.io/gorm"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/models"
)

type ExperienceClientRepository interface {
	FindByExperienceID(experienceID uint) ([]models.ExperienceClient, error)
	FindByID(id uint) (*models.ExperienceClient, error)
	Create(client *models.ExperienceClient) error
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
}

type experienceClientRepository struct {
	db *gorm.DB
}

func NewExperienceClientRepository(db *gorm.DB) ExperienceClientRepository {
	return &experienceClientRepository{db: db}
}

func (r *experienceClientRepository) FindByExperienceID(experienceID uint) ([]models.ExperienceClient, error) {
	var clients []models.ExperienceClient
	result := r.db.Where("experience_id = ? AND deleted_at IS NULL", experienceID).
		Order("start_date ASC").
		Find(&clients)
	if result.Error != nil {
		return nil, result.Error
	}
	return clients, nil
}

func (r *experienceClientRepository) FindByID(id uint) (*models.ExperienceClient, error) {
	var client models.ExperienceClient
	result := r.db.First(&client, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &client, nil
}

func (r *experienceClientRepository) Create(client *models.ExperienceClient) error {
	return r.db.Create(client).Error
}

func (r *experienceClientRepository) Update(id uint, updates map[string]interface{}) error {
	result := r.db.Model(&models.ExperienceClient{}).
		Where("id = ?", id).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *experienceClientRepository) Delete(id uint) error {
	result := r.db.Delete(&models.ExperienceClient{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
