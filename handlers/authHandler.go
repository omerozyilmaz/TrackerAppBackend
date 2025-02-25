package handlers

import (
	"job-tracker-api/config"
	"job-tracker-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user models.User

	// Parse the user data from the request
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the user to the database
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Kullanıcı başarıyla oluşturulduktan sonra
	if err := config.CreateUserTables(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user tables"})
		return
	}

	c.JSON(http.StatusCreated, user)
} 