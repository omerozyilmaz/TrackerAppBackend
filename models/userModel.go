package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Username      string         `json:"username" binding:"required" gorm:"unique"`
	Email         string         `json:"email" binding:"required,email" gorm:"unique"`
	Password      string         `json:"password" binding:"required"`
	FirstName     string         `json:"first_name"`
	LastName      string         `json:"last_name"`
	ProfilePicture string         `json:"profile_picture"`
	LinkedInToken string         `json:"linkedIn_token"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
} 