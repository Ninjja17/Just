package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/yourusername/10000hr/internal/models"
)

type OTPRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewOTPRepository(db *sql.DB, redis *redis.Client) *OTPRepository {
	return &OTPRepository{db: db, redis: redis}
}

func (r *OTPRepository) Create(ctx context.Context, otp *models.OTP) error {
	query := `
		INSERT INTO otps (id, user_id, code, purpose, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, otp.ID, otp.UserID, otp.Code, otp.Purpose, otp.ExpiresAt)
	if err != nil {
		return err
	}

	// Cache in Redis for faster lookup
	key := fmt.Sprintf("otp:%s:%s", otp.Purpose, otp.Code)
	return r.redis.Set(ctx, key, otp.UserID.String(), time.Until(otp.ExpiresAt)).Err()
}

func (r *OTPRepository) Verify(ctx context.Context, email, code, purpose string) (*models.OTP, error) {
	// First, get user by email
	var userID uuid.UUID
	err := r.db.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", email).Scan(&userID)
	if err != nil {
		return nil, err
	}

	// Check if OTP exists and is valid
	otp := &models.OTP{}
	query := `
		SELECT id, user_id, code, purpose, expires_at, used_at
		FROM otps
		WHERE user_id = $1 AND code = $2 AND purpose = $3 AND used_at IS NULL AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`
	err = r.db.QueryRowContext(ctx, query, userID, code, purpose).Scan(
		&otp.ID, &otp.UserID, &otp.Code, &otp.Purpose, &otp.ExpiresAt, &otp.UsedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid or expired OTP")
	}
	if err != nil {
		return nil, err
	}

	return otp, nil
}

func (r *OTPRepository) MarkUsed(ctx context.Context, otpID uuid.UUID) error {
	query := `UPDATE otps SET used_at = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, otpID, time.Now())
	return err
}

func (r *OTPRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM otps WHERE expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
