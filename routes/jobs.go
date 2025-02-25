package routes

import (
	"job-tracker-api/controllers"
	"job-tracker-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupJobRoutes(router *gin.Engine) {
	jobs := router.Group("/api/jobs")
	
	// Add CORS middleware before auth middleware
	jobs.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	
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