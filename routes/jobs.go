package routes

import (
	"job-tracker-api/controllers"
	"job-tracker-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupJobRoutes(router *gin.Engine) {
	jobs := router.Group("/jobs")
	
	// Use the centralized CORS middleware
	jobs.Use(middleware.CORSMiddleware())
	jobs.Use(middleware.AuthMiddleware())
	{
		jobs.GET("/", controllers.GetJobs)
		jobs.POST("/", controllers.CreateJob)
		jobs.GET("/:id", controllers.GetJob)
		jobs.PUT("/:id", controllers.UpdateJob)
		jobs.DELETE("/:id", controllers.DeleteJob)
		
		// Sadece status güncellemesi için özel endpoint (opsiyonel)
		jobs.PATCH("/:id/status", controllers.UpdateJobStatus)
		// Add new endpoint for moving jobs
		jobs.PATCH("/:id/move", controllers.MoveJob)
	}
} 