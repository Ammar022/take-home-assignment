package handlers

import (
	"net/http"
	"strconv"
	"take-home-assignment/internal/service"

	"github.com/gin-gonic/gin"
)

// VisitHandler handles visit-related HTTP requests
type VisitHandler struct {
	visitService *service.VisitService
}

// NewVisitHandler creates a new visit handler
func NewVisitHandler(visitService *service.VisitService) *VisitHandler {
	return &VisitHandler{
		visitService: visitService,
	}
}

// RecordVisit handles recording a visit to a link
func (h *VisitHandler) RecordVisit(c *gin.Context) {
	id := c.Param("id")
	userAgent := c.Request.UserAgent()
	ip := c.ClientIP()
	referrer := c.Request.Referer()

	link, err := h.visitService.RecordVisit(c.Request.Context(), id, userAgent, ip, referrer)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found or expired"})
		return
	}

	// Redirect to the actual URL
	c.Redirect(http.StatusFound, link.URL)
}

// GetVisitsForLink handles retrieving all visits for a link
func (h *VisitHandler) GetVisitsForLink(c *gin.Context) {
	id := c.Param("id")
	
	// Get pagination parameters
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("pageSize", "10"), 10, 64)

	visits, err := h.visitService.GetVisitsForLink(c.Request.Context(), id, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, visits)
}
