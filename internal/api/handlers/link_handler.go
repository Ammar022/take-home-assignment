// internal/api/handlers/link_handler.go
package handlers

import (
	"net/http"
	"strconv"
	"take-home-assignment/internal/models"
	"take-home-assignment/internal/service"

	"github.com/gin-gonic/gin"
)

// LinkHandler handles link-related HTTP requests
type LinkHandler struct {
	linkService *service.LinkService
}

// NewLinkHandler creates a new link handler
func NewLinkHandler(linkService *service.LinkService) *LinkHandler {
	return &LinkHandler{
		linkService: linkService,
	}
}

// internal/api/handlers/link_handler.go (continued)
// Create handles creating a new link
func (h *LinkHandler) Create(c *gin.Context) {
	var dto models.LinkCreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if exists {
		dto.UserID = userID.(string)
	}

	link, err := h.linkService.CreateLink(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, link)
}

// GetByID handles retrieving a link by ID
func (h *LinkHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	link, err := h.linkService.GetLinkByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	c.JSON(http.StatusOK, link)
}

// GetAll handles retrieving all links for a user
func (h *LinkHandler) GetAll(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("pageSize", "10"), 10, 64)

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	links, err := h.linkService.GetAllLinks(c.Request.Context(), userID.(string), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, links)
}

// Update handles updating an existing link
func (h *LinkHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var dto models.LinkUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.linkService.UpdateLink(c.Request.Context(), id, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete handles deleting a link
func (h *LinkHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.linkService.DeleteLink(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}


