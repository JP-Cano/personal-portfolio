package dto

import (
	"time"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/models"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/utils"
)

// ExperienceResponse represents the API response for an experience
type ExperienceResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Company     string     `json:"company"`
	URL         *string    `json:"url,omitempty"`
	Location    string     `json:"location"`
	Type        string     `json:"type"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate,omitempty"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// ExperienceRequest represents the API request for creating an experience
type ExperienceRequest struct {
	Title       string  `json:"title" binding:"required" validate:"required,min=1,max=255"`
	Company     string  `json:"company" binding:"required" validate:"required,min=1,max=255"`
	URL         *string `json:"url,omitempty" validate:"omitempty,url,max=500"`
	Location    string  `json:"location" validate:"omitempty,max=255"`
	Type        string  `json:"type" binding:"required" validate:"required,oneof=Remote 'On Site' Hybrid"`
	StartDate   string  `json:"start_date" binding:"required" validate:"required,date_format"`
	EndDate     string  `json:"end_date,omitempty" validate:"omitempty,date_format,after_start_date_str"`
	Description string  `json:"description" validate:"omitempty"`
}

// UpdateExperienceRequest represents the API request for updating an experience (all fields optional)
type UpdateExperienceRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Company     *string `json:"company,omitempty" validate:"omitempty,min=1,max=255"`
	URL         *string `json:"url,omitempty" validate:"omitempty,url,max=500"`
	Location    *string `json:"location,omitempty" validate:"omitempty,max=255"`
	Type        *string `json:"type,omitempty" validate:"omitempty,oneof=Remote 'On Site' Hybrid"`
	StartDate   *string `json:"start_date,omitempty" validate:"omitempty,date_format"`
	EndDate     *string `json:"end_date,omitempty" validate:"omitempty,date_format"`
	Description *string `json:"description,omitempty"`
}

// ToExperienceResponse converts a models.Experience to ExperienceResponse
func ToExperienceResponse(experience *models.Experience) ExperienceResponse {
	return ExperienceResponse{
		ID:          experience.ID,
		Title:       experience.Title,
		URL:         experience.URL,
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
func (req *ExperienceRequest) ToExperience() (*models.Experience, error) {
	startDate, err := utils.ParseDate(req.StartDate)
	if err != nil {
		return nil, err
	}

	var endDate *time.Time
	if req.EndDate != "" {
		parsed, err := utils.ParseDateToPtr(req.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = parsed
	}

	return &models.Experience{
		Title:       req.Title,
		Company:     req.Company,
		URL:         req.URL,
		Location:    req.Location,
		Type:        models.WorkType(req.Type),
		StartDate:   startDate,
		EndDate:     endDate,
		Description: req.Description,
	}, nil
}

// ToUpdateMap converts UpdateExperienceRequest to a map for partial updates
func (req *UpdateExperienceRequest) ToUpdateMap() (map[string]interface{}, error) {
	updates := make(map[string]interface{})

	if req.Title != nil {
		updates["title"] = *req.Title
	}

	if req.Company != nil {
		updates["company"] = *req.Company
	}

	if req.URL != nil {
		updates["url"] = req.URL
	}

	if req.Location != nil {
		updates["location"] = *req.Location
	}

	if req.Type != nil {
		updates["type"] = models.WorkType(*req.Type)
	}

	if req.StartDate != nil {
		startDate, err := utils.ParseDate(*req.StartDate)
		if err != nil {
			return nil, err
		}
		updates["start_date"] = startDate
	}

	if req.EndDate != nil {
		if *req.EndDate == "" {
			// Allow clearing the end date
			updates["end_date"] = nil
		} else {
			endDate, err := utils.ParseDateToPtr(*req.EndDate)
			if err != nil {
				return nil, err
			}
			updates["end_date"] = endDate
		}
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	return updates, nil
}
