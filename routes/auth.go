package routes

import (
	"job-tracker-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine) {
	auth := router.Group("/")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}
} 