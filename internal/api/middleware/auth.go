package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a mock authentication middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		
		// Check if the header is empty or doesn't start with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		
		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		
		// In a real app, validate the token here
		// For this mock, we'll just check if it's not empty
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		
		// Mock user ID - in a real app, this would be extracted from the token
		// For simplicity, we'll just use a fixed value
		c.Set("userId", "ammar-123")
		
		c.Next()
	}
}
