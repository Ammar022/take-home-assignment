package api

import (
	"take-home-assignment/internal/api/handlers"
	"take-home-assignment/internal/api/middleware"
	"take-home-assignment/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// SetupRouter configures the Gin router
func SetupRouter(linkService *service.LinkService, visitService *service.VisitService) *gin.Engine {
	// Create router
	r := gin.Default()
	
	// Apply global middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middleware.Logger())
	
	// Create rate limiter for visit endpoint
	visitLimiter := middleware.NewRateLimiter(rate.Limit(1000), 200)
	
	// Create handlers
	linkHandler := handlers.NewLinkHandler(linkService)
	visitHandler := handlers.NewVisitHandler(visitService)
	
	// Public routes
	r.GET("/visit/:id", visitLimiter.Middleware(), visitHandler.RecordVisit)
	
	// API routes (require authentication)
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		links := api.Group("/links")
		{
			links.GET("", linkHandler.GetAll)
			links.POST("", linkHandler.Create)
			links.GET("/:id", linkHandler.GetByID)
			links.PUT("/:id", linkHandler.Update)
			links.DELETE("/:id", linkHandler.Delete)
			
			// Visits for a specific link
			links.GET("/:id/visits", visitHandler.GetVisitsForLink)
		}
	}
	
	return r
}