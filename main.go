package main

import (
	"job-tracker-api/config"
	"job-tracker-api/middleware"
	"job-tracker-api/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	config.ConnectDB()

	// Initialize router
	r := gin.Default()

	// Use CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Setup routes
	routes.SetupJobRoutes(r)
	routes.SetupAuthRoutes(r)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	r.Run(":" + port)
}