package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/10000hr/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, google_id, name, avatar_url, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.GoogleID, user.Name,
		user.AvatarURL, user.IsVerified, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, google_id, name, avatar_url, is_verified, created_at, updated_at
		FROM users WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.GoogleID, &user.Name,
		&user.AvatarURL, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, google_id, name, avatar_url, is_verified, created_at, updated_at
		FROM users WHERE email = $1
	`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.GoogleID, &user.Name,
		&user.AvatarURL, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, google_id, name, avatar_url, is_verified, created_at, updated_at
		FROM users WHERE google_id = $1
	`
	err := r.db.QueryRowContext(ctx, query, googleID).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.GoogleID, &user.Name,
		&user.AvatarURL, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET email = $2, password_hash = $3, name = $4, avatar_url = $5, is_verified = $6, updated_at = $7
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Name, user.AvatarURL, user.IsVerified, time.Now(),
	)
	return err
}

func (r *UserRepository) SetVerified(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET is_verified = true, updated_at = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID, time.Now())
	return err
}
