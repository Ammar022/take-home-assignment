package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger middleware logs HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		status := c.Writer.Status()

		// Log request details (in a real app, use a proper logger)
		log.Printf("[%s] %s %d %v", method, path, status, latency)

		// For this example, use gin's built-in logger
		c.Set("latency", latency)
	}
}
