// cmd/api/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"take-home-assignment/internal/api"
	"take-home-assignment/internal/config"
	"take-home-assignment/internal/repo"
	"take-home-assignment/internal/service"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to MongoDB
	db, err := repo.NewMongoDBConnection(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Initialize repositories
	linkRepo := repo.NewLinkRepository(db)
	visitRepo := repo.NewVisitRepository(db)

	// Initialize services
	linkService := service.NewLinkService(linkRepo)
	visitService := service.NewVisitService(visitRepo, linkRepo)
	
	// Start background cleanup worker
	cleanupCtx, cleanupCancel := context.WithCancel(context.Background())
	cleanupService := service.NewCleanupService(linkRepo)
	go cleanupService.StartPeriodicCleanup(cleanupCtx, time.Hour*24)

	// Initialize HTTP router
	router := api.SetupRouter(linkService, visitService)

	// Configure HTTP server
	server := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server at %s", cfg.Server.Address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Cancel background workers
	cleanupCancel()

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}