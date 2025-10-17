package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/handlers/dto"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/services"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/constants"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// ProjectHandler handles HTTP requests for projects
type ProjectHandler struct {
	service services.ProjectService
}

// NewProjectHandler creates a new instance of ProjectHandler
func NewProjectHandler(service services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

// GetAllProjects godoc
// @Summary Get all projects
// @Description Retrieves all projects
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponse{data=[]dto.ProjectResponse} "List of projects"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /projects [get]
func (p *ProjectHandler) GetAllProjects(c *gin.Context) {
	projects, err := p.service.GetAllProjects()

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve projects", err)
		return
	}

	response := dto.ToProjectResponseList(projects)
	utils.RespondWithSuccess(c, http.StatusOK, response, "")
}

// GetProjectById godoc
// @Summary Get project by ID
// @Description Retrieves a single project by its ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} utils.SuccessResponse{data=dto.ProjectResponse} "Project details"
// @Failure 400 {object} utils.ErrorResponse "Invalid ID format"
// @Failure 404 {object} utils.ErrorResponse "Project not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /projects/{id} [get]
func (p *ProjectHandler) GetProjectById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid project ID", err)
		return
	}

	project, err := p.service.GetProjectByID(uint(id))

	if err != nil {
		if errors.Is(err, constants.ErrProjectNotFound) {
			utils.RespondWithError(c, http.StatusNotFound, "Project not found", err)
		}
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve project", err)
		return
	}

	response := dto.ToProjectResponse(project)
	utils.RespondWithSuccess(c, http.StatusOK, response, "")
}

// CreateProject godoc
// @Summary Create a new project
// @Description Creates a new project entry
// @Tags projects
// @Accept json
// @Produce json
// @Param project body dto.ProjectRequest true "Project data"
// @Success 201 {object} utils.SuccessResponse{data=dto.ProjectResponse} "Project created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /projects [post]
func (p *ProjectHandler) CreateProject(c *gin.Context) {
	req, exists := c.Get("validatedRequest")

	if !exists {
		utils.RespondWithError(c, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	projectReq := req.(dto.ProjectRequest)
	project, err := projectReq.ToProject()
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid date format", err)
		return
	}

	if err := p.service.CreateProject(project); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create project", err)
		return
	}

	response := dto.ToProjectResponse(project)
	utils.RespondWithSuccess(c, http.StatusCreated, response, "Project created successfully")
}

// UpdateProject godoc
// @Summary Update a project
// @Description Updates an existing project (partial update supported - send only fields to update)
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param project body dto.UpdateProjectRequest true "Fields to update"
// @Success 200 {object} utils.SuccessResponse{data=dto.ProjectResponse} "Project updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid ID or request body"
// @Failure 404 {object} utils.ErrorResponse "Project not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /projects/{id} [patch]
func (p *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid project ID", err)
		return
	}

	req, exists := c.Get("validatedRequest")
	if !exists {
		utils.RespondWithError(c, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	updateReq := req.(dto.UpdateProjectRequest)

	// Convert to update map
	updates, err := updateReq.ToUpdateMap()
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid update data", err)
		return
	}

	if err := p.service.UpdateProject(uint(id), updates); err != nil {
		if errors.Is(err, constants.ErrProjectNotFound) {
			utils.RespondWithError(c, http.StatusNotFound, "Project not found", err)
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update project", err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Project updated successfully")
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Soft deletes a project (sets deleted_at timestamp)
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} utils.SuccessResponse "Project deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid ID format"
// @Failure 404 {object} utils.ErrorResponse "Project not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /projects/{id} [delete]
func (p *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid project ID", err)
		return
	}

	if err := p.service.DeleteProject(uint(id)); err != nil {
		if errors.Is(err, constants.ErrProjectNotFound) {
			utils.RespondWithError(c, http.StatusNotFound, "Project not found", err)
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete project", err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Project deleted successfully")
	return
}
