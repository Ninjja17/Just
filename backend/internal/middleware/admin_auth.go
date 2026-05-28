package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AdminClaims struct {
	AdminID uuid.UUID `json:"admin_id"`
	Email   string    `json:"email"`
	Role    string    `json:"role"`
	jwt.RegisteredClaims
}

func adminSecret() []byte {
	s := os.Getenv("ADMIN_JWT_SECRET")
	if s == "" {
		s = os.Getenv("JWT_SECRET") + "-admin"
	}
	return []byte(s)
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		tok, err := jwt.ParseWithClaims(parts[1], &AdminClaims{}, func(t *jwt.Token) (interface{}, error) {
			return adminSecret(), nil
		})
		if err != nil || !tok.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}
		claims, ok := tok.Claims.(*AdminClaims)
		if !ok || claims.Issuer != "10000hr-admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid admin token"})
			return
		}
		c.Set("admin_id", claims.AdminID)
		c.Set("admin_email", claims.Email)
		c.Set("admin_role", claims.Role)
		c.Next()
	}
}

// RequireRole aborts the request unless the admin's role is one of the allowed roles.
func RequireRole(allowed ...string) gin.HandlerFunc {
	allowedSet := make(map[string]bool, len(allowed))
	for _, r := range allowed {
		allowedSet[r] = true
	}
	return func(c *gin.Context) {
		role, _ := c.Get("admin_role")
		s, _ := role.(string)
		if !allowedSet[s] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
