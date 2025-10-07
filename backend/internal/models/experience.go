package models

import "time"

type WorkType string

const (
	Remote WorkType = "Remote"
	OnSite WorkType = "On Site"
	Hybrid WorkType = "Hybrid"
)

type Experience struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string     `json:"title" gorm:"type:varchar(255);not null"`
	Company     string     `json:"company" gorm:"type:varchar(255);not null"`
	Location    string     `json:"location" gorm:"type:varchar(255)"`
	Type        WorkType   `json:"type" gorm:"type:varchar(50);not null"`
	StartDate   time.Time  `json:"start_date" gorm:"type:date;not null"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Description string     `json:"description" gorm:"type:text"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
