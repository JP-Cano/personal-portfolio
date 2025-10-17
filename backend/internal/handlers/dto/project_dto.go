package dto

import (
	"time"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/models"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/utils"
)

type ProjectResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	URL         string     `json:"url,omitempty"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type ProjectRequest struct {
	Name        string `json:"name" binding:"required" validate:"required,min=1,max=255"`
	Description string `json:"description" binding:"required" validate:"required,min=1,max=500"`
	URL         string `json:"url,omitempty" validate:"omitempty,url,max=500"`
	StartDate   string `json:"start_date" binding:"required" validate:"required,date_format"`
	EndDate     string `json:"end_date,omitempty" validate:"omitempty,date_format,after_start_date_str"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=1,max=500"`
	URL         *string `json:"url,omitempty" validate:"omitempty,url,max=500"`
	StartDate   *string `json:"start_date,omitempty" validate:"omitempty,date_format"`
	EndDate     *string `json:"end_date,omitempty" validate:"omitempty,date_format"`
}

func ToProjectResponse(project *models.Project) ProjectResponse {
	return ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		URL:         project.URL,
		StartDate:   project.StartDate,
		EndDate:     project.EndDate,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}
}

func ToProjectResponseList(projects []models.Project) []ProjectResponse {
	responses := make([]ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = ToProjectResponse(&project)
	}
	return responses
}

func (req *ProjectRequest) ToProject() (*models.Project, error) {
	// Parse start date
	startDate, err := utils.ParseDate(req.StartDate)
	if err != nil {
		return nil, err
	}

	// Parse end date if provided
	var endDate *time.Time
	if req.EndDate != "" {
		parsed, err := utils.ParseDateToPtr(req.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = parsed
	}

	return &models.Project{
		Name:        req.Name,
		Description: req.Description,
		URL:         req.URL,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

// ToUpdateMap converts UpdateProjectRequest to a map for partial updates
func (req *UpdateProjectRequest) ToUpdateMap() (map[string]interface{}, error) {
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if req.URL != nil {
		updates["url"] = *req.URL
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

	return updates, nil
}
