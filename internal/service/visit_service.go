package service

import (
	"context"
	"errors"
	"take-home-assignment/internal/models"
	"take-home-assignment/internal/repo"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VisitService handles visit business logic
type VisitService struct {
	visitRepo *repo.VisitRepository
	linkRepo  *repo.LinkRepository
}

// NewVisitService creates a new visit service
func NewVisitService(visitRepo *repo.VisitRepository, linkRepo *repo.LinkRepository) *VisitService {
	return &VisitService{
		visitRepo: visitRepo,
		linkRepo:  linkRepo,
	}
}

// RecordVisit records a new visit and increments link click count
func (s *VisitService) RecordVisit(ctx context.Context, linkID string, userAgent, ip, referrer string) (models.Link, error) {
	objectID, err := primitive.ObjectIDFromHex(linkID)
	if err != nil {
		return models.Link{}, errors.New("invalid link ID format")
	}

	// Get link details first to verify it exists
	link, err := s.linkRepo.GetByID(ctx, objectID)
	if err != nil {
		return models.Link{}, err
	}

	// Check if link is expired
	if !link.ExpiresAt.IsZero() && link.ExpiresAt.Before(time.Now()) {
		return models.Link{}, errors.New("link has expired")
	}

	// Use a channel to handle visit creation asynchronously
	errChan := make(chan error, 1)
	go func() {
		visit := models.Visit{
			LinkID:    objectID,
			Timestamp: time.Now(),
			UserAgent: userAgent,
			IP:        ip,
			Referrer:  referrer,
		}
		err := s.visitRepo.Create(context.Background(), visit)
		errChan <- err
	}()

	// Increment clicks synchronously
	if err := s.linkRepo.IncrementClicks(ctx, objectID); err != nil {
		return models.Link{}, err
	}

	// Update link with incremented click count before returning
	link.Clicks++

	// Check if there was an error recording the visit
	select {
	case err := <-errChan:
		if err != nil {
			// Log the error but don't fail the request
			// Just continue and return the link with updated click count
		}
	case <-ctx.Done():
		// Context was canceled, but we still return the link
	}

	return link, nil
}

// GetVisitsForLink retrieves all visits for a link
func (s *VisitService) GetVisitsForLink(ctx context.Context, linkID string, page, pageSize int64) ([]models.Visit, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.visitRepo.GetVisitsByLinkID(ctx, linkID, pageSize, offset)
}