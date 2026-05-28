package services

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/yourusername/10000hr/internal/middleware"
	"github.com/yourusername/10000hr/internal/models"
	"github.com/yourusername/10000hr/internal/repositories"
	"github.com/yourusername/10000hr/pkg/email"
)

type AuthService struct {
	userRepo *repositories.UserRepository
	otpRepo  *repositories.OTPRepository
}

func NewAuthService(userRepo *repositories.UserRepository, otpRepo *repositories.OTPRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		otpRepo:  otpRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	// Check if user already exists
	existing, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	passwordHashStr := string(hashedPassword)
	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: &passwordHashStr,
		Name:         req.Name,
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Send OTP for verification
	if err := s.SendOTP(ctx, user.ID, req.Email, "verification"); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if user.PasswordHash == nil {
		return nil, fmt.Errorf("please use Google login")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.IsVerified {
		return nil, fmt.Errorf("please verify your email first")
	}

	// Generate JWT token
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:  *user,
		Token: token,
	}, nil
}

func (s *AuthService) VerifyOTP(ctx context.Context, email, code string) (*models.AuthResponse, error) {
	otp, err := s.otpRepo.Verify(ctx, email, code, "verification")
	if err != nil {
		return nil, err
	}

	// Mark OTP as used
	if err := s.otpRepo.MarkUsed(ctx, otp.ID); err != nil {
		return nil, err
	}

	// Set user as verified
	if err := s.userRepo.SetVerified(ctx, otp.UserID); err != nil {
		return nil, err
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, otp.UserID)
	if err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:  *user,
		Token: token,
	}, nil
}

func (s *AuthService) SendOTP(ctx context.Context, userID uuid.UUID, userEmail, purpose string) error {
	// Generate 6-digit OTP
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	otp := &models.OTP{
		ID:        uuid.New(),
		UserID:    userID,
		Code:      code,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := s.otpRepo.Create(ctx, otp); err != nil {
		return err
	}

	// Send email
	emailService := email.NewEmailService()
	subject := "Your Verification Code"
	body := fmt.Sprintf("Your OTP code is: %s\n\nThis code will expire in 10 minutes.", code)

	return emailService.SendEmail(userEmail, subject, body)
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	expiryStr := os.Getenv("JWT_EXPIRY")
	if expiryStr == "" {
		expiryStr = "24h"
	}
	expiry, _ := time.ParseDuration(expiryStr)

	claims := &middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *AuthService) GoogleLogin(ctx context.Context, googleID, email, name, avatarURL string) (*models.AuthResponse, error) {
	// Check if user exists with Google ID
	user, err := s.userRepo.GetByGoogleID(ctx, googleID)
	if err != nil {
		return nil, err
	}

	// If not found, check by email
	if user == nil {
		user, err = s.userRepo.GetByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
	}

	// Create new user if not found
	if user == nil {
		user = &models.User{
			ID:         uuid.New(),
			Email:      email,
			GoogleID:   &googleID,
			Name:       name,
			AvatarURL:  &avatarURL,
			IsVerified: true, // Google accounts are pre-verified
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	// Generate token
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:  *user,
		Token: token,
	}, nil
}
