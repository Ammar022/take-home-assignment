package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter is a middleware for limiting request rates
type RateLimiter struct {
	limiter *rate.Limiter
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(limit, burst),
	}
}

// Middleware returns a Gin middleware function for rate limiting
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// In a real application, you would use the client IP or token to have
		// a separate limiter for each client
		
		if !rl.limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

