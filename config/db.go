package config

import (
	"job-tracker-api/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate işlemi burada yapılabilir
	// err = db.AutoMigrate(&models.User{}, &models.Job{})
	// if err != nil {
	//     return nil, err
	// }

	return db, nil
}

// Yeni kullanıcı oluşturulduğunda çağrılacak fonksiyon
func CreateUserTables(userID uint) error {
	return models.CreateJobTableForUser(DB, userID)
} 