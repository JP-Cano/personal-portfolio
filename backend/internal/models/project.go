package models

import "time"

type Project struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Description string     `json:"description" gorm:"type:text"`
	URL         string     `json:"url" gorm:"type:varchar(500)"`
	StartDate   time.Time  `json:"start_date" gorm:"type:date;not null"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
