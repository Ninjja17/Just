package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	PasswordHash *string    `json:"-"`
	GoogleID     *string    `json:"google_id,omitempty"`
	Name         string     `json:"name"`
	AvatarURL    *string    `json:"avatar_url,omitempty"`
	IsVerified   bool       `json:"is_verified"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Skill struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Name        string     `json:"name"`
	Category    string     `json:"category"`
	Color       string     `json:"color"`
	Icon        string     `json:"icon"`
	TargetHours int        `json:"target_hours"`
	IsArchived  bool       `json:"is_archived"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Session struct {
	ID             uuid.UUID  `json:"id"`
	SkillID        uuid.UUID  `json:"skill_id"`
	UserID         uuid.UUID  `json:"user_id"`
	StartTime      time.Time  `json:"start_time"`
	EndTime        *time.Time `json:"end_time,omitempty"`
	DurationMinutes int       `json:"duration_minutes"`
	Notes          *string    `json:"notes,omitempty"`
	SessionType    string     `json:"session_type"`
	QualityRating  *int       `json:"quality_rating,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type Goal struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	SkillID     *uuid.UUID `json:"skill_id,omitempty"`
	Type        string     `json:"type"` // daily, weekly, milestone
	TargetHours float64    `json:"target_hours"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	IsCompleted bool       `json:"is_completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type Follower struct {
	FollowerID  uuid.UUID `json:"follower_id"`
	FollowingID uuid.UUID `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Notification struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type OTP struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Code      string     `json:"code"`
	Purpose   string     `json:"purpose"` // verification, reset
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
}

// Request/Response DTOs
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

type CreateSkillRequest struct {
	Name        string `json:"name" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Color       string `json:"color" binding:"required"`
	Icon        string `json:"icon" binding:"required"`
	TargetHours int    `json:"target_hours"`
}

type CreateSessionRequest struct {
	SkillID         uuid.UUID  `json:"skill_id" binding:"required"`
	StartTime       time.Time  `json:"start_time" binding:"required"`
	EndTime         *time.Time `json:"end_time"`
	DurationMinutes int        `json:"duration_minutes"`
	Notes           *string    `json:"notes"`
	SessionType     string     `json:"session_type"`
	QualityRating   *int       `json:"quality_rating"`
}

type CreateGoalRequest struct {
	SkillID     *uuid.UUID `json:"skill_id"`
	Type        string     `json:"type" binding:"required"`
	TargetHours float64    `json:"target_hours" binding:"required"`
	Deadline    *time.Time `json:"deadline"`
}

type AuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
