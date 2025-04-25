package routes

import (
	"job-tracker-api/controllers"
	"job-tracker-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupProfileRoutes(router *gin.Engine) {
	// Support new URL structure (without /api prefix)
	profile := router.Group("/profile")
	profile.Use(middleware.AuthMiddleware())
	{
		profile.GET("", controllers.GetProfile)
		profile.GET("/", controllers.GetProfile)
		profile.POST("", controllers.UpdateProfile)
		profile.POST("/", controllers.UpdateProfile)
		profile.PUT("", controllers.UpdateProfile)
		profile.PUT("/", controllers.UpdateProfile)
	}

	// Also support old URL structure (with /api prefix)
	apiProfile := router.Group("/api/profile")
	apiProfile.Use(middleware.AuthMiddleware())
	{
		apiProfile.GET("", controllers.GetProfile)
		apiProfile.GET("/", controllers.GetProfile)
		apiProfile.POST("", controllers.UpdateProfile)
		apiProfile.POST("/", controllers.UpdateProfile)
		apiProfile.PUT("", controllers.UpdateProfile)
		apiProfile.PUT("/", controllers.UpdateProfile)
	}
} 