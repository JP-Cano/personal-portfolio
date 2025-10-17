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

// ProjectService defines a contract for managing Project resources in the application.
// GetAllProjects retrieves all projects from the data source.
// GetProjectByID fetches a project by its unique identifier.
// CreateProject adds a new project to the data source.
// UpdateProject modifies an existing project specified by its identifier.
// DeleteProject removes a project identified by its unique ID from the data source.
type ProjectService interface {
	GetAllProjects() ([]models.Project, error)
	GetProjectByID(id uint) (*models.Project, error)
	CreateProject(project *models.Project) error
	UpdateProject(id uint, updates map[string]interface{}) error
	DeleteProject(id uint) error
}

// projectService is a service struct that facilitates operations related to projects using the provided repository.
type projectService struct {
	repo repository.ProjectRepository
}

// NewProjectService initializes and returns a concrete implementation of ProjectService using the provided repository.
func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &projectService{repo: repo}
}

// GetAllProjects retrieves all projects from the repository.
// Returns a slice of Project models and an error if any occurred.
func (p *projectService) GetAllProjects() ([]models.Project, error) {
	logger.Debug("Fetching all projects")
	projects, err := p.repo.FindAll()

	if err != nil {
		logger.Error("Failed to fetch projects: %v", err)
		return nil, fmt.Errorf("failed to fetch projects: %w", err)
	}

	logger.Info("Successfully fetched %d projects", len(projects))
	return projects, nil
}

// GetProjectByID retrieves a project by its unique identifier.
// Returns the project instance and an error if any occurs during retrieval.
func (p *projectService) GetProjectByID(id uint) (*models.Project, error) {
	logger.Debug("Fetching project with ID: %d", id)
	project, err := p.repo.FindByID(id)

	if err != nil {
		logger.Error("Failed to fetch project with ID %d: %v", id, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrProjectNotFound
		}

		return nil, fmt.Errorf("failed to fetch project: %w", err)
	}

	logger.Info("Successfully fetched project with ID: %d", id)
	return project, nil
}

// CreateProject creates a new project record in the repository and logs the operation details. Returns an error if failed.
func (p *projectService) CreateProject(project *models.Project) error {
	logger.Debug("Creating project with name: %s", project.Name)
	if err := p.repo.Create(project); err != nil {
		logger.Error("Failed to create project: %v", err)
		return fmt.Errorf("failed to create project: %w", err)
	}
	logger.Info("Successfully created project with ID: %d", project.ID)
	return nil
}

// UpdateProject updates an existing project identified by its ID with the provided project data.
func (p *projectService) UpdateProject(id uint, updates map[string]interface{}) error {
	logger.Info("Updating project with ID: %d", id)
	if err := p.repo.Update(id, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Project not found for update: %d", id)
			return constants.ErrProjectNotFound
		}
		logger.Error("Failed to update project %d: %v", id, err)
		return fmt.Errorf("failed to update project: %w", err)
	}

	logger.Info("Successfully updated project: %d", id)
	return nil
}

// DeleteProject removes a project by its unique ID from the repository.
// Returns an error if the deletion fails or the project does not exist.
func (p *projectService) DeleteProject(id uint) error {
	logger.Info("Deleting project with ID: %d", id)
	err := p.repo.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Project not found for deletion: %d", id)
			return constants.ErrProjectNotFound
		}
		logger.Error("Failed to delete project %d: %v", id, err)
		return fmt.Errorf("failed to delete project: %w", err)
	}

	logger.Info("Successfully deleted project: %d", id)
	return nil
}
