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