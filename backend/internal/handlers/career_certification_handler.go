package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/handlers/dto"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/services"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/logger"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CareerCertificationHandler struct {
	service services.CareerCertificationService
}

func NewCareerCertificationHandler(service services.CareerCertificationService) *CareerCertificationHandler {
	return &CareerCertificationHandler{service: service}
}

// UploadAcademicCertificates handles both single and multiple file uploads with optional metadata
// @Summary Upload certification files
// @Description Upload one or multiple certification files with optional metadata
// @Tags certifications
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Certification file(s)"
// @Param workers query int false "Number of concurrent workers (default: 3, 0 = use default, max: 20)"
// @Param title formData string false "Certification title"
// @Param issuer formData string false "Issuer/Institution name"
// @Param issue_date formData string false "Issue date (DD/MM/YYYY)"
// @Param expiry_date formData string false "Expiry date (DD/MM/YYYY)"
// @Param credential_id formData string false "Credential ID"
// @Param credential_url formData string false "Credential verification URL"
// @Param description formData string false "Description"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /upload-certificates [post]
func (h *CareerCertificationHandler) UploadAcademicCertificates(c *gin.Context) {
	queryParams, _ := c.Get("validatedQuery")
	params := queryParams.(dto.UploadCertificatesRequest)

	form, err := c.MultipartForm()
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Failed to parse form data", err)
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "No files uploaded", nil)
		return
	}

	var metadata dto.CertificationMetadata
	if err := c.ShouldBind(&metadata); err != nil {
		logger.Warn("Failed to parse metadata: %v", err)
	}

	filesWithMetadata := make([]services.FileWithMetadata, len(files))
	for i, file := range files {
		filesWithMetadata[i] = services.FileWithMetadata{
			File:     file,
			Metadata: &metadata,
		}
	}

	timeout := 30 * time.Second
	if len(files) > 5 {
		timeout = 5 * time.Minute
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s/certifications", scheme, c.Request.Host)

	results := h.service.StoreBatch(ctx, filesWithMetadata, params.Workers, baseURL)

	var successful []map[string]interface{}
	var failed []map[string]string

	for _, result := range results {
		if result.Success {
			successful = append(successful, map[string]interface{}{
				"id":            result.Certification.ID,
				"title":         result.Certification.Title,
				"issuer":        result.Certification.Issuer,
				"issue_date":    result.Certification.IssueDate,
				"file_url":      result.Certification.FileURL,
				"file_name":     result.Certification.FileName,
				"original_name": result.Certification.OriginalName,
			})
		} else {
			failed = append(failed, map[string]string{
				"original": result.OriginalName,
				"error":    result.Error.Error(),
			})
		}
	}

	statusCode := http.StatusOK
	if len(successful) == 0 {
		statusCode = http.StatusInternalServerError
	} else if len(failed) > 0 {
		statusCode = http.StatusMultiStatus
	}

	utils.RespondWithSuccess(c, statusCode, map[string]interface{}{
		"total":      len(files),
		"successful": len(successful),
		"failed":     len(failed),
		"files":      successful,
		"errors":     failed,
	}, "Files upload completed")
}

// GetAllCertifications godoc
// @Summary Get all certifications
// @Description Retrieves all uploaded career certifications with metadata
// @Tags certifications
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponse{data=[]models.CareerCertification} "List of certifications"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /upload-certificates [get]
func (h *CareerCertificationHandler) GetAllCertifications(c *gin.Context) {
	certifications, err := h.service.GetAll()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve certifications", err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, certifications, "")
}

// GetCertificationByID godoc
// @Summary Get certification by ID
// @Description Retrieves a single certification by its ID
// @Tags certifications
// @Accept json
// @Produce json
// @Param id path int true "Certification ID"
// @Success 200 {object} utils.SuccessResponse{data=models.CareerCertification} "Certification details"
// @Failure 400 {object} utils.ErrorResponse "Invalid ID format"
// @Failure 404 {object} utils.ErrorResponse "Certification not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /upload-certificates/{id} [get]
func (h *CareerCertificationHandler) GetCertificationByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	certification, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Certification not found", err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, certification, "")
}

// DeleteCertification godoc
// @Summary Delete a certification
// @Description Soft deletes a certification and removes its associated file from storage
// @Tags certifications
// @Accept json
// @Produce json
// @Param id path int true "Certification ID"
// @Success 200 {object} utils.SuccessResponse "Certification deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid ID format"
// @Failure 404 {object} utils.ErrorResponse "Certification not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /upload-certificates/{id} [delete]
func (h *CareerCertificationHandler) DeleteCertification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete certification", err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Certification deleted successfully")
}
