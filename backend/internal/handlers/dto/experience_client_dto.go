package dto

import (
	"time"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/models"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/utils"
)

type ExperienceClientResponse struct {
	ID               uint       `json:"id"`
	ExperienceID     uint       `json:"experienceId"`
	Name             string     `json:"name"`
	URL              *string    `json:"url,omitempty"`
	StartDate        time.Time  `json:"startDate"`
	EndDate          *time.Time `json:"endDate,omitempty"`
	Description      string     `json:"description,omitempty"`
	Achievements     []string   `json:"achievements"`
	Responsibilities []string   `json:"responsibilities"`
	Technologies     []string   `json:"technologies"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

type ExperienceClientRequest struct {
	Name             string   `json:"name" binding:"required" validate:"required,min=1,max=255"`
	URL              *string  `json:"url" validate:"omitempty,url"`
	StartDate        string   `json:"start_date" binding:"required" validate:"required,date_format"`
	EndDate          *string  `json:"end_date" validate:"omitempty,date_format,after_start_date_str=StartDate"`
	Description      string   `json:"description"`
	Achievements     []string `json:"achievements"`
	Responsibilities []string `json:"responsibilities"`
	Technologies     []string `json:"technologies"`
}

type UpdateExperienceClientRequest struct {
	Name             *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	URL              *string  `json:"url,omitempty" validate:"omitempty,url"`
	StartDate        *string  `json:"start_date,omitempty" validate:"omitempty,date_format"`
	EndDate          *string  `json:"end_date,omitempty" validate:"omitempty,date_format,after_start_date_str=StartDate"`
	Description      *string  `json:"description,omitempty"`
	Achievements     []string `json:"achievements,omitempty"`
	Responsibilities []string `json:"responsibilities,omitempty"`
	Technologies     []string `json:"technologies,omitempty"`
}

func ToExperienceClientResponse(client *models.ExperienceClient) ExperienceClientResponse {
	resp := ExperienceClientResponse{
		ID:               client.ID,
		ExperienceID:     client.ExperienceID,
		Name:             client.Name,
		URL:              client.URL,
		StartDate:        client.StartDate.Time,
		Description:      client.Description,
		Achievements:     []string(client.Achievements),
		Responsibilities: []string(client.Responsibilities),
		Technologies:     []string(client.Technologies),
		CreatedAt:        client.CreatedAt,
		UpdatedAt:        client.UpdatedAt,
	}
	if client.EndDate != nil {
		resp.EndDate = &client.EndDate.Time
	}
	return resp
}

func ToExperienceClientResponseList(clients []models.ExperienceClient) []ExperienceClientResponse {
	responses := make([]ExperienceClientResponse, len(clients))
	for i, client := range clients {
		responses[i] = ToExperienceClientResponse(&client)
	}
	return responses
}

func (req *ExperienceClientRequest) ToExperienceClient() (*models.ExperienceClient, error) {
	startDate, err := utils.ParseToDate(req.StartDate)
	if err != nil {
		return nil, err
	}

	client := &models.ExperienceClient{
		Name:             req.Name,
		URL:              req.URL,
		StartDate:        startDate,
		Description:      req.Description,
		Achievements:     models.JSONStrings(req.Achievements),
		Responsibilities: models.JSONStrings(req.Responsibilities),
		Technologies:     models.JSONStrings(req.Technologies),
	}

	if req.EndDate != nil && *req.EndDate != "" {
		endDate, err := utils.ParseToDatePtr(*req.EndDate)
		if err != nil {
			return nil, err
		}
		client.EndDate = endDate
	}

	return client, nil
}

func (req *UpdateExperienceClientRequest) ToUpdateMap() (map[string]interface{}, error) {
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.URL != nil {
		updates["url"] = *req.URL
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Achievements != nil {
		updates["achievements"] = models.JSONStrings(req.Achievements)
	}
	if req.Responsibilities != nil {
		updates["responsibilities"] = models.JSONStrings(req.Responsibilities)
	}
	if req.Technologies != nil {
		updates["technologies"] = models.JSONStrings(req.Technologies)
	}
	if req.StartDate != nil && *req.StartDate != "" {
		startDate, err := utils.ParseToDate(*req.StartDate)
		if err != nil {
			return nil, err
		}
		updates["start_date"] = startDate
	}
	if req.EndDate != nil {
		if *req.EndDate == "" {
			updates["end_date"] = nil
		} else {
			endDate, err := utils.ParseToDatePtr(*req.EndDate)
			if err != nil {
				return nil, err
			}
			updates["end_date"] = endDate
		}
	}

	return updates, nil
}
