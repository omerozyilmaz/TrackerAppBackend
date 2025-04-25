package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimitMiddleware creates a rate limiter middleware
func RateLimitMiddleware() gin.HandlerFunc {
	// Create a rate limiter with 100 requests per minute
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}
	store := memory.NewStore()
	limiter := limiter.New(store, rate)

	return func(c *gin.Context) {
		context := c.Request.Context()
		ip := c.ClientIP()

		limiterCtx, err := limiter.Get(context, ip)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if limiterCtx.Reached {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}

// SecurityHeadersMiddleware adds security-related headers to responses
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent XSS attacks
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		
		// Enable HSTS
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Feature Policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		c.Next()
	}
}

// InputValidationMiddleware validates and sanitizes input
func InputValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add input validation logic here
		// For example, validate request body against a schema
		// or sanitize input data
		
		c.Next()
	}
} 