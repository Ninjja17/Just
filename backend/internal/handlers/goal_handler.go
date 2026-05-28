package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/10000hr/internal/middleware"
	"github.com/yourusername/10000hr/internal/models"
	"github.com/yourusername/10000hr/internal/services"
)

type GoalHandler struct {
	goalService *services.GoalService
}

func NewGoalHandler(goalService *services.GoalService) *GoalHandler {
	return &GoalHandler{goalService: goalService}
}

func (h *GoalHandler) CreateGoal(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req models.CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal, err := h.goalService.CreateGoal(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

func (h *GoalHandler) GetGoals(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	goals, err := h.goalService.GetGoals(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goals)
}

func (h *GoalHandler) UpdateGoal(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	var req models.CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal, err := h.goalService.UpdateGoal(c.Request.Context(), goalID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)
}

func (h *GoalHandler) DeleteGoal(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	if err := h.goalService.DeleteGoal(c.Request.Context(), goalID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goal deleted successfully"})
}
