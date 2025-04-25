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
	
	if dbURL != "" {
		log.Println("Connecting to database using URL...")
		DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	} else {
		log.Println("Connecting to database using individual credentials...")
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
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

	// Bağlantıyı test et
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Successfully connected to database")

	// Veritabanını migrate et
	log.Println("Starting database migration...")
	err = DB.AutoMigrate(&models.User{}, &models.Job{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}
	log.Println("Database migration completed successfully")
}

// Yeni kullanıcı oluşturulduğunda çağrılacak fonksiyon
func CreateUserTables(userID uint) error {
	log.Printf("Creating tables for user ID: %d\n", userID)
	err := models.CreateJobTableForUser(DB, userID)
	if err != nil {
		log.Printf("Error creating user tables: %v\n", err)
		return err
	}
	log.Println("User tables created successfully")
	return nil
} 