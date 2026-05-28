package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/10000hr/internal/middleware"
	"github.com/yourusername/10000hr/internal/services"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

func (h *AnalyticsHandler) GetOverview(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	overview, err := h.analyticsService.GetOverview(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, overview)
}

func (h *AnalyticsHandler) GetTrends(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	trends, err := h.analyticsService.GetTrends(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trends)
}

func (h *AnalyticsHandler) GetPredictions(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	predictions, err := h.analyticsService.GetPredictions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, predictions)
}
