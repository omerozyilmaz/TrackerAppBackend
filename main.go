package main

import (
	"bytes"
	"encoding/json"
	"job-tracker-api/config"
	"job-tracker-api/middleware"
	"job-tracker-api/routes"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	clientID     = os.Getenv("LINKEDIN_CLIENT_ID")
	clientSecret = os.Getenv("LINKEDIN_CLIENT_SECRET")
	redirectURI  = os.Getenv("LINKEDIN_REDIRECT_URI")
)

func main() {
	// Initialize database
	config.ConnectDB()

	// Initialize router
	r := gin.Default()
	
	// Disable automatic redirects
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	// Set trusted proxies for production
	if os.Getenv("GIN_MODE") == "release" {
		r.SetTrustedProxies([]string{"127.0.0.1"})
	} else {
		r.SetTrustedProxies(nil)
	}

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

func getAccessToken(code string) (string, error) {
	tokenURL := "https://www.linkedin.com/oauth/v2/accessToken"
	payload := []byte("grant_type=authorization_code&code=" + code +
		"&redirect_uri=" + redirectURI +
		"&client_id=" + clientID +
		"&client_secret=" + clientSecret)

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result["access_token"].(string), nil
}

func getLinkedInProfile(token string) (interface{}, error) {
	profileURL := "https://api.linkedin.com/v2/me"
	req, _ := http.NewRequest("GET", profileURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profile map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&profile)

	email, err := getLinkedInEmail(token)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"firstName": profile["localizedFirstName"],
		"lastName":  profile["localizedLastName"],
		"email":     email,
		"profilePicture": profile["profilePicture"],
	}, nil
}

func getLinkedInEmail(token string) (string, error) {
	emailURL := "https://api.linkedin.com/v2/emailAddress?q=members&projection=(elements*(handle~))"
	req, _ := http.NewRequest("GET", emailURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emailData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&emailData)

	return emailData["elements"].([]interface{})[0].(map[string]interface{})["handle~"].(map[string]interface{})["emailAddress"].(string), nil
}