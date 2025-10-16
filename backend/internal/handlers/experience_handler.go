package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/handlers/dto"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/services"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// ExperienceHandler handles HTTP requests for experiences
type ExperienceHandler struct {
	service services.ExperienceService
}

// NewExperienceHandler creates a new instance of ExperienceHandler
func NewExperienceHandler(service services.ExperienceService) *ExperienceHandler {
	return &ExperienceHandler{service: service}
}

// GetAllExperiences godoc
// @Summary Get all experiences
// @Description Retrieves all work experiences ordered by start date
// @Tags experiences
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponse{data=[]dto.ExperienceResponse} "List of experiences"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /experiences [get]
func (h *ExperienceHandler) GetAllExperiences(c *gin.Context) {
	experiences, err := h.service.GetAllExperiences()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve experiences", err)
		return
	}

	response := dto.ToExperienceResponseList(experiences)
	utils.RespondWithSuccess(c, http.StatusOK, response, "")
}

// GetExperienceByID godoc
// @Summary Get experience by ID
// @Description Retrieves a single work experience by its ID
// @Tags experiences
// @Accept json
// @Produce json
// @Param id path int true "Experience ID"
// @Success 200 {object} utils.SuccessResponse{data=dto.ExperienceResponse} "Experience details"
// @Failure 400 {object} utils.ErrorResponse "Invalid ID format"
// @Failure 404 {object} utils.ErrorResponse "Experience not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /experiences/{id} [get]
func (h *ExperienceHandler) GetExperienceByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid experience ID", err)
		return
	}

	experience, err := h.service.GetExperienceByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrExperienceNotFound) {
			utils.RespondWithError(c, http.StatusNotFound, "Experience not found", err)
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve experience", err)
		return
	}

	response := dto.ToExperienceResponse(experience)
	utils.RespondWithSuccess(c, http.StatusOK, response, "")
}

// CreateExperience godoc
// @Summary Create a new experience
// @Description Creates a new work experience entry
// @Tags experiences
// @Accept json
// @Produce json
// @Param experience body dto.ExperienceRequest true "Experience data"
// @Success 201 {object} utils.SuccessResponse{data=dto.ExperienceResponse} "Experience created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /experiences [post]
func (h *ExperienceHandler) CreateExperience(c *gin.Context) {
	var req dto.ExperienceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	experience := req.ToExperience()

	if err := h.service.CreateExperience(experience); err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			utils.RespondWithError(c, http.StatusBadRequest, "", err)
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create experience", err)
		return
	}

	response := dto.ToExperienceResponse(experience)
	utils.RespondWithSuccess(c, http.StatusCreated, response, "Experience created successfully")
}

// UpdateExperience godoc
// @Summary Update an experience
// @Description Updates an existing work experience
// @Tags experiences
// @Accept json
// @Produce json
// @Param id path int true "Experience ID"
// @Param experience body dto.ExperienceRequest true "Updated experience data"
// @Success 200 {object} utils.SuccessResponse{data=dto.ExperienceResponse} "Experience updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid ID or request body"
// @Failure 404 {object} utils.ErrorResponse "Experience not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /experiences/{id} [put]
func (h *ExperienceHandler) UpdateExperience(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid experience ID", err)
		return
	}

	var req dto.ExperienceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	experience := req.ToExperience()

	if err := h.service.UpdateExperience(uint(id), experience); err != nil {
		if errors.Is(err, services.ErrExperienceNotFound) {
			utils.RespondWithError(c, http.StatusNotFound, "Experience not found", err)
			return
		}
		if errors.Is(err, services.ErrInvalidInput) {
			utils.RespondWithError(c, http.StatusBadRequest, "", err)
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update experience", err)
		return
	}

	response := dto.ToExperienceResponse(experience)
	utils.RespondWithSuccess(c, http.StatusOK, response, "Experience updated successfully")
}

// DeleteExperience godoc
// @Summary Delete an experience
// @Description Soft deletes a work experience (sets deleted_at timestamp)
// @Tags experiences
// @Accept json
// @Produce json
// @Param id path int true "Experience ID"
// @Success 200 {object} utils.SuccessResponse "Experience deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid ID format"
// @Failure 404 {object} utils.ErrorResponse "Experience not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /experiences/{id} [delete]
func (h *ExperienceHandler) DeleteExperience(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid experience ID", err)
		return
	}

	if err := h.service.DeleteExperience(uint(id)); err != nil {
		if errors.Is(err, services.ErrExperienceNotFound) {
			utils.RespondWithError(c, http.StatusNotFound, "Experience not found", err)
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete experience", err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Experience deleted successfully")
}
