package dto

import (
	"time"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/models"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/utils"
)

type ProjectResponse struct {
	ID           uint       `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	URL          string     `json:"url,omitempty"`
	StartDate    time.Time  `json:"startDate"`
	EndDate      *time.Time `json:"endDate,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	Technologies string     `json:"technologies"`
}

type ProjectRequest struct {
	Name         string `json:"name" binding:"required" validate:"required,min=1,max=255"`
	Description  string `json:"description" binding:"required" validate:"required,min=1,max=500"`
	URL          string `json:"url,omitempty" validate:"omitempty,url,max=500"`
	StartDate    string `json:"start_date" binding:"required" validate:"required,date_format"`
	EndDate      string `json:"end_date,omitempty" validate:"omitempty,date_format,after_start_date_str"`
	Technologies string `json:"technologies" validate:"omitempty,max=500"`
}

type UpdateProjectRequest struct {
	Name         *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description  *string `json:"description,omitempty" validate:"omitempty,min=1,max=500"`
	URL          *string `json:"url,omitempty" validate:"omitempty,url,max=500"`
	StartDate    *string `json:"start_date,omitempty" validate:"omitempty,date_format"`
	EndDate      *string `json:"end_date,omitempty" validate:"omitempty,date_format"`
	Technologies *string `json:"technologies,omitempty" validate:"omitempty,max=500"`
}

func ToProjectResponse(project *models.Project) ProjectResponse {
	var endDate *time.Time
	if project.EndDate != nil {
		endDate = &project.EndDate.Time
	}

	return ProjectResponse{
		ID:           project.ID,
		Name:         project.Name,
		Description:  project.Description,
		URL:          project.URL,
		StartDate:    project.StartDate.Time,
		EndDate:      endDate,
		CreatedAt:    project.CreatedAt,
		UpdatedAt:    project.UpdatedAt,
		Technologies: project.Technologies,
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
	startDate, err := utils.ParseToDate(req.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := utils.ParseToDatePtr(req.EndDate)
	if err != nil {
		return nil, err
	}

	return &models.Project{
		Name:         req.Name,
		Description:  req.Description,
		URL:          req.URL,
		StartDate:    startDate,
		EndDate:      endDate,
		Technologies: req.Technologies,
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
		startDate, err := utils.ParseToDate(*req.StartDate)
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
			endDate, err := utils.ParseToDatePtr(*req.EndDate)
			if err != nil {
				return nil, err
			}
			updates["end_date"] = endDate
		}
	}

	if req.Technologies != nil {
		updates["technologies"] = *req.Technologies
	}

	return updates, nil
}
