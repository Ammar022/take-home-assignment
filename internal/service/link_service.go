package service

import (
	"context"
	"errors"
	"take-home-assignment/internal/models"
	"take-home-assignment/internal/repo"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LinkService handles link business logic
type LinkService struct {
	repo *repo.LinkRepository
}

// NewLinkService creates a new link service
func NewLinkService(repo *repo.LinkRepository) *LinkService {
	return &LinkService{
		repo: repo,
	}
}

// CreateLink creates a new link
func (s *LinkService) CreateLink(ctx context.Context, dto models.LinkCreateDTO) (models.Link, error) {
	link := models.Link{
		Title:     dto.Title,
		URL:       dto.URL,
		CreatedAt: time.Now(),
		ExpiresAt: dto.ExpiresAt,
		Clicks:    0,
		UserID:    dto.UserID,
	}

	return s.repo.Create(ctx, link)
}

// GetLinkByID retrieves a link by ID
func (s *LinkService) GetLinkByID(ctx context.Context, id string) (models.Link, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Link{}, errors.New("invalid link ID format")
	}

	return s.repo.GetByID(ctx, objectID)
}

// GetAllLinks retrieves all links for a user
func (s *LinkService) GetAllLinks(ctx context.Context, userID string, page, pageSize int64) ([]models.Link, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.repo.GetAll(ctx, userID, pageSize, offset)
}

// UpdateLink updates an existing link
func (s *LinkService) UpdateLink(ctx context.Context, id string, dto models.LinkUpdateDTO) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid link ID format")
	}

	return s.repo.Update(ctx, objectID, dto)
}

// DeleteLink deletes a link
func (s *LinkService) DeleteLink(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid link ID format")
	}

	return s.repo.Delete(ctx, objectID)
}
