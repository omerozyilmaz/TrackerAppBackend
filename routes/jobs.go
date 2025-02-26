package routes

import (
	"job-tracker-api/controllers"
	"job-tracker-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupJobRoutes(router *gin.Engine) {
	// Support new URL structure (without /api prefix)
	jobs := router.Group("/jobs")
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
	
	// Also support old URL structure (with /api prefix)
	apiJobs := router.Group("/api/jobs")
	apiJobs.Use(middleware.CORSMiddleware())
	apiJobs.Use(middleware.AuthMiddleware())
	{
		apiJobs.GET("/", controllers.GetJobs)
		apiJobs.POST("/", controllers.CreateJob)
		apiJobs.GET("/:id", controllers.GetJob)
		apiJobs.PUT("/:id", controllers.UpdateJob)
		apiJobs.DELETE("/:id", controllers.DeleteJob)
		
		// Sadece status güncellemesi için özel endpoint (opsiyonel)
		apiJobs.PATCH("/:id/status", controllers.UpdateJobStatus)
		// Add new endpoint for moving jobs
		apiJobs.PATCH("/:id/move", controllers.MoveJob)
	}
} 