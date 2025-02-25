package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Job struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	Title       string         `json:"title" gorm:"not null" binding:"required"`
	Company     string         `json:"company" gorm:"not null" binding:"required"`
	Status      string         `json:"status" gorm:"not null"`
	Description string         `json:"description"`
	Location    string         `json:"location"`
	JobURL      string         `json:"jobUrl"`
	Column      string         `json:"column" binding:"required,oneof=wishlist applied interview"`
	AddedTime   time.Time      `json:"addedTime"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func CreateJobTableForUser(db *gorm.DB, userID uint) error {
	tableName := fmt.Sprintf("jobs_user_%d", userID)
	// Her kullanıcı için ayrı tablo oluştur
	return db.Table(tableName).AutoMigrate(&Job{})
}

// Kullanıcıya özel tablo adını döndüren yardımcı fonksiyon
func GetJobTableName(userID uint) string {
	return fmt.Sprintf("jobs_user_%d", userID)
}

// Job oluşturma metodu
func (j *Job) Create(db *gorm.DB, userID uint) error {
	tableName := GetJobTableName(userID)
	return db.Table(tableName).Create(j).Error
}

// Jobları listeleme metodu
func GetJobs(db *gorm.DB, userID uint) ([]Job, error) {
	var jobs []Job
	tableName := GetJobTableName(userID)
	err := db.Table(tableName).Find(&jobs).Error
	return jobs, err
}

// Job güncelleme metodu
func (j *Job) Update(db *gorm.DB, userID uint) error {
	tableName := GetJobTableName(userID)
	return db.Table(tableName).Save(j).Error
}

// Job silme metodu
func DeleteJob(db *gorm.DB, userID uint, jobID uint) error {
	tableName := GetJobTableName(userID)
	return db.Table(tableName).Delete(&Job{}, jobID).Error
} 