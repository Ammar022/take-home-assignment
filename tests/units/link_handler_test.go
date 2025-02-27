package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"take-home-assignment/internal/api/handlers"
	"take-home-assignment/internal/models"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock service
type MockLinkService struct {
	mock.Mock
}

func (m *MockLinkService) CreateLink(ctx context.Context, dto models.LinkCreateDTO) (models.Link, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(models.Link), args.Error(1)
}

func (m *MockLinkService) GetLinkByID(ctx context.Context, id string) (models.Link, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Link), args.Error(1)
}

func (m *MockLinkService) GetAllLinks(ctx context.Context, userID string, page, pageSize int64) ([]models.Link, error) {
	args := m.Called(ctx, userID, page, pageSize)
	return args.Get(0).([]models.Link), args.Error(1)
}

func (m *MockLinkService) UpdateLink(ctx context.Context, id string, dto models.LinkUpdateDTO) error {
	args := m.Called(ctx, id, dto)
	return args.Error(0)
}

func (m *MockLinkService) DeleteLink(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestCreateLinkHandler(t *testing.T) {
	// Create mock service
	mockService := new(MockLinkService)
	
	// Create test data
	createDTO := models.LinkCreateDTO{
		Title:     "Test Link",
		URL:       "https://example.com",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	
	expectedLink := models.Link{
		ID:        primitive.NewObjectID(),
		Title:     createDTO.Title,
		URL:       createDTO.URL,
		CreatedAt: time.Now(),
		ExpiresAt: createDTO.ExpiresAt,
		Clicks:    0,
		UserID:    "user123",
	}
	
	// Set up mock expectations
	mockService.On("CreateLink", mock.Anything, mock.MatchedBy(func(dto models.LinkCreateDTO) bool {
		return dto.Title == createDTO.Title && dto.URL == createDTO.URL
	})).Return(expectedLink, nil)
	
	// Create handler with mock service
	handler := handlers.NewLinkHandler(mockService)
	
	// Setup router
	router := setupRouter()
	router.POST("/api/links", func(c *gin.Context) {
		// Mock authentication middleware
		c.Set("userId", "user123")
		handler.Create(c)
	})
	
	// Create request
	body, _ := json.Marshal(createDTO)
	req, _ := http.NewRequest("POST", "/api/links", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	
	// Assert results
	assert.Equal(t, http.StatusCreated, recorder.Code)
	
	var response models.Link
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedLink.ID, response.ID)
	assert.Equal(t, expectedLink.Title, response.Title)
	assert.Equal(t, expectedLink.URL, response.URL)
	
	// Verify that mock expectations were met
	mockService.AssertExpectations(t)
}