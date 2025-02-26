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
	}
	
	// Also support old URL structure (with /api prefix)
	apiAuth := router.Group("/api")
	{
		apiAuth.POST("/register", controllers.Register)
		apiAuth.POST("/login", controllers.Login)
	}
} 