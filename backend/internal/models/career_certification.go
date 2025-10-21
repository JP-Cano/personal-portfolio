package models

import (
	"time"

	"gorm.io/gorm"
)

type CareerCertification struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"type:varchar(255);not null" json:"title"`
	Issuer        string         `gorm:"type:varchar(255);not null" json:"issuer"`
	IssueDate     time.Time      `gorm:"not null" json:"issue_date"`
	ExpiryDate    *time.Time     `json:"expiry_date,omitempty"`
	CredentialID  *string        `gorm:"type:varchar(255)" json:"credential_id,omitempty"`
	CredentialURL *string        `gorm:"type:varchar(500)" json:"credential_url,omitempty"`
	FileURL       string         `gorm:"type:varchar(500);not null" json:"file_url"`
	FileName      string         `gorm:"type:varchar(255);not null" json:"file_name"`
	OriginalName  string         `gorm:"type:varchar(255);not null" json:"original_name"`
	FileSize      int64          `gorm:"not null" json:"file_size"`
	MimeType      string         `gorm:"type:varchar(100);not null" json:"mime_type"`
	Description   string         `gorm:"type:text" json:"description,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
