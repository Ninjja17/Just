package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/yourusername/10000hr/internal/models"
	"github.com/yourusername/10000hr/internal/repositories"
	"github.com/yourusername/10000hr/internal/services"
)

type AdminHandler struct {
	authSvc *services.AdminAuthService
	report  *repositories.AdminReportRepository
}

func NewAdminHandler(authSvc *services.AdminAuthService, report *repositories.AdminReportRepository) *AdminHandler {
	return &AdminHandler{authSvc: authSvc, report: report}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req models.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	admin, token, err := h.authSvc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.AdminAuthResponse{Admin: *admin, Token: token})
}

func (h *AdminHandler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id":    c.MustGet("admin_id"),
		"email": c.MustGet("admin_email"),
		"role":  c.MustGet("admin_role"),
	})
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	search := c.Query("search")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, total, err := h.report.ListUsers(c.Request.Context(), search, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users":  users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *AdminHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	d, err := h.report.GetUserDetail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if d == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.report.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *AdminHandler) Stats(c *gin.Context) {
	s, err := h.report.GlobalStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *AdminHandler) ExportUsersCSV(c *gin.Context) {
	rows, err := h.report.ExportUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `attachment; filename="users.csv"`)

	w := csv.NewWriter(c.Writer)
	_ = w.Write([]string{"id", "email", "name", "is_verified", "created_at", "skill_count", "session_count", "total_minutes"})
	for _, r := range rows {
		_ = w.Write([]string{
			r.ID.String(),
			r.Email,
			r.Name,
			fmt.Sprintf("%t", r.IsVerified),
			r.CreatedAt,
			strconv.Itoa(r.Skills),
			strconv.Itoa(r.Sessions),
			strconv.Itoa(r.TotalMin),
		})
	}
	w.Flush()
}

func (h *AdminHandler) ExportSessionsCSV(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10000"))
	rows, err := h.report.ExportSessions(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `attachment; filename="sessions.csv"`)

	w := csv.NewWriter(c.Writer)
	_ = w.Write([]string{"id", "user_email", "skill_name", "start_time", "end_time", "duration_minutes", "session_type"})
	for _, r := range rows {
		_ = w.Write([]string{
			r.ID.String(),
			r.UserEmail,
			r.SkillName,
			r.StartTime,
			r.EndTime,
			strconv.Itoa(r.Minutes),
			r.Type,
		})
	}
	w.Flush()
}
