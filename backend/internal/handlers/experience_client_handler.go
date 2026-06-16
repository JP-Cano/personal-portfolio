package handlers

import (
	"net/http"
	"strconv"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/handlers/dto"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/services"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/constants"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ExperienceClientHandler struct {
	service services.ExperienceClientService
}

func NewExperienceClientHandler(service services.ExperienceClientService) *ExperienceClientHandler {
	return &ExperienceClientHandler{service: service}
}

func (h *ExperienceClientHandler) GetClientsByExperienceID(c *gin.Context) {
	experienceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid experience ID", err)
		return
	}

	clients, err := h.service.GetClientsByExperienceID(uint(experienceID))
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get clients", err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, dto.ToExperienceClientResponseList(clients), "")
}

func (h *ExperienceClientHandler) GetClientByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("clientId"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid client ID", err)
		return
	}

	client, err := h.service.GetClientByID(uint(id))
	if err != nil {
		if err == constants.ErrExperienceClientNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Client not found", err)
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get client", err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, dto.ToExperienceClientResponse(client), "")
}

func (h *ExperienceClientHandler) CreateClient(c *gin.Context) {
	experienceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid experience ID", err)
		return
	}

	req, exists := c.Get("validatedRequest")
	if !exists {
		utils.RespondWithError(c, http.StatusBadRequest, "Missing request body", nil)
		return
	}

	clientRequest, ok := req.(dto.ExperienceClientRequest)
	if !ok {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request type", nil)
		return
	}

	client, err := clientRequest.ToExperienceClient()
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	if err := h.service.CreateClient(uint(experienceID), client); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create client", err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusCreated, dto.ToExperienceClientResponse(client), "Client created")
}

func (h *ExperienceClientHandler) UpdateClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("clientId"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid client ID", err)
		return
	}

	req, exists := c.Get("validatedRequest")
	if !exists {
		utils.RespondWithError(c, http.StatusBadRequest, "Missing request body", nil)
		return
	}

	updateRequest, ok := req.(dto.UpdateExperienceClientRequest)
	if !ok {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request type", nil)
		return
	}

	updates, err := updateRequest.ToUpdateMap()
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	if err := h.service.UpdateClient(uint(id), updates); err != nil {
		if err == constants.ErrExperienceClientNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Client not found", err)
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update client", err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Client updated")
}

func (h *ExperienceClientHandler) DeleteClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("clientId"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid client ID", err)
		return
	}

	if err := h.service.DeleteClient(uint(id)); err != nil {
		if err == constants.ErrExperienceClientNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Client not found", err)
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete client", err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Client deleted")
}
