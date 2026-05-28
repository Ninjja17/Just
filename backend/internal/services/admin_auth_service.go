package services

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/yourusername/10000hr/internal/models"
	"github.com/yourusername/10000hr/internal/repositories"
)

type AdminAuthService struct {
	repo *repositories.AdminRepository
}

func NewAdminAuthService(repo *repositories.AdminRepository) *AdminAuthService {
	return &AdminAuthService{repo: repo}
}

// AdminClaims is the JWT payload for admin tokens. Distinct issuer keeps
// admin tokens from being accepted by the user-facing API.
type AdminClaims struct {
	AdminID uuid.UUID `json:"admin_id"`
	Email   string    `json:"email"`
	Role    string    `json:"role"`
	jwt.RegisteredClaims
}

func (s *AdminAuthService) Login(ctx context.Context, email, password string) (*models.AdminUser, string, error) {
	admin, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if admin == nil {
		return nil, "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.issueToken(admin)
	if err != nil {
		return nil, "", err
	}
	return admin, token, nil
}

func (s *AdminAuthService) issueToken(admin *models.AdminUser) (string, error) {
	secret := os.Getenv("ADMIN_JWT_SECRET")
	if secret == "" {
		secret = os.Getenv("JWT_SECRET") + "-admin"
	}
	claims := AdminClaims{
		AdminID: admin.ID,
		Email:   admin.Email,
		Role:    admin.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "10000hr-admin",
			Subject:   admin.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(secret))
}

// EnsureBootstrapAdmin creates a default admin account from env vars if no
// admin exists. Useful for first run.
func (s *AdminAuthService) EnsureBootstrapAdmin(ctx context.Context) error {
	if err := s.repo.EnsureSchema(ctx); err != nil {
		return err
	}
	count, err := s.repo.Count(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	email := os.Getenv("ADMIN_BOOTSTRAP_EMAIL")
	password := os.Getenv("ADMIN_BOOTSTRAP_PASSWORD")
	name := os.Getenv("ADMIN_BOOTSTRAP_NAME")
	if email == "" {
		email = "admin@local.dev"
	}
	if password == "" {
		password = "changeme123"
	}
	if name == "" {
		name = "Bootstrap Admin"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.Create(ctx, &models.AdminUser{
		Email:        email,
		PasswordHash: string(hash),
		Name:         name,
		Role:         "superadmin",
	})
}
