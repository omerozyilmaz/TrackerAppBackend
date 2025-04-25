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

func ConnectDB() {
	var err error
	
	// Railway'de DATABASE_URL ortam değişkeni otomatik olarak sağlanır
	dbURL := os.Getenv("DATABASE_URL")
	
	// Eğer DATABASE_URL varsa, doğrudan onu kullan
	if dbURL != "" {
		DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	} else {
		// Yoksa, ayrı ayrı değişkenleri kullan (yerel geliştirme için)
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}
	
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create or update the tables
	err = DB.AutoMigrate(
		&models.User{},
		&models.Job{},
		&models.Profile{},
		&models.Skill{},
		&models.Education{},
		&models.Experience{},
		&models.Project{},
		&models.LinkedInAuth{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}
}

// Yeni kullanıcı oluşturulduğunda çağrılacak fonksiyon
func CreateUserTables(userID uint) error {
	// Create job table for user
	if err := models.CreateJobTableForUser(DB, userID); err != nil {
		return err
	}

	// Create profile for user
	profile := models.Profile{
		UserID: userID,
	}
	if err := DB.Create(&profile).Error; err != nil {
		return err
	}

	return nil
} 