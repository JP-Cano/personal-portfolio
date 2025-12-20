package models

import (
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/utils"

	"gorm.io/gorm"
)

type WorkType string

const (
	Remote WorkType = "Remote"
	OnSite WorkType = "On Site"
	Hybrid WorkType = "Hybrid"
)

type Experience struct {
	gorm.Model
	Title       string      `json:"title" gorm:"type:varchar(255);not null"`
	Company     string      `json:"company" gorm:"type:varchar(255);not null"`
	URL         *string     `json:"url,omitempty" gorm:"type:varchar(500)"`
	Location    string      `json:"location" gorm:"type:varchar(255)"`
	Type        WorkType    `json:"type" gorm:"type:varchar(50);not null"`
	StartDate   utils.Date  `json:"start_date" gorm:"type:date;not null"`
	EndDate     *utils.Date `json:"end_date,omitempty"`
	Description string      `json:"description" gorm:"type:text"`
}
