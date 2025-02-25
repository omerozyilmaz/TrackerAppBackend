package controllers

import (
	"job-tracker-api/config"
	"job-tracker-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetJobs(c *gin.Context) {
	var jobs []models.Job
	userID := c.GetUint("user_id")

	if err := config.DB.Where("user_id = ?", userID).Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

func CreateJob(c *gin.Context) {
	var job models.Job
	userID := c.GetUint("user_id")

	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job.UserID = userID
	job.AddedTime = time.Now()

	if err := config.DB.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	c.JSON(http.StatusCreated, job)
}

func UpdateJob(c *gin.Context) {
	var job models.Job
	jobID := c.Param("id")
	userID := c.GetUint("user_id")

	// Check if job exists and belongs to user
	if err := config.DB.Where("id = ? AND user_id = ?", jobID, userID).First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	// Update job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func DeleteJob(c *gin.Context) {
	jobID := c.Param("id")
	userID := c.GetUint("user_id")

	if err := config.DB.Where("id = ? AND user_id = ?", jobID, userID).Delete(&models.Job{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully"})
}

func MoveJob(c *gin.Context) {
	type MoveRequest struct {
		Column string `json:"column" binding:"required,oneof=wishlist applied interview"`
	}

	var moveReq MoveRequest
	var job models.Job
	jobID := c.Param("id")
	userID := c.GetUint("user_id")

	if err := c.ShouldBindJSON(&moveReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if job exists and belongs to user
	if err := config.DB.Where("id = ? AND user_id = ?", jobID, userID).First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	// Update column
	if err := config.DB.Model(&job).Update("column", moveReq.Column).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to move job"})
		return
	}

	c.JSON(http.StatusOK, job)
} 