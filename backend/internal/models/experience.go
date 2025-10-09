package models

import (
	"time"

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
	Title       string     `json:"title" gorm:"type:varchar(255);not null"`
	Company     string     `json:"company" gorm:"type:varchar(255);not null"`
	Location    string     `json:"location" gorm:"type:varchar(255)"`
	Type        WorkType   `json:"type" gorm:"type:varchar(50);not null"`
	StartDate   time.Time  `json:"start_date" gorm:"type:date;not null"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Description string     `json:"description" gorm:"type:text"`
}
