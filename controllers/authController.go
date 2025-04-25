package controllers

import (
	"encoding/json"
	"job-tracker-api/config"
	"job-tracker-api/models"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User

	// Kullanıcıdan gelen bilgileri kontrol et
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Eğer kullanıcı LinkedIn ile kayıt oluyorsa ve şifre boşsa, varsayılan bir şifre oluştur
	if user.LinkedInToken != "" && user.Password == "" {
		user.Password = "default_password" // veya başka bir varsayılan değer
	} else if user.Password == "" {
		// Manuel kayıt için şifre zorunlu
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required for manual registration"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Create user
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Kullanıcı başarıyla oluşturulduktan sonra
	if err := config.CreateUserTables(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user tables"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var loginReq models.LoginRequest
	var user models.User

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// LinkedIn'den alınan bilgileri güncelle (sadece LinkedIn token'ı varsa)
	if user.LinkedInToken != "" && (user.FirstName == "" || user.LastName == "") {
		// LinkedIn API'sinden profil bilgilerini al
		profile, err := getLinkedInProfile(user.LinkedInToken)
		if err == nil {
			user.FirstName = profile["firstName"].(string)
			user.LastName = profile["lastName"].(string)
			user.ProfilePicture = profile["profilePicture"].(string)

			// Kullanıcıyı güncelle
			config.DB.Save(&user)
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func LinkedInAuth(c *gin.Context) {
	clientID := os.Getenv("LINKEDIN_CLIENT_ID")
	redirectURI := os.Getenv("LINKEDIN_REDIRECT_URI")
	scope := "r_liteprofile r_emailaddress"

	authURL := "https://www.linkedin.com/oauth/v2/authorization" +
		"?response_type=code" +
		"&client_id=" + clientID +
		"&redirect_uri=" + redirectURI +
		"&scope=" + scope

	c.Redirect(http.StatusFound, authURL)
}

func LinkedInCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is required"})
		return
	}

	// Access token al
	accessToken, err := getAccessToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
		return
	}

	// LinkedIn profil bilgilerini al
	profile, err := getLinkedInProfile(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		return
	}

	// Kullanıcıyı bul veya oluştur
	var user models.User
	if err := config.DB.Where("email = ?", profile["email"]).First(&user).Error; err != nil {
		// Kullanıcı yoksa yeni kullanıcı oluştur
		user = models.User{
			Email:         profile["email"].(string),
			FirstName:     profile["firstName"].(string),
			LastName:      profile["lastName"].(string),
			ProfilePicture: profile["profilePicture"].(string),
			LinkedInToken: accessToken,
		}
		if err := config.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	} else {
		// Kullanıcı varsa LinkedIn bilgilerini güncelle
		user.FirstName = profile["firstName"].(string)
		user.LastName = profile["lastName"].(string)
		user.ProfilePicture = profile["profilePicture"].(string)
		user.LinkedInToken = accessToken
		if err := config.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
	}

	// JWT token oluştur
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func getAccessToken(code string) (string, error) {
	clientID := os.Getenv("LINKEDIN_CLIENT_ID")
	clientSecret := os.Getenv("LINKEDIN_CLIENT_SECRET")
	redirectURI := os.Getenv("LINKEDIN_REDIRECT_URI")

	resp, err := http.PostForm("https://www.linkedin.com/oauth/v2/accessToken", url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirectURI},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result["access_token"].(string), nil
}

func getLinkedInProfile(token string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://api.linkedin.com/v2/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profile map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	// Email adresini al
	email, err := getLinkedInEmail(token)
	if err != nil {
		return nil, err
	}
	profile["email"] = email

	return profile, nil
}

func getLinkedInEmail(token string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.linkedin.com/v2/emailAddress?q=members&projection=(elements*(handle~))", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	elements := result["elements"].([]interface{})
	if len(elements) == 0 {
		return "", nil
	}

	element := elements[0].(map[string]interface{})
	handle := element["handle~"].(map[string]interface{})
	return handle["emailAddress"].(string), nil
} 