package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/10000hr/internal/middleware"
	"github.com/yourusername/10000hr/internal/services"
)

type SocialHandler struct {
	socialService *services.SocialService
}

func NewSocialHandler(socialService *services.SocialService) *SocialHandler {
	return &SocialHandler{socialService: socialService}
}

func (h *SocialHandler) GetProfile(c *gin.Context) {
	profileUserID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.socialService.GetProfile(c.Request.Context(), profileUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *SocialHandler) Follow(c *gin.Context) {
	followerID, _ := middleware.GetUserID(c)
	followingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.socialService.Follow(c.Request.Context(), followerID, followingID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Followed successfully"})
}

func (h *SocialHandler) Unfollow(c *gin.Context) {
	followerID, _ := middleware.GetUserID(c)
	followingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.socialService.Unfollow(c.Request.Context(), followerID, followingID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
}

func (h *SocialHandler) GetFollowers(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	followers, err := h.socialService.GetFollowers(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, followers)
}

func (h *SocialHandler) GetFollowing(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	following, err := h.socialService.GetFollowing(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, following)
}

func (h *SocialHandler) GetLeaderboard(c *gin.Context) {
	// Placeholder for leaderboard logic
	c.JSON(http.StatusOK, gin.H{"message": "Leaderboard coming soon"})
}

func (h *SocialHandler) GetFeed(c *gin.Context) {
	// Placeholder for activity feed logic
	c.JSON(http.StatusOK, gin.H{"message": "Activity feed coming soon"})
}
