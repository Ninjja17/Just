package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/10000hr/internal/models"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// EnsureSchema creates the admin_users table if missing so the admin app works
// without a separate migration step in local dev.
func (r *AdminRepository) EnsureSchema(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS admin_users (
			id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			email         TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			name          TEXT NOT NULL,
			role          TEXT NOT NULL DEFAULT 'admin',
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_admin_users_email ON admin_users(email);
	`)
	return err
}

func (r *AdminRepository) Count(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM admin_users`).Scan(&n)
	return n, err
}

func (r *AdminRepository) Create(ctx context.Context, a *models.AdminUser) error {
	now := time.Now()
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	a.CreatedAt = now
	a.UpdatedAt = now
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO admin_users (id, email, password_hash, name, role, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		a.ID, a.Email, a.PasswordHash, a.Name, a.Role, a.CreatedAt, a.UpdatedAt,
	)
	return err
}

func (r *AdminRepository) GetByEmail(ctx context.Context, email string) (*models.AdminUser, error) {
	a := &models.AdminUser{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, name, role, created_at, updated_at
		 FROM admin_users WHERE email = $1`, email,
	).Scan(&a.ID, &a.Email, &a.PasswordHash, &a.Name, &a.Role, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return a, err
}

func (r *AdminRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.AdminUser, error) {
	a := &models.AdminUser{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, name, role, created_at, updated_at
		 FROM admin_users WHERE id = $1`, id,
	).Scan(&a.ID, &a.Email, &a.PasswordHash, &a.Name, &a.Role, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return a, err
}
