package dto

import (
	"time"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/models"
)

// ExperienceResponse represents the API response for an experience
type ExperienceResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Company     string     `json:"company"`
	Location    string     `json:"location"`
	Type        string     `json:"type"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ExperienceRequest represents the API request for creating/updating an experience
type ExperienceRequest struct {
	Title       string     `json:"title" binding:"required"`
	Company     string     `json:"company" binding:"required"`
	Location    string     `json:"location"`
	Type        string     `json:"type" binding:"required,oneof=Remote 'On Site' Hybrid"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Description string     `json:"description"`
}

// ToExperienceResponse converts a models.Experience to ExperienceResponse
func ToExperienceResponse(experience *models.Experience) ExperienceResponse {
	return ExperienceResponse{
		ID:          experience.ID,
		Title:       experience.Title,
		Company:     experience.Company,
		Location:    experience.Location,
		Type:        string(experience.Type),
		StartDate:   experience.StartDate,
		EndDate:     experience.EndDate,
		Description: experience.Description,
		CreatedAt:   experience.CreatedAt,
		UpdatedAt:   experience.UpdatedAt,
	}
}

// ToExperienceResponseList converts a slice of models.Experience to ExperienceResponse
func ToExperienceResponseList(experiences []models.Experience) []ExperienceResponse {
	responses := make([]ExperienceResponse, len(experiences))
	for i, exp := range experiences {
		responses[i] = ToExperienceResponse(&exp)
	}
	return responses
}

// ToExperience converts ExperienceRequest to models.Experience
func (req *ExperienceRequest) ToExperience() *models.Experience {
	return &models.Experience{
		Title:       req.Title,
		Company:     req.Company,
		Location:    req.Location,
		Type:        models.WorkType(req.Type),
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Description: req.Description,
	}
}
