// tests/unit/link_service_test.go
package unit

import (
	"context"
	"take-home-assignment/internal/models"
	"take-home-assignment/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock repository
type MockLinkRepository struct {
	mock.Mock
}

func (m *MockLinkRepository) Create(ctx context.Context, link models.Link) (models.Link, error) {
	args := m.Called(ctx, link)
	return args.Get(0).(models.Link), args.Error(1)
}

func (m *MockLinkRepository) GetByID(ctx context.Context, id primitive.ObjectID) (models.Link, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Link), args.Error(1)
}

func (m *MockLinkRepository) GetAll(ctx context.Context, userID string, limit, offset int64) ([]models.Link, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]models.Link), args.Error(1)
}

func (m *MockLinkRepository) Update(ctx context.Context, id primitive.ObjectID, link models.LinkUpdateDTO) error {
	args := m.Called(ctx, id, link)
	return args.Error(0)
}

func (m *MockLinkRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockLinkRepository) DeleteExpired(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLinkRepository) IncrementClicks(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateLink(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockLinkRepository)
	
	// Create test data
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	
	createDTO := models.LinkCreateDTO{
		Title:     "Test Link",
		URL:       "https://example.com",
		ExpiresAt: expiresAt,
		UserID:    "user123",
	}
	
	expectedLink := models.Link{
		ID:        primitive.NewObjectID(),
		Title:     createDTO.Title,
		URL:       createDTO.URL,
		CreatedAt: now,
		ExpiresAt: createDTO.ExpiresAt,
		Clicks:    0,
		UserID:    createDTO.UserID,
	}
	
	// Set up mock expectations
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("models.Link")).Return(expectedLink, nil)
	
	// Create service with mock repository
	service := service.NewLinkService(mockRepo)
	
	// Test CreateLink method
	result, err := service.CreateLink(context.Background(), createDTO)
	
	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, expectedLink.ID, result.ID)
	assert.Equal(t, createDTO.Title, result.Title)
	assert.Equal(t, createDTO.URL, result.URL)
	assert.Equal(t, createDTO.UserID, result.UserID)
	
	// Verify that mock expectations were met
	mockRepo.AssertExpectations(t)
}

func TestGetLinkByID(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockLinkRepository)
	
	// Create test data
	id := primitive.NewObjectID()
	expectedLink := models.Link{
		ID:        id,
		Title:     "Test Link",
		URL:       "https://example.com",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Clicks:    5,
		UserID:    "user123",
	}
	
	// Set up mock expectations
	mockRepo.On("GetByID", mock.Anything, id).Return(expectedLink, nil)
	
	// Create service with mock repository
	service := service.NewLinkService(mockRepo)
	
	// Test GetLinkByID method
	result, err := service.GetLinkByID(context.Background(), id.Hex())
	
	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, expectedLink.ID, result.ID)
	assert.Equal(t, expectedLink.Title, result.Title)
	assert.Equal(t, expectedLink.URL, result.URL)
	assert.Equal(t, expectedLink.Clicks, result.Clicks)
	
	// Verify that mock expectations were met
	mockRepo.AssertExpectations(t)
}
