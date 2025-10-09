package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Description string     `json:"description" gorm:"type:text"`
	URL         string     `json:"url" gorm:"type:varchar(500)"`
	StartDate   time.Time  `json:"start_date" gorm:"type:date;not null"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}
