package repository

import (
	"errors"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/models"
	"gorm.io/gorm"
)

// ProjectRepository defines a contract for operations on the Project model.
// It includes methods for CRUD operations and retrieval of project data.
// FindAll retrieves all projects, ordered by the default criteria.
// FindByID retrieves a single project by its unique identifier.
// Create adds a new project to the repository.
// Update modifies the details of an existing project identified by ID.
// Delete removes a project by its ID from the repository.
type ProjectRepository interface {
	FindAll() ([]models.Project, error)
	FindByID(id uint) (*models.Project, error)
	Create(project *models.Project) error
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
}

// projectRepository is a struct that interacts with the database to manage Project entities using gorm.DB.
type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new instance of ProjectRepository with the provided gorm.DB dependency.
func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

// FindAll retrieves all projects from the database, ordered by start_date in descending order. Returns a slice of projects or an error.
func (p *projectRepository) FindAll() ([]models.Project, error) {
	var projects []models.Project

	result := p.db.Order("start_date DESC").Find(&projects)

	if result.Error != nil {
		return nil, result.Error
	}

	return projects, nil
}

// FindByID retrieves a single project by its unique identifier (ID) from the database.
// Returns the project and nil if found, or nil and an error if not found or another database error occurs.
func (p *projectRepository) FindByID(id uint) (*models.Project, error) {
	var project models.Project

	result := p.db.First(&project, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}

	return &project, nil
}

// Create inserts a new project record into the database and returns an error if the operation fails.
func (p *projectRepository) Create(project *models.Project) error {
	result := p.db.Create(project)
	return result.Error
}

// Update modifies an existing project in the database using the provided project ID and updated project data.
func (p *projectRepository) Update(id uint, updates map[string]interface{}) error {
	result := p.db.Model(&models.Project{}).Where("id = ?", id).Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete removes a Project record from the database by its ID. It returns an error if deletion fails or the record is not found.
func (p *projectRepository) Delete(id uint) error {
	result := p.db.Delete(&models.Project{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
