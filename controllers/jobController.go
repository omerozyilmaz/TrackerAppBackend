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
	// Get job ID from URL
	jobID := c.Param("id")
	
	// Get user ID from context (auth middleware tarafından eklenir)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	// Find the job
	var job models.Job
	if err := config.DB.Where("id = ? AND user_id = ?", jobID, userID).First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found or not authorized"})
		return
	}
	
	// Bind the update data
	var input struct {
		Title       string `json:"title"`
		Company     string `json:"company"`
		Location    string `json:"location"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validate status if provided
	if input.Status != "" {
		if !models.IsValidStatus(input.Status) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be one of: wishlist, applied, interview"})
			return
		}
		job.Status = input.Status
	}
	
	// Update other fields if provided
	if input.Title != "" {
		job.Title = input.Title
	}
	if input.Company != "" {
		job.Company = input.Company
	}
	if input.Location != "" {
		job.Location = input.Location
	}
	if input.Description != "" {
		job.Description = input.Description
	}
	
	// Save the updated job
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

// UpdateJobStatus updates only the status of a job
func UpdateJobStatus(c *gin.Context) {
	// Get job ID from URL
	jobID := c.Param("id")
	
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	// Find the job
	var job models.Job
	if err := config.DB.Where("id = ? AND user_id = ?", jobID, userID).First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found or not authorized"})
		return
	}
	
	// Bind the status update
	var input struct {
		Status string `json:"status" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validate status
	if !models.IsValidStatus(input.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be one of: wishlist, applied, interview"})
		return
	}
	
	// Update status
	job.Status = input.Status
	
	// Save the updated job
	if err := config.DB.Save(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job status"})
		return
	}
	
	c.JSON(http.StatusOK, job)
}

// GetJob gets a single job by ID
func GetJob(c *gin.Context) {
	jobID := c.Param("id")
	userID := c.GetUint("user_id")
	
	var job models.Job
	if err := config.DB.Where("id = ? AND user_id = ?", jobID, userID).First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	
	c.JSON(http.StatusOK, job)
} 