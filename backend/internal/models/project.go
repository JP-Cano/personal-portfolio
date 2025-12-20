package models

import (
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/utils"

	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name         string      `json:"name" gorm:"type:varchar(255);not null"`
	Description  string      `json:"description" gorm:"type:text"`
	URL          string      `json:"url" gorm:"type:varchar(500)"`
	StartDate    utils.Date  `json:"start_date" gorm:"type:date;not null"`
	EndDate      *utils.Date `json:"end_date,omitempty"`
	Technologies string      `json:"technologies" gorm:"type:text"`
}
