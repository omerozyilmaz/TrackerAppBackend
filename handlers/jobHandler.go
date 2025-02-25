package handlers

import (
	"job-tracker-api/config"
	"job-tracker-api/models"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateJob(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var job models.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	job.UserID = userID
	job.AddedTime = time.Now()

	if err := config.DB.Create(&job).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create job"})
		return
	}

	c.JSON(201, job)
}

func GetJobs(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var jobs []models.Job
	if err := config.DB.Where("user_id = ?", userID).Find(&jobs).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch jobs"})
		return
	}

	c.JSON(200, jobs)
}

func UpdateJob(c *gin.Context) {
	userID := c.GetUint("user_id")
	jobID := c.Param("id")
	
	var job models.Job
	if err := config.DB.Where("id = ? AND user_id = ?", jobID, userID).First(&job).Error; err != nil {
		c.JSON(404, gin.H{"error": "Job not found"})
		return
	}

	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&job).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update job"})
		return
	}

	c.JSON(200, job)
}

func DeleteJob(c *gin.Context) {
	userID := c.GetUint("user_id")
	jobID := c.Param("id")

	if err := config.DB.Where("id = ? AND user_id = ?", jobID, userID).Delete(&models.Job{}).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete job"})
		return
	}

	c.JSON(200, gin.H{"message": "Job deleted successfully"})
} 