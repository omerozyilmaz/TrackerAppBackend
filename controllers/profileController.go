package controllers

import (
	"job-tracker-api/config"
	"job-tracker-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var profile models.Profile
	if err := config.DB.Preload("Skills").
		Preload("Education").
		Preload("Experience").
		Preload("Experience.Skills").
		Preload("Projects").
		Preload("Projects.Skills").
		Where("user_id = ?", userID).
		First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var profile models.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the profile belongs to the authenticated user
	profile.UserID = userID

	// Check if profile exists
	var existingProfile models.Profile
	if err := config.DB.Where("user_id = ?", userID).First(&existingProfile).Error; err != nil {
		// If profile doesn't exist, create new one
		if err := config.DB.Create(&profile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
			return
		}
	} else {
		// If profile exists, update it
		profile.ID = existingProfile.ID
		if err := config.DB.Save(&profile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
			return
		}
	}

	c.JSON(http.StatusOK, profile)
}

func UpdateProfilePicture(c *gin.Context) {
	userID := c.GetUint("user_id")

	// FormData'dan dosyayı al
	file, err := c.FormFile("profilePicture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	// Dosyayı kaydetme işlemi (örneğin, belirli bir dizine)
	if err := c.SaveUploadedFile(file, "./uploads/"+file.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Kullanıcı profilini güncelle (örneğin, veritabanında profil resmini güncelle)
	var profile models.Profile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	profile.ProfilePicture = file.Filename // veya dosya yolunu güncelleyin
	if err := config.DB.Save(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile picture"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func UpdateEducation(c *gin.Context) {
	userID := c.GetUint("user_id")

	var educationData []models.Education
	if err := c.ShouldBindJSON(&educationData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Eğitim bilgilerini güncelle
	for _, edu := range educationData {
		edu.ProfileID = userID // Kullanıcı ID'sini ayarla
		if err := config.DB.Save(&edu).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update education"})
			return
		}
	}

	c.JSON(http.StatusOK, educationData)
}

func UpdateExperience(c *gin.Context) {
	userID := c.GetUint("user_id")

	var experienceData []models.Experience
	if err := c.ShouldBindJSON(&experienceData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Deneyim bilgilerini güncelle
	for _, exp := range experienceData {
		exp.ProfileID = userID // Kullanıcı ID'sini ayarla
		if err := config.DB.Save(&exp).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update experience"})
			return
		}
	}

	c.JSON(http.StatusOK, experienceData)
} 