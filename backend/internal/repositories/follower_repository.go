package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/yourusername/10000hr/internal/models"
)

type FollowerRepository struct {
	db *sql.DB
}

func NewFollowerRepository(db *sql.DB) *FollowerRepository {
	return &FollowerRepository{db: db}
}

func (r *FollowerRepository) Create(ctx context.Context, followerID, followingID uuid.UUID) error {
	query := `INSERT INTO followers (follower_id, following_id, created_at) VALUES ($1, $2, NOW())`
	_, err := r.db.ExecContext(ctx, query, followerID, followingID)
	return err
}

func (r *FollowerRepository) Delete(ctx context.Context, followerID, followingID uuid.UUID) error {
	query := `DELETE FROM followers WHERE follower_id = $1 AND following_id = $2`
	_, err := r.db.ExecContext(ctx, query, followerID, followingID)
	return err
}

func (r *FollowerRepository) GetFollowers(ctx context.Context, userID uuid.UUID) ([]*models.User, error) {
	query := `
		SELECT u.id, u.email, u.name, u.avatar_url, u.created_at
		FROM users u
		INNER JOIN followers f ON u.id = f.follower_id
		WHERE f.following_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.AvatarURL, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *FollowerRepository) GetFollowing(ctx context.Context, userID uuid.UUID) ([]*models.User, error) {
	query := `
		SELECT u.id, u.email, u.name, u.avatar_url, u.created_at
		FROM users u
		INNER JOIN followers f ON u.id = f.following_id
		WHERE f.follower_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.AvatarURL, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *FollowerRepository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM followers WHERE follower_id = $1 AND following_id = $2)`
	err := r.db.QueryRowContext(ctx, query, followerID, followingID).Scan(&exists)
	return exists, err
}
