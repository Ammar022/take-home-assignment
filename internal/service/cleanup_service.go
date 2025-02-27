package service

import (
	"context"
	"take-home-assignment/internal/repo"
	"time"
)

// CleanupService handles background cleanup tasks
type CleanupService struct {
	linkRepo *repo.LinkRepository
}

// NewCleanupService creates a new cleanup service
func NewCleanupService(linkRepo *repo.LinkRepository) *CleanupService {
	return &CleanupService{
		linkRepo: linkRepo,
	}
}

// StartPeriodicCleanup starts a background goroutine to clean up expired links
func (s *CleanupService) StartPeriodicCleanup(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanupExpiredLinks(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// cleanupExpiredLinks removes expired links from the database
func (s *CleanupService) cleanupExpiredLinks(ctx context.Context) {
	// Create a new context with timeout for this operation
	cleanupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	count, err := s.linkRepo.DeleteExpired(cleanupCtx)
	if err != nil {
		// Log error (implement proper logging)
		return
	}

	if count > 0 {
		// Log cleanup results (implement proper logging)
	}
}