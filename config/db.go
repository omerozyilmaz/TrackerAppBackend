package config

import (
	"fmt"
	"job-tracker-api/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Create or update the jobs table
	err = DB.AutoMigrate(&models.User{}, &models.Job{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	return DB, nil
}

// Yeni kullanıcı oluşturulduğunda çağrılacak fonksiyon
func CreateUserTables(userID uint) error {
	return models.CreateJobTableForUser(DB, userID)
} 