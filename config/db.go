package config

import (
	"fmt"
	"job-tracker-api/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	
	// GORM logger yapılandırması
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Veritabanı yapılandırma seçenekleri
	config := &gorm.Config{
		Logger: newLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}
	
	// Railway'de DATABASE_URL ortam değişkeni otomatik olarak sağlanır
	dbURL := os.Getenv("DATABASE_URL")
	
	if dbURL != "" {
		log.Println("Connecting to database using URL...")
		DB, err = gorm.Open(postgres.Open(dbURL+"&sslmode=require"), config)
	} else {
		log.Println("Connecting to database using individual credentials...")
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require timezone=UTC",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
		DB, err = gorm.Open(postgres.Open(dsn), config)
	}
	
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Bağlantı havuzu ayarları
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// Bağlantı havuzu yapılandırması
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Bağlantıyı test et
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