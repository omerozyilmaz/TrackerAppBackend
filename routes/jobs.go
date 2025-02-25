package routes

import (
	"job-tracker-api/controllers"
	"job-tracker-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupJobRoutes(r *gin.Engine) {
	jobs := r.Group("jobs")
	jobs.Use(middleware.AuthMiddleware())
	{
		jobs.GET("", controllers.GetJobs)
		jobs.POST("", controllers.CreateJob)
		jobs.PUT("/:id", controllers.UpdateJob)
		jobs.DELETE("/:id", controllers.DeleteJob)
		jobs.PATCH("/:id/move", controllers.MoveJob)
	}
} 