package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"type:varchar(255);not null"`
	Password string `json:"-" gorm:"type:varchar(255);not null"`
}

type Session struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type AuthResponse struct {
	User    UserResponse `json:"user"`
	Message string       `json:"message"`
}
