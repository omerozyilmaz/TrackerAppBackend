package main

import (
	"job-tracker-api/config"
	"job-tracker-api/middleware"
	"job-tracker-api/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	config.ConnectDB()

	// Initialize router
	r := gin.Default()

	// Use CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Setup routes
	routes.SetupJobRoutes(r)
	routes.SetupAuthRoutes(r)

	// Start server
	r.Run(":8080")
}