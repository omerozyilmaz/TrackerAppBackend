package models

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	Headline      string         `json:"headline"`
	Summary       string         `json:"summary"`
	Location      string         `json:"location"`
	Website       string         `json:"website"`
	Skills        []Skill        `json:"skills" gorm:"many2many:profile_skills;"`
	Education     []Education    `json:"education"`
	Experience    []Experience   `json:"experience"`
	Projects      []Project      `json:"projects"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type Skill struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"unique;not null"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Education struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	ProfileID    uint           `json:"profile_id" gorm:"not null"`
	School       string         `json:"school" gorm:"not null"`
	Degree       string         `json:"degree"`
	FieldOfStudy string         `json:"field_of_study"`
	StartDate    time.Time      `json:"start_date"`
	EndDate      *time.Time     `json:"end_date"`
	Grade        string         `json:"grade"`
	Description  string         `json:"description"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type Experience struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ProfileID   uint           `json:"profile_id" gorm:"not null"`
	Title       string         `json:"title" gorm:"not null"`
	Company     string         `json:"company" gorm:"not null"`
	Location    string         `json:"location"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
	Description string         `json:"description"`
	Skills      []Skill        `json:"skills" gorm:"many2many:experience_skills;"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Project struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ProfileID   uint           `json:"profile_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
	URL         string         `json:"url"`
	Skills      []Skill        `json:"skills" gorm:"many2many:project_skills;"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type LinkedInAuth struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	AccessToken   string         `json:"access_token"`
	RefreshToken  string         `json:"refresh_token"`
	ExpiresIn     int64          `json:"expires_in"`
	TokenType     string         `json:"token_type"`
	Scope         string         `json:"scope"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
} 