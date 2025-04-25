package routes

import (
	"job-tracker-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine) {
	// Support new URL structure (without /api prefix)
	auth := router.Group("/")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.GET("/linkedin/auth", controllers.LinkedInAuth)
		auth.GET("/linkedin/callback", controllers.LinkedInCallback)
	}
	
	// Also support old URL structure (with /api prefix)
	apiAuth := router.Group("/api")
	{
		apiAuth.POST("/register", controllers.Register)
		apiAuth.POST("/login", controllers.Login)
		apiAuth.GET("/linkedin/auth", controllers.LinkedInAuth)
		apiAuth.GET("/linkedin/callback", controllers.LinkedInCallback)
	}
} 